package idgo

import (
	"errors"
	"sync"
)

// IDGenerator generate a id.
type IDGenerator struct {
	mutex *sync.Mutex

	allocatedIDStore AllocatedIDStore

	// nextTryID is the number to try next when allocate id.
	nextTryID int

	// allocatedIDCount is amount of currently allocated id.
	allocatedIDCount int

	// maxSize is the maximum value of allocatable id.
	maxSize int
}

// NewIDGenerator is IDGenerator constructed.
func NewIDGenerator(store AllocatedIDStore) (*IDGenerator, error) {
	if store.getMaxSize() <= 0 {
		return nil, errors.New("argument is negative number")
	}

	return &IDGenerator{
		mutex:            &sync.Mutex{},
		allocatedIDStore: store,
		maxSize:          store.getMaxSize(),
	}, nil
}

// Generate a new id.
func (g *IDGenerator) Generate() (int, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for {
		if g.allocatedIDCount >= g.maxSize {
			return 0, errors.New("id is exhausted")
		}

		// Set nextTryID to 0 when nextTryID is greater than maxSize.
		if g.nextTryID >= g.maxSize {
			g.nextTryID = 0
		}

		// Allocate and return nextTryID when nextTryID is not yet allocated.
		if !g.allocatedIDStore.isAllocated(g.nextTryID) {
			id := g.nextTryID
			g.allocatedIDStore.allocate(id)
			g.allocatedIDCount++
			g.nextTryID++
			return id, nil
		}

		// When nextTryID was assigned, increment nextTryID and try allocate id again.
		g.nextTryID++
	}
}

// Allocate a specified id.
func (g *IDGenerator) Allocate(id int) error {
	if id > g.maxSize {
		return errors.New("id exceeds the maximum value")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.allocatedIDStore.isAllocated(id) {
		return errors.New("id is already allocated")
	}

	g.allocatedIDStore.allocate(id)
	g.allocatedIDCount++
	return nil
}

// Free a allocated id.
func (g *IDGenerator) Free(id int) {
	if id > g.maxSize {
		return
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.allocatedIDStore.isAllocated(id) {
		g.allocatedIDStore.free(id)
		g.allocatedIDCount--
	}
}

// FreeAll free all allocated id.
func (g *IDGenerator) FreeAll() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.allocatedIDStore.freeAll()
	g.nextTryID = 0
	g.allocatedIDCount = 0
}

// IsAllocated check if specified id is allocated.
func (g *IDGenerator) IsAllocated(id int) bool {
	if id > g.maxSize {
		return false
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.allocatedIDStore.isAllocated(id)
}

// GetAllocatedIDCount is getter for allocatedIDCount.
func (g *IDGenerator) GetAllocatedIDCount() int {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.allocatedIDCount
}
