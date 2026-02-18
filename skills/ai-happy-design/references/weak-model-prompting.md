# Weak Model Prompting Guide

Use when delegating batch JSON generation to a smaller/faster model.
All patterns below were validated in a real parallel generation session (36 slides, 12 workers).

---

## The Correct Workflow

```
Smaller/faster model (creative, parallel)
  → generates {"ops": [...]} JSON per design
  → strip markdown fences (always needed, even with explicit "no markdown" prompt)
  → Python normalize/validate (fix common drift, reject garbage)
  → ai-happy-design validate output.json (final deterministic check)
  → ai-happy-design batch output.json (only reaches Figma if clean)
```

**Never send model output directly to Figma.** Always normalize first.

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

## Python Normalizer

Run on model output BEFORE sending to the CLI. Fixes common drift, rejects unrecoverable ops.

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

        # Drop ops with no command — unrecoverable
        if "command" in normalized and "params" in normalized:
            result.append(normalized)

    return result
```

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
import subprocess

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

with ThreadPoolExecutor(max_workers=8) as ex:
    futures = {ex.submit(generate_one, d): d["id"] for d in designs}
    for fut in as_completed(futures):
        slide_id, path = fut.result()
        print(f"✓ {slide_id} → {path}")

# Validate then send
for design in designs:
    path = f"/tmp/batch_{design['id']}.json"
    r = subprocess.run(["ai-happy-design", "validate", path], capture_output=True)
    if r.returncode == 0:
        subprocess.run(["ai-happy-design", "batch", path])
    else:
        print(f"✗ {design['id']} validation failed:\n{r.stderr.decode()}")
```
