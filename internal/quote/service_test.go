package quote

import (
	"crypto/ed25519"
	"testing"

	"github.com/blinkmarket/pmm-server/internal/sign"
)

func TestService_Quote_HappyPath(t *testing.T) {
	seed := make([]byte, 32)
	signer, _ := sign.NewEnvKeySigner(seed)
	var pmm [32]byte
	pmm[0] = 0xAB

	var archived int
	svc := New(Deps{
		Price:     func([32]byte, uint8, uint64) (uint64, error) { return 5000, nil },
		NextSeq:   func() (uint64, error) { return 7, nil },
		Signer:    signer,
		Archive:   func(Record) error { archived++; return nil },
		PMM:       pmm,
		QuoteTTL:  3000,
		NowMillis: func() uint64 { return 1000 },
	})

	var mkt [32]byte
	mkt[0] = 0x11
	resp, herr := svc.Quote(QuoteInput{MarketID: mkt, Side: 0, Size: 1_000_000})
	if herr != nil {
		t.Fatalf("unexpected err: %+v", herr)
	}
	if resp.SeqNumber != 7 || resp.PriceBps != 5000 || resp.ExpiresAt != 4000 {
		t.Fatalf("bad resp: %+v", resp)
	}
	if archived != 1 {
		t.Fatalf("archived=%d, want 1", archived)
	}
	// signature must verify over the canonical message
	msg := sign.EncodeMessage(mkt, 0, 5000, 1_000_000, 7, 4000, pmm)
	sigBytes := mustB64(t, resp.Signature)
	pubBytes := mustB64(t, resp.PMMPubkey)
	if !ed25519.Verify(ed25519.PublicKey(pubBytes), msg, sigBytes) {
		t.Fatal("response signature does not verify")
	}
}

func TestService_Quote_RejectsBadSide(t *testing.T) {
	svc := New(Deps{
		Price:     func([32]byte, uint8, uint64) (uint64, error) { return 1, nil },
		NextSeq:   func() (uint64, error) { return 1, nil },
		Signer:    mustSigner(t),
		Archive:   func(Record) error { return nil },
		QuoteTTL:  1000,
		NowMillis: func() uint64 { return 0 },
	})
	var mkt [32]byte
	_, herr := svc.Quote(QuoteInput{MarketID: mkt, Side: 9, Size: 10})
	if herr == nil || herr.Code != "INVALID_INPUT" {
		t.Fatalf("want INVALID_INPUT, got %+v", herr)
	}
}

func TestService_Quote_ArchiveFailureDoesNotBlock(t *testing.T) {
	svc := New(Deps{
		Price:     func([32]byte, uint8, uint64) (uint64, error) { return 5000, nil },
		NextSeq:   func() (uint64, error) { return 2, nil },
		Signer:    mustSigner(t),
		Archive:   func(Record) error { return errArchive },
		QuoteTTL:  1000,
		NowMillis: func() uint64 { return 0 },
	})
	var mkt [32]byte
	resp, herr := svc.Quote(QuoteInput{MarketID: mkt, Side: 1, Size: 10})
	if herr != nil {
		t.Fatalf("archive failure must not block: %+v", herr)
	}
	if resp.SeqNumber != 2 {
		t.Fatalf("bad resp seq: %d", resp.SeqNumber)
	}
}
