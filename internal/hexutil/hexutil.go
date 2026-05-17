// Package hexutil parses and formats fixed 32-byte hex values (0x-optional).
package hexutil

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func Parse32(s string) ([32]byte, error) {
	var out [32]byte
	s = strings.TrimPrefix(s, "0x")
	if len(s) != 64 {
		return out, fmt.Errorf("hexutil: want 64 hex chars, got %d", len(s))
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return out, fmt.Errorf("hexutil: %w", err)
	}
	copy(out[:], b)
	return out, nil
}

func Format32(a [32]byte) string {
	return "0x" + hex.EncodeToString(a[:])
}
