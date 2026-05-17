package seq

import "testing"

func TestConfigChainSeqReader_ReturnsConfigured(t *testing.T) {
	r := NewConfigChainSeqReader(99)
	var pmm [32]byte
	got, err := r.LastSeq(pmm)
	if err != nil {
		t.Fatalf("LastSeq: %v", err)
	}
	if got != 99 {
		t.Fatalf("got %d, want 99", got)
	}
}
