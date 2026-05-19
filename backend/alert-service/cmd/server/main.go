package main

import (
	"context"
	"log"
	"time"

	"sme-social-listening/alert-service/internal/config"
	"sme-social-listening/alert-service/internal/service"
	"sme-social-listening/alert-service/internal/store"
	"sme-social-listening/alert-service/internal/telegram"
	httptransport "sme-social-listening/alert-service/internal/transport/http"
	"sme-social-listening/alert-service/internal/upstream"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pgStore, err := store.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pgStore.Close()

	brandClient := upstream.NewBrandClient(cfg.BrandServiceURL)
	tgClient := telegram.NewClient(cfg.TelegramBotToken)

	alertService := service.NewAlertService(pgStore, nil, brandClient)
	batcher := telegram.NewBatcher(
		tgClient,
		alertService.LookupChatID,
		time.Duration(cfg.AlertBatchSeconds)*time.Second,
		cfg.AlertBatchMax,
	)
	alertService = service.NewAlertService(pgStore, batcher, brandClient)

	router := httptransport.NewRouter(alertService)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
