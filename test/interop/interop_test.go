package interop

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/blinkmarket/pmm-server/internal/sign"
)

type fixture struct {
	SeedHex            string `json:"private_key_seed_hex"`
	PubB64             string `json:"expected_pubkey_b64"`
	ExpectedMessageHex string `json:"expected_message_hex"`
	Inputs             struct {
		MarketIDHex string `json:"market_id_hex"`
		Side        uint8  `json:"side"`
		PriceBps    uint64 `json:"price_bps"`
		Size        uint64 `json:"size"`
		SeqNumber   uint64 `json:"seq_number"`
		ExpiresAt   uint64 `json:"expires_at"`
		PMMHex      string `json:"pmm_hex"`
	} `json:"inputs"`
}

func h32(t *testing.T, s string) [32]byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 32 {
		t.Fatalf("bad 32-hex %q: %v", s, err)
	}
	var a [32]byte
	copy(a[:], b)
	return a
}

func TestCrossEndSignatureInterop(t *testing.T) {
	raw, err := os.ReadFile("vectors.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var fx fixture
	if err := json.Unmarshal(raw, &fx); err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	seed := h32(t, fx.SeedHex)
	signer, err := sign.NewEnvKeySigner(seed[:])
	if err != nil {
		t.Fatalf("signer: %v", err)
	}
	in := fx.Inputs
	msg := sign.EncodeMessage(
		h32(t, in.MarketIDHex), in.Side, in.PriceBps, in.Size,
		in.SeqNumber, in.ExpiresAt, h32(t, in.PMMHex),
	)
	if got := hex.EncodeToString(msg); got != fx.ExpectedMessageHex {
		t.Fatalf("canonical message mismatch:\n got=%s\nwant=%s", got, fx.ExpectedMessageHex)
	}
	sig, pub, err := signer.Sign(msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	wantPub, err := base64.StdEncoding.DecodeString(fx.PubB64)
	if err != nil {
		t.Fatalf("decode want pub: %v", err)
	}
	if string(pub) != string(wantPub) {
		t.Fatalf("pubkey mismatch:\n got=%x\nwant=%x", pub, wantPub)
	}
	if !ed25519.Verify(ed25519.PublicKey(pub), msg, sig) {
		t.Fatal("signature does not verify over canonical message")
	}
}
