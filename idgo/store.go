package idgo

// AllocatedIDStore stores allocated id.
type AllocatedIDStore interface {
	initialize()
	isAllocated(int) bool
	allocate(int)
	free(int)
	freeAll()
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

func (l *LocalStore) initialize() {}

func (l *LocalStore) isAllocated(id int) bool {
	index := id / bits
	b := l.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b & mask
	return flag != 0
}

func (l *LocalStore) allocate(id int) {
	index := id / bits
	b := l.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b | mask
	l.allocatedID[index] = flag
}

func (l *LocalStore) free(id int) {
	index := id / bits
	b := l.allocatedID[index]
	shift := byte(id % bits)
	mask := byte(1 << shift)
	flag := b & ^mask
	l.allocatedID[index] = flag
}

func (l *LocalStore) freeAll() {
	l.allocatedID = make([]byte, l.maxSize/bits+1)
}

func (l *LocalStore) getMaxSize() int {
	return l.maxSize
}
