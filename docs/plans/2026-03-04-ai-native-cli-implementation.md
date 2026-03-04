# AI-Native CLI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Transform ai-happy-design from a doc-driven tool into a schema-driven CLI where every command is validated, auto-corrected, and design-linted before execution — so LLMs get great design on the first call.

**Architecture:** Define canonical schemas per command in Go structs. Wire schema validation and design lint into the CLI batch pipeline between normalization and execution. Expose missing Figma API features in the plugin. Remove MCP entirely. Shrink the catalog from 92KB to ~25KB. Generate llms.txt from schemas.

**Tech Stack:** Go (CLI + schemas + validation + lint), TypeScript/ES6 (Figma plugin handlers), JSON Schema output

---

## Phase 1: Schema System + Structured Errors

### Task 1: Create Schema Types and Registry

**Files:**
- Create: `internal/schema/types.go`
- Create: `internal/schema/registry.go`

**Step 1: Write the Schema and Param type definitions**

Create `internal/schema/types.go`:

```go
package schema

// Param defines a single parameter for a command schema.
type Param struct {
	Name           string      `json:"name"`
	Type           string      `json:"type"`           // "string", "number", "boolean", "array", "object"
	Required       bool        `json:"required,omitempty"`
	Aliases        []string    `json:"aliases,omitempty"`
	Desc           string      `json:"description"`
	Enum           []string    `json:"enum,omitempty"`
	Min            *float64    `json:"minimum,omitempty"`
	Max            *float64    `json:"maximum,omitempty"`
	Default        interface{} `json:"default,omitempty"`
	Pattern        string      `json:"pattern,omitempty"`
	SemanticTokens bool        `json:"semanticTokens,omitempty"` // allows token names like "hero", "body"
	AutoFix        string      `json:"autoFix,omitempty"`        // e.g. "lineHeightUnit:PERCENT"
}

// Schema defines the full schema for a command.action pair.
type Schema struct {
	Command     string  `json:"command"`
	Aliases     []string `json:"aliases,omitempty"`
	Description string  `json:"description"`
	Params      []Param `json:"params"`
}

// Ptr returns a pointer to a float64 (helper for Min/Max).
func Ptr(v float64) *float64 { return &v }
```

**Step 2: Create the registry**

Create `internal/schema/registry.go`:

```go
package schema

import "strings"

// All registered schemas. Populated by init() in each schema file.
var All []Schema

// Register adds a schema to the global registry.
func Register(s Schema) { All = append(All, s) }

// Lookup finds a schema by command name or alias. Returns nil if not found.
func Lookup(command string) *Schema {
	cmd := strings.ToLower(command)
	for i := range All {
		if strings.ToLower(All[i].Command) == cmd {
			return &All[i]
		}
		for _, alias := range All[i].Aliases {
			if strings.ToLower(alias) == cmd {
				return &All[i]
			}
		}
	}
	return nil
}

// LookupParam finds a param by name or alias within a schema.
func LookupParam(s *Schema, name string) *Param {
	lower := strings.ToLower(name)
	for i := range s.Params {
		if strings.ToLower(s.Params[i].Name) == lower {
			return &s.Params[i]
		}
		for _, alias := range s.Params[i].Aliases {
			if strings.ToLower(alias) == lower {
				return &s.Params[i]
			}
		}
	}
	return nil
}

// Commands returns all registered command names.
func Commands() []string {
	out := make([]string, len(All))
	for i, s := range All {
		out[i] = s.Command
	}
	return out
}
```

**Step 3: Verify compilation**

Run: `cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design" && go build ./internal/schema/...`
Expected: PASS (no errors)

**Step 4: Commit**

```bash
git add internal/schema/types.go internal/schema/registry.go
git commit -m "feat: add schema type system and registry"
```

---

### Task 2: Define Core Command Schemas (15 Most-Used)

**Files:**
- Create: `internal/schema/node_schemas.go`
- Create: `internal/schema/text_schemas.go`
- Create: `internal/schema/paint_schemas.go`
- Create: `internal/schema/shape_schemas.go`
- Create: `internal/schema/layout_schemas.go`
- Create: `internal/schema/effect_schemas.go`

**Step 1: Define node schemas**

Create `internal/schema/node_schemas.go`:

```go
package schema

func init() {
	Register(Schema{
		Command:     "node.create_frame",
		Aliases:     []string{"frame"},
		Description: "Create a frame (Figma's equivalent of a div)",
		Params: []Param{
			{Name: "name", Type: "string", Desc: "Semantic name for the frame (never leave as default)"},
			{Name: "parentId", Type: "string", Aliases: []string{"pid"}, Desc: "Parent node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "x", Type: "number", Desc: "X position"},
			{Name: "y", Type: "number", Desc: "Y position"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Desc: "Frame width", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Desc: "Frame height", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "color", Type: "string", Aliases: []string{"bg"}, Desc: "Fill color hex (NOT fillColor)", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Desc: "Corner radius", Min: Ptr(0), Max: Ptr(500)},
			{Name: "opacity", Type: "number", Desc: "Opacity", Min: Ptr(0), Max: Ptr(1)},
			{Name: "layoutMode", Type: "string", Desc: "Auto-layout direction", Enum: []string{"HORIZONTAL", "VERTICAL"}},
			{Name: "itemSpacing", Type: "number", Desc: "Spacing between auto-layout children", Min: Ptr(0), Max: Ptr(500)},
			{Name: "padding", Type: "number", Desc: "Uniform padding for all sides", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingTop", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingRight", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingBottom", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingLeft", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "primaryAxisAlignItems", Type: "string", Enum: []string{"MIN", "CENTER", "MAX", "SPACE_BETWEEN"}},
			{Name: "counterAxisAlignItems", Type: "string", Enum: []string{"MIN", "CENTER", "MAX", "BASELINE"}},
			{Name: "layoutSizingHorizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutSizingVertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutWrap", Type: "string", Enum: []string{"WRAP", "NO_WRAP"}},
			{Name: "clipsContent", Type: "boolean", Desc: "Whether children clip to frame bounds"},
			{Name: "structural", Type: "boolean", Desc: "Structural frame: removes default fill, enables clipping. Use for wrappers and auto-layout containers."},
			{Name: "stroke", Type: "string", Desc: "Stroke color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeWeight", Type: "number", Aliases: []string{"sw"}, Min: Ptr(0), Max: Ptr(100)},
			{Name: "minWidth", Type: "number", Desc: "Minimum width (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxWidth", Type: "number", Desc: "Maximum width (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "minHeight", Type: "number", Desc: "Minimum height (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxHeight", Type: "number", Desc: "Maximum height (auto-layout children only)", Min: Ptr(0), Max: Ptr(10000)},
		},
	})

	Register(Schema{
		Command:     "node.modify",
		Aliases:     []string{"modify"},
		Description: "Modify any property of an existing node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Desc: "Target node ID", Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "color", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "visible", Type: "boolean"},
			{Name: "name", Type: "string"},
			{Name: "rotation", Type: "number", Min: Ptr(0), Max: Ptr(360)},
			{Name: "text", Type: "string"},
			{Name: "fontSize", Type: "number", Aliases: []string{"sz"}, Min: Ptr(4), Max: Ptr(500), SemanticTokens: true},
			{Name: "fontFamily", Type: "string", Aliases: []string{"ff"}},
			{Name: "isMask", Type: "boolean"},
			{Name: "layoutSizingHorizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutSizingVertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "blendMode", Type: "string", Enum: []string{"NORMAL", "MULTIPLY", "SCREEN", "OVERLAY", "DARKEN", "LIGHTEN", "COLOR_DODGE", "COLOR_BURN", "HARD_LIGHT", "SOFT_LIGHT", "DIFFERENCE", "EXCLUSION", "HUE", "SATURATION", "COLOR", "LUMINOSITY"}},
			{Name: "minWidth", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxWidth", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "minHeight", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "maxHeight", Type: "number", Min: Ptr(0), Max: Ptr(10000)},
			{Name: "constrainProportions", Type: "boolean"},
		},
	})

	Register(Schema{
		Command:     "node.move",
		Description: "Move a node to new coordinates",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "x", Type: "number", Required: true}, {Name: "y", Type: "number", Required: true},
		},
	})

	Register(Schema{
		Command:     "node.resize",
		Description: "Resize a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "width", Type: "number", Required: true, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Required: true, Min: Ptr(1), Max: Ptr(10000)},
		},
	})

	Register(Schema{
		Command:     "node.delete",
		Description: "Delete a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "node.get_info",
		Description: "Get information about a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
		},
	})

	Register(Schema{
		Command:     "node.get_tree",
		Description: "Get the node tree hierarchy",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "depth", Type: "number", Desc: "Max depth to traverse", Min: Ptr(1), Max: Ptr(20)},
			{Name: "compact", Type: "boolean", Desc: "Return flat array instead of nested tree (3-5x fewer tokens)"},
		},
	})
}
```

**Step 2: Define text schemas**

Create `internal/schema/text_schemas.go`:

```go
package schema

func init() {
	Register(Schema{
		Command:     "text.create",
		Aliases:     []string{"text"},
		Description: "Create a text node",
		Params: []Param{
			{Name: "text", Type: "string", Required: true, Aliases: []string{"content"}, Desc: "Text content to display"},
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string", Desc: "Semantic name for the text node"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "fontSize", Type: "number", Aliases: []string{"sz"}, Desc: "Font size in pixels", Min: Ptr(4), Max: Ptr(500), Default: 16.0, SemanticTokens: true},
			{Name: "fontFamily", Type: "string", Aliases: []string{"ff"}, Desc: "Font family name", Default: "Inter"},
			{Name: "fontStyle", Type: "string", Aliases: []string{"fs"}, Desc: "Font style", Default: "Regular",
				Enum: []string{"Thin", "ExtraLight", "Light", "Regular", "Medium", "SemiBold", "Bold", "ExtraBold", "Black", "Italic", "Bold Italic"}},
			{Name: "color", Type: "string", Desc: "Text color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`, Default: "#000000"},
			{Name: "lineHeight", Type: "number", Aliases: []string{"lh"}, Desc: "Line height percentage (e.g. 150 = 150%%)", Min: Ptr(50), Max: Ptr(300), AutoFix: "lineHeightUnit:PERCENT"},
			{Name: "lineHeightUnit", Type: "string", Enum: []string{"PIXELS", "PERCENT", "AUTO"}},
			{Name: "letterSpacing", Type: "number", Aliases: []string{"ls"}, Min: Ptr(-5), Max: Ptr(20)},
			{Name: "width", Type: "number", Desc: "Text box width (enables wrapping)", Min: Ptr(1), Max: Ptr(10000)},
			{Name: "textAlign", Type: "string", Aliases: []string{"textAlignHorizontal"}, Enum: []string{"LEFT", "CENTER", "RIGHT", "JUSTIFIED"}},
			{Name: "textCase", Type: "string", Enum: []string{"ORIGINAL", "UPPER", "LOWER", "TITLE"}},
			{Name: "textDecoration", Type: "string", Enum: []string{"NONE", "UNDERLINE", "STRIKETHROUGH"}},
			{Name: "maxLines", Type: "number", Desc: "Maximum lines before truncation with ellipsis", Min: Ptr(1), Max: Ptr(100)},
			{Name: "layoutSizingHorizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "layoutSizingVertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
		},
	})

	Register(Schema{
		Command:     "text.set_content",
		Description: "Set text content of an existing text node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "content", Type: "string", Required: true, Aliases: []string{"text"}},
		},
	})
}
```

**Step 3: Define paint schemas**

Create `internal/schema/paint_schemas.go`:

```go
package schema

func init() {
	Register(Schema{
		Command:     "paint.set_solid",
		Aliases:     []string{"fill"},
		Description: "Set a solid fill color on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Required: true, Desc: "Hex color (#RRGGBB or #RGB)", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "opacity", Type: "number", Desc: "Opacity", Min: Ptr(0), Max: Ptr(1), Default: 1.0},
		},
	})

	Register(Schema{
		Command:     "paint.set_gradient",
		Aliases:     []string{"gradient"},
		Description: "Set a gradient fill on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "gradientType", Type: "string", Enum: []string{"LINEAR", "RADIAL", "ANGULAR", "DIAMOND"}, Default: "LINEAR"},
			{Name: "stops", Type: "string", Required: true, Desc: "JSON array of gradient stops: [{position:0, color:\"#FF0000\"}, {position:1, color:\"#0000FF\"}]"},
			{Name: "angle", Type: "number", Desc: "Gradient angle in degrees (0=top, 90=right)", Min: Ptr(0), Max: Ptr(360)},
		},
	})

	Register(Schema{
		Command:     "paint.set_stroke",
		Aliases:     []string{"stroke"},
		Description: "Set stroke on a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Required: true, Desc: "Stroke color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "strokeWeight", Type: "number", Aliases: []string{"sw"}, Min: Ptr(0), Max: Ptr(100), Default: 1.0},
			{Name: "strokeAlign", Type: "string", Enum: []string{"INSIDE", "CENTER", "OUTSIDE"}, Default: "INSIDE"},
			{Name: "dashPattern", Type: "array", Desc: "Dash pattern [dash, gap] e.g. [10, 5]"},
			{Name: "strokeCap", Type: "string", Enum: []string{"NONE", "ROUND", "SQUARE", "ARROW_LINES", "ARROW_EQUILATERAL"}},
			{Name: "strokeJoin", Type: "string", Enum: []string{"MITER", "BEVEL", "ROUND"}},
		},
	})

	Register(Schema{
		Command:     "paint.set_image",
		Description: "Set an image fill on a node from base64 data",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "imageData", Type: "string", Required: true, Desc: "Base64-encoded image data or data URL"},
			{Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}, Default: "FILL"},
		},
	})
}
```

**Step 4: Define shape schemas**

Create `internal/schema/shape_schemas.go`:

```go
package schema

func init() {
	Register(Schema{
		Command:     "shape.create_rectangle",
		Aliases:     []string{"rect"},
		Description: "Create a rectangle",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Aliases: []string{"bg"}, Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "opacity", Type: "number", Min: Ptr(0), Max: Ptr(1)},
			{Name: "constrainProportions", Type: "boolean", Desc: "Lock aspect ratio"},
		},
	})

	Register(Schema{
		Command:     "shape.create_ellipse",
		Description: "Create an ellipse or circle",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000), Default: 100.0},
			{Name: "color", Type: "string", Pattern: `^#[0-9A-Fa-f]{3,8}$`},
			{Name: "constrainProportions", Type: "boolean"},
		},
	})

	Register(Schema{
		Command:     "shape.create_image",
		Description: "Create a rectangle with an image fill (one-step convenience)",
		Params: []Param{
			{Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"}, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "name", Type: "string"},
			{Name: "x", Type: "number"}, {Name: "y", Type: "number"},
			{Name: "width", Type: "number", Aliases: []string{"w"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "height", Type: "number", Aliases: []string{"h"}, Min: Ptr(1), Max: Ptr(10000)},
			{Name: "imageData", Type: "string", Required: true, Desc: "Base64 image data, data URL, file path, or HTTP(S) URL"},
			{Name: "scaleMode", Type: "string", Enum: []string{"FILL", "FIT", "TILE", "CROP"}, Default: "FILL"},
			{Name: "cornerRadius", Type: "number", Aliases: []string{"r"}, Min: Ptr(0), Max: Ptr(500)},
			{Name: "constrainProportions", Type: "boolean"},
		},
	})
}
```

**Step 5: Define layout schemas**

Create `internal/schema/layout_schemas.go`:

```go
package schema

func init() {
	Register(Schema{
		Command:     "layout.set_auto_layout",
		Description: "Set auto-layout properties on a frame",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "direction", Type: "string", Required: true, Enum: []string{"HORIZONTAL", "VERTICAL", "NONE"}, Desc: "HORIZONTAL=row, VERTICAL=column, NONE=remove auto-layout"},
			{Name: "itemSpacing", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "padding", Type: "number", Min: Ptr(0), Max: Ptr(500)},
			{Name: "paddingTop", Type: "number", Min: Ptr(0)}, {Name: "paddingRight", Type: "number", Min: Ptr(0)},
			{Name: "paddingBottom", Type: "number", Min: Ptr(0)}, {Name: "paddingLeft", Type: "number", Min: Ptr(0)},
			{Name: "primaryAxisAlignItems", Type: "string", Aliases: []string{"primaryAxisAlign"}, Enum: []string{"MIN", "CENTER", "MAX", "SPACE_BETWEEN"}},
			{Name: "counterAxisAlignItems", Type: "string", Aliases: []string{"counterAxisAlign"}, Enum: []string{"MIN", "CENTER", "MAX", "BASELINE"}},
			{Name: "layoutWrap", Type: "string", Enum: []string{"WRAP", "NO_WRAP"}},
		},
	})

	Register(Schema{
		Command:     "layout.set_sizing",
		Description: "Set sizing behavior of a node within auto-layout",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "horizontal", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
			{Name: "vertical", Type: "string", Enum: []string{"FIXED", "HUG", "FILL"}},
		},
	})
}
```

**Step 6: Define effect schemas**

Create `internal/schema/effect_schemas.go`:

```go
package schema

func init() {
	Register(Schema{
		Command:     "effect.add_shadow",
		Aliases:     []string{"shadow"},
		Description: "Add a drop shadow to a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "color", Type: "string", Desc: "Shadow color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`, Default: "#00000040"},
			{Name: "offsetX", Type: "number", Desc: "Horizontal offset", Default: 0.0},
			{Name: "offsetY", Type: "number", Desc: "Vertical offset", Default: 4.0},
			{Name: "radius", Type: "number", Desc: "Blur radius", Min: Ptr(0), Max: Ptr(200), Default: 4.0},
			{Name: "spread", Type: "number", Desc: "Spread distance", Min: Ptr(-100), Max: Ptr(100), Default: 0.0},
			{Name: "type", Type: "string", Enum: []string{"DROP_SHADOW", "INNER_SHADOW"}, Default: "DROP_SHADOW"},
		},
	})

	Register(Schema{
		Command:     "effect.add_blur",
		Aliases:     []string{"blur"},
		Description: "Add blur to a node",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "radius", Type: "number", Min: Ptr(0), Max: Ptr(200), Default: 10.0},
			{Name: "type", Type: "string", Enum: []string{"LAYER_BLUR", "BACKGROUND_BLUR"}, Default: "LAYER_BLUR"},
		},
	})

	Register(Schema{
		Command:     "effect.apply_glass",
		Aliases:     []string{"glass"},
		Description: "Apply glass morphism effect (background blur + semi-transparent fill + stroke)",
		Params: []Param{
			{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`},
			{Name: "intensity", Type: "string", Enum: []string{"light", "medium", "heavy"}, Default: "medium"},
		},
	})
}
```

**Step 7: Verify compilation**

Run: `go build ./internal/schema/...`
Expected: PASS

**Step 8: Commit**

```bash
git add internal/schema/*_schemas.go
git commit -m "feat: define schemas for 15 core commands"
```

---

### Task 3: Build Schema Validator

**Files:**
- Create: `internal/validate/validator.go`
- Create: `internal/validate/colors.go`
- Create: `internal/validate/fuzzy.go`
- Create: `internal/validate/validator_test.go`

**Step 1: Write failing tests**

Create `internal/validate/validator_test.go`:

```go
package validate

import (
	"testing"

	_ "github.com/nerveband/ai-happy-design/internal/schema" // register schemas
)

func TestValidParam(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "fontSize": 48.0,
		}},
	}
	result := ValidateBatch(ops)
	if len(result.Errors) > 0 {
		t.Fatalf("expected no errors, got %v", result.Errors)
	}
}

func TestBelowMin(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "fontSize": -10.0,
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "BELOW_MIN" && w.Param == "fontSize" {
			found = true
			if w.Fix != 4.0 {
				t.Errorf("expected fix=4, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected BELOW_MIN warning for fontSize")
	}
}

func TestEnumFuzzy(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "fontStyle": "bold",
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "ENUM_INVALID" && w.Param == "fontStyle" {
			found = true
			if w.Fix != "Bold" {
				t.Errorf("expected fix=Bold, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected ENUM_INVALID warning with fuzzy fix")
	}
}

func TestNamedColor(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "paint.set_solid", "params": map[string]interface{}{
			"nodeId": "1:23", "color": "red",
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "PATTERN_MISMATCH" && w.Param == "color" {
			found = true
			if w.Fix != "#FF0000" {
				t.Errorf("expected fix=#FF0000, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected PATTERN_MISMATCH warning with color fix")
	}
}

func TestUnknownCommand(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.creat", "params": map[string]interface{}{}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "UNKNOWN_COMMAND" {
			found = true
			if w.Fix != "text.create" {
				t.Errorf("expected fix=text.create, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected UNKNOWN_COMMAND warning")
	}
}

func TestDependencyAutoFix(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "text.create", "params": map[string]interface{}{
			"text": "Hello", "parentId": "1:23", "lineHeight": 150.0,
		}},
	}
	result := ValidateBatch(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "DEPENDENCY_MISSING" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected DEPENDENCY_MISSING warning for lineHeightUnit")
	}
	// Check it was auto-applied
	params := ops[0]["params"].(map[string]interface{})
	if params["lineHeightUnit"] != "PERCENT" {
		t.Errorf("expected lineHeightUnit auto-set to PERCENT, got %v", params["lineHeightUnit"])
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/validate/... -v`
Expected: FAIL (package doesn't exist yet)

**Step 3: Create colors.go (named CSS color lookup)**

Create `internal/validate/colors.go`:

```go
package validate

import "strings"

// cssColors maps CSS named colors to hex values.
var cssColors = map[string]string{
	"black":       "#000000", "white":       "#FFFFFF", "red":         "#FF0000",
	"green":       "#008000", "blue":        "#0000FF", "yellow":      "#FFFF00",
	"cyan":        "#00FFFF", "magenta":     "#FF00FF", "silver":      "#C0C0C0",
	"gray":        "#808080", "grey":        "#808080", "maroon":      "#800000",
	"olive":       "#808000", "lime":        "#00FF00", "aqua":        "#00FFFF",
	"teal":        "#008080", "navy":        "#000080", "fuchsia":     "#FF00FF",
	"purple":      "#800080", "orange":      "#FFA500", "pink":        "#FFC0CB",
	"brown":       "#A52A2A", "coral":       "#FF7F50", "crimson":     "#DC143C",
	"gold":        "#FFD700", "indigo":      "#4B0082", "ivory":       "#FFFFF0",
	"khaki":       "#F0E68C", "lavender":    "#E6E6FA", "linen":       "#FAF0E6",
	"orchid":      "#DA70D6", "peru":        "#CD853F", "plum":        "#DDA0DD",
	"salmon":      "#FA8072", "sienna":      "#A0522D", "tan":         "#D2B48C",
	"thistle":     "#D8BFD8", "tomato":      "#FF6347", "turquoise":   "#40E0D0",
	"violet":      "#EE82EE", "wheat":       "#F5DEB3", "transparent": "#00000000",
	// Add more as needed — these cover the most common LLM outputs
}

// ResolveNamedColor converts a CSS named color to hex. Returns "" if not found.
func ResolveNamedColor(name string) string {
	return cssColors[strings.ToLower(strings.TrimSpace(name))]
}
```

**Step 4: Create fuzzy.go (Levenshtein + case-insensitive matching)**

Create `internal/validate/fuzzy.go`:

```go
package validate

import "strings"

// FuzzyMatchEnum finds the closest enum value for a given input.
// Returns the match and true if found, or "" and false.
func FuzzyMatchEnum(input string, enum []string) (string, bool) {
	lower := strings.ToLower(strings.TrimSpace(input))

	// Exact case-insensitive match first
	for _, v := range enum {
		if strings.ToLower(v) == lower {
			return v, true
		}
	}

	// Strip hyphens/underscores/spaces and retry
	normalized := strings.NewReplacer("-", "", "_", "", " ", "").Replace(lower)
	for _, v := range enum {
		norm := strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(v))
		if norm == normalized {
			return v, true
		}
	}

	// Levenshtein distance <= 3
	best := ""
	bestDist := 4
	for _, v := range enum {
		d := levenshtein(lower, strings.ToLower(v))
		if d < bestDist {
			bestDist = d
			best = v
		}
	}
	if bestDist <= 3 {
		return best, true
	}
	return "", false
}

// FuzzyMatchCommand finds the closest command name for a given input.
func FuzzyMatchCommand(input string, commands []string) (string, bool) {
	lower := strings.ToLower(input)

	// Exact match
	for _, c := range commands {
		if strings.ToLower(c) == lower {
			return c, true
		}
	}

	// Prefix/contains
	for _, c := range commands {
		cl := strings.ToLower(c)
		if strings.HasPrefix(cl, lower) || strings.Contains(cl, lower) {
			return c, true
		}
	}

	// Levenshtein
	best := ""
	bestDist := 4
	for _, c := range commands {
		d := levenshtein(lower, strings.ToLower(c))
		if d < bestDist {
			bestDist = d
			best = c
		}
	}
	if bestDist <= 3 {
		return best, true
	}
	return "", false
}

func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 { return lb }
	if lb == 0 { return la }
	prev := make([]int, lb+1)
	curr := make([]int, lb+1)
	for j := range prev { prev[j] = j }
	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] { cost = 0 }
			ins := curr[j-1] + 1
			del := prev[j] + 1
			sub := prev[j-1] + cost
			m := ins
			if del < m { m = del }
			if sub < m { m = sub }
			curr[j] = m
		}
		prev, curr = curr, prev
	}
	return prev[lb]
}
```

**Step 5: Create validator.go (main validation engine)**

Create `internal/validate/validator.go`:

```go
package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nerveband/ai-happy-design/internal/schema"
)

// Issue represents a single validation warning or error.
type Issue struct {
	Step     int         `json:"step"`
	Name     string      `json:"name,omitempty"`
	Phase    string      `json:"phase"`
	Code     string      `json:"code"`
	Param    string      `json:"param,omitempty"`
	Message  string      `json:"message"`
	Got      interface{} `json:"got,omitempty"`
	Expected interface{} `json:"expected,omitempty"`
	Fix      interface{} `json:"fix,omitempty"`
	Applied  bool        `json:"applied"`
}

// Result holds all validation results for a batch.
type Result struct {
	Warnings []Issue `json:"warnings,omitempty"`
	Errors   []Issue `json:"errors,omitempty"`
	Fixed    int     `json:"fixed"`
	Blocked  int     `json:"blocked"`
}

// ValidateBatch validates all operations against registered schemas.
// It mutates ops in-place when auto-fixes are applied (warn+fix mode).
func ValidateBatch(ops []map[string]interface{}) Result {
	var result Result
	commands := schema.Commands()

	for i, op := range ops {
		cmdRaw, _ := op["command"].(string)
		name, _ := op["name"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			params = map[string]interface{}{}
		}

		// Check command exists
		s := schema.Lookup(cmdRaw)
		if s == nil {
			fix, found := FuzzyMatchCommand(cmdRaw, commands)
			issue := Issue{
				Step: i, Name: name, Phase: "schema", Code: "UNKNOWN_COMMAND",
				Message: fmt.Sprintf("unknown command: %s", cmdRaw),
				Got:     cmdRaw,
			}
			if found {
				issue.Fix = fix
				issue.Message += fmt.Sprintf(". Did you mean: %s?", fix)
				issue.Applied = true
				op["command"] = fix
				s = schema.Lookup(fix)
				result.Fixed++
			} else {
				result.Blocked++
			}
			result.Warnings = append(result.Warnings, issue)
			if s == nil {
				continue
			}
		}

		// Validate each param against schema
		for key, val := range params {
			p := schema.LookupParam(s, key)
			if p == nil {
				// Unknown param — try fuzzy match
				bestMatch := ""
				bestDist := 4
				for _, sp := range s.Params {
					d := levenshtein(strings.ToLower(key), strings.ToLower(sp.Name))
					if d < bestDist {
						bestDist = d
						bestMatch = sp.Name
					}
				}
				issue := Issue{
					Step: i, Name: name, Phase: "schema", Code: "UNKNOWN_PARAM",
					Param: key, Got: key,
					Message: fmt.Sprintf("unknown param '%s' for %s", key, s.Command),
				}
				if bestDist <= 3 {
					issue.Fix = bestMatch
					issue.Message += fmt.Sprintf(". Did you mean: %s?", bestMatch)
				}
				result.Warnings = append(result.Warnings, issue)
				continue
			}

			// Type checking + bounds + enum + pattern
			issues := validateParam(i, name, p, val, params)
			for _, issue := range issues {
				if issue.Applied {
					result.Fixed++
				} else if issue.Code != "" && issue.Fix == nil {
					result.Blocked++
				}
				result.Warnings = append(result.Warnings, issue)
			}
		}

		// Check required params
		for _, p := range s.Params {
			if !p.Required {
				continue
			}
			if _, exists := params[p.Name]; exists {
				continue
			}
			// Check aliases too
			found := false
			for _, alias := range p.Aliases {
				if _, exists := params[alias]; exists {
					found = true
					break
				}
			}
			if !found {
				result.Warnings = append(result.Warnings, Issue{
					Step: i, Name: name, Phase: "schema", Code: "REQUIRED_MISSING",
					Param:   p.Name,
					Message: fmt.Sprintf("required param '%s' missing", p.Name),
				})
				result.Blocked++
			}
		}

		// Check auto-fix dependencies
		for _, p := range s.Params {
			if p.AutoFix == "" {
				continue
			}
			if _, exists := params[p.Name]; !exists {
				continue // param not set, skip dependency check
			}
			// Parse "key:value"
			parts := strings.SplitN(p.AutoFix, ":", 2)
			if len(parts) != 2 {
				continue
			}
			depKey, depVal := parts[0], parts[1]
			if _, exists := params[depKey]; !exists {
				params[depKey] = depVal
				result.Warnings = append(result.Warnings, Issue{
					Step: i, Name: name, Phase: "schema", Code: "DEPENDENCY_MISSING",
					Param:   depKey,
					Message: fmt.Sprintf("auto-set %s=%s (required when %s is set)", depKey, depVal, p.Name),
					Fix:     depVal,
					Applied: true,
				})
				result.Fixed++
			}
		}

		// Write back params (may have been mutated)
		op["params"] = params
	}
	return result
}

func validateParam(step int, name string, p *schema.Param, val interface{}, params map[string]interface{}) []Issue {
	var issues []Issue

	switch p.Type {
	case "number":
		num, ok := toFloat64(val)
		if !ok {
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "TYPE_MISMATCH",
				Param: p.Name, Got: val,
				Message: fmt.Sprintf("%s must be a number, got %T", p.Name, val),
			})
			return issues
		}

		if p.Min != nil && num < *p.Min {
			params[p.Name] = *p.Min
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "BELOW_MIN",
				Param: p.Name, Got: num,
				Expected: map[string]interface{}{"min": *p.Min, "max": p.Max},
				Fix:     *p.Min,
				Applied: true,
				Message: fmt.Sprintf("%s must be >= %.0f, got %.0f", p.Name, *p.Min, num),
			})
		} else if p.Max != nil && num > *p.Max {
			params[p.Name] = *p.Max
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "ABOVE_MAX",
				Param: p.Name, Got: num,
				Expected: map[string]interface{}{"min": p.Min, "max": *p.Max},
				Fix:     *p.Max,
				Applied: true,
				Message: fmt.Sprintf("%s must be <= %.0f, got %.0f", p.Name, *p.Max, num),
			})
		}

	case "string":
		str, ok := val.(string)
		if !ok {
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "TYPE_MISMATCH",
				Param: p.Name, Got: val,
				Message: fmt.Sprintf("%s must be a string, got %T", p.Name, val),
			})
			return issues
		}

		// Enum check with fuzzy matching
		if len(p.Enum) > 0 {
			found := false
			for _, e := range p.Enum {
				if e == str {
					found = true
					break
				}
			}
			if !found {
				fix, matched := FuzzyMatchEnum(str, p.Enum)
				issue := Issue{
					Step: step, Name: name, Phase: "schema", Code: "ENUM_INVALID",
					Param: p.Name, Got: str,
					Expected: map[string]interface{}{"enum": p.Enum},
					Message:  fmt.Sprintf("%s must be one of %v, got '%s'", p.Name, p.Enum, str),
				}
				if matched {
					issue.Fix = fix
					issue.Applied = true
					params[p.Name] = fix
				}
				issues = append(issues, issue)
			}
		}

		// Pattern check (hex colors, node IDs)
		if p.Pattern != "" {
			re, err := regexp.Compile(p.Pattern)
			if err == nil && !re.MatchString(str) {
				issue := Issue{
					Step: step, Name: name, Phase: "schema", Code: "PATTERN_MISMATCH",
					Param: p.Name, Got: str,
					Expected: map[string]interface{}{"pattern": p.Pattern},
					Message:  fmt.Sprintf("%s doesn't match expected pattern, got '%s'", p.Name, str),
				}
				// Try named color resolution
				if hex := ResolveNamedColor(str); hex != "" {
					issue.Fix = hex
					issue.Applied = true
					params[p.Name] = hex
				}
				issues = append(issues, issue)
			}
		}

	case "boolean":
		if _, ok := val.(bool); !ok {
			issues = append(issues, Issue{
				Step: step, Name: name, Phase: "schema", Code: "TYPE_MISMATCH",
				Param: p.Name, Got: val,
				Message: fmt.Sprintf("%s must be a boolean, got %T", p.Name, val),
			})
		}
	}
	return issues
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}
```

**Step 6: Run tests**

Run: `go test ./internal/validate/... -v`
Expected: ALL PASS

**Step 7: Commit**

```bash
git add internal/validate/
git commit -m "feat: schema-based batch validation with fuzzy matching and named colors"
```

---

### Task 4: Wire Schema Validation into CLI Batch Pipeline

**Files:**
- Modify: `cmd/ai-happy-design/main.go:486-525` (batch command, after loadBatchOperations)

**Step 1: Add import and validation call**

In `main.go`, after `loadBatchOperations` (line 486) and before the execution loop (line 521), add schema validation. Find the section between `ops, err := loadBatchOperations(...)` and `for i, op := range ops {`.

Insert after line 519 (after `preprocessBatchImageData`), before the execution loop at line 521:

```go
// Schema validation (warn+fix mode)
var schemaResult validate.Result
if !batchNoFix {
    validationOps := make([]map[string]interface{}, len(ops))
    for i, op := range ops {
        validationOps[i] = map[string]interface{}{
            "command": op.Command,
            "params":  op.Params,
            "name":    op.Name,
        }
    }
    schemaResult = validate.ValidateBatch(validationOps)
    // Write fixed values back to ops
    for i, vo := range validationOps {
        if cmd, ok := vo["command"].(string); ok {
            ops[i].Command = cmd
        }
        if p, ok := vo["params"].(map[string]interface{}); ok {
            ops[i].Params = p
        }
    }
    if schemaResult.Fixed > 0 || len(schemaResult.Warnings) > 0 {
        fmt.Fprintf(os.Stderr, "[schema] %d fixed, %d warnings, %d blocked\n",
            schemaResult.Fixed, len(schemaResult.Warnings), schemaResult.Blocked)
    }
    if batchStrictQuality && schemaResult.Blocked > 0 {
        out := map[string]interface{}{
            "ok": false,
            "preValidation": map[string]interface{}{
                "schema": map[string]interface{}{
                    "warnings": schemaResult.Warnings,
                    "errors":   schemaResult.Errors,
                    "fixed":    schemaResult.Fixed,
                    "blocked":  schemaResult.Blocked,
                },
            },
        }
        j, _ := json.MarshalIndent(out, "", "  ")
        fmt.Println(string(j))
        return fmt.Errorf("schema validation blocked %d issues", schemaResult.Blocked)
    }
}
```

Add the import: `"github.com/nerveband/ai-happy-design/internal/validate"`

Also include the `schemaResult` in the final batch output (in the response assembly section after line 640). Add to the `out` map:

```go
if schemaResult.Fixed > 0 || len(schemaResult.Warnings) > 0 {
    out["preValidation"] = map[string]interface{}{
        "schema": map[string]interface{}{
            "warnings": schemaResult.Warnings,
            "fixed":    schemaResult.Fixed,
            "blocked":  schemaResult.Blocked,
        },
    }
}
```

**Step 2: Verify build**

Run: `go build ./cmd/ai-happy-design/`
Expected: PASS

**Step 3: Manual test**

Run: `go run ./cmd/ai-happy-design/ batch '[{"command":"text.creat","params":{"text":"Hello","parentId":"1:23","fontSize":-10,"fontStyle":"bold","color":"red"}}]' --no-auto-relay 2>&1 | head -20`

Expected stderr: `[schema] N fixed, N warnings, 0 blocked`
Expected output includes `preValidation.schema.warnings` with fixes applied.

**Step 4: Commit**

```bash
git add cmd/ai-happy-design/main.go
git commit -m "feat: wire schema validation into CLI batch pipeline"
```

---

### Task 5: Add `schema` and `validate` CLI Commands

**Files:**
- Create: `cmd/ai-happy-design/schema_cmd.go`
- Create: `cmd/ai-happy-design/validate_cmd.go`

**Step 1: Create schema command**

Create `cmd/ai-happy-design/schema_cmd.go`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/nerveband/ai-happy-design/internal/schema"
	"github.com/spf13/cobra"
)

var schemaJSON bool
var schemaAllLLMSTxt bool

var schemaCmd = &cobra.Command{
	Use:   "schema [command.action]",
	Short: "Print parameter schema for a command",
	Long:  "Print the exact parameter schema for a command. LLMs can use this to generate valid JSON.\nUse --all --llms-txt to generate llms-full.txt for aihappydesign.com.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if schemaAllLLMSTxt {
			return printLLMSTxt()
		}

		if len(args) == 0 {
			// List all commands
			for _, s := range schema.All {
				aliases := ""
				if len(s.Aliases) > 0 {
					aliases = " (aliases: " + strings.Join(s.Aliases, ", ") + ")"
				}
				fmt.Printf("%-30s %s%s\n", s.Command, s.Description, aliases)
			}
			return nil
		}

		s := schema.Lookup(args[0])
		if s == nil {
			return fmt.Errorf("unknown command: %s. Run 'ai-happy-design schema' to list all commands.", args[0])
		}

		if schemaJSON {
			out, _ := json.MarshalIndent(s, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		// Human-readable table
		fmt.Printf("## %s\n\n%s\n", s.Command, s.Description)
		if len(s.Aliases) > 0 {
			fmt.Printf("Aliases: %s\n", strings.Join(s.Aliases, ", "))
		}
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "PARAM\tTYPE\tREQUIRED\tDEFAULT\tCONSTRAINTS\tDESCRIPTION\n")
		for _, p := range s.Params {
			req := ""
			if p.Required {
				req = "yes"
			}
			def := ""
			if p.Default != nil {
				def = fmt.Sprintf("%v", p.Default)
			}
			constraints := buildConstraints(p)
			desc := p.Desc
			if len(p.Aliases) > 0 {
				desc += " (alias: " + strings.Join(p.Aliases, "/") + ")"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", p.Name, p.Type, req, def, constraints, desc)
		}
		w.Flush()
		return nil
	},
}

func buildConstraints(p schema.Param) string {
	parts := []string{}
	if p.Min != nil && p.Max != nil {
		parts = append(parts, fmt.Sprintf("%.0f-%.0f", *p.Min, *p.Max))
	} else if p.Min != nil {
		parts = append(parts, fmt.Sprintf(">= %.0f", *p.Min))
	} else if p.Max != nil {
		parts = append(parts, fmt.Sprintf("<= %.0f", *p.Max))
	}
	if len(p.Enum) > 0 {
		parts = append(parts, strings.Join(p.Enum, "/"))
	}
	if p.Pattern != "" {
		parts = append(parts, "pattern: "+p.Pattern)
	}
	if p.SemanticTokens {
		parts = append(parts, "tokens: hero/title/heading/subheading/body/caption")
	}
	return strings.Join(parts, ", ")
}

func printLLMSTxt() error {
	fmt.Println("# AI Happy Design — Full Command Reference")
	fmt.Println()
	fmt.Println("Auto-generated from CLI schemas. Use `ai-happy-design schema <command> --json` for machine-readable format.")
	fmt.Println()

	for _, s := range schema.All {
		fmt.Printf("## %s\n\n", s.Command)
		fmt.Printf("%s\n\n", s.Description)
		if len(s.Aliases) > 0 {
			fmt.Printf("Aliases: `%s`\n\n", strings.Join(s.Aliases, "`, `"))
		}
		fmt.Println("### Parameters\n")
		fmt.Println("| Param | Type | Required | Default | Constraints | Description |")
		fmt.Println("|-------|------|----------|---------|-------------|-------------|")
		for _, p := range s.Params {
			req := ""
			if p.Required {
				req = "yes"
			}
			def := ""
			if p.Default != nil {
				def = fmt.Sprintf("%v", p.Default)
			}
			constraints := buildConstraints(p)
			desc := p.Desc
			if len(p.Aliases) > 0 {
				desc += " (alias: " + strings.Join(p.Aliases, "/") + ")"
			}
			fmt.Printf("| %s | %s | %s | %s | %s | %s |\n", p.Name, p.Type, req, def, constraints, desc)
		}
		fmt.Println()
	}
	return nil
}

func init() {
	schemaCmd.Flags().BoolVar(&schemaJSON, "json", false, "Output in JSON format")
	schemaCmd.Flags().BoolVar(&schemaAllLLMSTxt, "all", false, "Print all schemas (use with --llms-txt)")
	rootCmd.AddCommand(schemaCmd)
}
```

**Step 2: Create validate (dry-run) command**

Create `cmd/ai-happy-design/validate_cmd.go` — if one exists already, replace it. This runs schema validation + design lint without executing:

```go
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nerveband/ai-happy-design/internal/validate"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [operations-json | file.json]",
	Short: "Dry-run: validate batch JSON without executing (schema + design lint)",
	Long:  "Runs the full validation pipeline (schema checks, design lint, quality score) without sending any commands to Figma. Use this to check your batch JSON before execution.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ops, err := loadBatchOperations(
			getPositionalOrFlag(args, batchOperations),
			batchOperationsFile,
		)
		if err != nil {
			return err
		}
		if len(ops) == 0 {
			return fmt.Errorf("operations array is empty")
		}

		// Convert to validation format
		validationOps := make([]map[string]interface{}, len(ops))
		for i, op := range ops {
			validationOps[i] = map[string]interface{}{
				"command": op.Command,
				"params":  op.Params,
				"name":    op.Name,
			}
		}

		// Schema validation
		schemaResult := validate.ValidateBatch(validationOps)

		out := map[string]interface{}{
			"ok":    schemaResult.Blocked == 0,
			"total": len(ops),
			"preValidation": map[string]interface{}{
				"schema": map[string]interface{}{
					"warnings": schemaResult.Warnings,
					"fixed":    schemaResult.Fixed,
					"blocked":  schemaResult.Blocked,
				},
			},
		}

		j, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(j))

		if schemaResult.Blocked > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func getPositionalOrFlag(args []string, flag string) string {
	if len(args) > 0 {
		return args[0]
	}
	return flag
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
```

**Step 3: Verify build**

Run: `go build ./cmd/ai-happy-design/`
Expected: PASS

**Step 4: Test schema command**

Run: `go run ./cmd/ai-happy-design/ schema text.create`
Expected: Parameter table printed

Run: `go run ./cmd/ai-happy-design/ schema text.create --json`
Expected: JSON schema printed

**Step 5: Test validate command**

Run: `echo '[{"command":"text.create","params":{"text":"Hi","fontSize":-5,"color":"red"}}]' | go run ./cmd/ai-happy-design/ validate`
Expected: JSON with warnings showing BELOW_MIN and PATTERN_MISMATCH fixes

**Step 6: Commit**

```bash
git add cmd/ai-happy-design/schema_cmd.go cmd/ai-happy-design/validate_cmd.go
git commit -m "feat: add schema and validate CLI commands"
```

---

## Phase 2: Design Lint + Score

### Task 6: Build Design Lint Engine

**Files:**
- Create: `internal/designlint/lint.go`
- Create: `internal/designlint/contrast.go`
- Create: `internal/designlint/lint_test.go`

**Step 1: Write failing tests**

Create `internal/designlint/lint_test.go`:

```go
package designlint

import (
	"testing"
)

func TestTextTooSmall(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "node.create_frame", "params": map[string]interface{}{"width": 1080.0, "height": 1350.0, "name": "Slide"}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Hi", "parentId": "1:23", "fontSize": 14.0}},
	}
	result := Check(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "TEXT_TOO_SMALL" {
			found = true
			if w.Fix.(float64) < 30 {
				t.Errorf("expected fix >= caption tier, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected TEXT_TOO_SMALL warning")
	}
}

func TestContrastRatio(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "node.create_frame", "params": map[string]interface{}{"width": 1080.0, "height": 1350.0, "color": "#777777"}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Hi", "parentId": "1:23", "color": "#666666"}},
	}
	result := Check(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "LOW_CONTRAST" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected LOW_CONTRAST warning")
	}
}

func TestRadiusOverflow(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "shape.create_rectangle", "params": map[string]interface{}{
			"parentId": "1:23", "width": 40.0, "height": 40.0, "cornerRadius": 50.0,
		}},
	}
	result := Check(ops)
	found := false
	for _, w := range result.Warnings {
		if w.Code == "RADIUS_OVERFLOW" {
			found = true
			if w.Fix.(float64) != 20.0 {
				t.Errorf("expected fix=20, got %v", w.Fix)
			}
		}
	}
	if !found {
		t.Fatal("expected RADIUS_OVERFLOW warning")
	}
}

func TestScoreComputed(t *testing.T) {
	ops := []map[string]interface{}{
		{"command": "node.create_frame", "params": map[string]interface{}{"width": 1080.0, "height": 1350.0}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Title", "fontSize": 112.0}},
		{"command": "text.create", "params": map[string]interface{}{"text": "Body", "fontSize": 48.0}},
	}
	result := Check(ops)
	if result.Score.Overall == 0 {
		t.Fatal("expected non-zero design score")
	}
}
```

**Step 2: Create contrast.go (WCAG color contrast utilities)**

Create `internal/designlint/contrast.go`:

```go
package designlint

import (
	"math"
	"strconv"
	"strings"
)

// ContrastRatio computes WCAG 2.0 contrast ratio between two hex colors.
// Returns a value between 1 (no contrast) and 21 (max contrast).
func ContrastRatio(hex1, hex2 string) float64 {
	l1 := relativeLuminance(hex1)
	l2 := relativeLuminance(hex2)
	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)
	return (lighter + 0.05) / (darker + 0.05)
}

func relativeLuminance(hex string) float64 {
	r, g, b := parseHexRGB(hex)
	rr := linearize(float64(r) / 255.0)
	gg := linearize(float64(g) / 255.0)
	bb := linearize(float64(b) / 255.0)
	return 0.2126*rr + 0.7152*gg + 0.0722*bb
}

func linearize(v float64) float64 {
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func parseHexRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) < 6 {
		return 0, 0, 0
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 32)
	g, _ := strconv.ParseInt(hex[2:4], 16, 32)
	b, _ := strconv.ParseInt(hex[4:6], 16, 32)
	return int(r), int(g), int(b)
}

// AdjustForContrast returns a color adjusted to meet minimum contrast ratio against bg.
// It darkens or lightens the foreground color as needed.
func AdjustForContrast(fgHex, bgHex string, minRatio float64) string {
	ratio := ContrastRatio(fgHex, bgHex)
	if ratio >= minRatio {
		return fgHex
	}

	bgLum := relativeLuminance(bgHex)
	// Try darkening first (toward black), then lightening (toward white)
	if bgLum > 0.5 {
		return "#222222" // dark text on light bg
	}
	return "#FFFFFF" // light text on dark bg
}
```

**Step 3: Create lint.go (main lint engine)**

Create `internal/designlint/lint.go`:

```go
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
	if canvasW <= 0 || canvasH <= 0 {
		// No canvas detected — run structural checks only
		result.Score = Score{Readability: 10, Contrast: 10, Spacing: 10, Hierarchy: 10, Overall: 10}
		checkStructural(ops, &result)
		computeScore(&result, nil, canvasW)
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
	checkRadiusOverflow(ops, &result)
	checkStructural(ops, &result)
	computeScore(&result, textTokens, canvasW)
	return result
}

func detectCanvas(ops []map[string]interface{}) (float64, float64) {
	for _, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		// Root frame: has width/height, no parentId
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
	// Build a map of step index → background color (from frames)
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

	// Check text colors against nearest ancestor bg
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

		// Find nearest bg color (scan backwards)
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
		sidePadding = canvasW * 0.065 // fallback 6.5%
	}
	minPaddingRatio := 0.04

	for i, op := range ops {
		cmd, _ := op["command"].(string)
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			continue
		}
		name, _ := op["name"].(string)

		// Check root frame padding
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

func computeScore(result *Result, textTokens map[string]interface{}, canvasW float64) {
	// Start at 10, deduct per issue
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
		}
	}
	// Clamp to 0
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
```

**Step 4: Run tests**

Run: `go test ./internal/designlint/... -v`
Expected: ALL PASS

**Step 5: Commit**

```bash
git add internal/designlint/
git commit -m "feat: design lint engine with text sizing, contrast, spacing, and scoring"
```

---

### Task 7: Wire Design Lint into CLI Batch Pipeline

**Files:**
- Modify: `cmd/ai-happy-design/main.go` (batch command, after schema validation, before execution loop)

**Step 1: Add design lint call after schema validation**

In the batch command, after the schema validation block added in Task 4, and before the execution loop, add:

```go
// Design lint (pre-execution quality checks)
var lintResult designlint.Result
if batchLint && !batchNoLint {
    lintOps := make([]map[string]interface{}, len(ops))
    for i, op := range ops {
        lintOps[i] = map[string]interface{}{
            "command": op.Command,
            "params":  op.Params,
            "name":    op.Name,
        }
    }
    lintResult = designlint.Check(lintOps)
    // Write fixes back
    for i, lo := range lintOps {
        if p, ok := lo["params"].(map[string]interface{}); ok {
            ops[i].Params = p
        }
    }
    if lintResult.Fixed > 0 || len(lintResult.Warnings) > 0 {
        fmt.Fprintf(os.Stderr, "[design-lint] %d fixed, %d warnings, score: %.1f/10\n",
            lintResult.Fixed, len(lintResult.Warnings), lintResult.Score.Overall)
    }
    if batchStrictQuality && lintResult.Score.Overall < 7.0 {
        out := map[string]interface{}{
            "ok": false,
            "preValidation": map[string]interface{}{
                "designLint": map[string]interface{}{
                    "canvas":   lintResult.Canvas,
                    "warnings": lintResult.Warnings,
                    "fixed":    lintResult.Fixed,
                    "score":    lintResult.Score,
                },
            },
        }
        j, _ := json.MarshalIndent(out, "", "  ")
        fmt.Println(string(j))
        return fmt.Errorf("design quality score %.1f/10 below threshold 7.0", lintResult.Score.Overall)
    }
}
```

Add import: `"github.com/nerveband/ai-happy-design/internal/designlint"`

Also include lint result in the final batch output alongside schema result:

```go
if lintResult.Fixed > 0 || len(lintResult.Warnings) > 0 {
    preVal, _ := out["preValidation"].(map[string]interface{})
    if preVal == nil {
        preVal = map[string]interface{}{}
    }
    preVal["designLint"] = map[string]interface{}{
        "canvas":   lintResult.Canvas,
        "tokens":   lintResult.Tokens,
        "warnings": lintResult.Warnings,
        "fixed":    lintResult.Fixed,
        "score":    lintResult.Score,
    }
    out["preValidation"] = preVal
}
```

**Step 2: Verify build**

Run: `go build ./cmd/ai-happy-design/`
Expected: PASS

**Step 3: Commit**

```bash
git add cmd/ai-happy-design/main.go
git commit -m "feat: wire design lint with scoring into CLI batch pipeline"
```

---

## Phase 3: Plugin Feature Expansions

### Task 8: Add Min/Max Sizing, Clips Content, Structural Frame to Plugin

**Files:**
- Modify: `plugin/src/handlers/node.ts:141-181` (create_frame handler)

**Step 1: Add properties to create_frame**

In `node.ts`, inside the create_frame handler, after `clipsContent` (around line 141), add:

```typescript
// Min/max sizing (auto-layout children only)
if (params.minWidth != null) frame.minWidth = params.minWidth;
if (params.maxWidth != null) frame.maxWidth = params.maxWidth;
if (params.minHeight != null) frame.minHeight = params.minHeight;
if (params.maxHeight != null) frame.maxHeight = params.maxHeight;

// Constrain proportions
if (params.constrainProportions != null) frame.constrainProportions = params.constrainProportions;

// Structural frame shorthand: remove default fill, enable clipping
if (params.structural) {
    frame.fills = [];
    frame.clipsContent = true;
}
```

Also add these to the `modify` handler (search for the modify action section).

**Step 2: Build and verify plugin**

Run: `cd plugin && npm run check && cd ..`
Run: `grep -c '\?\.' plugin/dist/code.js` → must be 0
Run: `grep -c '\?\?' plugin/dist/code.js` → must be 0

**Step 3: Commit**

```bash
git add plugin/src/handlers/node.ts
cd plugin && npm run build && cd ..
git add plugin/dist/
git commit -m "feat: add min/max sizing, constrain proportions, structural frame to plugin"
```

---

### Task 9: Add Text Truncation to Plugin

**Files:**
- Modify: `plugin/src/handlers/text.ts:167-174` (create handler, after textCase)

**Step 1: Add maxLines and textTruncation**

In `text.ts`, inside the create handler, after textCase (around line 169):

```typescript
// Text truncation
if (params.maxLines != null) {
    textNode.maxLines = params.maxLines;
    textNode.textTruncation = 'ENDING';
}
if (params.textTruncation != null) {
    textNode.textTruncation = params.textTruncation;
}
```

**Step 2: Build and verify**

Run: `cd plugin && npm run check && cd ..`

**Step 3: Commit**

```bash
git add plugin/src/handlers/text.ts
cd plugin && npm run build && cd ..
git add plugin/dist/
git commit -m "feat: add text truncation (maxLines + ellipsis) to plugin"
```

---

### Task 10: Add Stroke Dash, Cap, Join to Plugin

**Files:**
- Modify: `plugin/src/handlers/paint.ts:273-274` (setStroke, after strokeAlign)

**Step 1: Add stroke properties**

In `paint.ts`, inside the setStroke function, after strokeAlign:

```typescript
// Dash pattern
if (params.dashPattern != null) {
    node.dashPattern = params.dashPattern;
}
// Stroke cap and join
if (params.strokeCap != null) {
    node.strokeCap = params.strokeCap;
}
if (params.strokeJoin != null) {
    node.strokeJoin = params.strokeJoin;
}
```

**Step 2: Build and verify**

Run: `cd plugin && npm run check && cd ..`

**Step 3: Commit**

```bash
git add plugin/src/handlers/paint.ts
cd plugin && npm run build && cd ..
git add plugin/dist/
git commit -m "feat: add dash pattern, stroke cap, stroke join to plugin"
```

---

### Task 11: Add Constrain Proportions to Shape Handlers

**Files:**
- Modify: `plugin/src/handlers/shape.ts:121-138` (createRectangle, createEllipse)

**Step 1: Add constrainProportions to shape creation**

In `shape.ts`, inside createRectangle (after cornerRadius setup) and createEllipse:

```typescript
if (params.constrainProportions != null) {
    node.constrainProportions = params.constrainProportions;
}
```

**Step 2: Build and verify**

Run: `cd plugin && npm run check && cd ..`

**Step 3: Commit**

```bash
git add plugin/src/handlers/shape.ts
cd plugin && npm run build && cd ..
git add plugin/dist/
git commit -m "feat: add constrain proportions to shape handlers"
```

---

## Phase 4: MCP Removal + Catalog Shrink + llms.txt

### Task 12: Remove MCP Server and Dependencies

**Files:**
- Delete: `internal/mcp/server.go`
- Delete: `internal/mcp/registry.go` (if exists)
- Modify: `cmd/ai-happy-design/main.go` — remove mcp command, register command, mcp imports
- Modify: `go.mod` — remove mcp-go dependency
- Modify: `internal/tools/bulk.go` — remove MCP registration, keep execution logic as internal function
- Modify: `internal/tools/design_tokens.go` — remove MCP registration, keep ComputeDesignTokens as exported function
- Modify: all `internal/tools/*.go` — remove `Register*Tool` functions that reference mcp-go, keep any utility functions

**Step 1: Delete MCP directory**

```bash
rm -rf internal/mcp/
```

**Step 2: Strip MCP registration from all tool files**

For each file in `internal/tools/` that has a `Register*Tool` function:
- Remove the `Register*Tool` function entirely
- Remove mcp-go imports
- Keep any utility/computation functions (like `ComputeDesignTokens`, `LLMCatalog`)

The tool files become pure logic — no MCP wiring.

**Step 3: Remove mcp and register commands from main.go**

Delete the `mcpCmd` variable and its handler. Delete `registerCmd` addition. Remove mcp-go imports.

**Step 4: Remove mcp-go from go.mod**

Run: `go mod tidy` (after removing all imports)

**Step 5: Verify build**

Run: `go build ./cmd/ai-happy-design/`
Expected: PASS

Run: `go test ./...`
Expected: ALL PASS

**Step 6: Commit**

```bash
git add -A
git commit -m "feat: remove MCP server, tool registration, and mcp-go dependency"
```

---

### Task 13: Shrink Catalog

**Files:**
- Modify: `internal/tools/catalog_llm.go`

**Step 1: Remove sections now encoded in schemas/lint**

Delete these sections from `catalog_llm.go`:
- Execution Rules (encoded in pipeline behavior)
- CSS Property Support (in cssnorm.go)
- Common Mistakes (every mistake is now a schema constraint)
- First-Pass Guardrails (design lint)
- Batch Observability (response format is self-documenting)

**Step 2: Keep creative guidance sections**

Keep:
- Design Thinking
- Visual Hierarchy
- Design Decisions (shadows, gradients, blur recipes)
- Layer Organization

Shrink:
- Workflow → one paragraph
- Design Quality Checklist → brief reminder (most enforced by lint)

**Step 3: Add `guide` CLI command**

Create `cmd/ai-happy-design/guide_cmd.go`:

```go
package main

import (
	"fmt"

	"github.com/nerveband/ai-happy-design/internal/tools"
	"github.com/spf13/cobra"
)

var guideSection string

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Print design intelligence guide for LLMs",
	RunE: func(cmd *cobra.Command, args []string) error {
		catalog := tools.LLMCatalog()
		fmt.Println(catalog)
		return nil
	},
}

func init() {
	guideCmd.Flags().StringVar(&guideSection, "section", "", "Specific section: hierarchy, decisions, thinking, organization")
	rootCmd.AddCommand(guideCmd)
}
```

**Step 4: Verify build and test**

Run: `go build ./cmd/ai-happy-design/`
Run: `go run ./cmd/ai-happy-design/ guide | wc -c`
Expected: ~25000-30000 bytes (down from ~92000)

**Step 5: Commit**

```bash
git add internal/tools/catalog_llm.go cmd/ai-happy-design/guide_cmd.go
git commit -m "feat: shrink catalog from 92KB to ~25KB, add guide CLI command"
```

---

### Task 14: Generate and Deploy llms.txt

**Files:**
- Create: `/Users/nerveband/Documents/GitHub/ai-happy-design-site/llms.txt` (hand-written)
- Generate: `/Users/nerveband/Documents/GitHub/ai-happy-design-site/llms-full.txt` (from schemas)

**Step 1: Write llms.txt**

Create the concise llms.txt at the site repo:

```markdown
# AI Happy Design

> CLI tool that gives LLMs full Figma canvas control via WebSocket relay to Figma plugin.

## Quick Start

1. Install: `brew install nerveband/tap/ai-happy-design` or `ai-happy-design upgrade`
2. Start relay: `ai-happy-design ws`
3. Open Figma plugin, connect to relay
4. Execute: `ai-happy-design batch '[{"command":"slide","params":{...}}]'`

## Key Commands

- `batch` — Execute design operations from JSON (main LLM path)
- `command` — Execute single operation
- `schema <command>` — Print exact JSON shape for a command
- `validate` — Dry-run: check JSON without executing
- `guide` — Design intelligence (visual hierarchy, composition, effects)
- `tools --json` — List all commands with schemas
- `examples [category]` — Pre-built example payloads

## Design Workflow

1. `ai-happy-design command design.compute_tokens '{"width":1080,"height":1350}'`
2. `ai-happy-design command document.find_free_space '{"width":1080,"height":1350}'`
3. Generate batch JSON using `ai-happy-design schema <command> --json`
4. `ai-happy-design batch ops.json`
5. `ai-happy-design command export.image '{"nodeId":"...","scale":2}'`

## Links

- [Full command reference](https://aihappydesign.com/llms-full.txt)
```

**Step 2: Generate llms-full.txt**

Run: `ai-happy-design schema --all > /Users/nerveband/Documents/GitHub/ai-happy-design-site/llms-full.txt`

**Step 3: Deploy**

Run: `cd /Users/nerveband/Documents/GitHub/ai-happy-design-site && netlify deploy --prod --dir=.`

**Step 4: Commit both repos**

```bash
# Main repo
git add -A
git commit -m "feat: complete AI-native CLI — schemas, validation, lint, plugin features, no MCP"

# Site repo
cd /Users/nerveband/Documents/GitHub/ai-happy-design-site
git add llms.txt llms-full.txt
git commit -m "feat: add llms.txt and llms-full.txt for AI discoverability"
git push
```

---

### Task 15: Update Project Documentation

**Files:**
- Modify: `AGENTS.md` — remove MCP references, update architecture, add schema/validate/guide commands
- Modify: `CLAUDE.md` — remove MCP build steps, update deploy workflow
- Modify: `~/.claude/skills/ai-happy-design/SKILL.md` — update for CLI-only workflow
- Modify: memory `MEMORY.md` — update architecture notes

**Step 1: Update AGENTS.md**

Remove MCP sections, update system components to show CLI → relay → plugin (no MCP layer), update build commands, update CLI command list.

**Step 2: Update CLAUDE.md**

Remove MCP from build/deploy, update `make deploy` if it references MCP.

**Step 3: Update SKILL.md**

Remove MCP tool calling, update to CLI-only workflow with `schema`, `validate`, `guide` commands.

**Step 4: Commit**

```bash
git add AGENTS.md CLAUDE.md
git commit -m "docs: update project docs for CLI-only architecture, remove MCP references"
```

---

## Final Verification

After all tasks complete:

```bash
# Full build
go build ./cmd/ai-happy-design/

# All tests
go test ./...
cd plugin && npm run check && cd ..

# Verify no MCP references
grep -r "mcp" internal/ --include="*.go" -l  # should be empty
grep -r "mcp-go" go.mod  # should be empty

# Verify schema system
ai-happy-design schema text.create --json
ai-happy-design schema --all | wc -l

# Verify validation
echo '[{"command":"text.creat","params":{"text":"Hi","fontSize":-5,"color":"red"}}]' | ai-happy-design validate

# Verify design lint
echo '[{"command":"node.create_frame","params":{"width":1080,"height":1350}},{"command":"text.create","params":{"text":"Hi","parentId":"$last","fontSize":10,"color":"#666666"}}]' | ai-happy-design batch

# Verify guide
ai-happy-design guide | wc -c  # ~25-30KB

# Deploy
make deploy
```
