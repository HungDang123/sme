package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"sme-social-listening/brand-service/internal/domain"
	"sme-social-listening/brand-service/internal/store"
)

var ErrBrandNotFound = errors.New("brand not found")

type BrandService struct {
	store *store.PostgresStore
}

func NewBrandService(store *store.PostgresStore) *BrandService {
	return &BrandService{store: store}
}

func (s *BrandService) List(ctx context.Context) ([]domain.Brand, error) {
	return s.store.ListBrands(ctx)
}

func (s *BrandService) Create(ctx context.Context, name string, keywords []string) (domain.Brand, error) {
	brand := domain.Brand{
		ID:        newID(),
		Name:      strings.TrimSpace(name),
		CreatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateBrand(ctx, brand); err != nil {
		return domain.Brand{}, err
	}
	for _, kw := range cleanKeywords(keywords) {
		if _, err := s.AddKeyword(ctx, brand.ID, kw); err != nil {
			return domain.Brand{}, err
		}
	}
	return s.store.GetBrand(ctx, brand.ID)
}

func (s *BrandService) AddKeyword(ctx context.Context, brandID, value string) (domain.Keyword, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return domain.Keyword{}, errors.New("keyword value is required")
	}
	if _, err := s.GetBrand(ctx, brandID); err != nil {
		return domain.Keyword{}, err
	}
	keyword := domain.Keyword{
		ID:        newID(),
		BrandID:   brandID,
		Value:     value,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateKeyword(ctx, keyword); err != nil {
		return domain.Keyword{}, err
	}
	return keyword, nil
}

func (s *BrandService) ListKeywords(ctx context.Context, brandID string) ([]domain.Keyword, error) {
	if _, err := s.GetBrand(ctx, brandID); err != nil {
		return nil, err
	}
	return s.store.ListKeywordsByBrand(ctx, brandID)
}

func (s *BrandService) ListActiveKeywords(ctx context.Context) ([]domain.ActiveKeyword, error) {
	return s.store.ListActiveKeywords(ctx)
}

func (s *BrandService) UpdateTelegramChatID(ctx context.Context, brandID, chatID string) error {
	if _, err := s.GetBrand(ctx, brandID); err != nil {
		return err
	}
	return s.store.UpdateTelegramChatID(ctx, brandID, strings.TrimSpace(chatID))
}

func (s *BrandService) GetBrand(ctx context.Context, brandID string) (domain.Brand, error) {
	brand, err := s.store.GetBrand(ctx, brandID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Brand{}, ErrBrandNotFound
		}
		return domain.Brand{}, err
	}
	return brand, nil
}

func cleanKeywords(input []string) []string {
	result := make([]string, 0, len(input))
	seen := map[string]bool{}
	for _, keyword := range input {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" || seen[keyword] {
			continue
		}
		seen[keyword] = true
		result = append(result, keyword)
	}
	return result
}

func newID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(b[:])
}
