package hexutil

import (
	"strings"
	"testing"
)

func TestParse32_StripsPrefixAndDecodes(t *testing.T) {
	in := "0x" + "11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"+"11"
	got, err := Parse32(in)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	for i, b := range got {
		if b != 0x11 {
			t.Fatalf("byte %d = %#x, want 0x11", i, b)
		}
	}
}

func TestParse32_RejectsWrongLength(t *testing.T) {
	if _, err := Parse32("0x1234"); err == nil {
		t.Fatal("expected error for short input")
	}
}

func TestFormat32_AddsPrefix(t *testing.T) {
	var a [32]byte
	a[31] = 0x01
	got := Format32(a)
	if got[:2] != "0x" || len(got) != 66 {
		t.Fatalf("bad format: %q", got)
	}
}

func TestParse32_NoPrefix(t *testing.T) {
	got, err := Parse32(strings.Repeat("ab", 32))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	for i, b := range got {
		if b != 0xab {
			t.Fatalf("byte %d = %#x, want 0xab", i, b)
		}
	}
}

func TestParse32_RejectsInvalidChars(t *testing.T) {
	if _, err := Parse32("0x" + strings.Repeat("zz", 32)); err == nil {
		t.Fatal("expected error for non-hex chars")
	}
}

func TestFormat32_RoundTripContent(t *testing.T) {
	var a [32]byte
	a[0] = 0xde
	a[31] = 0x01
	got := Format32(a)
	want := "0x" + strings.Repeat("00", 31)
	_ = want
	if got[:4] != "0xde" {
		t.Fatalf("prefix content wrong: %q", got[:4])
	}
	if got[len(got)-2:] != "01" {
		t.Fatalf("suffix content wrong: %q", got[len(got)-2:])
	}
	back, err := Parse32(got)
	if err != nil || back != a {
		t.Fatalf("round-trip failed: back=%x err=%v", back, err)
	}
}
