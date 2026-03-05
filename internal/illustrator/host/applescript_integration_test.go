//go:build darwin && integration

package host

import (
	"testing"
	"time"
)

func TestStatusIntegration(t *testing.T) {
	adapter := NewAdapter()
	status := adapter.Status()
	if !status.Supported {
		t.Skip("host integration only applies on macOS")
	}
	if !status.OSAScriptAvailable {
		t.Skip("osascript is unavailable on this runner")
	}
	if status.IllustratorAppFound && status.AppPath == "" {
		t.Fatal("expected resolved app path when Illustrator is installed")
	}
}

func TestExecuteJavaScriptIntegration(t *testing.T) {
	adapter := NewAdapter()
	status := adapter.Status()
	if !status.IllustratorAppFound {
		t.Skip("Illustrator not installed")
	}
	if !status.IllustratorRunning {
		t.Skip("Illustrator not running; skipping live script execution")
	}

	out, err := adapter.ExecuteJavaScript(`(function () { return "ok"; }())`, 15*time.Second)
	if err != nil {
		t.Fatalf("execute javascript: %v", err)
	}
	if out != "ok" {
		t.Fatalf("unexpected script output: %s", out)
	}
}
