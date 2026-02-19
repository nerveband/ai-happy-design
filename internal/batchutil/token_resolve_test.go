package batchutil

import (
	"math"
	"testing"
)

func TestResolveTokenAliases_AppliesKnownAliases(t *testing.T) {
	ops := []map[string]interface{}{
		{
			"name":    "root",
			"command": "frame",
			"params": map[string]interface{}{
				"name": "Design",
				"w":    1080,
				"h":    1350,
			},
		},
		{
			"name":    "content",
			"command": "frame",
			"params": map[string]interface{}{
				"pid":         "$root",
				"padding":     "side",
				"itemSpacing": "section",
				"r":           "card",
				"w":           "content",
			},
		},
		{
			"name":    "title",
			"command": "text",
			"params": map[string]interface{}{
				"pid": "$content",
				"sz":  "hero",
			},
		},
		{
			"name":    "cta",
			"command": "frame",
			"params": map[string]interface{}{
				"pid": "$content",
				"r":   "button",
			},
		},
		{
			"name":    "pill",
			"command": "frame",
			"params": map[string]interface{}{
				"pid":          "$content",
				"cornerRadius": "pill",
			},
		},
	}

	applied, rootWidth := ResolveTokenAliases(ops)
	if rootWidth != 1080 {
		t.Fatalf("root width = %d, want 1080", rootWidth)
	}
	if applied == 0 {
		t.Fatal("expected aliases to be resolved")
	}

	content := ops[1]["params"].(map[string]interface{})
	if got := asInt(content["padding"]); got != 72 {
		t.Fatalf("padding(side) = %d, want 72", got)
	}
	if got := asInt(content["itemSpacing"]); got != 48 {
		t.Fatalf("itemSpacing(section) = %d, want 48", got)
	}
	if got := asInt(content["r"]); got != 16 {
		t.Fatalf("r(card) = %d, want 16", got)
	}
	if got := asInt(content["w"]); got != 936 {
		t.Fatalf("w(content) = %d, want 936", got)
	}

	title := ops[2]["params"].(map[string]interface{})
	wantHero := int(math.Round(tokenSizes(1080)["hero"]))
	if got := asInt(title["sz"]); got != wantHero {
		t.Fatalf("sz(hero) = %d, want %d", got, wantHero)
	}

	cta := ops[3]["params"].(map[string]interface{})
	if got := asInt(cta["r"]); got != 32 {
		t.Fatalf("r(button) = %d, want 32", got)
	}

	pill := ops[4]["params"].(map[string]interface{})
	if got := asInt(pill["cornerRadius"]); got != 9999 {
		t.Fatalf("cornerRadius(pill) = %d, want 9999", got)
	}
}

func TestResolveTokenAliases_LeavesUnknownAndNumericValues(t *testing.T) {
	ops := []map[string]interface{}{
		{
			"command": "frame",
			"params":  map[string]interface{}{"w": 1080, "h": 1350},
		},
		{
			"command": "text",
			"params":  map[string]interface{}{"sz": "mega", "fontSize": 44},
		},
	}

	applied, _ := ResolveTokenAliases(ops)
	if applied != 0 {
		t.Fatalf("applied = %d, want 0", applied)
	}

	params := ops[1]["params"].(map[string]interface{})
	if got, _ := params["sz"].(string); got != "mega" {
		t.Fatalf("unknown alias mutated: %#v", params["sz"])
	}
	if !almostEqual(asFloat(params["fontSize"]), 44) {
		t.Fatalf("numeric fontSize mutated: %#v", params["fontSize"])
	}
}

func TestResolveTokenAliases_NoRootFrame_NoChanges(t *testing.T) {
	ops := []map[string]interface{}{
		{
			"command": "text",
			"params":  map[string]interface{}{"sz": "hero"},
		},
	}

	applied, rootWidth := ResolveTokenAliases(ops)
	if applied != 0 || rootWidth != 0 {
		t.Fatalf("expected no resolution, got applied=%d rootWidth=%d", applied, rootWidth)
	}
	if got, _ := ops[0]["params"].(map[string]interface{})["sz"].(string); got != "hero" {
		t.Fatalf("value changed unexpectedly: %q", got)
	}
}

func TestResolveTokenAliases_UsesFirstRootFrame(t *testing.T) {
	ops := []map[string]interface{}{
		{
			"command": "frame",
			"params":  map[string]interface{}{"w": 800, "h": 1000},
		},
		{
			"command": "frame",
			"params":  map[string]interface{}{"w": 1200, "h": 1200},
		},
		{
			"command": "text",
			"params":  map[string]interface{}{"w": "content", "sz": "body"},
		},
	}

	_, rootWidth := ResolveTokenAliases(ops)
	if rootWidth != 800 {
		t.Fatalf("root width = %d, want 800", rootWidth)
	}
	params := ops[2]["params"].(map[string]interface{})
	if got := asInt(params["w"]); got != 688 {
		t.Fatalf("w(content) = %d, want 688 for 800px root", got)
	}
}

func asInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(math.Round(n))
	default:
		return -1
	}
}

func asFloat(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	default:
		return math.NaN()
	}
}

func almostEqual(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	return math.Abs(a-b) < 1e-9
}
