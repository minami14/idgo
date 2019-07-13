package idgo

import (
	"log"
	"math"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

const address = ":4000"

func RunServer(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	s, err := NewServer(math.MaxInt16, tcpAddr, logger)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
}

func TestGenerateIDByServer(t *testing.T) {
	RunServer(t)

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	client := NewClient()
	if err := client.Connect(tcpAddr); err != nil {
		log.Fatal(err)
	}

	used := make([]bool, math.MaxInt16)
	m := &sync.Mutex{}
	for i := 0; i < math.MaxInt16; i++ {
		id, err := client.Generate()
		if err != nil {
			t.Error(err)
		}
		if !(id == i) {
			t.Errorf("invalid id %b", id)
			return
		}
		if err := client.Free(id); err != nil {
			t.Error(err)
		}
	}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				id, err := client.Generate()
				if err != nil {
					t.Error(err)
				}
				m.Lock()
				if used[id] {
					t.Errorf("used id %b", id)
				}
				used[id] = true
				m.Unlock()
			}
		}()

	}
}
