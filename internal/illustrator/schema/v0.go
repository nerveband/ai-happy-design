package schema

import "github.com/nerveband/ai-happy-design/internal/commonschema"

const (
	assetPattern = `^[A-Za-z0-9._:/ -]+$`
	colorPattern = `^#(?:[0-9A-Fa-f]{6}|[0-9A-Fa-f]{8})$`
)

func init() {
	registerApp()
	registerDocument()
	registerArtboard()
	registerLayer()
	registerSelection()
	registerPath()
	registerText()
	registerAppearance()
	registerAction()
	registerExport()
	registerInspect()
}

func registerApp() {
	commonschema.Register(commonschema.Command{
		Name:        "app.info",
		Aliases:     []string{"app.get_info"},
		Domain:      "app",
		Description: "Return application metadata and runtime host info.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.version",
		Domain:      "app",
		Description: "Return the Illustrator version string.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.select_tool",
		Domain:      "app",
		Description: "Select an Illustrator tool by name.",
		Mutating:    true,
		Params: []commonschema.Param{
			enumParam("tool", "Tool name to activate.", true, []string{"selectionTool", "penTool", "typeTool", "rectangleTool", "ellipseTool"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.execute_menu",
		Domain:      "app",
		Description: "Execute an Illustrator menu command string.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("menuCommand", "Menu command string.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.user_interaction_level",
		Domain:      "app",
		Description: "Set or inspect the user interaction level.",
		Mutating:    true,
		Params: []commonschema.Param{
			enumParam("mode", "Interaction level.", false, []string{"DISPLAYALERTS", "DONTDISPLAYALERTS"}, true),
		},
	})
}

func registerDocument() {
	commonschema.Register(commonschema.Command{
		Name:        "document.new",
		Domain:      "document",
		Description: "Create a new document.",
		Mutating:    true,
		Params: []commonschema.Param{
			numberParam("width", "Document width in points.", false, commonschema.Ptr(1), nil),
			numberParam("height", "Document height in points.", false, commonschema.Ptr(1), nil),
			numberParam("artboards", "Number of artboards.", false, commonschema.Ptr(1), commonschema.Ptr(100)),
			stringParam("colorSpace", "Document color space.", false),
			stringParam("preset", "Document preset name.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.open",
		Domain:      "document",
		Description: "Open an Illustrator document from disk.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Path to the .ai file to open.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.save",
		Domain:      "document",
		Description: "Save the current document.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("documentId", "Optional document identifier.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.save_as",
		Domain:      "document",
		Description: "Save the current document to a new path.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Destination path.", true),
			enumParam("format", "Output format.", false, []string{"ai", "pdf"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.close",
		Domain:      "document",
		Description: "Close the current or specified document.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("documentId", "Optional document identifier.", false),
			boolParam("save", "Save before closing.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.list",
		Aliases:     []string{"document.get_all"},
		Domain:      "document",
		Description: "List open documents.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.info",
		Domain:      "document",
		Description: "Return the active document summary.",
		Params: []commonschema.Param{
			stringParam("documentId", "Optional document identifier.", false),
		},
	})
}

func registerArtboard() {
	commonschema.Register(commonschema.Command{Name: "artboard.list", Domain: "artboard", Description: "List artboards in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "artboard.create",
		Domain:      "artboard",
		Description: "Create a new artboard.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Artboard name.", true),
			numberParam("left", "Left coordinate.", true, nil, nil),
			numberParam("top", "Top coordinate.", true, nil, nil),
			numberParam("right", "Right coordinate.", true, nil, nil),
			numberParam("bottom", "Bottom coordinate.", true, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "artboard.resize",
		Domain:      "artboard",
		Description: "Resize an artboard.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("artboardId", "Artboard identifier.", true),
			numberParam("width", "Artboard width.", true, commonschema.Ptr(1), nil),
			numberParam("height", "Artboard height.", true, commonschema.Ptr(1), nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "artboard.set_active",
		Domain:      "artboard",
		Description: "Set the active artboard.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("artboardId", "Artboard identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "artboard.fit_to_artwork",
		Domain:      "artboard",
		Description: "Fit the active artboard to its artwork bounds.",
		Mutating:    true,
		Params:      []commonschema.Param{},
	})
}

func registerLayer() {
	commonschema.Register(commonschema.Command{Name: "layer.list", Domain: "layer", Description: "List layers in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "layer.create",
		Domain:      "layer",
		Description: "Create a new layer.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Layer name.", true),
			stringParam("parentLayerId", "Optional parent layer identifier.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "layer.rename",
		Domain:      "layer",
		Description: "Rename a layer.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("layerId", "Layer identifier.", true),
			stringParam("name", "New layer name.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "layer.visibility",
		Domain:      "layer",
		Description: "Set layer visibility.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("layerId", "Layer identifier.", true),
			boolParam("visible", "Whether the layer is visible.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "layer.lock",
		Domain:      "layer",
		Description: "Set layer lock state.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("layerId", "Layer identifier.", true),
			boolParam("locked", "Whether the layer is locked.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "layer.reorder",
		Domain:      "layer",
		Description: "Move a layer before or after another layer.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("layerId", "Layer identifier.", true),
			stringParam("relativeTo", "Neighbor layer identifier.", true),
			enumParam("position", "Reorder position.", true, []string{"before", "after", "inside"}, true),
		},
	})
}

func registerSelection() {
	commonschema.Register(commonschema.Command{Name: "selection.get", Domain: "selection", Description: "Return the current selection.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{Name: "selection.clear", Domain: "selection", Description: "Clear the current selection.", Mutating: true, Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "selection.set_by_ids",
		Domain:      "selection",
		Description: "Select items by stable IDs.",
		Mutating:    true,
		Params: []commonschema.Param{
			arrayParam("ids", "Item identifiers to select.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "selection.select_by_name",
		Domain:      "selection",
		Description: "Select items by name.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Name to match.", true),
			boolParam("partial", "Allow partial name matches.", false),
		},
	})
}

func registerPath() {
	commonschema.Register(commonschema.Command{
		Name:        "path.create_rect",
		Domain:      "path",
		Description: "Create a rectangle path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Path name.", true),
			numberParam("left", "Left coordinate.", true, nil, nil),
			numberParam("top", "Top coordinate.", true, nil, nil),
			numberParam("width", "Width.", true, commonschema.Ptr(1), nil),
			numberParam("height", "Height.", true, commonschema.Ptr(1), nil),
			numberParam("cornerRadius", "Corner radius.", false, commonschema.Ptr(0), nil),
			stringParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "path.create_ellipse",
		Domain:      "path",
		Description: "Create an ellipse path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Path name.", true),
			numberParam("left", "Left coordinate.", true, nil, nil),
			numberParam("top", "Top coordinate.", true, nil, nil),
			numberParam("width", "Width.", true, commonschema.Ptr(1), nil),
			numberParam("height", "Height.", true, commonschema.Ptr(1), nil),
			stringParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "path.create_path",
		Domain:      "path",
		Description: "Create a path item from a set of points.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Path name.", true),
			arrayParam("points", "Point array in Illustrator coordinate space.", true),
			boolParam("closed", "Whether the path is closed.", false),
			stringParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "path.transform",
		Domain:      "path",
		Description: "Transform a path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Path item identifier.", true),
			numberParam("translateX", "Horizontal translation.", false, nil, nil),
			numberParam("translateY", "Vertical translation.", false, nil, nil),
			numberParam("scaleX", "Horizontal scale percentage.", false, commonschema.Ptr(0), nil),
			numberParam("scaleY", "Vertical scale percentage.", false, commonschema.Ptr(0), nil),
			numberParam("rotate", "Rotation in degrees.", false, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "path.duplicate",
		Domain:      "path",
		Description: "Duplicate a path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Path item identifier.", true),
			stringParam("destinationLayerId", "Optional destination layer identifier.", false),
		},
	})
}

func registerText() {
	commonschema.Register(commonschema.Command{
		Name:        "text.create",
		Domain:      "text",
		Description: "Create a point text item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Text item name.", true),
			{Name: "contents", Type: "string", Description: "Text contents.", Required: true, Aliases: []string{"text"}},
			numberParam("left", "Left coordinate.", true, nil, nil),
			numberParam("top", "Top coordinate.", true, nil, nil),
			stringParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "text.set_contents",
		Domain:      "text",
		Description: "Replace the contents of a text frame.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Text frame identifier.", true),
			stringParam("contents", "New contents.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "text.set_style",
		Domain:      "text",
		Description: "Apply basic typographic styling to a text frame.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Text frame identifier.", true),
			stringParam("fontFamily", "Font family name.", false),
			numberParam("fontSize", "Font size in points.", false, commonschema.Ptr(1), nil),
			numberParam("tracking", "Tracking value.", false, nil, nil),
			numberParam("leading", "Leading value.", false, nil, nil),
			stringParam("fillColor", "Hex fill color.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:           "text.outline",
		Domain:         "text",
		Description:    "Create outlines from a text frame.",
		Mutating:       true,
		PluginRequired: false,
		Params: []commonschema.Param{
			stringParam("itemId", "Text frame identifier.", true),
		},
	})
}

func registerAppearance() {
	commonschema.Register(commonschema.Command{
		Name:        "appearance.set_fill",
		Domain:      "appearance",
		Description: "Set a solid fill on an item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Target item identifier.", true),
			colorParam("color", "Solid fill color.", true),
			numberParam("opacity", "Fill opacity percent.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "appearance.set_stroke",
		Domain:      "appearance",
		Description: "Set stroke color and weight on an item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Target item identifier.", true),
			colorParam("color", "Stroke color.", true),
			numberParam("width", "Stroke width.", true, commonschema.Ptr(0), nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "appearance.set_gradient",
		Domain:      "appearance",
		Description: "Set a simple linear gradient on an item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Target item identifier.", true),
			arrayParam("stops", "Gradient stop definitions.", true),
			enumParam("type", "Gradient type.", false, []string{"linear", "radial"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "appearance.apply_graphic_style",
		Domain:      "appearance",
		Description: "Apply a named graphic style to an item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("itemId", "Target item identifier.", true),
			stringParam("styleName", "Graphic style name.", true),
		},
	})
}

func registerAction() {
	commonschema.Register(commonschema.Command{
		Name:        "action.load",
		Domain:      "action",
		Description: "Load an Illustrator action set from disk.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Path to the .aia file.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "action.run",
		Domain:      "action",
		Description: "Run an action by set and action name.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("setName", "Action set name.", true),
			stringParam("actionName", "Action name.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "action.unload",
		Domain:      "action",
		Description: "Unload an Illustrator action set.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("setName", "Action set name.", true),
		},
	})
}

func registerExport() {
	commonschema.Register(commonschema.Command{
		Name:        "export.png",
		Domain:      "export",
		Description: "Export the document or selection as PNG.",
		Mutating:    true,
		Params: []commonschema.Param{
			{Name: "outputPath", Type: "string", Description: "Destination path.", Required: true, SafePath: true, Aliases: []string{"path"}},
			numberParam("scale", "Scale multiplier.", false, commonschema.Ptr(0.1), commonschema.Ptr(10)),
			stringParam("artboardId", "Optional artboard identifier.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "export.jpg",
		Domain:      "export",
		Description: "Export the document or selection as JPG.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("outputPath", "Destination path.", true),
			numberParam("quality", "JPEG quality.", false, commonschema.Ptr(1), commonschema.Ptr(100)),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "export.svg",
		Domain:      "export",
		Description: "Export the document or selection as SVG.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("outputPath", "Destination path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "export.pdf",
		Domain:      "export",
		Description: "Export the active document as PDF.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("outputPath", "Destination path.", true),
			stringParam("preset", "PDF preset name.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "export.ai",
		Domain:      "export",
		Description: "Save the active document as .ai to a new path.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("outputPath", "Destination path.", true),
		},
	})
}

func registerInspect() {
	commonschema.Register(commonschema.Command{
		Name:        "inspect.tree",
		Domain:      "inspect",
		Description: "Return a tree view of the current document.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "inspect.styles",
		Domain:      "inspect",
		Description: "Return document style usage summaries.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "inspect.bounds",
		Domain:      "inspect",
		Description: "Return bounds for the current selection.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "inspect.fonts",
		Domain:      "inspect",
		Description: "Return font usage for the current document.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "inspect.summary",
		Domain:      "inspect",
		Description: "Return an overall document summary for agents.",
		Params:      []commonschema.Param{},
	})
}

func stringParam(name, description string, required bool) commonschema.Param {
	return commonschema.Param{Name: name, Type: "string", Description: description, Required: required}
}

func numberParam(name, description string, required bool, min, max *float64) commonschema.Param {
	return commonschema.Param{Name: name, Type: "number", Description: description, Required: required, Minimum: min, Maximum: max}
}

func boolParam(name, description string, required bool) commonschema.Param {
	return commonschema.Param{Name: name, Type: "boolean", Description: description, Required: required}
}

func arrayParam(name, description string, required bool) commonschema.Param {
	return commonschema.Param{Name: name, Type: "array", Description: description, Required: required}
}

func pathParam(name, description string, required bool) commonschema.Param {
	return commonschema.Param{Name: name, Type: "string", Description: description, Required: required, SafePath: true}
}

func colorParam(name, description string, required bool) commonschema.Param {
	return commonschema.Param{Name: name, Type: "string", Description: description, Required: required, Pattern: colorPattern}
}

func enumParam(name, description string, required bool, values []string, lowRisk bool) commonschema.Param {
	return commonschema.Param{Name: name, Type: "string", Description: description, Required: required, Enum: values, LowRiskFuzzy: lowRisk}
}
