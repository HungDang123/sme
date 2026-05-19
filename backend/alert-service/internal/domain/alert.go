package domain

import "time"

type CreateAlertRequest struct {
	BrandID    string `json:"brandId"`
	Keyword    string `json:"keyword"`
	Source     string `json:"source"`
	Sentiment  string `json:"sentiment"`
	Content    string `json:"content"`
	URL        string `json:"url"`
	AlertType  string `json:"alertType"`
	AuthorName string `json:"authorName"`
}

type Alert struct {
	CreateAlertRequest
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
