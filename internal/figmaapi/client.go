package figmaapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultBaseURL           = "https://api.figma.com"
	EventDevModeStatusUpdate = "DEV_MODE_STATUS_UPDATE"
)

var ErrMissingToken = errors.New("missing Figma token: set FIGMA_ACCESS_TOKEN or FIGMA_TOKEN, or pass token")

type Options struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("figma API error %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("figma API error %d", e.StatusCode)
}

func NewClient(opts Options) *Client {
	baseURL := strings.TrimRight(opts.BaseURL, "/")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		baseURL:    baseURL,
		token:      strings.TrimSpace(opts.Token),
		httpClient: httpClient,
	}
}

func ResolveToken(explicit string) string {
	if strings.TrimSpace(explicit) != "" {
		return strings.TrimSpace(explicit)
	}
	if token := strings.TrimSpace(os.Getenv("FIGMA_ACCESS_TOKEN")); token != "" {
		return token
	}
	return strings.TrimSpace(os.Getenv("FIGMA_TOKEN"))
}

func (c *Client) OEmbed(ctx context.Context, figmaURL string) (map[string]any, error) {
	q := url.Values{"url": []string{figmaURL}}
	return c.doJSON(ctx, http.MethodGet, "/v1/oembed?"+q.Encode(), nil)
}

func (c *Client) FileMetadata(ctx context.Context, fileKey string, query map[string]string) (map[string]any, error) {
	path := "/v1/files/" + url.PathEscape(fileKey)
	if len(query) > 0 {
		q := url.Values{}
		for k, v := range query {
			if v != "" {
				q.Set(k, v)
			}
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	return c.doJSON(ctx, http.MethodGet, path, nil)
}

func (c *Client) ListDevResources(ctx context.Context, fileKey string) (map[string]any, error) {
	return c.doJSON(ctx, http.MethodGet, "/v1/files/"+url.PathEscape(fileKey)+"/dev_resources", nil)
}

func (c *Client) CreateDevResource(ctx context.Context, fileKey string, payload map[string]any) (map[string]any, error) {
	return c.doJSON(ctx, http.MethodPost, "/v1/files/"+url.PathEscape(fileKey)+"/dev_resources", payload)
}

func (c *Client) UpdateDevResource(ctx context.Context, fileKey, resourceID string, payload map[string]any) (map[string]any, error) {
	return c.doJSON(ctx, http.MethodPut, "/v1/files/"+url.PathEscape(fileKey)+"/dev_resources/"+url.PathEscape(resourceID), payload)
}

func (c *Client) DeleteDevResource(ctx context.Context, fileKey, resourceID string) error {
	_, err := c.doJSON(ctx, http.MethodDelete, "/v1/files/"+url.PathEscape(fileKey)+"/dev_resources/"+url.PathEscape(resourceID), nil)
	return err
}

func (c *Client) ListWebhooks(ctx context.Context) (map[string]any, error) {
	return c.doJSON(ctx, http.MethodGet, "/v2/webhooks", nil)
}

func (c *Client) CreateWebhook(ctx context.Context, payload map[string]any) (map[string]any, error) {
	return c.doJSON(ctx, http.MethodPost, "/v2/webhooks", payload)
}

func (c *Client) GetWebhook(ctx context.Context, webhookID string) (map[string]any, error) {
	return c.doJSON(ctx, http.MethodGet, "/v2/webhooks/"+url.PathEscape(webhookID), nil)
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, payload map[string]any) (map[string]any, error) {
	return c.doJSON(ctx, http.MethodPut, "/v2/webhooks/"+url.PathEscape(webhookID), payload)
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	_, err := c.doJSON(ctx, http.MethodDelete, "/v2/webhooks/"+url.PathEscape(webhookID), nil)
	return err
}

func (c *Client) doJSON(ctx context.Context, method, path string, payload map[string]any) (map[string]any, error) {
	token := ResolveToken(c.token)
	if token == "" {
		return nil, ErrMissingToken
	}
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apiError(resp.StatusCode, raw)
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]any{"ok": true}, nil
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func apiError(status int, raw []byte) *APIError {
	body := strings.TrimSpace(string(raw))
	message := body
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err == nil {
		for _, key := range []string{"err", "message", "error"} {
			if value, ok := parsed[key].(string); ok {
				message = value
				break
			}
		}
	}
	return &APIError{StatusCode: status, Message: message, Body: body}
}
