package tools

import (
	"testing"
)

func TestComputeDesignTokens_SquareCanvas(t *testing.T) {
	tokens := ComputeDesignTokens(1080, 1080)
	assertTokenBasics(t, tokens, "square", 1080, 1080)
}

func TestComputeDesignTokens_PortraitCarousel(t *testing.T) {
	// 1350/1080 = 1.25 — below the 1.3 portrait threshold, so classified as "square"
	tokens := ComputeDesignTokens(1080, 1350)
	assertTokenBasics(t, tokens, "square", 1080, 1350)
}

func TestComputeDesignTokens_Story(t *testing.T) {
	tokens := ComputeDesignTokens(1080, 1920)
	assertTokenBasics(t, tokens, "portrait", 1080, 1920)
}

func TestComputeDesignTokens_EmailBanner(t *testing.T) {
	tokens := ComputeDesignTokens(600, 200)
	assertTokenBasics(t, tokens, "landscape", 600, 200)
}

func TestComputeDesignTokens_UltraTall(t *testing.T) {
	tokens := ComputeDesignTokens(500, 1200)
	assertTokenBasics(t, tokens, "ultra-tall", 500, 1200)
}

func TestComputeDesignTokens_LargerCanvasLargerFonts(t *testing.T) {
	small := ComputeDesignTokens(600, 600)
	large := ComputeDesignTokens(1080, 1080)

	smallText := small["text"].(map[string]interface{})
	largeText := large["text"].(map[string]interface{})

	smallHero := smallText["hero"].(int)
	largeHero := largeText["hero"].(int)
	if largeHero < smallHero {
		t.Errorf("larger canvas should produce larger (or equal) hero font: small=%d, large=%d", smallHero, largeHero)
	}

	smallBody := smallText["body"].(int)
	largeBody := largeText["body"].(int)
	if largeBody < smallBody {
		t.Errorf("larger canvas should produce larger (or equal) body font: small=%d, large=%d", smallBody, largeBody)
	}

	smallHeading := smallText["heading"].(int)
	largeHeading := largeText["heading"].(int)
	if largeHeading < smallHeading {
		t.Errorf("larger canvas should produce larger (or equal) heading font: small=%d, large=%d", smallHeading, largeHeading)
	}
}

func TestComputeDesignTokens_CardWidthsLessThanCanvas(t *testing.T) {
	sizes := []struct {
		w, h float64
	}{
		{1080, 1080},
		{1080, 1350},
		{1080, 1920},
		{600, 200},
		{500, 1200},
	}

	for _, sz := range sizes {
		tokens := ComputeDesignTokens(sz.w, sz.h)
		cards := tokens["cards"].(map[string]interface{})
		widths := cards["widths"].(map[string]interface{})

		for name, val := range widths {
			w := val.(int)
			if w >= int(sz.w) {
				t.Errorf("canvas %.0fx%.0f: card width %q (%d) should be less than canvas width (%.0f)",
					sz.w, sz.h, name, w, sz.w)
			}
			if w <= 0 {
				t.Errorf("canvas %.0fx%.0f: card width %q (%d) should be positive",
					sz.w, sz.h, name, w)
			}
		}
	}
}

func TestComputeDesignTokens_AllValuesPositive(t *testing.T) {
	sizes := []struct {
		w, h float64
	}{
		{1080, 1080},
		{1080, 1350},
		{1080, 1920},
		{600, 200},
		{500, 1200},
	}

	for _, sz := range sizes {
		tokens := ComputeDesignTokens(sz.w, sz.h)

		// Check text sizes
		text := tokens["text"].(map[string]interface{})
		for _, key := range []string{"hero", "heading", "subheading", "body", "caption", "numbers", "cta"} {
			val := text[key].(int)
			if val <= 0 {
				t.Errorf("canvas %.0fx%.0f: text.%s = %d, expected > 0", sz.w, sz.h, key, val)
			}
		}

		// Check spacing values
		spacing := tokens["spacing"].(map[string]interface{})
		for _, key := range []string{"sidePadding", "contentWidth", "framePadding", "cardPadding", "itemSpacing", "cardGap", "sectionGap", "topPadding"} {
			val := spacing[key].(int)
			if val <= 0 {
				t.Errorf("canvas %.0fx%.0f: spacing.%s = %d, expected > 0", sz.w, sz.h, key, val)
			}
		}
	}
}

func TestComputeDesignTokens_LayoutTypeClassification(t *testing.T) {
	tests := []struct {
		w, h       float64
		layoutType string
	}{
		{1080, 1080, "square"},      // ratio = 1.0
		{1080, 1350, "square"},      // ratio = 1.25 (below 1.3 portrait threshold)
		{1080, 1920, "portrait"},    // ratio = 1.78
		{600, 200, "landscape"},     // ratio = 0.33
		{500, 1200, "ultra-tall"},   // ratio = 2.4
		{1000, 900, "square"},       // ratio = 0.9
		{1000, 1300, "portrait"},    // ratio = 1.3 — exactly at portrait threshold (>= 1.3)
	}

	for _, tt := range tests {
		tokens := ComputeDesignTokens(tt.w, tt.h)
		summary := tokens["_summary"].(map[string]interface{})
		got := summary["layoutType"].(string)
		if got != tt.layoutType {
			t.Errorf("canvas %.0fx%.0f: expected layout %q, got %q", tt.w, tt.h, tt.layoutType, got)
		}
	}
}

func TestComputeDesignTokens_Round8(t *testing.T) {
	// All text sizes and spacing should be multiples of 8 (due to round8/max8)
	tokens := ComputeDesignTokens(1080, 1350)

	text := tokens["text"].(map[string]interface{})
	for _, key := range []string{"hero", "heading", "subheading", "body", "caption", "numbers", "cta"} {
		val := text[key].(int)
		if val%8 != 0 {
			t.Errorf("text.%s = %d, expected multiple of 8", key, val)
		}
	}

	spacing := tokens["spacing"].(map[string]interface{})
	for _, key := range []string{"sidePadding", "framePadding", "cardPadding", "itemSpacing", "cardGap", "sectionGap", "topPadding"} {
		val := spacing[key].(int)
		if val%8 != 0 {
			t.Errorf("spacing.%s = %d, expected multiple of 8", key, val)
		}
	}
}

// assertTokenBasics verifies the required structure and basic invariants.
func assertTokenBasics(t *testing.T, tokens map[string]interface{}, expectedLayout string, w, h float64) {
	t.Helper()

	// Check _summary exists with correct canvas dims
	summary, ok := tokens["_summary"].(map[string]interface{})
	if !ok {
		t.Fatal("missing _summary")
	}
	canvas := summary["canvas"].(map[string]interface{})
	if canvas["width"] != int(w) || canvas["height"] != int(h) {
		t.Errorf("canvas dimensions mismatch: got %v", canvas)
	}
	if summary["layoutType"] != expectedLayout {
		t.Errorf("expected layout %q, got %q", expectedLayout, summary["layoutType"])
	}

	// Check required top-level keys
	for _, key := range []string{"text", "spacing", "cards", "textRules"} {
		if tokens[key] == nil {
			t.Errorf("missing top-level key: %s", key)
		}
	}

	// Check text sizes are positive
	text := tokens["text"].(map[string]interface{})
	for _, key := range []string{"hero", "heading", "subheading", "body", "caption"} {
		val, ok := text[key].(int)
		if !ok || val <= 0 {
			t.Errorf("text.%s should be positive int, got %v", key, text[key])
		}
	}

	// Check spacing values are positive
	spacing := tokens["spacing"].(map[string]interface{})
	for _, key := range []string{"sidePadding", "contentWidth", "cardPadding", "itemSpacing"} {
		val, ok := spacing[key].(int)
		if !ok || val <= 0 {
			t.Errorf("spacing.%s should be positive int, got %v", key, spacing[key])
		}
	}

	// contentWidth should be less than canvas width
	contentWidth := spacing["contentWidth"].(int)
	if contentWidth >= int(w) {
		t.Errorf("contentWidth (%d) should be less than canvas width (%.0f)", contentWidth, w)
	}
}
