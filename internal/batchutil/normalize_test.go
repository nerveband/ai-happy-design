package batchutil

import "testing"

func TestNormalizeBatchParams_GenericAliases(t *testing.T) {
	in := map[string]interface{}{
		"w":   1200.0,
		"h":   628.0,
		"pid": "$frame",
	}
	out := NormalizeBatchParams("rect", in)

	if out["width"] != 1200.0 {
		t.Fatalf("expected width alias expansion, got %#v", out["width"])
	}
	if out["height"] != 628.0 {
		t.Fatalf("expected height alias expansion, got %#v", out["height"])
	}
	if out["parentId"] != "$frame" {
		t.Fatalf("expected parentId alias expansion, got %#v", out["parentId"])
	}
}

func TestNormalizeBatchParams_CanonicalWins(t *testing.T) {
	in := map[string]interface{}{
		"width": 900.0,
		"w":     1200.0,
	}
	out := NormalizeBatchParams("rect", in)
	if out["width"] != 900.0 {
		t.Fatalf("expected canonical width to win, got %#v", out["width"])
	}
}

func TestNormalizeBatchParams_TextAliasesAndDefaults(t *testing.T) {
	in := map[string]interface{}{
		"sz": 24.0,
		"ff": "Inter",
		"fs": "Bold",
		"lh": 140.0,
		"ls": 2.0,
	}
	out := NormalizeBatchParams("text", in)

	if out["fontSize"] != 24.0 || out["fontFamily"] != "Inter" || out["fontStyle"] != "Bold" {
		t.Fatalf("expected text aliases to expand, got %#v", out)
	}
	if out["lineHeight"] != 140.0 || out["letterSpacing"] != 2.0 {
		t.Fatalf("expected spacing aliases to expand, got %#v", out)
	}
	if out["lineHeightUnit"] != "PERCENT" {
		t.Fatalf("expected lineHeightUnit default PERCENT, got %#v", out["lineHeightUnit"])
	}
}

func TestNormalizeBatchParams_StrokeAliasByCommand(t *testing.T) {
	createOut := NormalizeBatchParams("ellipse", map[string]interface{}{"sw": 3.0})
	if createOut["strokeWidth"] != 3.0 {
		t.Fatalf("expected create command sw->strokeWidth, got %#v", createOut["strokeWidth"])
	}

	strokeOut := NormalizeBatchParams("stroke", map[string]interface{}{"sw": 4.0})
	if strokeOut["strokeWeight"] != 4.0 {
		t.Fatalf("expected stroke command sw->strokeWeight, got %#v", strokeOut["strokeWeight"])
	}
}

func TestNormalizeBatchParams_FrameRectOnlyAliases(t *testing.T) {
	frameOut := NormalizeBatchParams("frame", map[string]interface{}{"bg": "#111111", "r": 16.0})
	if frameOut["color"] != "#111111" {
		t.Fatalf("expected bg->color on frame create, got %#v", frameOut["color"])
	}
	if frameOut["cornerRadius"] != 16.0 {
		t.Fatalf("expected r->cornerRadius on frame create, got %#v", frameOut["cornerRadius"])
	}

	textOut := NormalizeBatchParams("text", map[string]interface{}{"r": 16.0})
	if _, ok := textOut["cornerRadius"]; ok {
		t.Fatalf("did not expect cornerRadius alias for text command: %#v", textOut)
	}
}

func TestNormalizeBatchParams_NilParams(t *testing.T) {
	out := NormalizeBatchParams("text", nil)
	if out == nil {
		t.Fatal("expected non-nil params map")
	}
	if len(out) != 0 {
		t.Fatalf("expected empty params map, got %#v", out)
	}
}
