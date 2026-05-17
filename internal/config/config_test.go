package config

import (
	"testing"
)

func TestLoad_DefaultsAndRequired(t *testing.T) {
	t.Setenv("PMM_PRIVATE_KEY_HEX", "00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00"+"00")
	t.Setenv("PMM_ADDRESS_HEX", "0x"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22"+"22")
	c, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.ListenAddr != ":8080" || c.QuoteTTLMillis != 3000 || c.StubPriceBps != 5000 {
		t.Fatalf("bad defaults: %+v", c)
	}
	if len(c.PrivateKeySeed) != 32 {
		t.Fatalf("seed len=%d", len(c.PrivateKeySeed))
	}
}

func TestLoad_MissingKeyFails(t *testing.T) {
	if _, err := Load(); err == nil {
		t.Fatal("expected error when PMM_PRIVATE_KEY_HEX unset")
	}
}
