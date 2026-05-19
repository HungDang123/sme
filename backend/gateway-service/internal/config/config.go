package config

import "os"

type Config struct {
	Port                string
	BrandServiceURL     string
	MentionServiceURL   string
	SentimentServiceURL string
	AlertServiceURL     string
	IngestServiceURL    string
}

func Load() Config {
	return Config{
		Port:                env("PORT", "8080"),
		BrandServiceURL:     env("BRAND_SERVICE_URL", "http://localhost:8081"),
		MentionServiceURL:   env("MENTION_SERVICE_URL", "http://localhost:8082"),
		SentimentServiceURL: env("SENTIMENT_SERVICE_URL", "http://localhost:8083"),
		AlertServiceURL:     env("ALERT_SERVICE_URL", "http://localhost:8084"),
		IngestServiceURL:    env("INGEST_SERVICE_URL", "http://localhost:8085"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
