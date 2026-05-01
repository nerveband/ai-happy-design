package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/creativeprojects/go-selfupdate"
	"github.com/spf13/cobra"
)

const repoOwner = "nerveband"
const repoName = "ai-happy-design"

type updateCheckCache struct {
	LastCheck      time.Time `json:"last_check"`
	LatestVersion  string    `json:"latest_version"`
	UpdateRequired bool      `json:"update_required"`
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".ai-happy-design")
}

func checkForUpdates() (bool, string, error) {
	cached, err := loadUpdateCache()
	if err == nil && time.Since(cached.LastCheck) < 24*time.Hour {
		return cached.UpdateRequired, cached.LatestVersion, nil
	}

	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return false, "", err
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source:    source,
		Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"},
	})
	if err != nil {
		return false, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	latest, found, err := updater.DetectLatest(ctx, selfupdate.NewRepositorySlug(repoOwner, repoName))
	if err != nil || !found {
		return false, "", err
	}

	hasUpdate := latest.GreaterThan(version)
	latestVer := latest.Version()

	saveUpdateCache(updateCheckCache{
		LastCheck:      time.Now(),
		LatestVersion:  latestVer,
		UpdateRequired: hasUpdate,
	})

	return hasUpdate, latestVer, nil
}

func notifyUpdateAvailable() {
	if os.Getenv("AHD_NO_UPDATE_CHECK") == "1" || os.Getenv("CI") == "true" {
		return
	}
	hasUpdate, latestVersion, err := checkForUpdates()
	if err != nil || !hasUpdate {
		return
	}
	fmt.Fprintf(os.Stderr, "\nNew version available: %s (current: %s)\n", latestVersion, version)
	fmt.Fprintf(os.Stderr, "Run '%s upgrade' to update\n\n", filepath.Base(os.Args[0]))
}

func loadUpdateCache() (*updateCheckCache, error) {
	dir := configDir()
	if dir == "" {
		return nil, fmt.Errorf("no home directory")
	}
	data, err := os.ReadFile(filepath.Join(dir, "update_cache.json"))
	if err != nil {
		return nil, err
	}
	var cache updateCheckCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

func saveUpdateCache(cache updateCheckCache) error {
	dir := configDir()
	if dir == "" {
		return fmt.Errorf("no home directory")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "update_cache.json"), data, 0644)
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version",
	Long:  "Check for and install the latest version from GitHub releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runUpgrade()
	},
}

func runUpgrade() error {
	fmt.Printf("Current version: %s\n", version)
	fmt.Printf("Checking for updates...\n")

	source, err := selfupdate.NewGitHubSource(selfupdate.GitHubConfig{})
	if err != nil {
		return fmt.Errorf("failed to create update source: %w", err)
	}

	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Source:    source,
		Validator: &selfupdate.ChecksumValidator{UniqueFilename: "checksums.txt"},
	})
	if err != nil {
		return fmt.Errorf("failed to create updater: %w", err)
	}

	latest, found, err := updater.DetectLatest(context.Background(), selfupdate.NewRepositorySlug(repoOwner, repoName))
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !found {
		fmt.Println("No releases found")
		return nil
	}

	if latest.LessOrEqual(version) {
		fmt.Printf("Already up to date (latest: %s)\n", latest.Version())
		return nil
	}

	fmt.Printf("New version available: %s\n", latest.Version())
	fmt.Printf("Downloading for %s/%s...\n", runtime.GOOS, runtime.GOARCH)

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	if err := updater.UpdateTo(context.Background(), latest, exe); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	// macOS requires ad-hoc signing or the binary gets SIGKILL'd
	if runtime.GOOS == "darwin" {
		fmt.Println("Signing binary for macOS...")
		if signErr := signBinary(exe); signErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: ad-hoc signing failed: %v\n", signErr)
			fmt.Fprintf(os.Stderr, "Run manually: codesign -s - -f %s\n", exe)
		}
	}

	fmt.Printf("Successfully upgraded to %s\n", latest.Version())
	return nil
}

func signBinary(path string) error {
	cmd := exec.Command("codesign", "-s", "-", "-f", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
