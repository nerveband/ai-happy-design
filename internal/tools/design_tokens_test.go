package tools

import (
	"strings"
	"testing"
)

func TestComputeDesignTokens_SquareCanvas(t *testing.T) {
	tokens := ComputeDesignTokens(1080, 1080, 0)
	assertTokenBasics(t, tokens, "square", 1080, 1080)
}

func TestComputeDesignTokens_PortraitCarousel(t *testing.T) {
	// 1350/1080 = 1.25 — now classified as "portrait" (threshold lowered to 1.15)
	tokens := ComputeDesignTokens(1080, 1350, 0)
	assertTokenBasics(t, tokens, "portrait", 1080, 1350)
}

func TestComputeDesignTokens_Story(t *testing.T) {
	tokens := ComputeDesignTokens(1080, 1920, 0)
	assertTokenBasics(t, tokens, "portrait", 1080, 1920)
}

func TestComputeDesignTokens_EmailBanner(t *testing.T) {
	tokens := ComputeDesignTokens(600, 200, 0)
	assertTokenBasics(t, tokens, "landscape", 600, 200)
}

func TestComputeDesignTokens_UltraTall(t *testing.T) {
	tokens := ComputeDesignTokens(500, 1200, 0)
	assertTokenBasics(t, tokens, "ultra-tall", 500, 1200)
}

func TestComputeDesignTokens_LargerCanvasLargerFonts(t *testing.T) {
	small := ComputeDesignTokens(600, 600, 0)
	large := ComputeDesignTokens(1080, 1080, 0)

	smallText := small["text"].(map[string]interface{})
	largeText := large["text"].(map[string]interface{})

	for _, key := range []string{"display", "hero", "title", "heading", "subheading", "body", "caption"} {
		s := smallText[key].(int)
		l := largeText[key].(int)
		if l < s {
			t.Errorf("larger canvas should produce larger (or equal) %s font: small=%d, large=%d", key, s, l)
		}
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
		tokens := ComputeDesignTokens(sz.w, sz.h, 0)
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
		tokens := ComputeDesignTokens(sz.w, sz.h, 0)

		// Check text sizes
		text := tokens["text"].(map[string]interface{})
		for _, key := range []string{"display", "hero", "title", "heading", "subheading", "body", "caption", "numbers", "cta"} {
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

		// Check button values
		button := tokens["button"].(map[string]interface{})
		for _, key := range []string{"fontSize", "height", "paddingH", "cornerPill", "cornerRounded"} {
			val := button[key].(int)
			if val <= 0 {
				t.Errorf("canvas %.0fx%.0f: button.%s = %d, expected > 0", sz.w, sz.h, key, val)
			}
		}
	}
}

func TestComputeDesignTokens_LayoutTypeClassification(t *testing.T) {
	tests := []struct {
		w, h       float64
		layoutType string
	}{
		{1080, 1080, "square"},    // ratio = 1.0
		{1080, 1350, "portrait"},  // ratio = 1.25 (now portrait, was square)
		{1080, 1920, "portrait"},  // ratio = 1.78
		{600, 200, "landscape"},   // ratio = 0.33
		{500, 1200, "ultra-tall"}, // ratio = 2.4
		{1000, 900, "square"},     // ratio = 0.9
		{1000, 1150, "portrait"},  // ratio = 1.15 — exactly at portrait threshold (>= 1.15)
		{1000, 850, "square"},     // ratio = 0.85 — exactly at square threshold (>= 0.85)
	}

	for _, tt := range tests {
		tokens := ComputeDesignTokens(tt.w, tt.h, 0)
		summary := tokens["_summary"].(map[string]interface{})
		got := summary["layoutType"].(string)
		if got != tt.layoutType {
			t.Errorf("canvas %.0fx%.0f: expected layout %q, got %q", tt.w, tt.h, tt.layoutType, got)
		}
	}
}

func TestComputeDesignTokens_TextRound4(t *testing.T) {
	// All text sizes should be multiples of 4 (round4 for finer type scale)
	tokens := ComputeDesignTokens(1080, 1350, 0)

	text := tokens["text"].(map[string]interface{})
	for _, key := range []string{"display", "hero", "title", "heading", "subheading", "body", "caption", "numbers", "cta"} {
		val := text[key].(int)
		if val%4 != 0 {
			t.Errorf("text.%s = %d, expected multiple of 4", key, val)
		}
	}
}

func TestComputeDesignTokens_SpacingRound8(t *testing.T) {
	// Spacing values should be multiples of 8
	tokens := ComputeDesignTokens(1080, 1350, 0)

	spacing := tokens["spacing"].(map[string]interface{})
	for _, key := range []string{"sidePadding", "framePadding", "cardPadding", "itemSpacing", "cardGap", "sectionGap", "topPadding"} {
		val := spacing[key].(int)
		if val%8 != 0 {
			t.Errorf("spacing.%s = %d, expected multiple of 8", key, val)
		}
	}
}

func TestComputeDesignTokens_ModularScale(t *testing.T) {
	// The type scale should be monotonically increasing
	tokens := ComputeDesignTokens(1080, 1350, 0)
	text := tokens["text"].(map[string]interface{})

	scale := []struct {
		name string
		val  int
	}{
		{"caption", text["caption"].(int)},
		{"body", text["body"].(int)},
		{"subheading", text["subheading"].(int)},
		{"heading", text["heading"].(int)},
		{"title", text["title"].(int)},
		{"hero", text["hero"].(int)},
		{"display", text["display"].(int)},
	}

	for i := 1; i < len(scale); i++ {
		if scale[i].val <= scale[i-1].val {
			t.Errorf("type scale not increasing: %s(%d) should be > %s(%d)",
				scale[i].name, scale[i].val, scale[i-1].name, scale[i-1].val)
		}
	}
}

func TestComputeDesignTokens_NumbersEqualsDisplay(t *testing.T) {
	tokens := ComputeDesignTokens(1080, 1350, 0)
	text := tokens["text"].(map[string]interface{})
	if text["numbers"] != text["display"] {
		t.Errorf("numbers (%d) should equal display (%d)", text["numbers"], text["display"])
	}
}

func TestComputeDesignTokens_ButtonSizing(t *testing.T) {
	tokens := ComputeDesignTokens(1080, 1350, 0)
	button := tokens["button"].(map[string]interface{})
	text := tokens["text"].(map[string]interface{})

	btnFont := button["fontSize"].(int)
	btnHeight := button["height"].(int)
	body := text["body"].(int)

	// Button font should equal body
	if btnFont != body {
		t.Errorf("button fontSize (%d) should equal body (%d)", btnFont, body)
	}

	// Button height should be at least 2× font size
	if btnHeight < btnFont*2 {
		t.Errorf("button height (%d) should be >= 2× fontSize (%d)", btnHeight, btnFont)
	}

	// Pill corner should be half the height
	pill := button["cornerPill"].(int)
	if pill != btnHeight/2 {
		t.Errorf("cornerPill (%d) should be height/2 (%d)", pill, btnHeight/2)
	}
}

func TestComputeDesignTokens_1080x1350_SpecificValues(t *testing.T) {
	// Verify the exact computed values for the most common social media format
	tokens := ComputeDesignTokens(1080, 1350, 0)
	text := tokens["text"].(map[string]interface{})

	expected := map[string]int{
		"display":    200,
		"hero":       152,
		"title":      112,
		"heading":    84,
		"subheading": 64,
		"body":       48,
		"caption":    36,
	}

	for name, want := range expected {
		got := text[name].(int)
		if got != want {
			t.Errorf("text.%s = %d, want %d", name, got, want)
		}
	}
}

func TestComputeDesignTokens_DefaultDPI(t *testing.T) {
	// dpi=0 and dpi=72 should produce identical results to the original behavior
	tokensDefault := ComputeDesignTokens(1080, 1080, 0)
	tokens72 := ComputeDesignTokens(1080, 1080, 72)

	textDefault := tokensDefault["text"].(map[string]interface{})
	text72 := tokens72["text"].(map[string]interface{})

	for _, key := range []string{"display", "hero", "title", "heading", "subheading", "body", "caption"} {
		d := textDefault[key].(int)
		s := text72[key].(int)
		if d != s {
			t.Errorf("dpi=0 vs dpi=72: text.%s differs: %d vs %d", key, d, s)
		}
	}

	// Neither should have _print metadata
	if tokensDefault["_print"] != nil {
		t.Error("dpi=0 should not have _print metadata")
	}
	if tokens72["_print"] != nil {
		t.Error("dpi=72 should not have _print metadata")
	}
}

func TestComputeDesignTokens_PrintDPI300(t *testing.T) {
	// 2550x3300 @ 300dpi = 8.5x11" letter
	tokens := ComputeDesignTokens(2550, 3300, 300)
	text := tokens["text"].(map[string]interface{})

	body := text["body"].(int)
	caption := text["caption"].(int)
	heading := text["heading"].(int)

	// Body should be around 50px (12pt at 300dpi)
	if body < 40 || body > 64 {
		t.Errorf("print body = %d, expected ~50px (12pt at 300dpi)", body)
	}

	// Caption should be smaller than body
	if caption >= body {
		t.Errorf("print caption (%d) should be < body (%d)", caption, body)
	}

	// Heading should be larger than body
	if heading <= body {
		t.Errorf("print heading (%d) should be > body (%d)", heading, body)
	}

	// All text sizes should still be on 4px grid
	for _, key := range []string{"display", "hero", "title", "heading", "subheading", "body", "caption"} {
		val := text[key].(int)
		if val%4 != 0 {
			t.Errorf("print text.%s = %d, expected multiple of 4", key, val)
		}
	}

	// Should have _print metadata
	printMeta, ok := tokens["_print"].(map[string]interface{})
	if !ok {
		t.Fatal("missing _print metadata for dpi=300")
	}
	if printMeta["dpi"] != 300.0 {
		t.Errorf("_print.dpi = %v, expected 300", printMeta["dpi"])
	}
	physW := printMeta["physicalWidth"].(float64)
	if physW < 8.4 || physW > 8.6 {
		t.Errorf("_print.physicalWidth = %v, expected ~8.5", physW)
	}
}

func TestComputeDesignTokens_PrintDPI150(t *testing.T) {
	// 1275x1650 @ 150dpi = 8.5x11" letter (same physical size)
	tokens := ComputeDesignTokens(1275, 1650, 150)
	text := tokens["text"].(map[string]interface{})

	body := text["body"].(int)
	// Body should be roughly half of 300dpi body (since canvas is half the pixels)
	// But same physical size → same pt → body ≈ 25px (12pt at 150dpi)
	if body < 20 || body > 36 {
		t.Errorf("150dpi body = %d, expected ~25px (12pt at 150dpi)", body)
	}

	printMeta := tokens["_print"].(map[string]interface{})
	physW := printMeta["physicalWidth"].(float64)
	if physW < 8.4 || physW > 8.6 {
		t.Errorf("_print.physicalWidth = %v, expected ~8.5", physW)
	}
}

func TestComputeDesignTokens_PrintVsScreen_SameCanvas(t *testing.T) {
	// Same canvas size, but dpi=300 gives SMALLER pixel values than screen
	screen := ComputeDesignTokens(2550, 3300, 0)
	print := ComputeDesignTokens(2550, 3300, 300)

	screenBody := screen["text"].(map[string]interface{})["body"].(int)
	printBody := print["text"].(map[string]interface{})["body"].(int)

	if printBody >= screenBody {
		t.Errorf("print body (%d) should be smaller than screen body (%d) for 2550px canvas",
			printBody, screenBody)
	}
}

func TestComputeDesignTokens_PrintMinimums(t *testing.T) {
	// Business card: 3.5x2" = 1050x600 @ 300dpi
	tokens := ComputeDesignTokens(1050, 600, 300)
	text := tokens["text"].(map[string]interface{})

	body := text["body"].(int)
	// Physical width = 3.5", bodyPt = clamp(3.5*1.4, 10, 14) = clamp(4.9, 10, 14) = 10pt minimum
	// body = 10 * 300/72 ≈ 42px
	if body < 36 || body > 52 {
		t.Errorf("business card body = %d, expected ~42px (10pt minimum at 300dpi)", body)
	}
}

func TestComputeDesignTokens_PrintSpacingRound8(t *testing.T) {
	tokens := ComputeDesignTokens(2550, 3300, 300)
	spacing := tokens["spacing"].(map[string]interface{})

	for _, key := range []string{"sidePadding", "framePadding", "cardPadding", "itemSpacing", "cardGap", "sectionGap", "topPadding"} {
		val := spacing[key].(int)
		if val%8 != 0 {
			t.Errorf("print spacing.%s = %d, expected multiple of 8", key, val)
		}
	}
}

func TestComputeDesignTokens_TipsIncludeFirstPassQualityGuidance(t *testing.T) {
	tokens := ComputeDesignTokens(1080, 1350, 0)
	tips, ok := tokens["tips"].([]string)
	if !ok {
		t.Fatalf("expected tips []string, got %T", tokens["tips"])
	}

	joined := ""
	for _, tip := range tips {
		joined += tip + " "
	}
	lower := strings.ToLower(joined)
	if !strings.Contains(lower, "--strict-quality") {
		t.Fatalf("expected tips to mention --strict-quality, got %v", tips)
	}
	if !strings.Contains(lower, "layoutpositioning:absolute") {
		t.Fatalf("expected tips to mention ABSOLUTE layout rule, got %v", tips)
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
	for _, key := range []string{"text", "spacing", "cards", "button", "textRules"} {
		if tokens[key] == nil {
			t.Errorf("missing top-level key: %s", key)
		}
	}

	// Check text sizes are positive
	text := tokens["text"].(map[string]interface{})
	for _, key := range []string{"display", "hero", "title", "heading", "subheading", "body", "caption"} {
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
