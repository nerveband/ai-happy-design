# Design Patterns

## Table of Contents
- [Design Tokens](#design-tokens)
- [Coordinate System](#coordinate-system)
- [8px Grid](#8px-grid)
- [Frames Not Shapes](#frames-not-shapes)
- [Auto-Layout](#auto-layout)
- [Building Cards](#building-cards)
- [Layout Structures](#layout-structures)
- [Balance Checklist](#balance-checklist)
- [Colors](#colors)
- [Shadows](#shadows)
- [Layer Organization](#layer-organization)
- [Composition Tips](#composition-tips)

## Design Tokens

Always start by calling `design.compute_tokens {width, height}` to get proportional sizing.

**Type scale** (perfect fourth = 1.333 ratio, base = width * 0.044):
- display → hero → title → heading → subheading → body → caption
- At 1080px: 200 → 152 → 112 → 84 → 64 → 48 → 36

**Spacing** (8px grid multiples):
- xs(8) → sm(16) → md(24) → lg(32) → xl(48) → xxl(64) → xxxl(96)

**Tier selection by content**:
- Short text (1-3 words): title/hero tier
- Medium text (1-2 sentences): heading/subheading tier
- Long paragraphs: body tier
- Labels/metadata: caption tier

**Button tokens**: height = 2.5 × fontSize, horizontalPadding = 1.25 × fontSize, pill radius = height/2

## Coordinate System

- Origin: top-left (0,0). X increases right, Y increases down.
- x/y are RELATIVE to the containing parent frame.
- Inside auto-layout frames, x/y are IGNORED for auto-positioned children.
- Always call `document.find_free_space {width, height}` before placing root frames.

## 8px Grid

Align ALL dimensions to 8px multiples: x, y, width, height, padding, spacing.

Common values: 8, 16, 24, 32, 48, 64, 96, 128. Use 4 only for text sizes (on 4px grid).

## Frames Not Shapes

Use `node.create_frame` as containers, NOT rectangles. Frames support auto-layout, clipping, nesting.

A "card" is a frame with auto-layout and children. Never a rectangle with floating text.

## Auto-Layout

Use for ALL structural containers: cards, rows, columns, stacks, headers, footers. Set on creation:

```json
{"command": "frame", "params": {
  "name": "Card", "w": 400, "h": 300,
  "bg": "#1a1a1a", "r": 16,
  "layoutMode": "VERTICAL", "itemSpacing": 16, "padding": 24,
  "primaryAxisAlign": "MIN", "counterAxisAlign": "CENTER",
  "primaryAxisSizing": "FIXED", "counterAxisSizing": "FIXED"
}}
```

**Sizing modes:**
- HUG (AUTO): Frame shrinks to wrap children. Use for buttons, badges.
- FILL: Child stretches to fill parent. Use for main content areas.
- FIXED: Maintains exact dimensions. Use for root frames, cards in rows.

**Child sizing** (set on children, not parent):
- `layoutSizingHorizontal: FILL` — child stretches horizontally
- `layoutSizingVertical: HUG` — child hugs vertically

**Axis guide:**
- VERTICAL frame: primaryAxisAlign = vertical, counterAxisAlign = horizontal
- HORIZONTAL frame: primaryAxisAlign = horizontal, counterAxisAlign = vertical

## Building Cards

One-step in v0.6.0 — a card is a single frame creation + text children:

```json
[
  {"name":"card","command":"frame","params":{"name":"Feature Card","w":320,"h":160,"bg":"#1a1a1a","r":16,"layoutMode":"VERTICAL","itemSpacing":8,"padding":24}},
  {"command":"text","params":{"text":"Title","parentId":"${{steps.card.result.id}}","sz":18,"fontStyle":"Bold","color":"#ffffff","w":280}},
  {"command":"text","params":{"text":"Description","parentId":"${{steps.card.result.id}}","sz":14,"color":"#999999","w":280}}
]
```

## Layout Structures

### Centered Content Page
Root frame (FIXED 1080x1080, bg color) → Content wrapper (VERTICAL, CENTER/CENTER, padding 64, noFill) → children stack and center.

### Two-Column Grid
Parent = HORIZONTAL auto-layout (noFill), children = card frames with same width and `primaryAxisSizing: FIXED`.

### Header/Content/Footer
Root (VERTICAL, FIXED) → Header (FIXED height) → Content (fills) → Footer (FIXED height).

## Balance Checklist

Run this for EVERY set of sibling elements:

- [ ] All cards: same width, same height (`primaryAxisSizing: FIXED`)
- [ ] All cards: same padding, same itemSpacing, same cornerRadius
- [ ] All text at same hierarchy level: same fontSize, fontStyle, color
- [ ] Row gap is from design system (8, 16, 24, 32)
- [ ] Total: `rowWidth = N * cardWidth + (N-1) * gap`
- [ ] Run `layout.check_overlaps {nodeId}` to verify no overlapping children

## Colors

Colors accept hex strings ("#1a1a1a") or `{r, g, b, a}` objects (0.0-1.0 floats).

**Dark theme palette:**
- Background: `"#1a1a1a"`
- Surface: `"#ffffff0a"` (white at 4%)
- Border: `"#ffffff12"` (white at 7%)
- Text primary: `"#ffffff"`
- Text muted: `"#999999"`

Glass/card effect: low-opacity white fill + 1px low-opacity white stroke + background blur.

## Shadows

- Subtle: `offsetY:2, radius:8, alpha:0.15`
- Medium: `offsetY:4, radius:16, alpha:0.25`
- Large: `offsetY:8, radius:32, alpha:0.3`
- Accent glow: accent color, `alpha:0.15, offsetY:4, radius:24`

## Layer Organization

**Naming**: Name everything by role. Never "Frame 47". Pattern: `[Type] - [Purpose]`.

**Grouping**:
```
Root Frame
  ├── Background (decorative shapes, gradients) — BACK
  ├── Content Frame (auto-layout, structured content) — MIDDLE
  │   ├── Header
  │   ├── Hero section
  │   ├── Body / Cards
  │   └── Footer
  └── Overlay elements (floating badges, stripes) — FRONT
```

**Layer order** = z-index. Back renders first (behind). Front renders last (on top).

## Composition Tips

1. Call `document.find_free_space` before placing root frames — never guess coordinates.
2. Call `design.compute_tokens` for proportional font/spacing sizes.
3. Create root frame first (1080x1080 for Instagram). Everything inside it.
4. Auto-layout for STRUCTURED content. Absolute x/y only for DECORATIVE elements.
5. Build cards as frames with auto-layout, not rectangles + floating text.
6. Use `noFill: true` on structural/wrapper frames — Figma adds white fill by default.
7. Decorative elements in root frame with computed x/y, sent to back.
8. Name all elements descriptively.
