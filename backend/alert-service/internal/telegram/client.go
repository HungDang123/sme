package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	token  string
	client *http.Client
}

func NewClient(token string) *Client {
	return &Client{token: token, client: &http.Client{}}
}

func (c *Client) Enabled() bool {
	return c.token != ""
}

func (c *Client) SendMessage(chatID, text string) error {
	if !c.Enabled() || chatID == "" {
		return nil
	}

	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)
	body, _ := json.Marshal(map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	})

	resp, err := c.client.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func EscapeHTML(s string) string {
	r := url.QueryEscape(s)
	// QueryEscape is not ideal for HTML; simple replace for MVP
	replacer := map[string]string{"&": "&amp;", "<": "&lt;", ">": "&gt;"}
	out := s
	for old, new := range replacer {
		out = replaceAll(out, old, new)
	}
	_ = r
	return out
}

func replaceAll(s, old, new string) string {
	for {
		i := indexOf(s, old)
		if i < 0 {
			return s
		}
		s = s[:i] + new + s[i+len(old):]
	}
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
