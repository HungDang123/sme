package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ActiveKeyword struct {
	ID        string `json:"id"`
	BrandID   string `json:"brandId"`
	Value     string `json:"value"`
	BrandName string `json:"brandName"`
}

type BrandClient struct {
	baseURL string
	client  *http.Client
}

func NewBrandClient(baseURL string) *BrandClient {
	return &BrandClient{baseURL: baseURL, client: &http.Client{}}
}

func (c *BrandClient) ListActiveKeywords() ([]ActiveKeyword, error) {
	resp, err := c.client.Get(c.baseURL + "/keywords/active")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("brand service status %d: %s", resp.StatusCode, string(body))
	}
	var out struct {
		Data []ActiveKeyword `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
