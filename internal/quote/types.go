package quote

// QuoteInput is the validated internal form of a quote request.
type QuoteInput struct {
	MarketID [32]byte
	Side     uint8
	Size     uint64
}

// QuoteRequest is the wire request body.
type QuoteRequest struct {
	MarketID string `json:"market_id"`
	Side     uint8  `json:"side"`
	Size     uint64 `json:"size"`
}

// QuoteResponse is the wire response body (fields map to execute_rfq, spec §5).
type QuoteResponse struct {
	MarketID  string `json:"market_id"`
	Side      uint8  `json:"side"`
	PriceBps  uint64 `json:"price_bps"`
	Size      uint64 `json:"size"`
	SeqNumber uint64 `json:"seq_number"`
	ExpiresAt uint64 `json:"expires_at"`
	PMM       string `json:"pmm"`
	PMMPubkey string `json:"pmm_pubkey"`
	Signature string `json:"signature"`
}

// Record mirrors archive.Record without importing it (avoids a cycle).
type Record struct {
	MarketID  string
	Side      uint8
	PriceBps  uint64
	Size      uint64
	SeqNumber uint64
	ExpiresAt uint64
	PMM       string
	Signature string
}

// HTTPError carries an API error code + HTTP status (spec §5 error schema).
type HTTPError struct {
	Status  int
	Code    string
	Message string
}

func (e *HTTPError) Error() string { return e.Code + ": " + e.Message }
