package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func compareCodeSpec(params map[string]interface{}) (interface{}, error) {
	specPath, _ := params["specPath"].(string)
	if specPath == "" {
		specPath, _ = params["codeSpecPath"].(string)
	}
	if specPath == "" {
		specPath, _ = params["file"].(string)
	}
	var raw []byte
	if inline, ok := params["codeSpec"]; ok {
		var err error
		raw, err = json.Marshal(inline)
		if err != nil {
			return nil, err
		}
	}
	if specPath == "" && len(raw) == 0 {
		return nil, fmt.Errorf("specPath is required")
	}
	if len(raw) == 0 {
		data, err := os.ReadFile(specPath)
		if err != nil {
			return nil, err
		}
		raw = data
	}
	var spec struct {
		Colors        map[string]interface{} `json:"colors"`
		Typography    map[string]interface{} `json:"typography"`
		Spacing       map[string]interface{} `json:"spacing"`
		Radii         map[string]interface{} `json:"radii"`
		Shadows       map[string]interface{} `json:"shadows"`
		Effects       map[string]interface{} `json:"effects"`
		Opacity       map[string]interface{} `json:"opacity"`
		Sizing        map[string]interface{} `json:"sizing"`
		Accessibility map[string]interface{} `json:"accessibility"`
		Nodes         []struct {
			Name       string                 `json:"name"`
			Type       string                 `json:"type"`
			Properties map[string]interface{} `json:"properties"`
		} `json:"nodes"`
	}
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, fmt.Errorf("invalid code spec JSON: %w", err)
	}
	findings := make([]map[string]interface{}, 0)
	for i, node := range spec.Nodes {
		if strings.TrimSpace(node.Name) == "" {
			findings = append(findings, map[string]interface{}{"index": i, "code": "MISSING_NAME", "message": "node spec is missing a semantic name"})
		}
		if strings.TrimSpace(node.Type) == "" {
			findings = append(findings, map[string]interface{}{"index": i, "code": "MISSING_TYPE", "message": "node spec is missing a type"})
		}
		for _, key := range []string{"colors", "typography", "spacing", "radii", "shadows", "effects", "opacity", "sizing", "accessibility"} {
			if node.Properties != nil {
				if _, ok := node.Properties[key]; ok {
					continue
				}
			}
			findings = append(findings, map[string]interface{}{"index": i, "code": "NODE_MISSING_" + strings.ToUpper(key), "message": "node spec does not declare " + key + " parity data"})
		}
	}
	categories := map[string]map[string]interface{}{
		"colors": spec.Colors, "typography": spec.Typography, "spacing": spec.Spacing, "radii": spec.Radii,
		"shadows": spec.Shadows, "effects": spec.Effects, "opacity": spec.Opacity, "sizing": spec.Sizing, "accessibility": spec.Accessibility,
	}
	for name, values := range categories {
		if len(values) == 0 {
			findings = append(findings, map[string]interface{}{"code": "MISSING_" + strings.ToUpper(name), "message": "code spec does not declare " + name})
		}
	}
	score := 100 - len(findings)*10
	if score < 0 {
		score = 0
	}
	threshold, _ := params["threshold"].(float64)
	if threshold <= 0 {
		threshold = 80
	}
	return map[string]interface{}{
		"ok":         score >= int(threshold),
		"score":      score,
		"threshold":  threshold,
		"checked":    len(spec.Nodes),
		"categories": []string{"colors", "typography", "spacing", "radii", "shadows", "effects", "opacity", "sizing", "accessibility"},
		"findings":   findings,
	}, nil
}

func exportTokenPreset(params map[string]interface{}) (interface{}, error) {
	config, err := loadTokenExportConfig(params)
	if err != nil {
		return nil, err
	}
	mergeTokenConfig(params, config)
	preset, _ := params["preset"].(string)
	if preset == "" {
		preset = "minimal"
	}
	format, _ := params["format"].(string)
	if format == "" {
		format = "json"
	}
	tokens := minimalTokenPreset()
	if variableTokens, ok, err := tokensFromVariables(params); err != nil {
		return nil, err
	} else if ok {
		tokens = variableTokens
		preset = "figma_variables"
	}
	rendered := renderTokens(tokens, format)
	out := map[string]interface{}{
		"preset":  preset,
		"format":  format,
		"tokens":  tokens,
		"content": rendered,
	}
	saved := map[string]string{}
	if path, _ := params["outputPath"].(string); path != "" {
		if err := writeTokenOutput(path, format, rendered, out); err != nil {
			return nil, err
		}
		saved[format] = path
	}
	if outputs, ok := params["outputs"].(map[string]interface{}); ok {
		for outputFormat, rawPath := range outputs {
			path, _ := rawPath.(string)
			if path == "" {
				continue
			}
			content := renderTokens(tokens, outputFormat)
			if err := writeTokenOutput(path, outputFormat, content, map[string]interface{}{"preset": preset, "format": outputFormat, "tokens": tokens, "content": content}); err != nil {
				return nil, err
			}
			saved[outputFormat] = path
		}
	}
	if len(saved) > 0 {
		out["saved"] = saved
	}
	return out, nil
}

func loadTokenExportConfig(params map[string]interface{}) (map[string]interface{}, error) {
	path, _ := params["configPath"].(string)
	if path == "" {
		path, _ = params["config"].(string)
	}
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("token export config must be JSON: %w", err)
	}
	return config, nil
}

func mergeTokenConfig(params, config map[string]interface{}) {
	for key, value := range config {
		if _, exists := params[key]; !exists {
			params[key] = value
		}
	}
}

func writeTokenOutput(path, format, rendered string, envelope map[string]interface{}) error {
	data := []byte(rendered)
	if format == "json" || rendered == "" {
		encoded, err := json.MarshalIndent(envelope, "", "  ")
		if err != nil {
			return err
		}
		data = encoded
	}
	return os.WriteFile(path, data, 0644)
}

func tokensFromVariables(params map[string]interface{}) (map[string]interface{}, bool, error) {
	var raw interface{}
	if variables, ok := params["variables"]; ok {
		raw = variables
	}
	if path, _ := params["variablesFile"].(string); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, false, err
		}
		var parsed interface{}
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, false, err
		}
		raw = parsed
	}
	if raw == nil {
		return nil, false, nil
	}
	var payload interface{} = raw
	if m, ok := raw.(map[string]interface{}); ok {
		if vars, ok := m["variables"]; ok {
			payload = vars
		}
		if result, ok := m["result"].(map[string]interface{}); ok {
			if vars, ok := result["variables"]; ok {
				payload = vars
			}
		}
	}
	var variables []interface{}
	switch typed := payload.(type) {
	case []interface{}:
		variables = typed
	default:
		return nil, false, fmt.Errorf("variables must be an array or variable.get_all response")
	}
	tokens := map[string]interface{}{
		"color":  map[string]string{},
		"space":  map[string]float64{},
		"radius": map[string]float64{},
		"type":   map[string]float64{},
		"raw":    map[string]interface{}{},
	}
	for _, item := range variables {
		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := v["name"].(string)
		if name == "" {
			continue
		}
		resolvedType, _ := v["resolvedType"].(string)
		value := firstVariableValue(v["valuesByMode"])
		normalized := tokenKey(name)
		tokens["raw"].(map[string]interface{})[normalized] = value
		switch resolvedType {
		case "COLOR":
			if hex, ok := figmaColorToHex(value); ok {
				tokens["color"].(map[string]string)[normalized] = hex
			}
		case "FLOAT":
			if n, ok := numericValue(value); ok {
				group := "space"
				if strings.Contains(strings.ToLower(name), "radius") {
					group = "radius"
				}
				if strings.Contains(strings.ToLower(name), "font") || strings.Contains(strings.ToLower(name), "type") {
					group = "type"
				}
				tokens[group].(map[string]float64)[normalized] = n
			}
		}
	}
	return tokens, true, nil
}

func firstVariableValue(raw interface{}) interface{} {
	if values, ok := raw.(map[string]interface{}); ok {
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if len(keys) > 0 {
			return values[keys[0]]
		}
	}
	return raw
}

func tokenKey(name string) string {
	name = strings.Trim(strings.ToLower(name), "/ ")
	name = strings.NewReplacer("/", "-", " ", "-", "_", "-").Replace(name)
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	return name
}

func figmaColorToHex(value interface{}) (string, bool) {
	m, ok := value.(map[string]interface{})
	if !ok {
		return "", false
	}
	r, rok := numericValue(m["r"])
	g, gok := numericValue(m["g"])
	b, bok := numericValue(m["b"])
	a, aok := numericValue(m["a"])
	if !rok || !gok || !bok {
		return "", false
	}
	if !aok {
		a = 1
	}
	if a < 1 {
		return fmt.Sprintf("#%02X%02X%02X%02X", colorByte(r), colorByte(g), colorByte(b), colorByte(a)), true
	}
	return fmt.Sprintf("#%02X%02X%02X", colorByte(r), colorByte(g), colorByte(b)), true
}

func colorByte(v float64) int {
	if v <= 1 {
		v *= 255
	}
	return int(math.Round(math.Max(0, math.Min(255, v))))
}

func minimalTokenPreset() map[string]interface{} {
	return map[string]interface{}{
		"color": map[string]string{
			"background": "#FFFFFF",
			"foreground": "#111827",
			"muted":      "#6B7280",
			"primary":    "#2563EB",
			"accent":     "#10B981",
			"danger":     "#DC2626",
		},
		"radius": map[string]float64{"sm": 4, "md": 8, "lg": 16, "xl": 24},
		"space":  map[string]float64{"xs": 4, "sm": 8, "md": 16, "lg": 24, "xl": 40},
		"type":   map[string]float64{"caption": 12, "body": 16, "heading": 24, "title": 36, "hero": 56},
	}
}

func renderTokens(tokens map[string]interface{}, format string) string {
	switch format {
	case "css":
		return renderCSSTokens(tokens)
	case "tailwind":
		data, _ := json.MarshalIndent(map[string]interface{}{"theme": map[string]interface{}{"extend": tokens}}, "", "  ")
		return string(data)
	case "swift":
		return "enum FigmaTokens {\n" +
			"  static let color = " + quotedJSON(tokens["color"]) + "\n" +
			"  static let radius = " + quotedJSON(tokens["radius"]) + "\n" +
			"  static let space = " + quotedJSON(tokens["space"]) + "\n" +
			"  static let type = " + quotedJSON(tokens["type"]) + "\n" +
			"}\n"
	case "android":
		return renderAndroidTokens(tokens)
	default:
		data, _ := json.MarshalIndent(tokens, "", "  ")
		return string(data)
	}
}

func renderCSSTokens(tokens map[string]interface{}) string {
	var b strings.Builder
	b.WriteString(":root {\n")
	for group, raw := range tokens {
		if values, ok := raw.(map[string]string); ok {
			for name, value := range values {
				b.WriteString(fmt.Sprintf("  --%s-%s: %s;\n", group, name, value))
			}
		}
		if values, ok := raw.(map[string]float64); ok {
			for name, value := range values {
				b.WriteString(fmt.Sprintf("  --%s-%s: %.0fpx;\n", group, name, value))
			}
		}
	}
	b.WriteString("}\n")
	return b.String()
}

func renderAndroidTokens(tokens map[string]interface{}) string {
	var b strings.Builder
	b.WriteString("<resources>\n")
	if colors, ok := tokens["color"].(map[string]string); ok {
		for name, value := range colors {
			b.WriteString(fmt.Sprintf("  <color name=\"figma_%s\">%s</color>\n", name, value))
		}
	}
	b.WriteString("</resources>\n")
	return b.String()
}

func quotedJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return strconv.Quote(string(data))
}

func auditBatchAccessibility(params map[string]interface{}) (interface{}, error) {
	path, _ := params["file"].(string)
	if path == "" {
		path, _ = params["batchFile"].(string)
	}
	if path == "" {
		return nil, fmt.Errorf("file is required for local accessibility audit")
	}
	ops, err := loadBatchOperations("", path)
	if err != nil {
		return nil, err
	}
	findings := make([]map[string]interface{}, 0)
	names := map[string]int{}
	for i, op := range ops {
		name, _ := op.Params["name"].(string)
		if strings.TrimSpace(name) == "" && createsVisibleNode(op.Command) {
			findings = append(findings, finding(i, op.Command, "WCAG_IMAGE_ALT", "created visible nodes should have semantic names/descriptions for assistive review", "warning"))
		}
		if name != "" {
			names[name]++
		}
		if size, ok := numericParam(op.Params, "fontSize"); ok && size < 12 {
			findings = append(findings, finding(i, op.Command, "WCAG_TEXT_CONTRAST", "fontSize below 12px hurts readability and contrast perception", "error"))
		}
		if strings.HasPrefix(op.Command, "text.") {
			if color, _ := op.Params["color"].(string); color != "" {
				if contrast, ok := textContrastAgainstBackground(color, op.Params); ok && contrast < 4.5 {
					findings = append(findings, finding(i, op.Command, "WCAG_TEXT_CONTRAST", fmt.Sprintf("text contrast %.2f is below WCAG AA 4.5:1", contrast), "error"))
				}
			}
		}
		if lineHeight, ok := numericParam(op.Params, "lineHeight"); ok && lineHeight > 0 && lineHeight < 120 {
			findings = append(findings, finding(i, op.Command, "WCAG_LINE_HEIGHT", "lineHeight below 120% can reduce readability", "warning"))
		}
		if createsVisibleNode(op.Command) {
			if opacity, ok := numericParam(op.Params, "opacity"); ok && opacity < 0.35 {
				findings = append(findings, finding(i, op.Command, "WCAG_NON_TEXT_CONTRAST", "very low opacity can make non-text content fail contrast requirements", "warning"))
			}
		}
		if w, ok := numericParam(op.Params, "width"); ok {
			if h, ok := numericParam(op.Params, "height"); ok && (w < 24 || h < 24) && strings.Contains(op.Command, "button") {
				findings = append(findings, finding(i, op.Command, "WCAG_TARGET_SIZE", "interactive target should be at least 24x24", "warning"))
			}
		}
		if strings.Contains(strings.ToLower(name), "error") || strings.Contains(strings.ToLower(name), "success") {
			findings = append(findings, finding(i, op.Command, "WCAG_COLOR_ONLY", "status-like layer names should not rely on color alone; include text or icon semantics", "warning"))
		}
	}
	duplicates := make([]string, 0)
	for name, count := range names {
		if count > 1 {
			duplicates = append(duplicates, name)
		}
	}
	sort.Strings(duplicates)
	for _, name := range duplicates {
		findings = append(findings, map[string]interface{}{"code": "WCAG_FOCUS_INDICATOR", "name": name, "severity": "warning", "message": "duplicate semantic layer names make focus and review targets ambiguous"})
	}
	codes := []string{"WCAG_TEXT_CONTRAST", "WCAG_NON_TEXT_CONTRAST", "WCAG_TARGET_SIZE", "WCAG_LINE_HEIGHT", "WCAG_COLOR_ONLY", "WCAG_FOCUS_INDICATOR", "WCAG_IMAGE_ALT"}
	score := 100 - len(findings)*5
	if score < 0 {
		score = 0
	}
	return map[string]interface{}{
		"ok":         len(findings) == 0,
		"score":      score,
		"operations": len(ops),
		"codes":      codes,
		"findings":   findings,
	}, nil
}

func textContrastAgainstBackground(color string, params map[string]interface{}) (float64, bool) {
	bg, _ := params["backgroundColor"].(string)
	if bg == "" {
		bg, _ = params["bg"].(string)
	}
	if bg == "" {
		return 0, false
	}
	fgRGB, ok := parseHexRGB(color)
	if !ok {
		return 0, false
	}
	bgRGB, ok := parseHexRGB(bg)
	if !ok {
		return 0, false
	}
	l1 := relativeLuminance(fgRGB)
	l2 := relativeLuminance(bgRGB)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05), true
}

func parseHexRGB(hex string) ([3]float64, bool) {
	var out [3]float64
	raw := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(raw) == 3 {
		raw = string([]byte{raw[0], raw[0], raw[1], raw[1], raw[2], raw[2]})
	}
	if len(raw) < 6 {
		return out, false
	}
	n, err := strconv.ParseUint(raw[:6], 16, 32)
	if err != nil {
		return out, false
	}
	out[0] = float64((n>>16)&0xff) / 255
	out[1] = float64((n>>8)&0xff) / 255
	out[2] = float64(n&0xff) / 255
	return out, true
}

func relativeLuminance(rgb [3]float64) float64 {
	var c [3]float64
	for i, v := range rgb {
		if v <= 0.03928 {
			c[i] = v / 12.92
		} else {
			c[i] = math.Pow((v+0.055)/1.055, 2.4)
		}
	}
	return 0.2126*c[0] + 0.7152*c[1] + 0.0722*c[2]
}

func createsVisibleNode(command string) bool {
	return strings.HasPrefix(command, "node.create_") || strings.HasPrefix(command, "shape.create_") || strings.HasPrefix(command, "text.create")
}

func numericParam(params map[string]interface{}, key string) (float64, bool) {
	v, ok := params[key]
	if !ok {
		return 0, false
	}
	return numericValue(v)
}

func numericValue(v interface{}) (float64, bool) {
	switch typed := v.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case json.Number:
		f, err := typed.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func finding(index int, command, code, message, severity string) map[string]interface{} {
	return map[string]interface{}{"index": index, "command": command, "code": code, "message": message, "severity": severity}
}
