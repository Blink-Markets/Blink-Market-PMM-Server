package sign

import (
	"crypto/ed25519"
	"fmt"
)

// Signer signs the canonical message and exposes the verifying pubkey.
type Signer interface {
	Sign(msg []byte) (sig []byte, pubkey []byte, err error)
}

type EnvKeySigner struct {
	priv ed25519.PrivateKey
	pub  ed25519.PublicKey
}

// NewEnvKeySigner builds a signer from a 32-byte Ed25519 seed.
func NewEnvKeySigner(seed []byte) (*EnvKeySigner, error) {
	if len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("sign: seed must be %d bytes, got %d", ed25519.SeedSize, len(seed))
	}
	priv := ed25519.NewKeyFromSeed(seed)
	return &EnvKeySigner{priv: priv, pub: priv.Public().(ed25519.PublicKey)}, nil
}

func (s *EnvKeySigner) Sign(msg []byte) ([]byte, []byte, error) {
	return ed25519.Sign(s.priv, msg), append([]byte(nil), s.pub...), nil
}
