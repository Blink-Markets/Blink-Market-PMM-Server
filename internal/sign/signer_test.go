package sign

import (
	"crypto/ed25519"
	"testing"
)

func TestEnvKeySigner_SignVerifies(t *testing.T) {
	seed := make([]byte, 32)
	seed[0] = 7
	s, err := NewEnvKeySigner(seed)
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	msg := []byte("hello")
	sig, pub, err := s.Sign(msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if len(sig) != 64 || len(pub) != 32 {
		t.Fatalf("sig=%d pub=%d", len(sig), len(pub))
	}
	if !ed25519.Verify(ed25519.PublicKey(pub), msg, sig) {
		t.Fatal("signature did not verify")
	}
}

func TestNewEnvKeySigner_RejectsBadSeed(t *testing.T) {
	if _, err := NewEnvKeySigner(make([]byte, 16)); err == nil {
		t.Fatal("expected error for 16-byte seed")
	}
}
