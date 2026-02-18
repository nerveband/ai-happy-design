package imgutil

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ResolveImageData checks if imageData is a file path (file:// URI, absolute path,
// or ~ path) and reads/encodes it as a base64 data URI. If already base64 or a
// data URI, it passes through unchanged.
func ResolveImageData(imageData string) (string, error) {
	if imageData == "" {
		return imageData, nil
	}

	// Already a data URI — pass through
	if strings.HasPrefix(imageData, "data:") {
		return imageData, nil
	}

	// HTTP/HTTPS URL — download and encode
	if strings.HasPrefix(imageData, "https://") || strings.HasPrefix(imageData, "http://") {
		return downloadAndEncode(imageData)
	}

	// Determine if this is a file reference
	var filePath string
	switch {
	case strings.HasPrefix(imageData, "file://"):
		filePath = strings.TrimPrefix(imageData, "file://")
	case strings.HasPrefix(imageData, "/"):
		filePath = imageData
	case strings.HasPrefix(imageData, "~/"):
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot expand ~: %w", err)
		}
		filePath = filepath.Join(home, imageData[2:])
	default:
		// Not a path — assume it's already base64, pass through
		return imageData, nil
	}

	// Expand ~ in file:// paths too (file:///Users/... won't have ~, but file://~/... might)
	if strings.HasPrefix(filePath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot expand ~: %w", err)
		}
		filePath = filepath.Join(home, filePath[2:])
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("imageData file not found: %s", filePath)
		}
		return "", fmt.Errorf("cannot read imageData file: %w", err)
	}

	mime := mimeFromExt(filepath.Ext(filePath))
	if mime == "image/svg+xml" {
		return "", fmt.Errorf("SVG files cannot be used as image fills — Figma's createImage only accepts raster formats (PNG/JPG/WebP/GIF). To import an SVG as vector nodes, use shape.create_from_svg with the SVG content instead")
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mime, encoded), nil
}

// ResolveParamsImageData scans a params map for imageData and resolves file paths.
// Returns the (possibly modified) params and whether any resolution occurred.
func ResolveParamsImageData(params map[string]interface{}) (map[string]interface{}, bool, error) {
	raw, ok := params["imageData"]
	if !ok {
		return params, false, nil
	}
	str, ok := raw.(string)
	if !ok || str == "" {
		return params, false, nil
	}

	resolved, err := ResolveImageData(str)
	if err != nil {
		return params, false, err
	}
	if resolved == str {
		return params, false, nil
	}

	// Shallow copy to avoid mutating original
	out := make(map[string]interface{}, len(params))
	for k, v := range params {
		out[k] = v
	}
	out["imageData"] = resolved
	return out, true, nil
}

// downloadAndEncode fetches an image from a URL and returns it as a base64 data URI.
func downloadAndEncode(url string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image from %s: HTTP %d", url, resp.StatusCode)
	}

	// Limit to 50MB
	data, err := io.ReadAll(io.LimitReader(resp.Body, 50*1024*1024))
	if err != nil {
		return "", fmt.Errorf("failed to read image from %s: %w", url, err)
	}

	// Determine MIME type: prefer Content-Type header, fall back to URL extension
	mime := ""
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		// Take just the media type, strip params like charset
		if idx := strings.Index(ct, ";"); idx > 0 {
			ct = strings.TrimSpace(ct[:idx])
		}
		if strings.HasPrefix(ct, "image/") {
			mime = ct
		}
	}
	if mime == "" {
		// Guess from URL path extension
		urlPath := url
		if idx := strings.Index(urlPath, "?"); idx > 0 {
			urlPath = urlPath[:idx]
		}
		mime = mimeFromExt(filepath.Ext(urlPath))
	}
	if mime == "" || mime == "application/octet-stream" {
		// Detect from magic bytes
		mime = detectMimeFromBytes(data)
	}
	if mime == "image/svg+xml" {
		return "", fmt.Errorf("SVG URLs cannot be used as image fills — Figma's createImage only accepts raster formats (PNG/JPG/WebP/GIF). To import an SVG as vector nodes, use shape.create_from_svg with the SVG content instead")
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mime, encoded), nil
}

func detectMimeFromBytes(data []byte) string {
	if len(data) < 4 {
		return "application/octet-stream"
	}
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "image/png"
	}
	if data[0] == 0xFF && data[1] == 0xD8 {
		return "image/jpeg"
	}
	if string(data[:4]) == "RIFF" && len(data) > 12 && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}
	if string(data[:3]) == "GIF" {
		return "image/gif"
	}
	return "application/octet-stream"
}

func mimeFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
