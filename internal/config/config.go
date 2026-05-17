package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/blinkmarket/pmm-server/internal/hexutil"
)

type Config struct {
	ListenAddr     string
	PrivateKeySeed []byte
	PMMAddress     [32]byte
	QuoteTTLMillis uint64
	StubPriceBps   uint64
	SeqFile        string
	ChainLastSeq   uint64
	ArchiveFile    string
	RateLimitQPS   float64
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func envU64(k string, def uint64) (uint64, error) {
	v := os.Getenv(k)
	if v == "" {
		return def, nil
	}
	return strconv.ParseUint(v, 10, 64)
}

func Load() (*Config, error) {
	keyHex := os.Getenv("PMM_PRIVATE_KEY_HEX")
	if keyHex == "" {
		return nil, fmt.Errorf("config: PMM_PRIVATE_KEY_HEX is required")
	}
	seed, err := hexutil.Parse32(keyHex)
	if err != nil {
		return nil, fmt.Errorf("config: PMM_PRIVATE_KEY_HEX: %w", err)
	}
	addrHex := os.Getenv("PMM_ADDRESS_HEX")
	if addrHex == "" {
		return nil, fmt.Errorf("config: PMM_ADDRESS_HEX is required")
	}
	addr, err := hexutil.Parse32(addrHex)
	if err != nil {
		return nil, fmt.Errorf("config: PMM_ADDRESS_HEX: %w", err)
	}
	ttl, err := envU64("PMM_QUOTE_TTL_MS", 3000)
	if err != nil {
		return nil, fmt.Errorf("config: PMM_QUOTE_TTL_MS: %w", err)
	}
	if ttl == 0 {
		return nil, fmt.Errorf("config: PMM_QUOTE_TTL_MS must be > 0")
	}
	price, err := envU64("PMM_STUB_PRICE_BPS", 5000)
	if err != nil {
		return nil, fmt.Errorf("config: PMM_STUB_PRICE_BPS: %w", err)
	}
	chainSeq, err := envU64("PMM_CHAIN_LAST_SEQ", 0)
	if err != nil {
		return nil, fmt.Errorf("config: PMM_CHAIN_LAST_SEQ: %w", err)
	}
	rlQPS := 5.0
	if v := os.Getenv("PMM_RATE_LIMIT_PER_MARKET_QPS"); v != "" {
		f, perr := strconv.ParseFloat(v, 64)
		if perr != nil {
			return nil, fmt.Errorf("config: PMM_RATE_LIMIT_PER_MARKET_QPS: %w", perr)
		}
		rlQPS = f
	}
	if rlQPS <= 0 {
		return nil, fmt.Errorf("config: PMM_RATE_LIMIT_PER_MARKET_QPS must be > 0, got %g", rlQPS)
	}
	return &Config{
		ListenAddr:     env("PMM_LISTEN_ADDR", ":8080"),
		PrivateKeySeed: seed[:],
		PMMAddress:     addr,
		QuoteTTLMillis: ttl,
		StubPriceBps:   price,
		SeqFile:        env("PMM_SEQ_FILE", "./pmm_seq.state"),
		ChainLastSeq:   chainSeq,
		ArchiveFile:    env("PMM_ARCHIVE_FILE", "./pmm_quotes.jsonl"),
		RateLimitQPS:   rlQPS,
	}, nil
}
