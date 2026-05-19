package domain

type AnalyzeRequest struct {
	BrandName string `json:"brandName"`
	Keyword   string `json:"keyword"`
	Content   string `json:"content"`
}

type AnalyzeResponse struct {
	Sentiment  string  `json:"sentiment"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}
