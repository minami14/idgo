package idgo

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"sync"
)

// IDGenerateServer generate id when requested by client.
type IDGenerateServer struct {
	generator *IDGenerator
	addr      *net.TCPAddr
	listener  *net.TCPListener
	mutex     *sync.Mutex
	logger    *log.Logger
	status    byte
}

const (
	stop = iota
	run
	pause
)

// NewServer is IDGenerateServer constructed.
func NewServer(maxSize int, tcpAddr *net.TCPAddr) (*IDGenerateServer, error) {
	gen, err := NewIDGenerator(maxSize)
	if err != nil {
		return nil, err
	}
	return &IDGenerateServer{
		generator: gen,
		addr:      tcpAddr,
		mutex:     &sync.Mutex{},
		logger:    &log.Logger{},
		status:    stop,
	}, nil
}

// SetLogger is setter for logger.
func (s *IDGenerateServer) SetLogger(logger *log.Logger) {
	s.logger = logger
}

// Run server.
func (s *IDGenerateServer) Run() error {
	s.mutex.Lock()
	if s.status == run {
		s.mutex.Unlock()
		return errors.New("server is already running")
	}
	s.status = run
	listener, err := net.ListenTCP("tcp", s.addr)
	if err != nil {
		s.mutex.Unlock()
		return err
	}
	s.listener = listener
	s.mutex.Unlock()
	s.run()
	return nil
}

// Pause server while maintaining allocated id.
func (s *IDGenerateServer) Pause() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.status != run {
		return errors.New("server is not running")
	}
	s.status = pause
	return nil
}

// Stop server and free all allocated id.
func (s *IDGenerateServer) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.generator.FreeAll()
	return nil
}

func (s *IDGenerateServer) run() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			s.logger.Println(err)
		}
		if err := s.serve(conn); err != nil {
			s.logger.Println(err)
		}
	}
}

const (
	generate = iota
	allocate
	free
	freeAll
	isAllocated
	allocatedIDCount
	ping
	pong
	disconnect
)

func (s *IDGenerateServer) serve(conn *net.TCPConn) error {
	buf := make([]byte, 8)
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Println(err)
		}
	}()
	for {
		if _, err := conn.Read(buf[:1]); err != nil {
			return err
		}
		switch buf[0] {
		case generate:
			if err := s.generate(conn); err != nil {
				return err
			}
		case allocate:
			if err := s.allocate(conn); err != nil {
				return err
			}
		case free:
			if err := s.free(conn); err != nil {
				return err
			}
		case freeAll:
			if err := s.freeAll(conn); err != nil {
				return err
			}
		case isAllocated:
			if err := s.isAllocated(conn); err != nil {
				return err
			}
		case allocatedIDCount:
			if err := s.getAllocatedIDCount(conn); err != nil {
				return err
			}
		case ping:
			if err := s.pong(conn); err != nil {
				return err
			}
		case disconnect:
			return nil
		default:
			return errors.New("invalid method number")
		}
	}
}

const (
	ok = iota
	exhausted
)

func (s *IDGenerateServer) generate(conn *net.TCPConn) error {
	id, err := s.generator.Generate()
	if err != nil {
		s.logger.Println(err)
		if _, err := conn.Write([]byte{exhausted}); err != nil {
			return err
		}
	}
	idByte := make([]byte, 9)
	idByte[0] = ok
	binary.LittleEndian.PutUint64(idByte[1:], uint64(id))
	if _, err := conn.Write(idByte); err != nil {
		return err
	}
	return nil
}

func (s *IDGenerateServer) allocate(conn *net.TCPConn) error {
	buf := make([]byte, 8)
	if _, err := conn.Read(buf); err != nil {
		return err
	}
	id := int(binary.LittleEndian.Uint64(buf))
	err := s.generator.Allocate(id)
	return err
}

func (s *IDGenerateServer) free(conn *net.TCPConn) error {
	buf := make([]byte, 8)
	if _, err := conn.Read(buf); err != nil {
		return err
	}
	id := int(binary.LittleEndian.Uint64(buf))
	s.generator.Free(id)
	return nil
}

func (s *IDGenerateServer) freeAll(conn *net.TCPConn) error {
	s.generator.FreeAll()
	return nil
}

func (s *IDGenerateServer) isAllocated(conn *net.TCPConn) error {
	buf := make([]byte, 8)
	if _, err := conn.Read(buf); err != nil {
		return err
	}
	id := int(binary.LittleEndian.Uint64(buf))
	buf = []byte{0}
	isAlloc := s.generator.IsAllocated(id)
	if isAlloc {
		buf[0] = 1
	}
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func (s *IDGenerateServer) getAllocatedIDCount(conn *net.TCPConn) error {
	count := uint64(s.generator.GetAllocatedIDCount())
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, count)
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	return nil
}

func (s *IDGenerateServer) pong(conn *net.TCPConn) error {
	if _, err := conn.Write([]byte{pong}); err != nil {
		return err
	}

	return nil
}
