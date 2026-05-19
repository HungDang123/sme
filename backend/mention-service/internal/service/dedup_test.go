package service

import (
	"testing"
	"time"
)

func TestBuildDedupKeyExternalID(t *testing.T) {
	key := BuildDedupKey("facebook", "post-123", "https://example.com", "content", time.Now())
	if key != "facebook:post-123" {
		t.Fatalf("expected external id key, got %s", key)
	}
}

func TestBuildDedupKeyURL(t *testing.T) {
	key := BuildDedupKey("youtube", "", "https://youtube.com/watch?v=abc", "content", time.Now())
	if len(key) < 20 || key[:10] != "youtube:ur" {
		t.Fatalf("unexpected url key: %s", key)
	}
}

func TestBuildDedupKeyContent(t *testing.T) {
	ts := time.Date(2026, 5, 19, 10, 0, 0, 0, time.UTC)
	key1 := BuildDedupKey("facebook", "", "", "same content", ts)
	key2 := BuildDedupKey("facebook", "", "", "same content", ts)
	if key1 != key2 {
		t.Fatal("content keys should match for same input")
	}
}
