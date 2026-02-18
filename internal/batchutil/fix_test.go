package batchutil_test

import (
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

func TestFixBatchOps_TypeToCommand(t *testing.T) {
	in := []byte(`[{"type":"frame","params":{}}]`)
	fixed, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) == 0 {
		t.Fatal("expected fixes")
	}
	if string(fixed) == string(in) {
		t.Fatal("data unchanged")
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
	if len(fixed) == 0 {
		t.Fatal("empty output")
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
