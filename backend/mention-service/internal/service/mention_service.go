package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"sme-social-listening/mention-service/internal/domain"
	"sme-social-listening/mention-service/internal/store"
)

type SentimentAnalyzer interface {
	Analyze(req domain.CreateMentionRequest) (domain.SentimentResult, error)
}

type AlertSender interface {
	SendMentionAlert(mention domain.Mention, alertType string) error
}

type MentionService struct {
	store             *store.PostgresStore
	sentimentAnalyzer SentimentAnalyzer
	alertSender       AlertSender
}

func NewMentionService(store *store.PostgresStore, sentimentAnalyzer SentimentAnalyzer, alertSender AlertSender) *MentionService {
	return &MentionService{
		store:             store,
		sentimentAnalyzer: sentimentAnalyzer,
		alertSender:       alertSender,
	}
}

func (s *MentionService) List(ctx context.Context, filter domain.MentionFilter) ([]domain.Mention, error) {
	return s.store.List(ctx, filter)
}

func (s *MentionService) Create(ctx context.Context, req domain.CreateMentionRequest) (domain.CreateResult, error) {
	sentiment, err := s.sentimentAnalyzer.Analyze(req)
	if err != nil {
		sentiment = domain.SentimentResult{
			Sentiment:  "neutral",
			Confidence: 0,
			Reason:     "sentiment service unavailable",
		}
	}

	mention := prepareMention(req, sentiment)
	mention.ID = newID()

	saved, created, err := s.store.Create(ctx, mention)
	if err != nil {
		return domain.CreateResult{}, err
	}

	if created {
		_ = s.alertSender.SendMentionAlert(saved, alertTypeFor(saved))
	}

	return domain.CreateResult{Mention: saved, Created: created}, nil
}

func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b[:])
}
