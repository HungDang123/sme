package telegram

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"sme-social-listening/alert-service/internal/domain"
)

type PendingItem struct {
	Alert   domain.CreateAlertRequest
	BrandID string
}

type Batcher struct {
	mu          sync.Mutex
	pending     map[string][]PendingItem
	telegram    *Client
	chatLookup  func(brandID string) string
	flushEvery  time.Duration
	maxBatch    int
	immediateNeg bool
}

func NewBatcher(
	telegram *Client,
	chatLookup func(brandID string) string,
	flushEvery time.Duration,
	maxBatch int,
) *Batcher {
	b := &Batcher{
		pending:    map[string][]PendingItem{},
		telegram:   telegram,
		chatLookup: chatLookup,
		flushEvery: flushEvery,
		maxBatch:   maxBatch,
	}
	go b.loop()
	return b
}

func (b *Batcher) Add(req domain.CreateAlertRequest) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if req.AlertType == "negative" {
		chatID := b.chatLookup(req.BrandID)
		_ = b.telegram.SendMessage(chatID, formatSingle(req, true))
		return
	}

	b.pending[req.BrandID] = append(b.pending[req.BrandID], PendingItem{Alert: req, BrandID: req.BrandID})
	if len(b.pending[req.BrandID]) >= b.maxBatch {
		b.flushBrand(req.BrandID)
	}
}

func (b *Batcher) loop() {
	ticker := time.NewTicker(b.flushEvery)
	defer ticker.Stop()
	for range ticker.C {
		b.FlushAll()
	}
}

func (b *Batcher) FlushAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for brandID := range b.pending {
		b.flushBrand(brandID)
	}
}

func (b *Batcher) flushBrand(brandID string) {
	items := b.pending[brandID]
	if len(items) == 0 {
		return
	}
	delete(b.pending, brandID)

	chatID := b.chatLookup(brandID)
	if chatID == "" {
		return
	}

	neg := 0
	for _, item := range items {
		if item.Alert.AlertType == "negative" {
			neg++
		}
	}

	text := fmt.Sprintf("<b>Social Listening Alert</b>\nTong %d mention moi\nTieu cuc: %d\n\n", len(items), neg)
	limit := 3
	if len(items) < limit {
		limit = len(items)
	}
	for i := 0; i < limit; i++ {
		text += formatSingle(items[i].Alert, false) + "\n\n"
	}
	_ = b.telegram.SendMessage(chatID, text)
}

func formatSingle(req domain.CreateAlertRequest, urgent bool) string {
	prefix := ""
	if urgent {
		prefix = "🚨 "
	}
	snippet := req.Content
	if len([]rune(snippet)) > 120 {
		runes := []rune(snippet)
		snippet = string(runes[:120]) + "..."
	}
	link := req.URL
	if link == "" {
		link = "(khong co link)"
	}
	return fmt.Sprintf("%s<b>%s</b> | %s | %s\nSentiment: %s\n%s\n%s",
		prefix, EscapeHTML(req.Keyword), EscapeHTML(req.Source), EscapeHTML(req.AlertType),
		EscapeHTML(req.Sentiment), EscapeHTML(snippet), link,
	)
}

func truncateLines(items []PendingItem) string {
	var b strings.Builder
	for _, item := range items {
		b.WriteString("- ")
		b.WriteString(item.Alert.Keyword)
		b.WriteString("\n")
	}
	return b.String()
}
