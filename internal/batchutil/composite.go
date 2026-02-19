package batchutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// compositeCommands lists commands that expand into multiple primitive ops.
var compositeCommands = map[string]bool{
	"slide":  true,
	"banner": true,
}

// IsComposite returns true if command is a composite that needs expansion.
func IsComposite(command string) bool {
	return compositeCommands[strings.ToLower(strings.TrimSpace(command))]
}

// ExpandComposite expands a single composite op into primitive batch ops.
// Non-composite commands pass through as single-element slices.
func ExpandComposite(op map[string]interface{}) ([]map[string]interface{}, error) {
	cmd, _ := op["command"].(string)
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	if !IsComposite(cmd) {
		return []map[string]interface{}{op}, nil
	}

	params, _ := op["params"].(map[string]interface{})
	if params == nil {
		params = map[string]interface{}{}
	}
	baseName, _ := op["name"].(string)
	if baseName == "" {
		baseName = "s1"
	}

	switch cmd {
	case "slide":
		return expandSlide(baseName, params)
	case "banner":
		return expandBanner(baseName, params)
	default:
		return []map[string]interface{}{op}, nil
	}
}

// ExpandAllComposites processes a batch array, expanding any composite ops.
func ExpandAllComposites(ops []map[string]interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for i, op := range ops {
		expanded, err := ExpandComposite(op)
		if err != nil {
			return nil, fmt.Errorf("op[%d]: %w", i, err)
		}
		result = append(result, expanded...)
	}
	return result, nil
}

// --- Design token helpers ---

const scaleRatio = 1.333

// tokenSizes computes the modular type scale for a given canvas width.
func tokenSizes(width float64) map[string]float64 {
	rawBase := width * 0.044
	return map[string]float64{
		"caption":    snap4(rawBase/scaleRatio, 12),
		"body":       snap4(rawBase, 14),
		"subheading": snap4(rawBase*scaleRatio, 16),
		"heading":    snap4(rawBase*math.Pow(scaleRatio, 2), 20),
		"title":      snap4(rawBase*math.Pow(scaleRatio, 3), 24),
		"hero":       snap4(rawBase*math.Pow(scaleRatio, 4), 32),
		"display":    snap4(rawBase*math.Pow(scaleRatio, 5), 40),
	}
}

// snap4 rounds to 4px grid, enforcing a minimum.
func snap4(v, min float64) float64 {
	s := math.Round(v/4) * 4
	if s < min {
		return min
	}
	return s
}

// snap8 rounds to 8px grid, enforcing a minimum.
func snap8(v, min float64) float64 {
	s := math.Round(v/8) * 8
	if s < min {
		return min
	}
	return s
}

// parseCanvas parses "WxH" into width, height.
func parseCanvas(canvas string) (float64, float64, error) {
	parts := strings.SplitN(strings.ToLower(canvas), "x", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid canvas %q: expected WxH", canvas)
	}
	w, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid canvas width %q: %w", parts[0], err)
	}
	h, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid canvas height %q: %w", parts[1], err)
	}
	return w, h, nil
}

// ref builds a step interpolation reference.
func ref(stepName string) string {
	return fmt.Sprintf("${{steps.%s.result.id}}", stepName)
}

// getString gets a string param with a default.
func getString(params map[string]interface{}, key, def string) string {
	if v, ok := params[key].(string); ok && v != "" {
		return v
	}
	return def
}

// getFloat gets a float param with a default.
func getFloat(params map[string]interface{}, key string, def float64) float64 {
	switch v := params[key].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return def
		}
		return f
	default:
		return def
	}
}

// getBool gets a bool param with a default.
func getBool(params map[string]interface{}, key string, def bool) bool {
	if v, ok := params[key].(bool); ok {
		return v
	}
	return def
}

// makeOp creates a batch operation map.
func makeOp(name, command string, params map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":    name,
		"command": command,
		"params":  params,
	}
}

// normalizeCompositeParams rewrites common aliases to canonical param names.
// This lets LLMs use "bg", "fillColor", or "pid" in composite params.
func normalizeCompositeParams(params map[string]interface{}) {
	aliases := map[string]string{
		"bg":        "color",
		"fillColor": "color",
		"pid":       "parentId",
	}
	for alias, canonical := range aliases {
		if v, ok := params[alias]; ok {
			if _, exists := params[canonical]; !exists {
				params[canonical] = v
			}
			delete(params, alias)
		}
	}
}

// --- Slide expansion ---

func expandSlide(baseName string, params map[string]interface{}) ([]map[string]interface{}, error) {
	normalizeCompositeParams(params)
	canvas := getString(params, "canvas", "1080x1350")
	w, h, err := parseCanvas(canvas)
	if err != nil {
		return nil, err
	}

	sizes := tokenSizes(w)
	margin := snap8(w*0.08, 16)
	contentWidth := w - 2*margin

	bgColor := getString(params, "color", "#FFFFFF")
	clipsContent := getBool(params, "clipsContent", true)

	var ops []map[string]interface{}

	// Root frame
	frameName := baseName + "_frame"
	frameParams := map[string]interface{}{
		"x":            getFloat(params, "x", 0),
		"y":            getFloat(params, "y", 0),
		"width":        w,
		"height":       h,
		"color":        bgColor,
		"clipsContent": clipsContent,
	}
	if pid, ok := params["parentId"]; ok {
		frameParams["parentId"] = pid
	}
	ops = append(ops, makeOp(frameName, "node.create_frame", frameParams))

	// Optional gradient
	if grad, ok := params["gradient"].(map[string]interface{}); ok {
		gradName := baseName + "_grad"
		gradParams := cloneParams(grad)
		gradParams["nodeId"] = ref(frameName)
		ops = append(ops, makeOp(gradName, "paint.set_gradient", gradParams))
	}

	// Process elements
	elements, _ := params["elements"].([]interface{})
	yPos := snap8(h*0.10, 16) // start with top padding

	for i, elem := range elements {
		elemMap, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}
		elemType, _ := elemMap["type"].(string)
		elemName := fmt.Sprintf("%s_e%d", baseName, i)

		expanded, newY, err := expandSlideElement(elemName, frameName, elemType, elemMap, w, h, margin, contentWidth, yPos, sizes)
		if err != nil {
			return nil, fmt.Errorf("element[%d] %q: %w", i, elemType, err)
		}
		ops = append(ops, expanded...)
		if newY > 0 {
			yPos = newY
		}
	}

	return ops, nil
}

func expandSlideElement(
	elemName, frameName, elemType string,
	params map[string]interface{},
	canvasW, canvasH, margin, contentWidth, yPos float64,
	sizes map[string]float64,
) ([]map[string]interface{}, float64, error) {

	switch elemType {
	case "eyebrow":
		return expandEyebrow(elemName, frameName, params, margin, contentWidth, yPos, sizes)
	case "headline":
		return expandHeadline(elemName, frameName, params, margin, contentWidth, yPos, sizes)
	case "body":
		return expandBody(elemName, frameName, params, margin, contentWidth, yPos, sizes)
	case "bar":
		return expandBar(elemName, frameName, params, margin, yPos, sizes)
	case "counter":
		return expandCounter(elemName, frameName, params, canvasW, margin, sizes)
	case "cta":
		return expandCTA(elemName, frameName, params, margin, contentWidth, yPos, sizes)
	case "url":
		return expandURL(elemName, frameName, params, canvasW, canvasH, contentWidth, sizes)
	case "stats":
		return expandStats(elemName, frameName, params, margin, contentWidth, yPos, sizes)
	case "progress":
		return expandProgress(elemName, frameName, params, margin, contentWidth, yPos, sizes)
	case "arabic":
		return expandArabic(elemName, frameName, params, margin, contentWidth, yPos, sizes)
	default:
		return nil, yPos, fmt.Errorf("unknown element type %q", elemType)
	}
}

func expandEyebrow(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	fontSize := sizes["caption"]
	text := getString(params, "text", "")
	color := getString(params, "color", "#999999")
	spacing := getFloat(params, "spacing", 16)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":      ref(frameName),
		"content":       text,
		"x":             margin,
		"y":             yPos,
		"width":         contentWidth,
		"fontSize":      fontSize,
		"fontFamily":    getString(params, "fontFamily", "Inter"),
		"fontStyle":     getString(params, "fontStyle", "Medium"),
		"color":         color,
		"letterSpacing": getFloat(params, "letterSpacing", 4),
		"textCase":      "UPPER",
		"lineHeight":    130.0,
		"lineHeightUnit": "PERCENT",
	})

	newY := yPos + fontSize*1.3 + spacing
	return []map[string]interface{}{op}, newY, nil
}

func expandHeadline(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	tier := getString(params, "tier", "title")
	fontSize, ok := sizes[tier]
	if !ok {
		fontSize = sizes["title"]
	}
	text := getString(params, "text", "")
	color := getString(params, "color", "#000000")
	spacing := getFloat(params, "spacing", 24)

	lineHeight := getFloat(params, "lineHeight", 118)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"x":              margin,
		"y":              yPos,
		"width":          contentWidth,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "Inter"),
		"fontStyle":      getString(params, "fontStyle", "Bold"),
		"color":          color,
		"lineHeight":     lineHeight,
		"lineHeightUnit": "PERCENT",
	})

	// Estimate text height: fontSize * lineHeight/100 * approximate lines
	textHeight := fontSize * lineHeight / 100
	newY := yPos + textHeight + spacing
	return []map[string]interface{}{op}, newY, nil
}

func expandBody(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	fontSize := sizes["body"]
	text := getString(params, "text", "")
	color := getString(params, "color", "#333333")
	spacing := getFloat(params, "spacing", 24)
	lineHeight := getFloat(params, "lineHeight", 150)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"x":              margin,
		"y":              yPos,
		"width":          contentWidth,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "DM Sans"),
		"fontStyle":      getString(params, "fontStyle", "Regular"),
		"color":          color,
		"lineHeight":     lineHeight,
		"lineHeightUnit": "PERCENT",
	})

	textHeight := fontSize * lineHeight / 100
	newY := yPos + textHeight + spacing
	return []map[string]interface{}{op}, newY, nil
}

func expandBar(name, frameName string, params map[string]interface{}, margin, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	barWidth := getFloat(params, "width", 108)
	barHeight := getFloat(params, "height", 4)
	color := getString(params, "color", "#FFD600")
	spacing := getFloat(params, "spacing", 24)

	op := makeOp(name, "shape.create_rectangle", map[string]interface{}{
		"parentId": ref(frameName),
		"x":        margin,
		"y":        yPos,
		"width":    barWidth,
		"height":   barHeight,
		"color":    color,
	})

	newY := yPos + barHeight + spacing
	return []map[string]interface{}{op}, newY, nil
}

func expandCounter(name, frameName string, params map[string]interface{}, canvasW, margin float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	current := getString(params, "current", "1")
	total := getString(params, "total", "5")
	text := current + " / " + total
	fontSize := sizes["caption"]
	color := getString(params, "color", "#999999")

	// Counter is positioned top-right, not in the flow
	counterY := margin

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"x":              canvasW - margin - 100,
		"y":              counterY,
		"width":          100.0,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "Inter"),
		"fontStyle":      getString(params, "fontStyle", "Regular"),
		"color":          color,
		"textAlignHorizontal": "RIGHT",
		"lineHeight":     130.0,
		"lineHeightUnit": "PERCENT",
	})

	// Counter doesn't advance yPos — it's absolutely positioned
	return []map[string]interface{}{op}, 0, nil
}

func expandCTA(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	text := getString(params, "text", "Learn More")
	fontSize := sizes["body"]
	bgColor := getString(params, "bgColor", "#000000")
	textColor := getString(params, "textColor", "#FFFFFF")
	style := getString(params, "style", "pill")
	spacing := getFloat(params, "spacing", 32)

	btnHeight := snap8(fontSize*2.5, 40)
	btnPadH := snap8(fontSize*1.25, 16)
	var cornerRadius float64
	if style == "pill" {
		cornerRadius = btnHeight / 2
	} else {
		cornerRadius = snap4(btnHeight*0.28, 8)
	}

	btnFrameName := name + "_btn"
	btnTextName := name + "_txt"

	// Button frame with auto-layout
	btnFrame := makeOp(btnFrameName, "node.create_frame", map[string]interface{}{
		"parentId":     ref(frameName),
		"x":            margin,
		"y":            yPos,
		"color":        bgColor,
		"cornerRadius": cornerRadius,
		"layoutMode":   "HORIZONTAL",
		"paddingLeft":  btnPadH,
		"paddingRight": btnPadH,
		"paddingTop":   snap8(fontSize*0.625, 8),
		"paddingBottom": snap8(fontSize*0.625, 8),
		"primaryAxisAlignItems":   "CENTER",
		"counterAxisAlignItems":   "CENTER",
		"primaryAxisSizingMode":   "AUTO",
		"counterAxisSizingMode":   "AUTO",
	})

	// Button text
	btnText := makeOp(btnTextName, "text.create", map[string]interface{}{
		"parentId":       ref(btnFrameName),
		"content":        text,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "Inter"),
		"fontStyle":      getString(params, "fontStyle", "Bold"),
		"color":          textColor,
		"lineHeight":     130.0,
		"lineHeightUnit": "PERCENT",
	})

	newY := yPos + btnHeight + spacing
	return []map[string]interface{}{btnFrame, btnText}, newY, nil
}

func expandURL(name, frameName string, params map[string]interface{}, canvasW, canvasH, contentWidth float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	text := getString(params, "text", "")
	fontSize := sizes["caption"]
	color := getString(params, "color", "#999999")
	margin := snap8(canvasW*0.08, 16)

	// URL is pinned near bottom
	urlY := canvasH - margin - fontSize*1.5

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":              ref(frameName),
		"content":               text,
		"x":                     margin,
		"y":                     urlY,
		"width":                 contentWidth,
		"fontSize":              fontSize,
		"fontFamily":            getString(params, "fontFamily", "Inter"),
		"fontStyle":             getString(params, "fontStyle", "Regular"),
		"color":                 color,
		"textAlignHorizontal":   "CENTER",
		"lineHeight":            130.0,
		"lineHeightUnit":        "PERCENT",
	})

	// URL doesn't advance yPos — it's bottom-pinned
	return []map[string]interface{}{op}, 0, nil
}

func expandStats(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	items, _ := params["items"].([]interface{})
	if len(items) == 0 {
		return nil, yPos, nil
	}

	spacing := getFloat(params, "spacing", 32)
	valueFontSize := sizes["display"]
	labelFontSize := sizes["caption"]
	valueColor := getString(params, "valueColor", "#000000")
	labelColor := getString(params, "labelColor", "#666666")

	colCount := len(items)
	gap := snap8(contentWidth*0.03, 8)
	colWidth := (contentWidth - float64(colCount-1)*gap) / float64(colCount)

	var ops []map[string]interface{}
	for j, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		value := getString(itemMap, "value", "0")
		label := getString(itemMap, "label", "")

		colX := margin + float64(j)*(colWidth+gap)

		// Value text
		valName := fmt.Sprintf("%s_v%d", name, j)
		ops = append(ops, makeOp(valName, "text.create", map[string]interface{}{
			"parentId":       ref(frameName),
			"content":        value,
			"x":              colX,
			"y":              yPos,
			"width":          colWidth,
			"fontSize":       valueFontSize,
			"fontFamily":     getString(params, "fontFamily", "Inter"),
			"fontStyle":      getString(params, "fontStyle", "Bold"),
			"color":          valueColor,
			"lineHeight":     118.0,
			"lineHeightUnit": "PERCENT",
		}))

		// Label text
		lblName := fmt.Sprintf("%s_l%d", name, j)
		labelY := yPos + valueFontSize*1.18 + 8
		ops = append(ops, makeOp(lblName, "text.create", map[string]interface{}{
			"parentId":       ref(frameName),
			"content":        label,
			"x":              colX,
			"y":              labelY,
			"width":          colWidth,
			"fontSize":       labelFontSize,
			"fontFamily":     getString(params, "fontFamily", "Inter"),
			"fontStyle":      getString(params, "fontStyle", "Regular"),
			"color":          labelColor,
			"lineHeight":     130.0,
			"lineHeightUnit": "PERCENT",
		}))
	}

	newY := yPos + valueFontSize*1.18 + 8 + labelFontSize*1.3 + spacing
	return ops, newY, nil
}

func expandProgress(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	current := getFloat(params, "current", 50)
	goal := getFloat(params, "goal", 100)
	trackColor := getString(params, "trackColor", "#E0E0E0")
	fillColor := getString(params, "fillColor", "#000000")
	barHeight := getFloat(params, "barHeight", 16)
	cornerRadius := barHeight / 2
	spacing := getFloat(params, "spacing", 32)
	labelFontSize := sizes["caption"]
	labelColor := getString(params, "labelColor", "#666666")

	var ops []map[string]interface{}

	// "Raised" label (current value) above left
	raisedName := name + "_raised"
	raisedText := getString(params, "raisedText", fmt.Sprintf("%.0f", current))
	ops = append(ops, makeOp(raisedName, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        raisedText,
		"x":              margin,
		"y":              yPos,
		"fontSize":       sizes["subheading"],
		"fontFamily":     getString(params, "fontFamily", "Inter"),
		"fontStyle":      "Bold",
		"color":          getString(params, "raisedColor", "#000000"),
		"lineHeight":     130.0,
		"lineHeightUnit": "PERCENT",
	}))

	// Goal label above right
	goalName := name + "_goal"
	goalText := getString(params, "goalText", fmt.Sprintf("Goal: %.0f", goal))
	ops = append(ops, makeOp(goalName, "text.create", map[string]interface{}{
		"parentId":              ref(frameName),
		"content":               goalText,
		"x":                     margin,
		"y":                     yPos,
		"width":                 contentWidth,
		"fontSize":              labelFontSize,
		"fontFamily":            getString(params, "fontFamily", "Inter"),
		"fontStyle":             "Regular",
		"color":                 labelColor,
		"textAlignHorizontal":   "RIGHT",
		"lineHeight":            130.0,
		"lineHeightUnit":        "PERCENT",
	}))

	labelRowHeight := sizes["subheading"]*1.3 + 8
	trackY := yPos + labelRowHeight

	// Track (background bar)
	trackName := name + "_track"
	ops = append(ops, makeOp(trackName, "shape.create_rectangle", map[string]interface{}{
		"parentId":     ref(frameName),
		"x":            margin,
		"y":            trackY,
		"width":        contentWidth,
		"height":       barHeight,
		"color":        trackColor,
		"cornerRadius": cornerRadius,
	}))

	// Fill (progress bar)
	pct := current / goal
	if pct > 1 {
		pct = 1
	}
	if pct < 0 {
		pct = 0
	}
	fillWidth := contentWidth * pct
	if fillWidth < barHeight {
		fillWidth = barHeight // min width for pill shape
	}

	fillName := name + "_fill"
	ops = append(ops, makeOp(fillName, "shape.create_rectangle", map[string]interface{}{
		"parentId":     ref(frameName),
		"x":            margin,
		"y":            trackY,
		"width":        fillWidth,
		"height":       barHeight,
		"color":        fillColor,
		"cornerRadius": cornerRadius,
	}))

	newY := trackY + barHeight + spacing
	return ops, newY, nil
}

func expandArabic(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	tier := getString(params, "tier", "hero")
	fontSize, ok := sizes[tier]
	if !ok {
		fontSize = sizes["hero"]
	}
	text := getString(params, "text", "")
	color := getString(params, "color", "#000000")
	spacing := getFloat(params, "spacing", 32)
	lineHeight := getFloat(params, "lineHeight", 150)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"x":              margin,
		"y":              yPos,
		"width":          contentWidth,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "Amiri"),
		"fontStyle":      getString(params, "fontStyle", "Regular"),
		"color":          color,
		"lineHeight":     lineHeight,
		"lineHeightUnit": "PERCENT",
	})

	textHeight := fontSize * lineHeight / 100
	newY := yPos + textHeight + spacing
	return []map[string]interface{}{op}, newY, nil
}

// --- Banner expansion ---

func expandBanner(baseName string, params map[string]interface{}) ([]map[string]interface{}, error) {
	normalizeCompositeParams(params)
	canvas := getString(params, "canvas", "1200x628")
	w, h, err := parseCanvas(canvas)
	if err != nil {
		return nil, err
	}

	sizes := tokenSizes(w)
	margin := snap8(w*0.08, 16)

	bgColor := getString(params, "color", "#FFFFFF")

	var ops []map[string]interface{}

	// Root frame
	frameName := baseName + "_frame"
	frameParams := map[string]interface{}{
		"x":            getFloat(params, "x", 0),
		"y":            getFloat(params, "y", 0),
		"width":        w,
		"height":       h,
		"color":        bgColor,
		"clipsContent": getBool(params, "clipsContent", true),
	}
	if pid, ok := params["parentId"]; ok {
		frameParams["parentId"] = pid
	}
	ops = append(ops, makeOp(frameName, "node.create_frame", frameParams))

	// Optional gradient
	if grad, ok := params["gradient"].(map[string]interface{}); ok {
		gradName := baseName + "_grad"
		gradParams := cloneParams(grad)
		gradParams["nodeId"] = ref(frameName)
		ops = append(ops, makeOp(gradName, "paint.set_gradient", gradParams))
	}

	// Optional vertical divider
	dividerX := getFloat(params, "dividerX", 0)
	contentStartX := margin
	contentWidth := w - 2*margin

	if dividerX > 0 {
		divName := baseName + "_div"
		divColor := getString(params, "dividerColor", "#E0E0E0")
		divWidth := getFloat(params, "dividerWidth", 2)
		ops = append(ops, makeOp(divName, "shape.create_rectangle", map[string]interface{}{
			"parentId": ref(frameName),
			"x":        dividerX,
			"y":        margin,
			"width":    divWidth,
			"height":   h - 2*margin,
			"color":    divColor,
		}))
		// Content goes after divider
		contentStartX = dividerX + getFloat(params, "dividerWidth", 2) + margin
		contentWidth = w - contentStartX - margin
	}

	// Process elements
	elements, _ := params["elements"].([]interface{})
	yPos := h / 2 // Start vertically centered for banners

	// Estimate total content height first for centering
	if len(elements) > 0 {
		yPos = snap8(h*0.3, 16)
	}

	for i, elem := range elements {
		elemMap, ok := elem.(map[string]interface{})
		if !ok {
			continue
		}
		elemType, _ := elemMap["type"].(string)
		elemName := fmt.Sprintf("%s_e%d", baseName, i)

		expanded, newY, err := expandBannerElement(elemName, frameName, elemType, elemMap, contentStartX, contentWidth, yPos, sizes)
		if err != nil {
			return nil, fmt.Errorf("element[%d] %q: %w", i, elemType, err)
		}
		ops = append(ops, expanded...)
		if newY > 0 {
			yPos = newY
		}
	}

	return ops, nil
}

func expandBannerElement(
	elemName, frameName, elemType string,
	params map[string]interface{},
	contentX, contentWidth, yPos float64,
	sizes map[string]float64,
) ([]map[string]interface{}, float64, error) {
	switch elemType {
	case "headline":
		return expandBannerHeadline(elemName, frameName, params, contentX, contentWidth, yPos, sizes)
	case "subtitle":
		return expandBannerSubtitle(elemName, frameName, params, contentX, contentWidth, yPos, sizes)
	default:
		return nil, yPos, fmt.Errorf("unknown banner element type %q", elemType)
	}
}

func expandBannerHeadline(name, frameName string, params map[string]interface{}, x, width, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	tier := getString(params, "tier", "heading")
	fontSize, ok := sizes[tier]
	if !ok {
		fontSize = sizes["heading"]
	}
	text := getString(params, "text", "")
	color := getString(params, "color", "#000000")
	spacing := getFloat(params, "spacing", 16)
	lineHeight := getFloat(params, "lineHeight", 118)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"x":              x,
		"y":              yPos,
		"width":          width,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "Inter"),
		"fontStyle":      getString(params, "fontStyle", "Bold"),
		"color":          color,
		"lineHeight":     lineHeight,
		"lineHeightUnit": "PERCENT",
	})

	newY := yPos + fontSize*lineHeight/100 + spacing
	return []map[string]interface{}{op}, newY, nil
}

func expandBannerSubtitle(name, frameName string, params map[string]interface{}, x, width, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	fontSize := sizes["subheading"]
	text := getString(params, "text", "")
	color := getString(params, "color", "#666666")
	spacing := getFloat(params, "spacing", 16)
	lineHeight := getFloat(params, "lineHeight", 140)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"x":              x,
		"y":              yPos,
		"width":          width,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "Inter"),
		"fontStyle":      getString(params, "fontStyle", "Regular"),
		"color":          color,
		"lineHeight":     lineHeight,
		"lineHeightUnit": "PERCENT",
	})

	newY := yPos + fontSize*lineHeight/100 + spacing
	return []map[string]interface{}{op}, newY, nil
}
