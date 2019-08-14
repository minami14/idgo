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

const (
	localStoreAddress = ":4000"
	redisStoreAddress = ":4001"
)

func RunServer(t *testing.T, store AllocatedIDStore, address string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		t.Fatal(err)
	}

	s, err := NewServer(store, tcpAddr)
	if err != nil {
		t.Fatal(err)
	}

	s.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	go func() {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
}

const maxSize = math.MaxInt8

func GenerateTest(t *testing.T, store AllocatedIDStore, address string) {
	RunServer(t, store, address)

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient()
	if err := client.Connect(tcpAddr); err != nil {
		t.Fatal(err)
	}

	used := make([]bool, maxSize)
	m := &sync.Mutex{}
	for i := 0; i < maxSize; i++ {
		id, err := client.Generate()
		if err != nil {
			t.Error(err)
		}
		if !(id == i) {
			t.Errorf("invalid id %v", id)
			return
		}
		if err := client.Free(id); err != nil {
			t.Error(err)
		}
	}

	maxSizeSqrt := int(math.Sqrt(float64(maxSize)))
	for i := 0; i < maxSizeSqrt; i++ {
		go func() {
			for j := 0; j < maxSizeSqrt; j++ {
				id, err := client.Generate()
				if err != nil {
					t.Error(err)
				}
				m.Lock()
				if used[id] {
					t.Errorf("used id %v", id)
				}
				used[id] = true
				m.Unlock()
			}
		}()

	}
}

func TestLocalStore(t *testing.T) {
	store, err := NewLocalStore(maxSize)
	if err != nil {
		t.Fatal(err)
	}

	GenerateTest(t, store, localStoreAddress)
}

func TestRedisStore(t *testing.T) {
	store, err := NewRedisStore("127.0.0.1:6379", "idgo-test", maxSize)
	if err != nil {
		t.Fatal(err)
	}

	if err := store.freeAll(); err != nil {
		t.Fatal(err)
	}

	GenerateTest(t, store, redisStoreAddress)
}

func BenchmarkLocalStore(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	store, err := NewLocalStore(maxSize)
	if err != nil {
		b.Fatal(err)
	}

	gen, err := NewIDGenerator(store)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < maxSize; j++ {
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

	store, err := NewRedisStore("127.0.0.1:6379", "idgo-bench", maxSize)
	if err != nil {
		b.Fatal(err)
	}

	if err := store.freeAll(); err != nil {
		b.Fatal(err)
	}

	gen, err := NewIDGenerator(store)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		for j := 0; j < maxSize; j++ {
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
