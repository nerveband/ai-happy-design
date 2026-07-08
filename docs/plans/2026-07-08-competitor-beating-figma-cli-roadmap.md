# Competitor-Beating Figma + CLI Roadmap

Date: 2026-07-08

## Goal

Make AI Happy Design the strongest Figma automation tool for AI agents by winning on two axes at once:

1. Deeper, current Figma control than competing MCP/plugin bridges.
2. A more agent-native CLI than competing Figma tools, including schema discipline, safety rails, proof gates, artifact handling, and compact introspection.

The durable moat remains: single Go binary, CLI-first batch workflows, schema validation, auto-correction, design intelligence, local/offline validation, MCP compatibility, and the Figma plugin bridge for full canvas write access.

## Competitive Baseline

Direct competitors to track:

- `southleft/figma-console-mcp`: broadest Figma MCP surface, cloud pairing, variables, comments, screenshots, FigJam, Slides, annotations, multi-file support, and agent skills.
- `vkhanhqui/figma-mcp-go`: Go MCP/plugin bridge, many read/write tools, prompts, variables, styles, screenshots, npm wrapper, and strong plugin write coverage.
- `silships/figma-cli`: CLI-oriented Figma automation with token presets, JSX-like render flows, accessibility checks, Iconify integration, FigJam, and screenshot verification.
- `gethopp/figma-mcp-bridge`: newer plugin bridge focused on streaming live Figma data with multiple-file support.
- `TranHoaiHung/figma-ui-mcp`: plugin + MCP bridge with read/write/screenshot loops.
- Southleft Figma skills repos: a newer category where design-system workflows are packaged as portable agent skills on top of Figma MCP.

Strategic implications:

- MCP is table stakes, but CLI remains the faster, cheaper, more reliable agent inner loop.
- Plugin bridge remains the winning architecture for full write access.
- Cloud pairing is a medium-term gap.
- Screenshot-based verification is becoming expected.
- Skills/workflow packaging is becoming a product surface, not just documentation.
- Thin endpoint wrappers are not enough; AHD should keep design intelligence and validation as differentiators.

## Phase 0: Fresh CLI Audit And Drift Report

Before feature work, run the full 85-point audit from the local `cli-best-practices` repo against the current `ahd-figma` binary.

Evidence commands:

```bash
ahd-figma --help
ahd-figma --version
ahd-figma schema --json
ahd-figma tools --llm --json
ahd-figma validate
ahd-figma status --json
ahd-figma batch --dry-run docs/examples/instagram-post-v3.json
```

Deliverables:

- Add `docs/research/cli-audit-2026-07.md`.
- Compare current behavior with the older external AHD scorecard.
- Record any drift, especially around MCP prompts, `commitUndo()`, structured errors, dry-run coverage, and JSON output consistency.
- Convert any failed audit item into a concrete implementation task.

## Phase 1: Figma API Catch-Up

Implement latest Figma Plugin API deltas behind runtime guards. Every beta/new runtime call should produce an explicit unsupported-feature error when unavailable.

### Motion

Add:

- `internal/schema/motion_schemas.go`
- `plugin/src/handlers/motion.ts`
- router wiring
- tests and syntax verification

Commands:

- `motion.get_styles`
- `motion.apply_style`
- `motion.remove_style`
- `motion.get_animations`
- `motion.apply_keyframes`
- `motion.remove_keyframes`
- `motion.set_timeline_duration`

Guard for runtime availability before touching Motion APIs.

### Shaders

Add guarded shader discovery/import/application:

- `effect.list_shaders`
- `effect.import_shader`
- `effect.apply_shader_effect`
- `fill.apply_shader`

Support paint and effect paths only where the runtime accepts them.

### Slots

Add component slot support:

- `component.create_slot`
- `component.reset_slot`
- `component.get_slots`

Update component/node serialization to include slot nodes and slot-relevant metadata where available.

### Grid Updates

Extend grid support:

- `layout.reorder_grid_rows`
- `layout.reorder_grid_columns`
- expose any current auto-flow/grid fields from `layout.get_grid_layout`

Keep existing grid commands backward-compatible.

### Noise And Texture Vectors

Extend existing effect commands:

- accept `noiseSizeVector`
- accept `noiseSizeX`
- accept `noiseSizeY`
- keep scalar `noiseSize` behavior

## Phase 2: Agent Reliability And Proof Loops

### Visual Verification

Add first-class screenshot/verification commands:

- `document.screenshot`
- `document.screenshot_selection`
- `verify.visual`

Return structured metadata and local artifact paths where possible. Support scale, format, node/selection/page targets, and output directory.

Document the agent loop:

1. Create or modify design.
2. Screenshot target.
3. Inspect artifact.
4. Apply corrections.
5. Screenshot again.

### Undo Safety

Call `figma.commitUndo()` centrally after write operations where available.

Requirements:

- Do not call for read-only commands.
- Keep response envelope unchanged.
- Add opt-out if needed: `commitUndo: false`.
- Add tests or live verification notes.

### Progress Events

For batch and long operations, emit optional progress messages:

- step index
- total steps
- command name
- elapsed time

Do not break existing response envelopes.

### Large Output Modes

Add or harden detail modes on high-volume reads:

- `detail: "minimal" | "compact" | "full"`
- `document.get_tree`
- `document.find_nodes`
- `node.get`
- `design_system.analyze`

## Phase 3: MCP Surface Upgrade

### MCP Prompts

Implement:

- `prompts/list`
- `prompts/get`

Back prompt content from `internal/tools/catalog_llm.go`.

Initial prompts:

- `design_strategy`
- `batch_strategy`
- `design_system_strategy`
- `visual_verification_strategy`
- `accessibility_strategy`
- `figma_api_guardrails`

### MCP Tool Metadata

Expose or encode:

- read/write/destructive safety
- idempotency
- requires Figma
- requires relay
- examples

Use existing schema metadata as the source.

### MCP Batch Discoverability

Make `batch` the obvious high-throughput path in MCP descriptions and `ahd_describe`. Agents should learn to use one batch for 20+ operations instead of many individual tool calls.

## Phase 4: CLI Hardening

This workstream is mandatory for beating competitors as a CLI, not just as a Figma bridge.

### Canonical Output Contract

Ensure every command obeys:

- data on stdout
- diagnostics/errors on stderr
- valid JSON for machine modes
- structured error shape:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "...",
  "hint": "...",
  "retryable": false
}
```

Add tests for stdout/stderr separation and JSON validity.

### `agent-context`

Add a compact machine-readable command:

```bash
ahd-figma agent-context --json
```

Return:

- version
- command groups
- recommended workflows
- output modes
- safety metadata
- exit codes
- config/env precedence
- dry-run and commit conventions
- artifact delivery support
- skill/docs pointers
- compact examples

This should be smaller and more workflow-oriented than `schema --all`.

### Schema Drift Gate

Add a contract verification target:

```bash
make verify-contracts
```

It should fail if these drift:

- CLI schema
- MCP tools
- docs command reference
- skill command docs
- website command docs

### Safety Metadata Enforcement

Every command schema should define:

- `safety`: read/write/destructive
- `idempotency`
- `supportsDryRun`
- `requiresFigma`
- `requiresRelay`
- `requiresAuth`

Add tests that fail on missing metadata.

### Dry-Run Everywhere

For every mutating command:

- support dry-run
- validate before side effects
- return what would happen
- show auto-fixes
- keep destructive work explicit

### Jobs Ledger

Add durable async job support:

- `jobs list`
- `jobs get <id>`
- `jobs resume <id>`
- `jobs cancel <id>`

Store job records under the AHD config directory as JSONL or SQLite.

Use jobs for:

- long exports
- large batches
- live validation runs
- future cloud relay operations

### Artifact Routing

Add consistent delivery sinks:

```bash
--deliver stdout
--deliver file:./out.json
--deliver dir:./artifacts
```

Use for screenshots, exports, audits, design-system reports, and parity reports.

### Profiles And Config Inspection

Add or harden:

- `profile list`
- `profile use`
- `profile inspect --redacted`
- `config sources`

Document precedence:

1. flags
2. environment
3. profile
4. config file
5. defaults

### Input Ergonomics

Support or verify:

- `@file.json` expansion
- `@data://base64,...` explicit encoding
- `--stdin`
- `--payload`
- `--payload-file`
- `--jq` or equivalent projection
- `--fields`
- `--limit`
- `--compact`
- `--detail`

### Proof Gates

Add first-class verification commands:

- `ahd-figma doctor`
- `ahd-figma doctor --json`
- `ahd-figma verify plugin`
- `ahd-figma verify syntax`
- `ahd-figma verify live`
- `ahd-figma verify release`

Agents should not need to remember scattered release/test commands.

### Feedback Loop

Add:

```bash
ahd-figma feedback "schema says X but command requires Y"
```

Store local JSONL feedback by default. Keep any upstream submission opt-in.

## Phase 5: Design-System Moat

### Token Presets

Add:

- `tokens.preset_tailwind`
- `tokens.preset_shadcn`
- `tokens.preset_material`
- `tokens.setup_system`

Support generating batch payloads or directly creating variables. Add DTCG import/export compatibility where practical.

### Design-System Health

Expand analysis into:

- `design_system.health`

Score:

- naming
- token usage
- component coverage
- accessibility
- spacing consistency
- typography consistency
- reusable styles

### Component-Set Analysis

Add:

- `component.analyze_set`
- `component.arrange_set`

Map variants to likely code states:

- hover
- focus
- disabled
- selected
- loading

Generate component documentation from Figma data.

### Design-Code Parity

Expand `parity.compare_code` into:

- `parity.audit_component`

Workflow:

1. inspect Figma node
2. inspect local component code
3. compare color, typography, spacing, radius, layout, and states
4. produce actionable diff

## Phase 6: FigJam, Slides, And Editor Coverage

Do not chase feature count blindly. Add editor coverage after core Design-mode and CLI hardening.

### Editor Context

Add:

- `document.get_editor_context`

Return:

- editor type
- supported domains
- unavailable APIs
- runtime feature flags

### Slides MVP

Add narrow support:

- `slides.get_current`
- `slides.create_slide`
- `slides.set_background`
- `slides.add_text`
- `slides.reorder`

### FigJam MVP

Add:

- `figjam.create_sticky`
- `figjam.create_shape`
- `figjam.create_connector`
- `figjam.get_board`

## Phase 7: Packaging And Distribution

### Skill Packaging

Generate portable skills from `internal/tools/catalog_llm.go`, preserving the catalog as the source of truth.

Targets:

- Claude skill
- Codex skill
- Cursor rules
- Gemini extension docs

### npm Wrapper Investigation

Evaluate an npm wrapper that downloads/runs the Go binary:

```bash
npx ahd-figma
```

This counters competitor install friction while preserving the single-binary model.

### Cloud Relay Design

Write a design proposal before implementation:

- six-character pairing code
- TTL
- single-use pairing
- WebSocket relay
- no file data persistence by default
- local-first default
- clear privacy model

## Phase 8: Test And Release Gates

Static gates:

```bash
go test ./...
go build ./...
cd plugin && npm run check
```

Plugin syntax gate:

```bash
grep -c '\?\.' plugin/dist/code.js
grep -c '\?\?' plugin/dist/code.js
grep -c '\.\.\.' plugin/dist/code.js
```

Contract gates:

```bash
make verify-contracts
ahd-figma schema --json
ahd-figma agent-context --json
```

Live Figma gate with an open doc:

1. Start relay and plugin.
2. Create a dedicated test page.
3. Run capability probes for Motion, Shaders, Slots, grid, and screenshots.
4. Run create -> screenshot -> modify -> screenshot loop.
5. Run design-system health.
6. Validate undo behavior.
7. Export artifacts for inspection.

## Proposed Build Order

1. Fresh 85-point CLI audit and drift report.
2. Fix scorecard regressions.
3. Add MCP prompts, central `commitUndo()`, and screenshot verification.
4. Add `agent-context`, schema safety metadata enforcement, and contract drift tests.
5. Implement Motion, Shaders, Slots, grid updates, and noise vectors.
6. Add jobs ledger, artifact delivery, profiles/config inspection, and proof gates.
7. Add design-system health and token presets.
8. Add Slides/FigJam MVPs.
9. Run live Figma validation against an open document.
10. Update public docs, skill files, and scorecard.

## Approval Checkpoint

Start with build order items 1-5 in one implementation branch. Validate against a live Figma document before expanding into design-system health, jobs, profiles, and Slides/FigJam.
