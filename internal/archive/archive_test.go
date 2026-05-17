package archive

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJSONLArchiver_AppendsLine(t *testing.T) {
	p := filepath.Join(t.TempDir(), "q.jsonl")
	a := NewJSONLArchiver(p)
	rec := Record{MarketID: "0xabc", Side: 1, PriceBps: 5000, Size: 10, SeqNumber: 7, ExpiresAt: 100}
	if err := a.Record(rec); err != nil {
		t.Fatalf("record: %v", err)
	}
	if err := a.Record(rec); err != nil {
		t.Fatalf("record 2: %v", err)
	}
	f, _ := os.Open(p)
	defer f.Close()
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if !strings.Contains(sc.Text(), `"seq_number":7`) {
			t.Fatalf("bad line: %s", sc.Text())
		}
		n++
	}
	if n != 2 {
		t.Fatalf("lines=%d, want 2", n)
	}
}

func TestNoopWalrusArchiver_OK(t *testing.T) {
	if err := NewNoopWalrusArchiver().Record(Record{}); err != nil {
		t.Fatalf("noop should not error: %v", err)
	}
}
