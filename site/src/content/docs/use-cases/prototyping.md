---
title: UI Prototyping
description: Build wireframes, UI prototypes, and inspect existing designs with AI Happy Design.
---

AI Happy Design is well-suited for rapid UI prototyping. Instead of dragging rectangles around in Figma, describe your layout as structured data and let the tool build it in milliseconds.

## Wireframe Components

Every wireframe is built from the same primitives: frames for containers, rectangles for placeholders, text for labels, and auto-layout for spacing.

### Mobile App Screen

Create a mobile app login screen in one batch:

```json
[
  {
    "name": "screen",
    "command": "frame",
    "params": {
      "name": "Login Screen",
      "w": 375, "h": 812,
      "bg": "#FFFFFF",
      "layoutMode": "VERTICAL",
      "paddingTop": 120,
      "paddingRight": 24,
      "paddingBottom": 48,
      "paddingLeft": 24,
      "itemSpacing": 24,
      "primaryAxisAlign": "MIN",
      "counterAxisAlign": "CENTER",
      "clipsContent": true
    }
  },
  {
    "name": "title",
    "command": "text",
    "params": {
      "pid": "${{steps.screen.result.id}}",
      "text": "Welcome back",
      "sz": 32,
      "ff": "Inter",
      "fontStyle": "Bold",
      "color": "#1a1a1a",
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "subtitle",
    "command": "text",
    "params": {
      "pid": "${{steps.screen.result.id}}",
      "text": "Sign in to continue",
      "sz": 16,
      "ff": "Inter",
      "color": "#888888",
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "email_field",
    "command": "frame",
    "params": {
      "pid": "${{steps.screen.result.id}}",
      "name": "Email Field",
      "layoutMode": "HORIZONTAL",
      "padding": 16,
      "bg": "#F5F5F5",
      "r": 12,
      "layoutSizingHorizontal": "FILL",
      "h": 52
    }
  },
  {
    "name": "email_placeholder",
    "command": "text",
    "params": {
      "pid": "${{steps.email_field.result.id}}",
      "text": "Email address",
      "sz": 16,
      "ff": "Inter",
      "color": "#888888"
    }
  },
  {
    "name": "password_field",
    "command": "frame",
    "params": {
      "pid": "${{steps.screen.result.id}}",
      "name": "Password Field",
      "layoutMode": "HORIZONTAL",
      "padding": 16,
      "bg": "#F5F5F5",
      "r": 12,
      "layoutSizingHorizontal": "FILL",
      "h": 52
    }
  },
  {
    "name": "password_placeholder",
    "command": "text",
    "params": {
      "pid": "${{steps.password_field.result.id}}",
      "text": "Password",
      "sz": 16,
      "ff": "Inter",
      "color": "#888888"
    }
  },
  {
    "name": "sign_in_btn",
    "command": "frame",
    "params": {
      "pid": "${{steps.screen.result.id}}",
      "name": "Sign In Button",
      "layoutMode": "HORIZONTAL",
      "padding": 16,
      "bg": "#0a0a0a",
      "r": 12,
      "layoutSizingHorizontal": "FILL",
      "h": 52,
      "primaryAxisAlign": "CENTER",
      "counterAxisAlign": "CENTER"
    }
  },
  {
    "name": "btn_text",
    "command": "text",
    "params": {
      "pid": "${{steps.sign_in_btn.result.id}}",
      "text": "Sign In",
      "sz": 16,
      "ff": "Inter",
      "fontStyle": "SemiBold",
      "color": "#FFFFFF"
    }
  }
]
```

### Wireframe Card

```json
[
  {
    "name": "card",
    "command": "frame",
    "params": {
      "name": "Card Wireframe",
      "w": 360, "h": 240,
      "bg": "#FFFFFF",
      "r": 12,
      "layoutMode": "VERTICAL",
      "padding": 16,
      "itemSpacing": 12
    }
  },
  {
    "name": "image_placeholder",
    "command": "rect",
    "params": {
      "pid": "${{steps.card.result.id}}",
      "name": "Image Placeholder",
      "h": 120,
      "bg": "#E0E0E0",
      "r": 8,
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "card_title",
    "command": "text",
    "params": {
      "pid": "${{steps.card.result.id}}",
      "text": "Card Title",
      "sz": 18,
      "ff": "Inter",
      "fontStyle": "Bold",
      "color": "#1a1a1a",
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "card_desc",
    "command": "text",
    "params": {
      "pid": "${{steps.card.result.id}}",
      "text": "Short description text goes here. Two lines max.",
      "sz": 14,
      "ff": "Inter",
      "color": "#666666",
      "layoutSizingHorizontal": "FILL",
      "lh": 150,
      "lineHeightUnit": "PERCENT"
    }
  }
]
```

## Auto-Layout for Responsive Prototypes

Auto-layout is the backbone of responsive prototyping. It maps directly to CSS flexbox concepts:

| Figma Property | CSS Equivalent | Values |
|---------------|---------------|--------|
| `layoutMode` | `flex-direction` | `HORIZONTAL`, `VERTICAL` |
| `itemSpacing` | `gap` | Pixel value |
| `padding` | `padding` | Pixel value (uniform) |
| `primaryAxisAlign` | `justify-content` | `MIN`, `CENTER`, `MAX`, `SPACE_BETWEEN` |
| `counterAxisAlign` | `align-items` | `MIN`, `CENTER`, `MAX` |
| `layoutSizingHorizontal` | `width` behavior | `FIXED`, `HUG`, `FILL` |
| `layoutSizingVertical` | `height` behavior | `FIXED`, `HUG`, `FILL` |

### Dashboard Wireframe with Auto-Layout

```json
[
  {
    "name": "wireframe",
    "command": "frame",
    "params": {
      "name": "Wireframe - Dashboard",
      "w": 1440, "h": 900,
      "bg": "#FFFFFF",
      "layoutMode": "HORIZONTAL",
      "itemSpacing": 0,
      "clipsContent": true
    }
  },
  {
    "name": "sidebar",
    "command": "frame",
    "params": {
      "pid": "${{steps.wireframe.result.id}}",
      "name": "Sidebar",
      "w": 240,
      "bg": "#F5F5F5",
      "layoutMode": "VERTICAL",
      "padding": 16,
      "itemSpacing": 8,
      "layoutSizingVertical": "FILL"
    }
  },
  {
    "name": "nav_item_1",
    "command": "frame",
    "params": {
      "pid": "${{steps.sidebar.result.id}}",
      "name": "Nav - Dashboard",
      "layoutMode": "HORIZONTAL",
      "padding": 12,
      "bg": "#E0E0E0",
      "r": 8,
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "nav_text_1",
    "command": "text",
    "params": {
      "pid": "${{steps.nav_item_1.result.id}}",
      "text": "Dashboard",
      "sz": 14,
      "ff": "Inter",
      "fontStyle": "Medium",
      "color": "#1a1a1a"
    }
  },
  {
    "name": "content",
    "command": "frame",
    "params": {
      "pid": "${{steps.wireframe.result.id}}",
      "name": "Content Area",
      "layoutMode": "VERTICAL",
      "padding": 32,
      "itemSpacing": 24,
      "noFill": true,
      "layoutSizingHorizontal": "FILL",
      "layoutSizingVertical": "FILL"
    }
  },
  {
    "name": "content_header",
    "command": "text",
    "params": {
      "pid": "${{steps.content.result.id}}",
      "text": "Dashboard Overview",
      "sz": 24,
      "ff": "Inter",
      "fontStyle": "Bold",
      "color": "#1a1a1a"
    }
  },
  {
    "name": "stat_row",
    "command": "frame",
    "params": {
      "pid": "${{steps.content.result.id}}",
      "name": "Stat Cards",
      "layoutMode": "HORIZONTAL",
      "itemSpacing": 16,
      "noFill": true,
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "stat_card_1",
    "command": "frame",
    "params": {
      "pid": "${{steps.stat_row.result.id}}",
      "name": "Revenue Card",
      "layoutMode": "VERTICAL",
      "padding": 16,
      "itemSpacing": 4,
      "bg": "#F9FAFB",
      "r": 12,
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "stat_label",
    "command": "text",
    "params": {
      "pid": "${{steps.stat_card_1.result.id}}",
      "text": "Monthly Revenue",
      "sz": 12,
      "ff": "Inter",
      "color": "#6B7280"
    }
  },
  {
    "name": "stat_value",
    "command": "text",
    "params": {
      "pid": "${{steps.stat_card_1.result.id}}",
      "text": "$12,450",
      "sz": 28,
      "ff": "Inter",
      "fontStyle": "Bold",
      "color": "#111827"
    }
  }
]
```

Key auto-layout patterns used:

- **`layoutSizingHorizontal: "FILL"`** on children to stretch and fill available width
- **`layoutSizingVertical: "FILL"`** on the content area to take remaining height
- **`noFill: true`** on structural wrapper frames (they should be transparent containers)
- **`SPACE_BETWEEN`** for navigation bars and tab rows
- **`clipsContent: true`** on the screen frame to simulate a viewport

## Design Tokens for Consistent Spacing

Use `design.compute_tokens` to generate spacing and sizing values tuned to your prototype canvas:

```bash
ai-happy-design command design.compute_tokens '{"width":375,"height":812}'
```

The returned spacing scale ensures consistent gaps throughout your prototype:

| Token | Typical Value (375px) | Use For |
|-------|----------------------|---------|
| `spacing.xs` | 4px | Inline element gaps |
| `spacing.sm` | 8px | Tight component padding |
| `spacing.md` | 16px | Card padding, list gaps |
| `spacing.lg` | 24px | Section spacing |
| `spacing.xl` | 32px | Page margins |
| `spacing.xxl` | 48px | Major section breaks |

Apply these tokens consistently across all frames for a cohesive feel. The font size scale is equally important -- use the returned `fonts.body`, `fonts.heading`, etc. values rather than guessing.

## Rapid Iteration Workflow

The speed of CLI-based prototyping enables a tight iteration loop.

### 1. Draft the Layout

Start with a rough batch file defining the structure:

```bash
ai-happy-design batch wireframe-v1.json --lint
```

### 2. Review in Figma

The frames appear instantly on your canvas. Check layout, spacing, and hierarchy visually.

### 3. Modify Individual Elements

Use `node.modify` to adjust specific nodes without recreating the whole layout:

```bash
# Change text content
ai-happy-design command text.set_content '{"nodeId":"42:248","text":"New headline"}'

# Change color
ai-happy-design command paint.set_solid '{"nodeId":"42:248","color":"#FF0000"}'

# Resize and reposition
ai-happy-design command node.modify \
  '{"nodeId":"42:248","width":400,"height":600,"y":120}'
```

### 4. Find and Modify

Search for nodes, then modify them:

```bash
# Find all text nodes named "Price"
ai-happy-design command document.find_nodes '{"query":"Price","nodeType":"TEXT"}'

# Update found nodes
ai-happy-design command node.modify '{"nodeId":"FOUND_ID","text":"$99/mo","color":"#FFD600"}'
```

### 5. Clone and Iterate

```bash
# Clone an existing card to create a variation
ai-happy-design command node.clone '{"nodeId":"42:248"}'

# Move the clone next to the original
ai-happy-design command node.move '{"nodeId":"CLONE_ID","x":500,"y":0}'

# Modify the clone
ai-happy-design command text.set_content '{"nodeId":"TEXT_IN_CLONE","text":"Variation B"}'
```

### 6. Export and Share

```bash
ai-happy-design command export.image \
  '{"nodeId":"42:248","format":"PNG","scale":2}'
```

## Inspecting Existing Designs

Use `node.get_tree` to reverse-engineer any existing Figma design:

```bash
# Full tree with nesting
ai-happy-design command node.get_tree '{"nodeId":"1:2","depth":3}'

# Compact flat array (3-5x fewer tokens for LLM context)
ai-happy-design command node.get_tree '{"nodeId":"1:2","compact":true}'
```

The compact format returns a flat array where each node includes its `parentId`:

```json
[
  {"id":"42:248","type":"FRAME","name":"Card","x":0,"y":0,"w":400,"h":300,"childCount":3,"depth":0},
  {"id":"42:249","type":"TEXT","name":"Title","x":32,"y":32,"w":336,"h":34,"parentId":"42:248","depth":1},
  {"id":"42:250","type":"TEXT","name":"Description","x":32,"y":80,"w":336,"h":48,"parentId":"42:248","depth":1}
]
```

### Search for Specific Elements

```bash
# Find all text nodes
ai-happy-design command document.find_nodes '{"nodeType":"TEXT"}'

# Find by name
ai-happy-design command document.find_nodes '{"query":"Button"}'

# Find by text content
ai-happy-design command document.find_nodes '{"textContent":"Submit"}'
```

## Wireframe to High-Fidelity

The transition from wireframe to high-fidelity is straightforward with `node.modify` and effect commands:

1. **Replace placeholder colors** -- swap `#E0E0E0` grays for brand colors
2. **Apply image fills** -- replace rectangles with actual images via `paint.set_image` or `paint.set_image_url`
3. **Add effects** -- drop shadows on cards, glass morphism on overlays
4. **Refine typography** -- adjust font sizes, weights, and line heights
5. **Add gradients** -- subtle background gradients for depth

```bash
# Add a shadow to a card
ai-happy-design command effect.add_shadow \
  '{"nodeId":"42:248","color":"#00000020","offsetX":0,"offsetY":4,"radius":12}'

# Apply glass morphism to an overlay
ai-happy-design command effect.apply_glass \
  '{"nodeId":"42:260","intensity":"light"}'

# Replace placeholder with image
ai-happy-design command paint.set_image_url \
  '{"nodeId":"42:252","url":"https://example.com/photo.jpg","scaleMode":"FILL"}'
```

## Tips for Prototyping

- **Use auto-layout** -- it handles spacing, alignment, and makes future changes easier
- **Start with wireframe grays** (#F5F5F5, #E0E0E0, #888888) and refine colors later
- **Set `clipsContent: true`** on screen frames to simulate device viewports
- **Use `FILL` sizing** on children that should stretch to match their parent width
- **Name every layer** -- "Header", "Content Area", "Tab Bar", not "Frame 47"
- **Use `node.get_tree compact:true`** to understand existing layouts before modifying
- **Use `design.compute_tokens`** to get consistent spacing and font sizes for any canvas
- **Export at 1x** for quick iteration, 2x for final review
- **Use `node.modify`** for quick multi-property updates on a single node
- **Clone and iterate** -- duplicate frames to explore variations side by side
- **Run `--lint`** after every batch to catch overflow and overlap issues

---

Made by [Ashraf Ali](https://ashrafali.net) | [GitHub](https://github.com/nerveband/ai-happy-design) | License: GPL-3.0
