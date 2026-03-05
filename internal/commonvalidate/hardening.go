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
