package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreprocessBatchImageData_DedupesAndResolves(t *testing.T) {
	// 1x1 PNG
	raw := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO3Z3YoAAAAASUVORK5CYII="
	pngBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		t.Fatalf("decode png: %v", err)
	}

	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "tiny.png")
	if err := os.WriteFile(imgPath, pngBytes, 0644); err != nil {
		t.Fatalf("write temp image: %v", err)
	}

	ops := []batchOperation{
		{
			Name:    "img1",
			Command: "paint.set_image",
			Params: map[string]interface{}{
				"nodeId":    "1:1",
				"imageData": imgPath,
			},
		},
		{
			Name:    "img2",
			Command: "paint.set_image",
			Params: map[string]interface{}{
				"nodeId":    "1:2",
				"imageData": imgPath,
			},
		},
		{
			Name:    "txt",
			Command: "text.create",
			Params: map[string]interface{}{
				"content": "hello",
			},
		},
	}

	prepped, summary := preprocessBatchImageData(ops, false, "", false)
	if summary.Candidates != 2 {
		t.Fatalf("candidates = %d, want 2", summary.Candidates)
	}
	if summary.Unique != 1 {
		t.Fatalf("unique = %d, want 1", summary.Unique)
	}
	if summary.CacheHits != 1 {
		t.Fatalf("cacheHits = %d, want 1", summary.CacheHits)
	}
	if summary.Failed != 0 {
		t.Fatalf("failed = %d, want 0", summary.Failed)
	}
	if summary.Resolved != 1 {
		t.Fatalf("resolved = %d, want 1", summary.Resolved)
	}

	p1, _ := prepped[0].Params["imageData"].(string)
	p2, _ := prepped[1].Params["imageData"].(string)
	if !strings.HasPrefix(p1, "data:image/png;base64,") {
		t.Fatalf("prepped imageData[0] not data URI: %q", p1)
	}
	if p1 != p2 {
		t.Fatalf("expected deduped payload to match, got different imageData values")
	}
}

func TestPreprocessBatchImageData_SkipsInterpolationTokens(t *testing.T) {
	ops := []batchOperation{
		{
			Name:    "img_interp",
			Command: "paint.set_image",
			Params: map[string]interface{}{
				"nodeId":    "1:1",
				"imageData": "${{steps.fetch.result.imageData}}",
			},
		},
	}

	prepped, summary := preprocessBatchImageData(ops, false, "", false)
	if summary.Candidates != 0 || summary.Unique != 0 {
		t.Fatalf("expected interpolation imageData to be skipped, got candidates=%d unique=%d", summary.Candidates, summary.Unique)
	}
	if prepped[0].Params["imageData"] != "${{steps.fetch.result.imageData}}" {
		t.Fatalf("expected interpolation imageData unchanged")
	}
}
