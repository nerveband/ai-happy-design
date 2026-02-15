// Package imgutil provides opt-in image compression for CLI image operations.
// Uses ImageMagick (convert/magick) when available for cross-platform support.
package imgutil

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
)

// Options controls image compression behavior.
type Options struct {
	Quality    int    // JPEG/WebP quality 1-100 (default 80)
	MaxWidth   int    // Max width in pixels (0 = no limit)
	MaxHeight  int    // Max height in pixels (0 = no limit)
	FormatHint string // Output format hint: "jpeg", "png", "webp" (default: auto-detect from input)
}

// DefaultOptions returns sensible compression defaults.
func DefaultOptions() Options {
	return Options{
		Quality:  80,
		MaxWidth: 2048,
	}
}

// HasImageMagick checks if ImageMagick is available on the system.
func HasImageMagick() bool {
	// Try "magick" first (ImageMagick 7+), then "convert" (ImageMagick 6)
	if _, err := exec.LookPath("magick"); err == nil {
		return true
	}
	if _, err := exec.LookPath("convert"); err == nil {
		return true
	}
	return false
}

// imageMagickCmd returns the correct ImageMagick command name.
func imageMagickCmd() string {
	if _, err := exec.LookPath("magick"); err == nil {
		return "magick"
	}
	return "convert"
}

// CompressBase64 takes a base64-encoded image (with or without data URL prefix),
// compresses it using ImageMagick, and returns the compressed base64 string.
// Returns the original string unchanged if ImageMagick is not available.
func CompressBase64(imageData string, opts Options) (string, error) {
	if !HasImageMagick() {
		return imageData, nil
	}

	// Strip data URL prefix if present
	rawBase64 := imageData
	prefix := ""
	if idx := strings.Index(imageData, ","); idx >= 0 && strings.HasPrefix(imageData, "data:") {
		prefix = imageData[:idx+1]
		rawBase64 = imageData[idx+1:]
	}

	// Decode base64 to bytes
	decoded, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		// Try RawStdEncoding (no padding)
		decoded, err = base64.RawStdEncoding.DecodeString(rawBase64)
		if err != nil {
			return imageData, fmt.Errorf("failed to decode base64: %w", err)
		}
	}

	originalSize := len(decoded)

	// Build ImageMagick command
	cmdName := imageMagickCmd()
	args := []string{"-"} // Read from stdin

	// Resize if max dimensions set
	if opts.MaxWidth > 0 || opts.MaxHeight > 0 {
		w := opts.MaxWidth
		h := opts.MaxHeight
		if w == 0 {
			w = 99999
		}
		if h == 0 {
			h = 99999
		}
		args = append(args, "-resize", fmt.Sprintf("%dx%d>", w, h))
	}

	// Strip metadata to save bytes
	args = append(args, "-strip")

	// Set quality
	quality := opts.Quality
	if quality <= 0 {
		quality = 80
	}
	args = append(args, "-quality", fmt.Sprintf("%d", quality))

	// Determine output format
	outFormat := opts.FormatHint
	if outFormat == "" {
		// Detect from input bytes (PNG magic: 89 50 4E 47)
		if len(decoded) > 4 && decoded[0] == 0x89 && decoded[1] == 0x50 {
			outFormat = "png"
		} else {
			outFormat = "jpeg"
		}
	}
	args = append(args, fmt.Sprintf("%s:-", outFormat)) // Write to stdout in format

	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = bytes.NewReader(decoded)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return imageData, fmt.Errorf("ImageMagick compression failed: %s (%w)", stderr.String(), err)
	}

	compressed := stdout.Bytes()

	// Only use compressed version if it's actually smaller
	if len(compressed) >= originalSize {
		return imageData, nil
	}

	// Re-encode to base64
	result := base64.StdEncoding.EncodeToString(compressed)

	// Restore data URL prefix if it was present, updating MIME type
	if prefix != "" {
		mime := "image/jpeg"
		switch outFormat {
		case "png":
			mime = "image/png"
		case "webp":
			mime = "image/webp"
		}
		result = fmt.Sprintf("data:%s;base64,%s", mime, result)
	}

	return result, nil
}

// CompressResult holds info about a compression operation.
type CompressResult struct {
	OriginalSize   int     `json:"originalSize"`
	CompressedSize int     `json:"compressedSize"`
	Ratio          float64 `json:"ratio"`
	Skipped        bool    `json:"skipped"`
}

// CompressBase64WithInfo is like CompressBase64 but also returns size info.
func CompressBase64WithInfo(imageData string, opts Options) (string, CompressResult, error) {
	originalB64Len := len(imageData)

	result, err := CompressBase64(imageData, opts)
	if err != nil {
		return imageData, CompressResult{Skipped: true}, err
	}

	compressedB64Len := len(result)
	ratio := 1.0
	if originalB64Len > 0 {
		ratio = float64(compressedB64Len) / float64(originalB64Len)
	}

	return result, CompressResult{
		OriginalSize:   originalB64Len,
		CompressedSize: compressedB64Len,
		Ratio:          ratio,
		Skipped:        compressedB64Len >= originalB64Len,
	}, nil
}

// CompressParamsImageData scans a params map for imageData keys and compresses them.
// Returns the modified params map and whether any compression was applied.
func CompressParamsImageData(params map[string]interface{}, opts Options) (map[string]interface{}, bool, error) {
	imageData, ok := params["imageData"]
	if !ok {
		return params, false, nil
	}

	str, ok := imageData.(string)
	if !ok || str == "" {
		return params, false, nil
	}

	compressed, err := CompressBase64(str, opts)
	if err != nil {
		return params, false, err
	}

	if compressed == str {
		return params, false, nil
	}

	// Make a shallow copy to avoid mutating the original
	out := make(map[string]interface{}, len(params))
	for k, v := range params {
		out[k] = v
	}
	out["imageData"] = compressed
	return out, true, nil
}
