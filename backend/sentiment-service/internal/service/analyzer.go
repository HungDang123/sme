package service

import (
	"context"
	"strings"
	"time"

	"sme-social-listening/sentiment-service/internal/domain"
)

type Analyzer struct {
	rules  *RuleAnalyzer
	gemini *GeminiClient
}

func NewAnalyzer(geminiAPIKey string) *Analyzer {
	return &Analyzer{
		rules:  NewRuleAnalyzer(),
		gemini: NewGeminiClient(geminiAPIKey),
	}
}

func (a *Analyzer) Analyze(req domain.AnalyzeRequest) domain.AnalyzeResponse {
	ruleResult := a.rules.Analyze(req.Content)
	if ruleResult.Sentiment != "neutral" && ruleResult.Confidence >= 0.8 {
		return ruleResult
	}

	if a.gemini.apiKey != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		geminiResult, err := a.gemini.Analyze(ctx, req)
		if err == nil {
			return geminiResult
		}
		ruleResult.Reason = "gemini unavailable: " + err.Error()
	}

	if ruleResult.Sentiment != "neutral" {
		return ruleResult
	}
	return domain.AnalyzeResponse{Sentiment: "neutral", Confidence: 0.6, Reason: ruleResult.Reason}
}

type RuleAnalyzer struct{}

func NewRuleAnalyzer() *RuleAnalyzer {
	return &RuleAnalyzer{}
}

func (a *RuleAnalyzer) Analyze(content string) domain.AnalyzeResponse {
	text := strings.ToLower(content)

	for _, keyword := range []string{"lua dao", "te", "kem", "that vong", "khong giao", "bao hanh kem", "phan nan", "khieu nai"} {
		if strings.Contains(text, keyword) {
			return domain.AnalyzeResponse{Sentiment: "negative", Confidence: 0.85, Reason: "matched negative keyword rule"}
		}
	}

	for _, keyword := range []string{"tot", "hai long", "ung ho", "dang tien", "rat thich", "tuyet voi", "re dep"} {
		if strings.Contains(text, keyword) {
			return domain.AnalyzeResponse{Sentiment: "positive", Confidence: 0.8, Reason: "matched positive keyword rule"}
		}
	}

	return domain.AnalyzeResponse{Sentiment: "neutral", Confidence: 0.6, Reason: "no strong rule matched"}
}
