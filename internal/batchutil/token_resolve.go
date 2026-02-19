package batchutil

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

// ResolveTokenAliases resolves semantic token aliases (e.g. sz:"hero",
// padding:"side", w:"content") to numeric values based on the first root
// frame's canvas size. It mutates ops in place and returns:
//   - applied: number of alias replacements made
//   - rootWidth: detected root frame width (0 if none found)
func ResolveTokenAliases(ops []map[string]interface{}) (applied int, rootWidth int) {
	rootW, rootH, ok := detectRootCanvas(ops)
	if !ok {
		return 0, 0
	}

	aliases := buildTokenAliasTable(rootW, rootH)
	for _, op := range ops {
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		applied += resolveFieldAlias(params, "fontSize", aliases.text)
		applied += resolveFieldAlias(params, "sz", aliases.text)

		applied += resolveFieldAlias(params, "padding", aliases.padding)
		applied += resolveFieldAlias(params, "paddingTop", aliases.padding)
		applied += resolveFieldAlias(params, "paddingRight", aliases.padding)
		applied += resolveFieldAlias(params, "paddingBottom", aliases.padding)
		applied += resolveFieldAlias(params, "paddingLeft", aliases.padding)

		applied += resolveFieldAlias(params, "itemSpacing", aliases.spacing)
		applied += resolveFieldAlias(params, "gap", aliases.spacing)

		applied += resolveFieldAlias(params, "cornerRadius", aliases.radius)
		applied += resolveFieldAlias(params, "r", aliases.radius)

		applied += resolveFieldAlias(params, "width", aliases.width)
		applied += resolveFieldAlias(params, "w", aliases.width)
	}

	return applied, int(math.Round(rootW))
}

type tokenAliasTable struct {
	text    map[string]float64
	padding map[string]float64
	spacing map[string]float64
	radius  map[string]float64
	width   map[string]float64
}

func buildTokenAliasTable(rootW, rootH float64) tokenAliasTable {
	if rootH <= 0 {
		rootH = rootW
	}

	text := tokenSizes(rootW)
	text["numbers"] = text["display"]
	text["cta"] = text["body"]

	sidePadding := snap8(rootW*0.065, 16)
	contentWidth := math.Max(0, rootW-2*sidePadding)
	framePadding := snap8(rootW*0.06, 16)
	cardPadding := snap8(rootW*0.035, 12)

	itemSpacing := snap8(rootW*0.015, 8)
	cardGap := snap8(rootW*0.022, 8)
	sectionGap := snap8(rootH*0.035, 16)

	ctaHeight := snap8(text["cta"]*2.5, 40)
	buttonRadius := snap4(float64(ctaHeight)*0.28, 8)
	cardRadius := snap8(rootW*0.015, 8)

	return tokenAliasTable{
		text: text,
		padding: map[string]float64{
			"side":  sidePadding,
			"frame": framePadding,
			"card":  cardPadding,
		},
		spacing: map[string]float64{
			"section": sectionGap,
			"card":    cardGap,
			"item":    itemSpacing,
			"tight":   8,
		},
		radius: map[string]float64{
			"card":   cardRadius,
			"button": buttonRadius,
			"pill":   9999,
		},
		width: map[string]float64{
			"content": contentWidth,
		},
	}
}

func detectRootCanvas(ops []map[string]interface{}) (width float64, height float64, ok bool) {
	for _, op := range ops {
		cmd, _ := op["command"].(string)
		switch strings.ToLower(strings.TrimSpace(cmd)) {
		case "frame", "node.create_frame", "create_frame":
		default:
			continue
		}
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		if hasParentRef(params) {
			continue
		}
		w, hasW := numericParam(params, "width", "w")
		if !hasW || w <= 0 {
			continue
		}
		h, _ := numericParam(params, "height", "h")
		return w, h, true
	}
	return 0, 0, false
}

func hasParentRef(params map[string]interface{}) bool {
	if v, ok := params["parentId"]; ok {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return true
		}
	}
	if v, ok := params["pid"]; ok {
		if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
			return true
		}
	}
	return false
}

func numericParam(params map[string]interface{}, canonical, shorthand string) (float64, bool) {
	if v, ok := params[canonical]; ok {
		if f, ok := toFloat64Strict(v); ok {
			return f, true
		}
	}
	if shorthand != "" {
		if v, ok := params[shorthand]; ok {
			if f, ok := toFloat64Strict(v); ok {
				return f, true
			}
		}
	}
	return 0, false
}

func resolveFieldAlias(params map[string]interface{}, key string, table map[string]float64) int {
	raw, ok := params[key]
	if !ok {
		return 0
	}
	s, ok := raw.(string)
	if !ok {
		return 0
	}
	alias := strings.ToLower(strings.TrimSpace(s))
	if alias == "" {
		return 0
	}
	resolved, ok := table[alias]
	if !ok {
		return 0
	}
	params[key] = numericOrFloat(resolved)
	return 1
}

func numericOrFloat(v float64) interface{} {
	if math.Abs(v-math.Round(v)) < 0.000001 {
		return int(math.Round(v))
	}
	return v
}

func toFloat64Strict(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		return f, err == nil
	default:
		return 0, false
	}
}
