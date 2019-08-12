package idgo

import "github.com/minami14/go-bitarray"

// AllocatedIDStore stores allocated id.
type AllocatedIDStore interface {
	isAllocated(int) (bool, error)
	allocate(int) error
	free(int) error
	freeAll() error
	getMaxSize() int
}

// LocalStore stores allocated id to byte slice.
type LocalStore struct {
	maxSize     int
	allocatedID *bitarray.BitArray
}

// NewLocalStore is LocalStore constructed.
func NewLocalStore(maxSize int) (*LocalStore, error) {
	allocatedID, err := bitarray.NewBitArray(maxSize)
	if err != nil {
		return nil, err
	}

	return &LocalStore{
		maxSize:     maxSize,
		allocatedID: allocatedID,
	}, nil
}

func (l *LocalStore) isAllocated(id int) (bool, error) {
	return l.allocatedID.Get(id)
}

func (l *LocalStore) allocate(id int) error {
	return l.allocatedID.Set(id)
}

func (l *LocalStore) free(id int) error {
	return l.allocatedID.Clear(id)
}

func (l *LocalStore) freeAll() error {
	l.allocatedID.Reset()
	return nil
}

func (l *LocalStore) getMaxSize() int {
	return l.maxSize
}
