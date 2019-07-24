package idgo

import (
	"errors"
	"sync"
)

// IDGenerator generate a id.
type IDGenerator struct {
	mutex *sync.Mutex

	// allocatedID is slice to judgement if id is allocated.
	allocatedID []byte

	// nextTryID is the number to try next when allocate id.
	nextTryID int

	// allocatedIDCount is amount of currently allocated id.
	allocatedIDCount int

	// maxSize is the maximum value of allocatable id.
	maxSize int

	// minSize is the minimum value of allocatable id.
	minSize int

	size int
}

const bits = 8

// NewIDGenerator is IDGenerator constructed.
func NewIDGenerator(minSize, maxSize int) (*IDGenerator, error) {
	if maxSize <= 0 {
		return nil, errors.New("maxSize can only be 1 or more")
	}

	if minSize < 0 {
		return nil, errors.New("minSize can only be 0 or more")
	}

	if maxSize <= minSize {
		return nil, errors.New("maxSize can only be minSize or more")
	}

	size := maxSize - minSize
	return &IDGenerator{
		mutex:       &sync.Mutex{},
		allocatedID: make([]byte, size/bits+1),
		maxSize:     maxSize,
		minSize:     minSize,
		size:        size,
	}, nil
}

// Generate a new id.
func (g *IDGenerator) Generate() (int, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for {
		if g.allocatedIDCount >= g.size {
			return 0, errors.New("id is exhausted")
		}

		// Set nextTryID to 0 when nextTryID is greater than maxSize.
		if g.nextTryID >= g.maxSize {
			g.nextTryID = g.minSize
		}

		// Allocate and return nextTryID when nextTryID is not yet allocated.
		if !g.isAllocated(g.nextTryID) {
			id := g.nextTryID
			g.allocate(id)
			g.nextTryID++
			return id, nil
		}

		// When nextTryID was assigned, increment nextTryID and try allocate id again.
		g.nextTryID++
	}
}

// Allocate a specified id.
func (g *IDGenerator) Allocate(id int) error {
	if id > g.maxSize || id < g.minSize {
		return errors.New("id is out of range")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.isAllocated(id) {
		return errors.New("id is already allocated")
	}

	g.allocate(id)
	return nil
}

// Free a allocated id.
func (g *IDGenerator) Free(id int) {
	if id > g.maxSize || id < g.minSize {
		return
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.isAllocated(id) {
		g.free(id)
	}
}

// FreeAll free all allocated id.
func (g *IDGenerator) FreeAll() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.allocatedID = make([]byte, g.size)
	g.nextTryID = 0
	g.allocatedIDCount = 0
}

// IsAllocated check if specified id is allocated.
func (g *IDGenerator) IsAllocated(id int) bool {
	if id > g.maxSize || id < g.minSize {
		return false
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.isAllocated(id)
}

// GetAllocatedIDCount is getter for allocatedIDCount.
func (g *IDGenerator) GetAllocatedIDCount() int {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.allocatedIDCount
}

func (g *IDGenerator) isAllocated(id int) bool {
	i := id - g.minSize
	index := i / bits
	b := g.allocatedID[index]
	shift := byte(i % bits)
	mask := byte(1 << shift)
	flag := b & mask
	return flag != 0
}

func (g *IDGenerator) allocate(id int) {
	i := id - g.minSize
	index := i / bits
	b := g.allocatedID[index]
	shift := byte(i % bits)
	mask := byte(1 << shift)
	flag := b | mask
	g.allocatedID[index] = flag
	g.allocatedIDCount++
}

func (g *IDGenerator) free(id int) {
	i := id - g.minSize
	index := i / bits
	b := g.allocatedID[index]
	shift := byte(i % bits)
	mask := byte(1 << shift)
	flag := b & ^mask
	g.allocatedID[index] = flag
	g.allocatedIDCount--
}
