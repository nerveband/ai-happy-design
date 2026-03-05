//go:build darwin && integration

package host

import (
	"strings"
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
}

func TestExecuteJavaScriptIntegration(t *testing.T) {
	adapter := NewAdapter()
	status := adapter.Status()
	if !status.IllustratorAppFound {
		t.Skip("Adobe Illustrator.app not installed")
	}
	if !status.IllustratorRunning {
		t.Skip("Illustrator not running; skipping live script execution")
	}

	out, err := adapter.ExecuteJavaScript(`(function () { return JSON.stringify({ "ok": true }); }())`, 5*time.Second)
	if err != nil {
		t.Fatalf("execute javascript: %v", err)
	}
	if !strings.Contains(out, `"ok": true`) && !strings.Contains(out, `"ok":true`) {
		t.Fatalf("unexpected script output: %s", out)
	}
}
