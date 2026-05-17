package quote

import (
	"encoding/base64"
	"log"

	"github.com/blinkmarket/pmm-server/internal/hexutil"
	"github.com/blinkmarket/pmm-server/internal/sign"
)

type Deps struct {
	Price     func(marketID [32]byte, side uint8, size uint64) (uint64, error)
	NextSeq   func() (uint64, error)
	Signer    sign.Signer
	Archive   func(Record) error
	PMM       [32]byte
	QuoteTTL  uint64 // milliseconds
	NowMillis func() uint64
}

type Service struct{ d Deps }

func New(d Deps) *Service { return &Service{d: d} }

func (s *Service) Quote(in QuoteInput) (*QuoteResponse, *HTTPError) {
	if in.Side != 0 && in.Side != 1 {
		return nil, &HTTPError{Status: 400, Code: "INVALID_INPUT", Message: "side must be 0 or 1"}
	}
	if in.Size == 0 {
		return nil, &HTTPError{Status: 400, Code: "INVALID_INPUT", Message: "size must be > 0"}
	}

	priceBps, err := s.d.Price(in.MarketID, in.Side, in.Size)
	if err != nil {
		return nil, &HTTPError{Status: 503, Code: "PRICER_UNAVAILABLE", Message: err.Error()}
	}

	seqNo, err := s.d.NextSeq()
	if err != nil {
		return nil, &HTTPError{Status: 500, Code: "SEQ_PERSIST_FAILED", Message: err.Error()}
	}

	expiresAt := s.d.NowMillis() + s.d.QuoteTTL

	msg := sign.EncodeMessage(in.MarketID, in.Side, priceBps, in.Size, seqNo, expiresAt, s.d.PMM)
	sig, pub, err := s.d.Signer.Sign(msg)
	if err != nil {
		return nil, &HTTPError{Status: 500, Code: "SIGN_FAILED", Message: err.Error()}
	}

	resp := &QuoteResponse{
		MarketID:  hexutil.Format32(in.MarketID),
		Side:      in.Side,
		PriceBps:  priceBps,
		Size:      in.Size,
		SeqNumber: seqNo,
		ExpiresAt: expiresAt,
		PMM:       hexutil.Format32(s.d.PMM),
		PMMPubkey: base64.StdEncoding.EncodeToString(pub),
		Signature: base64.StdEncoding.EncodeToString(sig),
	}

	// Best-effort archive: failures are logged, never block settlement (spec §9).
	if err := s.d.Archive(Record{
		MarketID: resp.MarketID, Side: resp.Side, PriceBps: resp.PriceBps,
		Size: resp.Size, SeqNumber: resp.SeqNumber, ExpiresAt: resp.ExpiresAt,
		PMM: resp.PMM, Signature: resp.Signature,
	}); err != nil {
		log.Printf("archive failed (non-blocking): seq=%d err=%v", seqNo, err)
	}

	return resp, nil
}
