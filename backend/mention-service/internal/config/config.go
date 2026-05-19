package config

import "os"

type Config struct {
	Port                string
	DatabaseURL         string
	SentimentServiceURL string
	AlertServiceURL     string
	IngestToken         string
}

func Load() Config {
	return Config{
		Port:                env("PORT", "8082"),
		DatabaseURL:         env("DATABASE_URL", "postgres://sme:sme_secret@localhost:5432/sme_listening?sslmode=disable"),
		SentimentServiceURL: env("SENTIMENT_SERVICE_URL", "http://localhost:8083"),
		AlertServiceURL:     env("ALERT_SERVICE_URL", "http://localhost:8084"),
		IngestToken:         env("INGEST_TOKEN", "dev-ingest-token"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
