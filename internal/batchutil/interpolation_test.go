package batchutil

import (
	"strings"
	"testing"
)

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

// --- New tests below ---

func TestInterpolateParams_NestedResultReference(t *testing.T) {
	// Result contains a nested map
	steps := []StepState{
		{Index: 0, Name: "frame1", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{
				"id": "10:20",
				"bounds": map[string]interface{}{
					"x":      100.0,
					"y":      200.0,
					"width":  1080.0,
					"height": 1350.0,
				},
			},
		},
	}
	params := map[string]interface{}{
		"parentId": "${{steps.frame1.result.id}}",
		"x":        "${{steps.frame1.result.bounds.x}}",
		"width":    "${{steps.frame1.result.bounds.width}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["parentId"]; got != "10:20" {
		t.Fatalf("unexpected parentId: %#v", got)
	}
	// When the whole string is a single placeholder, type is preserved
	if got, ok := out["x"].(float64); !ok || got != 100.0 {
		t.Fatalf("unexpected x: %#v (type %T)", out["x"], out["x"])
	}
	if got, ok := out["width"].(float64); !ok || got != 1080.0 {
		t.Fatalf("unexpected width: %#v (type %T)", out["width"], out["width"])
	}
}

func TestInterpolateParams_TypePreservation_Bool(t *testing.T) {
	steps := []StepState{
		{Index: 0, Command: "node.get_info", OK: true,
			Result: map[string]interface{}{"visible": true, "locked": false}},
	}
	params := map[string]interface{}{
		"visible": "${{steps.0.result.visible}}",
		"locked":  "${{steps.0.result.locked}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, ok := out["visible"].(bool); !ok || got != true {
		t.Fatalf("expected visible=true, got %#v (%T)", out["visible"], out["visible"])
	}
	if got, ok := out["locked"].(bool); !ok || got != false {
		t.Fatalf("expected locked=false, got %#v (%T)", out["locked"], out["locked"])
	}
}

func TestInterpolateParams_TypePreservation_NumberStaysNumber(t *testing.T) {
	steps := []StepState{
		{Index: 0, Command: "node.get_info", OK: true,
			Result: map[string]interface{}{"width": 1080.0, "height": 1350.0}},
	}
	params := map[string]interface{}{
		"width":  "${{steps.0.result.width}}",
		"height": "${{steps.0.result.height}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, ok := out["width"].(float64); !ok || got != 1080.0 {
		t.Fatalf("expected width=1080, got %#v (%T)", out["width"], out["width"])
	}
	if got, ok := out["height"].(float64); !ok || got != 1350.0 {
		t.Fatalf("expected height=1350, got %#v (%T)", out["height"], out["height"])
	}
}

func TestInterpolateParams_MissingStepReference(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"parentId": "${{steps.nonexistent.result.id}}",
	}

	_, err := InterpolateParams(params, steps)
	if err == nil {
		t.Fatal("expected error for missing step reference")
	}
	// Error should mention the path
	if got := err.Error(); got == "" {
		t.Fatal("error message should not be empty")
	}
}

func TestInterpolateParams_MissingFieldInResult(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"nodeId": "${{steps.frame.result.nonexistent_field}}",
	}

	_, err := InterpolateParams(params, steps)
	if err == nil {
		t.Fatal("expected error for missing field in result")
	}
}

func TestInterpolateParams_MultipleInterpolationsInOneString(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "a", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "10:20", "name": "Frame-A"}},
		{Index: 1, Name: "b", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "10:30", "name": "Frame-B"}},
	}
	params := map[string]interface{}{
		"label": "parent=${{steps.a.result.id}},child=${{steps.b.result.id}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["label"]; got != "parent=10:20,child=10:30" {
		t.Fatalf("unexpected label: %#v", got)
	}
}

func TestInterpolateParams_MultipleInterpolationsWithNonStringTypes(t *testing.T) {
	// When multiple placeholders exist in a string, non-string values get stringified
	steps := []StepState{
		{Index: 0, Name: "info", Command: "node.get_info", OK: true,
			Result: map[string]interface{}{"x": 100.0, "y": 200.0}},
	}
	params := map[string]interface{}{
		"position": "x=${{steps.info.result.x}},y=${{steps.info.result.y}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["position"]; got != "x=100,y=200" {
		t.Fatalf("unexpected position: %#v", got)
	}
}

func TestInterpolateParams_EmptyParams(t *testing.T) {
	steps := []StepState{
		{Index: 0, Command: "noop", OK: true, Result: map[string]interface{}{}},
	}
	params := map[string]interface{}{}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected empty params, got %v", out)
	}
}

func TestInterpolateParams_NoPlaceholders(t *testing.T) {
	steps := []StepState{
		{Index: 0, Command: "noop", OK: true, Result: map[string]interface{}{}},
	}
	params := map[string]interface{}{
		"name":  "TestFrame",
		"width": 1080.0,
		"x":     0.0,
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["name"] != "TestFrame" {
		t.Fatalf("expected name=TestFrame, got %v", out["name"])
	}
	if out["width"] != 1080.0 {
		t.Fatalf("expected width=1080, got %v", out["width"])
	}
}

func TestInterpolateParams_LastReference(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "first", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "10:20"}},
		{Index: 1, Name: "second", Command: "text.create", OK: true,
			Result: map[string]interface{}{"id": "10:30"}},
	}
	params := map[string]interface{}{
		"nodeId": "${{last.result.id}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["nodeId"]; got != "10:30" {
		t.Fatalf("expected last result id=10:30, got %#v", got)
	}
}

func TestInterpolateParams_NestedParamsMap(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"fills": map[string]interface{}{
			"nodeId": "${{steps.frame.result.id}}",
			"type":   "SOLID",
		},
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fills := out["fills"].(map[string]interface{})
	if fills["nodeId"] != "10:20" {
		t.Fatalf("expected nested nodeId=10:20, got %v", fills["nodeId"])
	}
	if fills["type"] != "SOLID" {
		t.Fatalf("expected type=SOLID, got %v", fills["type"])
	}
}

func TestInterpolateParams_ArrayInParams(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"nodeIds": []interface{}{
			"${{steps.frame.result.id}}",
			"static-id",
		},
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	nodeIds := out["nodeIds"].([]interface{})
	if nodeIds[0] != "10:20" {
		t.Fatalf("expected first nodeId=10:20, got %v", nodeIds[0])
	}
	if nodeIds[1] != "static-id" {
		t.Fatalf("expected second nodeId=static-id, got %v", nodeIds[1])
	}
}

func TestInterpolateParams_ShortSyntax_ByNameAndLast(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "node.create_frame", OK: true, Result: map[string]interface{}{"id": "10:20", "name": "Hero"}},
		{Index: 1, Name: "title", Command: "text.create", OK: true, Result: map[string]interface{}{"id": "10:21"}},
	}
	params := map[string]interface{}{
		"parentId": "$frame",
		"label":    "$frame.name",
		"lastId":   "$last",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["parentId"] != "10:20" {
		t.Fatalf("expected parentId=10:20, got %#v", out["parentId"])
	}
	if out["label"] != "Hero" {
		t.Fatalf("expected label=Hero, got %#v", out["label"])
	}
	if out["lastId"] != "10:21" {
		t.Fatalf("expected lastId=10:21, got %#v", out["lastId"])
	}
}

func TestInterpolateParams_ShortSyntax_MixedWithLongForm(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "node.create_frame", OK: true, Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"meta": "a=$frame,b=${{steps.frame.result.id}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["meta"] != "a=10:20,b=10:20" {
		t.Fatalf("unexpected meta: %#v", out["meta"])
	}
}

func TestInterpolateParams_CaseInsensitiveStepName(t *testing.T) {
	// Step names are lowercased by SanitizeStepName, but LLMs often reference
	// them in camelCase. The lookup should be case-insensitive.
	steps := []StepState{
		{Index: 0, Name: "createpage", Command: "page.create", OK: true,
			Result: map[string]interface{}{"id": "1:2"}},
		{Index: 1, Name: "emailnewsletter_frame", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "3:4"}},
	}

	params := map[string]interface{}{
		"pageId":   "${{steps.createPage.result.id}}",
		"parentId": "${{steps.emailNewsletter_frame.result.id}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["pageId"]; got != "1:2" {
		t.Fatalf("expected pageId=1:2, got %#v", got)
	}
	if got := out["parentId"]; got != "3:4" {
		t.Fatalf("expected parentId=3:4, got %#v", got)
	}
}

func TestInterpolateParams_CaseInsensitiveShortSyntax(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "myframe", Command: "node.create_frame", OK: true,
			Result: map[string]interface{}{"id": "5:6"}},
	}

	params := map[string]interface{}{
		"parentId": "$myFrame",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["parentId"]; got != "5:6" {
		t.Fatalf("expected parentId=5:6, got %#v", got)
	}
}

func TestInterpolateParams_MissingStepErrorIncludesSuggestion(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "create_page", Command: "page.create", OK: true, Result: map[string]interface{}{"id": "1:2"}},
	}
	params := map[string]interface{}{
		"pageId": "${{steps.createPage.result.id}}",
	}

	_, err := InterpolateParams(params, steps)
	if err == nil {
		t.Fatal("expected interpolation error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "available step names: create_page") {
		t.Fatalf("expected available step names in error, got: %s", msg)
	}
	if !strings.Contains(msg, "${{steps.create_page.result.id}}") {
		t.Fatalf("expected likely fix in error, got: %s", msg)
	}
}

func TestInterpolateParams_ShortSyntax_DoesNotRewriteCurrencyOrLongForm(t *testing.T) {
	steps := []StepState{
		{Index: 0, Name: "frame", Command: "node.create_frame", OK: true, Result: map[string]interface{}{"id": "10:20"}},
	}
	params := map[string]interface{}{
		"amount": "price is $100",
		"raw":    "${{steps.frame.result.id}}",
	}

	out, err := InterpolateParams(params, steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["amount"] != "price is $100" {
		t.Fatalf("expected currency string untouched, got %#v", out["amount"])
	}
	if out["raw"] != "10:20" {
		t.Fatalf("expected long form interpolation to resolve, got %#v", out["raw"])
	}
}
