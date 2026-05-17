package seq

// ChainSeqReader reports the last seq accepted on-chain for a PMM address.
// v1 uses ConfigChainSeqReader; a live Sui RPC reader is deferred (plan Task 13).
type ChainSeqReader interface {
	LastSeq(pmm [32]byte) (uint64, error)
}

type ConfigChainSeqReader struct{ last uint64 }

func NewConfigChainSeqReader(last uint64) *ConfigChainSeqReader {
	return &ConfigChainSeqReader{last: last}
}

func (r *ConfigChainSeqReader) LastSeq([32]byte) (uint64, error) { return r.last, nil }
