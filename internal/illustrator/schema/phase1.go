package schema

import "github.com/nerveband/ai-happy-design/internal/commonschema"

func init() {
	registerPhase1App()
	registerPhase1Document()
	registerWorkspace()
	registerPreference()
	registerView()
	registerMatrix()
	registerPerspective()
	registerPhase1Artboard()
	registerPhase1Selection()
	registerPageItem()
	registerPhase1Path()
	registerPhase1Text()
	registerStyle()
	registerSwatch()
	registerSpot()
	registerSymbol()
	registerPlaced()
	registerRaster()
	registerRepeat()
	registerDataset()
	registerVariable()
	registerTracePreset()
	registerCapture()
	registerPrint()
	registerPhase1Export()
}

func registerPhase1App() {
	for _, item := range []struct {
		name string
		desc string
	}{
		{name: "app.beep", desc: "Trigger Illustrator's system beep."},
		{name: "app.copy", desc: "Copy the current selection to the clipboard."},
		{name: "app.cut", desc: "Cut the current selection to the clipboard."},
		{name: "app.paste", desc: "Paste the clipboard contents into the active document."},
		{name: "app.undo", desc: "Undo the last Illustrator action."},
		{name: "app.redo", desc: "Redo the last undone Illustrator action."},
		{name: "app.redraw", desc: "Force Illustrator to redraw the current document."},
	} {
		commonschema.Register(commonschema.Command{
			Name:        item.name,
			Domain:      "app",
			Description: item.desc,
			Mutating:    item.name != "app.redraw" && item.name != "app.beep",
			Params:      []commonschema.Param{},
		})
	}

	commonschema.Register(commonschema.Command{
		Name:        "app.convert_sample_color",
		Domain:      "app",
		Description: "Convert a sample color array from one color space to another.",
		Params: []commonschema.Param{
			enumParam("sourceColorSpace", "Source color space.", true, []string{"CMYK", "RGB", "Grayscale", "LAB", "Separation", "DeviceN", "Indexed"}, true),
			arrayParam("sourceColor", "Source color component array.", true, commonschema.PtrInt(1), commonschema.PtrInt(8), &commonschema.Param{Type: "number"}),
			enumParam("destColorSpace", "Destination color space.", true, []string{"CMYK", "RGB", "Grayscale", "LAB", "Separation", "DeviceN", "Indexed"}, true),
			enumParam("purpose", "Color conversion purpose.", false, []string{"defaultpurpose", "previewpurpose", "exportpurpose", "dummypurpose"}, true),
			boolParam("sourceHasAlpha", "Whether the source sample includes alpha.", false),
			boolParam("destHasAlpha", "Whether the destination sample should include alpha.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.translate_placeholder_text",
		Domain:      "app",
		Description: "Generate placeholder text using Illustrator's built-in translation helper.",
		Params: []commonschema.Param{
			stringParam("text", "Source placeholder text.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.preset_lists",
		Domain:      "app",
		Description: "Return Illustrator preset, printer, and color-settings lists for runtime introspection.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.get_preset_file",
		Domain:      "app",
		Description: "Return the preset file path for a document preset type.",
		Params: []commonschema.Param{
			enumParam("presetType", "Illustrator document preset type.", true, []string{"BasicCMYK", "BasicRGB", "Print", "Mobile", "Video", "Web"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.get_preset_settings",
		Domain:      "app",
		Description: "Return typed preset settings for a named document preset.",
		Params: []commonschema.Param{
			stringParam("preset", "Preset name.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.load_color_settings",
		Domain:      "app",
		Description: "Load a color-settings file or explicitly disable color settings.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Color settings file path.", false),
			boolParam("disable", "Disable color settings instead of loading a file.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "app.show_presets",
		Domain:      "app",
		Description: "Reveal a preset file using Illustrator's preset browser flow.",
		Params: []commonschema.Param{
			pathParam("filePath", "Preset file path.", true),
		},
	})
}

func registerWorkspace() {
	commonschema.Register(commonschema.Command{
		Name:        "workspace.save",
		Domain:      "workspace",
		Description: "Save the current Illustrator workspace under a new name.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Workspace name.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "workspace.switch",
		Domain:      "workspace",
		Description: "Switch to a saved Illustrator workspace.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Workspace name.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "workspace.reset",
		Domain:      "workspace",
		Description: "Reset the current Illustrator workspace.",
		Mutating:    true,
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "workspace.delete",
		Domain:      "workspace",
		Description: "Delete a saved Illustrator workspace.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Workspace name.", true),
		},
	})
}

func registerPhase1Document() {
	commonschema.Register(commonschema.Command{
		Name:        "document.activate",
		Domain:      "document",
		Description: "Bring a document to the front and make it active.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("documentId", "Document identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.arrange",
		Domain:      "document",
		Description: "Arrange multiple open documents in Illustrator.",
		Mutating:    true,
		Params: []commonschema.Param{
			enumParam("layoutStyle", "Document layout arrangement style.", false, []string{"CASCADE", "HORIZONTALTILE", "VERTICALTILE"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.export_pdf_preset",
		Domain:      "document",
		Description: "Export the current document PDF preset values to a file.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("documentId", "Optional document identifier.", false),
			pathParam("filePath", "Destination PDF preset file path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.import_pdf_preset",
		Domain:      "document",
		Description: "Import PDF preset values from a file into the current document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("documentId", "Optional document identifier.", false),
			pathParam("filePath", "Source PDF preset file path.", true),
			boolParam("replacingPreset", "Whether to replace an existing imported preset.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.export_print_preset",
		Domain:      "document",
		Description: "Export the current document print preset values to a file.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("documentId", "Optional document identifier.", false),
			pathParam("filePath", "Destination print preset file path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "document.import_print_preset",
		Domain:      "document",
		Description: "Import a named print preset from a file into the current document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("documentId", "Optional document identifier.", false),
			pathParam("filePath", "Source print preset file path.", true),
			stringParam("printPreset", "Name of the print preset to import.", true),
		},
	})
}

func registerPreference() {
	commonschema.Register(commonschema.Command{
		Name:        "preference.get",
		Domain:      "preference",
		Description: "Read an Illustrator application preference key using the requested type accessor.",
		Params: []commonschema.Param{
			stringParam("key", "Preference key.", true),
			enumParam("valueType", "Typed accessor to use.", true, []string{"boolean", "integer", "real", "string"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "preference.set",
		Domain:      "preference",
		Description: "Set an Illustrator application preference key using the requested type accessor.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("key", "Preference key.", true),
			enumParam("valueType", "Typed setter to use.", true, []string{"boolean", "integer", "real", "string"}, true),
			objectParam("value", "Typed preference value wrapper.", true, []commonschema.Param{
				boolParam("boolean", "Boolean preference value.", false),
				numberParam("number", "Numeric preference value.", false, nil, nil),
				stringParam("string", "String preference value.", false),
			}),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "preference.delete",
		Domain:      "preference",
		Description: "Delete an Illustrator application preference key.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("key", "Preference key.", true),
		},
	})
}

func registerView() {
	commonschema.Register(commonschema.Command{
		Name:        "view.info",
		Domain:      "view",
		Description: "Return the active document view state.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        "view.set_screen_mode",
		Domain:      "view",
		Description: "Set the active document view screen mode.",
		Mutating:    true,
		Params: []commonschema.Param{
			enumParam("mode", "Illustrator screen mode.", true, []string{"DESKTOP", "MULTIWINDOW", "FULLSCREEN"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "view.set_zoom",
		Domain:      "view",
		Description: "Set the active document view zoom factor.",
		Mutating:    true,
		Params: []commonschema.Param{
			numberParam("zoom", "View zoom factor.", true, commonschema.Ptr(0.01), nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "view.set_ruler_visibility",
		Domain:      "view",
		Description: "Show or hide document rulers by normalizing the Illustrator menu toggle.",
		Mutating:    true,
		Params: []commonschema.Param{
			boolParam("visible", "Whether rulers should be visible.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "view.set_transparency_grid_visibility",
		Domain:      "view",
		Description: "Show or hide the transparency grid by normalizing the Illustrator menu toggle.",
		Mutating:    true,
		Params: []commonschema.Param{
			boolParam("visible", "Whether the transparency grid should be visible.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "view.set_center",
		Domain:      "view",
		Description: "Set the active document view center point.",
		Mutating:    true,
		Params: []commonschema.Param{
			numberParam("x", "Center x coordinate in points.", true, nil, nil),
			numberParam("y", "Center y coordinate in points.", true, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "view.rotate",
		Domain:      "view",
		Description: "Set the active document view rotation angle.",
		Mutating:    true,
		Params: []commonschema.Param{
			numberParam("angle", "Rotation angle in degrees.", true, nil, nil),
		},
	})
}

func registerMatrix() {
	commonschema.Register(commonschema.Command{Name: "matrix.identity", Domain: "matrix", Description: "Return the identity matrix.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.rotation",
		Domain:      "matrix",
		Description: "Return a rotation matrix.",
		Params: []commonschema.Param{
			numberParam("angle", "Rotation angle in degrees.", true, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.scale",
		Domain:      "matrix",
		Description: "Return a scale matrix.",
		Params: []commonschema.Param{
			numberParam("scaleX", "Horizontal scale percentage.", true, nil, nil),
			numberParam("scaleY", "Vertical scale percentage.", false, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.translation",
		Domain:      "matrix",
		Description: "Return a translation matrix.",
		Params: []commonschema.Param{
			numberParam("deltaX", "Horizontal translation in points.", true, nil, nil),
			numberParam("deltaY", "Vertical translation in points.", false, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.concatenate",
		Domain:      "matrix",
		Description: "Concatenate two matrix records.",
		Params: []commonschema.Param{
			matrixParam("matrix", "Base matrix.", true),
			matrixParam("secondMatrix", "Matrix to apply after the base matrix.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.concatenate_rotation",
		Domain:      "matrix",
		Description: "Concatenate a rotation onto a base matrix.",
		Params: []commonschema.Param{
			matrixParam("matrix", "Base matrix.", true),
			numberParam("angle", "Rotation angle in degrees.", true, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.concatenate_scale",
		Domain:      "matrix",
		Description: "Concatenate a scale onto a base matrix.",
		Params: []commonschema.Param{
			matrixParam("matrix", "Base matrix.", true),
			numberParam("scaleX", "Horizontal scale percentage.", true, nil, nil),
			numberParam("scaleY", "Vertical scale percentage.", false, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.concatenate_translation",
		Domain:      "matrix",
		Description: "Concatenate a translation onto a base matrix.",
		Params: []commonschema.Param{
			matrixParam("matrix", "Base matrix.", true),
			numberParam("deltaX", "Horizontal translation in points.", true, nil, nil),
			numberParam("deltaY", "Vertical translation in points.", false, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.invert",
		Domain:      "matrix",
		Description: "Invert a matrix record.",
		Params: []commonschema.Param{
			matrixParam("matrix", "Matrix to invert.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.compare",
		Domain:      "matrix",
		Description: "Compare two matrices for equality.",
		Params: []commonschema.Param{
			matrixParam("matrix", "Base matrix.", true),
			matrixParam("secondMatrix", "Matrix to compare.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "matrix.singular",
		Domain:      "matrix",
		Description: "Return whether a matrix is singular.",
		Params: []commonschema.Param{
			matrixParam("matrix", "Matrix to inspect.", true),
		},
	})
}

func registerPerspective() {
	commonschema.Register(commonschema.Command{Name: "perspective.show", Domain: "perspective", Description: "Show the active document perspective grid.", Mutating: true, Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{Name: "perspective.hide", Domain: "perspective", Description: "Hide the active document perspective grid.", Mutating: true, Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{Name: "perspective.get_active_plane", Domain: "perspective", Description: "Return the active perspective grid plane.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "perspective.set_active_plane",
		Domain:      "perspective",
		Description: "Set the active perspective grid plane.",
		Mutating:    true,
		Params: []commonschema.Param{
			enumParam("plane", "Perspective grid plane.", true, []string{"GRIDLEFTPLANETYPE", "GRIDRIGHTPLANETYPE", "GRIDFLOORPLANETYPE", "INVALIDGRIDPLANETYPE"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "perspective.select_preset",
		Domain:      "perspective",
		Description: "Select a named perspective grid preset.",
		Mutating:    true,
		Params: []commonschema.Param{
			enumParam("gridType", "Perspective grid type.", true, []string{"OnePointPerspectiveGridType", "TwoPointPerspectiveGridType", "ThreePointPerspectiveGridType", "InvalidPerspectiveGridType"}, true),
			stringParam("presetName", "Perspective preset name.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "perspective.import_preset",
		Domain:      "perspective",
		Description: "Import perspective grid preset definitions from a file.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Perspective preset file path.", true),
			stringParam("presetName", "Optional preset name to import.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "perspective.export_preset",
		Domain:      "perspective",
		Description: "Export the current perspective grid preset definitions to a file.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Perspective preset output file path.", true),
		},
	})
}

func registerPhase1Artboard() {
	commonschema.Register(commonschema.Command{
		Name:        "artboard.delete",
		Domain:      "artboard",
		Description: "Delete an artboard by identifier. Illustrator cannot remove the last artboard in a document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("artboardId", "Artboard identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "artboard.rearrange",
		Domain:      "artboard",
		Description: "Rearrange artboards in the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			enumParam("artboardLayout", "Layout for artboard rearrangement.", false, []string{"GridByRow", "GridByCol", "Row", "Column", "RLGridByRow", "RLGridByCol", "RLRow"}, true),
			numberParam("artboardRowsOrCols", "Rows or columns to use for multi-artboard layouts.", false, commonschema.Ptr(1), commonschema.Ptr(100)),
			numberParam("artboardSpacing", "Spacing between artboards.", false, commonschema.Ptr(0), nil),
			boolParam("moveArtwork", "Whether to move artwork with the artboards.", false),
		},
	})
}

func registerPhase1Selection() {
	commonschema.Register(commonschema.Command{
		Name:        "selection.select_active_artboard_objects",
		Domain:      "selection",
		Description: "Select objects on the currently active artboard.",
		Mutating:    true,
		Params:      []commonschema.Param{},
	})
}

func registerPageItem() {
	commonschema.Register(commonschema.Command{
		Name:        "page_item.remove",
		Domain:      "page_item",
		Description: "Delete a page item from the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Page item identifier.", true),
		},
	})
	for _, item := range []struct {
		name        string
		description string
		required    []commonschema.Param
		optional    []commonschema.Param
	}{
		{
			name:        "page_item.duplicate",
			description: "Duplicate a page item into another layer or relative position.",
			required:    []commonschema.Param{identifierParam("itemId", "Page item identifier.", true)},
			optional: []commonschema.Param{
				identifierParam("targetId", "Optional destination layer or page-item identifier.", false),
				enumParam("placement", "Illustrator placement mode.", false, []string{"inside", "before", "after", "at_beginning", "at_end"}, true),
			},
		},
		{
			name:        "page_item.move",
			description: "Move a page item relative to another layer or page item.",
			required: []commonschema.Param{
				identifierParam("itemId", "Page item identifier.", true),
				identifierParam("targetId", "Destination layer or page-item identifier.", true),
			},
			optional: []commonschema.Param{
				enumParam("placement", "Illustrator placement mode.", false, []string{"inside", "before", "after", "at_beginning", "at_end"}, true),
			},
		},
		{
			name:        "page_item.resize",
			description: "Resize a page item using Illustrator's generic page-item API.",
			required: []commonschema.Param{
				identifierParam("itemId", "Page item identifier.", true),
				numberParam("scaleX", "Horizontal scale percentage.", true, nil, nil),
				numberParam("scaleY", "Vertical scale percentage.", true, nil, nil),
			},
			optional: pageItemTransformOptions(true),
		},
		{
			name:        "page_item.rotate",
			description: "Rotate a page item using Illustrator's generic page-item API.",
			required: []commonschema.Param{
				identifierParam("itemId", "Page item identifier.", true),
				numberParam("angle", "Rotation angle in degrees.", true, nil, nil),
			},
			optional: pageItemTransformOptions(false),
		},
		{
			name:        "page_item.transform",
			description: "Apply a matrix transform to a page item.",
			required: []commonschema.Param{
				identifierParam("itemId", "Page item identifier.", true),
				matrixParam("matrix", "Transformation matrix.", true),
			},
			optional: pageItemTransformOptions(true),
		},
	} {
		params := append([]commonschema.Param{}, item.required...)
		params = append(params, item.optional...)
		commonschema.Register(commonschema.Command{
			Name:        item.name,
			Domain:      "page_item",
			Description: item.description,
			Mutating:    true,
			Params:      params,
		})
	}
	commonschema.Register(commonschema.Command{
		Name:        "page_item.translate",
		Domain:      "page_item",
		Description: "Translate a page item by delta values.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Page item identifier.", true),
			numberParam("deltaX", "Horizontal translation delta.", false, nil, nil),
			numberParam("deltaY", "Vertical translation delta.", false, nil, nil),
			boolParam("transformObjects", "Whether to move the object geometry.", false),
			boolParam("transformFillPatterns", "Whether to move fill patterns.", false),
			boolParam("transformFillGradients", "Whether to move fill gradients.", false),
			boolParam("transformStrokePatterns", "Whether to move stroke patterns.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "page_item.z_order",
		Domain:      "page_item",
		Description: "Reorder a page item in its parent's stacking order.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Page item identifier.", true),
			enumParam("method", "Illustrator z-order method.", true, []string{"BRINGTOFRONT", "BRINGFORWARD", "SENDBACKWARD", "SENDTOBACK"}, true),
		},
	})
}

func pageItemTransformOptions(includeLineWidths bool) []commonschema.Param {
	fields := []commonschema.Param{
		boolParam("changePositions", "Whether to move the item while transforming.", false),
		boolParam("changeFillPatterns", "Whether to transform fill patterns.", false),
		boolParam("changeFillGradients", "Whether to transform fill gradients.", false),
		boolParam("changeStrokePatterns", "Whether to transform stroke patterns.", false),
		enumParam("anchor", "Transformation anchor point.", false, []string{"BOTTOM", "BOTTOMLEFT", "BOTTOMRIGHT", "CENTER", "DOCUMENTORIGIN", "LEFT", "RIGHT", "TOP", "TOPLEFT", "TOPRIGHT"}, true),
	}
	if includeLineWidths {
		fields = append(fields, numberParam("changeLineWidths", "Optional line width scaling percentage.", false, nil, nil))
	}
	return fields
}

func registerPhase1Path() {
	commonschema.Register(commonschema.Command{
		Name:        "path.create_polygon",
		Domain:      "path",
		Description: "Create a polygon path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Path name.", true),
			numberParam("left", "Center x coordinate.", true, nil, nil),
			numberParam("top", "Center y coordinate.", true, nil, nil),
			numberParam("radius", "Polygon radius.", true, commonschema.Ptr(1), nil),
			numberParam("sides", "Polygon sides.", true, commonschema.Ptr(3), commonschema.Ptr(100)),
			identifierParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "path.create_rounded_rect",
		Domain:      "path",
		Description: "Create a rounded rectangle path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Path name.", true),
			numberParam("left", "Left coordinate.", true, nil, nil),
			numberParam("top", "Top coordinate.", true, nil, nil),
			numberParam("width", "Width.", true, commonschema.Ptr(1), nil),
			numberParam("height", "Height.", true, commonschema.Ptr(1), nil),
			numberParam("horizontalRadius", "Horizontal corner radius.", false, commonschema.Ptr(0), nil),
			numberParam("verticalRadius", "Vertical corner radius.", false, commonschema.Ptr(0), nil),
			boolParam("reversed", "Whether the path direction should be reversed.", false),
			identifierParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "path.create_star",
		Domain:      "path",
		Description: "Create a star path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Path name.", true),
			numberParam("left", "Center x coordinate.", true, nil, nil),
			numberParam("top", "Center y coordinate.", true, nil, nil),
			numberParam("outerRadius", "Outer radius.", true, commonschema.Ptr(1), nil),
			numberParam("innerRadius", "Inner radius.", true, commonschema.Ptr(1), nil),
			numberParam("points", "Star point count.", true, commonschema.Ptr(3), commonschema.Ptr(100)),
			identifierParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "path.set_entire_path",
		Domain:      "path",
		Description: "Replace the full point list of an existing path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Path item identifier.", true),
			arrayParam("points", "Point array in Illustrator coordinate space.", true, commonschema.PtrInt(2), nil, &commonschema.Param{
				Type:     "array",
				MinItems: commonschema.PtrInt(2),
				MaxItems: commonschema.PtrInt(2),
				Items: &commonschema.Param{
					Type: "number",
				},
			}),
			boolParam("closed", "Whether the path should be closed after replacing the point list.", false),
		},
	})
}

func registerPhase1Text() {
	commonschema.Register(commonschema.Command{
		Name:        "text.create_area",
		Domain:      "text",
		Description: "Create an area text frame from a rectangle path.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Text item name.", true),
			stringParam("contents", "Text contents.", true),
			numberParam("left", "Left coordinate.", true, nil, nil),
			numberParam("top", "Top coordinate.", true, nil, nil),
			numberParam("width", "Area text width.", true, commonschema.Ptr(1), nil),
			numberParam("height", "Area text height.", true, commonschema.Ptr(1), nil),
			identifierParam("layerId", "Optional target layer.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "text.create_on_path",
		Domain:      "text",
		Description: "Create a path text frame on an existing path item.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Text item name.", true),
			stringParam("contents", "Text contents.", true),
			identifierParam("pathItemId", "Existing path item identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "text.change_case",
		Domain:      "text",
		Description: "Change the capitalization of a text frame's full text range.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Text frame identifier.", true),
			enumParam("caseType", "Illustrator case-change type.", true, []string{"LOWERCASE", "SENTENCECASE", "TITLECASE", "UPPERCASE"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "text.thread",
		Domain:      "text",
		Description: "Thread one text frame into another.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Source text frame identifier.", true),
			identifierParam("nextItemId", "Destination text frame identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "text.convert_to_area",
		Domain:      "text",
		Description: "Convert a point text frame to area text.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Text frame identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "text.convert_to_point",
		Domain:      "text",
		Description: "Convert an area text frame to point text.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Text frame identifier.", true),
		},
	})
}

func registerStyle() {
	for _, kind := range []struct {
		domain string
		name   string
		desc   string
	}{
		{domain: "style.character", name: "style.character.list", desc: "List character styles in the active document."},
		{domain: "style.paragraph", name: "style.paragraph.list", desc: "List paragraph styles in the active document."},
		{domain: "style.graphic", name: "style.graphic.list", desc: "List graphic styles in the active document."},
	} {
		commonschema.Register(commonschema.Command{Name: kind.name, Domain: kind.domain, Description: kind.desc, Params: []commonschema.Param{}})
	}

	commonschema.Register(commonschema.Command{
		Name:        "style.character.apply",
		Domain:      "style.character",
		Description: "Apply a character style to a text frame.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("styleName", "Character style name.", true),
			identifierParam("itemId", "Text frame identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "style.character.import",
		Domain:      "style.character",
		Description: "Import character styles from an Illustrator document.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Source Illustrator document path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "style.paragraph.apply",
		Domain:      "style.paragraph",
		Description: "Apply a paragraph style to all paragraphs in a text frame.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("styleName", "Paragraph style name.", true),
			identifierParam("itemId", "Text frame identifier.", true),
			boolParam("clearOverrides", "Whether to clear local paragraph overrides.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "style.paragraph.import",
		Domain:      "style.paragraph",
		Description: "Import paragraph styles from an Illustrator document.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Source Illustrator document path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "style.graphic.apply",
		Domain:      "style.graphic",
		Description: "Apply a graphic style to a page item.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("styleName", "Graphic style name.", true),
			identifierParam("itemId", "Page item identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "style.graphic.merge",
		Domain:      "style.graphic",
		Description: "Merge a page item's appearance into an existing graphic style.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("styleName", "Graphic style name.", true),
			identifierParam("itemId", "Page item identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "style.graphic.remove",
		Domain:      "style.graphic",
		Description: "Remove a non-default graphic style from the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("styleName", "Graphic style name.", true),
		},
	})
}

func registerSwatch() {
	commonschema.Register(commonschema.Command{Name: "swatch.list", Domain: "swatch", Description: "List swatches in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "swatch.create",
		Domain:      "swatch",
		Description: "Create a document swatch from an RGB hex or CMYK value.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Swatch name.", true),
			enumParam("colorMode", "Swatch color mode.", true, []string{"RGB", "CMYK"}, true),
			colorParam("hex", "RGB swatch color.", false),
			numberParam("cyan", "CMYK cyan value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("magenta", "CMYK magenta value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("yellow", "CMYK yellow value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("black", "CMYK black value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "swatch.delete",
		Domain:      "swatch",
		Description: "Delete a swatch by name.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("swatchId", "Swatch name or identifier.", true),
		},
	})
}

func registerSpot() {
	commonschema.Register(commonschema.Command{Name: "spot.list", Domain: "spot", Description: "List spots in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "spot.create",
		Domain:      "spot",
		Description: "Create a spot color in the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("name", "Spot name.", true),
			enumParam("colorMode", "Spot color mode.", true, []string{"RGB", "CMYK"}, true),
			enumParam("colorType", "Spot color model.", false, []string{"SPOT", "PROCESS", "REGISTRATION"}, true),
			colorParam("hex", "RGB spot color.", false),
			numberParam("cyan", "CMYK cyan value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("magenta", "CMYK magenta value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("yellow", "CMYK yellow value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("black", "CMYK black value.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "spot.delete",
		Domain:      "spot",
		Description: "Delete a spot color by name.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("spotId", "Spot name or identifier.", true),
		},
	})
}

func registerSymbol() {
	commonschema.Register(commonschema.Command{Name: "symbol.list", Domain: "symbol", Description: "List symbols and symbol instances in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "symbol.create",
		Domain:      "symbol",
		Description: "Create a new symbol from an existing art item.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Source art item identifier.", true),
			stringParam("name", "Optional symbol name override.", false),
			enumParam("registrationPoint", "Symbol registration point.", false, []string{"SYMBOLTOPLEFTPOINT", "SYMBOLTOPMIDDLEPOINT", "SYMBOLTOPRIGHTPOINT", "SYMBOLMIDDLELEFTPOINT", "SYMBOLCENTERPOINT", "SYMBOLMIDDLERIGHTPOINT", "SYMBOLBOTTOMLEFTPOINT", "SYMBOLBOTTOMMIDDLEPOINT", "SYMBOLBOTTOMRIGHTPOINT"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "symbol.place",
		Domain:      "symbol",
		Description: "Place a symbol instance into the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("symbolId", "Symbol identifier.", true),
			stringParam("name", "Optional symbol item name.", false),
			numberParam("left", "Instance left coordinate.", false, nil, nil),
			numberParam("top", "Instance top coordinate.", false, nil, nil),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "symbol.break_link",
		Domain:      "symbol",
		Description: "Break the link from a symbol instance back to its source symbol.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Symbol item identifier.", true),
		},
	})
}

func registerPlaced() {
	commonschema.Register(commonschema.Command{Name: "placed.list", Domain: "placed", Description: "List linked placed items in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "placed.place",
		Domain:      "placed",
		Description: "Place linked artwork into the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Linked file path.", true),
			stringParam("name", "Optional placed item name.", false),
			numberParam("left", "Placed item left coordinate.", false, nil, nil),
			numberParam("top", "Placed item top coordinate.", false, nil, nil),
			boolParam("embed", "Embed the placed item immediately.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "placed.embed",
		Domain:      "placed",
		Description: "Embed a linked placed item into the current document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Placed item identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "placed.relink",
		Domain:      "placed",
		Description: "Relink a placed item to a new source file.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Placed item identifier.", true),
			pathParam("filePath", "New linked file path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "placed.trace",
		Domain:      "placed",
		Description: "Trace a placed item into a plugin tracing object.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Placed item identifier.", true),
			stringParam("name", "Optional traced plugin item name.", false),
			stringParam("presetName", "Optional tracing preset name.", false),
			boolParam("expand", "Expand the trace into a group item.", false),
			boolParam("viewed", "Preserve viewed tracing overlays on expand.", false),
		},
	})
}

func registerRaster() {
	commonschema.Register(commonschema.Command{Name: "raster.list", Domain: "raster", Description: "List raster items in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "raster.rasterize",
		Domain:      "raster",
		Description: "Rasterize an existing page item into a raster item.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Source art item identifier.", true),
			stringParam("name", "Optional raster item name.", false),
			rectParam("clipBounds", "Optional clip bounds for rasterization.", false),
			numberParam("resolution", "Raster output resolution.", false, commonschema.Ptr(72), commonschema.Ptr(2400)),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "raster.trace",
		Domain:      "raster",
		Description: "Trace a raster item into a plugin tracing object.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Raster item identifier.", true),
			stringParam("name", "Optional traced plugin item name.", false),
			stringParam("presetName", "Optional tracing preset name.", false),
			boolParam("expand", "Expand the trace into a group item.", false),
			boolParam("viewed", "Preserve viewed tracing overlays on expand.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "raster.colorize",
		Domain:      "raster",
		Description: "Colorize a raster item using an RGB or CMYK color payload.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Raster item identifier.", true),
			enumParam("colorMode", "Color payload mode.", false, []string{"RGB", "CMYK"}, true),
			colorParam("hex", "RGB color hex.", false),
			numberParam("cyan", "CMYK cyan.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("magenta", "CMYK magenta.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("yellow", "CMYK yellow.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
			numberParam("black", "CMYK black.", false, commonschema.Ptr(0), commonschema.Ptr(100)),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "raster.release_tracing",
		Domain:      "raster",
		Description: "Release a tracing plugin item back to its original raster or placed art.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Tracing plugin item identifier.", true),
		},
	})
}

func registerRepeat() {
	registerRepeatFamily("grid", "Grid repeat", []commonschema.Param{
		stringParam("columnFlipType", "Optional grid column flip type.", false),
		numberParam("horizontalSpacing", "Horizontal spacing.", false, nil, nil),
		stringParam("patternType", "Optional grid pattern type.", false),
		stringParam("rowFlipType", "Optional grid row flip type.", false),
		numberParam("verticalSpacing", "Vertical spacing.", false, nil, nil),
	}, []string{"GRIDALL", "HORIZONTALSPACING", "VERTICALSPACING"})

	registerRepeatFamily("radial", "Radial repeat", []commonschema.Param{
		numberParam("numberOfInstances", "Instance count.", false, commonschema.Ptr(1), commonschema.Ptr(1024)),
		numberParam("radius", "Radius of art.", false, nil, nil),
		boolParam("reverseOverlap", "Whether overlap direction is reversed.", false),
	}, []string{"NUMBEROFINSTANCES", "RADIALALL", "RADIUSOFART", "REVERSEOVERLAP"})

	registerRepeatFamily("symmetry", "Symmetry repeat", []commonschema.Param{
		numberParam("axisRotationAngleInRadians", "Axis rotation angle in radians.", false, nil, nil),
	}, []string{"AXISROTATION", "SYMMETRYALL"})
}

func registerRepeatFamily(kind, label string, configFields []commonschema.Param, states []string) {
	domain := "repeat." + kind
	commonschema.Register(commonschema.Command{
		Name:        domain + ".list",
		Domain:      domain,
		Description: "List " + label + " items in the active document.",
		Params:      []commonschema.Param{},
	})
	commonschema.Register(commonschema.Command{
		Name:        domain + ".create",
		Domain:      domain,
		Description: "Create a " + label + " from an existing art item.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("itemId", "Source art item identifier.", true),
			stringParam("name", "Optional repeat item name.", false),
			objectParam("config", label+" configuration.", false, configFields),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        domain + ".update",
		Domain:      domain,
		Description: "Update a " + label + " configuration.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("repeatId", "Repeat item identifier.", true),
			objectParam("config", label+" configuration.", true, configFields),
			enumParam("state", "Repeat update selector.", false, states, true),
		},
	})
}

func registerDataset() {
	commonschema.Register(commonschema.Command{Name: "dataset.list", Domain: "dataset", Description: "List datasets in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "dataset.create",
		Domain:      "dataset",
		Description: "Create a dataset from the current bound artwork state.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("name", "Optional dataset identifier.", false),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "dataset.apply",
		Domain:      "dataset",
		Description: "Display a dataset in the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("datasetId", "Dataset identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "dataset.update",
		Domain:      "dataset",
		Description: "Update a dataset from the current artwork bindings.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("datasetId", "Dataset identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "dataset.delete",
		Domain:      "dataset",
		Description: "Delete a dataset from the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("datasetId", "Dataset identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "dataset.import",
		Domain:      "dataset",
		Description: "Import variables and datasets from an XML library file.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Variables XML library path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "dataset.export",
		Domain:      "dataset",
		Description: "Export variables and datasets to an XML library file.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Variables XML output path.", true),
		},
	})
}

func registerVariable() {
	commonschema.Register(commonschema.Command{Name: "variable.list", Domain: "variable", Description: "List variables in the active document.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "variable.create",
		Domain:      "variable",
		Description: "Create a variable in the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("name", "Optional variable identifier.", false),
			enumParam("kind", "Variable kind.", true, []string{"GRAPH", "IMAGE", "VISIBILITY", "TEXTUAL"}, true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "variable.delete",
		Domain:      "variable",
		Description: "Delete a variable from the active document.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("variableId", "Variable identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "variable.bind_visibility",
		Domain:      "variable",
		Description: "Bind a visibility variable to an art item.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("variableId", "Variable identifier.", true),
			identifierParam("itemId", "Art item identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "variable.bind_text",
		Domain:      "variable",
		Description: "Bind a textual content variable to a text frame.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("variableId", "Variable identifier.", true),
			identifierParam("itemId", "Text frame identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "variable.bind_content",
		Domain:      "variable",
		Description: "Bind a content variable to a text, placed, or raster item that supports content bindings.",
		Mutating:    true,
		Params: []commonschema.Param{
			identifierParam("variableId", "Variable identifier.", true),
			identifierParam("itemId", "Content-bindable item identifier.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "variable.import",
		Domain:      "variable",
		Description: "Import variables and datasets from an XML library file.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Variables XML library path.", true),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "variable.export",
		Domain:      "variable",
		Description: "Export variables and datasets to an XML library file.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("filePath", "Variables XML output path.", true),
		},
	})
}

func registerTracePreset() {
	commonschema.Register(commonschema.Command{
		Name:        "trace.preset.list",
		Domain:      "trace.preset",
		Description: "List Illustrator tracing presets available to the scripting runtime.",
		Params:      []commonschema.Param{},
	})
}

func registerCapture() {
	commonschema.Register(commonschema.Command{
		Name:        "capture.image",
		Domain:      "capture",
		Description: "Capture artwork to a raster image file using Document.imageCapture.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("outputPath", "Raster image output path.", true),
			rectParam("clipBounds", "Optional capture clip bounds.", false),
			objectParam("options", "Image capture options.", false, []commonschema.Param{
				boolParam("antiAliasing", "Whether to anti-alias the capture.", false),
				boolParam("matte", "Whether to matte the artboard background.", false),
				colorParam("matteColor", "Optional matte color.", false),
				numberParam("resolution", "Capture resolution in PPI.", false, commonschema.Ptr(72), commonschema.Ptr(2400)),
				boolParam("transparency", "Whether to preserve transparency.", false),
			}),
		},
	})
	commonschema.Register(commonschema.Command{
		Name:        "capture.window",
		Domain:      "capture",
		Description: "Capture the current Illustrator window to a TIFF image file.",
		Mutating:    true,
		Params: []commonschema.Param{
			pathParam("outputPath", "Window capture TIFF path.", true),
			numberParam("width", "Capture width.", true, commonschema.Ptr(1), nil),
			numberParam("height", "Capture height.", true, commonschema.Ptr(1), nil),
		},
	})
}

func registerPrint() {
	commonschema.Register(commonschema.Command{Name: "print.presets", Domain: "print", Description: "List Illustrator print presets.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{Name: "print.devices", Domain: "print", Description: "List installed Illustrator printers and PPDs.", Params: []commonschema.Param{}})
	commonschema.Register(commonschema.Command{
		Name:        "print.run",
		Domain:      "print",
		Description: "Print the active document with typed print options.",
		Mutating:    true,
		Params: []commonschema.Param{
			stringParam("printerName", "Optional printer name.", false),
			{Name: "PPDName", Type: "string", Description: "Optional PPD name.", Aliases: []string{"ppdName"}},
			stringParam("printPreset", "Optional print preset name.", false),
			objectParam("jobOptions", "Optional print job options.", false, []commonschema.Param{
				stringParam("name", "Print job name.", false),
				numberParam("copies", "Number of copies.", false, commonschema.Ptr(1), nil),
				boolParam("collate", "Whether to collate pages.", false),
				boolParam("printAllArtboards", "Whether to print all artboards.", false),
				stringParam("artboardRange", "Artboard range when printAllArtboards is false.", false),
				boolParam("printAsBitmap", "Whether to print as bitmap.", false),
				boolParam("reversePages", "Whether to print pages in reverse order.", false),
			}),
		},
	})
}

func registerPhase1Export() {
	for _, item := range []struct {
		name   string
		desc   string
		params []commonschema.Param
	}{
		{name: "export.gif", desc: "Export the active document as GIF.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true), numberParam("scale", "Scale multiplier.", false, commonschema.Ptr(0.1), commonschema.Ptr(10)), identifierParam("artboardId", "Optional artboard identifier.", false)}},
		{name: "export.png8", desc: "Export the active document as PNG8.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true), numberParam("scale", "Scale multiplier.", false, commonschema.Ptr(0.1), commonschema.Ptr(10)), numberParam("colorCount", "Indexed color count.", false, commonschema.Ptr(2), commonschema.Ptr(256)), identifierParam("artboardId", "Optional artboard identifier.", false)}},
		{name: "export.tiff", desc: "Export the active document as TIFF.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true), numberParam("resolution", "Output resolution in DPI.", false, commonschema.Ptr(72), commonschema.Ptr(2400)), identifierParam("artboardId", "Optional artboard identifier.", false)}},
		{name: "export.webp", desc: "Export the active document as WebP.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true), boolParam("losslessCompression", "Whether to use lossless compression.", false), numberParam("imageQuality", "WebP image quality for lossy compression.", false, commonschema.Ptr(0), commonschema.Ptr(100)), boolParam("isTransparent", "Whether to preserve transparency.", false), numberParam("ppi", "WebP export resolution.", false, commonschema.Ptr(4), commonschema.Ptr(2400))}},
		{name: "export.photoshop", desc: "Export the active document as Photoshop PSD.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true), numberParam("resolution", "Output resolution in DPI.", false, commonschema.Ptr(72), commonschema.Ptr(2400)), boolParam("editableText", "Whether to keep text editable.", false), boolParam("maximumEditability", "Whether to preserve editability.", false), identifierParam("artboardId", "Optional artboard identifier.", false)}},
		{name: "export.autocad", desc: "Export the active document as AutoCAD DXF or DWG.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true)}},
		{name: "export.eps", desc: "Save the active document as EPS.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true)}},
		{name: "export.fxg", desc: "Save the active document as FXG.", params: []commonschema.Param{pathParam("outputPath", "Destination path.", true)}},
	} {
		commonschema.Register(commonschema.Command{
			Name:        item.name,
			Domain:      "export",
			Description: item.desc,
			Mutating:    true,
			Params:      item.params,
		})
	}
}

func objectParam(name, description string, required bool, fields []commonschema.Param) commonschema.Param {
	return commonschema.Param{Name: name, Type: "object", Description: description, Required: required, Fields: fields}
}

func rectParam(name, description string, required bool) commonschema.Param {
	return objectParam(name, description, required, []commonschema.Param{
		numberParam("left", "Left coordinate.", true, nil, nil),
		numberParam("top", "Top coordinate.", true, nil, nil),
		numberParam("right", "Right coordinate.", true, nil, nil),
		numberParam("bottom", "Bottom coordinate.", true, nil, nil),
	})
}

func matrixParam(name, description string, required bool) commonschema.Param {
	return objectParam(name, description, required, []commonschema.Param{
		numberParam("mValueA", "Matrix coefficient a.", true, nil, nil),
		numberParam("mValueB", "Matrix coefficient b.", true, nil, nil),
		numberParam("mValueC", "Matrix coefficient c.", true, nil, nil),
		numberParam("mValueD", "Matrix coefficient d.", true, nil, nil),
		numberParam("mValueTX", "Matrix translation x.", true, nil, nil),
		numberParam("mValueTY", "Matrix translation y.", true, nil, nil),
	})
}
