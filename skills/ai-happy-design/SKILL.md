---
name: ai-happy-design
description: Design in Figma using the ai-happy-design CLI. Use when asked to create, edit, or export Figma designs — social media posts, cards, layouts, UI mockups, posters, or any visual design. Triggers on "design in Figma", "create a Figma layout", "make an Instagram post", "export from Figma", "Figma design", or any task involving ai-happy-design. Also use when asked to batch-create multi-element Figma compositions. Always use the CLI (`ai-happy-design command` / `ai-happy-design batch`) unless the user explicitly requests MCP mode.
---

# AI Happy Design

Design in Figma via the `ai-happy-design` CLI. One Go binary provides CLI commands, MCP tools, and a WebSocket relay to a Figma plugin. **Always use CLI mode by default.** Only use MCP mode if the user explicitly asks for it.

## CLI Usage

**Single command:**
```bash
ai-happy-design command <action> '<json-params>'
```

**Batch (multiple operations):**
```bash
ai-happy-design batch '<json-array>'
# or from file:
ai-happy-design batch ops.json
# or piped:
cat ops.json | ai-happy-design batch
```

**Multi-batch (multiple files at once):**
```bash
ai-happy-design batch file1.json file2.json file3.json
ai-happy-design batch ./carousels/           # all .json files in directory
ai-happy-design batch slides/*.json          # glob pattern
ai-happy-design batch f1.json f2.json --parallel  # run concurrently (max 4)
```

**Useful flags:**
- `--live` — print progress events while running
- `--fail-fast` — stop at first failure
- `-o path.png` — save exported image to file
- `--port <N>` — use a custom relay port (default 3055)
- `--allow-overlap` — skip auto-placement (batch only)
- `--parallel` — run multi-file batches concurrently (max 4)
- `--compress-images` — compress imageData before sending (requires ImageMagick)

## Critical: Think in HTML/CSS First

Mentally draft every design as HTML/CSS first, then translate to Figma commands.

```
CSS property              → Figma command
display: flex             → node.create_frame {layoutMode: VERTICAL, ...}
gap: 16px                 → node.create_frame {itemSpacing: 16}
padding: 24px             → node.create_frame {padding: 24}
background-color          → node.create_frame {color: "#1a1a1a"}
border-radius: 16px       → node.create_frame {cornerRadius: 16}
border: 1px solid         → node.create_frame {stroke: "#333", strokeWidth: 1}
box-shadow                → effect.add_shadow
backdrop-filter: blur()   → effect.add_blur {blurType: "BACKGROUND_BLUR"}
mask                      → node.set_mask
```

Full CSS-to-Figma mapping: See [references/css-to-figma.md](references/css-to-figma.md)

## Workflow

1. **Get design tokens**: `ai-happy-design command design.compute_tokens '{"width": 1080, "height": 1080}'`
2. **Find free space**: `ai-happy-design command document.find_free_space '{"width": 1080, "height": 1080}'`
3. **Think in CSS**: Picture the design as a webpage. What HTML elements? What CSS?
4. **Create in one step**: Use `node.create_frame` with all params — no need for separate commands
5. **Batch create**: Build a JSON array of operations. Send via `ai-happy-design batch`
6. **Balance check**: All siblings MUST match — same height, padding, spacing, radius, text sizes
7. **Export & verify**: `ai-happy-design command export.image '{"nodeId": "...", "scale": 2}' -o output.png`

**Auto-placement**: Batch mode auto-calls `find_free_space` before placing root frames, so your design won't overlap existing work. Use `--allow-overlap` to skip.

## One-Step Frame Creation

Create fully-configured frames in a single command:

```bash
ai-happy-design command node.create_frame '{
  "name": "Card",
  "x": 0, "y": 0, "width": 400, "height": 300,
  "parentId": "123:456",
  "color": "#1a1a1a",
  "cornerRadius": 16,
  "stroke": "#333333",
  "strokeWidth": 1,
  "layoutMode": "VERTICAL",
  "itemSpacing": 16,
  "padding": 24,
  "primaryAxisAlign": "MIN",
  "counterAxisAlign": "CENTER",
  "primaryAxisSizing": "FIXED",
  "counterAxisSizing": "FIXED",
  "clipsContent": true,
  "noFill": false
}'
```

For structural/wrapper frames, use `"noFill": true` to remove the default white fill.

## Design Tokens (Proportional Sizing)

Always call `design.compute_tokens` first to get sizes proportional to the canvas:

```bash
ai-happy-design command design.compute_tokens '{"width": 1080, "height": 1080}'
```

Returns a modular type scale (perfect fourth = 1.333 ratio):
- **display** (200px at 1080w) — statement text
- **hero** (152px) — hero headlines
- **title** (112px) — section titles
- **heading** (84px) — card headings
- **subheading** (64px) — subtitles
- **body** (48px) — paragraph text
- **caption** (36px) — labels, metadata

Also returns spacing tokens (xs through xxxl) on 8px grid, and layout classification.

For print: add `"dpi": 300` to get point-based sizes.

## Batch Operations

> **Schema Contract** — this is a custom format. Read carefully before generating.
>
> Every element MUST have exactly these three keys:
> ```json
> {"name": "stepName", "command": "domain.action", "params": {"key": "value"}}
> ```
> **WRONG** (common mistakes that break silently):
> ```json
> {"type": "node.create_frame", "x": 0, "color": "#fff"}           ← "type" not "command", params not nested
> {"command": "frame", "name": "bg", "x": 0, "y": 0}               ← x/y at top level, not inside params
> {"command": "text.create", "params": {"text": "Hi"}, "name": ""}  ← empty name breaks interpolation
> ```
> **RIGHT**:
> ```json
> {"name": "bg", "command": "node.create_frame", "params": {"x": 0, "y": 0, "width": 1080, "height": 1080}}
> {"name": "title", "command": "text.create", "params": {"text": "Hello", "parentId": "${{steps.bg.result.id}}", "fontSize": 84}}
> ```

Build a JSON array with named steps and `${{steps.name.result.id}}` interpolation:

```bash
ai-happy-design batch '[
  {"name":"root","command":"frame","params":{"name":"Post","x":0,"y":0,"w":1080,"h":1080,"bg":"#1a1a1a"}},
  {"name":"content","command":"frame","params":{"name":"Content","pid":"${{steps.root.result.id}}","w":1080,"h":1080,"noFill":true,"layoutMode":"VERTICAL","itemSpacing":24,"padding":64,"primaryAxisAlign":"CENTER","counterAxisAlign":"CENTER"}},
  {"command":"text","params":{"text":"Hello World","pid":"${{steps.content.result.id}}","sz":80,"fontStyle":"Bold","color":"#ffffff","textAlign":"CENTER","w":952}}
]'
```

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
Use `ai-happy-design validate` to catch schema errors before sending.

### Batch Aliases

| Alias | Full Command |
|-------|-------------|
| `frame` | `node.create_frame` |
| `rect` | `shape.create_rectangle` |
| `ellipse` | `shape.create_ellipse` |
| `text` | `text.create` |
| `image` | `shape.create_image` |
| `fill` | `paint.set_solid` |
| `stroke` | `paint.set_stroke` |
| `gradient` | `paint.set_gradient` |
| `shadow` | `effect.add_shadow` |
| `blur` | `effect.add_blur` |
| `glass` | `effect.apply_glass` |
| `noise` | `effect.add_noise` |
| `mask` | `node.set_mask` |
| `modify` | `node.modify` |
| `autolayout` | `layout.set_auto_layout` |
| `opacity` | `node.set_opacity` |
| `nofill` | `paint.remove_fill` |

**Parameter aliases**: `w`=width, `h`=height, `sz`=fontSize, `ff`=fontFamily, `bg`=color, `r`=cornerRadius, `sw`=strokeWidth, `lh`=lineHeight, `pid`=parentId, `fs`=fontStyle, `ls`=letterSpacing. Also: `fillColor`=`color` (for cross-tool compatibility).

**Auto-fix**: `lineHeight` in batch auto-gets `lineHeightUnit: PERCENT`.

## Image Data: Files, URLs, and Base64

`imageData` accepts file paths, URLs, or raw base64 — the CLI auto-detects and encodes:

```bash
# File path (absolute or ~)
ai-happy-design command shape.create_image '{"parentId":"0:1","x":0,"y":0,"width":400,"height":300,"imageData":"/tmp/hero.png"}'

# file:// URI
ai-happy-design command shape.create_image '{"parentId":"0:1","x":0,"y":0,"width":400,"height":300,"imageData":"file:///tmp/hero.png"}'

# HTTP/HTTPS URL (downloaded and encoded automatically)
ai-happy-design command shape.create_image '{"parentId":"0:1","x":0,"y":0,"width":400,"height":300,"imageData":"https://picsum.photos/800/600"}'

# In batch operations too:
{"command":"image","params":{"pid":"$root","imageData":"https://example.com/photo.jpg","w":400,"h":300}}
```

Supported: `.png`, `.jpg`, `.jpeg`, `.webp`, `.gif`, `.svg`. URLs are downloaded (max 50MB, 30s timeout). Add `--compress-images` to optimize after loading.

## Connection Diagnostics

The CLI auto-checks plugin connectivity before sending commands. If no Figma plugin is connected:
- Shows diagnostic message with the channel name and suggestions
- Polls for up to 30s waiting for the plugin to connect
- Fails with a clear error instead of a 300s silent hang

## For large batches, use a file

When the JSON payload is large (10+ operations), write it to a temp file and pass the file path:

```bash
# Write ops to file, then batch from file
ai-happy-design batch /tmp/figma-ops.json
```

This avoids shell quoting issues with large JSON payloads.

## Advanced Effects

### Glass Morphism
One-step glass effect with native Figma GLASS or simulated fallback:
```json
{"command":"glass","params":{"nodeId":"$card","intensity":"medium","tint":"#FFFFFF"}}
```
Intensities: `light`, `medium`, `heavy`. Returns `{mode: "native"}` or `{mode: "simulated"}`.

For direct native glass control:
```json
{"command":"effect.add_glass","params":{"nodeId":"$card","lightIntensity":0.5,"lightAngle":45,"refraction":0.5,"depth":1.5,"dispersion":0.1,"radius":12}}
```

### Noise Texture (Organic Feel)
Add grain/noise overlay for depth and warmth:
```json
{"command":"noise","params":{"nodeId":"$bg","noiseType":"monotone","color":"#FFFFFF","noiseSize":100,"density":0.3}}
```
Types: `monotone` (single color), `duotone` (two-color with `secondaryColor`), `multitone`.

### Layered Shadows (Realistic Depth)
Stack 2-3 shadows for realistic elevation:
```json
[
  {"command":"shadow","params":{"nodeId":"$card","offsetY":1,"radius":2,"color":"#0000000D"}},
  {"command":"shadow","params":{"nodeId":"$card","offsetY":4,"radius":12,"color":"#0000001A"}},
  {"command":"shadow","params":{"nodeId":"$card","offsetY":16,"radius":48,"color":"#00000012"}}
]
```

**Shadow recipes**: subtle (offsetY:2, radius:4, #00000010), card (offsetY:4, radius:12, #0000001A), elevated (offsetY:8, radius:24, #00000026), floating (offsetY:16, radius:48, #00000033), glow (offsetY:0, radius:24, spread:4, accent color).

### Masking
Create shaped image crops, circular avatars, gradient fade-outs:
```json
[
  {"name":"avatar_img","command":"image","params":{"pid":"$root","x":40,"y":40,"w":120,"h":120,"imageData":"...","name":"avatar-photo"}},
  {"name":"avatar_mask","command":"ellipse","params":{"pid":"$root","x":40,"y":40,"w":120,"h":120,"name":"avatar-mask"}},
  {"command":"mask","params":{"nodeId":"${{steps.avatar_mask.result.id}}","targetIds":["${{steps.avatar_img.result.id}}"]}}
]
```

### Gradient Overlays (Text Readability)
Add gradient over images so text is readable:
```json
[
  {"name":"overlay","command":"rect","params":{"pid":"$hero","w":1080,"h":400,"y":680,"name":"gradient-overlay"}},
  {"command":"gradient","params":{"nodeId":"$overlay","type":"LINEAR","stops":[{"position":0,"color":"#00000000"},{"position":1,"color":"#000000CC"}]}}
]
```

### Modify Any Node
Quick property changes without knowing the node type:
```json
{"command":"modify","params":{"nodeId":"$card","color":"#FF0000","cornerRadius":24,"opacity":0.8,"width":500}}
```
Supports: x, y, width, height, color/fillColor, opacity, cornerRadius, visible, name, rotation, characters/text, fontSize, fontFamily, fontStyle.

### Find Nodes
Search existing designs:
```json
{"command":"document.find_nodes","params":{"query":"button","type":"FRAME"}}
```

## Creative Design Techniques

### Minimaximalism (2026 Trend)
Clean, minimal layouts infused with one bold focal point. Keep structure simple but make ONE element dramatically oversized or vibrant.
- Oversized hero number/word at display tier (200px+) in accent color
- Rest of content at body/caption tier, muted gray
- Generous whitespace (padding 80-120px)
- Single accent color against monochrome palette

### Editorial / Magazine Layout
Controlled asymmetry with overlapping elements for dynamic rhythm:
- Break the grid intentionally: one element crosses boundary or overlaps
- Use decorative shapes (low-opacity circles, lines) as ambient texture
- Mix font sizes dramatically (hero + caption in same composition)
- Rotated accent text or shapes at 5-15 degrees for energy

### Grainy / Textured Aesthetic
Combine noise overlays with soft gradients for warmth and authenticity:
```json
[
  {"command":"gradient","params":{"nodeId":"$bg","type":"LINEAR","stops":[{"position":0,"color":"#1a1a2e"},{"position":1,"color":"#16213e"}]}},
  {"command":"noise","params":{"nodeId":"$bg","noiseType":"monotone","color":"#FFFFFF","density":0.15,"noiseSize":80}}
]
```

### Glass Card Composition
Frosted glass cards over gradient backgrounds:
```json
[
  {"name":"bg","command":"frame","params":{"name":"Background","w":1080,"h":1080,"bg":"#0a0a1a"}},
  {"command":"gradient","params":{"nodeId":"$bg","type":"LINEAR","stops":[{"position":0,"color":"#667eea"},{"position":1,"color":"#764ba2"}]}},
  {"name":"card","command":"frame","params":{"name":"Glass Card","pid":"$bg","x":64,"y":300,"w":952,"h":400,"r":24,"noFill":true,"layoutMode":"VERTICAL","padding":48,"itemSpacing":16}},
  {"command":"glass","params":{"nodeId":"$card","intensity":"medium"}},
  {"command":"shadow","params":{"nodeId":"$card","offsetY":8,"radius":32,"color":"#00000033"}}
]
```

### Color Palette Strategies
- **Monochrome + accent**: One hue at different lightness levels + one contrasting accent
- **Dark + neon**: Deep blacks (#0a0a0a) with electric accents (#00ff88, #7c3aed, #f43f5e)
- **Warm earth**: #1a1a1a bg with cream (#f5f0e8), terracotta (#c4704a), sage (#8fbc8f)
- **Gradient hero**: Bold gradient background (2-3 stops), white/light text on top

### Inner Light Effect (Button Bevel)
Subtle top highlight makes elements feel 3D:
```json
{"command":"shadow","params":{"nodeId":"$btn","shadowType":"INNER_SHADOW","offsetY":-1,"radius":0,"color":"#FFFFFF30"}}
```

## Balance Rules (Critical)

The #1 cause of designs looking broken is unbalanced siblings.

- **Even card heights**: Set `primaryAxisSizing: FIXED` on every card in a row. Same height.
- **Consistent padding**: All sibling cards get identical padding values
- **Consistent spacing**: All sibling cards get identical itemSpacing
- **Consistent radius**: All sibling cards get the same cornerRadius
- **Text parity**: Same-level text across siblings = same fontSize, weight, color
- **Width formula**: `cardWidth = (rowWidth - (numCards - 1) * gap) / numCards`
- **Check overlaps**: `ai-happy-design command layout.check_overlaps '{"nodeId": "..."}'`

## Typography

Font sizes MUST be proportional to canvas width. Use `design.compute_tokens` for exact values.

**Text creation** supports all properties in one call:
```bash
ai-happy-design command text.create '{
  "text": "Hello World",
  "parentId": "...",
  "fontSize": 48,
  "fontFamily": "Inter",
  "fontStyle": "Bold",
  "color": "#ffffff",
  "width": 400,
  "lineHeight": 140,
  "letterSpacing": 1,
  "textCase": "UPPER",
  "textAlign": "CENTER"
}'
```

**List available fonts**: `ai-happy-design command text.list_fonts` or `ai-happy-design command text.list_fonts '{"family": "Inter"}'`

**Key rules**:
- Always set `width` on text nodes (prevents "hug contents" overflow)
- `lineHeight` is auto-set to PERCENT unit in batch mode
- `color` accepts hex string ("#ff0000") or object ({r:1,g:0,b:0})
- `fontStyle` is a string like "Regular", "Bold", "Medium" — NOT a weight number

## Visual Hierarchy

- **Primary**: One focal element. Largest size (title/hero tier), boldest, accent color
- **Secondary**: Supporting. Medium size (heading/subheading tier), semi-bold, white
- **Tertiary**: Details. Smaller (body tier), regular weight, muted gray
- **Ambient**: Background decoration. Very low opacity (0.05-0.15), never competes with content

## Design Decisions Quick Reference

| Element | Treatment |
|---------|-----------|
| Cards, CTAs | Layered shadows (2-3 stacked) for realistic depth |
| Hero backgrounds | Linear gradient + noise texture overlay |
| Glass effects | `glass` command (native or simulated auto-detected) |
| Card borders | stroke: "#ffffff12", strokeWidth: 1 |
| Structural frames | noFill: true (removes default white) |
| Image backgrounds | `image` command + gradient overlay for text readability |
| Circular avatars | `image` + `ellipse` + `mask` |
| Buttons | Frame + auto-layout + inner shadow bevel + layered shadows |

## Page Management

Pages are separate from document operations. Use `page.*` commands — NOT `document.*`.

```bash
ai-happy-design command page.get_all '{}'                                        # list all pages
ai-happy-design command page.create '{"name": "Slide 2"}'                       # new page
ai-happy-design command page.set_current '{"pageId": "3:0"}'                    # switch to page
ai-happy-design command page.rename '{"pageId": "3:0", "name": "Final"}'        # rename page
ai-happy-design command page.duplicate '{"pageId": "3:0"}'                      # duplicate page
ai-happy-design command page.delete '{"pageId": "3:0"}'                         # delete page
```

**Common mistake**: `document.list_pages` → correct command is `page.get_all`

## Relay & Port Configuration

The relay auto-starts on port 3055 and auto-shuts down after 15 minutes of inactivity. CLI auto-restarts it on next command.

**Custom port**: `ai-happy-design --port 4000 ws` (or `PORT=4000`). The CLI will guide you to update the Figma plugin relay URL.

**Check status**: `curl http://localhost:3055/status` — shows channels, uptime, idle timeout remaining.

**Env vars**: `PORT` (relay port), `AHD_IDLE_TIMEOUT` (e.g. "30m", "0" to disable), `AHD_CHANNEL` (default channel).

## Reference Files

- **[CSS-to-Figma map](references/css-to-figma.md)** — Complete property translation table
- **[Design patterns](references/design-patterns.md)** — Layout structures, card building, color systems
- **[Batch examples](references/batch-examples.md)** — Ready-to-use batch payloads for common designs
- **[Advanced effects](references/advanced-effects.md)** — Glass, noise, shadows, masking, creative techniques

## Key CLI Commands

```bash
ai-happy-design command design.compute_tokens '{"width": W, "height": H}'
ai-happy-design command document.find_free_space '{"width": W, "height": H}'
ai-happy-design command document.find_nodes '{"query": "...", "type": "FRAME"}'
ai-happy-design command text.list_fonts
ai-happy-design command text.list_fonts '{"family": "Inter"}'
ai-happy-design command export.image '{"nodeId": "...", "scale": 2}' -o output.png
ai-happy-design command layout.check_overlaps '{"nodeId": "..."}'
ai-happy-design command design_system.analyze '{}'                                    # scan file for existing styles/components/variables
ai-happy-design command node.get_tree '{"nodeId": "...", "compact": true}'            # flat tree summary (3-5x fewer tokens)
ai-happy-design command document.get_info '{}'                                        # document metadata
ai-happy-design command document.get_selection '{}'                                   # currently selected nodes
ai-happy-design command document.scan_by_type '{"nodeType": "FRAME"}'                 # scan all nodes of a type
ai-happy-design batch '<json-array-of-operations>'
ai-happy-design batch operations.json
ai-happy-design batch f1.json f2.json f3.json        # multi-file batch
ai-happy-design batch ./slides/                       # directory batch
ai-happy-design batch f1.json f2.json --parallel      # concurrent execution
ai-happy-design --port 4000 command ...               # custom port
ai-happy-design actions                               # list all available actions
ai-happy-design tools                                 # full tool catalog
```

## Common Pitfalls

- Never create cards as rectangles + floating text. Cards are frames with auto-layout.
- Never leave default names ("Frame 47"). Name everything descriptively.
- Never use text below minimum size for the canvas (use compute_tokens caption tier as minimum).
- Never skip the balance check on sibling elements.
- Always use 8px grid alignment for all dimensions.
- Always call `document.find_free_space` before placing frames — never guess coordinates (batch does this automatically).
- Always set `lineHeightUnit: PERCENT` when setting lineHeight manually (auto in batch).
- Structural frames (wrappers, groups) need `noFill: true` — Figma adds white fill by default.
- For large batch payloads, write JSON to a temp file to avoid shell quoting issues.
- Use `color` (not `fillColor`) for frame/text colors — `fillColor` is auto-aliased but `color` is canonical.
- Step names in batch must be snake_case (no spaces) — auto-sanitized but best to write correctly.
- **Wrong domain for pages**: Use `page.get_all` not `document.list_pages` — all page operations use `page.*` not `document.*`
- **Wrong domain for layer ops**: Use `layer.group` / `layer.bring_to_front` / `layer.send_to_back` — these are NOT under `node.*`
- **Design system scan**: `design_system.analyze` not `document.analyze` or `document.get_design_system`
- **Compact tree for discovery**: `node.get_tree {"nodeId":"...", "compact":true}` gives a flat summary array — 3-5x fewer tokens than default
- **`type` instead of `command`**: The batch format uses `"command"`, not `"type"`. Always `{"command": "node.create_frame", ...}`.
- **Top-level params**: ALL design properties go inside `"params"`: `{"command":"...","params":{"x":0,"color":"#fff"}}` — never `{"command":"...","x":0,"color":"#fff"}`.
- **Step name mismatch**: If a step is named `"bg"`, reference it as `${{steps.bg.result.id}}` — not `background`, `Background`, `root`. Mismatch fails silently.
- **Over-length batches**: Keep ≤15 ops per batch file. Schema compliance degrades at 30+ ops in production. Split large designs.
