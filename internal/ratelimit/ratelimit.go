// Package ratelimit is a simple per-key token bucket (spec §5: per-market limit).
package ratelimit

import (
	"sync"
	"time"
)

type bucket struct {
	tokens float64
	last   time.Time
}

type Limiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	capacity float64
	refill   float64 // tokens per second
}

// maxBuckets triggers a stale-bucket sweep so a flood of distinct keys
// (e.g. random valid market IDs) cannot grow the map without bound.
const maxBuckets = 10000

// New builds a limiter with capacity == refill-per-second == ratePerSec.
func New(ratePerSec float64) *Limiter {
	return &Limiter{
		buckets:  map[string]*bucket{},
		capacity: ratePerSec,
		refill:   ratePerSec,
	}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if len(l.buckets) > maxBuckets {
		l.sweepLocked(now)
	}
	b := l.buckets[key]
	if b == nil {
		b = &bucket{tokens: l.capacity, last: now}
		l.buckets[key] = b
	}
	b.tokens += now.Sub(b.last).Seconds() * l.refill
	if b.tokens > l.capacity {
		b.tokens = l.capacity
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// sweepLocked drops buckets idle long enough to have fully refilled. Since
// capacity == refill, a bucket untouched for capacity/refill seconds is back
// at full capacity and behaviorally identical to a fresh one, so deleting it
// is lossless. Caller must hold l.mu.
func (l *Limiter) sweepLocked(now time.Time) {
	ttl := time.Duration(l.capacity / l.refill * float64(time.Second))
	for k, b := range l.buckets {
		if now.Sub(b.last) >= ttl {
			delete(l.buckets, k)
		}
	}
}
