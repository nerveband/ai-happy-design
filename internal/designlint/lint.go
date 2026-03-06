package designlint

import (
	"fmt"
	"math"
	"strings"

	"github.com/nerveband/ai-happy-design/internal/tools"
)

// Issue represents a design lint warning.
type Issue struct {
	Step    int         `json:"step"`
	Name    string      `json:"name,omitempty"`
	Phase   string      `json:"phase"`
	Code    string      `json:"code"`
	Param   string      `json:"param,omitempty"`
	Message string      `json:"message"`
	Got     interface{} `json:"got,omitempty"`
	Fix     interface{} `json:"fix,omitempty"`
	Applied bool        `json:"applied"`
}

// Score holds per-axis quality scores (0-10).
type Score struct {
	Readability float64 `json:"readability"`
	Contrast    float64 `json:"contrast"`
	Spacing     float64 `json:"spacing"`
	Hierarchy   float64 `json:"hierarchy"`
	Overall     float64 `json:"overall"`
}

// Result holds lint results for a batch.
type Result struct {
	Canvas   map[string]interface{} `json:"canvas,omitempty"`
	Tokens   map[string]interface{} `json:"tokens,omitempty"`
	Warnings []Issue                `json:"warnings,omitempty"`
	Fixed    int                    `json:"fixed"`
	Score    Score                  `json:"score"`
}

// Check runs design lint on batch operations. Mutates ops for auto-fixes.
func Check(ops []map[string]interface{}) Result {
	var result Result

	// Detect canvas from first root frame
	canvasW, canvasH := detectCanvas(ops)
	// Always run structural and radius checks regardless of canvas
	checkRadiusOverflow(ops, &result)
	checkStructural(ops, &result)
	checkSizingTypes(ops, &result)

	if canvasW <= 0 || canvasH <= 0 {
		// No canvas detected — skip canvas-dependent checks
		computeScore(&result)
		return result
	}

	tokens := tools.ComputeDesignTokens(canvasW, canvasH, 0)
	result.Canvas = map[string]interface{}{
		"width": canvasW, "height": canvasH,
	}
	result.Tokens = tokens

	// Extract token values
	textTokens, _ := tokens["text"].(map[string]interface{})
	spacingTokens, _ := tokens["spacing"].(map[string]interface{})
	captionSize := getTokenFloat(textTokens, "caption")

	checkTextSizing(ops, &result, captionSize)
	checkContrast(ops, &result)
	checkSpacing(ops, &result, canvasW, spacingTokens)
	computeScore(&result)
	return result
}

func detectCanvas(ops []map[string]interface{}) (float64, float64) {
	for _, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		if strings.Contains(cmd, "frame") || cmd == "slide" || cmd == "banner" {
			_, hasParent := params["parentId"]
			if hasParent {
				continue
			}
			w, wOK := toFloat(params["width"])
			h, hOK := toFloat(params["height"])
			if wOK && hOK && w > 0 && h > 0 {
				return w, h
			}
		}
	}
	return 0, 0
}

func checkTextSizing(ops []map[string]interface{}, result *Result, captionSize float64) {
	if captionSize <= 0 {
		return
	}
	for i, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil || !strings.Contains(cmd, "text") {
			continue
		}
		fs, ok := toFloat(params["fontSize"])
		if !ok || fs <= 0 {
			continue
		}
		name, _ := op["name"].(string)
		if fs < captionSize {
			params["fontSize"] = captionSize
			result.Warnings = append(result.Warnings, Issue{
				Step: i, Name: name, Phase: "designLint", Code: "TEXT_TOO_SMALL",
				Param: "fontSize", Got: fs,
				Message: fmt.Sprintf("fontSize %.0f is below caption tier (%.0f) for this canvas", fs, captionSize),
				Fix:     captionSize,
				Applied: true,
			})
			result.Fixed++
		}
	}
}

func checkContrast(ops []map[string]interface{}, result *Result) {
	bgColors := map[int]string{}
	for i, op := range ops {
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		if color, ok := params["color"].(string); ok {
			if strings.HasPrefix(color, "#") {
				bgColors[i] = color
			}
		}
	}

	for i, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil || !strings.Contains(cmd, "text") {
			continue
		}
		textColor, ok := params["color"].(string)
		if !ok || !strings.HasPrefix(textColor, "#") {
			continue
		}

		bgColor := ""
		for j := i - 1; j >= 0; j-- {
			if c, exists := bgColors[j]; exists {
				bgColor = c
				break
			}
		}
		if bgColor == "" {
			continue
		}

		ratio := ContrastRatio(textColor, bgColor)
		name, _ := op["name"].(string)
		if ratio < 4.5 {
			fix := AdjustForContrast(textColor, bgColor, 4.5)
			params["color"] = fix
			result.Warnings = append(result.Warnings, Issue{
				Step: i, Name: name, Phase: "designLint", Code: "LOW_CONTRAST",
				Param: "color", Got: textColor,
				Message: fmt.Sprintf("text %s on background %s has contrast ratio %.1f:1 (minimum 4.5:1)", textColor, bgColor, ratio),
				Fix:     fix,
				Applied: true,
			})
			result.Fixed++
		}
	}
}

func checkSpacing(ops []map[string]interface{}, result *Result, canvasW float64, spacingTokens map[string]interface{}) {
	sidePadding := getTokenFloat(spacingTokens, "sidePadding")
	if sidePadding <= 0 {
		sidePadding = canvasW * 0.065
	}
	minPaddingRatio := 0.04

	for i, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		name, _ := op["name"].(string)

		if strings.Contains(cmd, "frame") {
			_, hasParent := params["parentId"]
			if !hasParent {
				for _, padKey := range []string{"padding", "paddingLeft", "paddingRight"} {
					pad, ok := toFloat(params[padKey])
					if ok && pad > 0 && pad/canvasW < minPaddingRatio {
						params[padKey] = sidePadding
						result.Warnings = append(result.Warnings, Issue{
							Step: i, Name: name, Phase: "designLint", Code: "PADDING_TOO_SMALL",
							Param: padKey, Got: pad,
							Message: fmt.Sprintf("%s %.0f is %.1f%% of canvas (minimum 4%%)", padKey, pad, pad/canvasW*100),
							Fix:     sidePadding,
							Applied: true,
						})
						result.Fixed++
					}
				}
			}
		}
	}
}

func checkRadiusOverflow(ops []map[string]interface{}, result *Result) {
	for i, op := range ops {
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		r, rOK := toFloat(params["cornerRadius"])
		w, wOK := toFloat(params["width"])
		h, hOK := toFloat(params["height"])
		if !rOK || !wOK || !hOK {
			continue
		}
		maxR := math.Min(w, h) / 2
		name, _ := op["name"].(string)
		if r > maxR {
			params["cornerRadius"] = maxR
			result.Warnings = append(result.Warnings, Issue{
				Step: i, Name: name, Phase: "designLint", Code: "RADIUS_OVERFLOW",
				Param: "cornerRadius", Got: r,
				Message: fmt.Sprintf("cornerRadius %.0f exceeds max %.0f (min(w,h)/2)", r, maxR),
				Fix:     maxR,
				Applied: true,
			})
			result.Fixed++
		}
	}
}

func checkStructural(ops []map[string]interface{}, result *Result) {
	names := map[string]int{}
	for i, op := range ops {
		name, _ := op["name"].(string)
		if name != "" {
			if prev, exists := names[name]; exists {
				result.Warnings = append(result.Warnings, Issue{
					Step: i, Phase: "designLint", Code: "DUPLICATE_STEP_NAME",
					Message: fmt.Sprintf("step name '%s' already used at step %d", name, prev),
					Got:     name,
				})
			}
			names[name] = i
		}
	}
}

// checkSizingTypes enforces sizing type constraints:
// - FILL requires parent to have auto-layout (layoutMode != NONE)
// - Auto-fixes FILL → HUG when parent is not auto-layout
func checkSizingTypes(ops []map[string]interface{}, result *Result) {
	// Build parent map: step name → layoutMode (from create_frame or set_auto_layout)
	parentLayout := map[string]string{} // step name → layoutMode
	for _, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		name, _ := op["name"].(string)

		if cmd == "node.create_frame" || cmd == "frame" {
			mode, _ := params["layoutMode"].(string)
			if name != "" && mode != "" {
				parentLayout[name] = mode
			}
		}
		if cmd == "layout.set_auto_layout" || cmd == "layout.auto_layout" {
			// Find the step name for this nodeId
			nodeId, _ := params["nodeId"].(string)
			dir, _ := params["direction"].(string)
			if dir == "" {
				dir, _ = params["layoutMode"].(string)
			}
			if dir != "" && dir != "NONE" {
				// Map from interpolation ref to step name
				targetName := extractStepNameFromRef(nodeId)
				if targetName != "" {
					parentLayout[targetName] = dir
				}
			}
			if name != "" && dir != "" && dir != "NONE" {
				parentLayout[name] = dir
			}
		}
	}

	// Check sizing operations for FILL on non-auto-layout parents
	for i, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		name, _ := op["name"].(string)

		if cmd != "layout.set_sizing" && cmd != "layout.sizing" {
			continue
		}

		nodeRef, _ := params["nodeId"].(string)
		targetName := extractStepNameFromRef(nodeRef)
		h, _ := params["horizontal"].(string)
		v, _ := params["vertical"].(string)

		// Check if target's parent has auto-layout
		// Look up what step created the target and find its parentId
		parentName := findParentStepName(ops, targetName)
		parentMode := parentLayout[parentName]

		if parentMode == "" || parentMode == "NONE" {
			if h == "FILL" {
				result.Warnings = append(result.Warnings, Issue{
					Step: i, Name: name, Phase: "designLint", Code: "FILL_NO_AUTOLAYOUT",
					Message: fmt.Sprintf("FILL sizing requires auto-layout parent. Parent '%s' has no auto-layout. Auto-fixed to HUG.", parentName),
					Got: "FILL", Fix: "HUG", Applied: true,
				})
				params["horizontal"] = "HUG"
				result.Fixed++
			}
			if v == "FILL" {
				result.Warnings = append(result.Warnings, Issue{
					Step: i, Name: name, Phase: "designLint", Code: "FILL_NO_AUTOLAYOUT",
					Message: fmt.Sprintf("FILL vertical sizing requires auto-layout parent. Parent '%s' has no auto-layout. Auto-fixed to HUG.", parentName),
					Got: "FILL", Fix: "HUG", Applied: true,
				})
				params["vertical"] = "HUG"
				result.Fixed++
			}
		}
	}
}

// extractStepNameFromRef extracts step name from interpolation like "${{steps.foo.result.id}}"
func extractStepNameFromRef(ref string) string {
	if !strings.HasPrefix(ref, "${{steps.") {
		return ""
	}
	// ${{steps.NAME.result.id}}
	trimmed := strings.TrimPrefix(ref, "${{steps.")
	idx := strings.Index(trimmed, ".")
	if idx < 0 {
		return ""
	}
	return trimmed[:idx]
}

// findParentStepName finds the step name of the parent frame for a given step
func findParentStepName(ops []map[string]interface{}, targetName string) string {
	for _, op := range ops {
		name, _ := op["name"].(string)
		if name != targetName {
			continue
		}
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			return ""
		}
		parentRef, _ := params["parentId"].(string)
		return extractStepNameFromRef(parentRef)
	}
	return ""
}

func computeScore(result *Result) {
	scores := map[string]float64{
		"readability": 10, "contrast": 10, "spacing": 10, "hierarchy": 10,
	}
	for _, w := range result.Warnings {
		switch w.Code {
		case "TEXT_TOO_SMALL":
			scores["readability"] -= 2
		case "TEXT_NO_HIERARCHY":
			scores["hierarchy"] -= 3
		case "LOW_CONTRAST":
			scores["contrast"] -= 3
		case "GRADIENT_ALPHA":
			scores["contrast"] -= 1
		case "PADDING_TOO_SMALL":
			scores["spacing"] -= 2
		case "SPACING_EXTREME":
			scores["spacing"] -= 2
		case "RADIUS_OVERFLOW":
			scores["spacing"] -= 1
		case "ELEMENT_DENSITY":
			scores["spacing"] -= 1
		case "FILL_NO_AUTOLAYOUT":
			scores["spacing"] -= 2
		}
	}
	for k, v := range scores {
		if v < 0 {
			scores[k] = 0
		}
	}
	result.Score = Score{
		Readability: scores["readability"],
		Contrast:    scores["contrast"],
		Spacing:     scores["spacing"],
		Hierarchy:   scores["hierarchy"],
		Overall:     scores["readability"]*0.3 + scores["contrast"]*0.25 + scores["spacing"]*0.25 + scores["hierarchy"]*0.2,
	}
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

func getTokenFloat(tokens map[string]interface{}, key string) float64 {
	if tokens == nil {
		return 0
	}
	v, ok := toFloat(tokens[key])
	if !ok {
		return 0
	}
	return v
}
