package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateBatchOps_Valid(t *testing.T) {
	input := `[
		{"name":"bg","command":"node.create_frame","params":{"x":0,"y":0,"width":1080,"height":1080}},
		{"name":"title","command":"text.create","params":{"text":"Hello","parentId":"${{steps.bg.result.id}}"}}
	]`
	errs := validateBatchOps([]byte(input))
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateBatchOps_TypeInsteadOfCommand(t *testing.T) {
	input := `[{"name":"bg","type":"node.create_frame","params":{"x":0}}]`
	errs := validateBatchOps([]byte(input))
	if len(errs) == 0 {
		t.Error("expected error for 'type' field, got none")
	}
}

func TestValidateBatchOps_TopLevelParams(t *testing.T) {
	input := `[{"name":"bg","command":"node.create_frame","params":{},"x":0,"color":"#fff"}]`
	errs := validateBatchOps([]byte(input))
	if len(errs) < 2 {
		t.Errorf("expected 2+ errors (x, color at top level), got %d: %v", len(errs), errs)
	}
}

func TestValidateBatchOps_BrokenInterpolation(t *testing.T) {
	// Step named "bg" but reference uses "background"
	input := `[
		{"name":"bg","command":"node.create_frame","params":{"x":0}},
		{"name":"title","command":"text.create","params":{"parentId":"${{steps.background.result.id}}"}}
	]`
	errs := validateBatchOps([]byte(input))
	found := false
	for _, e := range errs {
		if strings.Contains(e, "background") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for undefined step 'background', got: %v", errs)
	}
}

func TestValidateBatchOps_MissingParams(t *testing.T) {
	input := `[{"name":"bg","command":"node.create_frame"}]`
	errs := validateBatchOps([]byte(input))
	if len(errs) == 0 {
		t.Error("expected error for missing params, got none")
	}
}

func TestValidateBatchOps_InvalidJSON(t *testing.T) {
	errs := validateBatchOps([]byte("not json"))
	if len(errs) == 0 || !strings.Contains(errs[0], "Invalid JSON") {
		t.Errorf("expected Invalid JSON error, got: %v", errs)
	}
}

func TestValidateBatchOps_ShortAliasCommands(t *testing.T) {
	// Short aliases (frame, text, rect) are valid commands
	input := `[
		{"name":"bg","command":"frame","params":{"x":0,"y":0,"w":1080,"h":1080}},
		{"name":"title","command":"text","params":{"text":"Hi","pid":"${{steps.bg.result.id}}"}}
	]`
	errs := validateBatchOps([]byte(input))
	if len(errs) != 0 {
		t.Errorf("short alias commands should be valid, got: %v", errs)
	}
}

func TestStripMarkdownFences_WithFences(t *testing.T) {
	input := "```json\n[{\"name\":\"bg\"}]\n```"
	result := stripMarkdownFences([]byte(input))
	if strings.Contains(string(result), "```") {
		t.Errorf("expected fences stripped, got: %s", result)
	}
	if !strings.Contains(string(result), `"name"`) {
		t.Errorf("expected JSON content preserved, got: %s", result)
	}
}

func TestStripMarkdownFences_WithoutFences(t *testing.T) {
	input := `[{"name":"bg","command":"frame","params":{}}]`
	result := stripMarkdownFences([]byte(input))
	if string(result) != input {
		t.Errorf("expected unchanged input, got: %s", result)
	}
}

func TestFixBatchOps_TypeToCommand(t *testing.T) {
	input := `[{"name":"bg","type":"frame","params":{"x":0}}]`
	fixed, fixes, err := fixBatchOps([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fixes) == 0 {
		t.Error("expected at least one fix")
	}
	var ops []map[string]interface{}
	json.Unmarshal(fixed, &ops)
	if _, hasType := ops[0]["type"]; hasType {
		t.Error("expected 'type' to be removed")
	}
	if cmd, ok := ops[0]["command"].(string); !ok || cmd != "frame" {
		t.Errorf("expected command=frame, got: %v", ops[0]["command"])
	}
}

func TestFixBatchOps_HoistTopLevelProps(t *testing.T) {
	input := `[{"name":"bg","command":"frame","params":{},"x":10,"y":20,"color":"#fff"}]`
	fixed, fixes, err := fixBatchOps([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fixes) == 0 {
		t.Error("expected fixes for hoisted props")
	}
	var ops []map[string]interface{}
	json.Unmarshal(fixed, &ops)
	params := ops[0]["params"].(map[string]interface{})
	if params["x"] == nil || params["y"] == nil || params["color"] == nil {
		t.Errorf("expected x, y, color in params, got: %v", params)
	}
	if ops[0]["x"] != nil || ops[0]["color"] != nil {
		t.Error("expected x and color removed from top level")
	}
}

func TestFixBatchOps_StripFencesThenFix(t *testing.T) {
	// Model output with both fences AND type error
	input := "```json\n[{\"name\":\"bg\",\"type\":\"frame\",\"x\":0,\"params\":{}}]\n```"
	fixed, fixes, err := fixBatchOps([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fixes) == 0 {
		t.Error("expected fixes applied")
	}
	// After fix, should be valid
	errs := validateBatchOps(fixed)
	if len(errs) != 0 {
		t.Errorf("expected valid after fix, got errors: %v", errs)
	}
}

func TestFixBatchOps_UnwrapDictWrapper(t *testing.T) {
	// Model outputs {"ops": [...]} instead of bare [...]
	input := `{"ops":[{"name":"bg","command":"frame","params":{"x":0}}]}`
	fixed, fixes, err := fixBatchOps([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fixes) == 0 {
		t.Error("expected fix for dict unwrap")
	}
	errs := validateBatchOps(fixed)
	if len(errs) != 0 {
		t.Errorf("expected valid after unwrap, got: %v", errs)
	}
}

func TestValidateWithFix_EndToEnd(t *testing.T) {
	// Simulates the full --fix workflow: bad input → fixed → valid
	bad := `[{"name":"bg","type":"frame","x":0,"y":0,"color":"#1a1a1a","params":{}}]`
	fixed, fixes, err := fixBatchOps([]byte(bad))
	if err != nil {
		t.Fatalf("fixBatchOps error: %v", err)
	}
	if len(fixes) == 0 {
		t.Error("expected fixes")
	}
	errs := validateBatchOps(fixed)
	if len(errs) != 0 {
		t.Errorf("expected clean after fix, got: %v", errs)
	}
}
