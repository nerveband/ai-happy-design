# Auto-Normalize Batch Input Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development to implement this plan task-by-task.

**Goal:** Make the CLI automatically normalize LLM output on every `batch` run and `bulk.execute` MCP call — no separate validate step required, no model-tier distinction needed.

**Architecture:** Move the fix/normalize logic from `cmd/ai-happy-design/validate.go` into a shared `internal/batchutil/fix.go` package. Wire it into `loadBatchOperations()` in `main.go` (covers all CLI batch paths) and into `bulk.execute` in `internal/tools/bulk.go` (covers MCP path). Both paths auto-fix silently on success and print fix summaries to stderr so the model can learn. Fatal errors (missing command field, unparseable JSON) still fail hard. Soft warnings (broken interpolation refs) print to stderr but don't block execution — they'll fail at runtime anyway with a clear error. Add `--no-fix` flag for users who want strict mode.

**Tech Stack:** Go, cobra, existing `internal/batchutil` package

---

## Background: Current Architecture

- `fixBatchOps()` and `stripMarkdownFences()` live in `cmd/ai-happy-design/validate.go` (package `main`) — inaccessible from `internal/`
- `loadBatchOperations()` in `main.go:1095` reads raw bytes then calls `json.Unmarshal` directly — no normalization
- `bulk.execute` in `internal/tools/bulk.go:47` does `json.Unmarshal([]byte(opsStr), &ops)` directly — no normalization
- Auto-fix is currently only triggered by `ai-happy-design validate --fix` — an explicit opt-in step
- Goal: auto-fix becomes the default in every execution path

## What Auto-Fix Does (from existing validate.go)

1. **Strip markdown fences** — removes ` ```json ``` ` wrappers models add
2. **Unwrap dict wrapper** — `{"ops": [...]}` or any `{"key": [...]}` → bare `[...]`
3. **Rename `"type"` → `"command"`** — most common LLM schema drift
4. **Hoist top-level design props into `"params"`** — x, y, color, width, height, parentId, etc.

What it does NOT fix (reported as warnings/errors):
- Missing `"command"` field entirely
- Broken `${{steps.X.result.id}}` references where X doesn't match a step name

---

## Task 1: Move fix logic to `internal/batchutil/fix.go`

**Files:**
- Create: `internal/batchutil/fix.go`
- Modify: `cmd/ai-happy-design/validate.go` (use batchutil version)

**Step 1: Write the failing test**

File: `internal/batchutil/fix_test.go`

```go
package batchutil_test

import (
	"testing"

	"github.com/nerveband/ai-happy-design/internal/batchutil"
)

func TestStripMarkdownFences(t *testing.T) {
	in := "```json\n[{\"command\":\"frame\",\"params\":{}}]\n```"
	out := batchutil.StripMarkdownFences([]byte(in))
	if string(out) != "[{\"command\":\"frame\",\"params\":{}}]" {
		t.Fatalf("got %q", string(out))
	}
}

func TestFixBatchOps_TypeToCommand(t *testing.T) {
	in := []byte(`[{"type":"frame","params":{}}]`)
	fixed, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) == 0 {
		t.Fatal("expected fixes")
	}
	// Result must be valid JSON with "command" not "type"
	if string(fixed) == string(in) {
		t.Fatal("data unchanged")
	}
}

func TestFixBatchOps_HoistTopLevelProps(t *testing.T) {
	in := []byte(`[{"command":"frame","x":0,"y":0,"params":{}}]`)
	fixed, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) == 0 {
		t.Fatal("expected hoisting fix")
	}
	_ = fixed
}

func TestFixBatchOps_UnwrapDict(t *testing.T) {
	in := []byte(`{"ops":[{"command":"frame","params":{}}]}`)
	fixed, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) == 0 {
		t.Fatal("expected unwrap fix")
	}
	if len(fixed) == 0 {
		t.Fatal("empty output")
	}
}

func TestFixBatchOps_NoChanges(t *testing.T) {
	in := []byte(`[{"name":"bg","command":"frame","params":{"x":0,"y":0}}]`)
	_, fixes, err := batchutil.FixBatchOps(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(fixes) != 0 {
		t.Fatalf("expected no fixes but got: %v", fixes)
	}
}
```

**Step 2: Run to confirm it fails**

```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design"
go test ./internal/batchutil/ -run TestStripMarkdown -v
```
Expected: FAIL (StripMarkdownFences undefined)

**Step 3: Implement `internal/batchutil/fix.go`**

```go
package batchutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// knownTopLevelProps are design properties that belong inside "params", not at op root.
var knownTopLevelProps = []string{
	"x", "y", "width", "height", "w", "h",
	"color", "fillColor", "bg",
	"fontSize", "fontFamily", "fontStyle", "sz", "ff", "fs",
	"parentId", "pid",
	"cornerRadius", "r",
	"layoutMode", "itemSpacing", "padding", "opacity",
	"text", "imageData", "stroke", "strokeWidth",
}

// StripMarkdownFences removes ```json ... ``` or ``` ... ``` wrappers that
// models add even when explicitly told not to.
func StripMarkdownFences(data []byte) []byte {
	s := strings.TrimSpace(string(data))
	if !strings.HasPrefix(s, "```") {
		return data
	}
	lines := strings.Split(s, "\n")
	var filtered []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if i == 0 && strings.HasPrefix(trimmed, "```") {
			continue
		}
		if i == len(lines)-1 && trimmed == "```" {
			continue
		}
		filtered = append(filtered, line)
	}
	return []byte(strings.Join(filtered, "\n"))
}

// FixBatchOps applies auto-corrections to common LLM output drift.
// Returns: fixed JSON bytes, list of human-readable fix descriptions, error.
// Error is non-nil only if the input cannot be parsed at all.
func FixBatchOps(data []byte) ([]byte, []string, error) {
	data = StripMarkdownFences(data)

	// Unwrap {"ops": [...]} or any single-key dict wrapping an array
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, err
	}
	var fixes []string
	if obj, isObj := raw.(map[string]interface{}); isObj {
		for k, v := range obj {
			if arr, isArr := v.([]interface{}); isArr {
				fixes = append(fixes, fmt.Sprintf("unwrapped dict key %q to get ops array", k))
				b, _ := json.Marshal(arr)
				data = b
				break
			}
		}
	}

	var ops []map[string]interface{}
	if err := json.Unmarshal(data, &ops); err != nil {
		return nil, fixes, err
	}
	for i, op := range ops {
		label := fmt.Sprintf("op[%d]", i)
		if name, ok := op["name"].(string); ok && name != "" {
			label = fmt.Sprintf("op[%d] %q", i, name)
		}

		// Fix "type" → "command"
		if typeVal, hasType := op["type"]; hasType {
			if _, hasCmd := op["command"]; !hasCmd {
				op["command"] = typeVal
				fixes = append(fixes, label+`: renamed "type" to "command"`)
			}
			delete(op, "type")
		}

		// Ensure params exists
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			params = map[string]interface{}{}
		}

		// Hoist known top-level design props into params
		var hoisted []string
		for _, prop := range knownTopLevelProps {
			if val, ok := op[prop]; ok {
				params[prop] = val
				delete(op, prop)
				hoisted = append(hoisted, prop)
			}
		}
		if len(hoisted) > 0 {
			fixes = append(fixes, fmt.Sprintf(`%s: moved %s into "params"`, label, strings.Join(hoisted, ", ")))
		}

		op["params"] = params
		ops[i] = op
	}

	fixed, err := json.MarshalIndent(ops, "", "  ")
	if err != nil {
		return nil, fixes, err
	}
	// Preserve ${{...}} interpolation tokens — json.Marshal HTML-escapes < > & by default
	fixed = bytes.ReplaceAll(fixed, []byte(`\u0026`), []byte(`&`))
	fixed = bytes.ReplaceAll(fixed, []byte(`\u003c`), []byte(`<`))
	fixed = bytes.ReplaceAll(fixed, []byte(`\u003e`), []byte(`>`))
	return fixed, fixes, nil
}
```

**Step 4: Run tests**

```bash
go test ./internal/batchutil/ -v
```
Expected: all 5 new tests PASS, existing batchutil tests still PASS.

**Step 5: Update validate.go to delegate to batchutil**

In `cmd/ai-happy-design/validate.go`, replace the local `stripMarkdownFences` and `fixBatchOps` functions with thin wrappers that call `batchutil.StripMarkdownFences` and `batchutil.FixBatchOps`. Remove the `knownTopLevelProps` var (it moves to batchutil). Keep the `validateBatchOps` function and `interpolationRef` regex — those stay in validate.go.

Replace:
```go
var knownTopLevelProps = []string{...}  // DELETE this
func stripMarkdownFences(data []byte) []byte {...}  // DELETE this
func fixBatchOps(data []byte) ([]byte, []string, error) {...}  // DELETE this
```

Add import: `"github.com/nerveband/ai-happy-design/internal/batchutil"`

In the `--fix` path of `validateCmd.RunE`, change:
```go
fixed, fixes, fixErr := fixBatchOps(data)
```
to:
```go
fixed, fixes, fixErr := batchutil.FixBatchOps(data)
```

`validateBatchOps` still uses `knownTopLevelProps` locally — add a local copy or reference `batchutil.KnownTopLevelProps` if exported. Easiest: export `KnownTopLevelProps` from `batchutil/fix.go` and use it in validate.go.

Add to `internal/batchutil/fix.go`:
```go
// KnownTopLevelProps exposes the list for use in validate.go
var KnownTopLevelProps = knownTopLevelProps
```

**Step 6: Run all validate tests**

```bash
go test ./cmd/ai-happy-design/ -v -run TestValidate
```
Expected: all 14 existing validate tests PASS.

**Step 7: Build**

```bash
go build ./...
```
Expected: no errors.

**Step 8: Commit**

```bash
git add internal/batchutil/fix.go internal/batchutil/fix_test.go cmd/ai-happy-design/validate.go
git commit -m "refactor: move fix/normalize logic to internal/batchutil

Extracts StripMarkdownFences, FixBatchOps, knownTopLevelProps from
cmd/validate.go into internal/batchutil/fix.go so they can be reused
by the batch command and MCP bulk.execute handler.

validate.go updated to delegate to batchutil. All 14 validate tests pass.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 2: Auto-fix in CLI batch (`loadBatchOperations`)

**Files:**
- Modify: `cmd/ai-happy-design/main.go` — `loadBatchOperations` function (line ~1095) and flag init (line ~1575)

**Step 1: Write the failing test**

There are no unit tests for `loadBatchOperations` (it's in package main). We'll test via the existing test infrastructure and a manual smoke test. Skip to implementation.

**Step 2: Add `--no-fix` flag**

In the flags block near line 1575:

```go
var batchNoFix bool
// ...
batchCmd.Flags().BoolVar(&batchNoFix, "no-fix", false, "Skip automatic LLM output normalization (use if your JSON is already clean)")
```

**Step 3: Update `loadBatchOperations` to auto-fix**

Current code at line 1123:
```go
var ops []batchOperation
if err := json.Unmarshal(raw, &ops); err != nil {
    return nil, fmt.Errorf("invalid operations JSON: %w", err)
}
return ops, nil
```

Replace with:
```go
// Auto-normalize LLM output unless --no-fix is set
if !batchNoFix {
    fixed, fixes, fixErr := batchutil.FixBatchOps(raw)
    if fixErr != nil {
        return nil, fmt.Errorf("cannot parse operations JSON: %w", fixErr)
    }
    if len(fixes) > 0 {
        fmt.Fprintf(os.Stderr, "✎ Auto-normalized %d issue(s) in batch input:\n", len(fixes))
        for _, f := range fixes {
            fmt.Fprintf(os.Stderr, "  • %s\n", f)
        }
    }
    raw = fixed
}

var ops []batchOperation
if err := json.Unmarshal(raw, &ops); err != nil {
    return nil, fmt.Errorf("invalid operations JSON: %w", err)
}
return ops, nil
```

**Step 4: Smoke test — auto-fix in batch**

```bash
# Test 1: markdown fences get stripped automatically
echo '```json
[{"name":"bg","command":"frame","params":{"x":0,"y":0,"w":100,"h":100}}]
```' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch -
```
Expected: stderr shows "✎ Auto-normalized 1 issue(s)" and batch runs (or fails with plugin error, not parse error).

```bash
# Test 2: type→command auto-fix
echo '[{"name":"bg","type":"frame","params":{"x":0,"y":0,"w":100,"h":100}}]' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch -
```
Expected: stderr shows renamed "type" to "command", batch runs.

```bash
# Test 3: --no-fix rejects fenced input
echo '```json
[{"name":"bg","command":"frame","params":{}}]
```' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch --no-fix -
```
Expected: exits with error "invalid operations JSON".

**Step 5: Build**

```bash
go build ./...
```

**Step 6: Commit**

```bash
git add cmd/ai-happy-design/main.go
git commit -m "feat: auto-normalize LLM output in batch command

batch now auto-applies FixBatchOps before executing:
- strips markdown fences
- unwraps dict wrappers like {\"ops\":[...]}
- renames \"type\" to \"command\"
- hoists top-level design props (x, y, color, etc.) into \"params\"

Fixes are printed to stderr so the model can see what was corrected.
Use --no-fix to opt out (strict mode for clean JSON).

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 3: Auto-fix in MCP `bulk.execute`

**Files:**
- Modify: `internal/tools/bulk.go` — the `execute` case (line ~41)

**Step 1: Add auto-fix before Unmarshal**

Current code at line 46:
```go
var ops []BulkOperation
if err := json.Unmarshal([]byte(opsStr), &ops); err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("invalid operations JSON: %v", err)), nil
}
```

Replace with:
```go
opsBytes, fixes, fixErr := batchutil.FixBatchOps([]byte(opsStr))
if fixErr != nil {
    return mcp.NewToolResultError(fmt.Sprintf("cannot parse operations JSON: %v", fixErr)), nil
}

var ops []BulkOperation
if err := json.Unmarshal(opsBytes, &ops); err != nil {
    return mcp.NewToolResultError(fmt.Sprintf("invalid operations JSON: %v", err)), nil
}
```

Add import at top of bulk.go:
```go
"github.com/nerveband/ai-happy-design/internal/batchutil"
```

If fixes were applied, include them in the result. Add fixes to the `out` map at the summary:
```go
if len(fixes) > 0 {
    out["autoNormalized"] = fixes
}
```

This surfaces fixes in the MCP response so the LLM can see what was corrected.

**Step 2: Build**

```bash
go build ./...
```

**Step 3: Commit**

```bash
git add internal/tools/bulk.go
git commit -m "feat: auto-normalize LLM output in MCP bulk.execute

bulk.execute now applies FixBatchOps before parsing operations JSON.
Fixes applied are returned in the response as autoNormalized:[...].
This ensures MCP path handles the same LLM drift as the CLI batch path.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 4: Update SKILL.md and catalog_llm.go

**Files:**
- Modify: `skills/ai-happy-design/SKILL.md`
- Modify: `internal/tools/catalog_llm.go`

### SKILL.md changes

**Step 1: Update Workflow section**

Current step 6: `"Validate & fix: ai-happy-design validate --fix ops.json — auto-corrects schema drift..."`

Replace with (batch is now self-normalizing):
```markdown
6. **Send batch**: `ai-happy-design batch ops.json` — auto-normalizes LLM output before sending (fences, `type`→`command`, top-level props). Use `--no-fix` for strict mode.
```

Remove the now-redundant step, renumber 7→6 (balance) and 8→7 (export).

**Step 2: Update Batch Operations section**

In the `### Chunking Limit` section, replace:
```markdown
Use `ai-happy-design validate --fix` to auto-correct and validate before sending.
```
With:
```markdown
The `batch` command auto-normalizes input. Use `ai-happy-design validate --fix` to inspect what would be corrected without running.
```

**Step 3: Update Key CLI Commands section**

Update the validate lines to clarify their new role:
```bash
ai-happy-design validate batch.json          # inspect batch schema (batch auto-normalizes, this is for debugging)
ai-happy-design validate --fix batch.json    # inspect + preview what auto-fix would do
ai-happy-design batch --no-fix batch.json    # strict mode — no auto-normalization
```

**Step 4: Update Common Pitfalls**

Replace:
```markdown
- **Always run `validate --fix` on model output** before `batch` — catches and fixes common drift automatically.
```
With:
```markdown
- **`batch` auto-normalizes**: fences, `type`→`command`, top-level params are all corrected automatically. Use `validate` to inspect; use `--no-fix` to disable.
```

### catalog_llm.go changes

**Step 5: Update `workflow.create`**

Update the `how` string to remove the explicit validate step (it's now automatic):
```
4) Write ops to a file or inline. 5) Send: 'ai-happy-design batch ops.json' — auto-normalizes LLM output (fences, type→command, top-level params) before executing. Use --no-fix to disable.
```

**Step 6: Update `validateBeforeBatch` entry**

Update to reflect the new reality:
```go
"validateBeforeBatch": map[string]interface{}{
    "builtin": "batch auto-applies all normalization — no separate validate step needed.",
    "inspect": "Use 'ai-happy-design validate --fix ops.json' to preview corrections without running.",
    "strict":  "Use 'ai-happy-design batch --no-fix ops.json' to disable auto-normalization.",
    "what":    "Auto-corrected: markdown fences, type→command, top-level props into params, dict wrapper unwrap. MCP bulk.execute also auto-normalizes.",
},
```

**Step 7: Update quickPrompts**

Replace:
```
"VALIDATE ALWAYS: Run 'ai-happy-design validate --fix ops.json' before every batch..."
```
With:
```
"BATCH AUTO-NORMALIZES: batch and bulk.execute both auto-fix LLM output (fences, type→command, top-level params). Just write JSON and run batch — no separate validate step needed."
```

**Step 8: Build**

```bash
go build ./...
```

**Step 9: Commit**

```bash
git add skills/ai-happy-design/SKILL.md internal/tools/catalog_llm.go
git commit -m "docs: update skill to reflect auto-normalize in batch/bulk

batch and bulk.execute now auto-normalize LLM output.
SKILL.md workflow simplified: no explicit validate step needed.
validate --fix remains for inspection/debugging.
--no-fix flag documented for strict mode.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 5: Deploy and full smoke test

**Step 1: Run all tests**

```bash
go test ./... 2>&1
```
Expected: all pass.

**Step 2: Deploy**

```bash
make deploy
```

**Step 3: Full smoke test matrix**

```bash
# 1. Fenced input auto-fixed
echo '```json
[{"name":"t","command":"document.get_info","params":{}}]
```' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch -
# Expected: stderr shows auto-normalized, command runs

# 2. type→command auto-fixed
echo '[{"name":"t","type":"document.get_info","params":{}}]' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch -
# Expected: stderr shows renamed, command runs

# 3. top-level x/y auto-hoisted
echo '[{"name":"t","command":"document.get_info","params":{},"x":0,"y":0}]' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch -
# Expected: stderr shows hoisted x,y, command runs

# 4. dict wrapper auto-unwrapped
echo '{"ops":[{"name":"t","command":"document.get_info","params":{}}]}' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch -
# Expected: stderr shows unwrapped, command runs

# 5. --no-fix strict mode rejects fenced input
echo '```json
[{"name":"t","command":"document.get_info","params":{}}]
```' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch --no-fix -
# Expected: exits 1 with "invalid operations JSON"

# 6. clean input — no normalization message
echo '[{"name":"t","command":"document.get_info","params":{}}]' | AHD_CHANNEL=golden-waffle-25 ai-happy-design batch -
# Expected: no "Auto-normalized" message on stderr
```

**Step 4: Verify validate --fix still works as before**

```bash
echo '[{"type":"frame","x":0,"params":{}}]' | ai-happy-design validate --fix -
# Expected: still works, prints fixes, outputs corrected JSON to stdout
```

---

## Summary of Changes

| File | Change |
|------|--------|
| `internal/batchutil/fix.go` (new) | StripMarkdownFences, FixBatchOps, KnownTopLevelProps |
| `internal/batchutil/fix_test.go` (new) | 5 unit tests |
| `cmd/ai-happy-design/validate.go` | Delegate to batchutil, remove duplicate code |
| `cmd/ai-happy-design/main.go` | Auto-fix in loadBatchOperations, add --no-fix flag |
| `internal/tools/bulk.go` | Auto-fix in bulk.execute, surface fixes in response |
| `skills/ai-happy-design/SKILL.md` | Workflow simplified, batch is self-normalizing |
| `internal/tools/catalog_llm.go` | Update workflow + quickPrompts |

**No plugin changes needed.** All normalization is in the Go layer before commands reach the plugin.
