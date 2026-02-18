package main

import (
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
