package ratelimit

import "testing"

func TestTokenBucket_AllowsBurstThenBlocks(t *testing.T) {
	rl := New(2) // 2 tokens capacity, 2/sec refill
	k := "mkt1"
	if !rl.Allow(k) || !rl.Allow(k) {
		t.Fatal("first two should pass")
	}
	if rl.Allow(k) {
		t.Fatal("third should be blocked")
	}
}

func TestTokenBucket_PerKeyIsolated(t *testing.T) {
	rl := New(1)
	if !rl.Allow("a") {
		t.Fatal("a first should pass")
	}
	if !rl.Allow("b") {
		t.Fatal("b should be independent of a")
	}
}
