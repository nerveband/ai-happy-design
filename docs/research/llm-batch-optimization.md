# LLM Batch JSON Optimization Research

> Research conducted 2026-02-15 for ai-happy-design.
> Goal: Make batch JSON generation bulletproof and fast for LLMs operating via CLI/MCP.

---

## Table of Contents

1. [Problem Statement](#problem-statement)
2. [General LLM-to-JSON Techniques](#general-llm-to-json-techniques)
3. [Current ai-happy-design Batch Format](#current-batch-format)
4. [Figma API Color System](#figma-api-color-system)
5. [What Other Projects Do](#what-other-projects-do)
6. [Format Comparison: JSON vs YAML vs ASCII Trees vs CSS](#format-comparison)
7. [Proposed Optimizations](#proposed-optimizations)
8. [Revised Assessment & Recommendations](#revised-assessment)
9. [Sources](#sources)

---

## Problem Statement

When an LLM (e.g., Claude Code) creates Figma designs via ai-happy-design, it must generate batch JSON payloads token-by-token. A typical Instagram post batch is **21KB / 79 operations**. Most of that JSON is structural boilerplate, not design data.

The core issue: every `{`, `"`, `:` is a token the LLM must generate. More tokens = slower, more expensive, and more attention spent on syntax instead of design reasoning.

Research shows that **format restrictions cause 10-15% reasoning degradation** in LLMs — the model spends attention tracking bracket matching and comma placement instead of thinking about the actual design.

**Goal:** Reduce token count, reduce operation count, reduce error surface, so the LLM can focus on design quality.

---

## General LLM-to-JSON Techniques

Techniques people use for fast/reliable LLM-to-JSON, ranked by speed impact:

### 1. Parallel Field Generation (~10x speed)

Split the schema into independent key-value pairs and generate them as separate parallel requests in a single batch call. Since most JSON fields don't depend on each other, the LLM can fill them all simultaneously.

- **Applicability to us:** NOT applicable. Our batch ops are sequential (interpolation needs execution order).

### 2. Compressed FSM / Jump-Forward Decoding (~2-2.5x speed)

Convert JSON schema into a regex/FSM, identify "singular paths" where only one valid token can follow, and skip ahead multiple tokens at once. Used by SGLang/vLLM.

- **Applicability:** Requires local model hosting. Not relevant for API/Claude Code usage.

### 3. Schema Pre-filling / Jsonformer (~2-3x speed)

Walk the JSON schema, insert `{`, `"key":`, etc. yourself, and only call the LLM when you need an actual value. Dramatically reduces token count since the model skips predictable syntax.

- **Applicability:** Only works with local models (HuggingFace). Not available for API calls.

### 4. Provider Structured Output APIs (~1.5x + no retries)

Use `response_format: { type: "json_schema", schema: ... }` or tool use with schemas. The provider constrains decoding at the token level server-side.

- OpenAI: `response_format` with JSON schema — claims 100% schema compliance
- Anthropic: Tool use with input schemas (effectively the same — model fills a structured tool call)
- Google Gemini: `response_mime_type: "application/json"` + `response_schema`
- Local models via vLLM/SGLang: `guided_json` parameter

- **Applicability:** Not available in Claude Code (free-form generation, no schema constraints).

### 5. Batch API / Async Processing (throughput, not latency)

Submit a file of requests, get results back asynchronously. Typically 50% cheaper than real-time.

- **Applicability:** Different problem — we need real-time generation.

### 6. Generate Data, Not JSON

Don't ask the LLM to write JSON. Ask it to write a script/command that produces the JSON. The LLM outputs ~1 line of code, the shell produces hundreds of lines of perfect JSON instantly.

Examples:
- `jq` commands that generate JSON
- CSV/YAML that gets converted to JSON
- Python/Go scripts that output JSON
- Template files with value edits

- **Applicability:** HIGHLY relevant. This is essentially what our optimizations do — reduce the LLM's job from "generate 21KB of JSON" to "generate 6.5KB of compact operations."

### Most Reliable Approach (API)

Provider structured output APIs (tool use / function calling with strict schemas). Guaranteed valid JSON that matches your schema, zero retries.

### Most Reliable Approach (Claude Code / Non-API)

Reduce the generation surface: fewer operations, shorter syntax, sensible defaults. Let the LLM think about design, not JSON syntax.

---

## Current Batch Format

### Operations Array Schema

```json
[
  {
    "name": "stepName",
    "command": "domain.action",
    "params": { }
  }
]
```

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | No | Step identifier for interpolation: `${{steps.NAME.result.X}}` |
| `command` | string | Yes | `domain.action` format (e.g., `node.create_frame`, `paint.set_solid`) |
| `params` | object | Yes | Command-specific parameters |

### Interpolation System

```
${{steps.INDEX.result.FIELD}}   — by index (0-based)
${{steps.NAME.result.FIELD}}    — by named step
${{steps[INDEX].result.FIELD}}  — bracket notation
${{last.result.FIELD}}          — most recent step
```

Type preservation: when entire value is a single placeholder, type is preserved. When embedded in a larger string, converted to string.

### CLI Interface

```bash
# Positional arg (inline JSON)
ai-happy-design batch '[{"command":"node.create_frame","params":{"x":0,"y":0}}]'

# Positional arg (file path)
ai-happy-design batch operations.json

# Flags
ai-happy-design batch -o '[...]'
ai-happy-design batch -f operations.json

# Stdin
cat operations.json | ai-happy-design batch
```

| Flag | Default | Description |
|------|---------|-------------|
| `--fail-fast` | false | Stop at first failure |
| `--retries` | 1 | Retry attempts per op |
| `--retry-delay-ms` | 250 | Delay between retries |
| `--interpolate` | true | Enable placeholder interpolation |
| `--compress-images` | false | ImageMagick compression |
| `--live` | false | Stream progress updates |

### Batch Output Format

```json
{
  "ok": true,
  "summary": { "total": 3, "processed": 3, "succeeded": 3, "failed": 0 },
  "timing": { "totalMs": 1250, "avgMs": 416, "opsPerSec": 2 },
  "steps": [
    { "index": 0, "name": "createCard", "command": "...", "ok": true, "result": {...}, "elapsedMs": 450 }
  ]
}
```

### MCP Bulk Tool

```json
{
  "action": "execute",
  "operations": "[{...}]",
  "failFast": false,
  "retries": 1,
  "interpolate": true
}
```

### Common Batch Patterns

**Pattern: Create + Style (current — verbose)**
```json
{"name":"ring","command":"shape.create_ellipse","params":{"x":10,"y":10,"width":160,"height":160}},
{"command":"paint.set_stroke","params":{"nodeId":"${{steps.ring.result.id}}","color":{"r":1,"g":0.84,"b":0,"a":0.15},"width":3}},
{"command":"paint.remove_fill","params":{"nodeId":"${{steps.ring.result.id}}","index":0}},
{"command":"layer.move_to_parent","params":{"nodeId":"${{steps.ring.result.id}}","parentId":"${{steps.frame.result.id}}"}}
```
4 operations, ~490 characters for ONE ring element.

**What already works (no code changes):**
- `parentId` on all create commands (frame, rect, ellipse, text, line)
- `color` (fill) on frame, rect, ellipse, text
- `cornerRadius` on rect, image

Using these inline params, the 79-op batch could drop to ~45 ops immediately — the example batches were written before these params existed.

---

## Figma API Color System

### Figma's Native Color Types

From the [Plugin API](https://developers.figma.com/docs/plugins/api/Paint/):

| Paint Type | Key Fields |
|------------|-----------|
| `SolidPaint` | `color: RGB {r, g, b}` + `opacity: 0-1` (alpha is on paint, not color) |
| `GradientPaint` | `gradientStops: [{position, color: RGBA}]`, types: LINEAR, RADIAL, ANGULAR, DIAMOND |
| `ImagePaint` | `imageHash`, `scaleMode: FILL/FIT/CROP/TILE` |
| `VideoPaint` | `videoHash`, same scaleModes |
| `PatternPaint` (beta) | `sourceNodeId`, `tileType`, `spacing` |

### RGB/RGBA

All values are **floats 0-1** (not 0-255).

```typescript
interface RGB  { r: number; g: number; b: number }        // Used by SolidPaint
interface RGBA { r: number; g: number; b: number; a: number }  // Used by GradientStops
```

SolidPaint splits this: `color` uses RGB (no alpha), alpha goes to paint's `opacity`.

### Figma's Color Models (UI)

Figma supports five color models in the UI: **Hex** (default), **RGB**, **HSB**, **HSL**, **CSS** (rgba).

### Figma's Built-in Utilities

```typescript
figma.util.rgb('#FF00FF')      // → RGB object. Accepts: hex, rgb(), hsl(), lab()
figma.util.rgba('#FF00FF88')   // → RGBA object. Same formats.
figma.util.solidPaint('#FF00FF88') // → SolidPaint with color + opacity
```

**Key insight:** Figma's own utilities accept CSS color strings (hex, rgb(), hsl(), lab()) and throw on bad input. No named colors.

### Our Current `parseHexColor()`

Duplicated across **6 plugin handlers** (style.ts, text.ts, paint.ts, effect.ts, shape.ts, node.ts). Accepts hex strings (3/6/8 digit) and `{r,g,b}` objects.

**Critical bug: silent fallback to black on bad input.**

```typescript
if (typeof color !== 'string') return fallback;   // ← silently becomes black
if (!(raw.length === 3 || raw.length === 6 || raw.length === 8)) return fallback; // ← silently becomes black
if (Number.isNaN(n)) return fallback;             // ← silently becomes black
```

If an LLM sends `"gold"` or `"rgba(255,0,0,0.5)"` or `"#FFG700"`, it silently becomes black. The LLM never knows it failed.

### Color Precision Spectrum

| Format | Precision | Colors | Deterministic? | LLM Reliability |
|--------|-----------|--------|:--------------:|:---------------:|
| `#RRGGBB` (6-digit) | 24-bit | 16,777,216 | Yes | Very high |
| `#RRGGBBAA` (8-digit) | 32-bit | 4.29 billion | Yes | High |
| `#RGB` (3-digit) | 12-bit | 4,096 | Yes | Medium |
| `{r:0.84, g:0.2, b:0}` (Figma float) | ~24-bit | ~16.7M | Yes | Low (0-1 vs 0-255 confusion) |
| `rgb(255, 214, 0)` | 24-bit | 16,777,216 | Yes | High |
| `hsl(50, 100%, 50%)` | 24-bit equiv | 16,777,216 | Yes | Medium |
| `"gold"` (named) | Exact (CSS spec) | 1 | Yes but... | **Dangerous** |

### Color Input Policy

**Accept (precise, deterministic):**
- `"#FFD700"` — 6-digit hex (primary, always show in examples)
- `"#FFD70080"` — 8-digit hex with alpha
- `"#FD0"` — 3-digit hex (accept, expand internally)
- `{r:1, g:0.84, b:0}` — Figma native float (0-1)
- `{r:255, g:214, b:0}` — Auto-detect 0-255 range, normalize
- `"rgb(255, 214, 0)"` — CSS function
- `"rgba(255, 214, 0, 0.5)"` — CSS function with alpha
- `"hsl(50, 100%, 50%)"` — Perceptual color space (resolve to hex)

**Reject with clear error:**
- `"gold"` → "Named colors not accepted. Use hex: #FFD700"
- `"transparent"` → "Use alpha: #00000000 or opacity param"
- `"#GG0000"` → "Invalid hex digit 'G'. Hex uses 0-9 and A-F."
- malformed → "Could not parse color. Accepted: #RRGGBB, #RRGGBBAA, rgb(), hsl(), {r,g,b}"

**Rationale:** Hex IS precise. Named colors are not — they let the LLM be lazy ("golden-ish" → `gold`) instead of committing to an exact value (`#FFD700`). Figma's own `util.rgb()` takes the same approach: accepts hex/rgb/hsl, throws on bad input, no named colors.

### How Others Handle Color

| System | Named Colors | Hex | CSS rgb()/hsl() | Bad Input |
|--------|:-----------:|:---:|:---------------:|-----------|
| **Figma `util.rgb()`** | No | Yes | Yes (hsl, lab too) | **Throws error** |
| **Figma `util.solidPaint()`** | No | Yes | Yes | **Throws error** |
| **W3C Design Tokens** | No | Yes (canonical) | Yes (as strings) | Validation error |
| **Tokens Studio** | No (resolves→hex) | Yes | Yes | Error |
| **Our `parseHexColor()`** | No | Yes | **No** | **Silent fallback to black** |

### Future: Gradient Input

Not for v1. CSS gradient syntax (`linear-gradient(135deg, #FF0000, #0000FF)`) → Figma gradient transform matrices requires ~400 lines of coordinate math. Niche use case for social media designs. File for v2.

---

## What Other Projects Do

### The Fundamental Difference

Every other Figma optimization project solves the **opposite** problem:

| Direction | Problem | Projects |
|-----------|---------|----------|
| **Figma → LLM** (read) | Design data too verbose for context windows | Figma-Context-MCP, Figma Official MCP, UIKit Data Exporter, Figma Raw |
| **LLM → Figma** (write) | Batch JSON too verbose for token generation | **ai-happy-design (us)** |

Nobody else optimizes LLM output for Figma creation. We're in uncharted territory.

### Figma-Context-MCP (Framelink) — 13.1K stars

**What they do:** Simplify Figma API responses for LLM consumption.

**Key techniques:**
- ASCII tree format with inline properties → 95% token reduction (205K → 10K chars)
- YAML output instead of JSON
- CSS-like property names (`mode: row`, `gap: 16px`, `padding: 8px 16px`)
- Strip default/unchanged properties
- Convert colors from Figma `{r,g,b}` floats to hex/rgba
- Convert gradients to CSS `linear-gradient()` syntax (~400 lines of coordinate math)
- Convert layout to CSS flexbox equivalents

**Their color handling (from source code):**
```typescript
// style.ts — converts Figma RGBA to hex
export function convertColor(color: RGBA, opacity = 1): ColorValue {
  const r = Math.round(color.r * 255);
  const g = Math.round(color.g * 255);
  const b = Math.round(color.b * 255);
  const a = Math.round(opacity * color.a * 100) / 100;
  const hex = "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1).toUpperCase();
  return { hex, opacity: a };
}
```

They chose hex as canonical output because it's what LLMs understand best. Same reasoning supports hex as our canonical INPUT.

**Relevance:** Validates hex colors and CSS-like aliases. Their flagship technique (ASCII trees) doesn't transfer — it's for read, and LLMs are bad at generating ASCII art.

### Figma Official MCP Server

**What they do:** Provide design context to LLMs for code generation.

**Key techniques:**
- `get_design_context` → React+Tailwind structured representation
- `get_metadata` → compact XML tree for large designs (drill-in approach)
- `get_variable_defs` → extract design tokens (colors, spacing, typography)
- Code Connect → map Figma components to code components in your repo
- `get_screenshot` → visual reference for layout fidelity

**Relevance:** `get_metadata` (compact overview → drill in) is similar to our `compute_tokens` (compute sizes → then create). Their Code Connect could inspire a template library. Neither directly helps with batch JSON optimization.

### UIKit Data Exporter

**What they do:** Export design system data in JSON or "TOON" format.

**TOON format:** ~50% size reduction, LLM-optimized. Property flattening and shortening — similar to our param shorthand idea.

**Relevance:** Validates the param shorthand approach.

### Technique Transfer Assessment

| Technique | Source | Direction | Transfers? | Why |
|-----------|--------|-----------|:----------:|-----|
| ASCII tree format | Figma-Context-MCP | Read | **No** | Trees can't express sequential deps; LLMs bad at generating ASCII art |
| YAML output | Figma-Context-MCP | Read | **Partially** | Worth offering as optional input (see Format Comparison section) |
| CSS property aliases | Figma-Context-MCP | Read | **Partial** | A few intuitive aliases work. Not a full CSS layer |
| Property stripping | All read projects | Read | **Yes** | Same principle in reverse — we set good defaults on input |
| Hex as canonical color | Figma-Context-MCP | Read | **Yes** | They chose hex for LLM readability; validates our approach |
| CSS gradients | Figma-Context-MCP | Read | **Future** | Complex math; niche use case; v2 |
| Component templates | Figma Official MCP | Read | **Future** | Different scope; design systems |
| Compact overview | Figma Official MCP | Read | **No** | Already handled by `compute_tokens` |

---

## Format Comparison

### Research: JSON vs YAML vs ASCII Trees vs CSS for LLM Generation

#### YAML

**AlphaCodium paper** (widely cited): "YAML output is far better for code generation." Reasoning: JSON forces tracking nested bracket matching, quote pairing, and comma placement simultaneously. YAML only requires consistent indentation.

**ImprovingAgents benchmark** (Oct 2025, 1000 questions/format, 3 models):

| Format | GPT-5 Nano | Gemini 2.5 Flash | Llama 3.2 3B | Token Efficiency |
|--------|:----------:|:----------------:|:------------:|:----------------:|
| **YAML** | **62.1%** | **51.9%** | 49.1% | -10% vs JSON |
| Markdown | 54.3% | 48.2% | 48.0% | **-34% vs JSON** |
| JSON | 50.3% | 43.1% | **52.7%** | baseline |
| XML | 44.4% | 33.8% | 50.7% | +80% vs JSON |

YAML won 2/3 models for comprehension accuracy. These numbers measure how well LLMs READ formats, but generation follows similar patterns.

**Production reports:** 30-50% fewer parsing failures when switching from JSON to YAML output.

**YAML weakness:** Whitespace sensitivity. One wrong indent = broken parse. In streaming or long sequences, indentation drift happens. Error messages are often cryptic.

**For Claude specifically:** Claude's Structured Outputs (Nov 2025) compiles JSON schemas into grammars at inference time — guaranteed valid JSON. **No equivalent for YAML.** In Claude Code (our context), neither format gets schema constraints.

**Bottom line:** YAML is genuinely better for generation efficiency (fewer tokens, fewer bracket-matching errors). Worth offering as optional input. Not as the primary format.

#### ASCII Trees

**Research:** LLMs are bad at this.

- "LLMs are inherently biased toward processing sequential information rather than spatial patterns"
- ASCIIBench: GPT-4 achieves only **25.19% accuracy** on single-character recognition in ASCII art
- Figma-Context-MCP uses ASCII trees for OUTPUT (machine-generated, perfect). They never ask LLMs to PRODUCE them.

**For our use case:** ASCII trees would be unreliable. The LLM would frequently mix up `├─`/`└─`, break indentation alignment, lose nesting depth, generate invalid UTF-8 box-drawing characters.

**Verdict: Don't use ASCII trees for input.**

#### CSS

LLMs are excellent at CSS. Massively represented in training data, every LLM has seen billions of CSS property-value pairs, syntax is simple (`property: value;`).

But CSS-to-Figma translation is a different problem:
- Simple mappings work: `gap: 16px` → `itemSpacing: 16`
- Complex mappings are hard: `box-shadow` → Figma effects, `linear-gradient()` → Figma gradient transforms

**Verdict:** Accept CSS aliases for common simple properties. Don't build a full CSS-to-Figma translation layer.

#### The Reasoning Penalty

From research: "Significant decline in LLMs' reasoning abilities under format restrictions... stricter format constraints generally leading to greater performance degradation in reasoning tasks."

**Math/reasoning drops 10-15% when the LLM must simultaneously reason AND produce valid JSON.** Generating correct syntax competes for attention with the actual task.

This directly supports our optimizations:
1. **Fewer tokens to generate = more attention for reasoning.** Combined-ops and short-interpolation help here.
2. **Simpler syntax = less reasoning penalty.** YAML's simpler syntax means less attention on format compliance.

### Claude's Recommendations

From Claude's structured outputs docs and prompt engineering best practices:

- **Use JSON schema constraints** for guaranteed valid output (API only, not Claude Code)
- **Forcing JSON during reasoning degrades accuracy by 10-15%** — let the model think first, then format
- Claude is "highly versatile across JSON, YAML, XML, and CSV"
- When tracking structured information, "use JSON or other structured formats"
- Show examples of the desired format in the prompt — the LLM will mirror what it sees
- No blanket "JSON is best" recommendation

---

## Proposed Optimizations

### Priority 1: Short Interpolation Syntax

**Current:** `"${{steps.frame.result.id}}"` (30+ chars, easy to typo)

**Proposed:** `"$frame"` (6 chars)

```
$frame         → ${{steps.frame.result.id}}
$frame.name    → ${{steps.frame.result.name}}
$frame.width   → ${{steps.frame.result.width}}
$last          → ${{last.result.id}}
```

**Implementation:** Regex expansion in `interpolation.go` before `BuildContext`. Backward compatible — long form still works.

**Savings:** ~25 chars per reference × 40+ references = ~1000 chars per batch.

### Priority 2: Combined Create Operations

Fold `paint.set_stroke`, `paint.remove_fill`, `node.set_corner_radius`, `node.set_opacity` into creation commands.

**Already works inline (no changes needed):**
- `parentId` on all creates
- `color` (fill) on frame/rect/ellipse
- `cornerRadius` on rect

**Needs adding to plugin handlers:**
- `stroke` + `strokeWidth` on rect/ellipse (only on line currently)
- `noFill: true` to remove default white fill
- `opacity` on create
- `effects` (shadow) on create
- `letterSpacing`, `textCase` on `text.create`

**Before** (4 ops, ~490 chars):
```json
{"name":"ring","command":"shape.create_ellipse","params":{"x":10,"y":10,"width":160,"height":160}},
{"command":"paint.set_stroke","params":{"nodeId":"${{steps.ring.result.id}}","color":{"r":1,"g":0.84,"b":0,"a":0.15},"width":3}},
{"command":"paint.remove_fill","params":{"nodeId":"${{steps.ring.result.id}}","index":0}},
{"command":"layer.move_to_parent","params":{"nodeId":"${{steps.ring.result.id}}","parentId":"${{steps.frame.result.id}}"}}
```

**After** (1 op, ~200 chars):
```json
{"name":"ring","command":"shape.create_ellipse","params":{"x":10,"y":10,"width":160,"height":160,"parentId":"$frame","stroke":"#FFD70026","strokeWidth":3,"noFill":true}}
```

**Savings:** 79 → ~35 ops. 59% fewer chars per element. 75% fewer ops per element.

### Priority 3: Universal Color Parser

Move color parsing to **Go** (`internal/tools/helpers.go`). Accept any precise format, normalize to Figma-native `{r,g,b}` + `opacity` before sending to WebSocket.

Accept: hex (3/6/8), `{r,g,b}` objects (auto-detect 0-1 vs 0-255), `rgb()`, `rgba()`, `hsl()`.

Reject with clear error: named colors, unrecognized strings, malformed hex.

**Never silently fall back** — errors help the LLM self-correct.

**Implementation location:** Go side, not plugin. Errors caught early, one parser instead of 6 duplicates.

### Priority 4: Command Aliases

Lookup table in `command_routing.go`:

| Alias | Full Command | Saves |
|-------|-------------|-------|
| `frame` | `node.create_frame` | 12 chars |
| `text` | `text.create` | 7 chars |
| `rect` | `shape.create_rectangle` | 19 chars |
| `ellipse` | `shape.create_ellipse` | 16 chars |
| `line` | `shape.create_line` | 13 chars |
| `image` | `shape.create_image` | 14 chars |
| `fill` | `paint.set_solid` | 11 chars |
| `stroke` | `paint.set_stroke` | 12 chars |
| `parent` | `layer.move_to_parent` | 16 chars |
| `autolayout` | `layout.set_auto_layout` | 13 chars |

### Priority 5: Param Shorthand

Normalize in Go before sending to plugin:

| Short | Full | Context |
|-------|------|---------|
| `w` | `width` | all |
| `h` | `height` | all |
| `pid` | `parentId` | all creates |
| `r` | `cornerRadius` | rect/frame |
| `sz` | `fontSize` | text |
| `ff` | `fontFamily` | text |
| `fs` | `fontStyle` | text |
| `lh` | `lineHeight` | text |
| `ls` | `letterSpacing` | text |
| `sw` | `strokeWidth` | line/shape |
| `bg` | `color` (fill) | frame/rect |

### Priority 6: Sensible Defaults

| Parameter | Current Default | Proposed Default | Why |
|-----------|----------------|-----------------|-----|
| `lineHeightUnit` | `PIXELS` | `PERCENT` | #1 text bug — PIXELS is almost never what you want |
| `lineHeight` | none | `140` (when unit=PERCENT) | Reasonable body text default |
| `fontFamily` | `Roboto` | `Inter` | Most common for UI/social |
| `fontStyle` | varies | `Regular` | Safe default |
| `noFill` on create | frame gets white fill | if `noFill:true`, remove fill | Structural frames need this constantly |

### Priority 7: Input Validation with Clear Errors

Replace all silent fallbacks with actionable error messages:

```
Error: "color" value "gold" is a named color. Named colors not accepted — use hex: #FFD700
Error: Unknown command "rect.create". Did you mean "rect" (shape.create_rectangle)?
Error: Step "frame" not found for interpolation "$frame". Available: bg, title, subtitle.
Error: Could not parse color "#GG0000". Invalid hex digit 'G'. Hex uses 0-9 and A-F.
```

### Optional: YAML Input Support

Based on research, YAML is a defensible optional format:
- 10-15% fewer tokens than JSON
- Simpler syntax (no brackets/quotes to track)
- AlphaCodium: "YAML output is far better for code generation"
- 30-50% fewer parsing failures reported

Could be implemented as:
- Auto-detect YAML vs JSON input in batch command
- `batch --format yaml` explicit flag
- Accept YAML in MCP `operations` parameter

**Risk:** Adds Go YAML parser dependency. YAML has indentation sensitivity issues. If we implement priorities 1-7 first (which give 60-70% reduction), YAML adds only 10-15% on top of that — diminishing returns.

**Recommendation:** Implement as future enhancement after priorities 1-7 are validated.

---

## Revised Assessment

### Projected Token Savings

For the 79-op Instagram post batch (21KB):

| Optimization | Ops | Size | Cumulative Reduction |
|-------------|-----|------|:--------------------:|
| Original | 79 | 21KB | — |
| Use existing `parentId`/`color` on create | 52 | 15KB | 29% |
| + Combined ops (stroke, noFill, opacity) | 35 | 11KB | 48% |
| + Short interpolation (`$name`) | 35 | 9KB | 57% |
| + Hex colors everywhere | 35 | 8KB | 62% |
| + Command aliases | 35 | 7.5KB | 64% |
| + Param shorthand | 35 | 6.5KB | 69% |
| + YAML format (optional) | 35 | ~5.5KB | 74% |

**21KB → 6.5KB (JSON) or ~5.5KB (YAML) = 69-74% token reduction. 79 → 35 ops.**

### The Same Batch, Before and After

**Before (current format, 4 ops for one ring):**
```json
{"name":"ring","command":"shape.create_ellipse","params":{"x":1210,"y":10,"width":160,"height":160,"name":"Ring Outer"}},
{"command":"paint.set_stroke","params":{"nodeId":"${{steps.ring.result.id}}","color":{"r":1,"g":0.84,"b":0,"a":0.15},"width":3}},
{"command":"paint.remove_fill","params":{"nodeId":"${{steps.ring.result.id}}","index":0}},
{"command":"layer.move_to_parent","params":{"nodeId":"${{steps.ring.result.id}}","parentId":"${{steps.frame.result.id}}"}}
```

**After (optimized JSON, 1 op):**
```json
{"name":"ring","command":"ellipse","params":{"x":1210,"y":10,"w":160,"h":160,"name":"Ring Outer","pid":"$frame","stroke":"#FFD70026","sw":3,"noFill":true}}
```

**After (optimized YAML, 1 op):**
```yaml
- name: ring
  command: ellipse
  params:
    x: 1210
    y: 10
    w: 160
    h: 160
    name: Ring Outer
    pid: $frame
    stroke: "#FFD70026"
    sw: 3
    noFill: true
```

### Implementation Priority Order

1. **Short interpolation** (`$name`) — trivial Go change, huge QOL
2. **Combined create ops** (inline stroke/opacity/noFill) — plugin handler changes
3. **Universal color parser** — Go-side, replace 6 duplicate parsers
4. **Command aliases** — lookup table in `command_routing.go`
5. **Param shorthand** — normalization in Go before send
6. **Sensible defaults** — lineHeightUnit, fontFamily
7. **Error validation** — replace silent fallbacks with clear messages
8. **(Optional) YAML input** — auto-detect or `--format yaml`

### Implementation Locations

| Optimization | Go Side | Plugin Side |
|-------------|---------|-------------|
| Short interpolation | `internal/batchutil/interpolation.go` | — |
| Combined create ops | — | `plugin/src/handlers/shape.ts`, `node.ts` |
| Universal color parser | `internal/tools/helpers.go` (new) | Keep `parseHexColor()` as backup |
| Command aliases | `internal/ws/command_routing.go` | — |
| Param shorthand | `internal/tools/helpers.go` or `command_routing.go` | — |
| Sensible defaults | `internal/tools/text.go`, `catalog_llm.go` | `plugin/src/handlers/text.ts` |
| Error validation | `internal/tools/helpers.go` | Remove silent fallbacks |
| YAML input | `cmd/ai-happy-design/main.go` (batch cmd) | — |

### Key Files

**Go side:**
- `internal/batchutil/interpolation.go` — interpolation engine
- `internal/ws/command_routing.go` — command routing, aliases would go here
- `internal/tools/helpers.go` — shared helpers, color parser would go here
- `internal/tools/catalog_llm.go` — LLM design guidance, update examples
- `cmd/ai-happy-design/main.go` — batch command CLI

**Plugin side:**
- `plugin/src/handlers/shape.ts` — add inline stroke/opacity/noFill
- `plugin/src/handlers/node.ts` — add inline opacity
- `plugin/src/handlers/text.ts` — lineHeight default fix
- `plugin/src/utils/parseColor.ts` — potential shared color parser (replace 6 copies)

---

## Sources

### Figma API & Plugin Docs
- [Figma Plugin API - Global Objects](https://developers.figma.com/docs/plugins/api/global-objects/)
- [Figma Plugin API - Paint Types](https://developers.figma.com/docs/plugins/api/Paint/)
- [Figma Plugin API - RGB/RGBA](https://developers.figma.com/docs/plugins/api/RGB/)
- [Figma `util.rgb()` docs](https://www.figma.com/plugin-docs/api/properties/figma-util-rgb/) — accepts hex, rgb(), hsl(), lab(); throws on bad input
- [Figma Color Models](https://help.figma.com/hc/en-us/articles/360043042113-Color-models-in-Figma-design) — Hex, RGB, HSB, HSL, CSS

### Figma MCP & Design Tools
- [Figma MCP Server Docs](https://developers.figma.com/docs/figma-mcp-server/)
- [Figma MCP Server Guide (GitHub)](https://github.com/figma/mcp-server-guide/)
- [Framelink / Figma-Context-MCP](https://github.com/GLips/Figma-Context-MCP) — 13.1K stars, compact tree format
- [Figma-Context-MCP Explained (Skywork)](https://skywork.ai/blog/figma-context-mcp-mcp-server-ai-integration/)
- [UIKit Data Exporter - TOON format](https://www.figma.com/community/plugin/1567884041675515041/uikit-data-exporter)
- [Figma to AI JSON](https://www.figma.com/community/plugin/1587577656366372788/figma-to-ai-json)
- [Figma Raw: Export Design Data for AI/LLM](https://www.figma.com/community/plugin/1491678546144854232/figma-raw-export-design-data-for-ai-llm-agents)

### Design Token Standards
- [W3C Design Tokens Color Module 2025.10](https://www.designtokens.org/tr/drafts/color/)
- [Tokens Studio Color Handling](https://docs.tokens.studio/manage-tokens/token-types/color)

### LLM Structured Output Research
- [Claude Structured Outputs Docs](https://platform.claude.com/docs/en/build-with-claude/structured-outputs)
- [Claude Prompting Best Practices](https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/claude-4-best-practices)
- [Which Nested Data Format Do LLMs Understand Best? (ImprovingAgents, Oct 2025)](https://www.improvingagents.com/blog/best-nested-data-format/) — 1000-question benchmark, 3 models
- [AlphaCodium: Flow Engineering](https://arxiv.org/abs/2401.08500) — "YAML output is far better for code generation"
- [StructEval: Benchmarking LLM Structured Outputs](https://arxiv.org/html/2505.20139v1)
- [Format Restrictions Impact on LLM Performance](https://arxiv.org/html/2408.02442v1) — 10-15% reasoning penalty
- [Are LLMs Ready for TOON?](https://arxiv.org/html/2601.12014)

### LLM Format Comparison
- [LLM Reliability: JSON vs YAML (Sean Ryan)](https://medium.com/@mr.sean.ryan/llm-reliability-json-vs-yaml-22c58d7f51f6)
- [Structured Output: YAML vs JSON (Unalarming)](https://unalarming.com/structured-output-yaml-vs-json)
- [YAML Over JSON in LLM Applications (BlogOS)](https://blog.tashif.codes/blog/JSON-YAML-LLM)

### ASCII Art & LLM Limitations
- [Why LLMs Suck at ASCII Art (TDS)](https://medium.com/data-science/why-llms-suck-at-ascii-art-a9516cb880d5)
- [ASCIIBench: Evaluating ASCII Understanding](https://arxiv.org/html/2512.04125) — GPT-4 at 25.19% accuracy

### LLM JSON Generation Techniques
- [Super JSON Mode](https://github.com/varunshenoy/super-json-mode) — parallel field generation
- [Jsonformer](https://github.com/1rgs/jsonformer) — schema pre-filling
- [Fast JSON Decoding with Compressed FSM (SGLang)](https://lmsys.org/blog/2024-02-05-compressed-fsm/)
