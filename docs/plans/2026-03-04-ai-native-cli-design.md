# AI-Native CLI Design — Schema-Driven Validation & Structured Errors

**Date:** 2026-03-04
**Status:** Approved
**Philosophy:** ShipTypes — "Ship types, not docs." Make incorrect usage impossible. LLMs get it right on the first call.

## Summary

Transform ai-happy-design from a doc-driven tool (92KB prose catalog that LLMs must read and follow) into a schema-driven tool where the CLI validates, corrects, and enriches every command before execution. Remove MCP entirely — CLI via bash tool-use is the only LLM interaction path.

## Architecture

### Current
```
LLM → MCP (stdio) → Go tools → WebSocket relay → Figma plugin
LLM → CLI (bash) → Go commands → WebSocket relay → Figma plugin
```

### New
```
LLM → CLI (bash) → Schema validation → Design lint → WebSocket relay → Figma plugin
                         ↑                    ↑
                    Schema defs          compute_tokens
                  (source of truth)     (auto-injected)
```

### Removed
- `internal/mcp/` — entire MCP server
- `internal/tools/*.go` — MCP tool wrappers (node.go, text.go, paint.go, etc.)
- `internal/tools/registry.go` — MCP tool registration
- `internal/tools/describe.go` — MCP describe endpoint
- `mcp` and `register` CLI subcommands
- mcp-go dependency

### Added
- `internal/schema/` — canonical schema definitions per command
- `internal/validate/` — schema-based pre-execution validation
- `internal/designlint/` — design quality pre-checks
- Structured error/warning response format
- Plugin handler expansions (min/max, truncation, dash strokes, etc.)
- `llms.txt` and `llms-full.txt` on aihappydesign.com

### Stays
- `internal/ws/` — relay (server + client)
- `internal/batchutil/` — normalization pipeline (enhanced)
- `plugin/` — Figma plugin (enhanced)
- `cmd/ai-happy-design/` — CLI commands
- `internal/tools/catalog_llm.go` — design intelligence (shrinks from 92KB to ~25-30KB)
- `internal/tools/design_tokens.go` — compute_tokens
- `internal/tools/bulk.go` — batch execution orchestration

---

## 1. Schema System

### Schema Definition

Each command.action gets a Go struct definition:

```go
// internal/schema/text_create.go
var TextCreate = Schema{
    Command:     "text.create",
    Aliases:     []string{"text"},
    Description: "Create a text node",
    Params: []Param{
        {Name: "text", Type: "string", Required: true, Aliases: []string{"content"},
         Desc: "Text content to display"},
        {Name: "parentId", Type: "string", Required: true, Aliases: []string{"pid"},
         Desc: "Parent frame ID", Pattern: `^[0-9]+:[0-9]+$`},
        {Name: "fontSize", Type: "number", Aliases: []string{"sz"},
         Desc: "Font size in pixels", Min: ptr(4.0), Max: ptr(500.0), Default: 16.0,
         SemanticTokens: true},
        {Name: "fontFamily", Type: "string", Aliases: []string{"ff"},
         Desc: "Font family name", Default: "Inter"},
        {Name: "fontStyle", Type: "string", Aliases: []string{"fs"},
         Desc: "Font style", Default: "Regular",
         Enum: []string{"Thin", "ExtraLight", "Light", "Regular", "Medium", "SemiBold", "Bold", "ExtraBold", "Black"}},
        {Name: "color", Type: "string",
         Desc: "Text color hex", Pattern: `^#[0-9A-Fa-f]{3,8}$`, Default: "#000000"},
        {Name: "lineHeight", Type: "number", Aliases: []string{"lh"},
         Desc: "Line height percentage", Min: ptr(50.0), Max: ptr(300.0),
         AutoFix: "lineHeightUnit:PERCENT"},
        {Name: "width", Type: "number", Min: ptr(1.0), Max: ptr(10000.0),
         Desc: "Text box width (enables wrapping)"},
        {Name: "textAlign", Type: "string",
         Desc: "Horizontal alignment", Enum: []string{"LEFT", "CENTER", "RIGHT", "JUSTIFIED"}},
    },
}
```

### Schema Registry

```go
// internal/schema/registry.go
var All = []Schema{
    TextCreate, TextSetContent, TextSetFont,
    NodeCreateFrame, NodeMove, NodeResize, NodeModify,
    PaintSetSolid, PaintSetGradient, PaintSetImage, PaintSetStroke,
    ShapeCreateRectangle, ShapeCreateEllipse, ShapeCreateImage,
    LayoutSetAutoLayout, LayoutSetPadding, LayoutSetSizing,
    EffectAddShadow, EffectAddBlur, EffectApplyGlass,
    // ~50 total
}

func Lookup(command string) *Schema { ... } // by command or alias
```

### What Schemas Generate

| Output | How | Used By |
|---|---|---|
| CLI validation | Schema loader validates batch JSON params at runtime | `batch`, `command` |
| CLI help | `ai-happy-design help text.create` prints params table | Developers, LLMs |
| Schema dump | `ai-happy-design schema text.create --json` prints exact JSON shape | LLMs |
| Error messages | "fontSize must be 4-500, got -10. Default: 16" | LLM error recovery |
| Auto-fix rules | Clamp to min/max, apply defaults, resolve aliases | Warn+fix mode |
| Skill reference | Auto-generate command reference section of SKILL.md | Claude skill |
| llms-full.txt | `ai-happy-design schema --all --llms-txt` | Web discoverability |
| Examples | Generate valid example JSON from defaults + types | `examples` command |

### Pipeline Position

```
Batch JSON arrives
  → batchutil.Fix() (markdown, comments, dict unwrap)
  → batchutil.Normalize() (aliases, CSS)
  → batchutil.ResolveTokenAliases() (semantic tokens)
  → schema.Validate(ops) ← NEW
      1. Look up schema by command name
      2. Check required params present
      3. Check types match
      4. Check enums (with fuzzy matching)
      5. Check bounds (min/max)
      6. Check patterns (hex regex)
      7. Check dependencies (lineHeight → lineHeightUnit)
      8. Apply defaults for missing optional params
      9. Collect warnings (warn+fix mode)
  → designlint.Check(ops) ← NEW
  → Execute on Figma
```

---

## 2. Structured Errors

### Error Object Shape

Every error and warning uses the same structure:

```json
{
  "step": 3,
  "name": "card_title",
  "phase": "schema|designLint|execution",
  "code": "FONT_SIZE_BELOW_MIN",
  "param": "fontSize",
  "message": "fontSize must be 4-500, got -10",
  "got": -10,
  "expected": { "min": 4, "max": 500, "default": 16 },
  "fix": 16,
  "applied": true
}
```

### Error Codes

#### Schema Errors

| Code | When | Fix Strategy |
|---|---|---|
| `REQUIRED_MISSING` | Required param absent | `fix: null` |
| `TYPE_MISMATCH` | String where number expected | `fix:` parsed value if possible |
| `ENUM_INVALID` | Value not in allowed set | `fix:` closest Levenshtein match |
| `BELOW_MIN` | Number below minimum | `fix:` minimum value |
| `ABOVE_MAX` | Number above maximum | `fix:` maximum value |
| `PATTERN_MISMATCH` | String doesn't match regex | `fix:` best-effort parse (named colors → hex) |
| `UNKNOWN_COMMAND` | Command not in schema registry | `fix:` fuzzy match suggestion |
| `UNKNOWN_PARAM` | Param not in schema for this command | `fix:` Levenshtein match |
| `DEPENDENCY_MISSING` | Conditional param missing | `fix:` auto-apply dependency |

#### Design Lint Warnings

| Code | When | Fix Strategy |
|---|---|---|
| `TEXT_TOO_SMALL` | fontSize below caption tier | `fix:` caption size |
| `TEXT_NO_HIERARCHY` | All text sizes within 1 tier | `fix: null` (warn only) |
| `LOW_CONTRAST` | WCAG ratio below 4.5:1 | `fix:` adjusted color passing AA |
| `GRADIENT_ALPHA` | Gradient stop has accidental transparency | `fix:` opaque hex |
| `PADDING_TOO_SMALL` | Side padding below 4% canvas | `fix:` tokens.sidePadding |
| `SPACING_EXTREME` | itemSpacing outside token bounds | `fix:` nearest token value |
| `RADIUS_OVERFLOW` | cornerRadius > min(w,h)/2 | `fix:` clamped value |
| `ELEMENT_DENSITY` | Too many children in frame | `fix: null` (warn only) |
| `MISSING_NAME` | Node without semantic name | `fix:` auto-generated name |
| `DUPLICATE_STEP_NAME` | Two steps with same name | `fix:` appended index |
| `STRUCTURAL_FILL` | Auto-layout frame with default fill | `fix: null` (suggest structural:true) |

#### Execution Errors

| Code | When | Fix Strategy |
|---|---|---|
| `NODE_NOT_FOUND` | nodeId doesn't exist | `fix: null` |
| `FONT_NOT_FOUND` | fontFamily not available | `fix:` closest common font |
| `TIMEOUT` | Plugin didn't respond in 300s | `fix: null` |
| `PLUGIN_ERROR` | Figma plugin threw | `fix: null`, pass-through error |
| `INTERPOLATION_FAILED` | Step reference unresolved | `fix:` suggested step name |

### Special Auto-Fixes

**Named CSS colors:** "red" → "#FF0000", "white" → "#FFFFFF", etc. (140 entries)

**Enum fuzzy matching:** Case-insensitive + Levenshtein. "bold" → "Bold", "semi-bold" → "SemiBold", "HORIZONAL" → "HORIZONTAL"

### Batch Response Format

```json
{
  "ok": true,
  "preValidation": {
    "schema": { "warnings": [...], "errors": [...], "fixed": 2, "blocked": 0 },
    "designLint": {
      "canvas": { "width": 1080, "height": 1350, "type": "portrait" },
      "tokens": { "caption": 36, "body": 48, "sidePadding": 72 },
      "warnings": [...],
      "fixed": 3,
      "score": {
        "readability": 8,
        "contrast": 9,
        "spacing": 7,
        "hierarchy": 6,
        "overall": 7.5
      }
    }
  },
  "steps": [...],
  "summary": { "total": N, "succeeded": N, "failed": N },
  "timing": { "totalMs": N, "opsPerSec": N }
}
```

---

## 3. Design Lint

Pre-execution design quality checks. Operates on batch JSON + computed tokens. No Figma needed.

### Canvas Detection

1. Scan ops for first root frame (no `parentId`, has numeric `width` + `height`)
2. If found → auto-compute tokens for that canvas
3. If not found → skip token-dependent checks, run structural checks only
4. If `slide`/`banner` composite → extract canvas from params

### Checks

#### Tier 1: Text Readability
- **Text size floor:** `fontSize >= tokens.caption` → fix: bump to caption
- **Text size hierarchy:** 2+ text sizes must span ≥ 2 scale tiers → warn only
- **Line height range:** 100-200% for body, 100-150% for headings → clamp
- **Text width ratio:** `width <= tokens.contentWidth` → clamp
- **Letter spacing bounds:** -5 to 20 → clamp

#### Tier 2: Color & Contrast
- **Contrast ratio:** text color vs nearest ancestor background, WCAG AA (4.5:1) → fix: adjusted color
- **Gradient transparency:** hex with alpha suffix → warn with opaque fix
- **Near-identical colors:** deltaE < 5 → warn

#### Tier 3: Spacing & Layout
- **Side padding ratio:** ≥ 4% canvas width → fix: tokens.sidePadding
- **Spacing bounds:** itemSpacing within token range → clamp
- **Corner radius overflow:** ≤ min(w,h)/2 → clamp
- **Element density:** > 8 direct children → warn
- **Padding consistency:** siblings same padding → warn

#### Tier 4: Structural
- **Missing name:** auto-name from context
- **Structural frame fill:** auto-layout + no explicit color → suggest structural:true
- **Orphan nodes:** no parentId, not root frame → warn
- **Duplicate step names:** → error (breaks interpolation)

### Design Score

| Axis | Weight | What It Measures |
|---|---|---|
| readability | 30% | Text sizes vs scale, line heights |
| contrast | 25% | WCAG ratios for text/background pairs |
| spacing | 25% | Padding ratios, item spacing, grid alignment |
| hierarchy | 20% | Spread of text sizes, weight variation |
| **overall** | — | Weighted average |

### Flags

- `--lint` — enabled by default
- `--no-lint` — skip design lint
- `--strict-quality` — fail if overall score < 7
- `--lint-report` — detailed report to stderr (default: on)

---

## 4. Plugin Feature Expansions

Figma API capabilities that exist but the plugin doesn't expose.

### Tier 1: Prevent Broken Designs

| Feature | Handler | Figma API | Schema Params |
|---|---|---|---|
| Min/max sizing | node.ts, layout.ts | `minWidth`, `maxWidth`, `minHeight`, `maxHeight` | 4 numbers, min: 0, max: 10000 |
| Text truncation | text.ts | `textTruncation`, `maxLines` | integer maxLines (1-100), auto-sets truncation:ENDING |
| Clips content | node.ts | `clipsContent` | boolean, default true for auto-layout |
| Constrain proportions | shape.ts, node.ts | `constrainProportions` | boolean |

### Tier 2: Better Visual Quality

| Feature | Handler | Figma API | Schema Params |
|---|---|---|---|
| Stroke dash | paint.ts | `dashPattern` | number array [dash, gap] |
| Stroke cap | paint.ts, shape.ts | `strokeCap` | enum: NONE/ROUND/SQUARE/ARROW_LINES/ARROW_EQUILATERAL |
| Stroke join | paint.ts, shape.ts | `strokeJoin` | enum: MITER/BEVEL/ROUND |

### Structural Shorthand (Not Figma API — Plugin Convenience)

| Feature | Handler | Behavior |
|---|---|---|
| `structural: true` | node.ts create_frame | Removes default fill, sets clipsContent: true |

Total: ~60 lines of plugin TypeScript.

---

## 5. llms.txt

### aihappydesign.com/llms.txt (Hand-Maintained, ~500 Tokens)

Concise overview: what the tool is, quick start, key commands, design workflow, link to full reference.

### aihappydesign.com/llms-full.txt (Auto-Generated from Schemas)

Complete command reference with parameter tables, constraints, defaults, and examples.

Generated by: `ai-happy-design schema --all --llms-txt > llms-full.txt`

Deployed on every release alongside the binary.

---

## 6. Catalog Shrinkage

| Section | Current | After |
|---|---|---|
| Execution Rules (~2KB) | Prose rules | **Delete** — encoded in pipeline |
| CSS Property Support (~3KB) | Translation table | **Delete** — in cssnorm.go + schemas |
| Common Mistakes (~4KB) | Workarounds | **Delete entirely** — every mistake is now a schema constraint |
| First-Pass Guardrails (~3KB) | LLM instructions | **Delete** — becomes design lint |
| Batch Observability (~2KB) | Output docs | **Delete** — response format is self-documenting |
| Workflow (~2KB) | Create vs edit rules | **Shrink to ~500 bytes** |
| Design Thinking (~4KB) | Creative guidance | **Keep** |
| Visual Hierarchy (~3KB) | Composition theory | **Keep** |
| Design Decisions (~4KB) | When to use effects | **Keep** |
| Layer Organization (~2KB) | Naming patterns | **Keep** |
| Design Quality Checklist (~2KB) | Quality reminders | **Shrink to ~500 bytes** |

**92KB → ~25-30KB.** Creative guidance only.

Served via: `ai-happy-design guide [--section X] [--json]`

---

## 7. CLI Command Structure (Post-MCP)

```
ai-happy-design
  ├── ws                         # Start WebSocket relay
  ├── command <cmd> [json]       # Single command execution
  ├── batch [json|file]          # Batch execution (main LLM path)
  ├── tools [--json]             # List all commands (from schema registry)
  ├── schema <cmd> [--json]      # Print exact param schema for a command
  ├── schema --all --llms-txt    # Generate llms-full.txt
  ├── examples [category]        # Pre-built batch examples
  ├── guide [--section X]        # Design intelligence (from catalog)
  ├── help <cmd>                 # Parameter reference (from schema)
  ├── validate [json|file]       # Dry-run: schema + lint, no execution
  ├── relay start|stop|status|logs
  ├── setup                      # Extract embedded plugin
  ├── upgrade                    # Self-update
  ├── config set|get|reset
  └── benchmark exec|pipe|compare
```

---

## 8. Implementation Phases

### Phase 1: Schema System + Structured Errors (Foundation)

1. `internal/schema/` — Schema types, Param struct, registry
2. Define ~15 most-used command schemas
3. `internal/validate/` — Schema validator
4. Structured error format with codes, fix suggestions
5. Wire into batch pipeline after normalize, before execute
6. Named CSS color resolution (140 colors)
7. Enum fuzzy matching (case-insensitive + Levenshtein)
8. `schema` CLI command
9. `validate` dry-run CLI command

### Phase 2: Design Lint + Score

1. `internal/designlint/` — Lint engine
2. Text readability checks
3. Contrast checking (WCAG AA)
4. Spacing checks
5. Structural checks
6. Design score computation
7. Wire into pipeline after schema validation
8. `--no-lint`, `--strict-quality` flags

### Phase 3: Plugin Feature Expansions (Parallel with Phase 1)

1. Min/max sizing in node.ts, layout.ts
2. Text truncation in text.ts
3. Clips content in node.ts
4. Constrain proportions in shape.ts, node.ts
5. Stroke dash pattern in paint.ts
6. Stroke cap/join in paint.ts, shape.ts
7. Structural frame shorthand in node.ts
8. Schemas for all new params
9. Plugin build verification

### Phase 4: MCP Removal + Catalog Shrink + llms.txt

1. Delete `internal/mcp/`, MCP tool wrappers, mcp-go dependency
2. Remove `mcp` and `register` subcommands
3. Shrink catalog_llm.go (92KB → ~25-30KB)
4. Add `guide` CLI command
5. Add `schema --all --llms-txt` generation
6. Write llms.txt, generate and deploy llms-full.txt
7. Update SKILL.md, AGENTS.md, CLAUDE.md

### Dependency Graph

```
Phase 1 (schemas + errors) ──→ Phase 2 (design lint)
                            ──→ Phase 4 (MCP removal + llms.txt)

Phase 3 (plugin features)  ──→ Phase 4 (schemas for new features)
```

Phases 1 and 3 run in parallel. Phase 2 needs Phase 1. Phase 4 needs all three.

---

## Error Mode

**Default:** Warn and fix. Log the issue + auto-correct to nearest valid value. LLM sees what was fixed and learns.

**Backwards compatibility:** Break freely. LLMs regenerate JSON every time. Ship the better API.

## Key Insight

Every rule in the catalog that says "don't do X" becomes code that prevents X. The catalog shrinks to creative guidance only — color theory, visual hierarchy, composition. The mechanical rules become invisible, baked into the tools.
