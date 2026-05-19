package upstream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type envelope[T any] struct {
	Data T `json:"data"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func postJSON[T any](url string, payload any) (T, error) {
	var result T
	body, err := json.Marshal(payload)
	if err != nil {
		return result, err
	}

	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, fmt.Errorf("http %d from %s", resp.StatusCode, url)
	}

	return result, json.NewDecoder(resp.Body).Decode(&result)
}
