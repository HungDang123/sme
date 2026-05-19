package provider

import (
	"context"
	"encoding/json"
	"time"
)

type KeywordContext struct {
	BrandID   string
	BrandName string
	KeywordID string
	Keyword   string
}

type NormalizedMention struct {
	BrandID         string
	BrandName       string
	KeywordID       string
	Keyword         string
	Source          string
	ExternalID      string
	URL             string
	AuthorName      string
	Content         string
	PublishedAt     time.Time
	EngagementCount int
	RawPayload      json.RawMessage
}

type SourceProvider interface {
	Name() string
	Fetch(ctx context.Context, kw KeywordContext) ([]NormalizedMention, error)
}
