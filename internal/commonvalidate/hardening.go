package commonvalidate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SafePath resolves a user path under cwd and rejects traversal tricks.
func SafePath(cwd, raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	if strings.Contains(trimmed, "\x00") {
		return "", fmt.Errorf("path contains NUL byte")
	}
	lower := strings.ToLower(trimmed)
	if strings.Contains(lower, "://") || strings.HasPrefix(lower, "file:") {
		return "", fmt.Errorf("path must be a local filesystem path")
	}
	if strings.ContainsAny(trimmed, "?#") {
		return "", fmt.Errorf("path must not contain query or fragment syntax")
	}
	if strings.Contains(trimmed, `\`) {
		return "", fmt.Errorf("path must use forward slashes")
	}
	if decoded, changed := normalizePotentialEncodedPath(trimmed); changed {
		decodedLower := strings.ToLower(decoded)
		if strings.Contains(decodedLower, ".."+string(filepath.Separator)) || strings.HasPrefix(decodedLower, "..") || filepath.IsAbs(decoded) {
			return "", fmt.Errorf("path contains encoded traversal or absolute-path syntax")
		}
	}
	for _, snippet := range []string{"%2e", "%2f", "%5c", "%00"} {
		if strings.Contains(lower, snippet) {
			return "", fmt.Errorf("path contains encoded traversal syntax")
		}
	}
	base := cwd
	if base == "" {
		base = "."
	}
	clean := filepath.Clean(trimmed)
	if filepath.IsAbs(clean) {
		return "", fmt.Errorf("absolute paths require an explicit override")
	}
	resolved := filepath.Join(base, clean)
	rel, err := filepath.Rel(base, resolved)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes the working directory")
	}
	return resolved, nil
}
