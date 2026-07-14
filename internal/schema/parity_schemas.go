package schema

func init() {
	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	nodeIDs := Param{Name: "nodeIds", Type: "array", Required: true}
	parentID := Param{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`}
	color := Param{Name: "color", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`}
	sizeParams := []Param{{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000)}, {Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000)}}
	positionParams := []Param{{Name: "x", Type: "number"}, {Name: "y", Type: "number"}}

	for _, s := range []Schema{
		{Command: "node.rotate", Description: "Rotate a node", Params: []Param{nodeID, {Name: "rotation", Type: "number", Required: true}}},
		{Command: "node.get_css", Description: "Get CSS properties for a node using Figma getCSSAsync", Params: []Param{nodeID}},
		{Command: "node.set_opacity", Description: "Set node opacity", Params: []Param{nodeID, {Name: "opacity", Type: "number", Required: true, Min: Ptr(0), Max: Ptr(1)}}},
		{Command: "node.set_blend_mode", Description: "Set node blend mode", Params: []Param{nodeID, {Name: "blendMode", Type: "string"}}},
		{Command: "node.set_visibility", Description: "Set node visibility", Params: []Param{nodeID, {Name: "visible", Type: "boolean", Required: true}}},
		{Command: "node.set_locked", Description: "Set node locked state", Params: []Param{nodeID, {Name: "locked", Type: "boolean", Required: true}}},
		{Command: "node.rename", Description: "Rename a node", Params: []Param{nodeID, {Name: "name", Type: "string", Required: true}}},
		{Command: "node.clone", Description: "Clone a node", Params: []Param{nodeID}},
		{Command: "node.set_corner_radius", Description: "Set node corner radius", Params: []Param{nodeID, {Name: "cornerRadius", Type: "number", Aliases: []string{"r", "radius"}, Min: Ptr(0), Max: Ptr(500)}}},
		{Command: "node.set_mask", Description: "Create a mask group", Params: []Param{{Name: "maskNodeId", Type: "string", Required: true}, {Name: "targetIds", Type: "array", Required: true}}},
		{Command: "node.create_section", Description: "Create a section", Params: append(append([]Param{{Name: "name", Type: "string"}, parentID}, positionParams...), sizeParams...)},

		{Command: "layer.set_order", Description: "Set layer order", Params: []Param{nodeID, {Name: "index", Type: "number", Required: true, Min: Ptr(0)}}},
		{Command: "layer.bring_forward", Description: "Bring a layer forward", Params: []Param{nodeID}},
		{Command: "layer.send_backward", Description: "Send a layer backward", Params: []Param{nodeID}},
		{Command: "layer.bring_to_front", Description: "Bring a layer to front", Params: []Param{nodeID}},
		{Command: "layer.send_to_back", Description: "Send a layer to back", Params: []Param{nodeID}},
		{Command: "layer.group", Description: "Group nodes", Params: []Param{nodeIDs, parentID}},
		{Command: "layer.ungroup", Description: "Ungroup nodes", Params: []Param{nodeID}},
		{Command: "layer.move_to_parent", Description: "Move a node to a parent", Params: []Param{nodeID, {Name: "parentId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}, {Name: "index", Type: "number", Min: Ptr(0)}}},
		{Command: "layer.insert_child", Description: "Insert a child into a parent", Params: []Param{{Name: "parentId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}, {Name: "childId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}, {Name: "index", Type: "number", Min: Ptr(0)}}},

		{Command: "layout.set_padding", Description: "Set frame padding", Params: []Param{nodeID, {Name: "padding", Type: "number", Min: Ptr(0), Max: Ptr(500)}, {Name: "paddingTop", Type: "number"}, {Name: "paddingRight", Type: "number"}, {Name: "paddingBottom", Type: "number"}, {Name: "paddingLeft", Type: "number"}}},
		{Command: "layout.set_spacing", Description: "Set auto-layout spacing", Params: []Param{nodeID, {Name: "itemSpacing", Type: "number", Aliases: []string{"spacing"}, Min: Ptr(0), Max: Ptr(500)}}},
		{Command: "layout.set_alignment", Description: "Set auto-layout alignment", Params: []Param{nodeID, {Name: "primaryAxisAlign", Type: "string"}, {Name: "counterAxisAlign", Type: "string"}}},
		{Command: "layout.set_constraints", Description: "Set layout constraints", Params: []Param{nodeID, {Name: "constraintHorizontal", Type: "string"}, {Name: "constraintVertical", Type: "string"}}},
		{Command: "layout.set_layout_wrap", Description: "Set layout wrapping", Params: []Param{nodeID, {Name: "layoutWrap", Type: "string", Enum: []string{"WRAP", "NO_WRAP"}}}},
		{Command: "layout.remove_auto_layout", Description: "Remove auto-layout", Params: []Param{nodeID}},
		{Command: "layout.check_overlaps", Description: "Check child overlaps", Params: []Param{nodeID}},
		{Command: "layout.set_grid", Description: "Set layout grid", Params: []Param{nodeID, {Name: "gridType", Type: "string"}, {Name: "count", Type: "number"}, {Name: "size", Type: "number"}, {Name: "gutter", Type: "number"}}},
		{Command: "layout.get_grids", Description: "Get layout grids", Params: []Param{nodeID}},
		{Command: "layout.remove_grids", Description: "Remove layout grids", Params: []Param{nodeID}},

		{Command: "paint.set_image_url", Aliases: []string{"paint.set_image_fill_from_url"}, Description: "Set an image fill from a URL", Params: []Param{nodeID, {Name: "url", Type: "string", Required: true}, {Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}}}},
		{Command: "paint.add_fill", Description: "Add a fill", Params: []Param{nodeID, color, {Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)}}},
		{Command: "paint.remove_fill", Description: "Remove a fill", Params: []Param{nodeID, {Name: "fillIndex", Type: "number", Aliases: []string{"index"}, Min: Ptr(0)}}},
		{Command: "paint.get_fills", Description: "Get fills", Params: []Param{nodeID}},

		{Command: "shape.create_polygon", Description: "Create a polygon", Params: append(append([]Param{parentID, {Name: "name", Type: "string"}}, positionParams...), append(sizeParams, color, Param{Name: "pointCount", Type: "number", Min: Ptr(3), Max: Ptr(60)})...)},
		{Command: "shape.create_star", Description: "Create a star", Params: append(append([]Param{parentID, {Name: "name", Type: "string"}}, positionParams...), append(sizeParams, color, Param{Name: "pointCount", Type: "number", Min: Ptr(3), Max: Ptr(60)}, Param{Name: "innerRadius", Type: "number", Min: Ptr(0), Max: Ptr(1)})...)},
		{Command: "shape.create_line", Description: "Create a line", Params: append([]Param{parentID, {Name: "name", Type: "string"}}, append(positionParams, Param{Name: "width", Type: "number", Aliases: []string{"length"}, Min: Ptr(1)}, color, Param{Name: "strokeWeight", Type: "number", Min: Ptr(0), Max: Ptr(100)})...)},
		{Command: "shape.create_from_svg", Description: "Create a node from SVG", Params: append([]Param{parentID, {Name: "name", Type: "string"}, {Name: "svgPath", Type: "string", Required: true}}, positionParams...)},
		{Command: "shape.create_vector", Description: "Create a vector node", Params: append([]Param{parentID, {Name: "name", Type: "string"}}, positionParams...)},

		{Command: "text.set_font", Description: "Set text font", Params: []Param{nodeID, {Name: "fontFamily", Type: "string", Required: true}, {Name: "fontStyle", Type: "string"}}},
		{Command: "text.set_size", Description: "Set text size", Params: []Param{nodeID, {Name: "fontSize", Type: "number", Aliases: []string{"size"}, Required: true, Min: Ptr(4), Max: Ptr(500)}}},
		{Command: "text.set_weight", Description: "Set text weight", Params: []Param{nodeID, {Name: "fontWeight", Type: "number", Aliases: []string{"weight"}, Required: true}}},
		{Command: "text.set_color", Description: "Set text color", Params: []Param{nodeID, color}},
		{Command: "text.set_align", Description: "Set text alignment", Params: []Param{nodeID, {Name: "textAlign", Type: "string"}}},
		{Command: "text.set_spacing", Description: "Set text spacing", Params: []Param{nodeID, {Name: "letterSpacing", Type: "number"}, {Name: "lineHeight", Type: "number"}, {Name: "paragraphSpacing", Type: "number"}}},
		{Command: "text.set_line_height", Description: "Set line height", Params: []Param{nodeID, {Name: "lineHeight", Type: "number", Required: true}, {Name: "lineHeightUnit", Type: "string"}}},
		{Command: "text.set_letter_spacing", Description: "Set letter spacing", Params: []Param{nodeID, {Name: "letterSpacing", Type: "number", Required: true}, {Name: "letterSpacingUnit", Type: "string"}}},
		{Command: "text.set_decoration", Description: "Set text decoration", Params: []Param{nodeID, {Name: "textDecoration", Type: "string"}}},
		{Command: "text.set_case", Description: "Set text case", Params: []Param{nodeID, {Name: "textCase", Type: "string"}}},
		{Command: "text.set_paragraph_spacing", Description: "Set paragraph spacing", Params: []Param{nodeID, {Name: "paragraphSpacing", Type: "number", Required: true}}},
		{Command: "text.get_content", Description: "Get text content and aggregate properties. Mixed text properties are returned as 'mixed'; use text.get_segments for per-range formatting.", Params: []Param{nodeID}},
		{Command: "text.get_segments", Description: "Get styled text segments by character range. Use before replacing rich text and reapply styles with text.set_range_style afterward.", Params: []Param{nodeID}},
		{Command: "text.load_font", Description: "Load a font", Params: []Param{{Name: "fontFamily", Type: "string", Required: true}, {Name: "fontStyle", Type: "string"}}},
		{Command: "text.set_style_id", Description: "Apply a text style", Params: []Param{nodeID, {Name: "styleId", Type: "string", Required: true}}},
		{Command: "text.list_fonts", Description: "List available fonts", Params: []Param{{Name: "fontFamily", Type: "string"}}},
		{Command: "text.set_range_style", Description: "Apply styles to ranges", Params: []Param{nodeID, {Name: "ranges", Type: "array", Required: true}}},
		{Command: "text.set_opentype_features", Description: "Set OpenType features", Params: []Param{nodeID, {Name: "features", Type: "object", Required: true}}},
		{Command: "text.get_opentype_features", Description: "Get OpenType features", Params: []Param{nodeID}},
	} {
		Register(s)
	}
}
