package bridge

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeHost struct {
	responses []fakeResponse
	scripts   []string
}

type fakeResponse struct {
	raw string
	err error
}

func (f *fakeHost) ExecuteJavaScript(script string, timeout time.Duration) (string, error) {
	f.scripts = append(f.scripts, script)
	if len(f.responses) == 0 {
		return "", nil
	}
	next := f.responses[0]
	f.responses = f.responses[1:]
	return next.raw, next.err
}

func TestBuildRuntimeScriptUsesCustomSerializer(t *testing.T) {
	t.Parallel()

	script, err := buildRuntimeScript(runtimePayload{
		ID:     "req-1",
		Mode:   "script",
		Script: `(function () { return { ok: true }; }())`,
	})
	if err != nil {
		t.Fatalf("build runtime script: %v", err)
	}
	if !strings.Contains(script, "function ahdStringify") {
		t.Fatal("expected custom ahdStringify helper in runtime bridge")
	}
	if strings.Contains(script, "JSON.stringify") {
		t.Fatal("bridge script should not depend on JSON.stringify")
	}
}

func TestProbePluginClassifiesIllustratorSelectorError(t *testing.T) {
	t.Parallel()

	client := &Client{host: &fakeHost{
		responses: []fakeResponse{{
			raw: `{"v":"1.0","id":"probe-1","ok":false,"error":"Error: an Illustrator error occurred: 1344357988 ('dF!P')"}`,
		}},
	}}

	probe := client.ProbePlugin(2 * time.Second)
	if probe.Reachable {
		t.Fatal("expected probe to report plugin unreachable")
	}
	if got, want := probe.Code, "PLUGIN_REQUIRED"; got != want {
		t.Fatalf("unexpected probe code %q, want %q", got, want)
	}
}

func TestExecuteSelectorFailsFastWhenPluginMissing(t *testing.T) {
	t.Parallel()

	host := &fakeHost{
		responses: []fakeResponse{{
			raw: `{"v":"1.0","id":"probe-1","ok":false,"error":"Error: an Illustrator error occurred: 1344357988 ('dF!P')"}`,
		}},
	}
	client := &Client{host: host}

	_, err := client.ExecuteSelector("ahd.inspect", Request{Command: "inspect.summary"}, 3*time.Second)
	if err == nil {
		t.Fatal("expected selector execution to fail without plugin")
	}
	if got, want := err.Code, "PLUGIN_REQUIRED"; got != want {
		t.Fatalf("unexpected error code %q, want %q", got, want)
	}
	if len(host.scripts) != 1 {
		t.Fatalf("expected only the plugin probe to run, got %d host calls", len(host.scripts))
	}
}

func TestClassifyHostExecutionErrorKeepsScriptFailuresAsHostErrors(t *testing.T) {
	t.Parallel()

	err := classifyHostExecutionError(errors.New("HOST_EXEC_ERROR: some script issue"), false, "")
	if got, want := err.Code, "HOST_EXEC_ERROR"; got != want {
		t.Fatalf("unexpected code %q, want %q", got, want)
	}
}
