package idgo

// AllocatedIDStore stores allocated id.
type AllocatedIDStore interface {
	initialize() error
	isAllocated(int) (bool, error)
	allocate(int) error
	free(int) error
	freeAll() error
	getMaxSize() int
}

// LocalStore stores allocated id to byte slice.
type LocalStore struct {
	maxSize     int
	allocatedID []byte
}

const bits = 8

// NewLocalStore is LocalStore constructed.
func NewLocalStore(maxSize int) *LocalStore {
	return &LocalStore{
		maxSize:     maxSize,
		allocatedID: make([]byte, maxSize/bits+1)}
}

func (l *LocalStore) initialize() error {
	return nil
}

func (l *LocalStore) isAllocated(id int) (bool, error) {
	index := id / bits
	b := l.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b & mask
	return flag != 0, nil
}

func (l *LocalStore) allocate(id int) error {
	index := id / bits
	b := l.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b | mask
	l.allocatedID[index] = flag
	return nil
}

func (l *LocalStore) free(id int) error {
	index := id / bits
	b := l.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b & ^mask
	l.allocatedID[index] = flag
	return nil
}

func (l *LocalStore) freeAll() error {
	l.allocatedID = make([]byte, l.maxSize/bits+1)
	return nil
}

func (l *LocalStore) getMaxSize() int {
	return l.maxSize
}
