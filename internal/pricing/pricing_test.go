package pricing

import "testing"

func TestStubPricer_ReturnsConfigured(t *testing.T) {
	p := NewStubPricer(5000)
	var m [32]byte
	bps, err := p.Price(m, 0, 1_000_000)
	if err != nil {
		t.Fatalf("price: %v", err)
	}
	if bps != 5000 {
		t.Fatalf("got %d, want 5000", bps)
	}
}
