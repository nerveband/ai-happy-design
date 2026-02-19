package batchutil

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
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

// truncate shortens a string to max chars, appending "..." if truncated.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
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

func estimateWrappedLinesWithConfig(text string, width, fontSize, widthFactor float64, minCharsPerLine int) float64 {
	if strings.TrimSpace(text) == "" {
		return 1
	}
	if width <= 0 || fontSize <= 0 {
		return 1
	}
	if widthFactor <= 0 {
		widthFactor = 0.68
	}
	charsPerLine := int(math.Floor(width / (fontSize * widthFactor)))
	if minCharsPerLine < 1 {
		minCharsPerLine = 1
	}
	if charsPerLine < minCharsPerLine {
		charsPerLine = minCharsPerLine
	}

	totalLines := 0
	for _, paragraph := range strings.Split(text, "\n") {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			totalLines++
			continue
		}
		lineChars := 0
		lines := 1
		for _, word := range words {
			wlen := utf8.RuneCountInString(word)
			// Break very long tokens that exceed line width.
			if wlen > charsPerLine {
				if lineChars > 0 {
					lines++
					lineChars = 0
				}
				full := wlen / charsPerLine
				lines += full
				lineChars = wlen % charsPerLine
				if lineChars == 0 {
					lineChars = charsPerLine
					lines--
				}
				continue
			}
			needed := wlen
			if lineChars > 0 {
				needed++ // space
			}
			if lineChars+needed > charsPerLine {
				lines++
				lineChars = wlen
			} else {
				lineChars += needed
			}
		}
		totalLines += lines
	}
	if totalLines < 1 {
		totalLines = 1
	}
	return float64(totalLines)
}

// estimateWrappedLines approximates wrapped line count for larger display/body copy.
func estimateWrappedLines(text string, width, fontSize float64) float64 {
	// Keep a conservative floor for large headline/body text so layout stays safe by default.
	return estimateWrappedLinesWithConfig(text, width, fontSize, 0.68, 6)
}

func estimateWrappedLinesTight(text string, width, fontSize float64) float64 {
	// Tight estimate for narrow stat columns where single-token values should still shrink.
	return estimateWrappedLinesWithConfig(text, width, fontSize, 0.62, 1)
}

func estimateTextHeight(text string, width, fontSize, lineHeightPercent float64) float64 {
	lines := estimateWrappedLines(text, width, fontSize)
	return lines * fontSize * (lineHeightPercent / 100.0)
}

func fitHeadlineFontSize(text string, width, canvasH, wantedSize, lineHeightPercent float64) float64 {
	return fitHeadlineFontSizeWithBudget(text, width, canvasH, wantedSize, lineHeightPercent, 0.34, 5)
}

func fitHeadlineFontSizeWithBudget(text string, width, canvasH, wantedSize, lineHeightPercent, maxHeadlineRatio, maxHeadlineLines float64) float64 {
	if maxHeadlineRatio <= 0 || maxHeadlineRatio > 0.9 {
		maxHeadlineRatio = 0.34
	}
	if maxHeadlineLines < 1 {
		maxHeadlineLines = 5
	}
	// Reserve room for other slide elements; headline should not consume most of the canvas.
	maxHeadlineHeight := canvasH * maxHeadlineRatio
	size := wantedSize
	for size > 28 {
		lines := estimateWrappedLines(text, width, size)
		h := lines * size * (lineHeightPercent / 100.0)
		if h <= maxHeadlineHeight && lines <= maxHeadlineLines {
			return size
		}
		size = snap4(size-4, 28)
	}
	return size
}

var headlineTierOrder = []string{"display", "hero", "title", "heading", "subheading", "body", "caption"}

func adaptHeadlineTier(
	text string,
	width, canvasH float64,
	requestedTier string,
	sizes map[string]float64,
	lineHeightPercent, maxHeadlineRatio, maxHeadlineLines float64,
) (string, float64) {
	startTier := strings.ToLower(strings.TrimSpace(requestedTier))
	startIdx := -1
	for i, t := range headlineTierOrder {
		if t == startTier {
			startIdx = i
			break
		}
	}
	if startIdx < 0 {
		startIdx = 2 // title
	}
	if maxHeadlineRatio <= 0 || maxHeadlineRatio > 0.9 {
		maxHeadlineRatio = 0.34
	}
	if maxHeadlineLines < 1 {
		maxHeadlineLines = 5
	}
	maxHeadlineHeight := canvasH * maxHeadlineRatio

	fallbackTier := "title"
	fallbackSize, ok := sizes[fallbackTier]
	if !ok {
		fallbackSize = 64
	}

	for i := startIdx; i < len(headlineTierOrder); i++ {
		tier := headlineTierOrder[i]
		size, hasSize := sizes[tier]
		if !hasSize {
			continue
		}
		lines := estimateWrappedLines(text, width, size)
		height := lines * size * (lineHeightPercent / 100.0)
		if height <= maxHeadlineHeight && lines <= maxHeadlineLines {
			return tier, size
		}
		fallbackTier = tier
		fallbackSize = size
	}
	return fallbackTier, fallbackSize
}

func fitTextToMaxLines(text string, width, wantedSize, lineHeightPercent, maxLines, minSize float64) float64 {
	if maxLines < 1 {
		maxLines = 1
	}
	if minSize <= 0 || minSize > wantedSize {
		minSize = 20
	}
	size := wantedSize
	for size > minSize {
		lines := estimateWrappedLinesTight(text, width, size)
		height := lines * size * (lineHeightPercent / 100.0)
		if lines <= maxLines && height > 0 {
			return size
		}
		size = snap4(size-4, minSize)
	}
	return size
}

func gradientFillType(raw string) string {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "GRADIENT_LINEAR", "LINEAR":
		return "GRADIENT_LINEAR"
	case "GRADIENT_RADIAL", "RADIAL":
		return "GRADIENT_RADIAL"
	case "GRADIENT_ANGULAR", "ANGULAR":
		return "GRADIENT_ANGULAR"
	case "GRADIENT_DIAMOND", "DIAMOND":
		return "GRADIENT_DIAMOND"
	default:
		return "GRADIENT_LINEAR"
	}
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

	// Derive descriptive frame name from first headline element
	elements, _ := params["elements"].([]interface{})
	descriptiveName := baseName
	for _, elem := range elements {
		if em, ok := elem.(map[string]interface{}); ok {
			if et, _ := em["type"].(string); et == "headline" {
				if ht, _ := em["text"].(string); ht != "" {
					descriptiveName = "Slide — " + truncate(ht, 30)
					break
				}
			}
		}
	}

	// Root frame
	frameName := baseName + "_frame"
	frameParams := map[string]interface{}{
		"x":            getFloat(params, "x", 0),
		"y":            getFloat(params, "y", 0),
		"width":        w,
		"height":       h,
		"color":        bgColor,
		"clipsContent": clipsContent,
		"name":         descriptiveName,
	}
	if pid, ok := params["parentId"]; ok {
		frameParams["parentId"] = pid
	}
	ops = append(ops, makeOp(frameName, "node.create_frame", frameParams))

	// Optional background image (local path/URL/base64 supported by resolver).
	if bgImage := getString(params, "bgImage", ""); bgImage != "" {
		imgName := baseName + "_bgimg"
		imgParams := map[string]interface{}{
			"nodeId":    ref(frameName),
			"imageData": bgImage,
		}
		if scaleMode := strings.TrimSpace(getString(params, "bgImageScaleMode", "FILL")); scaleMode != "" {
			imgParams["scaleMode"] = strings.ToUpper(scaleMode)
		}
		ops = append(ops, makeOp(imgName, "paint.set_image", imgParams))
	}

	// Optional gradient
	if grad, ok := params["gradient"].(map[string]interface{}); ok {
		gradName := baseName + "_grad"
		gradParams := cloneParams(grad)
		gradParams["nodeId"] = ref(frameName)
		ops = append(ops, makeOp(gradName, "paint.set_gradient", gradParams))
	}

	// Optional photo overlay fill. This avoids creating overlap-prone child layers.
	if overlayGrad, ok := params["overlayGradient"].(map[string]interface{}); ok {
		overlayFill := map[string]interface{}{
			"nodeId": ref(frameName),
			"type":   gradientFillType(getString(overlayGrad, "type", "LINEAR")),
			"stops":  overlayGrad["stops"],
		}
		ops = append(ops, makeOp(baseName+"_overlay_fill", "paint.add_fill", overlayFill))
	} else if overlayColor := getString(params, "overlayColor", ""); overlayColor != "" {
		overlayFill := map[string]interface{}{
			"nodeId": ref(frameName),
			"type":   "SOLID",
			"color":  overlayColor,
		}
		if opacity := getFloat(params, "overlayOpacity", -1); opacity >= 0 {
			overlayFill["opacity"] = opacity
		}
		ops = append(ops, makeOp(baseName+"_overlay_fill", "paint.add_fill", overlayFill))
	}
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
		return expandHeadline(elemName, frameName, params, canvasH, margin, contentWidth, yPos, sizes)
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
	if explicit := getFloat(params, "fontSize", 0); explicit > 0 {
		fontSize = explicit
	}
	text := getString(params, "text", "")
	color := getString(params, "color", "#999999")
	spacing := getFloat(params, "spacing", 16)
	lineHeight := getFloat(params, "lineHeight", 130)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"name":           "Eyebrow — " + truncate(text, 25),
		"x":              margin,
		"y":              yPos,
		"width":          contentWidth,
		"fontSize":       fontSize,
		"fontFamily":     getString(params, "fontFamily", "Inter"),
		"fontStyle":      getString(params, "fontStyle", "Medium"),
		"color":          color,
		"letterSpacing":  getFloat(params, "letterSpacing", 4),
		"textCase":       "UPPER",
		"lineHeight":     lineHeight,
		"lineHeightUnit": "PERCENT",
	})

	newY := yPos + fontSize*(lineHeight/100.0) + spacing
	return []map[string]interface{}{op}, newY, nil
}

func expandHeadline(name, frameName string, params map[string]interface{}, canvasH, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	tier := getString(params, "tier", "title")
	fontSize, ok := sizes[tier]
	if !ok {
		tier = "title"
		fontSize = sizes["title"]
	}
	text := getString(params, "text", "")
	color := getString(params, "color", "#000000")
	spacing := getFloat(params, "spacing", 24)

	lineHeight := getFloat(params, "lineHeight", 118)
	maxHeadlineRatio := getFloat(params, "maxHeadlineRatio", 0.34)
	maxHeadlineLines := getFloat(params, "maxHeadlineLines", 5)
	explicitFontSize := getFloat(params, "fontSize", 0)
	if explicitFontSize > 0 {
		fontSize = explicitFontSize
	}
	adaptiveTier := getBool(params, "adaptiveTier", true)
	if explicitFontSize > 0 {
		if _, ok := params["adaptiveTier"]; !ok {
			adaptiveTier = false
		}
	}
	if adaptiveTier {
		tier, fontSize = adaptHeadlineTier(text, contentWidth, canvasH, tier, sizes, lineHeight, maxHeadlineRatio, maxHeadlineLines)
	}
	autoFit := getBool(params, "autoFit", true)
	if autoFit {
		fontSize = fitHeadlineFontSizeWithBudget(text, contentWidth, canvasH, fontSize, lineHeight, maxHeadlineRatio, maxHeadlineLines)
	}

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"name":           "Headline — " + truncate(text, 25),
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

	textHeight := estimateTextHeight(text, contentWidth, fontSize, lineHeight)
	newY := yPos + textHeight + spacing
	return []map[string]interface{}{op}, newY, nil
}

func expandBody(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	fontSize := sizes["body"]
	if explicit := getFloat(params, "fontSize", 0); explicit > 0 {
		fontSize = explicit
	}
	text := getString(params, "text", "")
	color := getString(params, "color", "#333333")
	spacing := getFloat(params, "spacing", 24)
	lineHeight := getFloat(params, "lineHeight", 150)

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":       ref(frameName),
		"content":        text,
		"name":           "Body — " + truncate(text, 25),
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

	textHeight := estimateTextHeight(text, contentWidth, fontSize, lineHeight)
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
		"name":     "Accent Bar",
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
	if explicit := getFloat(params, "fontSize", 0); explicit > 0 {
		fontSize = explicit
	}
	if fontSize > sizes["caption"] {
		fontSize = sizes["caption"]
	}
	color := getString(params, "color", "#999999")
	counterWidth := getFloat(params, "width", math.Max(140, fontSize*4))
	lineHeight := getFloat(params, "lineHeight", 130)

	// Counter is positioned top-right, not in the flow
	counterY := margin

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":            ref(frameName),
		"content":             text,
		"name":                "Counter",
		"x":                   canvasW - margin - counterWidth,
		"y":                   counterY,
		"width":               counterWidth,
		"fontSize":            fontSize,
		"fontFamily":          getString(params, "fontFamily", "Inter"),
		"fontStyle":           getString(params, "fontStyle", "Regular"),
		"color":               color,
		"textAlignHorizontal": "RIGHT",
		"lineHeight":          lineHeight,
		"lineHeightUnit":      "PERCENT",
	})

	// Counter doesn't advance yPos — it's absolutely positioned
	return []map[string]interface{}{op}, 0, nil
}

func expandCTA(name, frameName string, params map[string]interface{}, margin, contentWidth, yPos float64, sizes map[string]float64) ([]map[string]interface{}, float64, error) {
	text := getString(params, "text", "Learn More")
	fontSize := sizes["body"]
	if explicit := getFloat(params, "fontSize", 0); explicit > 0 {
		fontSize = explicit
	}
	bgColor := getString(params, "bgColor", "#000000")
	textColor := getString(params, "textColor", "#FFFFFF")
	style := getString(params, "style", "pill")
	spacing := getFloat(params, "spacing", 32)

	btnHeight := snap8(fontSize*2.5, 40)
	btnPadH := snap8(fontSize*1.25, 16)
	estimatedTextWidth := float64(utf8.RuneCountInString(text)) * (fontSize * 0.58)
	btnWidth := snap8(estimatedTextWidth+btnPadH*2, btnHeight*2)
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
		"parentId":              ref(frameName),
		"name":                  "CTA Button",
		"x":                     margin,
		"y":                     yPos,
		"width":                 btnWidth,
		"height":                btnHeight,
		"color":                 bgColor,
		"cornerRadius":          cornerRadius,
		"layoutMode":            "HORIZONTAL",
		"paddingLeft":           btnPadH,
		"paddingRight":          btnPadH,
		"paddingTop":            snap8(fontSize*0.625, 8),
		"paddingBottom":         snap8(fontSize*0.625, 8),
		"primaryAxisAlignItems": "CENTER",
		"counterAxisAlignItems": "CENTER",
		"primaryAxisSizingMode": "FIXED",
		"counterAxisSizingMode": "FIXED",
	})

	// Button text
	btnText := makeOp(btnTextName, "text.create", map[string]interface{}{
		"parentId":       ref(btnFrameName),
		"content":        text,
		"name":           "CTA Text — " + truncate(text, 25),
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
	if explicit := getFloat(params, "fontSize", 0); explicit > 0 {
		fontSize = explicit
	}
	color := getString(params, "color", "#999999")
	margin := snap8(canvasW*0.08, 16)
	lineHeight := getFloat(params, "lineHeight", 130)

	// URL is pinned near bottom
	urlY := canvasH - margin - fontSize*1.5

	op := makeOp(name, "text.create", map[string]interface{}{
		"parentId":            ref(frameName),
		"content":             text,
		"name":                "URL — " + truncate(text, 25),
		"x":                   margin,
		"y":                   urlY,
		"width":               contentWidth,
		"fontSize":            fontSize,
		"fontFamily":          getString(params, "fontFamily", "Inter"),
		"fontStyle":           getString(params, "fontStyle", "Regular"),
		"color":               color,
		"textAlignHorizontal": "CENTER",
		"lineHeight":          lineHeight,
		"lineHeightUnit":      "PERCENT",
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
	valueLineHeight := getFloat(params, "valueLineHeight", 118)
	labelLineHeight := getFloat(params, "labelLineHeight", 130)
	valueMaxLines := getFloat(params, "valueMaxLines", 1)
	labelMaxLines := getFloat(params, "labelMaxLines", 2)
	valueMinSize := getFloat(params, "valueMinSize", 20)
	labelMinSize := getFloat(params, "labelMinSize", 12)
	valueColor := getString(params, "valueColor", "#000000")
	labelColor := getString(params, "labelColor", "#666666")

	colCount := len(items)
	gap := snap8(contentWidth*0.03, 8)
	colWidth := (contentWidth - float64(colCount-1)*gap) / float64(colCount)
	if colWidth < 32 {
		colWidth = 32
	}

	sharedValueSize := valueFontSize
	sharedLabelSize := labelFontSize
	maxValueHeight := 0.0
	maxLabelHeight := 0.0
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		value := getString(itemMap, "value", "0")
		label := getString(itemMap, "label", "")
		valueFit := fitTextToMaxLines(value, colWidth, valueFontSize, valueLineHeight, valueMaxLines, valueMinSize)
		labelFit := fitTextToMaxLines(label, colWidth, labelFontSize, labelLineHeight, labelMaxLines, labelMinSize)
		if valueFit < sharedValueSize {
			sharedValueSize = valueFit
		}
		if labelFit < sharedLabelSize {
			sharedLabelSize = labelFit
		}
		valueHeight := estimateWrappedLinesTight(value, colWidth, valueFit) * valueFit * (valueLineHeight / 100.0)
		labelHeight := estimateWrappedLinesTight(label, colWidth, labelFit) * labelFit * (labelLineHeight / 100.0)
		if valueHeight > maxValueHeight {
			maxValueHeight = valueHeight
		}
		if labelHeight > maxLabelHeight {
			maxLabelHeight = labelHeight
		}
	}
	if maxValueHeight == 0 {
		maxValueHeight = sharedValueSize * (valueLineHeight / 100.0)
	}
	if maxLabelHeight == 0 {
		maxLabelHeight = sharedLabelSize * (labelLineHeight / 100.0)
	}
	labelY := yPos + maxValueHeight + 8

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
			"name":           "Stat Value — " + value,
			"x":              colX,
			"y":              yPos,
			"width":          colWidth,
			"fontSize":       sharedValueSize,
			"fontFamily":     getString(params, "fontFamily", "Inter"),
			"fontStyle":      getString(params, "fontStyle", "Bold"),
			"color":          valueColor,
			"lineHeight":     valueLineHeight,
			"lineHeightUnit": "PERCENT",
		}))

		// Label text
		lblName := fmt.Sprintf("%s_l%d", name, j)
		ops = append(ops, makeOp(lblName, "text.create", map[string]interface{}{
			"parentId":       ref(frameName),
			"content":        label,
			"name":           "Stat Label — " + label,
			"x":              colX,
			"y":              labelY,
			"width":          colWidth,
			"fontSize":       sharedLabelSize,
			"fontFamily":     getString(params, "fontFamily", "Inter"),
			"fontStyle":      getString(params, "fontStyle", "Regular"),
			"color":          labelColor,
			"lineHeight":     labelLineHeight,
			"lineHeightUnit": "PERCENT",
		}))
	}

	newY := yPos + maxValueHeight + 8 + maxLabelHeight + spacing
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
		"name":           "Progress Value",
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
		"parentId":            ref(frameName),
		"content":             goalText,
		"name":                "Progress Goal",
		"x":                   margin,
		"y":                   yPos,
		"width":               contentWidth,
		"fontSize":            labelFontSize,
		"fontFamily":          getString(params, "fontFamily", "Inter"),
		"fontStyle":           "Regular",
		"color":               labelColor,
		"textAlignHorizontal": "RIGHT",
		"lineHeight":          130.0,
		"lineHeightUnit":      "PERCENT",
	}))

	labelRowHeight := sizes["subheading"]*1.3 + 8
	trackY := yPos + labelRowHeight

	// Track (background bar)
	trackName := name + "_track"
	ops = append(ops, makeOp(trackName, "shape.create_rectangle", map[string]interface{}{
		"parentId":     ref(frameName),
		"name":         "Progress Track",
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
		"name":         "Progress Fill",
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
		"name":           "Arabic — " + truncate(text, 25),
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

	textHeight := estimateTextHeight(text, contentWidth, fontSize, lineHeight)
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

	// Derive descriptive frame name from first headline element
	bannerElements, _ := params["elements"].([]interface{})
	bannerDescName := baseName
	for _, elem := range bannerElements {
		if em, ok := elem.(map[string]interface{}); ok {
			if et, _ := em["type"].(string); et == "headline" {
				if ht, _ := em["text"].(string); ht != "" {
					bannerDescName = "Banner — " + truncate(ht, 30)
					break
				}
			}
		}
	}

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
		"name":         bannerDescName,
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
		"name":           "Headline — " + truncate(text, 25),
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
		"name":           "Subtitle — " + truncate(text, 25),
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
