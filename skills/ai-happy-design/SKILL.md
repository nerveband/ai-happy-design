---
name: ai-happy-design
description: Design in Figma using the ai-happy-design MCP server. Use when asked to create, edit, or export Figma designs — social media posts, cards, layouts, UI mockups, posters, or any visual design. Triggers on "design in Figma", "create a Figma layout", "make an Instagram post", "export from Figma", "Figma design", or any task involving the ai-happy-design CLI/MCP tools. Also use when asked to batch-create multi-element Figma compositions.
---

# AI Happy Design

Design in Figma via the `ai-happy-design` MCP server. One Go binary provides MCP tools, CLI commands, and a WebSocket relay to a Figma plugin.

## Critical: Think in HTML/CSS First

You already know HTML/CSS. Use that knowledge. Mentally draft every design as HTML/CSS first, then translate to Figma commands.

```
CSS property              → Figma command
display: flex             → node.create_frame {layoutMode: VERTICAL, ...}
gap: 16px                 → node.create_frame {itemSpacing: 16}
padding: 24px             → node.create_frame {padding: 24}
background-color          → node.create_frame {color: "#1a1a1a"}
border-radius: 16px       → node.create_frame {cornerRadius: 16}
border: 1px solid         → node.create_frame {stroke: "#333", strokeWidth: 1}
box-shadow                → effect.add_shadow
```

Full CSS-to-Figma mapping: See [references/css-to-figma.md](references/css-to-figma.md)

## Workflow

1. **Get design tokens**: Call `design.compute_tokens {width, height}` first — get proportional font sizes, spacing, and layout classification
2. **Find free space**: Call `document.find_free_space {width, height}` — get exact x,y for placement
3. **Think in CSS**: Picture the design as a webpage. What HTML elements? What CSS?
4. **Create in one step**: Use `node.create_frame` with all params (auto-layout, color, corner radius, stroke) — no need for separate commands
5. **Batch create**: Build a JSON array of operations. Send as one payload via `bulk.execute`
6. **Balance check**: All siblings MUST match — same height, padding, spacing, radius, text sizes
7. **Export & verify**: `export.image {nodeId, scale: 2}` — always 2x for crisp output

## One-Step Frame Creation (v0.6.0)

Create fully-configured frames in a single command — no need for separate layout/paint calls:

```json
{
  "command": "node.create_frame",
  "params": {
    "name": "Card",
    "x": 0, "y": 0, "width": 400, "height": 300,
    "parentId": "${{steps.root.result.id}}",
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
  }
}
```

For structural/wrapper frames, use `"noFill": true` to remove the default white fill.

## Design Tokens (Proportional Sizing)

Always call `design.compute_tokens` first to get sizes proportional to the canvas:

```json
{"command": "design.compute_tokens", "params": {"width": 1080, "height": 1080}}
```

Returns a modular type scale (perfect fourth = 1.333 ratio):
- **display** (200px at 1080w) — statement text
- **hero** (152px) — hero headlines
- **title** (112px) — section titles
- **heading** (84px) — card headings
- **subheading** (64px) — subtitles
- **body** (48px) — paragraph text
- **caption** (36px) — labels, metadata

Also returns spacing tokens (xs through xxxl) on 8px grid, and layout classification (portrait/square/landscape).

For print: add `"dpi": 300` to get point-based sizes.

## Batch Aliases (v0.6.0)

Use compact command names for faster batching:

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
| `autolayout` | `layout.set_auto_layout` |
| `opacity` | `node.set_opacity` |
| `nofill` | `paint.remove_fill` |

Parameter aliases also work: `w`=width, `h`=height, `sz`=fontSize, `ff`=fontFamily, `bg`=color, `r`=cornerRadius, `sw`=strokeWidth, `lh`=lineHeight.

**Auto-fix**: When `lineHeight` is set in batch, `lineHeightUnit: PERCENT` is auto-injected (Figma defaults to PIXELS otherwise).

## Creating vs Editing

| Task | Method |
|------|--------|
| New design (3+ elements) | Batch: `bulk.execute` with named steps and `${{steps.name.result.id}}` interpolation |
| Single change | Direct command: `paint.set_solid`, `node.resize`, etc. |

## Balance Rules (Critical)

The #1 cause of designs looking broken is unbalanced siblings.

- **Even card heights**: Set `primaryAxisSizing: FIXED` on every card in a row. Same height.
- **Consistent padding**: All sibling cards get identical padding values
- **Consistent spacing**: All sibling cards get identical itemSpacing
- **Consistent radius**: All sibling cards get the same cornerRadius
- **Text parity**: Same-level text across siblings = same fontSize, weight, color
- **Width formula**: `cardWidth = (rowWidth - (numCards - 1) * gap) / numCards`
- **Check overlaps**: Use `layout.check_overlaps {nodeId}` to verify children don't overlap

## Typography

Font sizes MUST be proportional to canvas width. Use `design.compute_tokens` for exact values.

**Text creation** supports all properties in one call:
```json
{
  "command": "text.create",
  "params": {
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
  }
}
```

**List available fonts**: `text.list_fonts` or `text.list_fonts {family: "Inter"}` to filter.

**Key rules**:
- Always set `width` on text nodes (prevents "hug contents" overflow)
- `lineHeight` is auto-set to PERCENT unit in batch mode
- `color` accepts hex string ("#ff0000") or object ({r:1,g:0,b:0})
- `fontStyle` is a string like "Regular", "Bold", "Medium" — NOT a weight number

## Visual Hierarchy

Every element has a level. Size, weight, and color reflect it:

- **Primary**: One focal element. Largest size (title/hero tier), boldest, accent color
- **Secondary**: Supporting. Medium size (heading/subheading tier), semi-bold, white
- **Tertiary**: Details. Smaller (body tier), regular weight, muted gray
- **Ambient**: Background decoration. Very low opacity (0.05-0.15), never competes with content

## Design Decisions Quick Reference

| Element | Treatment |
|---------|-----------|
| Cards, CTAs | Drop shadow (offsetY:4, radius:16, alpha:0.2) |
| Hero backgrounds | Linear gradient (dark → slightly lighter) |
| Glass effects | Background blur (16px) + semi-transparent fill + thin stroke |
| Card borders | stroke: "#ffffff12", strokeWidth: 1 |
| Structural frames | noFill: true (removes default white) |

## Reference Files

- **[CSS-to-Figma map](references/css-to-figma.md)** — Complete property translation table
- **[Design patterns](references/design-patterns.md)** — Layout structures, card building, color systems
- **[Batch examples](references/batch-examples.md)** — Ready-to-use batch payloads for common designs

## Key Commands

```
design.compute_tokens {width, height}  — Proportional sizing for any canvas
document.find_free_space {width, height} — Find placement coordinates
describe(action="catalog")              — Full tool catalog with examples
describe(action="design_guide")         — Design best practices and patterns
bulk.execute                            — Batch operations with interpolation
export.image {nodeId, scale: 2}         — Export at 2x resolution
text.list_fonts                         — Available fonts with filtering
layout.check_overlaps {nodeId}          — Detect overlapping children
```

## Common Pitfalls

- Never create cards as rectangles + floating text. Cards are frames with auto-layout.
- Never leave default names ("Frame 47"). Name everything descriptively.
- Never use text below minimum size for the canvas (use compute_tokens caption tier as minimum).
- Never skip the balance check on sibling elements.
- Always use 8px grid alignment for all dimensions.
- Always call `document.find_free_space` before placing frames — never guess coordinates.
- Always set `lineHeightUnit: PERCENT` when setting lineHeight manually (auto in batch).
- Structural frames (wrappers, groups) need `noFill: true` — Figma adds white fill by default.
