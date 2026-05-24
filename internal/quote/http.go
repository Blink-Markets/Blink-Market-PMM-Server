package quote

import (
	"encoding/json"
	"net/http"

	"github.com/blinkmarket/pmm-server/internal/hexutil"
)

type rateLimiter interface{ Allow(key string) bool }

type Handler struct {
	svc *Service
	rl  rateLimiter // nil disables limiting
}

func NewHandler(svc *Service, rl rateLimiter) *Handler { return &Handler{svc: svc, rl: rl} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErr(w, &HTTPError{Status: http.StatusMethodNotAllowed, Code: "METHOD_NOT_ALLOWED", Message: "use POST"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 512)
	var req QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, &HTTPError{Status: 400, Code: "INVALID_INPUT", Message: "malformed JSON body"})
		return
	}
	mkt, err := hexutil.Parse32(req.MarketID)
	if err != nil {
		writeErr(w, &HTTPError{Status: 400, Code: "INVALID_INPUT", Message: "market_id must be 32-byte hex"})
		return
	}
	if h.rl != nil && !h.rl.Allow(hexutil.Format32(mkt)) {
		writeErr(w, &HTTPError{Status: 429, Code: "RATE_LIMITED", Message: "too many quote requests for this market"})
		return
	}
	resp, herr := h.svc.Quote(QuoteInput{MarketID: mkt, Side: req.Side, Size: req.Size})
	if herr != nil {
		writeErr(w, herr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func writeErr(w http.ResponseWriter, e *HTTPError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{"code": e.Code, "message": e.Message},
	})
}
