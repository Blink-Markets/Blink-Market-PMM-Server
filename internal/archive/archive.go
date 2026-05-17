package archive

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Record is one signed quote, archived for verifiability (spec §1, §9).
type Record struct {
	MarketID  string `json:"market_id"`
	Side      uint8  `json:"side"`
	PriceBps  uint64 `json:"price_bps"`
	Size      uint64 `json:"size"`
	SeqNumber uint64 `json:"seq_number"`
	ExpiresAt uint64 `json:"expires_at"`
	PMM       string `json:"pmm"`
	Signature string `json:"signature"`
}

// Archiver persists signed quotes. Failures must not block settlement (spec §9).
type Archiver interface {
	Record(rec Record) error
}

type JSONLArchiver struct {
	mu   sync.Mutex
	path string
}

func NewJSONLArchiver(path string) *JSONLArchiver { return &JSONLArchiver{path: path} }

func (a *JSONLArchiver) Record(rec Record) error {
	line, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("archive marshal: %w", err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("archive open: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("archive write: %w", err)
	}
	return nil
}

// NoopWalrusArchiver is the deferred-Walrus placeholder (spec §2, §10).
type NoopWalrusArchiver struct{}

func NewNoopWalrusArchiver() *NoopWalrusArchiver { return &NoopWalrusArchiver{} }

func (*NoopWalrusArchiver) Record(Record) error { return nil }
