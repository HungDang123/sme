package upstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sme-social-listening/ingest-service/internal/provider"
)

type MentionClient struct {
	baseURL     string
	ingestToken string
	client      *http.Client
}

func NewMentionClient(baseURL, ingestToken string) *MentionClient {
	return &MentionClient{
		baseURL:     baseURL,
		ingestToken: ingestToken,
		client:      &http.Client{Timeout: 15 * time.Second},
	}
}

type mentionPayload struct {
	BrandID         string          `json:"brandId"`
	BrandName       string          `json:"brandName"`
	KeywordID       string          `json:"keywordId"`
	Keyword         string          `json:"keyword"`
	Source          string          `json:"source"`
	ExternalID      string          `json:"externalId"`
	URL             string          `json:"url"`
	AuthorName      string          `json:"authorName"`
	Content         string          `json:"content"`
	PublishedAt     time.Time       `json:"publishedAt"`
	EngagementCount int             `json:"engagementCount"`
	RawPayload      json.RawMessage `json:"rawPayload"`
}

func (c *MentionClient) Create(m provider.NormalizedMention) error {
	payload := mentionPayload{
		BrandID: m.BrandID, BrandName: m.BrandName, KeywordID: m.KeywordID, Keyword: m.Keyword,
		Source: m.Source, ExternalID: m.ExternalID, URL: m.URL, AuthorName: m.AuthorName,
		Content: m.Content, PublishedAt: m.PublishedAt, EngagementCount: m.EngagementCount,
		RawPayload: m.RawPayload,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/mentions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Ingest-Token", c.ingestToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("mention service status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
