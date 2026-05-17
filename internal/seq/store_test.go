package seq

import (
	"path/filepath"
	"testing"
)

func TestFileSeqStore_RoundTrip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "seq.state")
	s := NewFileSeqStore(p)

	v, err := s.Read()
	if err != nil {
		t.Fatalf("read empty: %v", err)
	}
	if v != 0 {
		t.Fatalf("empty store = %d, want 0", v)
	}
	if err := s.Write(42); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := NewFileSeqStore(p).Read()
	if err != nil {
		t.Fatalf("reread: %v", err)
	}
	if got != 42 {
		t.Fatalf("reread = %d, want 42", got)
	}
}
