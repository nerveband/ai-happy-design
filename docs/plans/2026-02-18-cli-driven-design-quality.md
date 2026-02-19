# CLI-Driven Design Quality — Revised Execution Plan (v2)

Date: 2026-02-19
Owner: CLI + Plugin Runtime
Status: Approved for implementation

## Goal

Make `ai-happy-design` produce better visual output by default for LLM-driven workflows, without relying on external docs or perfect prompt following.

The CLI output from each step should contain enough context for the next step.

## Success Criteria

1. LLM can generate a clean social design in one `batch` run with no manual token math.
2. The system catches common quality regressions (tiny text, overflow, generic names) automatically.
3. Guidance is embedded in command output (`compute_tokens`, batch warnings, examples), not external documents.
4. Existing numeric batch payloads continue to work unchanged.
5. No protocol contract changes and no plugin build target regressions.

## Reality Check: What Is Good vs. Weak in the Original Plan

### Keep (high ROI)

1. Token aliases resolved by CLI (`sz:"hero"`, `padding:"side"`, `w:"content"`).
2. `compute_tokens` returning starter template + actionable tips.
3. Post-batch lint surfaced as actionable warnings.
4. Better diagnostics for interpolation and command/action mistakes.
5. Stronger examples that teach by copy-paste.

### Remove / Defer (low ROI or high risk right now)

1. Full scaffold command (`scaffold ...`) now: deferred.
Reason: duplicates template system, adds maintenance/API complexity before core flow is stable.

2. Aggressive silent auto-fixes in a broad preprocessor.
Reason: hidden mutations can violate user intent and make debugging harder.
Replacement: deterministic alias resolution + explicit warnings.

3. Large one-shot examples expansion (6-8 new categories) in first pass.
Reason: content-heavy work can block core engine improvements.
Replacement: upgrade the highest-leverage examples first.

4. SKILL.md reduction as a deliverable.
Reason: not required to improve runtime quality; quality should come from CLI behavior.

## Design Principles for Implementation

1. Deterministic over magical: resolve only known aliases, no speculative rewrites.
2. Explicit warnings over silent correction: tell the user what is wrong and how to fix it.
3. Backward compatible: numeric values and existing commands remain valid.
4. Minimal surface-area change: improve current pipeline instead of adding new command families.
5. LLM-usable output: short, concrete, copyable fix hints.

## Final Scope (This Iteration)

### Phase 1: Core Behavior (must ship)

1. Add token alias resolver in batch pipeline.
2. Add token alias resolver in MCP `bulk.execute` pipeline for parity.
3. Make lint default-on for batch runs; add `--no-lint` escape hatch.
4. Improve batch lint output format with severity filtering and actionable fix snippets.
5. Improve interpolation failure messages with available step names and likely correction.

### Phase 2: Guidance Surfaces (must ship)

1. Extend `design.compute_tokens` output with:
- `template`: starter operations tuned by aspect ratio.
- `tips`: concise next-step guidance.
- `aliases`: quick reference with canonical names.
2. Upgrade `examples batch` and one multi-frame example (`carousel`) to use aliases + stronger structure.

### Phase 3: Hardening (must ship)

1. Unit tests for alias resolution behavior and edge cases.
2. Unit tests for improved interpolation diagnostics.
3. Build + test verification for Go and plugin checks.
4. Live Figma integration validation against connected plugin channel.

## Detailed Implementation Plan

### 1) Token Alias Resolver

New module: `internal/batchutil/token_resolve.go`

Supported aliases:

- Text size: `display|hero|title|heading|subheading|body|caption|numbers|cta`
- Spacing: `side|frame|card|section|item|tight`
- Radius: `card|button|pill`
- Width: `content`

Supported fields:

- `fontSize`, `sz`
- `padding`, `paddingTop`, `paddingRight`, `paddingBottom`, `paddingLeft`
- `itemSpacing`, `gap`
- `cornerRadius`, `r`
- `width`, `w`

Resolver behavior:

1. Detect root canvas width from first root frame create op.
2. Compute token table from existing `ComputeDesignTokens` math.
3. Replace only exact string aliases in supported fields.
4. Preserve unknown strings untouched.
5. Keep numeric values untouched.
6. For `pill`, resolve to large radius (`9999`) for shape-safe pill behavior.

Integration points:

1. CLI `batch`: after composite expansion and before per-op execution.
2. MCP `bulk.execute`: before per-op execution.

### 2) Batch Lint Default-On

Changes:

1. `--lint` becomes default `true`.
2. Add `--no-lint` flag to disable.
3. Run lint only when execution had no failed steps.
4. Report warning count with severity icons and compact lines.
5. Show up to N warnings (cap to avoid token noise), then summarize remainder.

Output style:

- Short heading.
- One line per warning.
- Optional copyable `node.modify` hint when useful (rename case).

### 3) Better Diagnostics

Interpolation errors:

1. If step path missing, include available step names.
2. If a case-insensitive match exists, include exact corrected reference.
3. Keep error machine-readable via existing `errorCode` flow.

Unknown command/action diagnostics:

1. Keep plugin-side available-actions list in unknown-action messages.
2. Improve CLI-side command routing errors with nearest known command when possible.

### 4) `compute_tokens` Guidance Upgrade

Add fields without removing existing keys:

1. `template`: starter batch array using token aliases.
2. `tips`: 3-5 short, actionable constraints.
3. `aliases`: quick map of alias groups and examples.

Aspect-ratio presets:

1. Portrait: centered vertical stack.
2. Square: centered stack with balanced spacing.
3. Landscape: split layout template.

### 5) Example Upgrade (Targeted)

Update examples with highest usage impact:

1. `batch`: canonical minimal starter using aliases.
2. `carousel`: cleaner hierarchy and naming for multi-slide behavior.

Keep current categories; improve quality of existing payloads instead of adding many new categories now.

## Out of Scope (This Iteration)

1. LLM-powered scaffold generation.
2. Heavy auto-correction that modifies user intent beyond deterministic aliases.
3. Expanding all examples categories in one pass.

## Test Plan

### Unit Tests

1. `token_resolve_test.go`
- Resolves all supported aliases for 1080 root width.
- Leaves unknown aliases untouched.
- Leaves numeric values untouched.
- No root frame => no changes.
- Multi-root => first root width wins.

2. `interpolation_test.go` additions
- Missing step reports available names.
- Case mismatch suggests corrected step key.

### Integration / Build

1. `go test ./...`
2. `go build ./...`
3. `cd plugin && npm run check`
4. Validate ES6 output constraints remain satisfied.

### Live Figma Validation (Connected Plugin)

1. Run `design.compute_tokens` and verify `template/tips/aliases` fields exist.
2. Run a batch using aliases only (`sz`, `padding`, `w:"content"`) and confirm execution success.
3. Intentionally create one lint issue (e.g., default name) and verify warning output.
4. Export a created frame and confirm output file is generated.

## Acceptance Checklist

1. Alias-only batch payload executes correctly in live Figma session.
2. Lint runs by default and can be disabled with `--no-lint`.
3. `compute_tokens` returns starter template and actionable tips.
4. No regressions in existing batch numeric payloads.
5. All tests and builds pass.

## Rollout / Risk Mitigation

1. All additions are additive and backward compatible.
2. Alias resolver is limited to explicit fields and exact-match strings.
3. Lint output is capped to prevent excessive stderr noise.
4. Keep changes behind existing command paths; no schema break.

## Implementation Order (Exact)

1. Implement alias resolver + tests.
2. Wire resolver into CLI batch and MCP bulk.
3. Flip lint default and improve lint formatting.
4. Improve interpolation diagnostics.
5. Extend `compute_tokens` with `template/tips/aliases`.
6. Refresh `examples batch` + `examples carousel`.
7. Run full test/build suite.
8. Run live Figma validation and document evidence.
