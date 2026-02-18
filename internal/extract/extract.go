// Package extract provides a lightweight HTML+CSS parser that converts
// styled HTML slide/banner documents into Figma batch JSON operations.
//
// The parser targets a specific HTML structure used for social media
// carousel previews, where CSS classes define backgrounds, typography,
// and layout, and the DOM structure maps to slides and email banners.
package extract

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Options configures the HTML-to-Figma extraction.
type Options struct {
	CanvasWidth  int // target Figma canvas width (e.g. 1080 for slides)
	CanvasHeight int // target Figma canvas height (e.g. 1350 for slides)
}

// FromHTML parses HTML+CSS and returns batch operations (composite commands).
// Each slide becomes a "slide" composite op, each email banner a "banner" op.
func FromHTML(r io.Reader, opts Options) ([]map[string]interface{}, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading HTML: %w", err)
	}
	htmlStr := string(raw)

	// 1. Parse CSS classes from <style> blocks
	styles := parseStyleBlocks(htmlStr)

	// 2. Parse DOM
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	// 3. Find slides and banners
	var ops []map[string]interface{}
	slideIdx := 0
	bannerIdx := 0

	// Walk DOM looking for slide rows and email banner containers
	walkDOM(doc, func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return true // continue
		}
		classes := getAttr(n, "class")

		// Check for .slides-row container
		if hasClass(classes, "slides-row") {
			slides := extractSlidesFromRow(n, styles, opts, &slideIdx)
			ops = append(ops, slides...)
			return false // don't recurse into this node
		}

		// Check for .email-banners container
		if hasClass(classes, "email-banners") {
			banners := extractBannersFromContainer(n, styles, opts, &bannerIdx)
			ops = append(ops, banners...)
			return false
		}

		// Check for standalone .slide not inside a .slides-row
		if hasClass(classes, "slide") {
			op := extractSlide(n, styles, opts, slideIdx)
			if op != nil {
				ops = append(ops, op)
				slideIdx++
			}
			return false
		}

		// Check for standalone .eb or .email-banner
		if hasClass(classes, "eb") || hasClass(classes, "email-banner") {
			op := extractBanner(n, styles, opts, bannerIdx)
			if op != nil {
				ops = append(ops, op)
				bannerIdx++
			}
			return false
		}

		return true // continue walking
	})

	return ops, nil
}

// walkDOM walks the DOM tree, calling fn for each node.
// If fn returns false, children are not visited.
func walkDOM(n *html.Node, fn func(*html.Node) bool) {
	if n == nil {
		return
	}
	if !fn(n) {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkDOM(c, fn)
	}
}

// extractSlidesFromRow finds .slide children inside a .slides-row.
func extractSlidesFromRow(row *html.Node, styles map[string]map[string]string, opts Options, idx *int) []map[string]interface{} {
	var ops []map[string]interface{}
	for c := row.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		classes := getAttr(c, "class")
		if hasClass(classes, "slide") {
			op := extractSlide(c, styles, opts, *idx)
			if op != nil {
				ops = append(ops, op)
				(*idx)++
			}
		}
	}
	return ops
}

// extractBannersFromContainer finds .eb children inside .email-banners.
func extractBannersFromContainer(container *html.Node, styles map[string]map[string]string, opts Options, idx *int) []map[string]interface{} {
	var ops []map[string]interface{}
	for c := container.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		classes := getAttr(c, "class")
		if hasClass(classes, "eb") {
			op := extractBanner(c, styles, opts, *idx)
			if op != nil {
				ops = append(ops, op)
				(*idx)++
			}
		}
	}
	return ops
}

// extractSlide converts a .slide DOM node into a composite "slide" command.
func extractSlide(node *html.Node, styles map[string]map[string]string, opts Options, idx int) map[string]interface{} {
	classes := getAttr(node, "class")
	inlineStyle := parseInlineStyle(getAttr(node, "style"))

	// Determine canvas size
	canvasW := opts.CanvasWidth
	canvasH := opts.CanvasHeight
	if canvasW == 0 {
		canvasW = 1080
	}
	if canvasH == 0 {
		canvasH = 1350
	}

	// Determine HTML preview size from CSS for scale factor
	htmlW := getHTMLDimension(classes, styles, "width", 280)
	scaleFactor := float64(canvasW) / htmlW

	// Resolve background
	merged := mergeClassStyles(classes, styles)
	for k, v := range inlineStyle {
		merged[k] = v
	}

	bgColor, gradient := resolveBackground(merged)

	params := map[string]interface{}{
		"canvas": fmt.Sprintf("%dx%d", canvasW, canvasH),
	}
	if bgColor != "" {
		params["color"] = bgColor
	}
	if gradient != nil {
		params["gradient"] = gradient
	}

	// Extract child elements
	elements := extractSlideElements(node, styles, scaleFactor)
	if len(elements) > 0 {
		params["elements"] = elements
	}

	return map[string]interface{}{
		"name":    fmt.Sprintf("s%d", idx+1),
		"command": "slide",
		"params":  params,
	}
}

// extractBanner converts a .eb DOM node into a composite "banner" command.
func extractBanner(node *html.Node, styles map[string]map[string]string, opts Options, idx int) map[string]interface{} {
	classes := getAttr(node, "class")
	inlineStyle := parseInlineStyle(getAttr(node, "style"))

	canvasW := 1200
	canvasH := 400
	if opts.CanvasWidth > 0 && opts.CanvasHeight > 0 {
		canvasW = opts.CanvasWidth
		canvasH = opts.CanvasHeight
	}

	htmlW := getHTMLDimension(classes, styles, "width", 600)
	scaleFactor := float64(canvasW) / htmlW

	merged := mergeClassStyles(classes, styles)
	for k, v := range inlineStyle {
		merged[k] = v
	}

	bgColor, gradient := resolveBackground(merged)

	// Also check child overlay elements (.eb-ov) for inline gradient
	if gradient == nil {
		walkDOM(node, func(n *html.Node) bool {
			if n == node {
				return true
			}
			if hasClass(getAttr(n, "class"), "eb-ov") {
				ovStyle := parseInlineStyle(getAttr(n, "style"))
				_, g := resolveBackground(ovStyle)
				if g != nil {
					gradient = g
				}
				return false
			}
			return true
		})
	}

	params := map[string]interface{}{
		"canvas": fmt.Sprintf("%dx%d", canvasW, canvasH),
	}
	if bgColor != "" {
		params["color"] = bgColor
	}
	if gradient != nil {
		params["gradient"] = gradient
	}

	// Extract child elements
	elements := extractBannerElements(node, styles, scaleFactor)
	if len(elements) > 0 {
		params["elements"] = elements
	}

	return map[string]interface{}{
		"name":    fmt.Sprintf("b%d", idx+1),
		"command": "banner",
		"params":  params,
	}
}

// extractSlideElements walks slide children and maps them to composite element types.
func extractSlideElements(node *html.Node, styles map[string]map[string]string, scale float64) []interface{} {
	var elements []interface{}

	walkDOM(node, func(n *html.Node) bool {
		if n.Type != html.ElementNode || n == node {
			return true
		}

		classes := getAttr(n, "class")
		inlineStyle := parseInlineStyle(getAttr(n, "style"))
		merged := mergeClassStyles(classes, styles)
		for k, v := range inlineStyle {
			merged[k] = v
		}

		// Counter (slide number)
		if hasClass(classes, "s-num") {
			text := textContent(n)
			parts := strings.SplitN(text, "/", 2)
			if len(parts) == 2 {
				elements = append(elements, map[string]interface{}{
					"type":    "counter",
					"current": strings.TrimSpace(parts[0]),
					"total":   strings.TrimSpace(parts[1]),
					"color":   resolveTextColor(merged),
				})
			}
			return false
		}

		// Eyebrow text (.ey class)
		if hasClass(classes, "ey") {
			elem := map[string]interface{}{
				"type": "eyebrow",
				"text": textContent(n),
			}
			applyTextStyle(elem, merged, scale)
			elements = append(elements, elem)
			return false
		}

		// Headline (.h class)
		if hasClass(classes, "h") {
			elem := map[string]interface{}{
				"type": "headline",
				"text": textContentWithBreaks(n),
			}
			applyTextStyle(elem, merged, scale)
			elem["tier"] = inferHeadlineTier(classes)
			if merged["font-weight"] == "700" || merged["font-weight"] == "bold" {
				elem["fontStyle"] = "Bold"
			}
			elements = append(elements, elem)
			return false
		}

		// Decorative bar (.bar class)
		if hasClass(classes, "bar") {
			elem := map[string]interface{}{
				"type": "bar",
			}
			if c := resolveBarColor(classes, merged); c != "" {
				elem["color"] = c
			}
			elements = append(elements, elem)
			return false
		}

		// URL element (.s-url class)
		if hasClass(classes, "s-url") {
			elem := map[string]interface{}{
				"type": "url",
				"text": textContent(n),
			}
			if c := resolveTextColor(merged); c != "" {
				elem["color"] = c
			}
			elements = append(elements, elem)
			return false
		}

		// Body text (.s-body or .body class or generic paragraph)
		if hasClass(classes, "s-body") || hasClass(classes, "body") || hasClass(classes, "s-p") {
			elem := map[string]interface{}{
				"type": "body",
				"text": textContentWithBreaks(n),
			}
			applyTextStyle(elem, merged, scale)
			elements = append(elements, elem)
			return false
		}

		// CTA / button (.s-cta or .cta class)
		if hasClass(classes, "s-cta") || hasClass(classes, "cta") || hasClass(classes, "btn") {
			elem := map[string]interface{}{
				"type": "cta",
				"text": textContent(n),
			}
			if bg := merged["background-color"]; bg != "" {
				hex := cssColorToHex(bg)
				if hex != "" {
					elem["bgColor"] = hex
				}
			} else if bg := merged["background"]; bg != "" && !strings.Contains(bg, "gradient") {
				hex := cssColorToHex(bg)
				if hex != "" {
					elem["bgColor"] = hex
				}
			}
			if c := resolveTextColor(merged); c != "" {
				elem["textColor"] = c
			}
			elements = append(elements, elem)
			return false
		}

		// Stats container (.stats)
		if hasClass(classes, "stats") || hasClass(classes, "s-stats") {
			items := extractStatsItems(n, styles, scale)
			if len(items) > 0 {
				elements = append(elements, map[string]interface{}{
					"type":  "stats",
					"items": items,
				})
			}
			return false
		}

		// Content wrapper (.s-c) — recurse into it
		if hasClass(classes, "s-c") {
			return true
		}

		return true
	})

	return elements
}

// extractBannerElements walks banner children and maps them to composite element types.
func extractBannerElements(node *html.Node, styles map[string]map[string]string, scale float64) []interface{} {
	var elements []interface{}

	walkDOM(node, func(n *html.Node) bool {
		if n.Type != html.ElementNode || n == node {
			return true
		}

		classes := getAttr(n, "class")
		inlineStyle := parseInlineStyle(getAttr(n, "style"))
		merged := mergeClassStyles(classes, styles)
		for k, v := range inlineStyle {
			merged[k] = v
		}

		// Banner headline (.eb-h)
		if hasClass(classes, "eb-h") {
			elem := map[string]interface{}{
				"type": "headline",
				"text": textContentWithBreaks(n),
			}
			applyTextStyle(elem, merged, scale)
			elem["tier"] = inferBannerHeadlineTier(classes)
			elements = append(elements, elem)
			return false
		}

		// Banner subtitle (.eb-sub)
		if hasClass(classes, "eb-sub") {
			elem := map[string]interface{}{
				"type": "subtitle",
				"text": textContent(n),
			}
			applyTextStyle(elem, merged, scale)
			elements = append(elements, elem)
			return false
		}

		// Banner text wrapper (.eb-text) — recurse
		if hasClass(classes, "eb-text") {
			return true
		}

		return true
	})

	return elements
}

// --- Helper functions ---

// getAttr returns the value of the named attribute, or "".
func getAttr(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// hasClass checks if the class string contains the given class name.
func hasClass(classes, name string) bool {
	for _, c := range strings.Fields(classes) {
		if c == name {
			return true
		}
	}
	return false
}

// textContent returns the concatenated text content of a node, trimmed.
func textContent(n *html.Node) string {
	var sb strings.Builder
	collectText(n, &sb, false)
	return strings.TrimSpace(sb.String())
}

// textContentWithBreaks returns text content, converting <br> to \n.
func textContentWithBreaks(n *html.Node) string {
	var sb strings.Builder
	collectText(n, &sb, true)
	return strings.TrimSpace(sb.String())
}

func collectText(n *html.Node, sb *strings.Builder, preserveBreaks bool) {
	if n.Type == html.TextNode {
		// Collapse whitespace
		text := strings.Join(strings.Fields(n.Data), " ")
		sb.WriteString(text)
	}
	if n.Type == html.ElementNode && n.Data == "br" && preserveBreaks {
		sb.WriteString("\n")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, sb, preserveBreaks)
	}
}

// mergeClassStyles merges CSS properties from all classes on an element.
func mergeClassStyles(classes string, styles map[string]map[string]string) map[string]string {
	result := make(map[string]string)
	for _, cls := range strings.Fields(classes) {
		key := "." + cls
		if props, ok := styles[key]; ok {
			for k, v := range props {
				result[k] = v
			}
		}
	}
	return result
}

// resolveBackground extracts background color and/or gradient from CSS properties.
func resolveBackground(props map[string]string) (string, map[string]interface{}) {
	// Check background-color first
	bgColor := ""
	if bc := props["background-color"]; bc != "" {
		bgColor = cssColorToHex(bc)
	}

	// Check background (could be color or gradient)
	bg := props["background"]
	if bg != "" {
		if strings.Contains(bg, "linear-gradient") {
			angle, stops, err := parseLinearGradient(bg)
			if err == nil && len(stops) >= 2 {
				// Use first stop color as base
				if bgColor == "" {
					bgColor = stops[0].Color
				}
				gradStops := make([]interface{}, len(stops))
				for i, s := range stops {
					gradStops[i] = map[string]interface{}{
						"color":    s.Color,
						"position": s.Position,
					}
				}
				gradient := map[string]interface{}{
					"type":  "LINEAR",
					"angle": angle,
					"stops": gradStops,
				}
				return bgColor, gradient
			}
		} else {
			// Plain color
			hex := cssColorToHex(bg)
			if hex != "" && bgColor == "" {
				bgColor = hex
			}
		}
	}

	if bgColor == "" {
		bgColor = "#FFFFFF"
	}
	return bgColor, nil
}

// resolveTextColor gets color from CSS properties.
func resolveTextColor(props map[string]string) string {
	if c := props["color"]; c != "" {
		return cssColorToHex(c)
	}
	return ""
}

// resolveBarColor determines bar color from CSS classes or properties.
func resolveBarColor(classes string, props map[string]string) string {
	// Check for specific bar color classes
	if hasClass(classes, "bar-grn") || hasClass(classes, "bar-green") {
		return "#22C55E" // green
	}
	if hasClass(classes, "bar-gold") || hasClass(classes, "bar-ylw") {
		return "#FFD600"
	}
	if hasClass(classes, "bar-sky") || hasClass(classes, "bar-blue") {
		return "#38BDF8"
	}
	if hasClass(classes, "bar-wht") || hasClass(classes, "bar-white") {
		return "#FFFFFF"
	}

	// Check CSS background
	if bg := props["background-color"]; bg != "" {
		return cssColorToHex(bg)
	}
	if bg := props["background"]; bg != "" && !strings.Contains(bg, "gradient") {
		return cssColorToHex(bg)
	}

	return "#FFD600" // default
}

// applyTextStyle populates element map with font properties from CSS.
func applyTextStyle(elem map[string]interface{}, props map[string]string, scale float64) {
	if ff := props["font-family"]; ff != "" {
		// Clean up font family (remove quotes, fallbacks)
		elem["fontFamily"] = cleanFontFamily(ff)
	}
	if fw := props["font-weight"]; fw != "" {
		elem["fontStyle"] = cssFontWeight(fw)
	}
	if c := props["color"]; c != "" {
		hex := cssColorToHex(c)
		if hex != "" {
			elem["color"] = hex
		}
	}
	// Note: font-size is NOT scaled here — the composite expander
	// uses design tokens to set sizes based on element type + tier.
	// We only use CSS font-size to infer tier.
}

// inferHeadlineTier guesses the design token tier from CSS classes.
func inferHeadlineTier(classes string) string {
	if hasClass(classes, "h-xxl") || hasClass(classes, "h-display") {
		return "display"
	}
	if hasClass(classes, "h-xl") {
		return "hero"
	}
	if hasClass(classes, "h-lg") {
		return "title"
	}
	if hasClass(classes, "h-md") {
		return "heading"
	}
	if hasClass(classes, "h-sm") {
		return "subheading"
	}
	return "title" // default
}

// inferBannerHeadlineTier guesses tier for banner headlines.
func inferBannerHeadlineTier(classes string) string {
	if hasClass(classes, "eb-h-xxl") {
		return "title"
	}
	if hasClass(classes, "eb-h-xl") {
		return "heading"
	}
	if hasClass(classes, "eb-h-lg") {
		return "subheading"
	}
	return "heading" // default
}

// cleanFontFamily extracts the primary font name from a CSS font-family declaration.
func cleanFontFamily(ff string) string {
	// Take first font in the list
	parts := strings.SplitN(ff, ",", 2)
	name := strings.TrimSpace(parts[0])
	// Remove quotes
	name = strings.Trim(name, "'\"")
	return name
}

// getHTMLDimension reads a CSS dimension (width/height) from the element's classes.
func getHTMLDimension(classes string, styles map[string]map[string]string, prop string, fallback float64) float64 {
	merged := mergeClassStyles(classes, styles)
	if v, ok := merged[prop]; ok {
		px := remToPx(v, 16)
		if px > 0 {
			return px
		}
	}
	return fallback
}

// extractStatsItems extracts stat value/label pairs from a stats container.
func extractStatsItems(node *html.Node, styles map[string]map[string]string, scale float64) []interface{} {
	var items []interface{}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		classes := getAttr(c, "class")
		if hasClass(classes, "stat") || hasClass(classes, "s-stat") {
			value := ""
			label := ""
			for gc := c.FirstChild; gc != nil; gc = gc.NextSibling {
				if gc.Type != html.ElementNode {
					continue
				}
				gcClasses := getAttr(gc, "class")
				if hasClass(gcClasses, "stat-val") || hasClass(gcClasses, "s-stat-val") {
					value = textContent(gc)
				}
				if hasClass(gcClasses, "stat-lbl") || hasClass(gcClasses, "s-stat-lbl") {
					label = textContent(gc)
				}
			}
			if value != "" {
				items = append(items, map[string]interface{}{
					"value": value,
					"label": label,
				})
			}
		}
	}
	return items
}

// letterSpacingToFigma converts CSS letter-spacing (e.g. ".08em") to px.
func letterSpacingToFigma(val string, fontSize float64) float64 {
	val = strings.TrimSpace(val)
	lsRe := regexp.MustCompile(`^([.\d]+)\s*em$`)
	if m := lsRe.FindStringSubmatch(val); m != nil {
		v := parseFloatDef(m[1], 0)
		return v * fontSize
	}
	// Try px
	return remToPx(val, 16)
}
