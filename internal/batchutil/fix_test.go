package batchutil_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/nerveband/ai-happy-design/internal/batchutil"
)

func TestStripMarkdownFences(t *testing.T) {
	in := "```json\n[{\"command\":\"frame\",\"params\":{}}]\n```"
	out := batchutil.StripMarkdownFences([]byte(in))
	if string(out) != "[{\"command\":\"frame\",\"params\":{}}]" {
		t.Fatalf("got %q", string(out))
	}
}

func TestStripMarkdownFences_TrailingNewline(t *testing.T) {
	// Model often appends a trailing newline after the closing fence
	in := "```json\n[{\"command\":\"frame\",\"params\":{}}]\n```\n"
	out := batchutil.StripMarkdownFences([]byte(in))
	got := strings.TrimSpace(string(out))
	want := `[{"command":"frame","params":{}}]`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFixBatchOps_TypeToCommand(t *testing.T) {
	in := []byte(`[{"type":"frame","params":{}}]`)
	fixed, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) == 0 {
		t.Fatal("expected fixes")
	}
	var ops []map[string]interface{}
	if err := json.Unmarshal(fixed, &ops); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, hasType := ops[0]["type"]; hasType {
		t.Fatal("output still has 'type' field")
	}
	if cmd, _ := ops[0]["command"].(string); cmd != "frame" {
		t.Fatalf("expected command=frame, got %q", cmd)
	}
}

func TestFixBatchOps_HoistTopLevelProps(t *testing.T) {
	in := []byte(`[{"command":"frame","x":0,"y":0,"params":{}}]`)
	_, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) == 0 {
		t.Fatal("expected hoisting fix")
	}
}

func TestFixBatchOps_UnwrapDict(t *testing.T) {
	in := []byte(`{"ops":[{"command":"frame","params":{}}]}`)
	fixed, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) == 0 {
		t.Fatal("expected unwrap fix")
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal(fixed, &arr); err != nil {
		t.Fatalf("output is not a valid JSON array: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("expected 1 element, got %d", len(arr))
	}
}

func TestFixBatchOps_NoChanges(t *testing.T) {
	in := []byte(`[{"name":"bg","command":"frame","params":{"x":0,"y":0}}]`)
	_, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) != 0 {
		t.Fatalf("expected no fixes but got: %v", fixes)
	}
}
