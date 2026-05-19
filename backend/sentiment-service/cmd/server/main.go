package main

import (
	"log"

	"sme-social-listening/sentiment-service/internal/config"
	httptransport "sme-social-listening/sentiment-service/internal/transport/http"
)

func main() {
	cfg := config.Load()
	router := httptransport.NewRouter(cfg)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
