package seq

import "sync"

// SeqAllocator issues strictly increasing seq numbers for a single PMM key.
// Startup: counter = max(chainLast, fileLast) + 1 (spec §7).
type SeqAllocator struct {
	mu      sync.Mutex
	counter uint64
	store   *FileSeqStore
}

func NewSeqAllocator(store *FileSeqStore, chain ChainSeqReader, pmm [32]byte) (*SeqAllocator, error) {
	chainLast, err := chain.LastSeq(pmm)
	if err != nil {
		return nil, err
	}
	fileLast, err := store.Read()
	if err != nil {
		return nil, err
	}
	start := chainLast
	if fileLast > start {
		start = fileLast
	}
	if err := store.Write(start); err != nil {
		return nil, err
	}
	return &SeqAllocator{counter: start, store: store}, nil
}

func (a *SeqAllocator) Next() (uint64, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	n := a.counter + 1
	if err := a.store.Write(n); err != nil {
		return 0, err
	}
	a.counter = n
	return n, nil
}
