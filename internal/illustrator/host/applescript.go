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
	defaultAppName = "Adobe Illustrator"
	defaultTimeout = 30 * time.Second
)

// Status describes the local Illustrator host state.
type Status struct {
	Platform            string `json:"platform"`
	Supported           bool   `json:"supported"`
	IllustratorRunning  bool   `json:"illustratorRunning"`
	IllustratorAppFound bool   `json:"illustratorAppFound"`
	OSAScriptAvailable  bool   `json:"osascriptAvailable"`
}

// Adapter drives Illustrator through macOS host tools.
type Adapter struct {
	AppName string
}

// NewAdapter constructs a host adapter.
func NewAdapter() *Adapter {
	return &Adapter{AppName: defaultAppName}
}

// Status inspects whether the required macOS tools and Illustrator app are available.
func (a *Adapter) Status() Status {
	status := Status{
		Platform:            runtime.GOOS,
		Supported:           runtime.GOOS == "darwin",
		IllustratorAppFound: appExists(a.AppName),
		OSAScriptAvailable:  hasCommand("osascript"),
	}
	if !status.Supported || !status.OSAScriptAvailable {
		return status
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := runCommand(ctx, "osascript", "-e", fmt.Sprintf(`tell application "System Events" to (name of processes) contains "%s"`, a.AppName))
	if err == nil {
		status.IllustratorRunning = strings.Contains(strings.ToLower(strings.TrimSpace(out)), "true")
	}
	return status
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
			return "", fmt.Errorf("PLUGIN_TIMEOUT: Illustrator script timed out after %s", timeout)
		}
		return "", fmt.Errorf("HOST_EXEC_ERROR: %w", err)
	}
	return strings.TrimSpace(out), nil
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func appExists(appName string) bool {
	_, err := os.Stat(filepath.Join("/Applications", appName+".app"))
	return err == nil
}
