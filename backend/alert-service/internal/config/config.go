package config

import "os"

type Config struct {
	Port            string
	DatabaseURL     string
	TelegramBotToken string
	BrandServiceURL string
	AlertBatchSeconds int
	AlertBatchMax     int
}

func Load() Config {
	return Config{
		Port:              env("PORT", "8084"),
		DatabaseURL:       env("DATABASE_URL", "postgres://sme:sme_secret@localhost:5432/sme_listening?sslmode=disable"),
		TelegramBotToken:  env("TELEGRAM_BOT_TOKEN", ""),
		BrandServiceURL:   env("BRAND_SERVICE_URL", "http://localhost:8081"),
		AlertBatchSeconds: envInt("ALERT_BATCH_SECONDS", 300),
		AlertBatchMax:     envInt("ALERT_BATCH_MAX", 10),
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
	n := 0
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
