package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FacebookProvider struct {
	apiKey string
	host   string
	client *http.Client
}

func NewFacebookProvider(apiKey, host string) *FacebookProvider {
	return &FacebookProvider{
		apiKey: apiKey,
		host:   host,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *FacebookProvider) Name() string { return "facebook" }

func (p *FacebookProvider) Fetch(ctx context.Context, kw KeywordContext) ([]NormalizedMention, error) {
	if p.apiKey == "" {
		return p.demoMentions(kw), nil
	}

	endpoint := fmt.Sprintf("https://%s/search/posts?query=%s", p.host, url.QueryEscape(kw.Keyword))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-rapidapi-key", p.apiKey)
	req.Header.Set("x-rapidapi-host", p.host)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("facebook api status %d: %s", resp.StatusCode, string(body))
	}

	return parseFacebookResponse(body, kw)
}

func parseFacebookResponse(body []byte, kw KeywordContext) ([]NormalizedMention, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	items, _ := payload["data"].([]any)
	if len(items) == 0 {
		items, _ = payload["posts"].([]any)
	}

	result := make([]NormalizedMention, 0, len(items))
	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		content := firstString(m, "message", "text", "content", "caption")
		if content == "" {
			continue
		}
		externalID := firstString(m, "post_id", "id")
		if externalID == "" {
			externalID = fmt.Sprintf("%s-fb-%d", kw.KeywordID, i)
		}
		raw, _ := json.Marshal(m)
		result = append(result, NormalizedMention{
			BrandID:     kw.BrandID,
			BrandName:   kw.BrandName,
			KeywordID:   kw.KeywordID,
			Keyword:     kw.Keyword,
			Source:      "facebook",
			ExternalID:  externalID,
			URL:         firstString(m, "url", "permalink", "link"),
			AuthorName:  firstString(m, "author", "username", "from"),
			Content:     content,
			PublishedAt: time.Now().UTC(),
			RawPayload:  raw,
		})
	}
	return result, nil
}

func (p *FacebookProvider) demoMentions(kw KeywordContext) []NormalizedMention {
	raw, _ := json.Marshal(map[string]string{"mode": "demo"})
	return []NormalizedMention{{
		BrandID: kw.BrandID, BrandName: kw.BrandName, KeywordID: kw.KeywordID, Keyword: kw.Keyword,
		Source: "facebook", ExternalID: "demo-fb-" + kw.KeywordID,
		URL: "https://facebook.com/demo/post", AuthorName: "Demo User",
		Content: fmt.Sprintf("[Demo Facebook] Co mention ve \"%s\" cho %s", kw.Keyword, kw.BrandName),
		PublishedAt: time.Now().UTC(), RawPayload: raw,
	}}
}

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key].(string); ok && strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
