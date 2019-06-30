package idgo

import (
	"errors"
	"sync"
)

// Generator generate a id.
type Generator struct {
	mutex *sync.Mutex

	// allocatedID is slice to judgement if id is allocated.
	allocatedID []uint64

	// nextTryID is the number to try next when allocate id.
	nextTryID int

	// allocatedIDCount is amount of currently allocated id.
	allocatedIDCount int

	// maxSize is the maximum value of allocatable id.
	maxSize int
}

const bits = 64

// NewGenerator is Generator constructed.
func NewGenerator(maxSize int) (*Generator, error) {
	if maxSize <= 0 {
		return nil, errors.New("argument is negative number")
	}

	return &Generator{
		mutex:       &sync.Mutex{},
		allocatedID: make([]uint64, maxSize/bits+1),
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
			g.allocatedIDCount++
			return id, nil
		}

		// When nextTryID was assigned, increment nextTryID and try allocate id again.
		g.nextTryID++
	}
}

// Free a used id
func (g *Generator) Free(id int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if g.isAllocated(id) {
		g.free(id)
		g.allocatedIDCount--
	}
}

func (g *Generator) isAllocated(id int) bool {
	index := id / bits
	u64 := g.allocatedID[index]
	shift := uint64(id % bits)
	mask := uint64(1 << shift)
	flag := u64 & mask
	return flag != 0
}

func (g *Generator) allocate(id int) {
	index := id / bits
	u64 := g.allocatedID[index]
	shift := uint64(id % bits)
	mask := uint64(1 << shift)
	flag := u64 | mask
	g.allocatedID[index] = flag
}

func (g *Generator) free(id int) {
	index := id / bits
	u64 := g.allocatedID[index]
	shift := uint64(id % bits)
	mask := uint64(1 << shift)
	flag := u64 & ^mask
	g.allocatedID[index] = flag
}
