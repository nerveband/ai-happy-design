---
title: Architecture Overview
description: System components, data flow, CLI modes, plugin handler architecture, and schema validation pipeline.
---

AI Happy Design is a single Go binary that bridges LLM agents and the Figma canvas. It combines a CLI, a WebSocket relay, an MCP server, and a design intelligence engine into one tool.

## System Components

The system has four major components:

```
┌─────────────────────────────────────────────────┐
│                  AI Agent / User                 │
│         (Claude Code, Cursor, CLI, etc.)         │
└─────────────┬───────────────────┬───────────────┘
              |                   |
         MCP (stdio)         CLI (command/batch)
              |                   |
              v                   v
┌─────────────────────────────────────────────────┐
│               Go Binary                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │
│  │ MCP      │  │ CLI      │  │ Schema +     │  │
│  │ Server   │  │ Commands │  │ Validation   │  │
│  └────┬─────┘  └────┬─────┘  │ + Design     │  │
│       │              │        │ Intelligence │  │
│       v              v        └──────────────┘  │
│  ┌─────────────────────────┐                    │
│  │   WebSocket Relay       │                    │
│  │   (ws://127.0.0.1:3055) │                    │
│  └────────────┬────────────┘                    │
└───────────────┼─────────────────────────────────┘
                |
           WebSocket
                |
                v
┌─────────────────────────────────────────────────┐
│            Figma Desktop App                     │
│  ┌─────────────────────────────────────────┐    │
│  │         AI Happy Design Plugin          │    │
│  │  ┌──────────┐  ┌─────────────────────┐  │    │
│  │  │ WS       │  │ Domain Handlers     │  │    │
│  │  │ Client   │  │ (14 handler files)  │  │    │
│  │  └──────────┘  └─────────────────────┘  │    │
│  └─────────────────────────────────────────┘    │
│                                                  │
│            Figma Plugin API                      │
└─────────────────────────────────────────────────┘
```

### Go Binary

The binary (`cmd/ai-happy-design/main.go`) is the central hub. It contains:

- **CLI frontend** -- parses commands using Cobra, manages batch orchestration, handles input/output
- **WebSocket relay** -- `internal/ws/server.go` accepts connections from both CLI clients and the Figma plugin, routing messages between them via channels
- **WebSocket client** -- `internal/ws/client.go` connects to the relay for sending commands and receiving results
- **MCP server** -- `internal/mcp/server.go` implements the Model Context Protocol over stdio, exposing all commands as MCP tools
- **Schema registry** -- `internal/schema/registry.go` holds JSON schemas for every command
- **Validator** -- `internal/validate/validator.go` validates parameters against schemas with fuzzy matching, named color resolution, and auto-correction
- **Design lint** -- `internal/designlint/lint.go` checks for text sizing, overflow, overlap, and naming issues
- **Design intelligence** -- `internal/tools/catalog_llm.go` is the single source of truth for all design rules, patterns, and guidance
- **Embedded plugin** -- `internal/plugin/files/` contains the compiled plugin JS and HTML, synced via `make sync-plugin`

### WebSocket Relay

The relay (`internal/ws/server.go`) runs on port 3055 and acts as a message broker:

- **Clients** (CLI, MCP) send commands to a channel
- **Plugins** (Figma) join a channel and receive commands
- **Responses** flow back from plugin to client through the relay
- **Multiple channels** support multiple Figma files simultaneously

Key properties:

| Property | Value |
|----------|-------|
| Bind address | 127.0.0.1 (local-only) |
| Max message size | 64 MB |
| Read timeout | 300 seconds |
| Write deadline | 120 seconds |
| Port | 3055 (fixed) |

The relay starts automatically in MCP mode. In CLI mode, run `ai-happy-design ws` to start it separately. The `ws` command auto-kills any stale relay on the same port.

### Figma Plugin

The plugin (`plugin/src/`) runs inside Figma's sandboxed environment (QuickJS/WASM):

- **Entry point**: `plugin/src/main.ts` -- routes incoming commands to the correct handler
- **UI relay**: `plugin/src/ws/client.ts` -- WebSocket client in the plugin UI iframe
- **Handlers**: `plugin/src/handlers/*.ts` -- 14 domain-specific handler files
- **Utilities**: `plugin/src/utils/*.ts` -- shared helpers for node lookup, stable IDs, effect sanitization

The plugin is built with esbuild targeting ES6 (required by Figma's QuickJS sandbox). The compiled output is a single `code.js` and `index.html` that get embedded into the Go binary.

### MCP Server

The MCP server (`internal/mcp/server.go`) communicates over stdio using the Model Context Protocol. AI editors (Claude Code, Cursor, Windsurf, VS Code, Zed) launch the binary with `ai-happy-design mcp` and interact with it through structured tool calls.

The MCP server exposes the same commands available via CLI, but as individual tools with JSON Schema parameter definitions.

## Data Flow

### CLI Command Flow

```
User types:  ai-happy-design command text.create '{"text":"Hello"}'
      |
      v
  CLI parses command + params
      |
      v
  Schema validation (registry.go + validator.go)
      |
      v
  WS client connects to relay (port 3055)
      |
      v
  Relay routes message to plugin channel
      |
      v
  Figma plugin receives command
      |
      v
  Handler executes (text.ts -> figma.createText())
      |
      v
  Plugin sends result back through relay
      |
      v
  CLI prints JSON result to stdout
```

### Batch Flow

```
User types:  ai-happy-design batch ops.json
      |
      v
  CLI reads JSON array from file/stdin/arg
      |
      v
  Auto-placement: calls document.find_free_space
      |
      v
  For each operation:
    1. Expand aliases (frame -> node.create_frame)
    2. Expand shorthands (pid -> parentId, w -> width)
    3. Interpolate step references (${{steps.NAME.result.id}})
    4. Validate against schema
    5. Send to plugin via relay
    6. Receive result
    7. Store in step results map
      |
      v
  Optional: run design lint (--lint)
      |
      v
  Print structured output {ok, steps, summary, timing}
```

### MCP Flow

```
Editor launches:  ai-happy-design mcp
      |
      v
  MCP server starts on stdio
      |
      v
  Relay starts automatically (background)
      |
      v
  Editor sends tool call (JSON-RPC over stdio)
      |
      v
  MCP server routes to WS client -> relay -> plugin
      |
      v
  Result flows back: plugin -> relay -> WS client -> MCP -> editor
```

## CLI Modes

The binary supports several execution modes:

| Mode | Command | Description |
|------|---------|-------------|
| `ws` | `ai-happy-design ws` | Start the WebSocket relay server only |
| `command` | `ai-happy-design command <cmd> '<json>'` | Execute a single Figma command |
| `batch` | `ai-happy-design batch <file/json>` | Execute multiple operations in sequence |
| `mcp` | `ai-happy-design mcp` | Run as MCP server (stdio transport) |
| `tools` | `ai-happy-design tools` | List available tools and their descriptions |
| `schema` | `ai-happy-design schema [cmd]` | Inspect command schemas |
| `validate` | `ai-happy-design validate` | Dry-run validation without executing |
| `guide` | `ai-happy-design guide` | Print design intelligence reference |
| `actions` | `ai-happy-design actions <domain>` | List actions in a domain |
| `register` | `ai-happy-design register` | Auto-configure MCP in AI editors |
| `extract` | `ai-happy-design extract <html>` | Convert HTML/CSS to batch JSON |
| `upgrade` | `ai-happy-design upgrade` | Self-update to latest release |

## Plugin Handler Architecture (14 Domains)

The plugin organizes commands into 14 handler domains. Each domain is a separate TypeScript file in `plugin/src/handlers/`:

| Domain | File | Key Commands | Purpose |
|--------|------|--------------|---------|
| `node` | `node.ts` | create_frame, get_info, get_tree, move, resize, modify, set_mask | Frame and node management |
| `text` | `text.ts` | create, set_content, set_font, set_size, set_color, set_align | Text creation and styling |
| `shape` | `shape.ts` | create_rectangle, create_ellipse, create_polygon, create_star, create_image | Geometric shapes and images |
| `paint` | `paint.ts` | set_solid, set_gradient, set_image, set_image_url, set_stroke | Fills, gradients, strokes |
| `effect` | `effect.ts` | add_shadow, add_blur, apply_glass, add_noise, add_texture | Shadows, blurs, glass morphism |
| `layout` | `layout.ts` | set_auto_layout, set_padding, set_spacing, set_sizing, check_overlaps | Auto-layout and constraints |
| `layer` | `layer.ts` | set_order, bring_forward, send_backward, group, ungroup | Z-ordering and grouping |
| `boolean` | `boolean.ts` | union, subtract, intersect, exclude, flatten | Boolean operations |
| `component` | `component.ts` | create, create_instance, create_set, get_local | Components and instances |
| `style` | `style.ts` | create_paint, create_text, create_effect, apply | Reusable styles |
| `variable` | `variable.ts` | create, set_value, bind, unbind, create_collection | Variables and collections |
| `page` | `page.ts` | create, delete, rename, duplicate, set_current | Page management |
| `document` | `document.ts` | get_info, find_nodes, find_free_space, scan_text, lint | Document-level operations |
| `export` | `export.ts` | image, svg, pdf, batch | Export as PNG, JPG, SVG, PDF |

Each handler receives a `params` object, executes Figma API calls, and returns a result object. The main router in `plugin/src/main.ts` dispatches commands based on the `domain.action` pattern.

### Shared Utilities

Handlers share common utilities in `plugin/src/utils/`:

| Utility | Purpose |
|---------|---------|
| `getNode.ts` | Safe node lookup with `loadAsync()` pre-call to handle cross-page access |
| `stableId.ts` | Resolves transient session IDs to stable committed IDs from parent's `.children` list |
| `sanitizeEffects.ts` | Strips internal Figma properties (boundVariables, etc.) before write-back |
| `sanitizeFills.ts` | Same for fill arrays -- prevents read-only property errors |
| `layoutGuard.ts` | Prevents `layoutPositioning:ABSOLUTE` on children of non-auto-layout parents |

## Schema Validation Pipeline

Every command passes through a multi-stage validation pipeline before reaching the plugin:

```
Input params
    |
    v
1. Alias expansion (frame -> node.create_frame)
    |
    v
2. Shorthand expansion (pid -> parentId, w -> width, sz -> fontSize)
    |
    v
3. Color resolution (named colors -> hex, e.g., "red" -> "#FF0000")
    |
    v
4. fillColor -> color alias resolution
    |
    v
5. JSON Schema validation (type checking, required fields, enums)
    |
    v
6. Fuzzy matching (suggest corrections for misspelled field names)
    |
    v
7. Auto-correction (fix common mistakes automatically)
    |
    v
Validated params -> send to plugin
```

### Schema Registry

The registry (`internal/schema/registry.go`) holds a JSON Schema definition for every command. Schemas are defined in domain-specific files:

```
schema/node_schemas.go     -> node.create_frame, node.move, etc.
schema/text_schemas.go     -> text.create, text.set_font, etc.
schema/shape_schemas.go    -> shape.create_rectangle, etc.
schema/paint_schemas.go    -> paint.set_solid, etc.
schema/effect_schemas.go   -> effect.add_shadow, etc.
schema/layout_schemas.go   -> layout.set_auto_layout, etc.
```

All schemas are registered at startup via `init()` functions in each file.

### Fuzzy Matching

When a parameter name does not match any known field, the validator uses Levenshtein distance to suggest corrections:

```
Error: unknown field "fonSize" in text.create
Did you mean: "fontSize"?
```

### Auto-Correction

Common mistakes are fixed automatically without user intervention:

- `fillColor` is silently mapped to `color` (the #1 silent failure from cross-tool training data)
- Named colors (`"red"`, `"blue"`) are resolved to hex values
- `lineHeight` without `lineHeightUnit` defaults to `PERCENT`

### Design Lint

Post-creation validation checks for quality issues:

| Check | Description |
|-------|-------------|
| `overflow` | Children extending beyond parent bounds |
| `overlap` | Sibling elements overlapping |
| `text_too_large` | Font size exceeds 50% of parent height |
| `text_too_small` | Font size below 12px |
| `default_name` | Nodes with Figma default names ("Frame 47") |
| `oversized_child` | Child larger than parent by 10%+ |
| `absolute_child_non_autolayout` | `layoutPositioning:ABSOLUTE` on non-auto-layout child |

Lint runs automatically with `--lint` in batch mode, or on-demand via `document.lint`.

## Design Intelligence

The design intelligence engine lives in `internal/tools/catalog_llm.go`. This is the single source of truth for all design rules -- nothing else defines design rules, everything else references the catalog.

The catalog provides:

- **Design thinking** -- CSS-to-Figma mental model for LLM agents
- **Visual hierarchy** -- primary, secondary, tertiary, ambient layers
- **Design patterns** -- auto-layout, cards, grids, typography, shadows
- **Balance rules** -- sibling consistency, padding ratios, spacing
- **Type scale** -- modular scale with perfect fourth ratio (1.333), base = `width * 0.044`
- **Playbook** -- 10-step workflow for creating designs

The `design.compute_tokens` tool (`internal/tools/design_tokens.go`) generates concrete pixel values for any canvas size, with text on a 4px grid and spacing on an 8px grid.

Agents access design intelligence via `ai-happy-design guide` (CLI) or the design guide MCP tools.

## Key Source Files

| Path | Purpose |
|------|---------|
| `cmd/ai-happy-design/main.go` | Binary entry point, CLI setup |
| `internal/ws/server.go` | WebSocket relay server |
| `internal/ws/client.go` | WebSocket client |
| `internal/ws/command_routing.go` | Legacy command name routing |
| `internal/mcp/server.go` | MCP server (stdio) |
| `internal/schema/registry.go` | Schema registry |
| `internal/schema/types.go` | Schema type definitions |
| `internal/validate/validator.go` | Validation with fuzzy match + auto-fix |
| `internal/designlint/lint.go` | Design quality checks |
| `internal/tools/catalog_llm.go` | Design intelligence (source of truth) |
| `internal/tools/describe.go` | Tool descriptions for LLMs |
| `internal/tools/design_tokens.go` | Token computation engine |
| `plugin/src/main.ts` | Plugin entry point and router |
| `plugin/src/handlers/*.ts` | 14 domain handler files |
| `plugin/src/ws/client.ts` | Plugin WebSocket client |
| `plugin/esbuild.config.mjs` | Plugin build config (target: es6) |

## Build System

| Command | What It Does |
|---------|-------------|
| `make build` | Build the Go binary |
| `make build-plugin` | Build the Figma plugin |
| `make sync-plugin` | Copy plugin dist into Go embed |
| `make deploy` | Full build + sign + install + restart relay |
| `go test ./...` | Run all Go tests |
| `cd plugin && npm run check` | Plugin typecheck + build + syntax verify |

---

Made by [Ashraf Ali](https://ashrafali.net) | [GitHub](https://github.com/nerveband/ai-happy-design) | License: GPL-3.0
