package quote

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/blinkmarket/pmm-server/internal/sign"
)

var errArchive = errors.New("archive down")

func mustSigner(t *testing.T) sign.Signer {
	t.Helper()
	s, err := sign.NewEnvKeySigner(make([]byte, 32))
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	return s
}

func mustB64(t *testing.T, s string) []byte {
	t.Helper()
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatalf("b64: %v", err)
	}
	return b
}
