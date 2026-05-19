package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"sme-social-listening/mention-service/internal/domain"
)

func BuildDedupKey(source, externalID, url, content string, publishedAt time.Time) string {
	source = strings.TrimSpace(source)
	externalID = strings.TrimSpace(externalID)
	url = strings.TrimSpace(url)

	if externalID != "" {
		return source + ":" + externalID
	}
	if url != "" {
		sum := sha256.Sum256([]byte(url))
		return source + ":url:" + hex.EncodeToString(sum[:])
	}
	payload := content + publishedAt.UTC().Format(time.RFC3339)
	sum := sha256.Sum256([]byte(payload))
	return source + ":content:" + hex.EncodeToString(sum[:])
}

func prepareMention(req domain.CreateMentionRequest, sentiment domain.SentimentResult) domain.Mention {
	publishedAt := time.Now().UTC()
	if req.PublishedAt != nil {
		publishedAt = req.PublishedAt.UTC()
	}

	source := fallback(strings.TrimSpace(req.Source), "manual")
	mention := domain.Mention{
		BrandID:             strings.TrimSpace(req.BrandID),
		KeywordID:           strings.TrimSpace(req.KeywordID),
		Keyword:             strings.TrimSpace(req.Keyword),
		Source:              source,
		ExternalID:          strings.TrimSpace(req.ExternalID),
		URL:                 strings.TrimSpace(req.URL),
		AuthorName:          strings.TrimSpace(req.AuthorName),
		Content:             strings.TrimSpace(req.Content),
		PublishedAt:         publishedAt,
		EngagementCount:     req.EngagementCount,
		Sentiment:           sentiment.Sentiment,
		SentimentConfidence: sentiment.Confidence,
		SentimentReason:     sentiment.Reason,
		RawPayload:          req.RawPayload,
		CreatedAt:           time.Now().UTC(),
	}
	mention.DedupKey = BuildDedupKey(mention.Source, mention.ExternalID, mention.URL, mention.Content, mention.PublishedAt)
	return mention
}

func fallback(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func alertTypeFor(mention domain.Mention) string {
	if mention.Sentiment == "negative" {
		return "negative"
	}
	return "new"
}
