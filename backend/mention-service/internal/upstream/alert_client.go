package upstream

import "sme-social-listening/mention-service/internal/domain"

type AlertClient struct {
	baseURL string
}

type alertRequest struct {
	BrandID    string `json:"brandId"`
	Keyword    string `json:"keyword"`
	Source     string `json:"source"`
	Sentiment  string `json:"sentiment"`
	Content    string `json:"content"`
	URL        string `json:"url"`
	AlertType  string `json:"alertType"`
	AuthorName string `json:"authorName"`
}

func NewAlertClient(baseURL string) *AlertClient {
	return &AlertClient{baseURL: baseURL}
}

func (c *AlertClient) SendMentionAlert(mention domain.Mention, alertType string) error {
	_, err := postJSON[envelope[any]](
		c.baseURL+"/alerts",
		alertRequest{
			BrandID:    mention.BrandID,
			Keyword:    mention.Keyword,
			Source:     mention.Source,
			Sentiment:  mention.Sentiment,
			Content:    mention.Content,
			URL:        mention.URL,
			AlertType:  alertType,
			AuthorName: mention.AuthorName,
		},
	)
	return err
}
