package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterEffectTool registers the "effect" tool for shadows, blurs, and effect styles.
func RegisterEffectTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("effect",
		mcp.WithDescription("Effect operations: set effects, add shadow, add blur, apply effect style, remove effects, add_noise (Figma Beta), add_texture (Figma Beta), apply_glass (glass morphism — tries native GLASS first, falls back to simulated), add_glass (native GLASS effect only)."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("set_effects", "add_shadow", "add_blur", "apply_style", "remove", "add_noise", "add_texture", "apply_glass", "add_glass")),
		mcp.WithString("nodeId", mcp.Required(), mcp.Description("Target node ID")),
		mcp.WithString("effects", mcp.Description("JSON array of effect objects")),
		mcp.WithString("shadowType", mcp.Description("Shadow type: DROP_SHADOW, INNER_SHADOW"),
			mcp.Enum("DROP_SHADOW", "INNER_SHADOW")),
		mcp.WithString("color", mcp.Description("Shadow/noise color as hex string")),
		mcp.WithNumber("offsetX", mcp.Description("Shadow X offset")),
		mcp.WithNumber("offsetY", mcp.Description("Shadow Y offset")),
		mcp.WithNumber("radius", mcp.Description("Blur or shadow radius")),
		mcp.WithNumber("spread", mcp.Description("Shadow spread")),
		mcp.WithString("blurType", mcp.Description("Blur type: LAYER_BLUR, BACKGROUND_BLUR"),
			mcp.Enum("LAYER_BLUR", "BACKGROUND_BLUR")),
		mcp.WithString("styleId", mcp.Description("Effect style ID to apply")),
		mcp.WithString("noiseType", mcp.Description("Noise type: monotone, duotone, multitone")),
		mcp.WithNumber("noiseSize", mcp.Description("Noise grain size (default 100)")),
		mcp.WithNumber("density", mcp.Description("Noise density 0-1 (default 0.3)")),
		mcp.WithString("secondaryColor", mcp.Description("Secondary color for duotone noise")),
		mcp.WithString("intensity", mcp.Description("Glass intensity: light, medium, heavy")),
		mcp.WithString("tint", mcp.Description("Glass tint color (default #FFFFFF)")),
		mcp.WithNumber("lightIntensity", mcp.Description("Native glass: specular highlight intensity 0-1 (default 0.5)")),
		mcp.WithNumber("lightAngle", mcp.Description("Native glass: light angle in degrees (default 45)")),
		mcp.WithNumber("refraction", mcp.Description("Native glass: refraction distortion 0-1 (default 0.5)")),
		mcp.WithNumber("depth", mcp.Description("Native glass: refraction depth >=1 (default 1)")),
		mcp.WithNumber("dispersion", mcp.Description("Native glass: chromatic aberration 0-1 (default 0)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")
		nodeId := getStringArg(args, "nodeId", "")

		switch action {
		case "set_effects":
			return sendCommand(commander, "set_effects", map[string]interface{}{
				"nodeId":  nodeId,
				"effects": getStringArg(args, "effects", "[]"),
			})

		case "add_shadow":
			return sendCommand(commander, "add_shadow", map[string]interface{}{
				"nodeId":     nodeId,
				"shadowType": getStringArg(args, "shadowType", "DROP_SHADOW"),
				"color":      getStringArg(args, "color", "#00000040"),
				"offsetX":    getFloat64Arg(args, "offsetX", 0),
				"offsetY":    getFloat64Arg(args, "offsetY", 4),
				"radius":     getFloat64Arg(args, "radius", 4),
				"spread":     getFloat64Arg(args, "spread", 0),
			})

		case "add_blur":
			return sendCommand(commander, "add_blur", map[string]interface{}{
				"nodeId":   nodeId,
				"blurType": getStringArg(args, "blurType", "LAYER_BLUR"),
				"radius":   getFloat64Arg(args, "radius", 4),
			})

		case "apply_style":
			styleId, errResult := requireStringArg(args, "styleId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_effect_style_id", map[string]interface{}{
				"nodeId":  nodeId,
				"styleId": styleId,
			})

		case "remove":
			return sendCommand(commander, "set_effects", map[string]interface{}{
				"nodeId":  nodeId,
				"effects": "[]",
			})

		case "add_noise":
			params := map[string]interface{}{
				"nodeId":    nodeId,
				"noiseType": getStringArg(args, "noiseType", "monotone"),
				"color":     getStringArg(args, "color", "#FFFFFF"),
				"noiseSize": getFloat64Arg(args, "noiseSize", 100),
				"density":   getFloat64Arg(args, "density", 0.3),
			}
			if hasArg(args, "secondaryColor") {
				params["secondaryColor"] = getStringArg(args, "secondaryColor", "")
			}
			if hasArg(args, "blendMode") {
				params["blendMode"] = getStringArg(args, "blendMode", "SOFT_LIGHT")
			}
			return sendCommand(commander, "effect.add_noise", params)

		case "add_texture":
			return sendCommand(commander, "effect.add_texture", map[string]interface{}{
				"nodeId":    nodeId,
				"noiseSize": getFloat64Arg(args, "noiseSize", 100),
				"radius":    getFloat64Arg(args, "radius", 0),
			})

		case "apply_glass":
			params := map[string]interface{}{
				"nodeId":    nodeId,
				"intensity": getStringArg(args, "intensity", "medium"),
			}
			if hasArg(args, "tint") {
				params["tint"] = getStringArg(args, "tint", "#FFFFFF")
			}
			return sendCommand(commander, "effect.apply_glass", params)

		case "add_glass":
			params := map[string]interface{}{
				"nodeId":         nodeId,
				"lightIntensity": getFloat64Arg(args, "lightIntensity", 0.5),
				"lightAngle":     getFloat64Arg(args, "lightAngle", 45),
				"refraction":     getFloat64Arg(args, "refraction", 0.5),
				"depth":          getFloat64Arg(args, "depth", 1),
				"dispersion":     getFloat64Arg(args, "dispersion", 0),
				"radius":         getFloat64Arg(args, "radius", 0),
			}
			return sendCommand(commander, "effect.add_glass", params)

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown effect action: %s", action)), nil
		}
	})
}
