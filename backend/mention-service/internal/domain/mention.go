package domain

import (
	"encoding/json"
	"time"
)

type Mention struct {
	ID                  string          `json:"id"`
	BrandID             string          `json:"brandId"`
	KeywordID           string          `json:"keywordId,omitempty"`
	Keyword             string          `json:"keyword"`
	Source              string          `json:"source"`
	ExternalID          string          `json:"externalId,omitempty"`
	URL                 string          `json:"url"`
	AuthorName          string          `json:"authorName"`
	Content             string          `json:"content"`
	PublishedAt         time.Time       `json:"publishedAt"`
	EngagementCount     int             `json:"engagementCount"`
	Sentiment           string          `json:"sentiment"`
	SentimentConfidence float64         `json:"sentimentConfidence"`
	SentimentReason     string          `json:"sentimentReason"`
	RawPayload          json.RawMessage `json:"rawPayload,omitempty"`
	DedupKey            string          `json:"-"`
	CreatedAt           time.Time       `json:"createdAt"`
}

type CreateMentionRequest struct {
	BrandID         string          `json:"brandId"`
	BrandName       string          `json:"brandName"`
	KeywordID       string          `json:"keywordId"`
	Keyword         string          `json:"keyword"`
	Source          string          `json:"source"`
	ExternalID      string          `json:"externalId"`
	URL             string          `json:"url"`
	AuthorName      string          `json:"authorName"`
	Content         string          `json:"content"`
	PublishedAt     *time.Time      `json:"publishedAt"`
	EngagementCount int             `json:"engagementCount"`
	RawPayload      json.RawMessage `json:"rawPayload"`
}

type MentionFilter struct {
	BrandID   string
	Keyword   string
	Source    string
	Sentiment string
	From      *time.Time
	To        *time.Time
	Limit     int
}

type SentimentResult struct {
	Sentiment  string
	Confidence float64
	Reason     string
}

type CreateResult struct {
	Mention Mention
	Created bool
}
