# ahd-figma Full Parity and Priority Roadmap Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `ahd-figma` fully internally consistent, competitive with current Figma CLI/MCP tooling, and aligned with current Figma Plugin/REST/MCP APIs.

**Architecture:** Keep the existing advantage: a single Go binary with a schema-validated CLI, WebSocket relay, and Figma plugin write bridge. Add executable MCP tools by wrapping the existing schema/CLI/relay contract instead of creating a second command system. Treat `internal/schema` as the canonical command contract, `internal/tools/catalog_llm.go` as the canonical design-intelligence contract, and plugin handlers as the only Figma Plugin API execution layer.

**Tech Stack:** Go, Cobra, `mark3labs/mcp-go`, TypeScript Figma plugin, esbuild, WebSocket relay, Figma Plugin API, optional REST API client for metadata/dev-resource/webhook features.

## Implementation Status

Implemented on branch `codex/full-figma-parity-execution` on 2026-05-01.

- Schema surface expanded to 184 commands.
- MCP exposes 185 tools, including `ahd_describe`.
- Contract parity tests pass.
- Local REST commands cover oEmbed, file metadata, Dev Resources, and Webhooks V2.
- Plugin runtime includes focused node, CSS readback, guarded Figma Draw, grid layout, extended variables, and safer component reset.
- CLI output supports `--output-format json|jsonl|text` and simple `--jq`.
- Local analysis commands include `document.accessibility_audit`, `tokens.export`, and `parity.compare_code`.
- Docs, skills, `llms.txt`, `llms-full.txt`, site runbook, support matrix, command parity evidence, and live acceptance artifacts were updated.
- `make deploy` completes locally and restarts the relay.
- Live Figma acceptance passed after deploy using frame `57:456`.

Not performed by design: no push to `main`, no git tag, no GitHub release, and no release publishing.

---

## Evidence Base

This plan is based on:

- Current repo inspection on May 1, 2026.
- Official Figma Plugin API updates and REST changelog.
- Context7 docs for Figma Plugin API, Figma REST API, and Figma MCP guide.
- Direct local clone and source inspection of representative competitors under `/tmp/ahd-figma-research`.
- Local docs already in this repo about better CLIs and agent DX.
- Direct local clone and inspection of `nerveband/cli-best-practices` under `/tmp/ahd-cli-best-practices`, including the 50-check Agent CLI Audit, Agent DX Scale, contract-first principle, structured-output pattern, and `scorecards/ahd-figma.md`.

### Local Repo Findings

- `go run ./cmd/ahd-figma schema --json | jq length` returns `144`.
- `go run ./cmd/ahd-figma tools --json | jq '[.[] | length] | add'` returns `115`.
- `go run ./cmd/ahd-figma validate docs/examples/batch-interpolation.json` currently fails on interpolation-looking `nodeId` values and missing `parentId`.
- `go run ./cmd/ahd-figma validate docs/examples/instagram-story-ahd-promo.json` currently fails because the file is wrapped as `{ "operations": [...] }`, while validation expects a plain array.
- `internal/mcp/server.go` currently exposes prompts only, not executable tools.
- `plugin/src/handlers/component.ts` still uses deprecated `InstanceNode.resetOverrides()`.
- `plugin/src/handlers/variable.ts` supports core local variable CRUD, but not extended variable collections or overrides.
- `plugin/src/handlers/layout.ts` supports layout grids, but not the newer grid container/child APIs as first-class schema commands.
- `cmd/ahd-figma/main.go` and `cmd/ai-happy-design/main.go` are duplicated enough that feature parity can drift.
- `internal/tools/describe.go` is a separate hand-written catalog, so schema/catalog/plugin parity can drift.

### Local Better-CLI References

- `docs/llm-friendly-cli-mcp-notes-2026-02-13.md` says the target is discoverable schemas, deterministic machine output, stable errors, chaining, progress, pagination, `--json`, `--jq`, and MCP resources/tools/prompts.
- `docs/superpowers/specs/2026-04-02-ahd-figma-audit-and-improvement-plan.md` records the Agent DX and Agent CLI audit framing.
- `docs/superpowers/plans/2026-04-02-ahd-figma-best-in-class.md` already planned parts of CLI parity and should be treated as an input, not a source of truth where current code has moved on.
- `docs/research/competitive-landscape-2026-04.md` is useful but partially stale: this repo now has more variable/component/MCP scaffolding than the older notes implied.
- `/tmp/ahd-cli-best-practices/scorecards/agent-cli-audit.md` is the external runnable checklist to use as the acceptance gate: 50 checks across discoverability, structured output, input flexibility, safety rails, error handling, context discipline, predictability, agent knowledge, resilience, and distribution.
- `/tmp/ahd-cli-best-practices/principles/agent-dx-scale.md` is the 21-point subjective but evidence-backed scale that must remain 21/21 after the new work.
- `/tmp/ahd-cli-best-practices/principles/contract-first.md` backs this plan's central decision: schema is the canonical contract and docs/help/MCP must be generated from or checked against it.
- `/tmp/ahd-cli-best-practices/patterns/structured-output.md` backs the requirement for JSON/NDJSON, structured errors, meaningful exit codes, and field filtering.

### Competitor Code Studied

This was a direct implementation-level pass, not only README/metadata:

- `GLips/Figma-Context-MCP`
  - Inspected `src/services/get-figma-data.ts`, `src/mcp/tools/*`, `src/extractors/*`, `src/transformers/*`, tests.
  - Key implementation lesson: separate fetch, simplify, serialize, and metrics. Their hooks report progress and per-phase metrics without mixing telemetry into core logic.
- `vkhanhqui/figma-mcp-go`
  - Inspected `internal/tools.go`, `internal/tools_write_variables.go`, `internal/tools_write*.go`, `internal/leader.go`, `internal/bridge.go`, tests.
  - Key implementation lesson: executable MCP tools can be thin typed wrappers over a bridge. Useful notes include free-plan Figma variable-mode guidance and safe output-path handling for screenshots.
- `southleft/figma-console-mcp`
  - Inspected `src/index.ts`, `src/local.ts`, `src/core/*-tools.ts`, `src/apps/*`, and tests including accessibility, schema compatibility, write tools, plugin version sync, and design-code parity.
  - Key implementation lesson: tool modules grouped by workflow, strong tests for registered tool counts and schema compatibility, accessibility scans, design/code parity, token browser/dashboard MCP apps, plugin version sync.
- `grab/cursor-talk-to-figma-mcp`
  - Inspected `src/socket.ts`, generated server output, plugin manifest/plugin handlers.
  - Key implementation lesson: popular plugin-bridge ergonomics, but less deterministic/schema-first than AHD.
- `silships/figma-cli`
  - Inspected `src/index.js`, `src/figma-client.js`, `src/daemon.js`, `plugin/ui.html`, `plugin/code.js`, `src/shadcn.js`, README/CLAUDE.
  - Key implementation lesson: CLI-shaped design-system automation and JSX/shadcn ergonomics are interesting, but CDP/app patching is a path to avoid. Their safe plugin mode validates AHD's bridge approach.
- `RedMadRobot/figma-export`
  - Inspected README, package layout, examples, Swift export modules.
  - Key implementation lesson: production token/asset export needs config-driven target outputs and deterministic generated code, not just ad-hoc export commands.

### Official Figma API Updates to Incorporate

- Plugin API Update 124, March 26, 2026: Dev Mode focused node via `figma.currentPage.focusedNode`.
- Plugin API Update 123, January 26, 2026: Figma Draw support: `TEXT_PATH`, `TRANSFORM_GROUP`, `createTextPath`, `transformGroup`, `loadBrushesAsync`, variable-width strokes, dynamic/brush strokes, pattern fills/strokes.
- Plugin API Update 121, November 20, 2025: extended variable collections and overrides.
- Plugin API Update 120, November 6, 2025: Grid layout API updates and `HUG` grid tracks.
- REST API March 25, 2026: oEmbed API.
- REST API May 28, 2025: Webhooks V2 with `DEV_MODE_STATUS_UPDATE`.
- REST Dev Resources API: node-linked source/doc URLs for Dev Mode.
- Plugin API shared nodes expose `getCSSAsync()`.

---

## Current Strategic Position

`ahd-figma` should not become a clone of Figma Console MCP or TalkToFigma. Its strongest position is:

- CLI-first, batch-first, deterministic execution.
- Single binary distribution.
- Plugin bridge for full local read/write without REST rate limits.
- Schema validation, fuzzy correction, design linting, and design-intelligence catalog.
- Agent-friendly payload generation and validation before Figma is touched.

The missing piece is parity and surface coherence: the repo has many capabilities, but they are split across schema, tools catalog, plugin handlers, docs, and duplicate binaries.

---

## Definition of Full Parity

Full parity means four things:

1. **Schema parity:** every public command has one canonical `internal/schema.Schema` entry.
2. **Discovery parity:** every schema command appears in `ahd-figma tools --json`, `actions`, `schema`, `guide`, and MCP tools/resources when appropriate.
3. **Runtime parity:** every executable schema command routes through `internal/ws/command_routing.go` to a plugin handler or local command.
4. **Docs/examples parity:** every shipped example validates, every old binary name is clearly legacy, and docs use the current `ahd-figma` command unless explicitly documenting compatibility.

---

## File Map

### Core Contracts

- Modify: `internal/schema/types.go`
- Modify: `internal/schema/registry.go`
- Modify: `internal/schema/*_schemas.go`
- Modify: `internal/tools/describe.go`
- Modify: `internal/tools/catalog_llm.go`
- Modify: `internal/ws/command_routing.go`
- Create: `internal/contract/parity.go`
- Create: `internal/contract/parity_test.go`

### CLI and Shared Command Execution

- Modify: `cmd/ahd-figma/main.go`
- Modify: `cmd/ai-happy-design/main.go`
- Create: `internal/figmacli/command.go`
- Create: `internal/figmacli/batch.go`
- Create: `internal/figmacli/output.go`
- Create: `internal/figmacli/output_test.go`

### MCP

- Modify: `internal/mcp/server.go`
- Create: `internal/mcp/tools.go`
- Create: `internal/mcp/resources.go`
- Create: `internal/mcp/execute.go`
- Modify: `internal/mcp/server_test.go`

### Plugin Runtime

- Modify: `plugin/src/main.ts`
- Modify: `plugin/src/handlers/component.ts`
- Modify: `plugin/src/handlers/variable.ts`
- Modify: `plugin/src/handlers/layout.ts`
- Modify: `plugin/src/handlers/node.ts`
- Modify: `plugin/src/handlers/effect.ts`
- Create: `plugin/src/handlers/devmode.ts`
- Create: `plugin/src/handlers/draw.ts`
- Create: `plugin/src/handlers/rest-metadata.ts` only if REST commands are routed through plugin; otherwise skip.
- Modify/Create tests under `plugin/src/**/*.test.ts`.

### REST Features

- Create: `internal/figmaapi/client.go`
- Create: `internal/figmaapi/client_test.go`
- Create: `internal/schema/rest_schemas.go`
- Modify: `cmd/ahd-figma/config_cmd.go`
- Modify: `cmd/ai-happy-design/config_cmd.go`

### Docs and Examples

- Modify: `README.md`
- Modify: `docs/getting-started.md`
- Modify: `docs/llm-integration.md`
- Modify: `docs/cli-batch-and-payloads.md`
- Modify: `docs/examples/batch-interpolation.json`
- Modify: `docs/examples/instagram-story-ahd-promo.json`
- Modify: `skills/ai-happy-design/SKILL.md`
- Modify: `skills/ai-happy-design/ai-happy-design.skill.md`
- Create: `docs/reference/current-figma-api-support.md`
- Create: `docs/reference/command-surface-parity.md`
- Create: `docs/research/figma-tool-implementation-notes-2026-05.md`

---

## Chunk 1: Contract Parity and Drift Prevention

### Task 1: Add a Parity Model

**Files:**
- Create: `internal/contract/parity.go`
- Create: `internal/contract/parity_test.go`
- Modify: `internal/tools/describe.go`
- Modify: `internal/schema/registry.go`

- [ ] **Step 1: Add a contract package that can compare schema commands, tool catalog actions, and routing entries.**

Implement:

```go
package contract

type CommandSurface struct {
	Command     string
	Domain      string
	Action      string
	Aliases     []string
	ReadOnly    bool
	Destructive bool
	Source      string
}
```

Expose functions:

```go
func FromSchemas() []CommandSurface
func FromToolCatalog() []CommandSurface
func Diff(schema, catalog []CommandSurface) []Finding
```

- [ ] **Step 2: Make `internal/tools/describe.go` generated from schema where possible.**

Do not keep a separate hand-maintained list for commands that already have schemas. Keep custom descriptions only for design-intelligence sections.

- [ ] **Step 3: Add failing tests for current drift.**

Run:

```bash
go test ./internal/contract -run TestSchemaToolCatalogParity -v
```

Expected first result: FAIL showing `144` schema commands versus `115` tool actions.

- [ ] **Step 4: Update tool catalog generation until the test passes.**

The catalog should include every schema command by default, grouped by domain.

- [ ] **Step 5: Add route parity tests.**

For every non-local schema command, assert `internal/ws/command_routing.go` can resolve it.

Run:

```bash
go test ./internal/contract ./internal/ws -v
```

Expected: PASS.

- [ ] **Step 6: Commit.**

```bash
git add internal/contract internal/tools/describe.go internal/schema/registry.go
git commit -m "test: enforce Figma command surface parity"
```

### Task 2: Normalize the Two Figma Binaries

**Files:**
- Create: `internal/figmacli/command.go`
- Create: `internal/figmacli/batch.go`
- Create: `internal/figmacli/output.go`
- Modify: `cmd/ahd-figma/main.go`
- Modify: `cmd/ai-happy-design/main.go`

- [ ] **Step 1: Extract shared command execution out of both `main.go` files.**

Target shared entry points:

```go
type CommandOptions struct {
	BinaryName     string
	Version        string
	Command        string
	Params         map[string]interface{}
	Channel        string
	Live           bool
	DryRun         bool
	Fields         []string
	Limit          int
	TimeoutSeconds int
}

func ExecuteCommand(ctx context.Context, opts CommandOptions) (map[string]interface{}, error)
```

- [ ] **Step 2: Port all existing flags from `ahd-figma` and `ai-happy-design` into shared option mapping.**

No binary should have a command feature the other lacks unless explicitly documented as deprecated.

- [ ] **Step 3: Add tests comparing help/flag surfaces.**

Run:

```bash
go test ./cmd/ahd-figma ./cmd/ai-happy-design ./internal/figmacli -v
```

Expected: PASS.

- [ ] **Step 4: Commit.**

```bash
git add cmd/ahd-figma cmd/ai-happy-design internal/figmacli
git commit -m "refactor: share Figma CLI command execution"
```

---

## Chunk 2: Fix Existing User-Facing Drift

### Task 3: Fix Batch Interpolation Validation

**Files:**
- Modify: `internal/validate/validator.go`
- Modify: `internal/validate/*_test.go`
- Modify: `docs/examples/batch-interpolation.json`
- Modify: `docs/cli-batch-and-payloads.md`

- [ ] **Step 1: Add a failing test for interpolation placeholders.**

Inputs like these must pass pre-validation when interpolation is enabled:

```json
"nodeId": "${{steps.create_card.result.id}}"
```

Expected behavior:

- Skip regex pattern validation for strings matching supported interpolation expressions.
- Still validate required presence.
- Still fail malformed interpolation syntax.

- [ ] **Step 2: Decide and document the required parent rule.**

Current example uses `node.create_frame` without `parentId`, but schema blocks it. Pick one:

- Allow root-level frame creation without `parentId`, matching plugin behavior.
- Or update examples to pass an explicit current-page/root parent mechanism.

Recommended: allow root-level frame creation for `node.create_frame` because the plugin supports appending to `figma.currentPage`.

- [ ] **Step 3: Run validation.**

```bash
go run ./cmd/ahd-figma validate docs/examples/batch-interpolation.json
```

Expected: `ok: true`.

- [ ] **Step 4: Commit.**

```bash
git add internal/validate docs/examples/batch-interpolation.json docs/cli-batch-and-payloads.md
git commit -m "fix: allow documented batch interpolation validation"
```

### Task 4: Accept or Normalize `{operations: [...]}` Example Files

**Files:**
- Modify: `cmd/ahd-figma/validate.go`
- Modify: `cmd/ai-happy-design/validate.go`
- Modify: `internal/batchutil/*`
- Modify: `docs/examples/instagram-story-ahd-promo.json`
- Modify: tests for validate command.

- [ ] **Step 1: Add tests for both accepted batch file shapes.**

Supported inputs:

```json
[
  {"command": "node.create_frame", "params": {}}
]
```

and:

```json
{
  "operations": [
    {"command": "node.create_frame", "params": {}}
  ]
}
```

- [ ] **Step 2: Implement a single batch loader helper.**

Add:

```go
func DecodeBatchOperations(raw []byte) ([]map[string]interface{}, error)
```

Use it in `batch`, `validate`, tests, and docs examples.

- [ ] **Step 3: Run validation.**

```bash
go run ./cmd/ahd-figma validate docs/examples/instagram-story-ahd-promo.json
```

Expected: `ok: true`, or clear schema warnings unrelated to wrapper shape.

- [ ] **Step 4: Commit.**

```bash
git add cmd internal/batchutil docs/examples/instagram-story-ahd-promo.json
git commit -m "fix: normalize wrapped batch operation files"
```

### Task 5: Remove Deprecated Component Reset API Use

**Files:**
- Modify: `plugin/src/handlers/component.ts`
- Add/modify: plugin component tests.

- [ ] **Step 1: Add a unit test for `component.reset_instance` behavior shape.**

Expected plugin logic:

```ts
const instanceAny = node as any;
if (typeof instanceAny.removeOverrides === 'function') {
  instanceAny.removeOverrides();
} else {
  instanceAny.resetOverrides();
}
```

- [ ] **Step 2: Replace direct `resetOverrides()` call.**

This aligns with Figma's deprecation of `resetOverrides` in favor of `removeOverrides`.

- [ ] **Step 3: Run plugin checks.**

```bash
cd plugin && npm run check
```

Expected: PASS.

- [ ] **Step 4: Commit.**

```bash
git add plugin/src/handlers/component.ts plugin/src/**/*.test.ts
git commit -m "fix: use current Figma instance override reset API"
```

---

## Chunk 3: Real MCP Tool Surface

### Task 6: Expose Schema Commands as MCP Tools

**Files:**
- Modify: `internal/mcp/server.go`
- Create: `internal/mcp/tools.go`
- Create: `internal/mcp/execute.go`
- Modify: `internal/mcp/server_test.go`

- [ ] **Step 1: Add failing MCP tool registration tests.**

Expected:

- `tools/list` contains at least all safe/read-only schema commands.
- Mutating tools are exposed with descriptions that include safety notes.
- `bulk.execute` exists.
- `design.compute_tokens` exists and executes locally without Figma.

- [ ] **Step 2: Generate MCP tools from `internal/schema.All`.**

Use `mcp.NewTool` and convert `schema.Param` into JSON Schema input. Preserve:

- required fields
- enum
- min/max
- pattern
- default in description if direct default support is limited
- `ReadOnly`, `Destructive`, `Idempotent` in tool annotations/descriptions

- [ ] **Step 3: Execute tools through shared `internal/figmacli.ExecuteCommand`.**

Do not create a second execution path. MCP should call the same command/batch validation pipeline as CLI.

- [ ] **Step 4: Add MCP resource endpoints.**

Expose:

- `ahd://schema/all`
- `ahd://schema/{command}`
- `ahd://tools/catalog`
- `ahd://guide/design`
- `ahd://examples/batch`

- [ ] **Step 5: Keep existing prompts.**

The current prompts in `internal/mcp/server.go` remain valuable. MCP should expose prompts, resources, and tools.

- [ ] **Step 6: Run tests.**

```bash
go test ./internal/mcp ./internal/schema ./internal/contract -v
```

Expected: PASS.

- [ ] **Step 7: Commit.**

```bash
git add internal/mcp internal/figmacli internal/contract
git commit -m "feat: expose Figma schema commands as MCP tools"
```

### Task 7: Add MCP Tool Output Discipline

**Files:**
- Create/modify: `internal/figmacli/output.go`
- Modify: `internal/mcp/execute.go`
- Modify: `docs/llm-friendly-cli-mcp-notes-2026-02-13.md`

- [ ] **Step 1: Define one machine output envelope.**

Use:

```json
{
  "ok": true,
  "command": "node.get_info",
  "result": {},
  "meta": {
    "durationMs": 12,
    "channel": "active",
    "warnings": []
  }
}
```

Errors:

```json
{
  "ok": false,
  "error": {
    "code": "NODE_NOT_FOUND",
    "message": "...",
    "retryable": false,
    "hint": "..."
  }
}
```

- [ ] **Step 2: Add `--output json|jsonl|text` and `--jq` to CLI if missing.**

This completes the better-CLI notes in `docs/llm-friendly-cli-mcp-notes-2026-02-13.md`.

- [ ] **Step 3: Add golden output tests.**

Run:

```bash
go test ./internal/figmacli ./cmd/ahd-figma ./cmd/ai-happy-design -v
```

Expected: PASS.

- [ ] **Step 4: Commit.**

```bash
git add internal/figmacli cmd docs/llm-friendly-cli-mcp-notes-2026-02-13.md
git commit -m "feat: standardize CLI and MCP output envelopes"
```

---

## Chunk 4: Current Figma API Feature Coverage

### Task 8: Add Read-Only Dev Mode Helpers

**Files:**
- Create: `plugin/src/handlers/devmode.ts`
- Modify: `plugin/src/main.ts`
- Create: `internal/schema/devmode_schemas.go`
- Modify: `internal/ws/command_routing.go`
- Modify: `internal/tools/catalog_llm.go`
- Create: plugin tests.

- [ ] **Step 1: Add schemas.**

Commands:

- `devmode.get_focused_node`
- `devmode.get_context`
- `devmode.get_selection_context`

- [ ] **Step 2: Implement plugin handler.**

Use:

```ts
const focused = figma.currentPage.focusedNode;
```

Return null if unavailable or not in Dev Mode.

- [ ] **Step 3: Ensure mutating commands are not exposed as Dev Mode-compatible.**

Dev Mode plugin APIs are read-only except metadata/export. Keep this boundary explicit in docs.

- [ ] **Step 4: Run checks.**

```bash
go test ./internal/schema ./internal/ws ./internal/contract -v
cd plugin && npm run check
```

- [ ] **Step 5: Commit.**

```bash
git add internal/schema internal/ws internal/tools plugin/src
git commit -m "feat: add Dev Mode focused-node helpers"
```

### Task 9: Add `node.get_css` via `getCSSAsync()`

**Files:**
- Modify: `plugin/src/handlers/node.ts`
- Modify: `internal/schema/node_schemas.go`
- Modify: `internal/ws/command_routing.go`
- Modify: `internal/tools/catalog_llm.go`

- [ ] **Step 1: Add `node.get_css` schema.**

Params:

- `nodeId` required
- `includeChildren` optional boolean, default false

- [ ] **Step 2: Implement plugin call.**

Use:

```ts
const css = await (node as any).getCSSAsync();
```

Fallback to useful error if the method is unavailable.

- [ ] **Step 3: Add compact output mode.**

Return CSS as a map and optionally as text:

```json
{
  "id": "1:2",
  "name": "Button",
  "css": {
    "display": "flex"
  }
}
```

- [ ] **Step 4: Test.**

```bash
go test ./internal/schema ./internal/ws ./internal/contract -v
cd plugin && npm run check
```

- [ ] **Step 5: Commit.**

```bash
git add plugin/src/handlers/node.ts internal/schema internal/ws internal/tools
git commit -m "feat: expose Figma node CSS extraction"
```

### Task 10: Add Figma Draw Advanced Commands

**Files:**
- Create: `plugin/src/handlers/draw.ts`
- Modify: `plugin/src/main.ts`
- Create: `internal/schema/draw_schemas.go`
- Modify: `internal/ws/command_routing.go`
- Modify: `internal/tools/catalog_llm.go`

- [ ] **Step 1: Add guarded schemas.**

Commands:

- `draw.create_text_path`
- `draw.create_transform_group`
- `draw.load_brushes`
- `draw.set_variable_width_stroke`
- `draw.set_pattern_fill`
- `draw.set_pattern_stroke`

Each schema should be marked experimental or include `Beta Figma API` in the description.

- [ ] **Step 2: Implement runtime guards.**

Every handler checks whether the Figma API method exists before calling it:

```ts
if (typeof (figma as any).createTextPath !== 'function') {
  throw new Error('Figma API createTextPath is unavailable in this runtime');
}
```

- [ ] **Step 3: Keep ES6 build compatibility unless runtime testing proves modern syntax safe.**

Official docs now say modern JS is supported, but this repo has a hard-won QuickJS/WASM ES6 target rule. Do not relax it in this task.

- [ ] **Step 4: Test.**

```bash
cd plugin && npm run check
go test ./internal/schema ./internal/ws ./internal/contract -v
```

- [ ] **Step 5: Commit.**

```bash
git add plugin/src internal/schema internal/ws internal/tools
git commit -m "feat: add guarded Figma Draw commands"
```

### Task 11: Expand Grid Layout Commands

**Files:**
- Modify: `plugin/src/handlers/layout.ts`
- Modify: `internal/schema/layout_schemas.go`
- Modify: `internal/tools/guide_layout.go`

- [ ] **Step 1: Add schemas for grid containers and children.**

Commands:

- `layout.set_grid_container`
- `layout.set_grid_tracks`
- `layout.set_grid_child_position`
- `layout.get_grid_layout`

Params should include:

- `gridRowCount`
- `gridColumnCount`
- `gridRowGap`
- `gridColumnGap`
- `gridRowsSizing`
- `gridColumnsSizing`
- `gridRowSpan`
- `gridColumnSpan`
- `gridRowAnchorIndex`
- `gridColumnAnchorIndex`
- `gridChildHorizontalAlign`
- `gridChildVerticalAlign`

- [ ] **Step 2: Implement handler with property-existence guards.**

- [ ] **Step 3: Add guide entries explaining when to use Figma grid layout versus layout-grid overlays.**

- [ ] **Step 4: Test.**

```bash
go test ./internal/schema ./internal/tools ./internal/contract -v
cd plugin && npm run check
```

- [ ] **Step 5: Commit.**

```bash
git add plugin/src/handlers/layout.ts internal/schema/layout_schemas.go internal/tools/guide_layout.go
git commit -m "feat: support Figma grid layout APIs"
```

### Task 12: Extended Variable Collections and Overrides

**Files:**
- Modify: `plugin/src/handlers/variable.ts`
- Modify: `internal/schema/variable_schemas.go`
- Modify: `internal/tools/catalog_llm.go`
- Modify: `docs/reference/current-figma-api-support.md`

- [ ] **Step 1: Add schemas.**

Commands:

- `variable.extend_collection`
- `variable.extend_library_collection`
- `variable.get_values_for_collection`
- `variable.remove_mode_override`
- `variable.remove_collection_overrides`
- `variable.get_overrides`

- [ ] **Step 2: Implement plugin handlers with guards.**

Use:

- `figma.variables.extendLibraryCollectionByKeyAsync`
- `variableCollection.extend`
- `variable.valuesByModeForCollectionAsync`
- `variable.removeOverrideForMode`
- `extendedVariableCollection.removeOverridesForVariable`

- [ ] **Step 3: Improve free-plan/mode error handling.**

Borrow the useful guidance pattern from `figma-mcp-go`: if Figma returns a one-mode/free-plan limit, stop retrying and return a specific hint.

- [ ] **Step 4: Test.**

```bash
go test ./internal/schema ./internal/ws ./internal/contract -v
cd plugin && npm run check
```

- [ ] **Step 5: Commit.**

```bash
git add plugin/src/handlers/variable.ts internal/schema/variable_schemas.go internal/tools docs/reference/current-figma-api-support.md
git commit -m "feat: support extended variable collections"
```

---

## Chunk 5: REST Metadata, Dev Resources, and Webhooks

### Task 13: Add a Minimal REST Client

**Files:**
- Create: `internal/figmaapi/client.go`
- Create: `internal/figmaapi/client_test.go`
- Modify: `cmd/ahd-figma/config_cmd.go`
- Modify: `cmd/ai-happy-design/config_cmd.go`

- [ ] **Step 1: Add token resolution.**

Order:

1. `FIGMA_TOKEN`
2. `AHD_FIGMA_TOKEN`
3. config file key, if already supported

Never print tokens.

- [ ] **Step 2: Add scoped request helper.**

Only use `https://api.figma.com`. Do not support `http://api.figma.com`.

- [ ] **Step 3: Add typed error mapping.**

Map:

- 401/403 to `FIGMA_AUTH_ERROR`
- 404 to `FIGMA_NOT_FOUND`
- 429 to `FIGMA_RATE_LIMITED`
- 5xx to `FIGMA_API_ERROR`

- [ ] **Step 4: Test with mocked HTTP server.**

```bash
go test ./internal/figmaapi -v
```

- [ ] **Step 5: Commit.**

```bash
git add internal/figmaapi cmd
git commit -m "feat: add scoped Figma REST client"
```

### Task 14: Add REST Metadata and oEmbed Commands

**Files:**
- Create: `internal/schema/rest_schemas.go`
- Modify: `cmd/ahd-figma/main.go` or shared `internal/figmacli`
- Modify: `internal/tools/describe.go`
- Create: docs reference entries.

- [ ] **Step 1: Add local-only commands.**

Commands:

- `rest.file_metadata`
- `rest.oembed`

- [ ] **Step 2: Implement without plugin dependency.**

These commands should work even when no relay/plugin is running.

- [ ] **Step 3: Document required scopes.**

- `file_metadata:read` for metadata/oEmbed.

- [ ] **Step 4: Test.**

```bash
go test ./internal/figmaapi ./internal/schema ./cmd/ahd-figma -v
```

- [ ] **Step 5: Commit.**

```bash
git add internal/schema internal/figmacli cmd docs/reference
git commit -m "feat: add Figma REST metadata commands"
```

### Task 15: Add Dev Resources Commands

**Files:**
- Modify: `internal/schema/rest_schemas.go`
- Modify: `internal/figmaapi/client.go`
- Modify: shared CLI execution.

- [ ] **Step 1: Add schemas.**

Commands:

- `dev_resource.list`
- `dev_resource.create`
- `dev_resource.update`
- `dev_resource.delete`

- [ ] **Step 2: Implement REST calls.**

Use the REST Dev Resources endpoints. Require file key and node id for create/update.

- [ ] **Step 3: Add source-link helper.**

Support:

```bash
ahd-figma command dev_resource.create '{"fileKey":"...","nodeId":"1:2","name":"Button.tsx","url":"https://github.com/org/repo/blob/main/src/Button.tsx"}'
```

- [ ] **Step 4: Test.**

```bash
go test ./internal/figmaapi ./internal/schema ./internal/figmacli -v
```

- [ ] **Step 5: Commit.**

```bash
git add internal/figmaapi internal/schema internal/figmacli docs/reference
git commit -m "feat: manage Figma Dev Resources"
```

### Task 16: Add Webhooks V2 Management

**Files:**
- Modify: `internal/schema/rest_schemas.go`
- Modify: `internal/figmaapi/client.go`
- Modify: docs.

- [ ] **Step 1: Add schemas.**

Commands:

- `webhook.create`
- `webhook.list`
- `webhook.get`
- `webhook.update`
- `webhook.delete`
- `webhook.requests`

- [ ] **Step 2: Support `DEV_MODE_STATUS_UPDATE`.**

This is useful for Ready for Dev/Completed automation.

- [ ] **Step 3: Add safety warnings.**

Webhook commands touch external endpoints. Validate URLs and never echo passcodes.

- [ ] **Step 4: Test.**

```bash
go test ./internal/figmaapi ./internal/schema ./internal/figmacli -v
```

- [ ] **Step 5: Commit.**

```bash
git add internal/figmaapi internal/schema docs/reference
git commit -m "feat: manage Figma webhooks v2"
```

---

## Chunk 6: Accessibility, Design-Code Parity, and Token Export

### Task 17: Expand Design Lint to Accessibility Audit

**Files:**
- Modify: `internal/designlint/lint.go`
- Modify: `internal/designlint/*_test.go`
- Modify: `plugin/src/handlers/document.ts`
- Modify: `internal/schema/document_schemas.go`
- Modify: `internal/tools/catalog_llm.go`

- [ ] **Step 1: Add WCAG lint codes.**

Add:

- `WCAG_TEXT_CONTRAST`
- `WCAG_NON_TEXT_CONTRAST`
- `WCAG_TARGET_SIZE`
- `WCAG_LINE_HEIGHT`
- `WCAG_COLOR_ONLY`
- `WCAG_FOCUS_INDICATOR`
- `WCAG_IMAGE_ALT`

- [ ] **Step 2: Implement at least deterministic static checks.**

Start with checks that can be computed from node geometry, fills, text size, and names/descriptions.

- [ ] **Step 3: Add `document.accessibility_audit`.**

Return findings grouped by severity with WCAG reference and remediation hint.

- [ ] **Step 4: Test.**

```bash
go test ./internal/designlint ./internal/schema ./internal/tools -v
cd plugin && npm run check
```

- [ ] **Step 5: Commit.**

```bash
git add internal/designlint internal/schema internal/tools plugin/src/handlers/document.ts
git commit -m "feat: add accessibility audit lint rules"
```

### Task 18: Add Design-Code Parity Report

**Files:**
- Create: `internal/parity/design_code.go`
- Create: `internal/parity/design_code_test.go`
- Create: `internal/schema/parity_schemas.go`
- Modify: `cmd/ahd-figma/main.go` or shared local command handler.

- [ ] **Step 1: Define a code spec JSON shape.**

Fields:

- colors
- typography
- spacing
- radii
- shadows/effects
- opacity
- sizing
- accessibility

- [ ] **Step 2: Add command `parity.compare_code`.**

Inputs:

- `nodeId`
- `codeSpecPath` or inline `codeSpec`
- `threshold`

- [ ] **Step 3: Compare Figma node tree to code spec.**

Return scored differences. This does not need AST parsing in v1.

- [ ] **Step 4: Test with fixtures.**

```bash
go test ./internal/parity ./internal/schema ./internal/figmacli -v
```

- [ ] **Step 5: Commit.**

```bash
git add internal/parity internal/schema internal/figmacli
git commit -m "feat: add design-code parity reports"
```

### Task 19: Add Token Export Pipelines

**Files:**
- Create: `internal/exporttokens/export.go`
- Create: `internal/exporttokens/export_test.go`
- Create: `internal/schema/token_export_schemas.go`
- Modify: `plugin/src/handlers/variable.ts` if live variable data needs richer shape.
- Create: docs reference.

- [ ] **Step 1: Add command schemas.**

Commands:

- `tokens.export_json`
- `tokens.export_css`
- `tokens.export_tailwind`
- `tokens.export_swift`
- `tokens.export_android`

- [ ] **Step 2: Use Figma variables as source.**

Prefer plugin live variable data for local files. Optionally support REST variables later with explicit Enterprise-plan caveats.

- [ ] **Step 3: Add deterministic config file support.**

Example:

```yaml
outputs:
  css: ./tokens/figma.css
  tailwind: ./tokens/tailwind.tokens.json
  swift: ./Sources/DesignTokens/FigmaTokens.swift
```

- [ ] **Step 4: Test generated snapshots.**

```bash
go test ./internal/exporttokens ./internal/schema -v
```

- [ ] **Step 5: Commit.**

```bash
git add internal/exporttokens internal/schema docs/reference
git commit -m "feat: export Figma variables to platform tokens"
```

---

## Chunk 7: CLI Agent-DX Completion

### Task 20: Add Pagination, Limits, and Streaming to Large Results

**Files:**
- Modify: `internal/figmacli/output.go`
- Modify: `plugin/src/handlers/node.ts`
- Modify: `plugin/src/handlers/document.ts`
- Modify: `internal/schema/*_schemas.go`

- [ ] **Step 1: Identify large-output commands.**

At minimum:

- `node.get_tree`
- `document.find_nodes`
- `document.scan_text`
- `document.scan_by_type`
- `variable.get_all`
- `component.get_local`
- `style.get_all`

- [ ] **Step 2: Add common params.**

- `limit`
- `offset`
- `cursor`
- `depth`
- `fields`
- `compact`

- [ ] **Step 3: Implement output shaping.**

Do not make agents post-process huge JSON when the CLI can return the exact fields.

- [ ] **Step 4: Add `--output jsonl` for batch progress.**

Each line should be a standalone JSON object.

- [ ] **Step 5: Test.**

```bash
go test ./internal/figmacli ./internal/schema ./internal/contract -v
cd plugin && npm run check
```

- [ ] **Step 6: Commit.**

```bash
git add internal/figmacli internal/schema plugin/src
git commit -m "feat: add pagination and streaming output controls"
```

### Task 21: Add Plugin/CLI Version Mismatch Detection

**Files:**
- Modify: `plugin/src/ws/client.ts`
- Modify: `internal/ws/server.go`
- Modify: `internal/ws/server_test.go`
- Modify: `cmd/ahd-figma/status.go`

- [ ] **Step 1: Ensure plugin sends version on join.**

It should include:

```json
{
  "role": "plugin",
  "version": "0.12.0"
}
```

- [ ] **Step 2: Relay compares against binary version.**

Return warning in `/status` and CLI command metadata if mismatched.

- [ ] **Step 3: Add test for stale plugin warning.**

- [ ] **Step 4: Commit.**

```bash
git add plugin/src/ws/client.ts internal/ws cmd/ahd-figma
git commit -m "feat: detect Figma plugin version drift"
```

### Task 22: Make Port/Channel Multiplexing Safer

**Files:**
- Modify: `internal/ws/server.go`
- Modify: `internal/ws/client.go`
- Modify: `cmd/ahd-figma/main.go`
- Modify: `plugin/src/ws/client.ts`

- [ ] **Step 1: Replace auto-kill default with safer leader/follower behavior.**

Borrow the concept from `figma-mcp-go`: one leader owns the bridge, followers proxy or report clearly.

- [ ] **Step 2: Keep `--force-stop-relay` for explicit cleanup.**

Do not silently kill non-current user work.

- [ ] **Step 3: Test port conflict cases.**

```bash
go test ./internal/ws ./internal/relay ./cmd/ahd-figma -v
```

- [ ] **Step 4: Commit.**

```bash
git add internal/ws internal/relay cmd/ahd-figma plugin/src/ws
git commit -m "feat: harden relay leader and channel handling"
```

---

## Chunk 8: Docs, Skills, and Release

### Task 23: Write Current API Support Matrix

**Files:**
- Create: `docs/reference/current-figma-api-support.md`
- Create: `docs/reference/command-surface-parity.md`
- Create: `docs/research/figma-tool-implementation-notes-2026-05.md`

- [ ] **Step 1: Document supported Figma API areas.**

Sections:

- Plugin API basics
- Dynamic page access
- Images/video constraints
- Variables
- Components
- Grid layout
- Figma Draw
- Dev Mode
- REST metadata/dev resources/webhooks
- MCP

- [ ] **Step 2: Document competitor implementation findings.**

Include what was actually inspected and what ideas were adopted/rejected.

- [ ] **Step 3: Document command parity.**

Generate a table from `ahd-figma schema --json` and `tools --json` after parity fixes.

- [ ] **Step 4: Commit.**

```bash
git add docs/reference docs/research/figma-tool-implementation-notes-2026-05.md
git commit -m "docs: document Figma API and competitor implementation findings"
```

### Task 24: Update User-Facing Docs and Skills

**Files:**
- Modify: `README.md`
- Modify: `docs/getting-started.md`
- Modify: `docs/llm-integration.md`
- Modify: `skills/ai-happy-design/SKILL.md`
- Modify: `skills/ai-happy-design/ai-happy-design.skill.md`
- Modify: `site/` source files if present; otherwise regenerate or replace the docs-site source before building `site/dist`.

- [ ] **Step 1: Make `ahd-figma` the default command in docs.**

Mention `ai-happy-design` only as a legacy compatibility name.

- [ ] **Step 2: Add MCP tool workflow.**

Document:

```bash
ahd-figma register
ahd-figma mcp
ahd-figma tools --json
ahd-figma schema --all --json
```

- [ ] **Step 3: Update skill guidance.**

The skill should say:

- CLI is preferred for batch creation.
- MCP tools are valid for integrated AI editors.
- Always validate generated batch before executing.
- Use `node.get_css` and Dev Resources for design-to-code workflows.

- [ ] **Step 4: Commit.**

```bash
git add README.md docs skills/ai-happy-design
git commit -m "docs: refresh ahd-figma CLI and MCP guidance"
```

### Task 25: Update and Verify the Website

**Files:**
- Modify: `site/` source files if available.
- Modify: `site/dist/llms.txt`
- Modify: `site/dist/llms-full.txt`
- Modify: generated `site/dist/**` only as part of a documented site build/deploy flow.
- Create: `docs/reference/site-update-runbook.md` if the current site source is missing or external.

- [ ] **Step 1: Locate the real website source.**

Current working tree has `site/dist` and `site/node_modules`, but no obvious `site/src`, `site/content`, or `site/pages` source directory. Before implementation, determine whether the source lives outside this repo, was omitted, or needs to be restored.

- [ ] **Step 2: Document the website source of truth.**

If source is missing, create `docs/reference/site-update-runbook.md` explaining:

- where the source lives,
- how to build,
- how `llms.txt` and `llms-full.txt` are generated,
- how docs changes reach production.

- [ ] **Step 3: Update website content for new features.**

The website must cover:

- executable MCP tools, not just prompts,
- schema/catalog/plugin parity,
- current Figma API support matrix,
- Dev Mode helpers,
- `node.get_css`,
- Figma Draw commands,
- extended variables,
- REST metadata, Dev Resources, and webhooks,
- accessibility audit,
- token export,
- live-Figma verification workflow,
- current `ahd-figma` naming with legacy `ai-happy-design` notes.

- [ ] **Step 4: Rebuild site.**

Use the discovered site build command. If it is Astro, expected command is likely:

```bash
cd site && npm run build
```

Expected: `site/dist` updates cleanly and `llms.txt` / `llms-full.txt` include the new command surface.

- [ ] **Step 5: Verify website locally.**

Run a local preview server and inspect with Browser Use or Playwright:

```bash
cd site && npm run preview -- --host 127.0.0.1
```

Acceptance checks:

- home page loads,
- search works,
- new API support page exists,
- command reference includes new commands,
- `llms.txt` and `llms-full.txt` are reachable,
- no broken internal links in changed pages.

- [ ] **Step 6: Commit.**

```bash
git add site docs/reference/site-update-runbook.md
git commit -m "docs: update website for new ahd-figma features"
```

### Task 26: Run External CLI Best-Practices Audit

**Files:**
- Modify: `docs/reference/command-surface-parity.md`
- Modify: `docs/research/figma-tool-implementation-notes-2026-05.md`
- Modify: `README.md` if scores change.

- [ ] **Step 1: Clone or refresh the external audit repo.**

```bash
rm -rf /tmp/ahd-cli-best-practices
gh repo clone nerveband/cli-best-practices /tmp/ahd-cli-best-practices -- --depth=1
```

- [ ] **Step 2: Run the 50-check audit manually against the current binary.**

Use `/tmp/ahd-cli-best-practices/scorecards/agent-cli-audit.md` as the checklist. Save evidence for each category in `docs/reference/command-surface-parity.md`.

Minimum command evidence:

```bash
ahd-figma --help
ahd-figma --version
ahd-figma schema --json
ahd-figma tools --json
ahd-figma command design.compute_tokens '{"width":1080,"height":1350}' --dry-run
ahd-figma command design.compute_tokens '{"width":1080,"height":1350}' --fields scale,spacing
ahd-figma command not.real '{}' 2>/tmp/ahd-error.json; echo $?
cat docs/examples/batch-interpolation.json | ahd-figma validate -
```

- [ ] **Step 3: Re-score the Agent DX Scale.**

Use `/tmp/ahd-cli-best-practices/principles/agent-dx-scale.md`. Expected target remains `21/21`.

- [ ] **Step 4: Fail the release if any gate regresses.**

Required targets:

- Agent CLI Audit: `50/50`.
- Agent DX Scale: `21/21`.
- No stale scorecard claims without current command evidence.

- [ ] **Step 5: Commit the audit evidence.**

```bash
git add docs/reference/command-surface-parity.md docs/research/figma-tool-implementation-notes-2026-05.md README.md
git commit -m "docs: record current agent CLI audit evidence"
```

### Task 27: Live Figma App Acceptance Loop

**Files:**
- Create: `docs/reference/live-figma-acceptance.md`
- Create: `docs/examples/live-acceptance-full-parity.json`
- Modify code/docs only if this acceptance run finds defects.

- [ ] **Step 1: Start clean local build and relay.**

```bash
make build
ahd-figma setup --force
ahd-figma ws
```

Expected:

- plugin files extracted,
- relay starts on the configured port,
- no stale version warning once plugin is reopened.

- [ ] **Step 2: Open/reopen the Figma desktop plugin.**

Use Computer Use if needed to:

- open Figma Desktop,
- load the development plugin from the current manifest,
- click/connect the plugin,
- verify connected channel in the plugin UI.

- [ ] **Step 3: Verify relay and plugin status through CLI.**

```bash
ahd-figma status
ahd-figma command document.get_info '{}'
```

Expected:

- status shows active plugin channel,
- document info returns wrapped JSON,
- no CLI/plugin version mismatch.

- [ ] **Step 4: Run a full live parity batch.**

Create and run `docs/examples/live-acceptance-full-parity.json` covering:

- root frame creation,
- text creation,
- shape/image fill,
- auto-layout and grid layout,
- variable collection + variable bind,
- component create + instance + override,
- `node.get_css`,
- design lint/accessibility audit,
- export image,
- selection/focus,
- delete/cleanup or clearly named acceptance frame.

Run:

```bash
ahd-figma validate docs/examples/live-acceptance-full-parity.json
ahd-figma batch -f docs/examples/live-acceptance-full-parity.json --strict-quality --live
```

Expected:

- validation passes,
- batch succeeds,
- no failed steps,
- exported image exists.

- [ ] **Step 5: Visually verify in Figma.**

Use Computer Use or screenshots to confirm:

- nodes appear on canvas,
- names are semantic,
- no obvious overlaps,
- component instance exists,
- variable-bound style is visible,
- exported PNG matches expected content.

- [ ] **Step 6: Verify MCP execution against the same live plugin.**

Use an MCP client or stdio harness to confirm:

- `tools/list` exposes generated tools,
- `design.compute_tokens` works without Figma,
- a read command works against live Figma,
- a small write command works against live Figma,
- errors return structured tool results.

- [ ] **Step 7: Document defects and keep looping.**

If any step fails:

1. write the failing command and exact error into `docs/reference/live-figma-acceptance.md`,
2. fix the code,
3. rerun the relevant unit tests,
4. rerun this live step,
5. continue until all acceptance checks pass.

- [ ] **Step 8: Commit acceptance artifacts.**

```bash
git add docs/reference/live-figma-acceptance.md docs/examples/live-acceptance-full-parity.json
git commit -m "test: add live Figma acceptance proof"
```

### Task 28: Full Verification and Release Pipeline

**Files:**
- Modify only if verification finds issues.

- [ ] **Step 1: Run Go tests.**

```bash
go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run plugin checks.**

```bash
cd plugin && npm run check
```

Expected: PASS, including syntax verification.

- [ ] **Step 3: Build.**

```bash
make build
```

Expected: `bin/ahd-figma` and legacy binary build successfully.

- [ ] **Step 4: Validate examples.**

```bash
for f in docs/examples/*.json; do
  go run ./cmd/ahd-figma validate "$f" || exit 1
done
```

Expected: PASS for all examples intended to be executable batch payloads. Non-batch examples should be moved or documented.

- [ ] **Step 5: Verify command surface counts.**

```bash
go run ./cmd/ahd-figma schema --json | jq length
go run ./cmd/ahd-figma tools --json | jq '[.[] | length] | add'
```

Expected: counts match or a documented allowlist explains differences.

- [ ] **Step 6: Run deployment pipeline if pushing code.**

Follow AGENTS.md:

```bash
make deploy
git add <changed-files>
git commit -m "..."
git push origin main
git tag v<next-version>
git push origin v<next-version>
GITHUB_TOKEN=$(gh auth token) goreleaser release --clean
ahd-figma upgrade
skillshare sync
```

- [ ] **Step 7: Run external CLI audit gate.**

Run Task 26 completely. Expected: 50/50 and 21/21 with current evidence.

- [ ] **Step 8: Run live Figma acceptance gate.**

Run Task 27 completely. Expected: all CLI/plugin/MCP/live visual checks pass.

- [ ] **Step 9: Reopen Figma plugin.**

Required after plugin build changes.

---

## Recommended Implementation Order

Do this order, even if tempted to start with new features:

1. Chunk 1: Contract parity and drift prevention.
2. Chunk 2: Existing user-facing drift fixes.
3. Chunk 3: Real MCP tool surface.
4. Chunk 7 Task 21: Plugin/CLI version mismatch detection.
5. Chunk 4 Task 9: `node.get_css`.
6. Chunk 4 Task 8: Dev Mode helpers.
7. Chunk 4 Task 12: extended variables.
8. Chunk 4 Tasks 10-11: Figma Draw and grid expansion.
9. Chunk 5: REST metadata/dev resources/webhooks.
10. Chunk 6: accessibility, design-code parity, token export.
11. Chunk 7 remaining CLI-DX work.
12. Chunk 8 docs/release.

Reason: parity tests prevent future regressions, then MCP and new APIs can be added without increasing drift.

---

## Explicit Non-Goals

- Do not adopt CDP/app.asar patching from `silships/figma-cli`; keep plugin bridge as the safe path.
- Do not duplicate design rules outside `internal/tools/catalog_llm.go`.
- Do not relax ES6 build target just because current Figma docs mention modern JS support. Revisit only after direct Figma desktop runtime testing.
- Do not make REST the main write path. REST is for metadata, Dev Resources, webhooks, oEmbed, and Enterprise-gated variable workflows only.
- Do not create a second MCP-only command registry. MCP must derive from schema.

---

## Success Criteria

- `schema` and `tools` command counts match, or differences are explained by a checked allowlist.
- Every schema command routes locally, through REST, or through WebSocket/plugin.
- MCP clients can list and execute real AHD tools, not only prompts.
- All examples validate.
- Plugin build and syntax checks pass.
- Go tests pass.
- Docs consistently prefer `ahd-figma`.
- Current Figma API support matrix exists and names unsupported/beta APIs honestly.
- Competitive feature pressure is addressed without sacrificing AHD's deterministic CLI-first architecture.
- External CLI quality gates pass against `nerveband/cli-best-practices`: Agent CLI Audit remains 50/50 and Agent DX Scale remains 21/21 with current evidence, not stale scorecard claims.
- Live end-to-end Figma proof exists for CLI, MCP, relay, and plugin: create, inspect, modify, export, validate, and undo-safe operations all work in a real Figma desktop session.
- Website/docs are updated and rebuilt so public docs, `llms.txt`, and `llms-full.txt` describe the new features and current `ahd-figma` command surface.
