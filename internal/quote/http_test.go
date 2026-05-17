package quote

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	return New(Deps{
		Price:     func([32]byte, uint8, uint64) (uint64, error) { return 5000, nil },
		NextSeq:   func() (uint64, error) { return 1, nil },
		Signer:    mustSigner(t),
		Archive:   func(Record) error { return nil },
		QuoteTTL:  3000,
		NowMillis: func() uint64 { return 0 },
	})
}

func TestHandler_OK(t *testing.T) {
	h := NewHandler(newTestService(t), nil)
	body, _ := json.Marshal(QuoteRequest{
		MarketID: "0x" + "11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11",
		Side: 0, Size: 100,
	})
	r := httptest.NewRequest(http.MethodPost, "/v1/quote", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var resp QuoteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.SeqNumber != 1 {
		t.Fatalf("seq=%d", resp.SeqNumber)
	}
}

func TestHandler_BadMarketID(t *testing.T) {
	h := NewHandler(newTestService(t), nil)
	body, _ := json.Marshal(QuoteRequest{MarketID: "0x12", Side: 0, Size: 1})
	r := httptest.NewRequest(http.MethodPost, "/v1/quote", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("status=%d, want 400", w.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := NewHandler(newTestService(t), nil)
	r := httptest.NewRequest(http.MethodGet, "/v1/quote", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d, want 405", w.Code)
	}
}

type stubLimiter struct{ allow bool }

func (s stubLimiter) Allow(string) bool { return s.allow }

func TestHandler_RateLimited(t *testing.T) {
	h := NewHandler(newTestService(t), stubLimiter{allow: false})
	body, _ := json.Marshal(QuoteRequest{
		MarketID: "0x" + "11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11",
		Side: 0, Size: 100,
	})
	r := httptest.NewRequest(http.MethodPost, "/v1/quote", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("status=%d, want 429", w.Code)
	}
	var env struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if env.Error.Code != "RATE_LIMITED" {
		t.Fatalf("error code=%q, want RATE_LIMITED", env.Error.Code)
	}
}

func TestHandler_BadMarketID_EnvelopeShape(t *testing.T) {
	h := NewHandler(newTestService(t), nil)
	body, _ := json.Marshal(QuoteRequest{MarketID: "0x12", Side: 0, Size: 1})
	r := httptest.NewRequest(http.MethodPost, "/v1/quote", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("status=%d, want 400", w.Code)
	}
	var env struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if env.Error.Code != "INVALID_INPUT" {
		t.Fatalf("error code=%q, want INVALID_INPUT", env.Error.Code)
	}
}
