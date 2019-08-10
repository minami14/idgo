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
		t.Fatal(err)
	}

	store := NewLocalStore(math.MaxInt16)
	s, err := NewServer(store, tcpAddr)
	if err != nil {
		t.Fatal(err)
	}

	s.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	go func() {
		if err := s.Run(); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
}

func TestGenerateIDByServer(t *testing.T) {
	RunServer(t)

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient()
	if err := client.Connect(tcpAddr); err != nil {
		t.Fatal(err)
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

func BenchmarkLocalStore(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	store := NewLocalStore(math.MaxInt16)
	gen, err := NewIDGenerator(store)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < math.MaxInt16; j++ {
			id, err := gen.Generate()
			if err != nil {
				b.Error(err)
			}
			if err := gen.Free(id); err != nil {
				b.Error(err)
			}
		}
		if err := gen.FreeAll(); err != nil {
			b.Error()
		}
	}
}

func BenchmarkRedisStore(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	store, err := NewRedisStore("127.0.0.1:6379", "idgo", math.MaxInt16)
	if err != nil {
		b.Fatal(err)
	}
	gen, err := NewIDGenerator(store)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < math.MaxInt16; j++ {
			id, err := gen.Generate()
			if err != nil {
				b.Error(err)
			}
			if err := gen.Free(id); err != nil {
				b.Error(err)
			}
		}
		if err := gen.FreeAll(); err != nil {
			b.Error()
		}
	}
}
