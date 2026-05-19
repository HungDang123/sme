package main

import (
	"context"
	"log"

	"sme-social-listening/brand-service/internal/config"
	"sme-social-listening/brand-service/internal/service"
	"sme-social-listening/brand-service/internal/store"
	httptransport "sme-social-listening/brand-service/internal/transport/http"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pgStore, err := store.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgStore.Close()

	brandService := service.NewBrandService(pgStore)
	router := httptransport.NewRouter(brandService)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
