# LLM Batch Optimization Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Reduce LLM batch payload verbosity and operation count while preserving backwards compatibility for CLI/MCP and plugin protocol envelopes.

**Architecture:** Add a normalization layer in batch execution paths (CLI + MCP bulk), expand interpolation syntax and command aliases in routing/interpolation layers, and move more styling concerns into create handlers so batch programs emit fewer operations. Keep transport protocol and legacy commands stable.

**Tech Stack:** Go (CLI/MCP/relay), TypeScript (Figma plugin), Cobra, gorilla/websocket, esbuild (ES6 target).

---

## Scope

1. Implement short interpolation syntax (`$name`, `$name.field`, `$last`) while preserving existing `${{...}}` syntax.
2. Add concise command aliases (`frame`, `rect`, `ellipse`, `text`, `fill`, `stroke`, etc.) through routing.
3. Implement command-aware shorthand parameter normalization for batch paths only (CLI `batch` + MCP `bulk.execute`).
4. Add combined create-operation capabilities in plugin handlers (`stroke`, `strokeWidth`, `noFill`, `opacity`, plus text create extras).
5. Wire Go MCP tool forwarding for new create params.
6. Expand tests for interpolation, routing aliases, and normalization behavior.
7. Verify with Go tests/build and plugin build + ES6 syntax checks.

## Non-Goals (this pass)

1. Parallel execution of batch operations.
2. YAML input support.
3. Transport redesign (streaming/chunking exports).

## Compatibility Rules

1. Keep existing command names and long interpolation syntax fully functional.
2. Do not change relay envelope format.
3. Additive params only for plugin handlers.
4. Preserve ES6-compatible plugin output.

## Phase 1: Go Batch Surface Reduction

### Task 1.1: Short interpolation syntax

**Files:**
- Modify: `internal/batchutil/interpolation.go`
- Test: `internal/batchutil/interpolation_test.go`

**Work:**
1. Add short syntax expansion before current placeholder matching.
2. Support:
   - `$name` -> `steps.name.result.id`
   - `$name.field` -> `steps.name.result.field`
   - `$last` -> `last.result.id`
   - `$last.field` -> `last.result.field`
3. Prevent accidental expansion of `${{...}}` and numeric/currency strings.

**Verification:**
- `go test ./internal/batchutil`

### Task 1.2: Command aliases

**Files:**
- Modify: `internal/ws/command_routing.go`
- Test: `internal/ws/command_routing_test.go`

**Work:**
1. Add aliases for compact command names:
   - `frame`, `rect`, `ellipse`, `line`, `image`, `text`
   - `fill`, `stroke`, `gradient`
   - `parent`, `autolayout`, `opacity`, `nofill`, `shadow`, `blur`
2. Keep dot-notation precedence unchanged.

**Verification:**
- `go test ./internal/ws`

### Task 1.3: Batch/bulk parameter normalization

**Files:**
- Add: `internal/batchutil/normalize.go`
- Add: `internal/batchutil/normalize_test.go`
- Modify: `cmd/ai-happy-design/main.go`
- Modify: `internal/tools/bulk.go`

**Work:**
1. Add command-aware normalization for shorthand params:
   - General: `w->width`, `h->height`, `pid->parentId`
   - Text: `sz->fontSize`, `ff->fontFamily`, `fs->fontStyle`, `lh->lineHeight`, `ls->letterSpacing`
   - Shape/node create: `sw->strokeWidth`, `bg->color`, `r->cornerRadius` (command-scoped only)
2. Add text default in normalized params:
   - If `lineHeight` is present without `lineHeightUnit` in text spacing/create contexts, default to `PERCENT`.
3. Apply normalization in CLI `batch` loop and MCP `bulk.execute` loop before interpolation/send.

**Verification:**
- `go test ./internal/batchutil ./internal/tools`

## Phase 2: Combined Create Ops in Plugin + Go Tool Forwarding

### Task 2.1: Shape create handlers

**Files:**
- Modify: `plugin/src/handlers/shape.ts`

**Work:**
1. Add helper(s) to apply optional stroke, opacity, and no-fill for geometry nodes.
2. Support in `create_rectangle`, `create_ellipse`, `create_polygon`, `create_star`:
   - `stroke`/`strokeColor`
   - `strokeWidth`
   - `noFill`
   - `opacity`
3. Keep existing behavior unchanged when new params absent.

### Task 2.2: Frame create handler

**Files:**
- Modify: `plugin/src/handlers/node.ts`

**Work:**
1. In `createFrame`, support:
   - `stroke`/`strokeColor`
   - `strokeWidth`
   - `noFill`
   - `opacity`
   - `cornerRadius`
2. Preserve existing auto-layout behavior.

### Task 2.3: Text create handler

**Files:**
- Modify: `plugin/src/handlers/text.ts`

**Work:**
1. In `createText`, support:
   - `opacity`
   - `letterSpacing` + `letterSpacingUnit`
   - `textCase`
2. Keep existing width/auto-resize and auto-layout child logic.

### Task 2.4: Go tool forwarding for new params

**Files:**
- Modify: `internal/tools/shape.go`
- Modify: `internal/tools/node.go`
- Modify: `internal/tools/text.go`

**Work:**
1. Add schema fields for new params where missing.
2. Forward optional params only when explicitly provided.

## Phase 3: Docs and Guidance Update

### Task 3.1: Catalog and CLI help updates

**Files:**
- Modify: `internal/tools/catalog_llm.go`
- Modify: `cmd/ai-happy-design/main.go`

**Work:**
1. Update examples to show compact syntax (`$name`, aliases, shorthand).
2. Document additive create params for operation collapse.

## Phase 4: Verification & Safety Gates

### Commands

1. `go test ./internal/batchutil ./internal/ws ./internal/tools`
2. `go build ./...`
3. `cd plugin && npm run build && cd ..`
4. `grep -c '\?\.' plugin/dist/code.js`
5. `grep -c '\?\?' plugin/dist/code.js`
6. `grep -c '\.\.\.' plugin/dist/code.js`

### Manual behavior checks

1. Run one batch using old syntax (`${{steps...}}`, full commands) and confirm unchanged behavior.
2. Run one batch using new syntax (`$name`, aliases, shorthand, combined create params) and confirm equivalent layout output.

