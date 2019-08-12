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

	// maxSize is the maximum value of allocatable id.
	maxSize int
}

// NewIDGenerator is IDGenerator constructed.
func NewIDGenerator(store AllocatedIDStore) (*IDGenerator, error) {
	if store.getMaxSize() <= 0 {
		return nil, errors.New("maxsize is 0 or less")
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
		if g.allocatedIDStore.getAllocatedIDCount() >= g.maxSize {
			return 0, errors.New("id is exhausted")
		}

		// Set nextTryID to 0 when nextTryID is greater than maxSize.
		if g.nextTryID >= g.maxSize {
			g.nextTryID = 0
		}

		// Allocate and return nextTryID when nextTryID is not yet allocated.
		isAlloc, err := g.allocatedIDStore.isAllocated(g.nextTryID)
		if err != nil {
			return 0, err
		}

		if !isAlloc {
			id := g.nextTryID
			if err := g.allocatedIDStore.allocate(id); err != nil {
				return 0, err
			}
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

	isAlloc, err := g.allocatedIDStore.isAllocated(id)
	if err != nil {
		return err
	}

	if isAlloc {
		return errors.New("id is already allocated")
	}

	if err := g.allocatedIDStore.allocate(id); err != nil {
		return err
	}
	return nil
}

// Free a allocated id.
func (g *IDGenerator) Free(id int) error {
	if id > g.maxSize {
		return errors.New("greater than max size")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	isAlloc, err := g.allocatedIDStore.isAllocated(id)
	if err != nil {
		return err
	}

	if isAlloc {
		if err := g.allocatedIDStore.free(id); err != nil {
			return err
		}
	}

	return nil
}

// FreeAll free all allocated id.
func (g *IDGenerator) FreeAll() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if err := g.allocatedIDStore.freeAll(); err != nil {
		return err
	}
	g.nextTryID = 0

	return nil
}

// IsAllocated check if specified id is allocated.
func (g *IDGenerator) IsAllocated(id int) (bool, error) {
	if id > g.maxSize {
		return false, errors.New("greater than max size")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.allocatedIDStore.isAllocated(id)
}

// GetAllocatedIDCount is getter for allocatedIDCount.
func (g *IDGenerator) GetAllocatedIDCount() int {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.allocatedIDStore.getAllocatedIDCount()
}
