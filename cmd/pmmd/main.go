package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/blinkmarket/pmm-server/internal/archive"
	"github.com/blinkmarket/pmm-server/internal/config"
	"github.com/blinkmarket/pmm-server/internal/pricing"
	"github.com/blinkmarket/pmm-server/internal/quote"
	"github.com/blinkmarket/pmm-server/internal/ratelimit"
	"github.com/blinkmarket/pmm-server/internal/seq"
	"github.com/blinkmarket/pmm-server/internal/sign"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if os.Getenv("PMM_CHAIN_LAST_SEQ") == "" {
		log.Printf("WARN: PMM_CHAIN_LAST_SEQ not set; seq will resync from the local file only. " +
			"Ensure this matches the on-chain last seq for this PMM or every quote will be rejected on-chain.")
	}
	signer, err := sign.NewEnvKeySigner(cfg.PrivateKeySeed)
	if err != nil {
		log.Fatalf("signer: %v", err)
	}
	store := seq.NewFileSeqStore(cfg.SeqFile)
	chain := seq.NewConfigChainSeqReader(cfg.ChainLastSeq)
	alloc, err := seq.NewSeqAllocator(store, chain, cfg.PMMAddress)
	if err != nil {
		log.Fatalf("seq allocator: %v", err)
	}
	pricer := pricing.NewStubPricer(cfg.StubPriceBps)
	jsonl := archive.NewJSONLArchiver(cfg.ArchiveFile)

	svc := quote.New(quote.Deps{
		Price:   pricer.Price,
		NextSeq: alloc.Next,
		Signer:  signer,
		Archive: func(r quote.Record) error {
			return jsonl.Record(archive.Record{
				MarketID: r.MarketID, Side: r.Side, PriceBps: r.PriceBps,
				Size: r.Size, SeqNumber: r.SeqNumber, ExpiresAt: r.ExpiresAt,
				PMM: r.PMM, Signature: r.Signature,
			})
		},
		PMM:       cfg.PMMAddress,
		QuoteTTL:  cfg.QuoteTTLMillis,
		NowMillis: func() uint64 { return uint64(time.Now().UnixMilli()) },
	})

	mux := http.NewServeMux()
	mux.Handle("/v1/quote", quote.NewHandler(svc, ratelimit.New(cfg.RateLimitQPS)))
	log.Printf("pmmd listening on %s", cfg.ListenAddr)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, mux))
}
