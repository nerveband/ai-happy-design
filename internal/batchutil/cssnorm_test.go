package batchutil

import (
	"testing"
)

func TestFlexDirectionColumn(t *testing.T) {
	params := map[string]interface{}{"flexDirection": "column"}
	result := NormalizeCSSProps(params)
	if result["layoutMode"] != "VERTICAL" {
		t.Errorf("expected layoutMode VERTICAL, got %v", result["layoutMode"])
	}
	if _, ok := result["flexDirection"]; ok {
		t.Error("flexDirection should be removed after normalization")
	}
}

func TestFlexDirectionRow(t *testing.T) {
	params := map[string]interface{}{"flexDirection": "row"}
	result := NormalizeCSSProps(params)
	if result["layoutMode"] != "HORIZONTAL" {
		t.Errorf("expected layoutMode HORIZONTAL, got %v", result["layoutMode"])
	}
}

func TestGapToItemSpacing(t *testing.T) {
	params := map[string]interface{}{"gap": float64(16)}
	result := NormalizeCSSProps(params)
	if result["itemSpacing"] != float64(16) {
		t.Errorf("expected itemSpacing 16, got %v", result["itemSpacing"])
	}
	if _, ok := result["gap"]; ok {
		t.Error("gap should be removed after normalization")
	}
}

func TestBackgroundToColor(t *testing.T) {
	params := map[string]interface{}{"background": "#FF0000"}
	result := NormalizeCSSProps(params)
	if result["color"] != "#FF0000" {
		t.Errorf("expected color #FF0000, got %v", result["color"])
	}
	if _, ok := result["background"]; ok {
		t.Error("background should be removed after normalization")
	}
}

func TestBackgroundColorToColor(t *testing.T) {
	params := map[string]interface{}{"backgroundColor": "#00FF00"}
	result := NormalizeCSSProps(params)
	if result["color"] != "#00FF00" {
		t.Errorf("expected color #00FF00, got %v", result["color"])
	}
}

func TestPaddingNumber(t *testing.T) {
	params := map[string]interface{}{"padding": float64(24)}
	result := NormalizeCSSProps(params)
	for _, key := range []string{"paddingTop", "paddingRight", "paddingBottom", "paddingLeft"} {
		if result[key] != float64(24) {
			t.Errorf("expected %s=24, got %v", key, result[key])
		}
	}
}

func TestPaddingStringNumber(t *testing.T) {
	params := map[string]interface{}{"padding": "24"}
	result := NormalizeCSSProps(params)
	for _, key := range []string{"paddingTop", "paddingRight", "paddingBottom", "paddingLeft"} {
		if result[key] != float64(24) {
			t.Errorf("expected %s=24, got %v", key, result[key])
		}
	}
}

func TestPaddingTwoValues(t *testing.T) {
	params := map[string]interface{}{"padding": "24 32"}
	result := NormalizeCSSProps(params)
	if result["paddingTop"] != float64(24) {
		t.Errorf("expected paddingTop=24, got %v", result["paddingTop"])
	}
	if result["paddingBottom"] != float64(24) {
		t.Errorf("expected paddingBottom=24, got %v", result["paddingBottom"])
	}
	if result["paddingRight"] != float64(32) {
		t.Errorf("expected paddingRight=32, got %v", result["paddingRight"])
	}
	if result["paddingLeft"] != float64(32) {
		t.Errorf("expected paddingLeft=32, got %v", result["paddingLeft"])
	}
}

func TestPaddingFourValues(t *testing.T) {
	params := map[string]interface{}{"padding": "24 32 16 32"}
	result := NormalizeCSSProps(params)
	if result["paddingTop"] != float64(24) {
		t.Errorf("expected paddingTop=24, got %v", result["paddingTop"])
	}
	if result["paddingRight"] != float64(32) {
		t.Errorf("expected paddingRight=32, got %v", result["paddingRight"])
	}
	if result["paddingBottom"] != float64(16) {
		t.Errorf("expected paddingBottom=16, got %v", result["paddingBottom"])
	}
	if result["paddingLeft"] != float64(32) {
		t.Errorf("expected paddingLeft=32, got %v", result["paddingLeft"])
	}
}

func TestBorderShorthand(t *testing.T) {
	params := map[string]interface{}{"border": "1px solid #333"}
	result := NormalizeCSSProps(params)
	if result["strokeWeight"] != float64(1) {
		t.Errorf("expected strokeWeight=1, got %v", result["strokeWeight"])
	}
	if result["stroke"] != "#333" {
		t.Errorf("expected stroke=#333, got %v", result["stroke"])
	}
}

func TestBorderShorthandThick(t *testing.T) {
	params := map[string]interface{}{"border": "2px solid #FF0000"}
	result := NormalizeCSSProps(params)
	if result["strokeWeight"] != float64(2) {
		t.Errorf("expected strokeWeight=2, got %v", result["strokeWeight"])
	}
	if result["stroke"] != "#FF0000" {
		t.Errorf("expected stroke=#FF0000, got %v", result["stroke"])
	}
}

func TestJustifyContentCenter(t *testing.T) {
	params := map[string]interface{}{"justifyContent": "center"}
	result := NormalizeCSSProps(params)
	if result["primaryAxisAlignItems"] != "CENTER" {
		t.Errorf("expected primaryAxisAlignItems=CENTER, got %v", result["primaryAxisAlignItems"])
	}
	if _, ok := result["justifyContent"]; ok {
		t.Error("justifyContent should be removed after normalization")
	}
}

func TestJustifyContentSpaceBetween(t *testing.T) {
	params := map[string]interface{}{"justifyContent": "space-between"}
	result := NormalizeCSSProps(params)
	if result["primaryAxisAlignItems"] != "SPACE_BETWEEN" {
		t.Errorf("expected primaryAxisAlignItems=SPACE_BETWEEN, got %v", result["primaryAxisAlignItems"])
	}
}

func TestAlignItemsFlexEnd(t *testing.T) {
	params := map[string]interface{}{"alignItems": "flex-end"}
	result := NormalizeCSSProps(params)
	if result["counterAxisAlignItems"] != "MAX" {
		t.Errorf("expected counterAxisAlignItems=MAX, got %v", result["counterAxisAlignItems"])
	}
	if _, ok := result["alignItems"]; ok {
		t.Error("alignItems should be removed after normalization")
	}
}

func TestAlignItemsFlexStart(t *testing.T) {
	params := map[string]interface{}{"alignItems": "flex-start"}
	result := NormalizeCSSProps(params)
	if result["counterAxisAlignItems"] != "MIN" {
		t.Errorf("expected counterAxisAlignItems=MIN, got %v", result["counterAxisAlignItems"])
	}
}

func TestWidthPercent(t *testing.T) {
	params := map[string]interface{}{"width": "100%"}
	result := NormalizeCSSProps(params)
	if result["layoutSizingHorizontal"] != "FILL" {
		t.Errorf("expected layoutSizingHorizontal=FILL, got %v", result["layoutSizingHorizontal"])
	}
	if _, ok := result["width"]; ok {
		t.Error("width should be removed when converted to layoutSizingHorizontal")
	}
}

func TestWidthAuto(t *testing.T) {
	params := map[string]interface{}{"width": "auto"}
	result := NormalizeCSSProps(params)
	if result["layoutSizingHorizontal"] != "HUG" {
		t.Errorf("expected layoutSizingHorizontal=HUG, got %v", result["layoutSizingHorizontal"])
	}
}

func TestHeightPercent(t *testing.T) {
	params := map[string]interface{}{"height": "100%"}
	result := NormalizeCSSProps(params)
	if result["layoutSizingVertical"] != "FILL" {
		t.Errorf("expected layoutSizingVertical=FILL, got %v", result["layoutSizingVertical"])
	}
}

func TestHeightAuto(t *testing.T) {
	params := map[string]interface{}{"height": "auto"}
	result := NormalizeCSSProps(params)
	if result["layoutSizingVertical"] != "HUG" {
		t.Errorf("expected layoutSizingVertical=HUG, got %v", result["layoutSizingVertical"])
	}
}

func TestBorderRadiusToCornerRadius(t *testing.T) {
	params := map[string]interface{}{"borderRadius": float64(8)}
	result := NormalizeCSSProps(params)
	if result["cornerRadius"] != float64(8) {
		t.Errorf("expected cornerRadius=8, got %v", result["cornerRadius"])
	}
	if _, ok := result["borderRadius"]; ok {
		t.Error("borderRadius should be removed after normalization")
	}
}

func TestBorderRadiusString(t *testing.T) {
	params := map[string]interface{}{"borderRadius": "12"}
	result := NormalizeCSSProps(params)
	if result["cornerRadius"] != float64(12) {
		t.Errorf("expected cornerRadius=12, got %v", result["cornerRadius"])
	}
}

func TestNoOverwriteExistingColor(t *testing.T) {
	params := map[string]interface{}{
		"color":      "#EXISTING",
		"background": "#SHOULD_NOT_WIN",
	}
	result := NormalizeCSSProps(params)
	if result["color"] != "#EXISTING" {
		t.Errorf("expected color=#EXISTING (not overwritten), got %v", result["color"])
	}
}

func TestNoOverwriteExistingLayoutMode(t *testing.T) {
	params := map[string]interface{}{
		"layoutMode":    "HORIZONTAL",
		"flexDirection": "column",
	}
	result := NormalizeCSSProps(params)
	if result["layoutMode"] != "HORIZONTAL" {
		t.Errorf("expected layoutMode=HORIZONTAL (not overwritten), got %v", result["layoutMode"])
	}
}

func TestNoOverwriteExistingPadding(t *testing.T) {
	params := map[string]interface{}{
		"paddingTop": float64(10),
		"padding":    float64(24),
	}
	result := NormalizeCSSProps(params)
	if result["paddingTop"] != float64(10) {
		t.Errorf("expected paddingTop=10 (not overwritten), got %v", result["paddingTop"])
	}
	// Other sides should still be set
	if result["paddingRight"] != float64(24) {
		t.Errorf("expected paddingRight=24, got %v", result["paddingRight"])
	}
}

func TestBoxShadowRGBA(t *testing.T) {
	params := map[string]interface{}{
		"boxShadow": "0 4px 12px rgba(0,0,0,0.1)",
	}
	result := NormalizeCSSProps(params)
	shadows, ok := result["_shadows"].([]map[string]interface{})
	if !ok || len(shadows) != 1 {
		t.Fatalf("expected 1 shadow, got %v", result["_shadows"])
	}
	s := shadows[0]
	if s["offsetX"] != float64(0) {
		t.Errorf("expected offsetX=0, got %v", s["offsetX"])
	}
	if s["offsetY"] != float64(4) {
		t.Errorf("expected offsetY=4, got %v", s["offsetY"])
	}
	if s["radius"] != float64(12) {
		t.Errorf("expected radius=12, got %v", s["radius"])
	}
	if s["color"] != "rgba(0,0,0,0.1)" {
		t.Errorf("expected color=rgba(0,0,0,0.1), got %v", s["color"])
	}
}

func TestBoxShadowHex(t *testing.T) {
	params := map[string]interface{}{
		"boxShadow": "0 4px 12px #0000001A",
	}
	result := NormalizeCSSProps(params)
	shadows, ok := result["_shadows"].([]map[string]interface{})
	if !ok || len(shadows) != 1 {
		t.Fatalf("expected 1 shadow, got %v", result["_shadows"])
	}
	if shadows[0]["color"] != "#0000001A" {
		t.Errorf("expected color=#0000001A, got %v", shadows[0]["color"])
	}
}

func TestBoxShadowMultiple(t *testing.T) {
	params := map[string]interface{}{
		"boxShadow": "0 2px 4px rgba(0,0,0,0.1), 0 8px 16px rgba(0,0,0,0.2)",
	}
	result := NormalizeCSSProps(params)
	shadows, ok := result["_shadows"].([]map[string]interface{})
	if !ok || len(shadows) != 2 {
		t.Fatalf("expected 2 shadows, got %v", result["_shadows"])
	}
	if shadows[0]["radius"] != float64(4) {
		t.Errorf("expected first shadow radius=4, got %v", shadows[0]["radius"])
	}
	if shadows[1]["radius"] != float64(16) {
		t.Errorf("expected second shadow radius=16, got %v", shadows[1]["radius"])
	}
}

func TestDisplayFlexRemoved(t *testing.T) {
	params := map[string]interface{}{
		"display":       "flex",
		"flexDirection": "row",
	}
	result := NormalizeCSSProps(params)
	if _, ok := result["display"]; ok {
		t.Error("display:flex should be removed")
	}
	if result["layoutMode"] != "HORIZONTAL" {
		t.Errorf("expected layoutMode=HORIZONTAL, got %v", result["layoutMode"])
	}
}

func TestNilParams(t *testing.T) {
	result := NormalizeCSSProps(nil)
	if result != nil {
		t.Errorf("expected nil for nil input, got %v", result)
	}
}

func TestEmptyParams(t *testing.T) {
	params := map[string]interface{}{}
	result := NormalizeCSSProps(params)
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestGapStringValue(t *testing.T) {
	params := map[string]interface{}{"gap": "16"}
	result := NormalizeCSSProps(params)
	if result["itemSpacing"] != float64(16) {
		t.Errorf("expected itemSpacing=16, got %v", result["itemSpacing"])
	}
}
