package config

import "os"

type Config struct {
	Port                string
	BrandServiceURL     string
	MentionServiceURL   string
	IngestToken         string
	RapidAPIKey         string
	RapidAPIFacebookHost string
	RapidAPIYouTubeHost  string
	CrawlCron           string
	WorkerConcurrency   int
}

func Load() Config {
	return Config{
		Port:                 env("PORT", "8085"),
		BrandServiceURL:      env("BRAND_SERVICE_URL", "http://localhost:8081"),
		MentionServiceURL:    env("MENTION_SERVICE_URL", "http://localhost:8082"),
		IngestToken:          env("INGEST_TOKEN", "dev-ingest-token"),
		RapidAPIKey:          env("RAPIDAPI_KEY", ""),
		RapidAPIFacebookHost: env("RAPIDAPI_FACEBOOK_HOST", "facebook-scraper3.p.rapidapi.com"),
		RapidAPIYouTubeHost:  env("RAPIDAPI_YOUTUBE_HOST", "youtube-v31.p.rapidapi.com"),
		CrawlCron:            env("CRAWL_CRON", "*/30 * * * *"),
		WorkerConcurrency:    envInt("WORKER_CONCURRENCY", 5),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var n int
	for _, c := range v {
		if c < '0' || c > '9' {
			return fallback
		}
		n = n*10 + int(c-'0')
	}
	if n <= 0 {
		return fallback
	}
	return n
}
