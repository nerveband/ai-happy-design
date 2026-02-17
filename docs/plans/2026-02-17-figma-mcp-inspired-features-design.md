# Figma MCP-Inspired Features Design

**Date**: 2026-02-17
**Context**: After reviewing Figma's official MCP server documentation, three actionable features were identified that improve LLM usability of ai-happy-design.

---

## Feature 1: Export to Temp File (Inline Screenshot)

### Problem
LLMs like Claude Code can "see" images via their file Read tool, but our export currently requires the user to specify an output path or only works inline via MCP ImageContent. There's no convenient way for an LLM to export a node, get a path back, and visually inspect it.

### Design
- When `export.image` (or `export.svg`/`export.pdf`) is called without an explicit output path, save to `/tmp/ahd-export-<nodeId>-<timestamp>.<ext>`
- Return JSON: `{"path": "/tmp/ahd-export-...", "id": "1:23", "name": "My Frame", "format": "PNG", "size": 42816}`
- MCP: keep existing ImageContent behavior, add the temp file path to the metadata text block
- Plugin: no changes (already returns base64 data)

### Changes
| File | Change |
|------|--------|
| `internal/tools/export.go` | Add temp file write logic when no output path specified. Build path as `/tmp/ahd-export-<nodeId>-<timestamp>.<ext>`. Return path in result JSON. |
| `cmd/ai-happy-design/main.go` | CLI export subcommand: if no `-o` flag, use temp file path. Print JSON with path to stdout. |

### Why
Claude Code can `Read /tmp/ahd-export-1:23-1708200000.png` and visually verify the design. No manual file management needed.

---

## Feature 2: Compact Tree Output

### Problem
`node.get_tree` returns full serialized node data (fills, strokes, effects, text content, etc.). For structural discovery — finding node IDs, understanding hierarchy — this wastes tokens. Figma's `get_metadata` returns lightweight XML for exactly this reason.

### Design
- Add `compact` boolean param to `node.get_tree` (default false)
- When `compact: true`, plugin returns a flat JSON array instead of nested tree
- Each entry: `{id, type, name, x, y, w, h, childCount, parentId, depth}`
- `depth` param still controls scan depth (default 3)
- Both CLI and MCP return the same flat array

### Example Output
```json
[
  {"id": "1:2", "type": "FRAME", "name": "Card", "x": 0, "y": 0, "w": 400, "h": 300, "childCount": 3, "parentId": null, "depth": 0},
  {"id": "1:3", "type": "TEXT", "name": "Title", "x": 24, "y": 24, "w": 352, "h": 40, "childCount": 0, "parentId": "1:2", "depth": 1},
  {"id": "1:4", "type": "RECTANGLE", "name": "Avatar", "x": 24, "y": 80, "w": 64, "h": 64, "childCount": 0, "parentId": "1:2", "depth": 1}
]
```

### Changes
| File | Change |
|------|--------|
| `plugin/src/handlers/node.ts` | Add `serializeNodeCompact()` function. When `compact` param is truthy, call it instead of `serializeNode()`. Walk tree, collect flat array with `{id, type, name, x, y, w, h, childCount, parentId, depth}`. |
| `internal/tools/node.go` | Add `compact` boolean param to `get_tree` action definition. Pass through to plugin. |

### Why
Flat array is ~3-5x fewer tokens than full nested serialization. LLMs can scan it to find node IDs, then call targeted commands on specific nodes. Mirrors Figma MCP's metadata-first → drill-down pattern.

---

## Feature 3: Design System Rules Generator

### Problem
When an LLM creates designs in an existing Figma file, it doesn't know what styles, variables, and components already exist. It creates new colors instead of reusing existing paint styles, new text nodes instead of applying text styles, new shapes instead of instantiating components. This produces inconsistent designs.

### Design
- New standalone tool: `design_system` with initial action `analyze`
- CLI: `ai-happy-design command design_system.analyze` (no params required)
- Go layer orchestrates three existing plugin calls: `style.get_all`, `variable.get_all`, `component.get_local`
- Aggregates results into a structured rules document

### Output Format
```json
{
  "colors": {
    "styles": [{"name": "Primary/Blue", "id": "S:abc"}],
    "variables": [{"name": "color/primary", "value": "#7c3aed"}],
    "rule": "Use these existing colors. Do not introduce new colors without adding them as variables."
  },
  "typography": {
    "styles": [{"name": "Heading/H1", "id": "S:def"}],
    "rule": "Apply text styles by ID when available. Match font families and weights."
  },
  "spacing": {
    "variables": [{"name": "spacing/sm", "value": 8}, {"name": "spacing/md", "value": 16}],
    "rule": "Use spacing variables for padding and gaps. Snap to the existing scale."
  },
  "components": {
    "available": [{"name": "Button", "id": "1:23"}, {"name": "Card", "id": "1:45"}],
    "rule": "Instantiate existing components instead of rebuilding from scratch."
  },
  "summary": "This file has 12 paint styles, 3 text styles, 8 variables, and 6 components."
}
```

### Changes
| File | Change |
|------|--------|
| `internal/tools/design_system.go` | New file. Define `design_system` tool with `analyze` action. Orchestrate three plugin calls (`style.get_all`, `variable.get_all`, `component.get_local`). Categorize results: colors (paint styles + COLOR variables), typography (text styles), spacing (FLOAT variables with spacing/size-like names), components. Generate rules text per category. |
| `internal/tools/register.go` | Register the new `design_system` tool. |
| `internal/tools/catalog_llm.go` | Add `design_system` to LLM catalog with usage guidance. |

### Future Actions (not in this iteration)
- `design_system.lint` — check a node against the file's design rules
- `design_system.export_tokens` — export variables as CSS custom properties or JSON tokens
- `design_system.suggest` — given a design intent, suggest which existing styles/components to use

### Why
LLMs producing designs in existing files will maintain consistency. The rules document gives them actionable guidance: "use style S:abc for blue" instead of guessing hex values.

---

## Implementation Priority

1. **Feature 2 (Compact Tree)** — smallest change, immediate token savings, plugin + Go
2. **Feature 1 (Export Temp File)** — Go-only change, quick win
3. **Feature 3 (Design System Analyzer)** — new tool, Go-only (orchestrates existing plugin handlers), most impactful for design quality
