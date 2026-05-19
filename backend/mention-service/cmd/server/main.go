package main

import (
	"context"
	"log"

	"sme-social-listening/mention-service/internal/config"
	"sme-social-listening/mention-service/internal/service"
	"sme-social-listening/mention-service/internal/store"
	httptransport "sme-social-listening/mention-service/internal/transport/http"
	"sme-social-listening/mention-service/internal/upstream"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	mentionStore, err := store.NewPostgresStore(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer mentionStore.Close()

	sentimentClient := upstream.NewSentimentClient(cfg.SentimentServiceURL)
	alertClient := upstream.NewAlertClient(cfg.AlertServiceURL)
	mentionService := service.NewMentionService(mentionStore, sentimentClient, alertClient)
	router := httptransport.NewRouter(mentionService, cfg)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
