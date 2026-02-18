package batchutil

import (
	"testing"
)

func TestIsComposite(t *testing.T) {
	tests := []struct {
		cmd  string
		want bool
	}{
		{"slide", true},
		{"Slide", true},
		{"SLIDE", true},
		{"banner", true},
		{"Banner", true},
		{"node.create_frame", false},
		{"text.create", false},
		{"", false},
	}
	for _, tt := range tests {
		got := IsComposite(tt.cmd)
		if got != tt.want {
			t.Errorf("IsComposite(%q) = %v, want %v", tt.cmd, got, tt.want)
		}
	}
}

func TestExpandSlideCommand(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s1",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#0C1E2C",
			"elements": []interface{}{
				map[string]interface{}{
					"type": "eyebrow",
					"text": "CHAPTER 1",
				},
				map[string]interface{}{
					"type": "headline",
					"text": "Big Title",
					"tier": "hero",
				},
				map[string]interface{}{
					"type": "body",
					"text": "Some paragraph text here.",
				},
				map[string]interface{}{
					"type": "bar",
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	// Should have: frame + 4 elements = 5 ops minimum
	if len(ops) < 5 {
		t.Fatalf("expected at least 5 ops, got %d", len(ops))
	}

	// First op should be the root frame
	frame := ops[0]
	if frame["command"] != "node.create_frame" {
		t.Errorf("first op command = %v, want node.create_frame", frame["command"])
	}
	if frame["name"] != "s1_frame" {
		t.Errorf("frame name = %v, want s1_frame", frame["name"])
	}
	frameParams, _ := frame["params"].(map[string]interface{})
	if frameParams["width"] != 1080.0 {
		t.Errorf("frame width = %v, want 1080", frameParams["width"])
	}
	if frameParams["height"] != 1350.0 {
		t.Errorf("frame height = %v, want 1350", frameParams["height"])
	}
	if frameParams["color"] != "#0C1E2C" {
		t.Errorf("frame color = %v, want #0C1E2C", frameParams["color"])
	}

	// Second op should be the eyebrow
	eyebrow := ops[1]
	if eyebrow["command"] != "text.create" {
		t.Errorf("eyebrow command = %v, want text.create", eyebrow["command"])
	}
	if eyebrow["name"] != "s1_e0" {
		t.Errorf("eyebrow name = %v, want s1_e0", eyebrow["name"])
	}
	ebParams, _ := eyebrow["params"].(map[string]interface{})
	if ebParams["textCase"] != "UPPER" {
		t.Errorf("eyebrow textCase = %v, want UPPER", ebParams["textCase"])
	}
	if ebParams["parentId"] != "${{steps.s1_frame.result.id}}" {
		t.Errorf("eyebrow parentId = %v, want ref to s1_frame", ebParams["parentId"])
	}

	// Third op should be headline
	headline := ops[2]
	if headline["name"] != "s1_e1" {
		t.Errorf("headline name = %v, want s1_e1", headline["name"])
	}
	hlParams, _ := headline["params"].(map[string]interface{})
	// Hero tier at 1080px width: base=47.52 → hero = snap4(47.52 * 1.333^4) = snap4(119.2) = 120
	heroSize := hlParams["fontSize"].(float64)
	if heroSize < 100 || heroSize > 160 {
		t.Errorf("hero fontSize = %v, expected ~120 range for 1080px canvas", heroSize)
	}

	// Fourth op should be body
	body := ops[3]
	bodyParams, _ := body["params"].(map[string]interface{})
	if bodyParams["fontFamily"] != "DM Sans" {
		t.Errorf("body fontFamily = %v, want DM Sans", bodyParams["fontFamily"])
	}
	if bodyParams["lineHeight"] != 150.0 {
		t.Errorf("body lineHeight = %v, want 150", bodyParams["lineHeight"])
	}

	// Fifth op should be bar
	bar := ops[4]
	if bar["command"] != "shape.create_rectangle" {
		t.Errorf("bar command = %v, want shape.create_rectangle", bar["command"])
	}
}

func TestExpandSlideWithGradient(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s2",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#000000",
			"gradient": map[string]interface{}{
				"type": "LINEAR",
				"stops": []interface{}{
					map[string]interface{}{"color": "#FF0000", "position": 0.0},
					map[string]interface{}{"color": "#0000FF", "position": 1.0},
				},
			},
			"elements": []interface{}{},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	// frame + gradient = 2 ops
	if len(ops) < 2 {
		t.Fatalf("expected at least 2 ops, got %d", len(ops))
	}

	if ops[1]["command"] != "paint.set_gradient" {
		t.Errorf("gradient op command = %v, want paint.set_gradient", ops[1]["command"])
	}
	if ops[1]["name"] != "s2_grad" {
		t.Errorf("gradient op name = %v, want s2_grad", ops[1]["name"])
	}
	gradParams, _ := ops[1]["params"].(map[string]interface{})
	if gradParams["nodeId"] != "${{steps.s2_frame.result.id}}" {
		t.Errorf("gradient nodeId = %v, want ref to s2_frame", gradParams["nodeId"])
	}
}

func TestExpandBannerCommand(t *testing.T) {
	op := map[string]interface{}{
		"name":    "b1",
		"command": "banner",
		"params": map[string]interface{}{
			"canvas":       "1200x628",
			"color":        "#FFFFFF",
			"dividerX":     400.0,
			"dividerColor": "#CCCCCC",
			"elements": []interface{}{
				map[string]interface{}{
					"type": "headline",
					"text": "Banner Title",
				},
				map[string]interface{}{
					"type": "subtitle",
					"text": "Supporting text",
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	// frame + divider + headline + subtitle = 4 ops
	if len(ops) < 4 {
		t.Fatalf("expected at least 4 ops, got %d", len(ops))
	}

	// First op: root frame
	if ops[0]["command"] != "node.create_frame" {
		t.Errorf("frame command = %v", ops[0]["command"])
	}
	frameParams, _ := ops[0]["params"].(map[string]interface{})
	if frameParams["width"] != 1200.0 {
		t.Errorf("frame width = %v, want 1200", frameParams["width"])
	}

	// Second op: divider
	if ops[1]["command"] != "shape.create_rectangle" {
		t.Errorf("divider command = %v, want shape.create_rectangle", ops[1]["command"])
	}
	divParams, _ := ops[1]["params"].(map[string]interface{})
	if divParams["x"] != 400.0 {
		t.Errorf("divider x = %v, want 400", divParams["x"])
	}

	// Headline should reference frame
	hlParams, _ := ops[2]["params"].(map[string]interface{})
	if hlParams["parentId"] != "${{steps.b1_frame.result.id}}" {
		t.Errorf("headline parentId = %v", hlParams["parentId"])
	}
	// Content X should be after divider
	hlX := hlParams["x"].(float64)
	if hlX <= 400 {
		t.Errorf("headline x = %v, should be after divider at 400", hlX)
	}
}

func TestExpandStatsElement(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s1",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#FFFFFF",
			"elements": []interface{}{
				map[string]interface{}{
					"type": "stats",
					"items": []interface{}{
						map[string]interface{}{"value": "500+", "label": "Clients"},
						map[string]interface{}{"value": "98%", "label": "Satisfaction"},
						map[string]interface{}{"value": "$2M", "label": "Revenue"},
					},
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	// frame + 3 value texts + 3 label texts = 7 ops
	if len(ops) < 7 {
		t.Fatalf("expected at least 7 ops, got %d", len(ops))
	}

	// Check that values and labels alternate
	// ops[1] should be first value (s1_e0_v0)
	if ops[1]["name"] != "s1_e0_v0" {
		t.Errorf("first stats value name = %v, want s1_e0_v0", ops[1]["name"])
	}
	v0Params, _ := ops[1]["params"].(map[string]interface{})
	if v0Params["content"] != "500+" {
		t.Errorf("first value content = %v, want 500+", v0Params["content"])
	}

	// ops[2] should be first label (s1_e0_l0)
	if ops[2]["name"] != "s1_e0_l0" {
		t.Errorf("first stats label name = %v, want s1_e0_l0", ops[2]["name"])
	}
	l0Params, _ := ops[2]["params"].(map[string]interface{})
	if l0Params["content"] != "Clients" {
		t.Errorf("first label content = %v, want Clients", l0Params["content"])
	}

	// Check columns are spread horizontally
	v1Params, _ := ops[3]["params"].(map[string]interface{})
	if v1Params["x"].(float64) <= v0Params["x"].(float64) {
		t.Errorf("second column x should be > first column x")
	}
}

func TestExpandProgressElement(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s1",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#FFFFFF",
			"elements": []interface{}{
				map[string]interface{}{
					"type":       "progress",
					"current":    75.0,
					"goal":       100.0,
					"trackColor": "#EEE",
					"fillColor":  "#00AA00",
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	// frame + raised_label + goal_label + track + fill = 5 ops
	if len(ops) < 5 {
		t.Fatalf("expected at least 5 ops, got %d", len(ops))
	}

	// Find track and fill rects
	var trackOp, fillOp map[string]interface{}
	for _, o := range ops {
		if o["name"] == "s1_e0_track" {
			trackOp = o
		}
		if o["name"] == "s1_e0_fill" {
			fillOp = o
		}
	}

	if trackOp == nil {
		t.Fatal("track op not found")
	}
	if fillOp == nil {
		t.Fatal("fill op not found")
	}

	trackParams, _ := trackOp["params"].(map[string]interface{})
	fillParams, _ := fillOp["params"].(map[string]interface{})

	// Fill should be 75% of track width
	trackW := trackParams["width"].(float64)
	fillW := fillParams["width"].(float64)
	ratio := fillW / trackW
	if ratio < 0.70 || ratio > 0.80 {
		t.Errorf("fill/track ratio = %v, expected ~0.75", ratio)
	}

	if fillParams["color"] != "#00AA00" {
		t.Errorf("fill color = %v, want #00AA00", fillParams["color"])
	}
}

func TestNonCompositePassthrough(t *testing.T) {
	op := map[string]interface{}{
		"name":    "myframe",
		"command": "node.create_frame",
		"params": map[string]interface{}{
			"x":     0.0,
			"y":     0.0,
			"width": 500.0,
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	if len(ops) != 1 {
		t.Fatalf("expected 1 op for non-composite, got %d", len(ops))
	}

	if ops[0]["command"] != "node.create_frame" {
		t.Errorf("passthrough command = %v", ops[0]["command"])
	}
}

func TestExpandAllComposites(t *testing.T) {
	ops := []map[string]interface{}{
		{
			"name":    "pre",
			"command": "node.create_frame",
			"params":  map[string]interface{}{"x": 0.0},
		},
		{
			"name":    "s1",
			"command": "slide",
			"params": map[string]interface{}{
				"canvas": "1080x1350",
				"color":  "#000",
				"elements": []interface{}{
					map[string]interface{}{
						"type": "headline",
						"text": "Hello",
					},
				},
			},
		},
		{
			"name":    "post",
			"command": "text.create",
			"params":  map[string]interface{}{"content": "after"},
		},
	}

	result, err := ExpandAllComposites(ops)
	if err != nil {
		t.Fatalf("ExpandAllComposites failed: %v", err)
	}

	// pre + (frame + headline) + post = 4 ops
	if len(result) < 4 {
		t.Fatalf("expected at least 4 ops, got %d", len(result))
	}

	// First should be the passthrough
	if result[0]["command"] != "node.create_frame" {
		t.Errorf("first op should be passthrough, got %v", result[0]["command"])
	}

	// Last should be the passthrough
	last := result[len(result)-1]
	if last["command"] != "text.create" {
		t.Errorf("last op should be passthrough, got %v", last["command"])
	}

	// Middle should contain expanded slide ops
	foundSlideFrame := false
	for _, op := range result {
		if op["name"] == "s1_frame" && op["command"] == "node.create_frame" {
			foundSlideFrame = true
		}
	}
	if !foundSlideFrame {
		t.Error("expanded slide frame not found in result")
	}
}

func TestTokenSizes(t *testing.T) {
	sizes := tokenSizes(1080)
	// Base = 1080 * 0.044 = 47.52
	// body = snap4(47.52, 14) = 48
	if sizes["body"] != 48 {
		t.Errorf("body at 1080px = %v, want 48", sizes["body"])
	}
	// caption = snap4(47.52/1.333, 12) = snap4(35.6) = 36
	if sizes["caption"] != 36 {
		t.Errorf("caption at 1080px = %v, want 36", sizes["caption"])
	}
	// subheading = snap4(47.52*1.333, 16) = snap4(63.3) = 64
	if sizes["subheading"] != 64 {
		t.Errorf("subheading at 1080px = %v, want 64", sizes["subheading"])
	}
}

func TestParseCanvas(t *testing.T) {
	w, h, err := parseCanvas("1080x1350")
	if err != nil {
		t.Fatal(err)
	}
	if w != 1080 || h != 1350 {
		t.Errorf("parseCanvas(1080x1350) = %v, %v", w, h)
	}

	_, _, err = parseCanvas("invalid")
	if err == nil {
		t.Error("expected error for invalid canvas")
	}
}

func TestSnap4(t *testing.T) {
	if snap4(47.5, 12) != 48 {
		t.Errorf("snap4(47.5, 12) = %v", snap4(47.5, 12))
	}
	if snap4(5, 12) != 12 {
		t.Errorf("snap4(5, 12) = %v, want 12 (minimum)", snap4(5, 12))
	}
}

func TestExpandCTAElement(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s1",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#FFF",
			"elements": []interface{}{
				map[string]interface{}{
					"type":      "cta",
					"text":      "Get Started",
					"bgColor":   "#FF0000",
					"textColor": "#FFFFFF",
					"style":     "pill",
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	// frame + btn_frame + btn_text = 3 ops
	if len(ops) < 3 {
		t.Fatalf("expected at least 3 ops, got %d", len(ops))
	}

	// Find button frame
	var btnFrame map[string]interface{}
	for _, o := range ops {
		if o["name"] == "s1_e0_btn" {
			btnFrame = o
		}
	}
	if btnFrame == nil {
		t.Fatal("CTA button frame not found")
	}

	btnParams, _ := btnFrame["params"].(map[string]interface{})
	if btnParams["layoutMode"] != "HORIZONTAL" {
		t.Errorf("CTA layoutMode = %v, want HORIZONTAL", btnParams["layoutMode"])
	}
	if btnParams["color"] != "#FF0000" {
		t.Errorf("CTA bg color = %v, want #FF0000", btnParams["color"])
	}
}

func TestExpandCounterElement(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s1",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#FFF",
			"elements": []interface{}{
				map[string]interface{}{
					"type":    "counter",
					"current": "3",
					"total":   "10",
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	// frame + counter = 2 ops
	if len(ops) < 2 {
		t.Fatalf("expected at least 2 ops, got %d", len(ops))
	}

	counterParams, _ := ops[1]["params"].(map[string]interface{})
	if counterParams["content"] != "3 / 10" {
		t.Errorf("counter content = %v, want '3 / 10'", counterParams["content"])
	}
	if counterParams["textAlignHorizontal"] != "RIGHT" {
		t.Errorf("counter align = %v, want RIGHT", counterParams["textAlignHorizontal"])
	}
}

func TestExpandURLElement(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s1",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#FFF",
			"elements": []interface{}{
				map[string]interface{}{
					"type": "url",
					"text": "www.example.com",
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	if len(ops) < 2 {
		t.Fatalf("expected at least 2 ops, got %d", len(ops))
	}

	urlParams, _ := ops[1]["params"].(map[string]interface{})
	if urlParams["textAlignHorizontal"] != "CENTER" {
		t.Errorf("url align = %v, want CENTER", urlParams["textAlignHorizontal"])
	}
	// URL Y should be near bottom
	urlY := urlParams["y"].(float64)
	if urlY < 1200 {
		t.Errorf("url y = %v, expected near bottom of 1350px canvas", urlY)
	}
}

func TestExpandArabicElement(t *testing.T) {
	op := map[string]interface{}{
		"name":    "s1",
		"command": "slide",
		"params": map[string]interface{}{
			"canvas": "1080x1350",
			"color":  "#FFF",
			"elements": []interface{}{
				map[string]interface{}{
					"type": "arabic",
					"text": "بسم الله الرحمن الرحيم",
					"tier": "hero",
				},
			},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	if len(ops) < 2 {
		t.Fatalf("expected at least 2 ops, got %d", len(ops))
	}

	arabicParams, _ := ops[1]["params"].(map[string]interface{})
	if arabicParams["fontFamily"] != "Amiri" {
		t.Errorf("arabic fontFamily = %v, want Amiri", arabicParams["fontFamily"])
	}
}

func TestDefaultBaseName(t *testing.T) {
	// When name is missing, should default to "s1"
	op := map[string]interface{}{
		"command": "slide",
		"params": map[string]interface{}{
			"canvas":   "1080x1350",
			"elements": []interface{}{},
		},
	}

	ops, err := ExpandComposite(op)
	if err != nil {
		t.Fatalf("ExpandComposite failed: %v", err)
	}

	if ops[0]["name"] != "s1_frame" {
		t.Errorf("default frame name = %v, want s1_frame", ops[0]["name"])
	}
}
