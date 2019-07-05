package idgo

import (
	"encoding/binary"
	"errors"
	"net"
	"sync"
)

type IDGenerateClient struct {
	mutex *sync.Mutex
	conn  *net.TCPConn
}

// NewClient is IDGenerateClient constructed.
func NewClient() *IDGenerateClient {
	return &IDGenerateClient{
		mutex: &sync.Mutex{},
	}
}

// Connect to server
func (c *IDGenerateClient) Connect(addr *net.TCPAddr) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// Generate a new id
func (c *IDGenerateClient) Generate() (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := c.conn.Write([]byte{generate}); err != nil {
		return 0, err
	}

	buf := make([]byte, 8)
	if _, err := c.conn.Read(buf[:1]); err != nil {
		return 0, err
	}

	switch buf[0] {
	case ok:
	case exhausted:
		return 0, errors.New("id is exhausted")
	default:
		return 0, errors.New("failed to generate id")
	}

	if _, err := c.conn.Read(buf); err != nil {
		return 0, err
	}

	id := int(binary.LittleEndian.Uint64(buf))
	return id, nil
}

// Allocate a specified id
func (c *IDGenerateClient) Allocate(id int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := make([]byte, 9)
	buf[0] = allocate
	binary.LittleEndian.PutUint64(buf[1:], uint64(id))
	if _, err := c.conn.Write(buf); err != nil {
		return err
	}

	return nil
}

// Free a allocated id
func (c *IDGenerateClient) Free(id int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := make([]byte, 9)
	buf[0] = free
	binary.LittleEndian.PutUint64(buf[1:], uint64(id))
	if _, err := c.conn.Write(buf); err != nil {
		return err
	}

	return nil
}

// FreeAll free all allocated id
func (c *IDGenerateClient) FreeAll() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := make([]byte, 1, 8)
	buf[0] = freeAll
	if _, err := c.conn.Write(buf); err != nil {
		return err
	}

	return nil
}

// IsAllocated check if specified id is allocated
func (c *IDGenerateClient) IsAllocated(id int) (bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := make([]byte, 9)
	buf[0] = isAllocated
	binary.LittleEndian.PutUint64(buf[1:], uint64(id))
	if _, err := c.conn.Write(buf); err != nil {
		return false, err
	}

	if _, err := c.conn.Read(buf[:1]); err != nil {
		return false, err
	}

	return buf[0] == 1, nil
}

// GetAllocatedIDCount is getter for allocatedIDCount
func (c *IDGenerateClient) GetAllocatedIDCount(id int) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := c.conn.Write([]byte{allocatedIDCount}); err != nil {
		return 0, err
	}

	buf := make([]byte, 8)
	if _, err := c.conn.Read(buf); err != nil {
		return 0, err
	}

	count := int(binary.LittleEndian.Uint64(buf))
	return count, nil
}

// Close connection
func (c *IDGenerateClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := []byte{disconnect}
	_, err := c.conn.Write(buf)
	return err
}
