package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Brand struct {
	ID             string `json:"id"`
	TelegramChatID string `json:"telegramChatId"`
}

type BrandClient struct {
	baseURL string
	client  *http.Client
}

func NewBrandClient(baseURL string) *BrandClient {
	return &BrandClient{baseURL: baseURL, client: &http.Client{}}
}

func (c *BrandClient) GetBrand(id string) (Brand, error) {
	resp, err := c.client.Get(c.baseURL + "/brands/" + id)
	if err != nil {
		return Brand{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Brand{}, err
	}
	if resp.StatusCode == 404 {
		return Brand{}, fmt.Errorf("brand not found")
	}
	if resp.StatusCode >= 400 {
		return Brand{}, fmt.Errorf("brand service status %d: %s", resp.StatusCode, string(body))
	}
	var out struct {
		Data Brand `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return Brand{}, err
	}
	return out.Data, nil
}
