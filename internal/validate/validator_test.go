package validate

import (
	"testing"

	_ "github.com/nerveband/ai-happy-design/internal/schema" // register schemas
)

func TestValidParam(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "fontSize": 48.0,
		}},
	}
	result := ValidateBatch(ops)
	if len(result.Errors) > 0 {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}
}

func TestBelowMin(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "fontSize": -10.0,
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "BELOW_MIN" && w.Param == "fontSize" {
			found = true
			if w.Fix != 4.0 {
				t.Errorf("expected fix=4, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected BELOW_MIN warning for fontSize")
	}
}

func TestEnumFuzzy(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "fontStyle": "bold",
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "ENUM_INVALID" && w.Param == "fontStyle" {
			found = true
			if w.Fix != "Bold" {
				t.Errorf("expected fix=Bold, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected ENUM_INVALID warning with fuzzy fix")
	}
}

func TestNamedColor(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "paint.set_solid", "params": map[string]interface{}{
			"nodeId": "1:23", "color": "red",
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "PATTERN_MISMATCH" && w.Param == "color" {
			found = true
			if w.Fix != "#FF0000" {
				t.Errorf("expected fix=#FF0000, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected PATTERN_MISMATCH warning with color fix")
	}
}

func TestUnknownCommand(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.creat", "params": map[string]interface{}{}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "UNKNOWN_COMMAND" {
			found = true
			if w.Fix != "text.create" {
				t.Errorf("expected fix=text.create, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected UNKNOWN_COMMAND warning")
	}
}

func TestControlChars(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello\x00World", "parentId": "1:23",
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "CONTROL_CHARS" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected CONTROL_CHARS warning for null byte in text")
	}
}

func TestControlCharsAllowsNewlines(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello\nWorld\ttab", "parentId": "1:23",
		}},
	}
	result := ValidateBatch(ops)
	for _, w := range result.Warnings {
		if w.Code == "CONTROL_CHARS" {
			t.Fatal("should not flag newlines/tabs as control chars")
		}
	}
}

func TestPathTraversal(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "../../etc/passwd", "parentId": "1:23",
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "PATH_TRAVERSAL" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected PATH_TRAVERSAL warning for .. sequence")
	}
}

func TestPathTraversalPercentEncoded(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "%2e%2e%2fetc%2fpasswd", "parentId": "1:23",
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "PATH_TRAVERSAL" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected PATH_TRAVERSAL warning for percent-encoded traversal")
	}
}

func TestDependencyAutoFix(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "lineHeight": 150.0,
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "DEPENDENCY_MISSING" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected DEPENDENCY_MISSING warning for lineHeightUnit")
	}
	// Check it was auto-applied
	params := ops[0]["params"].(map[string]interface{})
	if params["lineHeightUnit"] != "PERCENT" {
		t.Errorf("expected lineHeightUnit auto-set to PERCENT, got %v", params["lineHeightUnit"])
	}
}
