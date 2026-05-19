package main

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"

	"sme-social-listening/ingest-service/internal/config"
	"sme-social-listening/ingest-service/internal/crawler"
	"sme-social-listening/ingest-service/internal/provider"
	"sme-social-listening/ingest-service/internal/ratelimit"
	httptransport "sme-social-listening/ingest-service/internal/transport/http"
	"sme-social-listening/ingest-service/internal/upstream"
	"sme-social-listening/ingest-service/internal/worker"
)

func main() {
	cfg := config.Load()

	brandClient := upstream.NewBrandClient(cfg.BrandServiceURL)
	mentionClient := upstream.NewMentionClient(cfg.MentionServiceURL, cfg.IngestToken)

	providers := []provider.SourceProvider{
		provider.NewFacebookProvider(cfg.RapidAPIKey, cfg.RapidAPIFacebookHost),
		provider.NewYouTubeProvider(cfg.RapidAPIKey, cfg.RapidAPIYouTubeHost),
	}

	crawl := crawler.New(
		brandClient,
		mentionClient,
		providers,
		ratelimit.New(1, 2),
		worker.New(cfg.WorkerConcurrency),
	)

	c := cron.New()
	if _, err := c.AddFunc(cfg.CrawlCron, func() {
		if err := crawl.Run(context.Background()); err != nil {
			log.Printf("scheduled crawl error: %v", err)
		}
	}); err != nil {
		log.Fatal(err)
	}
	c.Start()

	router := httptransport.NewRouter(crawl)
	log.Printf("ingest-service listening on :%s cron=%s", cfg.Port, cfg.CrawlCron)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
