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
	"set_multiple_text_contents": {Domain: "text", Action: "set_content"},

	// Layout
	"set_auto_layout":   {Domain: "layout", Action: "set_auto_layout"},
	"set_padding":       {Domain: "layout", Action: "set_padding"},
	"set_item_spacing":  {Domain: "layout", Action: "set_spacing"},
	"set_layout_align":  {Domain: "layout", Action: "set_alignment"},
	"set_layout_sizing": {Domain: "layout", Action: "set_sizing"},
	"set_layout_wrap":   {Domain: "layout", Action: "set_layout_wrap"},
	"set_constraints":   {Domain: "layout", Action: "set_constraints"},

	// Node
	"get_node_info":     {Domain: "node", Action: "get_info"},
	"get_node_tree":     {Domain: "node", Action: "get_tree"},
	"create_frame":      {Domain: "node", Action: "create_frame"},
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

	// AI-friendly aliases (common guesses that should work)
	"list_pages":            {Domain: "page", Action: "get_all"},
	"list_styles":           {Domain: "style", Action: "get_all"},
	"get_all_styles":        {Domain: "style", Action: "get_all"},
	"list_variables":        {Domain: "variable", Action: "get_all"},
	"get_all_variables":     {Domain: "variable", Action: "get_all"},
	"list_components":       {Domain: "component", Action: "get_local"},
	"get_components":        {Domain: "component", Action: "get_local"},
	"list_collections":      {Domain: "variable", Action: "get_collections"},
	"get_collections":       {Domain: "variable", Action: "get_collections"},
	"list_children":         {Domain: "node", Action: "get_tree"},
	"get_children":          {Domain: "node", Action: "get_tree"},
	"get_tree":              {Domain: "node", Action: "get_tree"},
	"export_image":          {Domain: "export", Action: "image"},
	"export_svg":            {Domain: "export", Action: "svg"},
	"export_pdf":            {Domain: "export", Action: "pdf"},
	"export_json":           {Domain: "export", Action: "json"},
	"export_node":           {Domain: "export", Action: "image"},
	"get_layout":            {Domain: "layout", Action: "set_auto_layout"},
	"get_effects":           {Domain: "effect", Action: "get_effects"},
	"create_union":          {Domain: "boolean", Action: "union"},
	"create_subtract":       {Domain: "boolean", Action: "subtract"},
	"create_intersect":      {Domain: "boolean", Action: "intersect"},
	"create_exclude":        {Domain: "boolean", Action: "exclude"},
	"list_local_components":  {Domain: "component", Action: "get_local"},
	"list_remote_components": {Domain: "component", Action: "get_remote"},
	"get_document_selection": {Domain: "document", Action: "get_selection"},
	"get_current_page":      {Domain: "page", Action: "get_current"},
	"remove_auto_layout":    {Domain: "layout", Action: "remove_auto_layout"},
}

func resolveCommandRoute(command string, params map[string]interface{}) (string, string, error) {
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

	route, ok := legacyCommandRoutes[command]
	if !ok {
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
