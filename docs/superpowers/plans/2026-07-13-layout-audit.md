# Layout Audit Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a read-only `layout.audit` command that measures a Figma subtree, reports actionable layout issues, and gives bounded repair commands without relying on screenshots or guesswork.

**Architecture:** Add pure geometry/audit helpers in a focused plugin utility, expose a new `layout.audit` action through the existing layout handler, and keep all recommendations as data-only command plans. Use a temporary text clone only for measurement and remove it before returning. Register the command in Go schema and LLM descriptions so agents discover the audit-first workflow.

**Tech Stack:** Go schema/catalog, TypeScript Figma plugin, Vitest, existing `getNodeById` utilities, existing layout/document lint conventions.

---

## Chunk 1: Pure geometry and audit contracts

**Files:**
- Create: `plugin/src/utils/layoutAudit.ts`
- Create: `plugin/src/utils/layoutAudit.test.ts`

- [ ] **Step 1: Write failing tests for rectangle geometry**
  - Cover intersection area/rectangle, containment, edge-touching as non-overlap, and minimum gap calculation.
  - Use plain rectangles only so these tests run without Figma globals.

- [ ] **Step 2: Run the geometry tests and verify they fail for the expected missing exports**

  Run: `cd plugin && npm test -- src/utils/layoutAudit.test.ts`

- [ ] **Step 3: Implement pure geometry helpers**
  - Define `Bounds`, `Point`, `LayoutIssue`, `FixRecommendation`, and confidence/severity unions.
  - Implement `intersection`, `contains`, `gapBetween`, `expand`, and safe numeric helpers.
  - Treat touching edges as zero overlap.

- [ ] **Step 4: Run the geometry tests and verify they pass**

- [ ] **Step 5: Add failing tests for recommendation builders**
  - Verify text overflow produces a minimum resize recommendation.
  - Verify sibling overlap recommends moving the later node by the measured amount plus the requested gap.
  - Verify every recommendation includes node IDs, strategy, and concrete command objects.

- [ ] **Step 6: Implement recommendation builders and rerun tests**

- [ ] **Step 7: Commit the pure audit utility**

---

## Chunk 2: Figma subtree traversal and measurements

**Files:**
- Modify: `plugin/src/handlers/layout.ts`
- Modify: `plugin/src/utils/layoutAudit.ts`
- Create: `plugin/src/handlers/layoutAudit.test.ts`

- [ ] **Step 1: Write failing mocked-plugin tests for recursive traversal**
  - Build a mock frame with nested children and verify all descendants are visited up to the requested depth.
  - Verify the root itself is included in the audit context.
  - Verify missing/invalid `nodeId` errors are useful.

- [ ] **Step 2: Add mocked text fixtures for the real regression cases**
  - A fixed-size text node whose declared height is smaller than its natural height.
  - A long title whose natural width exceeds its declared width.
  - Two neighboring text nodes with intersecting rendered bounds.
  - A childcare/outdoor-like pair where the background box is smaller than the text measurement.

- [ ] **Step 3: Run the handler tests and verify they fail**

- [ ] **Step 4: Implement safe node measurement**
  - Load the requested node through `getNodeById`.
  - Read `absoluteBoundingBox` when present; otherwise accumulate local coordinates from parents.
  - Capture node ID, name, type, parent ID, bounds, layout mode, and text sizing properties.
  - For text nodes with fixed sizing, clone temporarily, set natural height sizing, read the measured bounds, and remove the clone in a `finally` block.
  - Never append persistent nodes or return a result before temporary clones are removed.
  - Use low confidence and an explanatory evidence field when clone measurement is unavailable.

- [ ] **Step 5: Implement recursive parent overflow and clipping checks**
  - Detect child bounds outside parent bounds.
  - Detect text natural width/height outside its usable bounds.
  - Include actual, available, delta, and confidence evidence.

- [ ] **Step 6: Implement sibling overlap and tight-gap checks**
  - Compare siblings in the same parent using rendered bounds.
  - Report intentional absolute-positioning separately rather than silently ignoring it.
  - Use a default minimum gap of 4px for direct neighboring content and allow `minGap` override.
  - Avoid reporting duplicate parent/child findings for the same collision.

- [ ] **Step 7: Implement manual-layout and auto-layout risk checks**
  - Flag non-auto-layout containers with multiple manually positioned content children as an informational risk.
  - Flag auto-layout children with fixed sizing that prevents expected reflow.
  - Do not claim auto-layout is always the correct fix; recommendation confidence must reflect structural compatibility.

- [ ] **Step 8: Implement compact and detailed response shapes**
  - Default response: `ok`, `nodeId`, `summary`, and actionable issue list.
  - `compact:true`: omit repeated geometry while retaining codes, severity, node IDs, deltas, and fixes.
  - `detailed:true`: include full bounds and measurement evidence.
  - Include `truncated`/`visitedCount` if depth or node limits are reached.

- [ ] **Step 9: Add `layout.audit` dispatch and run the handler tests**
  - Add action aliases only if they do not create ambiguity.
  - Preserve existing `layout.check_overlaps` behavior.

- [ ] **Step 10: Commit the plugin audit implementation**

---

## Chunk 3: Schema, catalog, and agent workflow

**Files:**
- Modify: `internal/schema/layout_schemas.go`
- Modify: `internal/tools/describe.go`
- Modify: `internal/tools/catalog_llm.go`
- Modify: `AGENTS.md` only if a non-rule architectural reference is needed; do not duplicate design rules
- Test: `internal/schema/*_test.go` or existing schema tests if contract coverage exists

- [ ] **Step 1: Write failing schema/discovery tests**
  - Verify `layout.audit` is registered.
  - Verify `nodeId` is required and matches the node ID pattern.
  - Verify `depth`, `compact`, `detailed`, and `minGap` are discoverable.

- [ ] **Step 2: Register the schema**
  - Description must say the command is read-only.
  - Description must require audit → intentional batch → re-audit before screenshot.
  - Document that screenshots are not the primary diagnostic signal.

- [ ] **Step 3: Add LLM action description and workflow guidance**
  - Add `layout.audit` to action descriptions with concrete example output and repair workflow.
  - Update layout guidance to prefer measurable audit evidence before visual iteration.
  - Keep the catalog as the source of design rules; do not duplicate unrelated rules elsewhere.

- [ ] **Step 4: Run Go schema/tools tests and verify discovery output**

  Run:
  - `go test ./internal/schema ./internal/tools`
  - `go run ./cmd/ai-happy-design schema layout.audit`

- [ ] **Step 5: Commit schema and agent guidance**

---

## Chunk 4: Regression, build, and live proof

**Files:**
- Modify: `docs/reference/live-figma-acceptance.md` if a concise acceptance command belongs there
- No production changes outside the audit scope

- [ ] **Step 1: Add compact-output contract tests**
  - Verify compact output is valid JSON and contains no screenshots or verbose diagnostics.
  - Verify every issue has a bounded fix strategy or an explicit `fix:null` reason.

- [ ] **Step 2: Add regression coverage for temporary measurement cleanup**
  - Assert the page child count is unchanged after a text measurement audit.
  - Assert clone/remove is performed even when measurement throws.

- [ ] **Step 3: Run complete verification**

  Run:
  - `go test ./...`
  - `cd plugin && npm run check && npm test`
  - `git diff --check`

- [ ] **Step 4: Build and install the local CLI/plugin artifacts**
  - `make build`
  - Update local `~/bin` binaries as required.
  - Reopen/reload the Figma plugin before live testing.

- [ ] **Step 5: Run live acceptance against the affected subtree**
  - `ahd-figma command layout.audit '{"nodeId":"2260:460","compact":true}'`
  - Verify the childcare/outdoor overflow and note collision are reported with evidence.
  - Apply no repair during acceptance; confirm the command is read-only.
  - Re-run the audit after any user-approved repair batch and take one final screenshot.

- [ ] **Step 6: Review diff and report remaining limitations**
  - Explicitly report cases where Figma cannot provide exact text measurement.
  - Do not claim visual correctness solely from a passing audit.
