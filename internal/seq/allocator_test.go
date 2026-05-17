package seq

import (
	"path/filepath"
	"sync"
	"testing"
)

func newTestAllocator(t *testing.T, chainLast, fileLast uint64) *SeqAllocator {
	t.Helper()
	p := filepath.Join(t.TempDir(), "seq.state")
	st := NewFileSeqStore(p)
	// fileLast == 0 is treated as "no pre-existing file": FileSeqStore.Read() returns 0 for a missing file, so seeding 0 would be equivalent.
	if fileLast != 0 {
		if err := st.Write(fileLast); err != nil {
			t.Fatalf("seed file: %v", err)
		}
	}
	var pmm [32]byte
	a, err := NewSeqAllocator(st, NewConfigChainSeqReader(chainLast), pmm)
	if err != nil {
		t.Fatalf("new allocator: %v", err)
	}
	return a
}

func TestSeqAllocator_StartsAtMaxPlusOne(t *testing.T) {
	a := newTestAllocator(t, 100, 40) // max(100,40)+1
	n, err := a.Next()
	if err != nil {
		t.Fatalf("next: %v", err)
	}
	if n != 101 {
		t.Fatalf("got %d, want 101", n)
	}
}

func TestSeqAllocator_FileAhead(t *testing.T) {
	a := newTestAllocator(t, 5, 70) // max(5,70)+1
	n, err := a.Next()
	if err != nil {
		t.Fatalf("next: %v", err)
	}
	if n != 71 {
		t.Fatalf("got %d, want 71", n)
	}
}

func TestSeqAllocator_StrictlyIncreasingConcurrent(t *testing.T) {
	a := newTestAllocator(t, 0, 0)
	const n = 200
	seen := make([]uint64, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, err := a.Next()
			if err != nil {
				t.Errorf("next: %v", err)
				return
			}
			seen[i] = v
		}(i)
	}
	wg.Wait()
	dup := map[uint64]bool{}
	for _, v := range seen {
		if v == 0 || dup[v] {
			t.Fatalf("dup or zero seq: %d", v)
		}
		dup[v] = true
	}
}
