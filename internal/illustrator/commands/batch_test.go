package commands

import (
	"fmt"
	"testing"

	"github.com/nerveband/ai-happy-design/internal/commoncli"
)

func TestResolveBatchOpInterpolatesStepResults(t *testing.T) {
	t.Parallel()

	op := commoncli.BatchOp{
		Name:    "second",
		Command: "document.info",
		Params: map[string]any{
			"documentId": "${{steps.first.result.params.width}}",
		},
	}
	resolved, err := ResolveBatchOp(op, []commoncli.BatchStep{
		{
			Name: "first",
			Result: map[string]any{
				"params": map[string]any{"width": 1200},
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected interpolation error: %v", err)
	}
	if fmt.Sprint(resolved.Params["documentId"]) != "1200" {
		t.Fatalf("expected numeric interpolation, got %#v", resolved.Params["documentId"])
	}
}

func TestResolveBatchOpRejectsMissingPaths(t *testing.T) {
	t.Parallel()

	op := commoncli.BatchOp{
		Name:    "second",
		Command: "export.png",
		Params: map[string]any{
			"outputPath": "${{steps.first.result.name}}",
		},
	}
	_, err := ResolveBatchOp(op, []commoncli.BatchStep{{Name: "first", Result: map[string]any{}}})
	if err == nil {
		t.Fatal("expected missing result path to fail")
	}
}
