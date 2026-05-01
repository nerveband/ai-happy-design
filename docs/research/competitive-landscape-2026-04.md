# Figma CLI & AI Tool Competitive Landscape

**Date:** April 2, 2026
**Purpose:** Competitive analysis of Figma CLI/MCP tools and actionable recommendations for ai-happy-design

---

## Executive Summary

The Figma automation space has exploded around MCP (Model Context Protocol). There are now 20+ tools, but three direct competitors stand out against ai-happy-design:

| Tool | Stars | Language | Architecture | Read/Write |
|------|-------|----------|-------------|------------|
| **Figma Console MCP** (southleft) | 1,376 | TypeScript | MCP + Plugin WS + Cloud relay | Full |
| **silships/figma-cli** | 499 | Node.js | CLI + Plugin WS (+ CDP hack) | Full |
| **figma-mcp-go** (vkhanhqui) | 300 | Go | MCP + Plugin WS | Full |
| **ai-happy-design** (us) | — | Go | CLI + MCP + Plugin WS | Full |

All three use the same plugin-bridge architecture we do. The market has validated our approach. The question is: where are our gaps, and where is our moat?

**Our moat:** Single Go binary + design intelligence (catalog, lint, tokens, scoring) + schema validation with fuzzy matching/auto-correction + CLI batch with step interpolation. Nobody else combines execution with built-in design knowledge.

**Our gaps:** Variable/token CRUD, accessibility auditing, FigJam support, component library presets, multi-instance relay, undo support.

---

## Competitor Deep Dives

### 1. Figma Console MCP — The Feature Monster

**Repo:** https://github.com/southleft/figma-console-mcp
**Stars:** 1,376 | **Author:** Southleft (agency) | **License:** MIT

#### Architecture
- TypeScript, Node.js
- 4 connection modes: local stdio+WS, cloud (Cloudflare Workers + Durable Objects), remote SSE, remote SSE + cloud pair
- Plugin uses a bootloader architecture — thin shell loads fresh UI from MCP server on every open, never needs re-importing
- Multi-instance support via port scanning (9223-9232) with heartbeat files and orphan cleanup

#### 92+ Tools by Category

**Navigation & Status (3):** `figma_navigate`, `figma_get_status`, `figma_reconnect`

**Console Debugging (3):** `figma_get_console_logs` (filtered), `figma_watch_console` (real-time streaming, 300s), `figma_clear_console`

**Visual Debugging (2):** `figma_take_screenshot` (plugin/full-page/viewport, PNG/JPEG), `figma_reload_plugin`

**Design System Extraction (10):**
- `figma_get_variables` — tokens with enrichment (CSS/Tailwind/Sass exports, usage, dependencies)
- `figma_get_styles` — color, text, effect, grid styles with code exports
- `figma_get_component` — metadata or reconstruction spec format
- `figma_get_component_for_development` — component + rendered image at 2x
- `figma_get_component_image` — render as PNG/JPG/SVG/PDF
- `figma_get_file_data` — file structure with verbosity control (summary/standard/full)
- `figma_get_file_for_plugin` — optimized for plugin dev
- `figma_get_design_system_kit` — **full design system in one call** (tokens + components + styles + visual specs + images)
- `figma_get_design_system_summary` — high-level overview
- `figma_get_token_values` — variable values by collection and mode

**Design Creation (3):** `figma_execute` (arbitrary Plugin API JS), `figma_arrange_component_set`, `figma_set_description`

**Component Management (7):** search (local + published libraries), library discovery, instantiate with overrides, add/edit/delete component properties

**Variable Management (11):**
- Full CRUD: create, read, update, rename, delete for variables and collections
- Mode management: add and rename modes (Light, Dark, High Contrast, etc.)
- Batch: `figma_batch_create_variables` and `figma_batch_update_variables` (up to 100 per call, 10-50x faster)
- `figma_setup_design_tokens` — atomic one-shot: collection + modes + all variables in one call

**Node Manipulation (8):** resize, move, clone, delete, rename, set text, set fills, set strokes, create child

**Design-Code Parity (2):**
- `figma_check_design_parity` — compare Figma specs vs code across 8 dimensions with scored diff
- `figma_generate_component_doc` — markdown docs merging Figma data + code info

**Deep Component Analysis (2):**
- `figma_get_component_for_development_deep` — unlimited-depth tree, resolved token names, prototype reactions, annotations
- `figma_analyze_component_set` — variant state machine with CSS pseudo-class mappings (hover → `:hover`, focus → `:focus-visible`)

**Design Lint (1):** `figma_lint_design` — WCAG accessibility and design quality checks

**FigJam (9):** stickies (single + batch up to 200), connectors, shapes with text (diamond, ellipse, triangle, parallelogram, database, queue, file, folder), tables, code blocks, auto-arrange, board reading, connection graph

**Slides (17):** full CRUD, transitions (22 styles, 8 easing curves), grid layout, content manipulation, background, text styles, view modes

**Comments (3):** get, post (with node pinning), delete

**Annotations (3):** get (with child traversal), set (markdown, pinned properties, categories), list categories

**Cloud Relay (1):** `figma_pair_plugin` — 6-character pairing code for cloud connection

**MCP Apps (2, experimental):** Design System Dashboard (Lighthouse-style scoring), Token Browser (interactive)

#### Cloud Mode Architecture
```
AI Client → Cloud MCP Server (Cloudflare Worker) → PluginRelayDO (Durable Object) → WebSocket → Plugin → Figma
```
- 6-character pairing code, 5-minute TTL, single-use
- Enables Claude.ai, v0, Replit, Lovable to write to Figma without local binary
- 43 tools available in cloud mode (no real-time monitoring)

#### Notable Technical Details
- Adaptive response compression: auto-reduces payload at 100KB/200KB/500KB/1MB thresholds
- Enrichment pipeline: style resolution, relationship mapping, CSS/Tailwind/Sass export generation
- Design System Health Dashboard: Lighthouse-style scoring across 6 weighted categories (naming, token architecture, component metadata, accessibility, consistency, coverage)
- Codebase-aware component scanning: reads local `src/components/` to cross-reference against Figma components

#### Strengths vs Us
- Cloud mode (web AI clients can write to Figma)
- FigJam + Slides support (26 tools)
- Design-code parity checker (novel)
- Component library search across published libraries
- 92+ tools surface area
- Bootloader plugin auto-updates

#### Weaknesses vs Us
- No CLI mode (MCP only)
- No batch with step interpolation
- No design intelligence engine (no LLM catalog, no design tokens computation)
- No schema validation with fuzzy matching/auto-correction
- No single binary distribution
- `figma_execute` (raw JS) undermines structured safety

---

### 2. figma-mcp-go — The Go Rival

**Repo:** https://github.com/vkhanhqui/figma-mcp-go
**Stars:** 300 | **Author:** vkhanhqui | **License:** MIT | **Created:** March 2026

#### Architecture
- Single Go binary, stdio MCP only (no CLI)
- Flat `internal/` package structure (bridge.go, leader.go, follower.go, tools_read.go, tools_write.go, election.go, schema.go)
- Plugin: Svelte 5 + Vite, distributed as separate zip (not embedded)
- Distribution: npm wrapper (`npx @vkhanhqui/figma-mcp-go`) around cross-compiled Go binaries (6 platform combos)
- Dependencies: `coder/websocket`, `mark3labs/mcp-go`

#### 37 Tools (18 Read, 19 Write)

**Read (18):** `get_document`, `get_pages`, `get_metadata`, `get_selection`, `get_node`, `get_nodes_info` (batch), `get_design_context` (with detail levels: minimal/compact/full), `search_nodes`, `scan_text_nodes`, `scan_nodes_by_types`, `get_viewport`, `get_fonts`, `get_styles`, `get_variable_defs`, `get_local_components`, `get_annotations`, `get_reactions`, `get_screenshot`/`save_screenshots`

**Write (19):** `create_frame`, `create_rectangle`, `create_ellipse`, `create_text`, `import_image`, `set_text`, `set_fills`, `set_strokes`, `move_nodes`, `resize_nodes`, `rename_node`, `clone_node`, `set_auto_layout`, `delete_nodes`, `create_paint_style`, `create_text_style`, `create_effect_style`, `create_grid_style`, `update_paint_style`, `delete_style`, `create_variable_collection`, `add_variable_mode`, `create_variable`, `set_variable_value`, `delete_variable`

**6 MCP Prompts:** `read_design_strategy`, `design_strategy`, `text_replacement_strategy`, `annotation_conversion_strategy`, `swap_overrides_instances`, `reaction_to_connector_strategy`

#### Leader/Follower Election (Standout Feature)
```
Startup → try bind port 1994 → success? → Leader (owns WS to plugin)
                                → fail? → ping leader → healthy? → Follower (proxy via HTTP /rpc)
Background monitor (3-5s with jitter) → follower detects dead leader → takeover attempt
```
- Solves multi-AI-tool contention for single plugin connection
- Automatic failover with jitter to prevent thundering herd
- Smarter than our auto-kill approach

#### Other Notable Features
- `get_design_context` with detail levels: minimal (~5% tokens), compact (~30%), full (100%)
- `save_screenshots`: Go binary writes PNG/SVG to disk, returns local file path
- Progress reporting: plugin sends incremental updates for long ops, extends timeout
- `figma.commitUndo()` on every write — all mutations undoable via Cmd+Z
- Full variable and style CRUD
- Input validation: node ID format, enum values, required fields
- Node ID normalization: auto-converts hyphen format (LLM artifact) to colon format
- 100 MB WebSocket read limit

#### Strengths vs Us
- Leader/follower relay election (multi-tool support)
- `commitUndo()` on every write (safety net)
- Detail levels on read operations (token savings)
- npm distribution via `npx` (easy install)
- Progress reporting for long operations

#### Weaknesses vs Us
- No CLI mode at all (MCP only, no batch, no interactive)
- No design intelligence (no catalog, no lint, no tokens, no scoring, no auto-fix)
- No schema registry with JSON schema output
- No fuzzy matching or auto-correction
- No gradients, shadows, blur, glass, noise, texture, masks
- No `node.modify` (unified update)
- No `find_free_space`
- No embedded plugin (separate download)
- No `register` command
- No command aliases or parameter shorthands

---

### 3. silships/figma-cli — Batteries Included

**Repo:** https://github.com/silships/figma-cli
**Stars:** 499 | **Author:** Sil Bormueller | **License:** MIT | **Created:** Feb 2026

#### Architecture
- Pure Node.js, JavaScript (not TypeScript)
- Single monolithic `src/index.js` (~8,300 lines) + supporting modules
- 4 npm deps: commander, ws, chalk, ora
- Two connection modes:
  - **Yolo Mode:** Patches Figma's `app.asar` to re-enable CDP remote debugging (fragile)
  - **Safe Mode:** Plugin WebSocket bridge on ports 3456-3460 via daemon
- Distribution: git clone only (not on npm, no binaries)
- Delegates to `figma-use` (by dannote) for some Yolo Mode features

#### Full Command List

**Connection:** `connect`, `connect --safe`, `init`, `setup`, `status`, `daemon` (status/start/stop/restart/diagnose)

**Design Tokens:**
- `tokens preset shadcn` — 276 variables (primitives + semantic, Light/Dark modes)
- `tokens tailwind` — 242 variables (22 color families, 50-950 shades)
- `tokens ds` — IDS Base (71 variables, 5 collections)
- `tokens spacing`, `tokens radii`, `tokens components`
- `tokens import <file>` — from JSON

**Variables:** `var list`, `var create`, `var find`, `var visualize`, `var delete-all`, `var delete-batch`

**Collections:** `col list`, `col create`

**Create Elements:** `create frame/rect/ellipse/text/line/icon/image/group/component/autolayout`

**Batch:** `create-batch`, `delete-batch`, `bind-batch`, `set-batch`, `rename-batch`, `render-batch`

**Modify:** `set fill/stroke/radius/size/pos/opacity/autolayout/name`

**Layout:** `sizing hug/fill/fixed`, `padding`, `gap`, `align`

**Variable Binding:** `bind fill/stroke/radius/gap/padding`, `bind list`

**shadcn/ui (30 components):**
- `shadcn list`, `shadcn add [names...]`, `shadcn add --all`
- Components: Button, Card, Input, Select, Checkbox, Radio, Switch, Slider, Dialog, Sheet, Dropdown Menu, Popover, Tooltip, Avatar, Badge, Alert, Table, Tabs, Accordion, Separator, Label, Textarea, Progress, Skeleton, Toast, Command, Breadcrumb, Pagination, Calendar, Navigation Menu
- JSX templates manually authored to match official shadcn specs
- All use `var:` prefix for variable binding

**JSX Render Engine:**
- `render '<Frame bg="var:card" padding={16}><Text>Hello</Text></Frame>'`
- `render-batch` — array of JSX, arranged in row/column
- `var:` prefix syntax enables inline variable binding at creation
- Icons via `<Icon name="lucide:star" size={24} />`

**Export:** CSS variables, Tailwind config, JSX, Storybook stories, PNG, SVG, screenshot

**Analysis:** `analyze colors/typography/spacing/clusters`

**Lint (8+ rules):** `no-default-names`, `no-deeply-nested`, `no-empty-frames`, `prefer-auto-layout`, `no-hardcoded-colors`, `color-contrast`, `touch-target-size`, `min-text-size`. Presets: recommended, strict, accessibility, design-system. `--fix` auto-fix.

**Accessibility (full suite):**
- `a11y contrast` — WCAG AA/AAA, luminance calculation, large text detection
- `a11y vision` — color blindness simulation (Brettel matrices: protanopia, deuteranopia, tritanopia, achromatopsia), confusable color pair detection
- `a11y touch` — WCAG 2.5.8 (24x24 critical, 44x44 recommended)
- `a11y text` — font size min (12px), line height ratio (WCAG 1.4.12 ≥ 1.5x), ALL CAPS warning
- `a11y focus` — reading/focus order analysis
- `a11y audit` — full combined audit, letter grade (A+ through D)

**FigJam:** `fj list/info/nodes/sticky/shape/text/connect/delete/move/update/eval`

**Other:** `verify` (screenshot for AI validation), `blocks` (pre-built layouts), `combos` (variant generation), `sizes` (Small/Medium/Large scaling), `slots` (instance swap management), `raw query/select/export`, `eval`, `run`

#### Iconify Integration
- `<Icon name="lucide:star" size={24} color="var:foreground" />` in JSX
- Fetches SVG from Iconify API at runtime → `figma.createNodeFromSvg()`
- 150,000+ icons across all Iconify sets (lucide, mdi, heroicons, etc.)
- Color binding works: fills/strokes in SVG tree bound to variables

#### AI-Specific Features
- CLAUDE.md (500+ lines) as the primary AI instruction set
- `verify` command: screenshot current selection, return base64 for AI visual validation
- `.claude/MEMORY.md` for cross-session learned patterns
- No MCP server — all interaction via CLI commands Claude runs in terminal

#### Strengths vs Us
- shadcn/ui component library (30 components with variable bindings)
- Token presets (shadcn 276 vars, Tailwind 242 vars)
- 150K+ icons via Iconify
- Full accessibility suite with letter-grade audit
- JSX render engine for declarative component creation
- Color blindness simulation
- FigJam support
- `verify` for AI visual validation

#### Weaknesses vs Us
- No MCP server (CLI only, no structured tool definitions)
- No schema system, no JSON schema output, no validation
- No design intelligence engine (guidance is text in CLAUDE.md)
- No batch interpolation (no step references)
- No single binary distribution (git clone only)
- No `register` command
- No image compression
- CDP patching is fragile and a security concern
- Monolithic 8,300-line single file

---

## Broader Landscape (Tier 2+)

### MCP Servers

| Tool | Stars | Notes |
|------|-------|-------|
| **Framelink** (GLips/Figma-Context-MCP) | 14,111 | Read-only, REST API, context simplification for design-to-code. Most popular overall. |
| **Figma Official MCP** (figma/mcp-server-guide) | 926 | Hosted remote, read + write (beta). Rate limited: 6 calls/mo free, 200/day Pro. |
| **claude-talk-to-figma-mcp** (arinspunk) | 555 | TalkToFigma variant for Claude Desktop/Code. |
| **figma-mcp-go** (vkhanhqui) | 300 | Covered above. |
| **MatthewDailey/figma-mcp** | 209 | REST API-based, read + comments. Early entrant. |
| **TimHolden/figma-mcp-server** | 146 | Read-only, architecture for future variable management. |
| **Antonytm/figma-mcp-server** | 135 | Community write MCP, streaming HTTP. Hit Product Hunt #1. |
| **gethopp/figma-mcp-bridge** | 118 | Plugin streaming, "for the rest of us" (free users). |
| **uSpec** (redongreen) | 116 | Design system documentation generation. |
| **thirdstrandstudio/mcp-figma** | 67 | Full Figma REST API wrapper (28+ tools). |
| **magic-spells/figma-mcp-bridge** | 46 | 62 ops via plugin bridge, port 3055. |

### Non-MCP CLIs

| Tool | Stars | Notes |
|------|-------|-------|
| **Figma Code Connect** (figma/code-connect) | 1,449 | Official. Maps Figma components to codebase components. |
| **RedMadRobot/figma-export** | 811 | Swift CLI for iOS/Android asset export (colors, typography, icons). |
| **marcomontalbano/figma-export** | 338 | Web export (SVG, React components). |
| **alexchantastic/figma-export** | 129 | Bulk .fig/.jam file export via Playwright. |
| **opral/parrot-figcd** | 44 | CI/CD for Figma plugin publishing. |

### Design System Tools

| Tool | Stars | Notes |
|------|-------|-------|
| **Tokens Studio** (tokens-studio/figma-plugin) | 1,562 | Figma plugin for design token management. GitHub sync. |
| **Figma SDS** (figma/sds) | 758 | Official base design system (Variables, Styles, Components, Code Connect). |
| **generate-design-system** (natdexterra) | 19 | MCP skill for 5-phase design system generation. |

### Official Figma Developer Tools

| Tool | Stars | Notes |
|------|-------|-------|
| **figma/rest-api-spec** | 203 | OpenAPI spec + TypeScript types for REST API. |
| **figma/ai-plugin-template** | 143 | Template for plugins that talk to OpenAI. |
| **figma/figma-make-local-runner** | 105 | Run Figma Make code locally. |

---

## Two Architectural Camps

The market has split into two clear approaches:

### Camp 1: REST API (read-mostly)
- Uses Figma REST API with personal access tokens
- Subject to rate limits (6 calls/month free, 200/day Pro)
- Primarily read operations (design-to-code)
- Examples: Framelink, official Figma MCP, thirdstrandstudio

### Camp 2: Plugin Bridge (full read/write)
- WebSocket connection to a Figma plugin running in the desktop app
- No API rate limits, no token required
- Full read AND write access via Plugin API
- Examples: **ai-happy-design**, TalkToFigma, Figma Console MCP, figma-mcp-go, figma-cli

**We are firmly in Camp 2**, which is the more powerful approach. The official Figma MCP validates MCP as the standard, but its rate limits push power users toward plugin-bridge tools.

---

## Recommendations

### Priority 1: Table Stakes (Close Gaps)

These features exist in all three competitors. Not having them is a visible gap.

#### 1.1 Variable/Token CRUD Tools
**What:** Create collections, modes, variables. Read, update, rename, delete. Batch create/update.
**Why:** All 3 competitors have full variable management. Design system workflows require it. Currently, we can read variables (`design_system.analyze`) but can't create or modify them.
**Effort:** Medium — new handler in `plugin/src/handlers/variable.ts`, new schemas in `internal/schema/`, new MCP tools in `internal/tools/`.
**Reference:** figma-mcp-go's `tools_write.go` (clean implementation), Console MCP's batch operations (up to 100 per call).
**Scope:**
- `variable.create_collection` — name, modes array
- `variable.add_mode` — collectionId, name
- `variable.create` — collectionId, name, type (COLOR/FLOAT/STRING/BOOLEAN), values per mode
- `variable.update` — variableId, modeId, value
- `variable.rename` — variableId, newName
- `variable.delete` — variableId
- `variable.delete_collection` — collectionId
- `variable.batch_create` — array of variable definitions (up to 100)
- `variable.batch_update` — array of {variableId, modeId, value}
- `variable.setup_tokens` — atomic: collection + modes + all variables in one call

#### 1.2 `figma.commitUndo()` on All Write Operations
**What:** Call `figma.commitUndo()` after every mutation so users can Cmd+Z.
**Why:** Safety net. figma-mcp-go does this on every write. One-line addition per handler.
**Effort:** Low — add to each write handler's return path.
**Scope:** All handlers in `plugin/src/handlers/` that mutate the document.

#### 1.3 Response Detail Levels
**What:** `minimal/compact/full` verbosity on read operations.
**Why:** Token savings compound fast. figma-mcp-go's `get_design_context` saves 95% tokens at minimal level.
**Effort:** Low-Medium — extend `node.get_tree compact:true` pattern to all read tools.
**Scope:** `node.get_tree`, `document.scan_by_type`, `design_system.analyze`, `document.find_nodes`.

---

### Priority 2: Competitive Advantage (Strengthen Moat)

These extend our existing strengths into areas competitors haven't reached.

#### 2.1 Accessibility Audit Suite
**What:** Expand `designlint` into a proper accessibility audit: contrast (WCAG AA/AAA), touch targets (WCAG 2.5.8), text readability, focus order, color blindness detection.
**Why:** figma-cli has a full suite with letter grades. Our `designlint` already checks contrast and text sizing — extend it.
**Effort:** Medium — build on existing `internal/designlint/lint.go` infrastructure.
**Scope:**
- `a11y.contrast` — WCAG AA/AAA with large text detection
- `a11y.touch_targets` — 24x24 critical, 44x44 recommended
- `a11y.text_readability` — min font size, line height ratio, ALL CAPS warning
- `a11y.audit` — combined report with letter grade (A+ through D)
- Integrate into existing `validate` pipeline

#### 2.2 Design-Code Parity Checker
**What:** Compare Figma design specs against actual code implementation across dimensions (colors, spacing, typography, effects, corner radius, opacity, layout, sizing).
**Why:** Console MCP has this and nobody else does. Combined with our design intelligence engine, this could be our killer differentiator.
**Effort:** High — needs to read both Figma state and local code files.
**Scope:**
- `design.check_parity` — nodeId + code file path → scored diff report
- Dimensions: colors, spacing, typography, corner radius, effects, opacity, layout mode, sizing
- Output: per-dimension score, overall score, specific mismatches with fix suggestions

#### 2.3 Token Presets (One-Command Design System Bootstrap)
**What:** Ship pre-built token sets that create a full variable system in one command.
**Why:** figma-cli's `tokens preset shadcn` (276 vars) and `tokens tailwind` (242 vars) are hugely popular for onboarding. Combined with our `design.compute_tokens`, we can generate canvas-aware presets.
**Effort:** Medium — JSON definitions + `variable.setup_tokens` tool.
**Scope:**
- `tokens.preset shadcn` — 276 variables (primitives + semantic, Light/Dark)
- `tokens.preset tailwind` — 242 variables (22 color families)
- `tokens.preset minimal` — 50 essential tokens (colors, spacing, radii, typography)
- `tokens.preset from_canvas` — auto-generate tokens from `design.compute_tokens` for current canvas size

---

### Priority 3: Feature Expansion (New Use Cases)

#### 3.1 Leader/Follower Relay
**What:** When multiple AI tools each spawn a relay, elect one leader that owns the plugin WebSocket. Others proxy through it. Auto-failover.
**Why:** figma-mcp-go's implementation solves a real user pain point (Cursor + Claude Code + another tool competing for the plugin). Our auto-kill approach loses the other tool's connection.
**Effort:** Medium-High — refactor `internal/ws/server.go`.
**Scope:**
- Leader election on startup (try bind port → leader; fail → ping → follower)
- Follower HTTP proxy to leader's WebSocket
- Background health monitor with jitter (3-5s)
- Automatic failover on leader death
- Port advertisement files for discovery

#### 3.2 FigJam Support
**What:** Create and read FigJam boards — stickies, connectors, shapes, tables, code blocks.
**Why:** Both Console MCP (9 tools) and figma-cli have this. Expands use case beyond design into brainstorming, diagramming, planning.
**Effort:** Medium — new handler `plugin/src/handlers/figjam.ts`, new schemas.
**Scope:**
- `figjam.create_sticky` — text, color, position
- `figjam.create_stickies` — batch (up to 200)
- `figjam.create_connector` — source + target nodes, label
- `figjam.create_shape` — diamond, ellipse, triangle, etc. with text
- `figjam.create_table` — rows, columns, cell data
- `figjam.get_board` — read all board content
- `figjam.get_connections` — read connection graph

#### 3.3 Iconify Integration
**What:** Fetch any icon from 150K+ via Iconify API, render as vector nodes in Figma.
**Why:** figma-cli has this and it's extremely useful for design work. Icons are a constant need.
**Effort:** Medium — HTTP fetch in Go, SVG pass-through to plugin's `figma.createNodeFromSvg()`.
**Scope:**
- `icon.create` — name (e.g., "lucide:star"), size, color, position, parentId
- `icon.search` — query Iconify API for icon discovery
- `icon.batch` — create multiple icons at once
- Support `var:` binding for icon colors

#### 3.4 JSX-to-Figma Render Engine
**What:** Declarative markup that creates Figma components: `<Frame bg="var:card" padding={16}><Text>Hello</Text></Frame>`
**Why:** figma-cli's render engine is their most powerful feature. Dramatically simplifies complex component creation.
**Effort:** High — JSX parser + code generator.
**Scope:**
- Parse JSX-like syntax in Go
- Map to existing batch operations
- Support `var:` prefix for variable binding
- Support `<Icon>` elements via Iconify
- `render` command + MCP tool

---

### Priority 4: Strategic (Market Positioning)

#### 4.1 Cloud Relay Mode
**What:** Let web AI clients (Claude.ai, v0, Replit) write to Figma without a local binary.
**Why:** Console MCP's cloud mode via Cloudflare Workers + Durable Objects opens a massive new user base. Most AI tool users are on web clients.
**Effort:** High — requires Cloudflare Workers deployment, pairing flow, Durable Object relay.
**Scope:** Separate infrastructure project. Could be a hosted service.

#### 4.2 MCP Prompts
**What:** Expose design guidance as MCP prompts (not just `guide` command).
**Why:** figma-mcp-go ships 6 MCP prompts. It's a standard MCP feature we don't use. Our `catalog_llm.go` content is far richer — exposing it as prompts makes it discoverable via MCP protocol.
**Effort:** Low — `mcp-go` supports prompts natively, just wire up existing catalog content.
**Scope:**
- `design_strategy` — from catalog playbook
- `batch_workflow` — from catalog batch patterns
- `token_computation` — from design tokens guide
- `common_mistakes` — from catalog pitfalls

#### 4.3 Adaptive Response Compression
**What:** Auto-reduce response payload when exceeding context window limits.
**Why:** Console MCP's adaptive compression prevents context exhaustion. Important for large Figma files.
**Effort:** Low-Medium — implement in Go response pipeline.
**Scope:** Threshold-based compression (100KB → truncate details, 500KB → summary only, 1MB → emergency).

#### 4.4 Visual Verify for AI
**What:** `verify` command that screenshots current state and returns it for AI visual validation.
**Why:** figma-cli uses this to let Claude visually confirm what it created. We already have export — add a lightweight verify path.
**Effort:** Low — wrapper around existing export functionality.
**Scope:** `verify` command/tool that exports current selection at 1x, returns base64 or saves to temp file.

---

### Skip (Not Worth the Effort)

| Feature | Why Skip |
|---------|----------|
| `figma_execute` (raw JS execution) | Defeats the purpose of our schema validation. Our structured approach IS the moat. |
| CDP/app.asar patching | Fragile, breaks on Figma updates, security concern. Plugin bridge is correct. |
| Figma Slides support | Niche use case. Wait for user demand. |
| Bootloader plugin | Nice but our `go:embed` approach is cleaner for single-binary distribution. |
| Reconstruction spec format | Very niche (component version control). Not worth the complexity. |
| Codebase scanning | Tight coupling to user's project structure. Better left to the AI agent. |

---

## Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)
- [ ] Variable/token CRUD tools (1.1)
- [ ] `commitUndo()` on all writes (1.2)
- [ ] Response detail levels (1.3)
- [ ] MCP Prompts from catalog (4.2)

### Phase 2: Intelligence (Weeks 3-4)
- [ ] Accessibility audit suite (2.1)
- [ ] Token presets — shadcn, Tailwind, minimal (2.3)
- [ ] Visual verify command (4.4)
- [ ] Adaptive response compression (4.3)

### Phase 3: Expansion (Weeks 5-8)
- [ ] Leader/follower relay (3.1)
- [ ] FigJam support (3.2)
- [ ] Iconify integration (3.3)
- [ ] Design-code parity checker (2.2)

### Phase 4: Strategic (Weeks 9+)
- [ ] JSX-to-Figma render engine (3.4)
- [ ] Cloud relay mode (4.1)

---

## Key Takeaways

1. **MCP is the standard.** Every new Figma tool ships as an MCP server. Our dual CLI+MCP approach is actually an advantage — we serve both AI agents (MCP) and power users (CLI batch).

2. **Plugin bridge is the winning architecture.** Figma's official MCP validated MCP but pushed power users to plugin-bridge tools with its rate limits. We're on the right side.

3. **Design intelligence is our moat.** Nobody else combines execution with built-in design knowledge (catalog, lint, tokens, scoring, auto-correction). Every feature we add should leverage this.

4. **Variable management is table stakes.** All three competitors have it. This is our most visible gap and should be closed first.

5. **The Go binary advantage is real.** Single binary, fast startup, easy distribution, cross-platform. figma-mcp-go validates this approach but lacks our depth.

6. **Don't chase feature count.** Console MCP has 92+ tools but many are thin wrappers. Our 40+ tools with schema validation, design lint, and auto-correction provide more value per tool.
