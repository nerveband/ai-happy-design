# Batch Payload Examples

## How Batch Works

Send a JSON array of operations via the CLI. Each operation can have a `name` for reference by later steps via `${{steps.name.result.id}}`.

**Inline:**
```bash
ahd-figma batch '[{"name":"root","command":"frame","params":{...}}, ...]'
```

**From file (recommended for 10+ operations):**
```bash
ahd-figma batch ops.json
```

**With options:**
```bash
ahd-figma batch --fail-fast --live ops.json
```

## Example: Single Card

3 ops — one-step frame creation + text children:

```json
[
  {"name":"card","command":"frame","params":{"name":"Feature Card","x":0,"y":0,"w":320,"h":160,"bg":"#1a1a1a","r":16,"layoutMode":"VERTICAL","itemSpacing":8,"padding":24}},
  {"name":"title","command":"text","params":{"text":"Feature Title","pid":"${{steps.card.result.id}}","sz":18,"fontStyle":"Bold","color":"#ffffff","w":280}},
  {"command":"text","params":{"text":"Card description goes here","pid":"${{steps.card.result.id}}","sz":14,"color":"#999999","w":280}}
]
```

Compact aliases: `frame` = `node.create_frame`, `text` = `text.create`, `w` = width, `h` = height, `sz` = fontSize, `bg` = color, `r` = cornerRadius, `pid` = parentId.

## Example: Centered Content Page

```json
[
  {"name":"root","command":"frame","params":{"name":"Post","x":0,"y":0,"w":1080,"h":1080,"bg":"#1a1a1a"}},
  {"name":"content","command":"frame","params":{"name":"Content","pid":"${{steps.root.result.id}}","w":1080,"h":1080,"noFill":true,"layoutMode":"VERTICAL","itemSpacing":24,"padding":64,"primaryAxisAlign":"CENTER","counterAxisAlign":"CENTER"}},
  {"command":"text","params":{"text":"Hello World","pid":"${{steps.content.result.id}}","sz":80,"fontStyle":"Bold","color":"#ffffff","textAlign":"CENTER","w":952}}
]
```

## Example: Two-Column Card Row

```json
[
  {"name":"row","command":"frame","params":{"name":"Card Row","w":680,"h":160,"noFill":true,"layoutMode":"HORIZONTAL","itemSpacing":16}},
  {"name":"col1","command":"frame","params":{"name":"Card 1","pid":"${{steps.row.result.id}}","w":332,"h":160,"bg":"#222222","r":12,"layoutMode":"VERTICAL","padding":16,"primaryAxisSizing":"FIXED"}},
  {"name":"col2","command":"frame","params":{"name":"Card 2","pid":"${{steps.row.result.id}}","w":332,"h":160,"bg":"#222222","r":12,"layoutMode":"VERTICAL","padding":16,"primaryAxisSizing":"FIXED"}}
]
```

## Example: Card with Border and Shadow

```json
[
  {"name":"card","command":"frame","params":{"name":"Glass Card","w":400,"h":200,"bg":"#ffffff0a","r":16,"stroke":"#ffffff12","strokeWidth":1,"layoutMode":"VERTICAL","padding":24,"itemSpacing":12}},
  {"command":"shadow","params":{"nodeId":"${{steps.card.result.id}}","shadowType":"DROP_SHADOW","offsetY":4,"radius":16,"color":{"r":0,"g":0,"b":0,"a":0.25}}},
  {"command":"blur","params":{"nodeId":"${{steps.card.result.id}}","blurType":"BACKGROUND_BLUR","radius":16}}
]
```

## Key Patterns

- Name the root frame step so children can reference its ID
- Use `noFill: true` for layout-only wrapper frames
- Set auto-layout ON creation with `layoutMode` param
- Create text INSIDE frames using `pid` (parentId)
- Set text properties (fontStyle, color, fontSize) ON creation
- For FIXED-height cards in rows: set `primaryAxisSizing: FIXED` on each card
- Use `--fail-fast` flag when debugging
- Always set `width` on text nodes to enable wrapping
- `lineHeight` auto-gets `lineHeightUnit: PERCENT` in batch mode
- For large payloads, write JSON to a file and pass the file path

## Command Alias Quick Reference

| Alias | Full | Alias | Full |
|-------|------|-------|------|
| `frame` | `node.create_frame` | `fill` | `paint.set_solid` |
| `rect` | `shape.create_rectangle` | `stroke` | `paint.set_stroke` |
| `ellipse` | `shape.create_ellipse` | `gradient` | `paint.set_gradient` |
| `text` | `text.create` | `shadow` | `effect.add_shadow` |
| `image` | `shape.create_image` | `blur` | `effect.add_blur` |
| `autolayout` | `layout.set_auto_layout` | `nofill` | `paint.remove_fill` |
| `opacity` | `node.set_opacity` | `parent` | `node.set_parent` |

## Parameter Alias Quick Reference

| Short | Full | Short | Full |
|-------|------|-------|------|
| `w` | `width` | `sz` | `fontSize` |
| `h` | `height` | `ff` | `fontFamily` |
| `bg` | `color` | `fs` | `fontStyle` |
| `r` | `cornerRadius` | `lh` | `lineHeight` |
| `sw` | `strokeWidth` | `ls` | `letterSpacing` |
| `pid` | `parentId` | | |
