//go:build darwin && integration

package bridge

import (
	"testing"
	"time"

	"github.com/nerveband/ai-happy-design/internal/illustrator/host"
)

func TestExecuteScriptIntegration(t *testing.T) {
	adapter := host.NewAdapter()
	status := adapter.Status()
	if !status.IllustratorAppFound {
		t.Skip("Illustrator not installed")
	}
	if !status.IllustratorRunning {
		t.Skip("Illustrator not running")
	}

	client := NewClient(adapter)
	response, err := client.ExecuteScript(`(function () { return { version: app.version, name: app.name }; }())`, 15*time.Second)
	if err != nil {
		t.Fatalf("execute script: %v", err)
	}
	result, ok := response.Result.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %#v", response.Result)
	}
	if result["version"] == "" {
		t.Fatalf("expected version in result, got %#v", result)
	}
}

func TestProbePluginIntegration(t *testing.T) {
	adapter := host.NewAdapter()
	status := adapter.Status()
	if !status.IllustratorAppFound {
		t.Skip("Illustrator not installed")
	}
	if !status.IllustratorRunning {
		t.Skip("Illustrator not running")
	}

	probe := NewClient(adapter).ProbePlugin(5 * time.Second)
	if !probe.Reachable && probe.Code != "PLUGIN_REQUIRED" {
		t.Fatalf("expected plugin probe to be reachable or deterministically plugin-required, got %+v", probe)
	}
}
