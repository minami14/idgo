package idgo

import (
	"math"
	"sync"
	"testing"
)

func TestGenerateID(t *testing.T) {
	gen, err := NewIDGenerator(math.MaxInt16)
	if err != nil {
		t.Error(err)
	}
	used := make([]bool, math.MaxInt16)
	m := &sync.Mutex{}
	for i := 0; i < math.MaxInt16; i++ {
		id, err := gen.Generate()
		if err != nil {
			t.Error(err)
		}
		if !(id == i) {
			t.Errorf("invalid id %b", id)
			return
		}
		gen.Free(id)
	}
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				id, err := gen.Generate()
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
