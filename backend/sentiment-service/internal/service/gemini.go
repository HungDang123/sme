package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sme-social-listening/sentiment-service/internal/domain"
)

type GeminiClient struct {
	apiKey string
	client *http.Client
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{apiKey: apiKey, client: &http.Client{Timeout: 20 * time.Second}}
}

func (g *GeminiClient) Analyze(ctx context.Context, req domain.AnalyzeRequest) (domain.AnalyzeResponse, error) {
	if g.apiKey == "" {
		return domain.AnalyzeResponse{}, fmt.Errorf("gemini api key not configured")
	}

	prompt := fmt.Sprintf(`Phan tich sentiment tieng Viet cho mention social listening.
Brand: %s
Keyword: %s
Content: %s

Tra ve JSON duy nhat voi format:
{"sentiment":"positive|neutral|negative","confidence":0.0,"reason":"...","detected_language":"vi"}`,
		req.BrandName, req.Keyword, req.Content,
	)

	body, _ := json.Marshal(map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	})

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + g.apiKey
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return domain.AnalyzeResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return domain.AnalyzeResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.AnalyzeResponse{}, err
	}
	if resp.StatusCode >= 400 {
		return domain.AnalyzeResponse{}, fmt.Errorf("gemini status %d: %s", resp.StatusCode, string(respBody))
	}

	return parseGeminiResponse(respBody)
}

func parseGeminiResponse(body []byte) (domain.AnalyzeResponse, error) {
	var payload struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return domain.AnalyzeResponse{}, err
	}
	if len(payload.Candidates) == 0 || len(payload.Candidates[0].Content.Parts) == 0 {
		return domain.AnalyzeResponse{}, fmt.Errorf("empty gemini response")
	}

	text := strings.TrimSpace(payload.Candidates[0].Content.Parts[0].Text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result domain.AnalyzeResponse
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return domain.AnalyzeResponse{}, err
	}
	result.Sentiment = normalizeSentiment(result.Sentiment)
	if result.Confidence <= 0 {
		result.Confidence = 0.7
	}
	return result, nil
}

func normalizeSentiment(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "positive", "negative", "neutral":
		return strings.ToLower(strings.TrimSpace(s))
	default:
		return "neutral"
	}
}
