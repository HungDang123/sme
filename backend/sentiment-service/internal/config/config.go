package config

import "os"

type Config struct {
	Port         string
	GeminiAPIKey string
}

func Load() Config {
	return Config{
		Port:         env("PORT", "8083"),
		GeminiAPIKey: env("GEMINI_API_KEY", ""),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
