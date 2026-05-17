package pricing

// Pricer computes the quoted probability in basis points (0..10000).
type Pricer interface {
	Price(marketID [32]byte, side uint8, size uint64) (uint64, error)
}

type StubPricer struct{ bps uint64 }

func NewStubPricer(bps uint64) *StubPricer { return &StubPricer{bps: bps} }

func (p *StubPricer) Price([32]byte, uint8, uint64) (uint64, error) { return p.bps, nil }
