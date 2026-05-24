package seq

import (
	"fmt"
	"os"
	"path/filepath"
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
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("seqstore write tmp: %w", err)
	}
	_, werr := fmt.Fprintf(f, "%d", v)
	serr := f.Sync()
	cerr := f.Close()
	if werr != nil {
		return fmt.Errorf("seqstore write tmp: %w", werr)
	}
	if serr != nil {
		return fmt.Errorf("seqstore fsync tmp: %w", serr)
	}
	if cerr != nil {
		return fmt.Errorf("seqstore close tmp: %w", cerr)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("seqstore rename: %w", err)
	}
	// fsync the parent directory so the rename itself survives a crash:
	// without this the seq could regress on restart and a signed seq be
	// reissued, producing quotes the chain will reject as stale.
	df, err := os.Open(filepath.Dir(s.path))
	if err != nil {
		return fmt.Errorf("seqstore fsync dir: %w", err)
	}
	derr := df.Sync()
	df.Close()
	if derr != nil {
		return fmt.Errorf("seqstore fsync dir: %w", derr)
	}
	return nil
}
