# Speed Optimization: Composite Commands, HTML Extractor & E2E Benchmarks

**Date**: 2026-02-18
**Status**: Approved
**Goal**: Reduce prompt-to-Figma time from ~7 minutes to under 2 minutes

## Problem

Creating 36 frames (30 slides + 6 email banners) from an HTML spec took ~7 minutes. Bottlenecks:
1. **LLM deliberation** (3-4 min) — deciding approach, writing Python scripts
2. **Verbose batch JSON** — each element requires 5-10 fields, no higher-level abstractions
3. **No HTML pre-digestion** — LLM must manually translate CSS to Figma params
4. **No end-to-end measurement** — can't track or optimize the full pipeline

## Design: Three Implementation Phases

### Phase 1: Composite Batch Commands

Higher-level commands that expand into multiple primitive batch operations during normalization.

#### Slide & Banner (Container Templates)

```json
{"command": "slide", "params": {
  "name": "P1 - Slide 1",
  "canvas": "1080x1350",
  "bg": "#0C1E2C",
  "gradient": {"angle": 150, "stops": [{"pos": 0, "color": "#0C1E2C"}, {"pos": 1, "color": "#14344A"}]},
  "elements": [
    {"type": "eyebrow", "text": "AMC  ·  Ramadan 2026", "color": "#7FBCD2"},
    {"type": "headline", "text": "The Care\nBehind\nthe Care", "tier": "hero"},
    {"type": "bar", "color": "#029056"},
    {"type": "body", "text": "Supporting tagline here", "color": "#FAFCFB6B"}
  ]
}}
```

```json
{"command": "banner", "params": {
  "name": "Email 1 — Launch",
  "canvas": "1200x400",
  "bg": "#0C1E2C",
  "gradient": {"angle": 150, "stops": [...]},
  "dividerX": 200,
  "elements": [
    {"type": "headline", "text": "The Care Behind the Care"},
    {"type": "subtitle", "text": "AMC  ·  Ramadan 2026", "textCase": "UPPER"}
  ]
}}
```

#### Extended Element Types

| Type | Description | Key Params |
|------|-------------|------------|
| `eyebrow` | Small uppercase label | text, color, letterSpacing |
| `headline` | Main heading | text, tier (hero/title/heading), color |
| `body` | Body/paragraph text | text, color |
| `bar` | Accent divider bar | color, width (default 108px) |
| `cta` | Call-to-action button | text, bg, color, style (pill/rounded) |
| `url` | URL text (muted, centered) | text |
| `counter` | Slide counter "1 / 5" | current, total |
| `stats` | Stats row (250+ / Chaplains) | items: [{value, label}] |
| `progress` | Progress bar with labels | raised, goal, color |
| `avatar` | Circular masked image | imageData, size |
| `pattern` | Geometric/decorative pattern | style (dots/lines/circles), color, opacity |
| `stars` | Star rating row | count (1-5), color, size |
| `arabic` | Arabic calligraphy text | text, size, color, fontFamily (default Amiri) |

#### text.set_range_style

Style character ranges within existing text nodes:

```json
{"command": "text.set_range_style", "params": {
  "nodeId": "123:456",
  "ranges": [
    {"match": "250+", "bold": true, "color": "#7FBCD2", "fontSize": 64},
    {"match": "chaplains", "italic": true},
    {"start": 0, "end": 5, "color": "#C88B0A"}
  ]
}}
```

Match by string (`match`) or position (`start`/`end`). Supports: bold, italic, color, fontSize, fontFamily, fontStyle, letterSpacing, lineHeight, textDecoration.

#### Expansion

Composite commands expand during `batchutil/normalize.go` processing:
- One `slide` command → 8-15 primitive ops (frame + gradient + elements)
- Design tokens auto-computed from `canvas` param
- Element positioning uses token-based layout rules
- Result: LLM generates ~10 JSON objects instead of ~40

### Phase 2: HTML Extractor

New CLI command: `ai-happy-design extract`

Two modes:

#### Mode A: Lightweight Go Parser (Default)

```bash
ai-happy-design extract social-posts.html --canvas 1080x1350
```

Uses `golang.org/x/net/html` + regex CSS parser:
- Parses HTML structure (div hierarchy, classes, text content)
- Extracts inline styles and `<style>` block rules
- Maps CSS properties to Figma params using existing css-to-figma mapping
- Outputs batch JSON ready for execution

Handles: backgrounds, gradients, text styling, flexbox layout, padding, margins, border-radius, shadows.
Limitations: No computed styles, no JS-dependent layouts, no web fonts.

#### Mode B: Headless Chrome (`--computed`)

```bash
ai-happy-design extract social-posts.html --computed --canvas 1080x1350
```

Uses `chromedp` (headless Chrome):
- Renders HTML in real browser
- Extracts `getComputedStyle()` for every element
- Gets actual bounding boxes (`getBoundingClientRect()`)
- Resolves CSS variables, calc(), media queries
- Handles font loading, flexbox/grid computed positions

Trade-off: Requires Chrome installed, ~2-3s overhead per page, but handles 100% of CSS.

#### Output

Both modes output the same batch JSON format. The `--computed` mode is more accurate but slower.

```bash
# Extract → Execute pipeline
ai-happy-design extract input.html --canvas 1080x1350 | ai-happy-design batch --stdin

# Or save intermediate JSON
ai-happy-design extract input.html -o /tmp/batch.json
ai-happy-design batch /tmp/batch.json
```

### Phase 3: End-to-End Benchmark Harness

New CLI command: `ai-happy-design benchmark`

#### The Pipeline

```
Prompt → [LLM generates JSON] → [CLI executes batch] → [Figma renders] → Done
           Phase A                  Phase B                Phase C
```

All three phases timed independently AND as total.

#### Subcommands

```bash
# End-to-end: prompt → Figma (uses Cerebras by default)
ai-happy-design benchmark e2e \
  --prompt "Create a 5-slide carousel about Muslim chaplains" \
  --canvas 1080x1350 \
  --runs 3

# Just extraction: HTML → JSON
ai-happy-design benchmark extract \
  --input social-posts.html \
  --mode lightweight \
  --runs 5

# Just execution: JSON → Figma
ai-happy-design benchmark exec \
  --input /tmp/slides/post1.json \
  --runs 5

# Compare modes side-by-side
ai-happy-design benchmark compare \
  --input social-posts.html \
  --runs 3
```

#### E2E Timing

The `e2e` subcommand orchestrates the full pipeline:

1. **Phase A (LLM Generation)**: CLI calls configured LLM API (Cerebras default). Sends prompt + design tokens + batch format spec. Timer: API call start → valid JSON parsed.
2. **Phase B (CLI Execution)**: Returned JSON goes straight into batch runner. Timer: full batch execution including WebSocket round-trips.
3. **Phase C (Verification)**: Optional export + visual check. Timer covers export call.

#### Piped Mode (External LLMs)

When using Claude, GPT, or any external LLM:

```bash
# Manual phase-a timing
START=$(date +%s%N)
BATCH_JSON=$(curl -sS cerebras-api... | jq -r '.choices[0].message.content')
LLM_MS=$(( ($(date +%s%N) - START) / 1000000 ))

echo "$BATCH_JSON" | ai-happy-design benchmark pipe --phase-a-ms $LLM_MS
```

#### Output Format

```
╭──────────────────────────────────────────────╮
│  AI Happy Design — End-to-End Benchmark      │
├──────────────────────────────────────────────┤
│  Prompt: "Create a 5-slide carousel..."      │
│  Canvas: 1080×1350  │  Runs: 3               │
├──────────────────────────────────────────────┤
│  Phase A (LLM Gen)     │  avg 4.2s  ± 0.3s   │
│    └ tokens: 2,847 in / 8,420 out             │
│    └ model: qwen-3-235b                       │
│  Phase B (CLI Exec)    │  avg 6.1s  ± 0.4s    │
│    └ ops: 42  │  throughput: 6.9 ops/s        │
│  Phase C (Verify)      │  avg 1.8s  ± 0.2s    │
│                                               │
│  TOTAL E2E             │  avg 12.1s ± 0.6s    │
│  Errors: 0/3 runs                             │
╰──────────────────────────────────────────────╯
```

#### Comparison Table

```
┌─────────────────┬──────────┬──────────┬──────────┬──────────┐
│ Method          │ Extract  │ LLM Gen  │ Execute  │ Total    │
├─────────────────┼──────────┼──────────┼──────────┼──────────┤
│ HTML→lightweight│ 0.2s     │ —        │ 6.1s     │ 6.3s     │
│ HTML→computed   │ 2.8s     │ —        │ 6.0s     │ 8.8s     │
│ Prompt→Cerebras │ —        │ 4.2s     │ 6.1s     │ 10.3s    │
│ Prompt→Claude   │ —        │ 8.5s     │ 5.9s     │ 14.4s    │
└─────────────────┴──────────┴──────────┴──────────┴──────────┘
```

#### Metrics Captured

| Metric | Source |
|--------|--------|
| LLM latency | API call start → valid JSON parsed |
| LLM tokens | Input/output token counts from API response |
| LLM cost | Calculated from model pricing |
| Extraction time | HTML parse start → JSON output |
| Batch execution | First WS message → last response |
| Ops throughput | Total ops ÷ execution seconds |
| Error rate | Failed ops ÷ total ops across runs |
| Total E2E | Prompt input → Figma complete |

## Implementation Order

1. Phase 1 first — immediate impact on every workflow
2. Phase 3 next — need benchmarks before optimizing Phase 2
3. Phase 2 last — validated against benchmarks

## Success Criteria

- 36-frame AMC campaign: < 2 minutes end-to-end (down from 7)
- Single carousel (5 slides): < 30 seconds
- Benchmark reproducible across runs (< 15% variance)
