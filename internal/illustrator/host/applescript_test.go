package host

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestStatusResolvesBundlePathAndVersion(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("host status bundle resolution is only used on macOS")
	}

	oldLookPath := lookPathFunc
	lookPathFunc = func(name string) (string, error) {
		if name == "osascript" {
			return "/usr/bin/osascript", nil
		}
		return "", exec.ErrNotFound
	}
	t.Cleanup(func() {
		lookPathFunc = oldLookPath
	})

	oldRunCommand := runCommand
	runCommand = func(ctx context.Context, name string, args ...string) (string, error) {
		switch {
		case name == "osascript" && containsArg(args, `path to application id "com.adobe.Illustrator"`):
			return "/Applications/Adobe Illustrator 2026/Adobe Illustrator.app/\n", nil
		case name == "osascript" && containsArg(args, `tell application "System Events"`):
			return "true\n", nil
		case name == "/usr/libexec/PlistBuddy":
			return "30.2.1\n", nil
		default:
			return "", errors.New("unexpected command")
		}
	}
	t.Cleanup(func() {
		runCommand = oldRunCommand
	})

	status := NewAdapter().Status()
	if !status.IllustratorAppFound {
		t.Fatal("expected bundle-id app lookup to mark Illustrator as installed")
	}
	if !status.IllustratorRunning {
		t.Fatal("expected mocked running status to be true")
	}
	if got, want := status.AppPath, "/Applications/Adobe Illustrator 2026/Adobe Illustrator.app"; got != want {
		t.Fatalf("unexpected app path %q, want %q", got, want)
	}
	if got, want := status.Version, "30.2.1"; got != want {
		t.Fatalf("unexpected version %q, want %q", got, want)
	}
}

func TestNormalizeAppPath(t *testing.T) {
	t.Parallel()

	if got, want := normalizeAppPath(" /Applications/Adobe Illustrator 2026/Adobe Illustrator.app/\n"), "/Applications/Adobe Illustrator 2026/Adobe Illustrator.app"; got != want {
		t.Fatalf("unexpected normalized path %q, want %q", got, want)
	}
	if got := normalizeAppPath(" \n "); got != "" {
		t.Fatalf("expected empty path for blank input, got %q", got)
	}
}

func containsArg(args []string, needle string) bool {
	for _, arg := range args {
		if strings.Contains(arg, needle) {
			return true
		}
	}
	return false
}
