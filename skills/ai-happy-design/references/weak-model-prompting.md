# Weak Model Prompting Guide

Use when delegating batch JSON generation to a smaller/faster model.
All patterns below were validated in a real parallel generation session (36 slides, 12 workers).

---

## The Correct Workflow

```
Smaller/faster model (creative, parallel)
  → generates {"ops": [...]} JSON per design
  → ai-happy-design validate --fix output.json   (fixes fences, type→command, top-level props)
  → ai-happy-design batch output.json            (sends to Figma only if valid)
```

The CLI handles all normalization. **The model's only job is generating JSON.**

`validate --fix` auto-corrects:
- Markdown fences (` ```json ... ``` ` wrappers models add despite being told not to)
- `"type"` → `"command"` rename
- Top-level design props (`x`, `y`, `color`, etc.) hoisted into `"params"`

What `validate --fix` cannot fix (model must get right):
- Step name mismatches in `${{steps.X.result.id}}` references
- Missing `"command"` entirely
- Structurally invalid JSON

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

WRONG - "type" key instead of "command":
  {"type":"frame","name":"root","x":0,"y":0}
WRONG - params at top level instead of inside "params":
  {"name":"root","command":"frame","x":0,"y":0,"w":1080}
WRONG - invented step names not in manifest:
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
        "temperature": 0
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

## Normalization

Use the CLI — no Python needed:

```bash
# Fix common issues in-place, then validate
ai-happy-design validate --fix output.json

# Fix from stdin (e.g. piped from model call), output corrected JSON
echo '{"ops":[...]}' | python3 -c "import sys,json; d=json.load(sys.stdin); print(json.dumps(d['ops'] if isinstance(d,dict) else d))" | ai-happy-design validate --fix -

# Full workflow
ai-happy-design validate --fix output.json && ai-happy-design batch output.json
```

The `--fix` flag handles what the Python `normalize_ops()` function used to do. No custom normalization code needed.

---

## Step Name Manifest Pattern

Define your manifest BEFORE generating. Pass it in the user message:

```python
manifest = ["root", "overlay", "headline", "body", "cta", "cta_label"]

user_msg = f"""STEP NAME MANIFEST (use ONLY these names, exactly as written): {', '.join(manifest)}
...
"""
# NOTE: In Python f-strings, the interpolation syntax needs quadruple braces:
# f"${{{{steps.root.result.id}}}}" renders as ${{steps.root.result.id}}
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
- **Chunk 1**: Frame structure only (no text, no effects) — ~8-10 ops
- **Chunk 2**: Text nodes — reference real IDs from chunk 1 results
- **Chunk 3**: Effects, shadows, gradients

After chunk 1 runs, use real committed IDs for chunk 2. More reliable than pre-defining all step names.

---

## Validation Self-Check (Optional Second Call)

After generation, fire a second cheap call:

```
Review this JSON. For each element verify:
1. Has "name", "command", "params" keys (not "type")
2. All x/y/width/height/color/etc are INSIDE "params", not top-level
3. All ${{steps.X.result.id}} use names from: [YOUR MANIFEST]

Return ONLY the corrected JSON object {"ops": [...]}. Unchanged if already correct.
```

Adds ~1s. Catches residual drift before the CLI gate.

---

## Full Workflow Example

```python
from concurrent.futures import ThreadPoolExecutor, as_completed
import subprocess, json, os

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
    path = f"/tmp/batch_{design['id']}.json"
    with open(path, "w") as f:
        json.dump(ops, f, indent=2)
    return design["id"], path

# Generate all in parallel
with ThreadPoolExecutor(max_workers=8) as ex:
    futures = {ex.submit(generate_one, d): d["id"] for d in designs}
    for fut in as_completed(futures):
        slide_id, path = fut.result()
        print(f"✓ generated {slide_id} → {path}")

# Fix, validate, send — CLI handles all normalization
for design in designs:
    path = f"/tmp/batch_{design['id']}.json"
    # --fix auto-corrects fences, type→command, top-level props
    r = subprocess.run(["ai-happy-design", "validate", "--fix", path],
                       capture_output=True, text=True)
    if r.returncode == 0:
        subprocess.run(["ai-happy-design", "batch", path])
        print(f"✓ sent {design['id']}")
    else:
        print(f"✗ {design['id']} still invalid after fix:\n{r.stderr}")
```

No Python normalization code. The model generates JSON, the CLI fixes and validates it.
