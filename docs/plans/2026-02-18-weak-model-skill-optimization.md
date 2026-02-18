# Weak Model Skill Optimization — ai-happy-design

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make the ai-happy-design skill work reliably with smaller/faster LLMs that have weaker instruction-following, by (a) making the schema crystal-clear in SKILL.md, (b) providing a prompting reference file for agents using weaker models — validated in a real production session, and (c) adding a CLI `validate` command so the generate→validate→fix→send workflow is automated.

**Architecture:**
Three-layer approach mirroring how distilled models fail. Layer 1 (SKILL.md): rewrite the batch section to be schema-contract-first with WRONG/RIGHT examples, explicit interpolation rules, chunking limits, and an anti-patterns block. Layer 2 (new reference file): a complete prompting guide with system prompt templates, step-name manifest pattern, two-pass generation — all patterns battle-tested in a real 36-slide parallel generation session. Layer 3 (CLI): a new `validate` command that checks batch JSON against the schema contract — enabling the generate→validate→send workflow without ever hitting the Figma plugin.

**Tech Stack:** Go (CLI validate command), Markdown (SKILL.md + reference file)

---

## Root Cause Analysis

Smaller/distilled models fail on batch JSON for 5 reasons:

| Failure | Symptom | Fix |
|---------|---------|-----|
| Schema mismatch | `type` instead of `command`, params at top-level | Show WRONG/RIGHT in SKILL.md |
| Interpolation drift | `steps.background.result.id` when step is named `bg` | Pre-define step name manifest |
| Volume degradation | Schema compliance falls at op 30+ | Chunk limit ≤15 ops documented |
| Novel format | Model never saw this format in training | Show complete examples, not just schema |
| No grammar enforcement | Free-text generation drifts | CLI validate as post-generation guard |

## Production-Validated Findings

The following was confirmed in a real 36-slide parallel generation session. **These are not theories** — treat them as hard constraints:

### 1. `json_schema` strict mode does NOT work for batch ops

`response_format: { type: "json_schema", json_schema: { strict: true } }` with `additionalProperties: false` fails because `params` is an open object. The API rejects the schema definition, not the model output. **Use plain prompting + Python validation instead.** This is the working pattern.

### 2. Models output markdown fences even with explicit "no markdown" prompts

Even `"No markdown fences, no prose"` in the system prompt doesn't prevent:
````
```json
[{"name":"bg",...}]
```
````
**You MUST strip fences programmatically** — checking if content starts with `` ` `` and filtering those lines.

### 3. Large curl payloads need temp files

Shell argument escaping breaks with large JSON payloads. Use temp file pattern:
```python
import tempfile, os
with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
    json.dump(payload, f)
    tmpfile = f.name
try:
    result = subprocess.run(["curl", "-sS", endpoint,
        "-H", "Content-Type: application/json",
        "-H", f"Authorization: Bearer {api_key}",
        "-d", f"@{tmpfile}"],  # ← @filepath not inline JSON
        capture_output=True, text=True, timeout=60)
finally:
    os.unlink(tmpfile)
```

### 4. API can return error OR empty — handle both

```python
if not result.stdout.strip():
    raise RuntimeError("Empty response")
resp = json.loads(result.stdout)
if "error" in resp and "choices" not in resp:
    raise RuntimeError(f"API error: {resp.get('message', resp)}")
content = resp["choices"][0]["message"]["content"].strip()
```

### 5. Model may return `{"ops": [...]}` or bare `[...]`

Parse defensively:
```python
data = json.loads(content)
if isinstance(data, dict):
    ops = data.get("ops", list(data.values())[0] if data else [])
else:
    ops = data
```

### 6. f-string brace escaping for interpolation syntax

`${{steps.X.result.id}}` in Python f-strings needs quadruple braces:
```python
f"Reference: ${{{{steps.{step_name}.result.id}}}}"
# → renders as: Reference: ${{steps.bg.result.id}}
```

### 7. Wrap output as `{"ops": [...]}` not bare array

Ask for `{"ops": [...]}` not `[...]`. Single-key dicts are easier to extract, and the model is less likely to add commentary before or after.

### 8. `temperature: 0` is important

Schema compliance is a deterministic task. Always use `temperature: 0` for batch JSON generation.

---

## Task 1: Rewrite SKILL.md Batch Section

**Files:**
- Modify: `skills/ai-happy-design/SKILL.md`

The current batch section (lines ~122-158) has an example but lacks: WRONG/RIGHT contrast, interpolation pitfalls, chunking limits, anti-patterns. Rewrite it to be weak-model-safe.

**Step 1: Replace the `## Batch Operations` section opening**

Find (lines 122-133):
```markdown
## Batch Operations

Build a JSON array with named steps and `${{steps.name.result.id}}` interpolation:
```

Replace with:
```markdown
## Batch Operations

> **Schema Contract** — this is a custom format. Read carefully before generating.
>
> Every element MUST have exactly these three keys:
> ```json
> {"name": "stepName", "command": "domain.action", "params": {"key": "value"}}
> ```
> **WRONG** (common mistakes that break silently):
> ```json
> {"type": "node.create_frame", "x": 0, "color": "#fff"}          ← "type" not "command", params not nested
> {"command": "frame", "name": "bg", "x": 0, "y": 0}              ← x/y at top level, not inside params
> {"command": "text.create", "params": {"text": "Hi"}, "name": ""}← empty name breaks interpolation
> ```
> **RIGHT**:
> ```json
> {"name": "bg", "command": "node.create_frame", "params": {"x": 0, "y": 0, "width": 1080, "height": 1080}}
> {"name": "title", "command": "text.create", "params": {"text": "Hello", "parentId": "${{steps.bg.result.id}}", "fontSize": 84}}
> ```

Build a JSON array with named steps and `${{steps.name.result.id}}` interpolation:
```

**Step 2: After the main batch example, add interpolation rules + chunking**

After the batch example code block (after line ~132), add:

```markdown
### Interpolation Rules

Reference earlier step results by exact name:
```
${{steps.EXACT_STEP_NAME.result.id}}
```

**Rules:**
1. Step names must be unique, snake_case, no spaces: `bg`, `card_1`, `title_text`
2. You can only reference steps that appear EARLIER in the array
3. The name in `${{steps.X.result.id}}` must match the `"name"` field EXACTLY — `bg` ≠ `background` ≠ `Background`
4. **Define your step manifest before writing ops and use only those names.** Never invent a name mid-array.

### Chunking Limit

**Keep batches to ≤15 operations per file.** Schema compliance degrades on longer outputs — confirmed in production at 30+ ops.
For designs needing 30+ ops, split into multiple files and run sequentially:
```bash
ai-happy-design batch phase1.json && ai-happy-design batch phase2.json
```
Or use the batch validator before sending (see `ai-happy-design validate`).
```

**Step 3: Add anti-patterns to Common Pitfalls section**

At the end of `## Common Pitfalls` (after line ~422), add:

```markdown
- **`type` instead of `command`**: The batch format uses `"command"`, not `"type"`. Always `{"command": "node.create_frame", ...}`.
- **Top-level params**: ALL design properties go inside `"params"`: `{"command":"...","params":{"x":0,"color":"#fff"}}` not `{"command":"...","x":0,"color":"#fff"}`.
- **Step name mismatch**: If a step is named `"bg"`, reference it as `${{steps.bg.result.id}}` — not `background`, `Background`, `root`, or any variation. Mismatch fails silently.
- **Over-length batches**: Keep ≤15 ops per batch file. Schema compliance degrades at 30+ ops. Split large designs.
- **Empty or missing names**: Every op that will be referenced MUST have a `"name"` field. Ops at the end that reference nothing can omit it.
```

**Step 4: Verify line count**

```bash
wc -l "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/skills/ai-happy-design/SKILL.md"
```

Target: under ~490 lines.

**Step 5: Commit**

```bash
git add skills/ai-happy-design/SKILL.md
git commit -m "docs: make batch schema weak-model-safe

Add WRONG/RIGHT contrast, interpolation rules, step name manifest
pattern, chunking limit, and anti-patterns. Prevents schema drift
on smaller/faster distilled models.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 2: New Reference File — Weak Model Prompting Guide

**Files:**
- Create: `skills/ai-happy-design/references/weak-model-prompting.md`
- Modify: `skills/ai-happy-design/SKILL.md` (add link to reference)

This file contains everything validated in the real 36-slide session. Copy-paste these patterns directly when orchestrating a weaker model.

**Step 1: Create the reference file at `skills/ai-happy-design/references/weak-model-prompting.md`**

Full content to write:

```markdown
# Weak Model Prompting Guide

Use when delegating batch JSON generation to a smaller/faster model.
All patterns below were validated in a real parallel generation session (36 slides, 12 workers).

---

## The Correct Workflow

```
Smaller/faster model (creative, parallel)
  → generates {"ops": [...]} JSON per design
  → strip markdown fences
  → Python normalize/validate (fix common drift, reject garbage)
  → ai-happy-design validate output.json (final deterministic check)
  → ai-happy-design batch output.json (only reaches Figma if clean)
```

Do NOT send model output directly to Figma. Always normalize first.

---

## System Prompt (Validated)

```
You output ONLY a JSON object with key "ops" containing an array of operations.
No markdown fences, no prose, no explanations. Just the JSON object.

EXACT SCHEMA - every array element MUST match this shape:
  {"name": "STEPNAME", "command": "COMMAND", "params": {"key": value}}

VALID COMMANDS (short aliases): frame, rect, text, image, gradient, shadow, blur, glass, mask, modify, fill, stroke, noise

CORRECT EXAMPLES:
  {"name":"root","command":"frame","params":{"x":0,"y":0,"w":1080,"h":1080,"bg":"#14344A","name":"My Post"}}
  {"name":"content","command":"frame","params":{"x":86,"y":86,"w":908,"h":908,"pid":"${{steps.root.result.id}}","noFill":true,"layoutMode":"VERTICAL","padding":64,"itemSpacing":24}}
  {"name":"headline","command":"text","params":{"text":"Hello World","pid":"${{steps.content.result.id}}","sz":84,"ff":"Inter","fs":"Bold","color":"#FFFFFF","w":780}}
  {"command":"gradient","params":{"nodeId":"${{steps.root.result.id}}","type":"LINEAR","stops":[{"position":0,"color":"#14344A"},{"position":1,"color":"#1A5F7A"}]}}

WRONG — "type" key instead of "command":
  {"type":"frame","name":"root","x":0,"y":0}
WRONG — params at top level instead of inside "params":
  {"name":"root","command":"frame","x":0,"y":0,"w":1080}
WRONG — invented step names not in manifest:
  {"command":"text","params":{"pid":"${{steps.background.result.id}}"}}

PARAM ALIASES (use these short forms):
  w=width, h=height, pid=parentId, sz=fontSize, ff=fontFamily, fs=fontStyle,
  bg=color, lh=lineHeight, sw=strokeWidth, ls=letterSpacing, r=cornerRadius

REFERENCE SYNTAX: ${{steps.STEPNAME.result.id}}
Only use step names from the provided manifest. Never invent new names.

NEVER exceed 15 operations.
```

---

## User Message Template (Per Design)

```
STEP NAME MANIFEST (use ONLY these names, exactly as written): {manifest}

Use ONLY names from the manifest above in any pid/parentId/nodeId references.

SPEC:
{design_spec}

Output ONLY the JSON object {"ops": [...]}. Max 15 ops. No markdown, no text outside JSON.
```

Replace `{manifest}` with comma-separated step names, e.g.: `root, overlay, headline, body, cta, cta_label`

---

## API Call Pattern (Python)

```python
import json, subprocess, tempfile, os

def call_model(api_key: str, endpoint: str, model: str,
               system_prompt: str, user_msg: str) -> list:
    """Returns normalized ops list. Raises on API error or parse failure."""
    payload = {
        "model": model,
        "messages": [
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_msg}
        ],
        "temperature": 0   # deterministic — schema compliance is not creative
        # DO NOT use response_format json_schema with strict:true —
        # it fails because params is an open object (additionalProperties conflict)
    }

    # Must use temp file — large JSON payloads break shell arg escaping
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        json.dump(payload, f)
        tmpfile = f.name
    try:
        result = subprocess.run(
            ["curl", "-sS", endpoint,
             "-H", "Content-Type: application/json",
             "-H", f"Authorization: Bearer {api_key}",
             "-d", f"@{tmpfile}"],   # ← @filepath, not inline -d '...'
            capture_output=True, text=True, timeout=60
        )
    finally:
        os.unlink(tmpfile)

    if not result.stdout.strip():
        raise RuntimeError("Empty API response")

    resp = json.loads(result.stdout)
    if "error" in resp and "choices" not in resp:
        raise RuntimeError(f"API error: {resp.get('message', resp)}")

    content = resp["choices"][0]["message"]["content"].strip()

    # Strip markdown fences — models output these even with explicit "no markdown" prompts
    if content.startswith("`"):
        lines = content.split("\n")
        content = "\n".join(l for l in lines if not l.strip().startswith("`"))

    # Parse — model may return {"ops": [...]} or bare [...]
    data = json.loads(content)
    if isinstance(data, dict):
        ops = data.get("ops", list(data.values())[0] if data else [])
    else:
        ops = data

    return normalize_ops(ops)
```

---

## Python Normalizer

Run this on model output BEFORE sending to the CLI. Fixes common drift, rejects garbage.

```python
TOP_LEVEL_PARAMS = {
    "x","y","width","height","w","h","color","bg","fillColor","cornerRadius","r",
    "parentId","pid","fontSize","sz","fontFamily","ff","fontStyle","fs","lineHeight","lh",
    "text","opacity","stroke","strokeWidth","sw","letterSpacing","ls","textCase","textAlign",
    "layoutMode","itemSpacing","padding","primaryAxisAlign","counterAxisAlign","imageData",
    "noFill","clipsContent","primaryAxisSizing","counterAxisSizing",
}

def normalize_ops(ops: list) -> list:
    """Normalize model output: fix type→command, hoist top-level params into params dict."""
    result = []
    for op in ops:
        if not isinstance(op, dict):
            continue
        normalized = {}

        # Fix "type" → "command"
        if "type" in op and "command" not in op:
            normalized["command"] = op.pop("type")

        # Copy known op-level fields
        for k in ("name", "command", "params"):
            if k in op:
                normalized[k] = op[k]

        # Hoist misplaced design props into params
        params = dict(normalized.get("params") or {})
        for k, v in op.items():
            if k in TOP_LEVEL_PARAMS:
                params[k] = v
        if params:
            normalized["params"] = params

        if "command" in normalized and "params" in normalized:
            result.append(normalized)
        # silently drop ops that are unrecoverable (no command)

    return result
```

---

## Step Name Manifest Pattern

Define your manifest BEFORE generating. Pass it explicitly in the user message:

```python
# Define the manifest for this design
manifest = ["root", "overlay", "headline", "body", "cta", "cta_label"]

# In the user message:
user_msg = f"""STEP NAME MANIFEST (use ONLY these names, exactly as written): {', '.join(manifest)}
...
"""
# NOTE: In f-strings, the interpolation syntax needs quadruple braces:
# f"${{{{steps.root.result.id}}}}" → renders as ${{steps.root.result.id}}
```

The model cannot invent `background` if the manifest says `root`.

---

## Two-Pass Generation (For Complex Designs)

When the design spec is complex, separate design thinking from encoding:

**Pass 1 — Structure only (no JSON):**
```
Describe the layout for [DESIGN SPEC]:
- What frames are needed?
- What text content?
- Colors and hierarchy?
Brief bullet points only.
```

**Pass 2 — Encode to JSON:**
```
Convert this plan to batch JSON. Use the schema above.
Plan: [PASTE PASS 1 OUTPUT]
Output ONLY {"ops": [...]}
```

---

## Chunking for Large Designs

Max 15 ops per request — confirmed in production. Schema compliance degrades beyond this.

For 30-op designs, split by layer:
- **Chunk 1**: Frame structure only (no text, no effects) — maybe 8-10 ops
- **Chunk 2**: Text nodes — uses real IDs from chunk 1 results
- **Chunk 3**: Effects, shadows, gradients

After chunk 1 runs, use real committed IDs from the result for chunk 2. More reliable than pre-defining all step names.

---

## Validation Self-Check (Optional Second Call)

Fire a cheap second call on the generated JSON:

```
Review this JSON. For each element verify:
1. Has "name", "command", "params" keys (not "type")
2. All x/y/width/height/color/etc are INSIDE "params", not top-level
3. All ${{steps.X.result.id}} use names from: [YOUR MANIFEST]

Return ONLY the corrected JSON object {"ops": [...]}. Unchanged if already correct.
```

Adds ~1s. Catches residual drift.

---

## Full Workflow Example

```python
from concurrent.futures import ThreadPoolExecutor, as_completed

designs = [
    {"id": "post_1", "spec": "Hero image post, brand blue, white text",
     "manifest": ["root", "img", "overlay", "headline", "sub"]},
    {"id": "post_2", "spec": "Quote card, dark background, yellow accent",
     "manifest": ["root", "card", "quote", "author", "badge"]},
]

def generate_one(design):
    manifest_str = ", ".join(design["manifest"])
    user_msg = f"""STEP NAME MANIFEST (use ONLY these): {manifest_str}
SPEC: {design['spec']}
Output ONLY {{"ops": [...]}}. Max 15 ops."""

    ops = call_model(API_KEY, ENDPOINT, MODEL, SYSTEM_PROMPT, user_msg)

    # Write to file
    path = f"/tmp/batch_{design['id']}.json"
    with open(path, "w") as f:
        json.dump(ops, f, indent=2)
    return design["id"], path

with ThreadPoolExecutor(max_workers=8) as ex:
    futures = {ex.submit(generate_one, d): d["id"] for d in designs}
    for fut in as_completed(futures):
        slide_id, path = fut.result()
        print(f"✓ {slide_id} → {path}")

# Validate all, then batch-send
import subprocess
for design in designs:
    path = f"/tmp/batch_{design['id']}.json"
    r = subprocess.run(["ai-happy-design", "validate", path], capture_output=True)
    if r.returncode == 0:
        subprocess.run(["ai-happy-design", "batch", path])
    else:
        print(f"✗ {design['id']} failed validation:\n{r.stderr.decode()}")
```
```

**Step 2: Add link to SKILL.md reference section**

Find `## Reference Files` and add:
```markdown
- **[Weak Model Prompting](references/weak-model-prompting.md)** — System prompt templates, Python normalizer, two-pass generation, chunking strategy for smaller/faster models (production-validated)
```

**Step 3: Commit**

```bash
git add skills/ai-happy-design/references/weak-model-prompting.md
git add skills/ai-happy-design/SKILL.md
git commit -m "docs: add production-validated weak model prompting guide

System prompt template, Python normalizer, API call pattern with temp
file + fence stripping, step name manifest, chunking strategy. All
patterns validated in a real 36-slide parallel generation session.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 3: CLI `validate` Command

**Files:**
- Create: `cmd/ai-happy-design/validate.go`
- Modify: `cmd/ai-happy-design/main.go` (register command)

This is the deterministic gate before Figma. The Python normalizer in Task 2 does best-effort fixing; the Go validator is the final pass that rejects anything unfixable.

**Step 1: Write `cmd/ai-happy-design/validate.go`**

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file.json or '-' for stdin]",
	Short: "Validate batch JSON against the ai-happy-design schema",
	Long: `Checks batch JSON for common schema errors before sending to Figma.

Detects:
  - Using 'type' instead of 'command'
  - Missing 'params' field
  - Design properties at top level instead of inside params
  - Broken ${{steps.X.result.id}} references (X not defined as a step name)

Exit code 0 = valid, 1 = validation errors found.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var data []byte
		var err error

		if len(args) == 0 || args[0] == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(args[0])
		}
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		errs := validateBatchOps(data)
		if len(errs) == 0 {
			fmt.Println("✓ Valid — 0 errors found")
			return nil
		}
		fmt.Fprintf(os.Stderr, "✗ %d error(s) found:\n\n", len(errs))
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "  "+e)
		}
		os.Exit(1)
		return nil
	},
}

// knownTopLevelProps are keys that belong inside "params", not at the op root.
var knownTopLevelProps = []string{
	"x", "y", "width", "height", "w", "h",
	"color", "fillColor", "bg",
	"fontSize", "fontFamily", "fontStyle", "sz", "ff", "fs",
	"parentId", "pid",
	"cornerRadius", "r",
	"layoutMode", "itemSpacing", "padding", "opacity",
	"text", "imageData", "stroke", "strokeWidth",
}

var interpolationRef = regexp.MustCompile(`\$\{\{steps\.([a-zA-Z0-9_]+)\.result\.[a-z]+\}\}`)

func validateBatchOps(data []byte) []string {
	var ops []map[string]interface{}
	if err := json.Unmarshal(data, &ops); err != nil {
		return []string{fmt.Sprintf("Invalid JSON: %v", err)}
	}
	if len(ops) == 0 {
		return []string{"Batch is empty"}
	}

	var errs []string
	definedNames := map[string]bool{}

	for i, op := range ops {
		label := fmt.Sprintf("op[%d]", i)
		if name, ok := op["name"].(string); ok && name != "" {
			definedNames[name] = true
			label = fmt.Sprintf("op[%d] %q", i, name)
		}

		// "type" instead of "command"
		if _, hasType := op["type"]; hasType {
			errs = append(errs, label+`: use "command" not "type"`)
		}

		// Missing "command"
		if cmd, ok := op["command"].(string); !ok || cmd == "" {
			errs = append(errs, label+`: missing "command" field`)
		}

		// Missing or wrong-type "params"
		if p, hasParams := op["params"]; !hasParams {
			errs = append(errs, label+`: missing "params" field`)
		} else if _, ok := p.(map[string]interface{}); !ok {
			errs = append(errs, label+`: "params" must be an object`)
		}

		// Design props at top level
		for _, prop := range knownTopLevelProps {
			if _, ok := op[prop]; ok {
				errs = append(errs, fmt.Sprintf(`%s: %q must be inside "params", not at top level`, label, prop))
			}
		}
	}

	// Second pass: validate all interpolation references
	raw, _ := json.Marshal(ops)
	for _, m := range interpolationRef.FindAllSubmatch(raw, -1) {
		refName := string(m[1])
		if !definedNames[refName] {
			errs = append(errs, fmt.Sprintf(`interpolation ${{steps.%s.result.id}} references undefined step name`, refName))
		}
	}

	return errs
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
```

**Step 2: Write tests at `cmd/ai-happy-design/validate_test.go`**

```go
package main

import (
	"strings"
	"testing"
)

func TestValidateBatchOps_Valid(t *testing.T) {
	input := `[
		{"name":"bg","command":"node.create_frame","params":{"x":0,"y":0,"width":1080,"height":1080}},
		{"name":"title","command":"text.create","params":{"text":"Hello","parentId":"${{steps.bg.result.id}}"}}
	]`
	errs := validateBatchOps([]byte(input))
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateBatchOps_TypeInsteadOfCommand(t *testing.T) {
	input := `[{"name":"bg","type":"node.create_frame","params":{"x":0}}]`
	errs := validateBatchOps([]byte(input))
	if len(errs) == 0 {
		t.Error("expected error for 'type' field, got none")
	}
}

func TestValidateBatchOps_TopLevelParams(t *testing.T) {
	input := `[{"name":"bg","command":"node.create_frame","params":{},"x":0,"color":"#fff"}]`
	errs := validateBatchOps([]byte(input))
	if len(errs) < 2 {
		t.Errorf("expected 2+ errors (x, color at top level), got %d: %v", len(errs), errs)
	}
}

func TestValidateBatchOps_BrokenInterpolation(t *testing.T) {
	// Step is named "bg" but reference uses "background"
	input := `[
		{"name":"bg","command":"node.create_frame","params":{"x":0}},
		{"name":"title","command":"text.create","params":{"parentId":"${{steps.background.result.id}}"}}
	]`
	errs := validateBatchOps([]byte(input))
	found := false
	for _, e := range errs {
		if strings.Contains(e, "background") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for undefined step 'background', got: %v", errs)
	}
}

func TestValidateBatchOps_MissingParams(t *testing.T) {
	input := `[{"name":"bg","command":"node.create_frame"}]`
	errs := validateBatchOps([]byte(input))
	if len(errs) == 0 {
		t.Error("expected error for missing params, got none")
	}
}

func TestValidateBatchOps_InvalidJSON(t *testing.T) {
	errs := validateBatchOps([]byte("not json"))
	if len(errs) == 0 || !strings.Contains(errs[0], "Invalid JSON") {
		t.Errorf("expected Invalid JSON error, got: %v", errs)
	}
}

func TestValidateBatchOps_ShortAliasCommands(t *testing.T) {
	// Short aliases (frame, text, rect) should not cause errors
	input := `[
		{"name":"bg","command":"frame","params":{"x":0,"y":0,"w":1080,"h":1080}},
		{"name":"title","command":"text","params":{"text":"Hi","pid":"${{steps.bg.result.id}}"}}
	]`
	errs := validateBatchOps([]byte(input))
	if len(errs) != 0 {
		t.Errorf("short alias commands should be valid, got: %v", errs)
	}
}
```

**Step 3: Run tests**

```bash
cd "/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design"
go test ./cmd/ai-happy-design/... -run TestValidate -v
```
Expected: all 6 tests pass.

**Step 4: Manual smoke test**

```bash
go build -o /tmp/ahd-test ./cmd/ai-happy-design/

# Should pass
echo '[{"name":"bg","command":"node.create_frame","params":{"x":0,"y":0}}]' | /tmp/ahd-test validate -
# Expected: ✓ Valid — 0 errors found

# Should fail with 3 errors
echo '[{"name":"bg","type":"frame","x":0,"color":"#fff"}]' | /tmp/ahd-test validate -
# Expected: ✗ 3 error(s) found (type, missing params, top-level x, top-level color)
```

**Step 5: Update SKILL.md Key CLI Commands**

Add to the `## Key CLI Commands` section:
```bash
ai-happy-design validate batch.json                  # validate schema before sending to Figma
cat model-output.json | ai-happy-design validate -   # validate from stdin
```

**Step 6: Deploy**

```bash
make deploy
```

**Step 7: Commit**

```bash
git add cmd/ai-happy-design/validate.go cmd/ai-happy-design/validate_test.go
git add skills/ai-happy-design/SKILL.md
git commit -m "feat: add validate command for batch JSON schema checking

Catches type-vs-command, top-level params, broken step name references
before Figma is touched. Enables generate→validate→fix→send workflow.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Summary

| Task | Files | Impact |
|------|-------|--------|
| 1. SKILL.md batch rewrite | `SKILL.md` | LLMs see WRONG/RIGHT upfront, don't guess schema |
| 2. Weak model reference | `references/weak-model-prompting.md` | Agents get battle-tested system prompts + Python normalizer |
| 3. CLI validate command | `cmd/validate.go` + tests | Deterministic guard — stops bad JSON before Figma |

**Execution order**: Tasks 1+2 are docs (no build needed), Task 3 requires `make deploy`.

---

## What This Enables

```
Smaller/faster model (parallel, cheap)
  → system prompt from references/weak-model-prompting.md
  → generates {"ops": [...]} × N designs in parallel
  → strip markdown fences (always needed)
  → Python normalize_ops() (fix common drift)
  → ai-happy-design validate output.json  (gate: exit 1 if broken)
  → if errors: re-prompt with error text  (tight feedback loop)
  → ai-happy-design batch output.json     (hits Figma only when clean)
```

This routes around model weaknesses with deterministic tooling instead of fighting them.
The key insight from production: **never send model output directly to Figma**. Always normalize first.
