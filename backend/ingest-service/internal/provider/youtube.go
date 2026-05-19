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

type YouTubeProvider struct {
	apiKey string
	host   string
	client *http.Client
}

func NewYouTubeProvider(apiKey, host string) *YouTubeProvider {
	return &YouTubeProvider{
		apiKey: apiKey,
		host:   host,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *YouTubeProvider) Name() string { return "youtube" }

func (p *YouTubeProvider) Fetch(ctx context.Context, kw KeywordContext) ([]NormalizedMention, error) {
	if p.apiKey == "" {
		return p.demoMentions(kw), nil
	}

	endpoint := fmt.Sprintf("https://%s/search?q=%s&part=snippet", p.host, url.QueryEscape(kw.Keyword))
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
		return nil, fmt.Errorf("youtube api status %d: %s", resp.StatusCode, string(body))
	}

	return parseYouTubeResponse(body, kw)
}

func parseYouTubeResponse(body []byte, kw KeywordContext) ([]NormalizedMention, error) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	items, _ := payload["items"].([]any)
	result := make([]NormalizedMention, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id := firstString(m, "id")
		snippet, _ := m["snippet"].(map[string]any)
		if snippet == nil {
			continue
		}
		title := firstString(snippet, "title")
		desc := firstString(snippet, "description")
		content := strings.TrimSpace(title + " " + desc)
		if content == "" {
			continue
		}
		channel := firstString(snippet, "channelTitle")
		raw, _ := json.Marshal(m)
		videoURL := "https://www.youtube.com/watch?v=" + id
		result = append(result, NormalizedMention{
			BrandID: kw.BrandID, BrandName: kw.BrandName, KeywordID: kw.KeywordID, Keyword: kw.Keyword,
			Source: "youtube", ExternalID: id, URL: videoURL, AuthorName: channel, Content: content,
			PublishedAt: time.Now().UTC(), RawPayload: raw,
		})
	}
	return result, nil
}

func (p *YouTubeProvider) demoMentions(kw KeywordContext) []NormalizedMention {
	raw, _ := json.Marshal(map[string]string{"mode": "demo"})
	return []NormalizedMention{{
		BrandID: kw.BrandID, BrandName: kw.BrandName, KeywordID: kw.KeywordID, Keyword: kw.Keyword,
		Source: "youtube", ExternalID: "demo-yt-" + kw.KeywordID,
		URL: "https://youtube.com/watch?v=demo", AuthorName: "Demo Channel",
		Content: fmt.Sprintf("[Demo YouTube] Video review ve \"%s\" - %s", kw.Keyword, kw.BrandName),
		PublishedAt: time.Now().UTC(), RawPayload: raw,
	}}
}
