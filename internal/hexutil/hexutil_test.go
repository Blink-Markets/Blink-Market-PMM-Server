package hexutil

import "testing"

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
