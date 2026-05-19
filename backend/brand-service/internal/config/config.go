package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port:        env("PORT", "8081"),
		DatabaseURL: env("DATABASE_URL", "postgres://sme:sme_secret@localhost:5432/sme_listening?sslmode=disable"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
