package idgo

import (
	"errors"
	"sync"
)

// Generator generate a id.
type Generator struct {
	mutex *sync.Mutex

	// allocatedID is slice to judgement if id is allocated.
	allocatedID []byte

	// nextTryID is the number to try next when allocate id.
	nextTryID int

	// allocatedIDCount is amount of currently allocated id.
	allocatedIDCount int

	// maxSize is the maximum value of allocatable id.
	maxSize int
}

const bits = 8

// NewGenerator is Generator constructed.
func NewGenerator(maxSize int) (*Generator, error) {
	if maxSize <= 0 {
		return nil, errors.New("argument is negative number")
	}

	return &Generator{
		mutex:       &sync.Mutex{},
		allocatedID: make([]byte, maxSize/bits+1),
		maxSize:     maxSize,
	}, nil
}

// Generate a new id
func (g *Generator) Generate() (int, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for {
		if g.allocatedIDCount >= g.maxSize {
			return 0, errors.New("id is exhausted")
		}

		// Set nextTryID to 0 when nextTryID is greater than maxSize
		if g.nextTryID >= g.maxSize {
			g.nextTryID = 0
		}

		// Allocate and return nextTryID when nextTryID is not yet allocated
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

// Allocate a specified id
func (g *Generator) Allocate(id int) error {
	if id > g.maxSize {
		return errors.New("id exceeds the maximum value")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.isAllocated(id) {
		return errors.New("id is already allocated")
	}

	g.allocate(id)
	return nil
}

// Free a allocated id
func (g *Generator) Free(id int) {
	if id > g.maxSize {
		return
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.isAllocated(id) {
		g.free(id)
	}
}

// FreeAll free all allocated id
func (g *Generator) FreeAll() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.allocatedID = make([]byte, g.maxSize)
	g.nextTryID = 0
	g.allocatedIDCount = 0
}

// IsAllocated check if specified id is allocated
func (g *Generator) IsAllocated(id int) bool {
	if id > g.maxSize {
		return false
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.isAllocated(id)
}

func (g *Generator) isAllocated(id int) bool {
	index := id / bits
	b := g.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b & mask
	return flag != 0
}

func (g *Generator) allocate(id int) {
	index := id / bits
	b := g.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b | mask
	g.allocatedID[index] = flag
	g.allocatedIDCount++
}

func (g *Generator) free(id int) {
	index := id / bits
	b := g.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b & ^mask
	g.allocatedID[index] = flag
	g.allocatedIDCount--
}
