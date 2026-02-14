package batchutil

import "testing"

func TestInterpolateParams_ByIndexAndName(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "shape.create_rectangle", OK: true, Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"nodeId": "${{steps.0.result.id}}",
		"meta":   "rect-${{steps.frame.result.id}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, ok := out["nodeId"].(string); !ok || got != "10:20" {
		t.Fatalf("unexpected nodeId: %#v", out["nodeId"])
	}
	if got, ok := out["meta"].(string); !ok || got != "rect-10:20" {
		t.Fatalf("unexpected meta: %#v", out["meta"])
	}
}

func TestInterpolateParams_TypedValue(t *testing.T) {
	steps := []StepState{
		{Index: 0, Command: "shape.create_rectangle", OK: true, Result: map[string]interface{}{"opacity": 0.75}},
	}
	params := map[string]interface{}{
		"opacity": "${{steps.0.result.opacity}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, ok := out["opacity"].(float64); !ok || got != 0.75 {
		t.Fatalf("unexpected opacity: %#v", out["opacity"])
	}
}

func TestInterpolateParams_PathNotFound(t *testing.T) {
	steps := []StepState{{Index: 0, Command: "noop", OK: false}}
	params := map[string]interface{}{
		"nodeId": "${{steps.0.result.id}}",
	}

	_, err := InterpolateParams(params, steps)
	if err == nil {
		t.Fatal("expected interpolation error")
	}
}

func TestInterpolateParams_BracketPath(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "created", Command: "shape.create_rectangle", OK: true, Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"nodeId": "${{steps[0].result.id}}",
		"label":  "id=${{steps[created].result.id}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["nodeId"]; got != "10:20" {
		t.Fatalf("unexpected nodeId: %#v", got)
	}
	if got := out["label"]; got != "id=10:20" {
		t.Fatalf("unexpected label: %#v", got)
	}
}
