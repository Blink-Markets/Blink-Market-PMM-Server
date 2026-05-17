package seq

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FileSeqStore persists the last issued seq via atomic temp-file + rename.
type FileSeqStore struct{ path string }

func NewFileSeqStore(path string) *FileSeqStore { return &FileSeqStore{path: path} }

func (s *FileSeqStore) Read() (uint64, error) {
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("seqstore read: %w", err)
	}
	v, err := strconv.ParseUint(strings.TrimSpace(string(b)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("seqstore parse: %w", err)
	}
	return v, nil
}

func (s *FileSeqStore) Write(v uint64) error {
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strconv.FormatUint(v, 10)), 0o600); err != nil {
		return fmt.Errorf("seqstore write tmp: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("seqstore rename: %w", err)
	}
	return nil
}
