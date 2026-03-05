package host

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	defaultAppName  = "Adobe Illustrator"
	defaultBundleID = "com.adobe.Illustrator"
	defaultTimeout  = 30 * time.Second
)

var lookPathFunc = exec.LookPath
var runCommand = runCommandImpl

// Status describes the local Illustrator host state.
type Status struct {
	Platform            string `json:"platform"`
	Supported           bool   `json:"supported"`
	IllustratorRunning  bool   `json:"illustratorRunning"`
	IllustratorAppFound bool   `json:"illustratorAppFound"`
	OSAScriptAvailable  bool   `json:"osascriptAvailable"`
	AppPath             string `json:"appPath,omitempty"`
	Version             string `json:"version,omitempty"`
}

// Adapter drives Illustrator through macOS host tools.
type Adapter struct {
	AppName  string
	BundleID string
}

// NewAdapter constructs a host adapter.
func NewAdapter() *Adapter {
	return &Adapter{
		AppName:  defaultAppName,
		BundleID: defaultBundleID,
	}
}

// Status inspects whether the required macOS tools and Illustrator app are available.
func (a *Adapter) Status() Status {
	status := Status{
		Platform:           runtime.GOOS,
		Supported:          runtime.GOOS == "darwin",
		OSAScriptAvailable: hasCommand("osascript"),
	}
	if !status.Supported || !status.OSAScriptAvailable {
		return status
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	status.AppPath = a.resolveAppPath(ctx)
	status.IllustratorAppFound = status.AppPath != ""
	if status.IllustratorAppFound {
		status.Version = readBundleVersion(ctx, status.AppPath)
	}

	out, err := runCommand(ctx, "osascript", "-e", fmt.Sprintf(`tell application "System Events" to (name of processes) contains "%s"`, a.AppName))
	if err == nil {
		status.IllustratorRunning = strings.Contains(strings.ToLower(strings.TrimSpace(out)), "true")
	}
	return status
}

func (a *Adapter) resolveAppPath(ctx context.Context) string {
	if ctx == nil {
		ctx = context.Background()
	}
	out, err := runCommand(ctx, "osascript", "-e", fmt.Sprintf(`POSIX path of (path to application id "%s")`, a.BundleID))
	if err != nil {
		return ""
	}
	return normalizeAppPath(out)
}

func normalizeAppPath(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	return strings.TrimSuffix(trimmed, "/")
}

// Open launches Illustrator using macOS open.
func (a *Adapter) Open() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("ILLUSTRATOR_NOT_RUNNING: Illustrator host is only supported on macOS")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := runCommand(ctx, "open", "-a", a.AppName)
	if err != nil {
		return fmt.Errorf("HOST_EXEC_ERROR: failed to open Illustrator: %w", err)
	}
	return nil
}

// Quit exits Illustrator gracefully.
func (a *Adapter) Quit() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("ILLUSTRATOR_NOT_RUNNING: Illustrator host is only supported on macOS")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := runCommand(ctx, "osascript", "-e", fmt.Sprintf(`tell application "%s" to quit`, a.AppName))
	if err != nil {
		return fmt.Errorf("HOST_EXEC_ERROR: failed to quit Illustrator: %w", err)
	}
	return nil
}

// ExecuteJavaScript runs a JSX snippet through Illustrator's do javascript AppleScript surface.
func (a *Adapter) ExecuteJavaScript(script string, timeout time.Duration) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("ILLUSTRATOR_NOT_RUNNING: Illustrator host is only supported on macOS")
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	tmpDir, err := os.MkdirTemp("", "ahd-illustrator-*.jsx")
	if err != nil {
		return "", fmt.Errorf("HOST_EXEC_ERROR: create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := filepath.Join(tmpDir, "command.jsx")
	if err := os.WriteFile(scriptPath, []byte(script), 0o644); err != nil {
		return "", fmt.Errorf("HOST_EXEC_ERROR: write temp jsx: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := runCommand(ctx, "osascript", "-e", fmt.Sprintf(`tell application "%s" to do javascript POSIX file "%s"`, a.AppName, scriptPath))
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("HOST_EXEC_ERROR: Illustrator script timed out after %s", timeout)
		}
		return "", fmt.Errorf("HOST_EXEC_ERROR: %w", err)
	}
	return strings.TrimSpace(out), nil
}

func runCommandImpl(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func hasCommand(name string) bool {
	_, err := lookPathFunc(name)
	return err == nil
}

func readBundleVersion(ctx context.Context, appPath string) string {
	if ctx == nil || strings.TrimSpace(appPath) == "" {
		return ""
	}
	infoPlist := strings.TrimSuffix(appPath, "/") + "/Contents/Info.plist"
	out, err := runCommand(ctx, "/usr/libexec/PlistBuddy", "-c", "Print :CFBundleShortVersionString", infoPlist)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}
