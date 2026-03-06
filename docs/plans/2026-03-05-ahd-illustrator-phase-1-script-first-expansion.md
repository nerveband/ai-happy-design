# AHD Illustrator Phase 1: Script-First Expansion

## Goal
- Expand `ahd-illustrator` aggressively using Illustrator's native scripting surface only.
- Keep the product usable with just:
  - Adobe Illustrator on macOS
  - the `ahd-illustrator` CLI
- Do not require the Adobe Illustrator SDK.
- Do not require the optional native plugin.
- This phase is implementation-only. It is not an exploration pass, prototype pass, or partial scaffold pass.
- Phase 1 is complete only when the script-backed Illustrator surface exposed by the Adobe 2025 scripting references has been fully integrated into the CLI in an agent-first form.

## Why Phase 1 Comes First
- The 2025 Illustrator JavaScript, AppleScript, and VBScript references expose significantly more scriptable functionality than the current CLI already ships.
- That makes Phase 1 the fastest path to more user value.
- This keeps the product aligned with agent-first CLI principles:
  - typed schema first
  - raw JSON payload first
  - strong validation and hardening
  - deterministic envelopes
  - runtime introspection via `tools` and `schema`

## What We Can Add Without a Plugin

### High-value domains
- `preference.*`
  - get, set, delete typed Illustrator preferences
- `view.*`
  - inspect view state, screen mode, ruler visibility, transparency-grid visibility
- `matrix.*`
  - identity, rotation, scale, translation, concatenate, invert, compare, singular
- `swatch.*`
  - list, create, delete
- `spot.*`
  - list, create, delete
- `style.character.*`
  - list, apply, import
- `style.paragraph.*`
  - list, apply, import
- `symbol.*`
  - list, create, place, break link
- `placed.*`
  - list, place, embed, relink, trace
- `raster.*`
  - list, rasterize, trace, release tracing
- `repeat.grid.*`
  - create, update
- `repeat.radial.*`
  - create, update
- `repeat.symmetry.*`
  - create, update
- `dataset.*`
  - list, create, apply, update, delete, import, export
- `variable.*`
  - list, create, delete, bind visibility, bind text, import, export
- `capture.*`
  - image capture / window capture
- `print.*`
  - list presets, run print with typed options
- richer `export.*`
  - GIF, PNG8, TIFF, WebP, Photoshop, AutoCAD, EPS, FXG

### Additional host/document surfaces that must also be completed in Phase 1
- `workspace.*`
  - save, reset, switch, delete
- richer `app.*`
  - undo, redo, redraw, copy, cut, paste, beep, convert_sample_color
- richer `document.*`
  - image capture, cloud-document operations where scriptable on macOS, library/export helpers where scriptable
- `perspective.*`
  - show, hide, select preset, get/set active plane, grid configuration where scriptable
- richer `selection.*`
  - active-artboard selection flows and normalization for text-range/insertion-point selection cases
- richer `path.*`
  - polygon, star, rounded rectangle, and efficient whole-path geometry setters where scriptable
- richer `text.*`
  - area/point/path text conversions, threading, range/paragraph/character-level operations where safely representable in schema

### Runtime corrections discovered during implementation
- `workspace.list` should not be treated as part of Phase 1. The Adobe scripting references expose workspace save/switch/reset/delete, but not a documented workspace enumeration API.
- `export.for_screens` should not be treated as part of the agent-first Phase 1 contract yet. A hidden live `document.exportForScreens` method exists, but current `do javascript` probing was not deterministic enough for safe CLI exposure.

## Feature-Complete Definition
Phase 1 is feature-complete only when all of the following are true:
- Every Illustrator scripting capability from the local JavaScript and AppleScript references that is relevant to an agent-first CLI and does not require the native SDK/plugin has a typed command family in `ahd-illustrator`.
- Those commands are exposed through:
  - `tools`
  - `schema`
  - `schema --all --llms-txt`
  - `examples`
- All commands return stable machine-readable envelopes.
- All mutating commands support `--dry-run`.
- All large-result commands support field filtering or similarly bounded output.
- The CLI normalizes Illustrator host quirks instead of leaking them:
  - coordinate-system inconsistencies
  - selection shape inconsistencies
  - export extension quirks
  - clipboard foreground requirements
  - tracing async behavior
- Live integration tests exist for each major command family, not just a representative subset.
- Docs and skill docs describe the expanded script-only surface accurately.

### Hardening and runtime normalization required
- Treat points as the canonical geometry unit.
- Normalize coordinate-system quirks at the CLI boundary.
- Normalize selection shape differences:
  - JS `null`
  - AppleScript `{}`
  - VBScript `Empty`
  - text insertion-point / text-range cases
- Make export behavior format-aware:
  - extension rules
  - artboard targeting
  - preset compatibility
- Frontmost-gate clipboard commands because AppleScript requires Illustrator to be active.
- Make tracing redraw-aware and timeout-aware.
- Validate color payloads against document color space.

## Scope

### In scope
- Exhaustive schema-backed commands for scriptable surfaces
- Input hardening for all new commands
- Live macOS integration coverage using scratch documents and temp outputs
- Expanded docs and examples
- Better `tools` / `schema --all --llms-txt` output

### Out of scope
- Illustrator SDK-native plugin implementation
- CEP/panel UI
- Windows host implementation
- Deferring major scriptable command families to a later “cleanup” phase

## Files Likely To Change
- `cmd/ahd-illustrator/main.go`
- `internal/commonschema/*`
- `internal/commonvalidate/*`
- `internal/illustrator/schema/v0.go`
- `internal/illustrator/commands/builders.go`
- `internal/illustrator/commands/executor.go`
- `internal/illustrator/commands/live_integration_test.go`
- `docs/illustrator/*`
- `skills/ahd-illustrator/SKILL.md`

## Constraints
- Do not edit these unrelated dirty files:
  - `internal/designlint/lint.go`
  - `internal/tools/catalog_llm.go`
  - `plugin/src/handlers/effect.ts`
  - `plugin/src/handlers/shape.ts`
- Keep stable response envelopes.
- Keep `--dry-run`, `--fields`, and NDJSON support intact.
- Never weaken identifier/path hardening for convenience.
- Do not regress the existing script-backed command set.

## Acceptance Criteria
- The CLI exposes the full scriptable command surface with typed schemas, not a prioritized subset.
- `tools` and `schema` fully describe the expanded surface.
- Live integration tests pass for all major script-backed command families.
- Standard users still need no SDK and no plugin.

## Completion Gate
Do not mark Phase 1 done if any major documented scripting family remains only as:
- a TODO
- a stub
- a placeholder schema
- a partial handler that does not execute live
- a docs-only promise without implementation

## Suggested Work Order
1. Add `preference.*`, `view.*`, and `matrix.*`.
2. Add `swatch.*`, `spot.*`, and style commands.
3. Add `placed.*`, `raster.*`, and tracing commands.
4. Add repeat commands.
5. Add dataset/variable commands.
6. Add richer export/capture/print support.
7. Expand docs, examples, and llms output.
8. Run full Go tests and live Illustrator integration tests.

## AI Coder Prompt
```text
You are in implementation mode in the ai-happy-design monorepo root.

Work only in:
/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design

Implement Phase 1 from:
/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/docs/plans/2026-03-05-ahd-illustrator-phase-1-script-first-expansion.md

Context:
- This repo already has a working script-first ahd-illustrator CLI.
- The local Adobe scripting references show we can add significantly more without any native plugin or SDK.
- Follow an agent-first CLI approach:
  - schema and tools introspection are canonical
  - raw JSON payloads first
  - deterministic envelopes
  - strict input hardening against hallucinations and malformed targets
  - dry-run support on mutations

Requirements:
1. Implement the Phase 1 script-first expansion end to end without pausing unless blocked.
2. Do not edit unrelated dirty files:
   - internal/designlint/lint.go
   - internal/tools/catalog_llm.go
   - plugin/src/handlers/effect.ts
   - plugin/src/handlers/shape.ts
3. Start with the suggested work order in the phase doc.
4. Preserve stable CLI envelopes and existing command behavior.
5. Add or update tests for schema, validation, and live Illustrator flows.
6. Update docs and skill docs for the new script-backed surfaces.
7. Run relevant tests and report pass/fail clearly.
8. At the end, print:
   - files changed
   - test results
   - remaining risks

Implementation bar:
- Favor typed schemas over bespoke flags.
- Favor shared validation utilities over per-command ad hoc parsing.
- Normalize Illustrator host quirks at the CLI boundary instead of leaking them to callers.
- Keep the CLI usable without any plugin or SDK.
```
