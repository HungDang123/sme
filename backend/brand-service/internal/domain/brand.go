package domain

import "time"

type Brand struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Keywords         []string  `json:"keywords"`
	TelegramChatID   string    `json:"telegramChatId,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
}

type Keyword struct {
	ID        string    `json:"id"`
	BrandID   string    `json:"brandId"`
	Value     string    `json:"value"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type ActiveKeyword struct {
	Keyword
	BrandName string `json:"brandName"`
}
