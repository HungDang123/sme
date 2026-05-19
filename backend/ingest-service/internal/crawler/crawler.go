package crawler

import (
	"context"
	"log"

	"sme-social-listening/ingest-service/internal/provider"
	"sme-social-listening/ingest-service/internal/ratelimit"
	"sme-social-listening/ingest-service/internal/upstream"
	"sme-social-listening/ingest-service/internal/worker"
)

type Crawler struct {
	brandClient   *upstream.BrandClient
	mentionClient *upstream.MentionClient
	providers     []provider.SourceProvider
	limiter       *ratelimit.Limiter
	pool          *worker.Pool
}

func New(
	brandClient *upstream.BrandClient,
	mentionClient *upstream.MentionClient,
	providers []provider.SourceProvider,
	limiter *ratelimit.Limiter,
	pool *worker.Pool,
) *Crawler {
	return &Crawler{
		brandClient:   brandClient,
		mentionClient: mentionClient,
		providers:     providers,
		limiter:       limiter,
		pool:          pool,
	}
}

func (c *Crawler) Run(ctx context.Context) error {
	keywords, err := c.brandClient.ListActiveKeywords()
	if err != nil {
		return err
	}
	if len(keywords) == 0 {
		log.Println("crawl: no active keywords")
		return nil
	}

	jobs := make([]worker.Job, 0, len(keywords))
	for _, kw := range keywords {
		kwCopy := kw
		jobs = append(jobs, func(ctx context.Context) error {
			return c.crawlKeyword(ctx, kwCopy)
		})
	}

	errs := c.pool.Run(ctx, jobs)
	for _, err := range errs {
		if err != nil {
			log.Printf("crawl keyword error: %v", err)
		}
	}
	return nil
}

func (c *Crawler) crawlKeyword(ctx context.Context, kw upstream.ActiveKeyword) error {
	kwCtx := provider.KeywordContext{
		BrandID: kw.BrandID, BrandName: kw.BrandName, KeywordID: kw.ID, Keyword: kw.Value,
	}
	for _, p := range c.providers {
		if err := c.limiter.Wait(ctx, p.Name()); err != nil {
			return err
		}
		mentions, err := p.Fetch(ctx, kwCtx)
		if err != nil {
			log.Printf("provider %s keyword %s error: %v", p.Name(), kw.Value, err)
			continue
		}
		for _, m := range mentions {
			if err := c.mentionClient.Create(m); err != nil {
				log.Printf("save mention %s/%s error: %v", p.Name(), kw.Value, err)
			}
		}
	}
	return nil
}
