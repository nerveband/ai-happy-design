package ws

import (
	"fmt"
	"strings"
)

type commandRoute struct {
	Domain string
	Action string
}

var legacyCommandRoutes = map[string]commandRoute{
	// Paint
	"set_fill_color":          {Domain: "paint", Action: "set_solid"},
	"set_gradient":            {Domain: "paint", Action: "set_gradient"},
	"set_gradient_fill":       {Domain: "paint", Action: "set_gradient"},
	"set_image_fill":          {Domain: "paint", Action: "set_image"},
	"set_image_fill_url":      {Domain: "paint", Action: "set_image_url"},
	"set_image_fill_from_url": {Domain: "paint", Action: "set_image_url"},
	"add_fill":                {Domain: "paint", Action: "add_fill"},
	"remove_fill":             {Domain: "paint", Action: "remove_fill"},
	"get_fills":               {Domain: "paint", Action: "get_fills"},
	"set_stroke_color":        {Domain: "paint", Action: "set_stroke"},

	// Shape
	"create_rectangle": {Domain: "shape", Action: "create_rectangle"},
	"create_ellipse":   {Domain: "shape", Action: "create_ellipse"},
	"create_polygon":   {Domain: "shape", Action: "create_polygon"},
	"create_star":      {Domain: "shape", Action: "create_star"},
	"create_line":      {Domain: "shape", Action: "create_line"},
	"create_from_svg":  {Domain: "shape", Action: "create_from_svg"},
	"create_image":     {Domain: "shape", Action: "create_image"},
	"create_vector":    {Domain: "shape", Action: "create_vector"},

	// Text
	"create_text":                {Domain: "text", Action: "create"},
	"set_text_content":           {Domain: "text", Action: "set_content"},
	"set_font_name":              {Domain: "text", Action: "set_font"},
	"set_font_size":              {Domain: "text", Action: "set_size"},
	"set_font_weight":            {Domain: "text", Action: "set_weight"},
	"set_text_align":             {Domain: "text", Action: "set_align"},
	"set_text_spacing":           {Domain: "text", Action: "set_spacing"},
	"set_text_case":              {Domain: "text", Action: "set_case"},
	"set_text_decoration":        {Domain: "text", Action: "set_decoration"},
	"get_styled_text_segments":   {Domain: "text", Action: "get_segments"},
	"load_font_async":            {Domain: "text", Action: "load_font"},
	"set_text_style_id":          {Domain: "text", Action: "set_style_id"},
	"set_letter_spacing":         {Domain: "text", Action: "set_letter_spacing"},
	"set_line_height":            {Domain: "text", Action: "set_line_height"},
	"set_paragraph_spacing":      {Domain: "text", Action: "set_paragraph_spacing"},
	"list_fonts":                 {Domain: "text", Action: "list_fonts"},
	"list_available_fonts":       {Domain: "text", Action: "list_fonts"},
	"set_multiple_text_contents": {Domain: "text", Action: "set_content"},
	"set_opentype_features":      {Domain: "text", Action: "set_opentype_features"},
	"get_opentype_features":      {Domain: "text", Action: "get_opentype_features"},

	// Layout
	"set_auto_layout":   {Domain: "layout", Action: "set_auto_layout"},
	"set_padding":       {Domain: "layout", Action: "set_padding"},
	"set_item_spacing":  {Domain: "layout", Action: "set_spacing"},
	"set_layout_align":  {Domain: "layout", Action: "set_alignment"},
	"set_layout_sizing": {Domain: "layout", Action: "set_sizing"},
	"set_layout_wrap":   {Domain: "layout", Action: "set_layout_wrap"},
	"set_constraints":   {Domain: "layout", Action: "set_constraints"},
	"check_overlaps":    {Domain: "layout", Action: "check_overlaps"},
	"set_grid":          {Domain: "layout", Action: "set_grid"},
	"set_layout_grid":   {Domain: "layout", Action: "set_grid"},
	"get_layout_grids":  {Domain: "layout", Action: "get_grids"},
	"remove_layout_grids": {Domain: "layout", Action: "remove_grids"},

	// Node
	"get_node_info":     {Domain: "node", Action: "get_info"},
	"get_node_tree":     {Domain: "node", Action: "get_tree"},
	"create_frame":      {Domain: "node", Action: "create_frame"},
	"create_section":    {Domain: "node", Action: "create_section"},
	"move_node":         {Domain: "node", Action: "move"},
	"resize_node":       {Domain: "node", Action: "resize"},
	"rotate_node":       {Domain: "node", Action: "rotate"},
	"set_rotation":      {Domain: "node", Action: "rotate"},
	"set_opacity":       {Domain: "node", Action: "set_opacity"},
	"set_blend_mode":    {Domain: "node", Action: "set_blend_mode"},
	"set_visibility":    {Domain: "node", Action: "set_visibility"},
	"set_locked":        {Domain: "node", Action: "set_locked"},
	"rename_node":       {Domain: "node", Action: "rename"},
	"delete_node":       {Domain: "node", Action: "delete"},
	"clone_node":        {Domain: "node", Action: "clone"},
	"set_corner_radius": {Domain: "node", Action: "set_corner_radius"},

	// Layer
	"set_layer_order": {Domain: "layer", Action: "set_order"},
	"bring_forward":   {Domain: "layer", Action: "bring_forward"},
	"send_backward":   {Domain: "layer", Action: "send_backward"},
	"bring_to_front":  {Domain: "layer", Action: "bring_to_front"},
	"send_to_back":    {Domain: "layer", Action: "send_to_back"},
	"group_nodes":     {Domain: "layer", Action: "group"},
	"ungroup_nodes":   {Domain: "layer", Action: "ungroup"},
	"insert_child":    {Domain: "layer", Action: "move_to_parent"},

	// Component
	"create_component_from_node": {Domain: "component", Action: "create"},
	"create_component_instance":  {Domain: "component", Action: "create_instance"},
	"create_component_set":       {Domain: "component", Action: "create_set"},
	"get_local_components":       {Domain: "component", Action: "get_local"},
	"get_remote_components":      {Domain: "component", Action: "get_remote"},
	"get_overrides":              {Domain: "component", Action: "get_overrides"},
	"set_overrides":              {Domain: "component", Action: "set_overrides"},
	"get_property_definitions":   {Domain: "component", Action: "get_property_definitions"},
	"add_component_property":     {Domain: "component", Action: "add_property_definition"},
	"delete_component_property":  {Domain: "component", Action: "delete_property_definition"},

	// Style
	"create_paint_style":  {Domain: "style", Action: "create_paint"},
	"create_text_style":   {Domain: "style", Action: "create_text"},
	"create_effect_style": {Domain: "style", Action: "create_effect"},
	"apply_style":         {Domain: "style", Action: "apply"},
	"get_styles":          {Domain: "document", Action: "get_styles"},
	"remove_style":        {Domain: "style", Action: "remove"},

	// Variable
	"create_variable":            {Domain: "variable", Action: "create"},
	"get_variables":              {Domain: "variable", Action: "get_all"},
	"set_variable_value":         {Domain: "variable", Action: "set_value"},
	"bind_variable":              {Domain: "variable", Action: "bind"},
	"unbind_variable":            {Domain: "variable", Action: "unbind"},
	"create_variable_collection": {Domain: "variable", Action: "create_collection"},
	"resolve_variable":           {Domain: "variable", Action: "resolve_for_consumer"},
	"resolve_for_consumer":       {Domain: "variable", Action: "resolve_for_consumer"},
	"add_variable_mode":          {Domain: "variable", Action: "add_mode"},
	"rename_variable_mode":       {Domain: "variable", Action: "rename_mode"},
	"delete_variable_mode":       {Domain: "variable", Action: "delete_mode"},

	// Effect
	"set_effects":         {Domain: "effect", Action: "set_effects"},
	"add_shadow":          {Domain: "effect", Action: "add_shadow"},
	"add_blur":            {Domain: "effect", Action: "add_blur"},
	"set_effect_style_id": {Domain: "effect", Action: "apply_style"},

	// Boolean
	"boolean_union":     {Domain: "boolean", Action: "union"},
	"boolean_subtract":  {Domain: "boolean", Action: "subtract"},
	"boolean_intersect": {Domain: "boolean", Action: "intersect"},
	"boolean_exclude":   {Domain: "boolean", Action: "exclude"},
	"flatten_node":      {Domain: "boolean", Action: "flatten"},

	// Page
	"create_page":      {Domain: "page", Action: "create"},
	"delete_page":      {Domain: "page", Action: "delete"},
	"rename_page":      {Domain: "page", Action: "rename"},
	"duplicate_page":   {Domain: "page", Action: "duplicate"},
	"set_current_page": {Domain: "page", Action: "set_current"},
	"get_pages":        {Domain: "page", Action: "get_all"},

	// Document
	"get_document_info": {Domain: "document", Action: "get_info"},
	"get_selection":     {Domain: "document", Action: "get_selection"},
	"set_selection":     {Domain: "document", Action: "set_selection"},
	"scan_text_nodes":   {Domain: "document", Action: "scan_text"},
	"scan_by_type":      {Domain: "document", Action: "find_by_type"},
	"focus_node":        {Domain: "document", Action: "zoom_to"},
	"find_free_space":   {Domain: "document", Action: "find_free_space"},

	// AI-friendly aliases (common guesses that should work)
	"modify":                 {Domain: "node", Action: "modify"},
	"find":                   {Domain: "document", Action: "find_nodes"},
	"find_nodes":             {Domain: "document", Action: "find_nodes"},
	"search_nodes":           {Domain: "document", Action: "find_nodes"},
	"mask":                   {Domain: "node", Action: "set_mask"},
	"set_mask":               {Domain: "node", Action: "set_mask"},
	"noise":                  {Domain: "effect", Action: "add_noise"},
	"texture":                {Domain: "effect", Action: "add_texture"},
	"glass":                  {Domain: "effect", Action: "apply_glass"},
	"native_glass":           {Domain: "effect", Action: "add_glass"},
	"frame":                  {Domain: "node", Action: "create_frame"},
	"rect":                   {Domain: "shape", Action: "create_rectangle"},
	"ellipse":                {Domain: "shape", Action: "create_ellipse"},
	"line":                   {Domain: "shape", Action: "create_line"},
	"image":                  {Domain: "shape", Action: "create_image"},
	"vector":                 {Domain: "shape", Action: "create_vector"},
	"section":                {Domain: "node", Action: "create_section"},
	"grid":                   {Domain: "layout", Action: "set_grid"},
	"text":                   {Domain: "text", Action: "create"},
	"fill":                   {Domain: "paint", Action: "set_solid"},
	"stroke":                 {Domain: "paint", Action: "set_stroke"},
	"gradient":               {Domain: "paint", Action: "set_gradient"},
	"shadow":                 {Domain: "effect", Action: "add_shadow"},
	"blur":                   {Domain: "effect", Action: "add_blur"},
	"parent":                 {Domain: "layer", Action: "move_to_parent"},
	"autolayout":             {Domain: "layout", Action: "set_auto_layout"},
	"opacity":                {Domain: "node", Action: "set_opacity"},
	"nofill":                 {Domain: "paint", Action: "remove_fill"},
	"list_pages":             {Domain: "page", Action: "get_all"},
	"list_styles":            {Domain: "style", Action: "get_all"},
	"get_all_styles":         {Domain: "style", Action: "get_all"},
	"list_variables":         {Domain: "variable", Action: "get_all"},
	"get_all_variables":      {Domain: "variable", Action: "get_all"},
	"list_components":        {Domain: "component", Action: "get_local"},
	"get_components":         {Domain: "component", Action: "get_local"},
	"list_collections":       {Domain: "variable", Action: "get_collections"},
	"get_collections":        {Domain: "variable", Action: "get_collections"},
	"list_children":          {Domain: "node", Action: "get_tree"},
	"get_children":           {Domain: "node", Action: "get_tree"},
	"get_tree":               {Domain: "node", Action: "get_tree"},
	"export_image":           {Domain: "export", Action: "image"},
	"export_svg":             {Domain: "export", Action: "svg"},
	"export_pdf":             {Domain: "export", Action: "pdf"},
	"export_json":            {Domain: "export", Action: "json"},
	"batch_export":           {Domain: "export", Action: "batch_export"},
	"export_batch":           {Domain: "export", Action: "batch_export"},
	"export_node":            {Domain: "export", Action: "image"},
	"get_layout":             {Domain: "layout", Action: "set_auto_layout"},
	"get_effects":            {Domain: "effect", Action: "get_effects"},
	"create_union":           {Domain: "boolean", Action: "union"},
	"create_subtract":        {Domain: "boolean", Action: "subtract"},
	"create_intersect":       {Domain: "boolean", Action: "intersect"},
	"create_exclude":         {Domain: "boolean", Action: "exclude"},
	"list_local_components":  {Domain: "component", Action: "get_local"},
	"list_remote_components": {Domain: "component", Action: "get_remote"},
	"get_document_selection": {Domain: "document", Action: "get_selection"},
	"get_current_page":       {Domain: "page", Action: "get_current"},
	"remove_auto_layout":     {Domain: "layout", Action: "remove_auto_layout"},
}

// compoundAliases intercepts common LLM hallucinations for dotted commands
// (e.g. "document.list_pages") before they reach the wrong domain handler.
// These are checked BEFORE splitting on "." in resolveCommandRoute.
var compoundAliases = map[string]commandRoute{
	// Page operations — LLMs often guess document.* for these
	"document.list_pages":  {Domain: "page", Action: "get_all"},
	"document.get_pages":   {Domain: "page", Action: "get_all"},
	"document.pages":       {Domain: "page", Action: "get_all"},
	"document.create_page": {Domain: "page", Action: "create"},
	"document.delete_page": {Domain: "page", Action: "delete"},
	"document.rename_page": {Domain: "page", Action: "rename"},
	"document.switch_page": {Domain: "page", Action: "set_current"},
	"document.set_page":    {Domain: "page", Action: "set_current"},
	// Layer operations — LLMs often guess node.* for these
	"node.group":          {Domain: "layer", Action: "group"},
	"node.ungroup":        {Domain: "layer", Action: "ungroup"},
	"node.bring_to_front": {Domain: "layer", Action: "bring_to_front"},
	"node.send_to_back":   {Domain: "layer", Action: "send_to_back"},
	"node.bring_forward":  {Domain: "layer", Action: "bring_forward"},
	"node.send_backward":  {Domain: "layer", Action: "send_backward"},
	"node.move_to_parent": {Domain: "layer", Action: "move_to_parent"},
	"node.reorder":        {Domain: "layer", Action: "set_order"},
	// Design system — LLMs might try document.*
	"document.analyze":           {Domain: "design_system", Action: "analyze"},
	"document.get_design_system": {Domain: "design_system", Action: "analyze"},
	"document.design_system":     {Domain: "design_system", Action: "analyze"},
}

func resolveCommandRoute(command string, params map[string]interface{}) (string, string, error) {
	command = strings.TrimSpace(command)
	// Check compound aliases first — catches LLM hallucinations like "document.list_pages"
	if route, ok := compoundAliases[command]; ok {
		return route.Domain, route.Action, nil
	}

	if dot := strings.Index(command, "."); dot > 0 && dot < len(command)-1 {
		return command[:dot], command[dot+1:], nil
	}

	if command == "export_node_as_image" {
		format := strings.ToUpper(stringArg(params, "format"))
		switch format {
		case "SVG":
			return "export", "svg", nil
		case "PDF":
			return "export", "pdf", nil
		default:
			return "export", "image", nil
		}
	}

	// shape.create with type param → route to specific sub-command.
	// LLMs frequently guess shape.create {type:"RECTANGLE"} instead of shape.create_rectangle.
	if command == "shape.create" {
		shapeType := strings.ToUpper(stringArg(params, "type"))
		switch shapeType {
		case "RECTANGLE", "RECT":
			return "shape", "create_rectangle", nil
		case "ELLIPSE", "CIRCLE", "OVAL":
			return "shape", "create_ellipse", nil
		case "POLYGON":
			return "shape", "create_polygon", nil
		case "STAR":
			return "shape", "create_star", nil
		case "LINE":
			return "shape", "create_line", nil
		case "SVG":
			return "shape", "create_from_svg", nil
		case "IMAGE":
			return "shape", "create_image", nil
		case "VECTOR":
			return "shape", "create_vector", nil
		default:
			return "shape", "create_rectangle", nil // sensible default
		}
	}

	route, ok := legacyCommandRoutes[command]
	if !ok {
		if suggestion := suggestCommand(command); suggestion != "" {
			return "", "", fmt.Errorf("unknown command: %s. did you mean: %s", command, suggestion)
		}
		return "", "", fmt.Errorf("unknown command: %s", command)
	}
	return route.Domain, route.Action, nil
}

func stringArg(params map[string]interface{}, key string) string {
	if params == nil {
		return ""
	}
	if v, ok := params[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func suggestCommand(command string) string {
	query := normalizeToken(command)
	if query == "" {
		return ""
	}

	candidates := knownCommands()
	best := ""

	// First pass: prefix/contains matches on normalized tokens.
	for _, c := range candidates {
		n := normalizeToken(c)
		if n == query {
			return c
		}
		if strings.HasPrefix(n, query) || strings.HasPrefix(query, n) || strings.Contains(n, query) {
			if best == "" || len(c) < len(best) {
				best = c
			}
		}
	}
	if best != "" {
		return best
	}

	// Fallback: nearest small edit distance.
	bestDist := 4
	for _, c := range candidates {
		n := normalizeToken(c)
		d := levenshtein(query, n)
		if d < bestDist {
			bestDist = d
			best = c
		}
	}
	if bestDist <= 3 {
		return best
	}
	return ""
}

func knownCommands() []string {
	seen := make(map[string]struct{}, len(legacyCommandRoutes)+len(compoundAliases))
	out := make([]string, 0, len(legacyCommandRoutes)+len(compoundAliases))
	for k := range legacyCommandRoutes {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, k)
	}
	for k := range compoundAliases {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, k)
	}
	return out
}

func normalizeToken(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func levenshtein(a, b string) int {
	if a == b {
		return 0
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	prev := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr := make([]int, len(b)+1)
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			curr[j] = minInt(del, minInt(ins, sub))
		}
		prev = curr
	}
	return prev[len(b)]
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
