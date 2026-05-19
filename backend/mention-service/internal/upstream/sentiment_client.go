package upstream

import "sme-social-listening/mention-service/internal/domain"

type SentimentClient struct {
	baseURL string
}

type sentimentRequest struct {
	BrandName string `json:"brandName"`
	Keyword   string `json:"keyword"`
	Content   string `json:"content"`
}

type sentimentResponse struct {
	Sentiment  string  `json:"sentiment"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

func NewSentimentClient(baseURL string) *SentimentClient {
	return &SentimentClient{baseURL: baseURL}
}

func (c *SentimentClient) Analyze(req domain.CreateMentionRequest) (domain.SentimentResult, error) {
	resp, err := postJSON[envelope[sentimentResponse]](
		c.baseURL+"/sentiment/analyze",
		sentimentRequest{BrandName: req.BrandName, Keyword: req.Keyword, Content: req.Content},
	)
	if err != nil {
		return domain.SentimentResult{}, err
	}

	return domain.SentimentResult{
		Sentiment:  resp.Data.Sentiment,
		Confidence: resp.Data.Confidence,
		Reason:     resp.Data.Reason,
	}, nil
}
