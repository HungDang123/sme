package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sme-social-listening/brand-service/internal/domain"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) CreateBrand(ctx context.Context, brand domain.Brand) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO brands (id, name, telegram_chat_id, created_at) VALUES ($1, $2, NULLIF($3,''), $4)`,
		brand.ID, brand.Name, brand.TelegramChatID, brand.CreatedAt,
	)
	return err
}

func (s *PostgresStore) GetBrand(ctx context.Context, id string) (domain.Brand, error) {
	var brand domain.Brand
	var chatID *string
	err := s.pool.QueryRow(ctx,
		`SELECT id, name, telegram_chat_id, created_at FROM brands WHERE id = $1`, id,
	).Scan(&brand.ID, &brand.Name, &chatID, &brand.CreatedAt)
	if err != nil {
		return domain.Brand{}, err
	}
	if chatID != nil {
		brand.TelegramChatID = *chatID
	}
	keywords, err := s.ListKeywordsByBrand(ctx, id)
	if err != nil {
		return domain.Brand{}, err
	}
	brand.Keywords = keywordValues(keywords)
	return brand, nil
}

func (s *PostgresStore) ListBrands(ctx context.Context) ([]domain.Brand, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, telegram_chat_id, created_at FROM brands ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var brands []domain.Brand
	for rows.Next() {
		var brand domain.Brand
		var chatID *string
		if err := rows.Scan(&brand.ID, &brand.Name, &chatID, &brand.CreatedAt); err != nil {
			return nil, err
		}
		if chatID != nil {
			brand.TelegramChatID = *chatID
		}
		keywords, err := s.ListKeywordsByBrand(ctx, brand.ID)
		if err != nil {
			return nil, err
		}
		brand.Keywords = keywordValues(keywords)
		brands = append(brands, brand)
	}
	return brands, rows.Err()
}

func (s *PostgresStore) CreateKeyword(ctx context.Context, keyword domain.Keyword) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO keywords (id, brand_id, value, is_active, created_at) VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (brand_id, value) DO UPDATE SET is_active = true`,
		keyword.ID, keyword.BrandID, keyword.Value, keyword.IsActive, keyword.CreatedAt,
	)
	return err
}

func (s *PostgresStore) ListKeywordsByBrand(ctx context.Context, brandID string) ([]domain.Keyword, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, brand_id, value, is_active, created_at FROM keywords
		 WHERE brand_id = $1 AND is_active = true ORDER BY created_at ASC`, brandID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Keyword
	for rows.Next() {
		var k domain.Keyword
		if err := rows.Scan(&k.ID, &k.BrandID, &k.Value, &k.IsActive, &k.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, k)
	}
	return result, rows.Err()
}

func (s *PostgresStore) ListActiveKeywords(ctx context.Context) ([]domain.ActiveKeyword, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT k.id, k.brand_id, k.value, k.is_active, k.created_at, b.name
		 FROM keywords k
		 JOIN brands b ON b.id = k.brand_id
		 WHERE k.is_active = true
		 ORDER BY b.name ASC, k.value ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.ActiveKeyword
	for rows.Next() {
		var ak domain.ActiveKeyword
		if err := rows.Scan(&ak.ID, &ak.BrandID, &ak.Value, &ak.IsActive, &ak.CreatedAt, &ak.BrandName); err != nil {
			return nil, err
		}
		result = append(result, ak)
	}
	return result, rows.Err()
}

func (s *PostgresStore) UpdateTelegramChatID(ctx context.Context, brandID, chatID string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE brands SET telegram_chat_id = NULLIF($2,'') WHERE id = $1`, brandID, chatID,
	)
	return err
}

func keywordValues(keywords []domain.Keyword) []string {
	values := make([]string, 0, len(keywords))
	for _, k := range keywords {
		values = append(values, k.Value)
	}
	return values
}
