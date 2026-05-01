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
		{"command": "text.crate", "params": map[string]interface{}{}},
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

func TestInterpolationReferenceSkipsPatternValidation(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "paint.set_solid", "params": map[string]interface{}{
			"nodeId": "${{steps.frame.result.id}}", "color": "#FF0000",
		}},
	}
	result := ValidateBatch(ops)
	for _, w := range result.Warnings {
		if w.Code == "PATTERN_MISMATCH" && w.Param == "nodeId" {
			t.Fatalf("expected interpolation nodeId to pass pre-execution pattern validation: %+v", w)
		}
	}
}

func TestFigmaRGBColorObjectAutoConvertsToHex(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "paint.set_solid", "params": map[string]interface{}{
			"nodeId": "1:23",
			"color":  map[string]interface{}{"r": 1.0, "g": 0.5, "b": 0.0},
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "TYPE_MISMATCH" && w.Param == "color" {
			found = true
			if w.Fix != "#FF8000" || !w.Applied {
				t.Fatalf("expected applied #FF8000 fix, got %+v", w)
			}
		}
	}
	if !found {
		t.Fatal("expected color object conversion warning")
	}
	params := ops[0]["params"].(map[string]interface{})
	if params["color"] != "#FF8000" {
		t.Fatalf("expected color param converted, got %v", params["color"])
	}
}

func TestHighLevelCommandValidation(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.measure", "params": map[string]interface{}{
			"text": "Measure this", "width": 320.0, "fontSize": 18.0,
		}},
		{"command": "text.fit_box", "params": map[string]interface{}{
			"text": "Fit this", "width": 240.0, "height": 80.0, "minFontSize": 12.0, "maxFontSize": 28.0,
		}},
		{"command": "text.rich", "params": map[string]interface{}{
			"pid": "1:23", "width": 360.0, "heading": "Gold", "price": "$500", "bullets": []interface{}{"Logo", "Table"},
		}},
		{"command": "layout.pricing_grid", "params": map[string]interface{}{
			"pid": "1:23", "width": 960.0, "columns": 3.0, "cards": []interface{}{
				map[string]interface{}{"title": "Sponsor", "price": "$500"},
			},
		}},
	}

	result := ValidateBatch(ops)
	if len(result.Errors) > 0 {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}
	if len(result.Warnings) > 0 {
		t.Fatalf("expected no warnings, got %+v", result.Warnings)
	}
}

func TestHighLevelCommandValidationRequiresCoreParams(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.measure", "params": map[string]interface{}{}},
		{"command": "text.fit_box", "params": map[string]interface{}{"text": "Fit"}},
		{"command": "text.create_rich_block", "params": map[string]interface{}{"heading": "Gold"}},
		{"command": "layout.pricing_grid", "params": map[string]interface{}{"width": 960.0}},
	}

	result := ValidateBatch(ops)
	required := map[string]bool{
		"text.measure:text":            false,
		"text.fit_box:width":           false,
		"text.fit_box:height":          false,
		"text.create_rich_block:width": false,
		"layout.pricing_grid:cards":    false,
	}

	for _, issue := range result.Warnings {
		if issue.Code != "REQUIRED_MISSING" {
			continue
		}
		cmd, _ := ops[issue.Step]["command"].(string)
		key := cmd + ":" + issue.Param
		if _, ok := required[key]; ok {
			required[key] = true
		}
	}

	for key, found := range required {
		if !found {
			t.Fatalf("expected REQUIRED_MISSING warning for %s; got %+v", key, result.Warnings)
		}
	}
}
