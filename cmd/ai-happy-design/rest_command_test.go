package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunLocalCommandFigmaOEmbed(t *testing.T) {
	t.Setenv("FIGMA_ACCESS_TOKEN", "token")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/oembed" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"title": "Embed"})
	}))
	defer server.Close()
	t.Setenv("AHD_FIGMA_API_BASE_URL", server.URL)

	got, err := runLocalCommand("figma.oembed", map[string]any{"url": "https://figma.com/file/abc"})
	if err != nil {
		t.Fatal(err)
	}
	if got.(map[string]any)["title"] != "Embed" {
		t.Fatalf("result = %#v", got)
	}
}

func TestRunLocalCommandWebhookCreateDevModeEvent(t *testing.T) {
	t.Setenv("FIGMA_ACCESS_TOKEN", "token")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v2/webhooks" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["event_type"] != "DEV_MODE_STATUS_UPDATE" {
			t.Fatalf("event_type = %v", body["event_type"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "wh1"})
	}))
	defer server.Close()
	t.Setenv("AHD_FIGMA_API_BASE_URL", server.URL)

	got, err := runLocalCommand("figma.webhook_create", map[string]any{
		"eventType": "DEV_MODE_STATUS_UPDATE",
		"context":   "file",
		"contextId": "abc",
		"endpoint":  "https://example.com/hook",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.(map[string]any)["id"] != "wh1" {
		t.Fatalf("result = %#v", got)
	}
}
