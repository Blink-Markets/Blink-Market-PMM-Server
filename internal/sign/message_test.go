package sign

import (
	"encoding/hex"
	"testing"
)

func TestEncodeMessage_GoldenVector(t *testing.T) {
	var marketID, pmm [32]byte
	for i := range marketID {
		marketID[i] = 0x11
		pmm[i] = 0x22
	}
	got := EncodeMessage(marketID, 1, 5000, 1000000, 42, 1000, pmm)

	want := ""
	for i := 0; i < 32; i++ {
		want += "11"
	}
	want += "01"
	want += "8813000000000000"
	want += "40420f0000000000"
	want += "2a00000000000000"
	want += "e803000000000000"
	for i := 0; i < 32; i++ {
		want += "22"
	}

	if hex.EncodeToString(got) != want {
		t.Fatalf("message mismatch:\n got=%s\nwant=%s", hex.EncodeToString(got), want)
	}
	if len(got) != 97 {
		t.Fatalf("len=%d, want 97", len(got))
	}
}
