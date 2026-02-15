# CSS-to-Figma Command Translation

Translate CSS properties you know into Figma commands. Most properties can now be set on `node.create_frame` or `text.create` directly.

## Layout (prefer one-step on create_frame)

| CSS | Figma Command |
|-----|---------------|
| `display: flex; flex-direction: column` | `node.create_frame {layoutMode: VERTICAL, ...}` |
| `display: flex; flex-direction: row` | `node.create_frame {layoutMode: HORIZONTAL, ...}` |
| `justify-content: center` | `node.create_frame {primaryAxisAlign: CENTER}` |
| `justify-content: space-between` | `node.create_frame {primaryAxisAlign: SPACE_BETWEEN}` |
| `align-items: center` | `node.create_frame {counterAxisAlign: CENTER}` |
| `gap: 16px` | `node.create_frame {itemSpacing: 16}` |
| `padding: 24px` | `node.create_frame {padding: 24}` |
| `padding: 16px 24px` | `node.create_frame {paddingTop:16, paddingBottom:16, paddingLeft:24, paddingRight:24}` |
| `flex-wrap: wrap` | `node.create_frame {layoutWrap: WRAP}` |
| `width: 100%; height: auto` | `layout.set_sizing {layoutSizingHorizontal: FILL, layoutSizingVertical: HUG}` |
| `position: absolute` | Set x/y on child (outside auto-layout flow) |
| `overflow: hidden` | `node.create_frame {clipsContent: true}` |

For existing frames, use `layout.set_auto_layout`, `layout.set_padding`, `layout.set_alignment`, `layout.set_sizing`.

## Colors & Fills

| CSS | Figma Command |
|-----|---------------|
| `background-color: #1a1a1a` | `node.create_frame {color: "#1a1a1a"}` or `paint.set_solid {color}` |
| `background: none` | `node.create_frame {noFill: true}` or `paint.remove_fill {index: 0}` |
| `background: linear-gradient(...)` | `paint.set_gradient {type:LINEAR, stops:[{position:0,color},{position:1,color}]}` |
| `background: radial-gradient(...)` | `paint.set_gradient {type:RADIAL, stops:[...]}` |
| `opacity: 0.5` | `node.set_opacity {opacity: 0.5}` |
| `mix-blend-mode: multiply` | `node.set_blend_mode {blendMode: MULTIPLY}` |

## Borders & Shapes

| CSS | Figma Command |
|-----|---------------|
| `border: 1px solid #333` | `node.create_frame {stroke: "#333", strokeWidth: 1}` or `paint.set_stroke` |
| `border-radius: 16px` | `node.create_frame {cornerRadius: 16}` or `node.set_corner_radius` |

## Effects

| CSS | Figma Command |
|-----|---------------|
| `box-shadow: 0 4px 16px rgba(...)` | `effect.add_shadow {shadowType:DROP_SHADOW, offsetX:0, offsetY:4, radius:16, color}` |
| `box-shadow: inset ...` | `effect.add_shadow {shadowType:INNER_SHADOW, ...}` |
| `filter: blur(8px)` | `effect.add_blur {blurType:LAYER_BLUR, radius:8}` |
| `backdrop-filter: blur(16px)` | `effect.add_blur {blurType:BACKGROUND_BLUR, radius:16}` |

## Typography (prefer one-step on text.create)

| CSS | Figma Command |
|-----|---------------|
| `color: white` | `text.create {color: "#ffffff"}` or `text.set_color` |
| `font-weight: bold` | `text.create {fontStyle: "Bold"}` or `text.set_weight` |
| `font-size: 24px` | `text.create {fontSize: 24}` or `text.set_size` |
| `font-family: Inter` | `text.create {fontFamily: "Inter"}` or `text.set_font` |
| `text-transform: uppercase` | `text.create {textCase: "UPPER"}` or `text.set_case` |
| `letter-spacing: 2px` | `text.create {letterSpacing: 2}` or `text.set_letter_spacing` |
| `line-height: 140%` | `text.create {lineHeight: 140}` (auto PERCENT in batch) |
| `text-align: center` | `text.create {textAlign: "CENTER"}` or `text.set_alignment` |
| `text-decoration: underline` | `text.set_decoration {decoration: UNDERLINE}` |

## Layering

| CSS | Figma Command |
|-----|---------------|
| `z-index` (higher) | `layer.bring_to_front` |
| `z-index` (lower) | `layer.send_to_back` |
