package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"sme-social-listening/alert-service/internal/domain"
	"sme-social-listening/alert-service/internal/store"
	"sme-social-listening/alert-service/internal/telegram"
	"sme-social-listening/alert-service/internal/upstream"
)

type AlertService struct {
	store       *store.PostgresStore
	batcher     *telegram.Batcher
	brandClient *upstream.BrandClient
}

func NewAlertService(
	store *store.PostgresStore,
	batcher *telegram.Batcher,
	brandClient *upstream.BrandClient,
) *AlertService {
	return &AlertService{store: store, batcher: batcher, brandClient: brandClient}
}

func (s *AlertService) List(ctx context.Context) ([]domain.Alert, error) {
	return s.store.List(ctx, 100)
}

func (s *AlertService) Create(ctx context.Context, req domain.CreateAlertRequest) (domain.Alert, error) {
	if req.AlertType == "" {
		if req.Sentiment == "negative" {
			req.AlertType = "negative"
		} else {
			req.AlertType = "new"
		}
	}

	alert := domain.Alert{
		CreateAlertRequest: req,
		ID:                 newID(),
		CreatedAt:          time.Now().UTC(),
	}
	if err := s.store.Create(ctx, alert); err != nil {
		return domain.Alert{}, err
	}

	s.batcher.Add(req)
	return alert, nil
}

func (s *AlertService) LookupChatID(brandID string) string {
	brand, err := s.brandClient.GetBrand(brandID)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(brand.TelegramChatID)
}

func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b[:])
}
