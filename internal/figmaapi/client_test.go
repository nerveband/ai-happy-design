package figmaapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClientUsesBearerTokenAndReadsOEmbed(t *testing.T) {
	t.Setenv("FIGMA_ACCESS_TOKEN", "env-token")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer env-token" {
			t.Fatalf("Authorization = %q", got)
		}
		if r.URL.Path != "/v1/oembed" || r.URL.Query().Get("url") != "https://www.figma.com/file/abc/Test" {
			t.Fatalf("unexpected request: %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"title": "Test file", "type": "rich"})
	}))
	defer server.Close()

	client := NewClient(Options{BaseURL: server.URL})
	got, err := client.OEmbed(context.Background(), "https://www.figma.com/file/abc/Test")
	if err != nil {
		t.Fatal(err)
	}
	if got["title"] != "Test file" {
		t.Fatalf("title = %v", got["title"])
	}
}

func TestClientReturnsTypedAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"err":"missing scope"}`, http.StatusForbidden)
	}))
	defer server.Close()

	_, err := NewClient(Options{BaseURL: server.URL, Token: "token"}).FileMetadata(context.Background(), "abc", nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.Message != "missing scope" {
		t.Fatalf("api error = %#v", apiErr)
	}
}

func TestFileMetadataAndDevResourcesRequests(t *testing.T) {
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Method+" "+r.URL.RequestURI())
		switch r.Method + " " + r.URL.Path {
		case "GET /v1/files/abc":
			_ = json.NewEncoder(w).Encode(map[string]any{"name": "File"})
		case "GET /v1/files/abc/dev_resources":
			_ = json.NewEncoder(w).Encode(map[string]any{"dev_resources": []any{map[string]any{"id": "dr1"}}})
		case "POST /v1/files/abc/dev_resources":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "dr2"})
		case "PUT /v1/files/abc/dev_resources/dr2":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "dr2", "name": "Updated"})
		case "DELETE /v1/files/abc/dev_resources/dr2":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(Options{BaseURL: server.URL, Token: "token"})
	if _, err := client.FileMetadata(context.Background(), "abc", map[string]string{"ids": "1:2"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ListDevResources(context.Background(), "abc"); err != nil {
		t.Fatal(err)
	}
	if _, err := client.CreateDevResource(context.Background(), "abc", map[string]any{"name": "Spec"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.UpdateDevResource(context.Background(), "abc", "dr2", map[string]any{"name": "Updated"}); err != nil {
		t.Fatal(err)
	}
	if err := client.DeleteDevResource(context.Background(), "abc", "dr2"); err != nil {
		t.Fatal(err)
	}
	wantLast := "DELETE /v1/files/abc/dev_resources/dr2"
	if seen[len(seen)-1] != wantLast {
		t.Fatalf("last request = %q, want %q", seen[len(seen)-1], wantLast)
	}
}

func TestWebhooksV2SupportsDevModeStatusUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/webhooks" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["event_type"] != EventDevModeStatusUpdate {
			t.Fatalf("event_type = %v", body["event_type"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "wh1"})
	}))
	defer server.Close()

	got, err := NewClient(Options{BaseURL: server.URL, Token: "token"}).CreateWebhook(context.Background(), map[string]any{
		"event_type": EventDevModeStatusUpdate,
		"context":    "file",
		"context_id": "abc",
		"endpoint":   "https://example.com/hook",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got["id"] != "wh1" {
		t.Fatalf("id = %v", got["id"])
	}
}

func TestMissingTokenIsReportedBeforeHTTPRequest(t *testing.T) {
	_ = os.Unsetenv("FIGMA_ACCESS_TOKEN")
	_ = os.Unsetenv("FIGMA_TOKEN")

	_, err := NewClient(Options{BaseURL: "http://127.0.0.1:1"}).OEmbed(context.Background(), "https://figma.com/file/abc")
	if err != ErrMissingToken {
		t.Fatalf("err = %v, want ErrMissingToken", err)
	}
}
