# AGENTS.md

Guidance for AI coding agents working on AI Happy Design v2.

## Project Goal

A single Go binary + Figma plugin that gives LLMs full Figma canvas access through:
- CLI for direct operations (schema-validated, design-linted)
- WebSocket relay to Figma plugin
- Built-in design intelligence (catalog + design guide)
- Schema system with auto-correction and fuzzy matching

## System Components

### 1) Go Binary
- Entry: `cmd/ai-happy-design/main.go`
- Modes: `ws` (relay only), `command`, `batch`, `tools`, `schema`, `validate`, `guide`

### 2) Relay Layer
- Server: `internal/ws/server.go`
- Client: `internal/ws/client.go`
- Legacy command routing: `internal/ws/command_routing.go`

### 3) Schema + Validation Layer
- Schema types: `internal/schema/types.go`
- Schema registry: `internal/schema/registry.go`
- Command schemas: `internal/schema/*_schemas.go`
- Validator: `internal/validate/validator.go` (fuzzy matching, named colors, auto-fix)
- Design lint: `internal/designlint/lint.go` (text sizing, contrast, spacing, scoring)
- **LLM catalog (SOURCE OF TRUTH)**: `internal/tools/catalog_llm.go`
- Describe tool: `internal/tools/describe.go`

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

| CLI Command | Returns |
|-------------|---------|
| `ai-happy-design schema` | List all commands with descriptions |
| `ai-happy-design schema <command> --json` | Exact JSON schema for a command |
| `ai-happy-design validate` | Dry-run validation (schema + design lint) |
| `ai-happy-design guide` | Design intelligence (visual hierarchy, composition, effects) |
| `ai-happy-design schema --all` | Full command reference (for llms-full.txt) |

### When updating design rules:

1. Edit ONLY `internal/tools/catalog_llm.go`
2. Run `go build ./...` to verify compilation
3. Rebuild binary: `make build && cp bin/ai-happy-design ~/bin/`
4. Restart relay if running
5. **Do NOT duplicate rules** into SKILL.md, AGENTS.md, or reference files — they all point to the CLI

### What references the catalog (but does NOT define rules):
- **Claude skill** (`~/.claude/skills/ai-happy-design/SKILL.md`) — workflow + "call design_guide for rules"
- **Skill reference files** (`references/design-patterns.md`) — quick offline fallback only
- **README.md** — user-facing overview, links to CLI commands
- **This file (AGENTS.md)** — architecture + development practices

## Build/Test Commands

```bash
make build                              # Go binary
go test ./...                           # Go tests (schema, validate, designlint, tools)
go build ./...                          # Verify compilation
cd plugin && npm run check && cd ..     # Plugin typecheck + build + syntax verification
ai-happy-design schema text.create      # Verify schema system
ai-happy-design validate                # Verify validation pipeline
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
2. Validate from plugin dir: `cd plugin && npm run check`
3. Verify no unsupported syntax in `dist/code.js`
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

## CLI-First Development

CLI is the only interface — fastest for bulk ops and what LLM agents use for heavy design work.

**Compact aliases** for batch/CLI: `frame`, `rect`, `text`, `fill`, `stroke`, `gradient`, `shadow`, `blur`, `glass`, `noise`, `texture`, `modify`, `mask`, `find`.

**Parameter shorthands**: `pid`=parentId, `w`=width, `h`=height, `sz`=fontSize, `ff`=fontFamily, `lh`=lineHeight (auto-PERCENT), `bg`=color, `r`=cornerRadius, `fillColor`=color.

**Step naming**: Use `snake_case` for batch step names. Names are auto-sanitized but clean names prevent interpolation issues.

**Semantic naming**: ALWAYS name every Figma layer by its content/role. Never leave defaults like 'Frame 47'.

## Supported Advanced Features

- **node.modify**: Unified "update any node" — pass nodeId + bag of properties (x, y, width, height, color, opacity, cornerRadius, visible, name, rotation, text, fontSize, isMask, etc.)
- **node.set_mask**: Create mask groups — pass mask shape nodeId + targetIds array
- **effect.apply_glass**: One-call glass morphism (fill + background blur + stroke) with intensity presets (light/medium/heavy)
- **effect.add_noise**: Noise overlay effects (monotone/duotone/multitone) — requires Figma Beta API
- **effect.add_texture**: Texture effects — requires Figma Beta API
- **document.find_nodes**: Unified search by name (query), type (nodeType), text content (textContent)
- **fillColor → color alias**: Automatically resolved everywhere (the #1 silent failure from cross-tool training data)

## Editing Rules for Agents

1. Preserve backward compatibility for legacy command names via `internal/ws/command_routing.go`
2. Keep error messages useful and pass-through to CLI
3. **Design rule changes go ONLY in `catalog_llm.go`** — never duplicate elsewhere
4. Add/adjust tests when routing or envelope logic changes
5. Do not silently change message schemas
6. Always target `es6` for plugin builds
7. Verify balance rules when generating design payloads

## Release Workflow (Full Pipeline)

When pushing code changes, follow this complete pipeline:

```bash
# 1. Build, sign, install, restart relay
make deploy

# 2. Commit and push to git
git add <changed-files>
git commit -m "Description of changes"
git push origin main

# 3. Tag and release via goreleaser
git tag v<next-version>
git push origin v<next-version>
GITHUB_TOKEN=$(gh auth token) goreleaser release --clean

# 4. Upgrade local binary to the released version
ai-happy-design upgrade

# 5. Update the ai-happy-design skill (if CLI features changed)
# - Edit ~/.claude/skills/ai-happy-design/SKILL.md with new features
# - Copy to repo: cp ~/.claude/skills/ai-happy-design/SKILL.md skills/ai-happy-design/SKILL.md
# - The skill is already committed in step 2 if you update it before committing

# 6. Sync skill to all AI CLI tools via skillshare
skillshare sync

# 7. Reopen Figma plugin to load new code.js
```

**When to run the full pipeline**: Any time code is pushed to git, all 7 steps should be executed. This ensures the released binary, local install, and skill docs are all in sync across every AI tool.

**Version numbering**: Check `git tag --sort=-v:refname | head -1` for the latest tag, then increment appropriately (patch for fixes, minor for features, major for breaking changes).

## External References
- Figma Plugin docs: https://developers.figma.com/docs/plugins/
- Figma Plugin API: https://developers.figma.com/docs/plugins/api/api-reference/
- Figma plugin manifest: https://developers.figma.com/docs/plugins/manifest
- Plugin typings: https://github.com/figma/plugin-typings
