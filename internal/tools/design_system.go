package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

func RegisterDesignSystemTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("design_system",
		mcp.WithDescription("Analyze the current Figma file's styles, variables, and components to generate design consistency rules. Use before creating designs in existing files."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("analyze")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		action := getStringArg(req.GetArguments(), "action", "")
		switch action {
		case "analyze":
			return analyzeDesignSystem(commander)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown design_system action: %s", action)), nil
		}
	})
}

func analyzeDesignSystem(commander *figma.Commander) (*mcp.CallToolResult, error) {
	stylesRaw, err := commander.SendCommand("get_all_styles", map[string]interface{}{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get styles: %s", err)), nil
	}

	varsRaw, err := commander.SendCommand("get_all_variables", map[string]interface{}{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get variables: %s", err)), nil
	}

	compsRaw, err := commander.SendCommand("get_local_components", map[string]interface{}{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get components: %s", err)), nil
	}

	styles := dsToMap(stylesRaw)
	vars := dsToMap(varsRaw)
	comps := dsToMap(compsRaw)

	result := buildDesignRules(styles, vars, comps)
	out, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(out)), nil
}

func buildDesignRules(styles, vars, comps map[string]interface{}) map[string]interface{} {
	rules := map[string]interface{}{}

	paintStyles := dsExtractArray(styles, "paintStyles")
	if len(paintStyles) > 0 {
		colorSection := map[string]interface{}{"styles": paintStyles}
		colorVars := dsFilterVarsByType(vars, "COLOR")
		if len(colorVars) > 0 {
			colorSection["variables"] = colorVars
		}
		colorSection["rule"] = "Use these existing colors. Do not introduce new hex values — apply paint styles by ID or reference color variables."
		rules["colors"] = colorSection
	} else if colorVars := dsFilterVarsByType(vars, "COLOR"); len(colorVars) > 0 {
		rules["colors"] = map[string]interface{}{
			"variables": colorVars,
			"rule":      "Use these existing color variables. Do not introduce new hex values.",
		}
	}

	textStyles := dsExtractArray(styles, "textStyles")
	if len(textStyles) > 0 {
		rules["typography"] = map[string]interface{}{
			"styles": textStyles,
			"rule":   "Apply text styles by ID when available. Match font families and weights from existing styles.",
		}
	}

	spacingVars := dsFilterVarsByType(vars, "FLOAT")
	if len(spacingVars) > 0 {
		rules["spacing"] = map[string]interface{}{
			"variables": spacingVars,
			"rule":      "Use spacing variables for padding, gaps, and margins. Snap to the existing scale.",
		}
	}

	effectStyles := dsExtractArray(styles, "effectStyles")
	if len(effectStyles) > 0 {
		rules["effects"] = map[string]interface{}{
			"styles": effectStyles,
			"rule":   "Apply effect styles by ID for shadows, blurs, and other effects instead of creating new ones.",
		}
	}

	components := dsExtractArray(comps, "components")
	if len(components) > 0 {
		rules["components"] = map[string]interface{}{
			"available": components,
			"rule":      "Instantiate existing components (component.create_instance) instead of rebuilding from scratch.",
		}
	}

	var summary []string
	if n := len(paintStyles); n > 0 {
		summary = append(summary, fmt.Sprintf("%d paint styles", n))
	}
	if n := len(textStyles); n > 0 {
		summary = append(summary, fmt.Sprintf("%d text styles", n))
	}
	if n := len(effectStyles); n > 0 {
		summary = append(summary, fmt.Sprintf("%d effect styles", n))
	}
	if allVars := dsExtractArray(vars, "variables"); len(allVars) > 0 {
		summary = append(summary, fmt.Sprintf("%d variables", len(allVars)))
	}
	if n := len(components); n > 0 {
		summary = append(summary, fmt.Sprintf("%d components", n))
	}

	if len(summary) > 0 {
		rules["summary"] = fmt.Sprintf("This file has %s. Follow the rules above for consistency.", strings.Join(summary, ", "))
	} else {
		rules["summary"] = "This file has no styles, variables, or components defined yet. You may create designs freely."
	}

	return rules
}

func dsToMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	b, err := json.Marshal(v)
	if err != nil {
		return map[string]interface{}{}
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return map[string]interface{}{}
	}
	return m
}

func dsExtractArray(m map[string]interface{}, key string) []interface{} {
	if arr, ok := m[key].([]interface{}); ok {
		return arr
	}
	return nil
}

func dsFilterVarsByType(vars map[string]interface{}, typeName string) []interface{} {
	allVars := dsExtractArray(vars, "variables")
	var result []interface{}
	for _, v := range allVars {
		if vm, ok := v.(map[string]interface{}); ok {
			if rt, _ := vm["resolvedType"].(string); rt == typeName {
				result = append(result, v)
			}
		}
	}
	return result
}
