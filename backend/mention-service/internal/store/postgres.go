package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sme-social-listening/mention-service/internal/domain"
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

func (s *PostgresStore) Create(ctx context.Context, mention domain.Mention) (domain.Mention, bool, error) {
	row := s.pool.QueryRow(ctx,
		`INSERT INTO mentions (
			id, brand_id, keyword_id, keyword, source, external_id, url, author_name,
			content, published_at, engagement_count, sentiment, sentiment_confidence,
			sentiment_reason, raw_payload, dedup_key, created_at
		) VALUES (
			$1,$2,NULLIF($3,''),$4,$5,NULLIF($6,''),NULLIF($7,''),NULLIF($8,''),
			$9,$10,$11,$12,$13,NULLIF($14,''),$15,$16,$17
		)
		ON CONFLICT (dedup_key) DO NOTHING
		RETURNING id`,
		mention.ID, mention.BrandID, mention.KeywordID, mention.Keyword, mention.Source,
		mention.ExternalID, mention.URL, mention.AuthorName, mention.Content, mention.PublishedAt,
		mention.EngagementCount, mention.Sentiment, mention.SentimentConfidence, mention.SentimentReason,
		nullRaw(mention.RawPayload), mention.DedupKey, mention.CreatedAt,
	)

	var insertedID string
	err := row.Scan(&insertedID)
	if err == pgx.ErrNoRows {
		existing, err := s.GetByDedupKey(ctx, mention.DedupKey)
		if err != nil {
			return domain.Mention{}, false, err
		}
		return existing, false, nil
	}
	if err != nil {
		return domain.Mention{}, false, err
	}
	mention.ID = insertedID
	return mention, true, nil
}

func (s *PostgresStore) GetByDedupKey(ctx context.Context, dedupKey string) (domain.Mention, error) {
	return s.scanOne(ctx,
		`SELECT id, brand_id, COALESCE(keyword_id,''), keyword, source, COALESCE(external_id,''),
			COALESCE(url,''), COALESCE(author_name,''), content, published_at, engagement_count,
			sentiment, COALESCE(sentiment_confidence,0), COALESCE(sentiment_reason,''),
			raw_payload, dedup_key, created_at
		 FROM mentions WHERE dedup_key = $1`, dedupKey,
	)
}

func (s *PostgresStore) List(ctx context.Context, filter domain.MentionFilter) ([]domain.Mention, error) {
	query := `SELECT id, brand_id, COALESCE(keyword_id,''), keyword, source, COALESCE(external_id,''),
		COALESCE(url,''), COALESCE(author_name,''), content, published_at, engagement_count,
		sentiment, COALESCE(sentiment_confidence,0), COALESCE(sentiment_reason,''),
		raw_payload, dedup_key, created_at
		FROM mentions WHERE 1=1`
	args := []any{}
	argN := 1

	if filter.BrandID != "" {
		query += fmt.Sprintf(" AND brand_id = $%d", argN)
		args = append(args, filter.BrandID)
		argN++
	}
	if filter.Keyword != "" {
		query += fmt.Sprintf(" AND keyword = $%d", argN)
		args = append(args, filter.Keyword)
		argN++
	}
	if filter.Source != "" {
		query += fmt.Sprintf(" AND source = $%d", argN)
		args = append(args, filter.Source)
		argN++
	}
	if filter.Sentiment != "" {
		query += fmt.Sprintf(" AND sentiment = $%d", argN)
		args = append(args, filter.Sentiment)
		argN++
	}
	if filter.From != nil {
		query += fmt.Sprintf(" AND published_at >= $%d", argN)
		args = append(args, *filter.From)
		argN++
	}
	if filter.To != nil {
		query += fmt.Sprintf(" AND published_at <= $%d", argN)
		args = append(args, *filter.To)
		argN++
	}

	query += " ORDER BY published_at DESC"
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	query += fmt.Sprintf(" LIMIT $%d", argN)
	args = append(args, limit)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Mention
	for rows.Next() {
		m, err := scanMentionRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (s *PostgresStore) scanOne(ctx context.Context, query string, args ...any) (domain.Mention, error) {
	row := s.pool.QueryRow(ctx, query, args...)
	return scanMentionRow(row)
}

type scannable interface {
	Scan(dest ...any) error
}

func scanMentionRow(row scannable) (domain.Mention, error) {
	var m domain.Mention
	var raw []byte
	err := row.Scan(
		&m.ID, &m.BrandID, &m.KeywordID, &m.Keyword, &m.Source, &m.ExternalID,
		&m.URL, &m.AuthorName, &m.Content, &m.PublishedAt, &m.EngagementCount,
		&m.Sentiment, &m.SentimentConfidence, &m.SentimentReason,
		&raw, &m.DedupKey, &m.CreatedAt,
	)
	if err != nil {
		return domain.Mention{}, err
	}
	if len(raw) > 0 {
		m.RawPayload = raw
	}
	return m, nil
}

func nullRaw(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	return raw
}
