package imgutil

import (
	"encoding/base64"
	"testing"
)

func TestHasImageMagick(t *testing.T) {
	// Just verify it doesn't panic — result depends on system
	_ = HasImageMagick()
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.Quality != 80 {
		t.Errorf("expected quality 80, got %d", opts.Quality)
	}
	if opts.MaxWidth != 2048 {
		t.Errorf("expected maxWidth 2048, got %d", opts.MaxWidth)
	}
}

func TestCompressBase64_NoImageMagick(t *testing.T) {
	// If ImageMagick isn't installed, should return original unchanged
	if HasImageMagick() {
		t.Skip("ImageMagick is installed, skipping no-magick test")
	}

	input := "iVBORw0KGgoAAAANSUhEUg=="
	result, err := CompressBase64(input, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != input {
		t.Fatal("expected unchanged output when ImageMagick not available")
	}
}

func TestCompressBase64_InvalidBase64(t *testing.T) {
	if !HasImageMagick() {
		t.Skip("ImageMagick not installed")
	}

	// Completely invalid base64 — shouldn't crash, should return original
	_, err := CompressBase64("not-valid-base64!!!", DefaultOptions())
	if err == nil {
		// Either returns error or returns original
		t.Log("returned nil error (acceptable)")
	}
}

func TestCompressBase64_DataURLPrefix(t *testing.T) {
	if !HasImageMagick() {
		t.Skip("ImageMagick not installed")
	}

	// Create a tiny valid PNG (1x1 red pixel)
	pngBytes := createTinyPNG()
	b64 := base64.StdEncoding.EncodeToString(pngBytes)
	input := "data:image/png;base64," + b64

	result, err := CompressBase64(input, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still have a data URL prefix (or be unchanged if not smaller)
	if result != input {
		if len(result) > 5 && result[:5] != "data:" {
			t.Error("expected data URL prefix to be preserved")
		}
	}
}

func TestCompressParamsImageData_NoImageData(t *testing.T) {
	params := map[string]interface{}{
		"x":    100,
		"name": "test",
	}
	out, compressed, err := CompressParamsImageData(params, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if compressed {
		t.Error("expected no compression for params without imageData")
	}
	if out["x"] != 100 || out["name"] != "test" {
		t.Error("expected unchanged params")
	}
}

func TestCompressParamsImageData_EmptyImageData(t *testing.T) {
	params := map[string]interface{}{
		"imageData": "",
	}
	_, compressed, err := CompressParamsImageData(params, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if compressed {
		t.Error("expected no compression for empty imageData")
	}
}

func TestCompressBase64WithInfo(t *testing.T) {
	if !HasImageMagick() {
		t.Skip("ImageMagick not installed")
	}

	pngBytes := createTinyPNG()
	b64 := base64.StdEncoding.EncodeToString(pngBytes)

	_, info, err := CompressBase64WithInfo(b64, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.OriginalSize <= 0 {
		t.Error("expected positive original size")
	}
}

// createTinyPNG returns bytes for a minimal valid 1x1 red pixel PNG (69 bytes).
func createTinyPNG() []byte {
	return []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
		0x00, 0x03, 0x01, 0x01, 0x00, 0xc9, 0xfe, 0x92, 0xef, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
		0x44, 0xae, 0x42, 0x60, 0x82,
	}
}
