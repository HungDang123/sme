package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sme-social-listening/alert-service/internal/domain"
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

func (s *PostgresStore) Create(ctx context.Context, alert domain.Alert) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO alerts (id, brand_id, keyword, source, sentiment, content, url, alert_type, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NULLIF($7,''),$8,$9)`,
		alert.ID, alert.BrandID, alert.Keyword, alert.Source, alert.Sentiment,
		alert.Content, alert.URL, alert.AlertType, alert.CreatedAt,
	)
	return err
}

func (s *PostgresStore) List(ctx context.Context, limit int) ([]domain.Alert, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.pool.Query(ctx,
		`SELECT id, brand_id, keyword, source, sentiment, content, COALESCE(url,''), alert_type, created_at
		 FROM alerts ORDER BY created_at DESC LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Alert
	for rows.Next() {
		var a domain.Alert
		if err := rows.Scan(&a.ID, &a.BrandID, &a.Keyword, &a.Source, &a.Sentiment,
			&a.Content, &a.URL, &a.AlertType, &a.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}
