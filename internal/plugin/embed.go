package plugin

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed files/manifest.json files/dist/code.js files/dist/ui.html
var pluginFiles embed.FS

// DefaultPluginDir returns ~/.ahd-figma/plugin, falling back to the legacy
// ~/.ai-happy-design/plugin directory when present.
func DefaultPluginDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	newDir := filepath.Join(home, ".ahd-figma", "plugin")
	legacyDir := filepath.Join(home, ".ai-happy-design", "plugin")
	if _, err := os.Stat(newDir); err == nil {
		return newDir, nil
	}
	if _, err := os.Stat(legacyDir); err == nil {
		return legacyDir, nil
	}
	return newDir, nil
}

// NeedsUpdate returns true if the plugin at dir is missing or has a different version.
func NeedsUpdate(dir, version string) bool {
	data, err := os.ReadFile(filepath.Join(dir, ".version"))
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(data)) != version
}

// ExtractTo writes the embedded plugin files to dir and stamps .version.
func ExtractTo(dir, version string) error {
	err := fs.WalkDir(pluginFiles, "files", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel("files", path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		target := filepath.Join(dir, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := pluginFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded %s: %w", path, err)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		return os.WriteFile(target, data, 0644)
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, ".version"), []byte(version), 0644)
}
