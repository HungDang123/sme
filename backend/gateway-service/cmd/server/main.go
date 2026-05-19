package main

import (
	"log"

	"sme-social-listening/gateway-service/internal/config"
	httptransport "sme-social-listening/gateway-service/internal/transport/http"
)

func main() {
	cfg := config.Load()
	router := httptransport.NewRouter(cfg)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
