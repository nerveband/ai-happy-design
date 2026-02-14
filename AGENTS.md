# AGENTS.md

Guidance for AI coding agents working on AI Happy Design v2.

## Project Goal

A single Go binary + Figma plugin that gives LLMs full Figma canvas access through:
- MCP server for AI tool calls
- CLI for direct operations
- WebSocket relay to Figma plugin
- Built-in design intelligence (catalog + design guide)

## System Components

### 1) Go Binary
- Entry: `cmd/ai-happy-design/main.go`
- Modes: `mcp` (stdio MCP + embedded relay), `ws` (relay only), `command`, `batch`, `tools`

### 2) Relay Layer
- Server: `internal/ws/server.go`
- Client: `internal/ws/client.go`
- Legacy command routing: `internal/ws/command_routing.go`

### 3) MCP Tool Layer
- Registry: `internal/tools/registry.go`
- Domain tools: `internal/tools/*.go`
- **LLM catalog (SOURCE OF TRUTH)**: `internal/tools/catalog_llm.go`
- Describe tool: `internal/tools/describe.go`
- Bulk tool: `internal/tools/bulk.go`

### 4) Plugin Runtime
- Entry: `plugin/src/main.ts`
- Domain handlers: `plugin/src/handlers/*.ts`
- UI relay client: `plugin/src/ws/client.ts`
- UI: `plugin/src/ui/*`

## Figma Plugin Build Target (CRITICAL)

Figma's plugin sandbox uses QuickJS/WASM. Build target MUST be `es6`.

### Unsupported syntax (causes "Unexpected token" errors):
- `?.` optional chaining
- `??` nullish coalescing
- `{...obj}` object spread
- `for await...of` async iteration
- `?.()` optional call
- `??=`, `||=`, `&&=` logical assignment
- Class fields and private class features

### Required build config:

**esbuild.config.mjs**: `target: 'es6'`

**tsconfig.json**: `"target": "ES6"`, `"lib": ["ES2015", "ES2017"]`

### Post-build verification:
```bash
grep -c '\?\.' dist/code.js    # 0
grep -c '\?\?' dist/code.js    # 0
grep -c '\.\.\.' dist/code.js  # 0
```

### Plugin UI images:
Figma blocks `data:` URIs on `<img>` tags. Use `<div>` with CSS `background-image` instead.

## Protocol Contracts (Do Not Break)

### Response envelope
All responses wrapped: `{"type":"message","channel":"<ch>","message":{"id":"<id>","result":{...}}}`.
Errors: `{"type":"message","channel":"<ch>","message":{"id":"<id>","error":"..."}}`.
Never send bare `{id,error}`.

### Dynamic page access
Use `await figma.getNodeByIdAsync(...)`. Avoid deprecated sync getters.

### Image fill flow
- `set_image_fill_from_url`: try `figma.createImageAsync(url)` → fallback `fetch(url)` → bytes → `figma.createImage(bytes)` → set IMAGE fill
- `set_image_fill`: decode base64/data URL → `figma.createImage(bytes)` → set IMAGE fill

## Channel Resolution Order
1. Positional argument
2. `--channel` flag
3. `AHD_CHANNEL` env var
4. Relay preferred/active channel

## Design Intelligence — Central Source of Truth

### The Rule

**`internal/tools/catalog_llm.go` is the SINGLE source of truth for ALL design rules.** Nothing else should define design rules. Everything else references the catalog.

### What lives in catalog_llm.go:
- Design thinking (CSS-to-Figma, visual hierarchy, design decisions, layer organization)
- Design patterns (coordinates, grid, auto-layout, cards, typography, balance, scaling, aspect ratio, frame positioning)
- Playbook (12-step process)
- Workflow (batch vs single command)

### Discovery endpoints:

| Endpoint | Returns |
|----------|---------|
| `describe(action="catalog")` | Full catalog: tools + examples + design patterns + playbook |
| `describe(action="design_guide")` | Focused: design thinking + patterns + playbook |
| `describe(action="setup")` | Installation and connection instructions |

### When updating design rules:

1. Edit ONLY `internal/tools/catalog_llm.go`
2. Run `go build ./...` to verify compilation
3. Rebuild binary: `make build && cp bin/ai-happy-design ~/bin/`
4. Restart relay if running
5. **Do NOT duplicate rules** into SKILL.md, AGENTS.md, or reference files — they all point to the MCP

### What references the catalog (but does NOT define rules):
- **Claude skill** (`~/.claude/skills/ai-happy-design/SKILL.md`) — workflow + "call design_guide for rules"
- **Skill reference files** (`references/design-patterns.md`) — quick offline fallback only
- **README.md** — user-facing overview, links to MCP actions
- **This file (AGENTS.md)** — architecture + development practices

## Build/Test Commands

```bash
make build                              # Go binary
go test ./...                           # Go tests
go build ./...                          # Verify compilation
cd plugin && npm run build && cd ..     # Plugin build
```

## Development Practices (Learned)

### When modifying the catalog:
1. Edit `catalog_llm.go`
2. Run `go build ./...` to verify compilation
3. Rebuild binary: `make build`
4. Copy to bin: `cp bin/ai-happy-design ~/bin/`
5. Restart relay if running

### When modifying plugin handlers:
1. Edit `plugin/src/handlers/*.ts`
2. Build from plugin dir: `cd plugin && node esbuild.config.mjs`
3. Verify no ES2018+ syntax in `dist/code.js`
4. Reload plugin in Figma

### Batch testing:
- Create payload JSON in `docs/examples/`
- Run: `ai-happy-design batch -f docs/examples/payload.json`
- Check for step failures in output
- Export result: `ai-happy-design command export.image -p '{"nodeId":"...","scale":2}'`

### Common gotchas:
- esbuild MUST run from `plugin/` directory (relative paths)
- Batch JSON must be a plain array `[{...}]`, NOT wrapped in `{"operations":[...]}`
- `lineHeight` must be a plain number (e.g., `110` for 110%), NOT `{value, unit}` object
- Plugin auto-connects on startup but needs relay running first
- Default export scale is 2x (changed from 1x for quality)
- Large exports (e.g., 2160x3840 at 2x) may hang — use 1x for very large frames

## Editing Rules for Agents

1. Preserve backward compatibility for legacy command names via `internal/ws/command_routing.go`
2. Keep error messages useful and pass-through to CLI + MCP
3. **Design rule changes go ONLY in `catalog_llm.go`** — never duplicate elsewhere
4. Add/adjust tests when routing or envelope logic changes
5. Do not silently change message schemas
6. Always target `es6` for plugin builds
7. Verify balance rules when generating design payloads

## External References
- Figma Plugin docs: https://developers.figma.com/docs/plugins/
- Figma Plugin API: https://developers.figma.com/docs/plugins/api/api-reference/
- Figma plugin manifest: https://developers.figma.com/docs/plugins/manifest
- Plugin typings: https://github.com/figma/plugin-typings
- MCP specification: https://modelcontextprotocol.io/specification/latest/basic/lifecycle
