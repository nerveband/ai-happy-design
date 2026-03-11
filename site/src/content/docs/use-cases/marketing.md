---
title: Marketing Materials
description: Create email banners, ad creatives, social media kits, and event graphics at scale with AI Happy Design.
---

AI Happy Design streamlines marketing asset production -- from email headers to multi-size ad creatives. Define a template once, then batch-produce every size your campaign needs.

## Common Canvas Sizes

| Format | Size (px) | Use Case |
|--------|----------|----------|
| Email Banner | 600 x 200 | Newsletter headers |
| Email Hero | 600 x 400 | Feature announcements |
| Google Display Ad | 300 x 250 | Medium rectangle |
| Google Display Ad | 728 x 90 | Leaderboard |
| Google Display Ad | 160 x 600 | Wide skyscraper |
| Facebook Ad | 1200 x 628 | Feed ad |
| Instagram Ad | 1080 x 1080 | Square feed ad |
| LinkedIn Sponsored | 1200 x 627 | Sponsored content |
| OG Image | 1200 x 630 | Link previews |
| Event Banner | 1920 x 1080 | Web hero / event page |

## Email Banners

Use the built-in banner composite for quick email headers:

```json
[
  {
    "command": "bulk.banner",
    "params": {
      "canvas": "600x200",
      "bg": "#FFD600",
      "elements": [
        { "type": "headline", "text": "Spring Sale — 40% Off Everything" },
        { "type": "subtitle", "text": "Use code SPRING40 at checkout" }
      ]
    }
  }
]
```

### Custom Email Header

For more control, build it manually with a CTA button:

```json
[
  {
    "name": "header",
    "command": "frame",
    "params": {
      "name": "Email Header",
      "w": 600, "h": 200,
      "bg": "#0C1E2C",
      "layoutMode": "HORIZONTAL",
      "padding": 32,
      "primaryAxisAlign": "SPACE_BETWEEN",
      "counterAxisAlign": "CENTER",
      "itemSpacing": 24
    }
  },
  {
    "name": "text_group",
    "command": "frame",
    "params": {
      "pid": "${{steps.header.result.id}}",
      "name": "Text Group",
      "layoutMode": "VERTICAL",
      "itemSpacing": 8,
      "noFill": true,
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "headline",
    "command": "text",
    "params": {
      "pid": "${{steps.text_group.result.id}}",
      "text": "Your Weekly Design Digest",
      "sz": 28,
      "ff": "Space Grotesk",
      "fontStyle": "Bold",
      "color": "#FFFFFF"
    }
  },
  {
    "name": "subtitle",
    "command": "text",
    "params": {
      "pid": "${{steps.text_group.result.id}}",
      "text": "The latest tools, tips, and trends in design",
      "sz": 16,
      "ff": "Inter",
      "color": "#888888"
    }
  },
  {
    "name": "cta_btn",
    "command": "frame",
    "params": {
      "pid": "${{steps.header.result.id}}",
      "name": "CTA Button",
      "layoutMode": "HORIZONTAL",
      "padding": 12,
      "paddingLeft": 24,
      "paddingRight": 24,
      "bg": "#FFD600",
      "r": 8,
      "primaryAxisAlign": "CENTER",
      "counterAxisAlign": "CENTER"
    }
  },
  {
    "name": "cta_label",
    "command": "text",
    "params": {
      "pid": "${{steps.cta_btn.result.id}}",
      "text": "Read Now",
      "sz": 14,
      "ff": "Inter",
      "fontStyle": "Bold",
      "color": "#0C1E2C"
    }
  }
]
```

## Ad Creatives at Multiple Sizes

A typical display ad campaign requires the same creative in several sizes. Compute tokens per canvas to get the right font sizes automatically.

### Step 1: Compute Tokens Per Size

```bash
ai-happy-design command design.compute_tokens '{"width":300,"height":250}'
ai-happy-design command design.compute_tokens '{"width":728,"height":90}'
ai-happy-design command design.compute_tokens '{"width":1080,"height":1080}'
```

Each call returns font sizes and spacing tuned to that canvas. Use the returned values in your batch JSON.

### Step 2: Create Per-Size Batch Files

```bash
# Instagram (1080x1080)
ai-happy-design batch ad-square.json

# Facebook (1200x628)
ai-happy-design batch ad-landscape.json

# Story (1080x1920)
ai-happy-design batch ad-story.json
```

### Multi-Size in One Batch

Create all sizes in a single batch run. Each frame auto-places on the canvas:

```json
[
  {
    "name": "sq",
    "command": "frame",
    "params": { "name": "Ad - Square", "w": 1080, "h": 1080, "bg": "#0C1E2C" }
  },
  {
    "name": "sq_headline",
    "command": "text",
    "params": {
      "pid": "${{steps.sq.result.id}}",
      "text": "Try AI Happy Design",
      "sz": 64, "ff": "Space Grotesk", "fontStyle": "Bold",
      "color": "#FFFFFF", "x": 80, "y": 200, "w": 920
    }
  },
  {
    "name": "sq_body",
    "command": "text",
    "params": {
      "pid": "${{steps.sq.result.id}}",
      "text": "CLI-first Figma automation for teams and AI agents.",
      "sz": 36, "ff": "Inter",
      "color": "#888888", "x": 80, "y": 320, "w": 920,
      "lh": 150, "lineHeightUnit": "PERCENT"
    }
  },
  {
    "name": "fb",
    "command": "frame",
    "params": { "name": "Ad - Facebook", "w": 1200, "h": 628, "bg": "#0C1E2C" }
  },
  {
    "name": "fb_headline",
    "command": "text",
    "params": {
      "pid": "${{steps.fb.result.id}}",
      "text": "Try AI Happy Design",
      "sz": 48, "ff": "Space Grotesk", "fontStyle": "Bold",
      "color": "#FFFFFF", "x": 60, "y": 120, "w": 600
    }
  },
  {
    "name": "fb_body",
    "command": "text",
    "params": {
      "pid": "${{steps.fb.result.id}}",
      "text": "CLI-first Figma automation for teams and AI agents.",
      "sz": 24, "ff": "Inter",
      "color": "#888888", "x": 60, "y": 200, "w": 600,
      "lh": 150, "lineHeightUnit": "PERCENT"
    }
  },
  {
    "name": "lb",
    "command": "frame",
    "params": { "name": "Ad - Leaderboard", "w": 728, "h": 90, "bg": "#0C1E2C",
      "layoutMode": "HORIZONTAL", "padding": 20, "itemSpacing": 24,
      "primaryAxisAlign": "SPACE_BETWEEN", "counterAxisAlign": "CENTER" }
  },
  {
    "name": "lb_headline",
    "command": "text",
    "params": {
      "pid": "${{steps.lb.result.id}}",
      "text": "Design at Scale with AI Happy Design",
      "sz": 20, "ff": "Space Grotesk", "fontStyle": "Bold",
      "color": "#FFFFFF"
    }
  }
]
```

## Social Media Kits

When a campaign needs the same creative across Instagram, LinkedIn, Twitter, and Facebook, batch all formats in a single run or as separate files.

### Kit Batch Strategy

Save each format as a separate JSON file:

```
campaign/
  instagram-post.json    (1080x1080)
  instagram-story.json   (1080x1920)
  linkedin-post.json     (1200x628)
  twitter-card.json      (1200x675)
  facebook-post.json     (1200x630)
```

Run them all:

```bash
for file in campaign/*.json; do
  ai-happy-design batch "$file" --lint
done
```

### Adapting Content Per Format

The key to multi-format kits is adapting text tiers to each canvas:

| Canvas | Headline Tier | Body Tier |
|--------|--------------|-----------|
| 1080x1080 (square) | heading (84px) | body (48px) |
| 1080x1920 (story) | hero (152px) | subheading (64px) |
| 1200x628 (landscape) | heading (52px) | body (28px) |
| 728x90 (banner) | subheading (20px) | caption (12px) |

Call `design.compute_tokens` for each canvas to get exact values rather than guessing.

## Event Graphics

Event graphics typically include a web hero, social announcements, and promotional assets.

### Event Hero (1920x1080)

```json
[
  {
    "name": "hero",
    "command": "frame",
    "params": {
      "name": "Event Hero",
      "w": 1920, "h": 1080,
      "bg": "#0C1E2C"
    }
  },
  {
    "name": "hero_gradient",
    "command": "gradient",
    "params": {
      "nodeId": "${{steps.hero.result.id}}",
      "gradientType": "LINEAR",
      "stops": [
        { "position": 0, "color": "#0C1E2C" },
        { "position": 1, "color": "#1a3a5c" }
      ]
    }
  },
  {
    "name": "hero_eyebrow",
    "command": "text",
    "params": {
      "pid": "${{steps.hero.result.id}}",
      "text": "March 15, 2026 / San Francisco",
      "sz": 24,
      "ff": "Inter",
      "color": "#FFD600",
      "x": 120, "y": 320
    }
  },
  {
    "name": "hero_title",
    "command": "text",
    "params": {
      "pid": "${{steps.hero.result.id}}",
      "text": "Design Systems Conference",
      "sz": 84,
      "ff": "Space Grotesk",
      "fontStyle": "Bold",
      "color": "#FFFFFF",
      "x": 120, "y": 380,
      "w": 1200,
      "lh": 110,
      "lineHeightUnit": "PERCENT"
    }
  },
  {
    "name": "hero_subtitle",
    "command": "text",
    "params": {
      "pid": "${{steps.hero.result.id}}",
      "text": "Three days of talks, workshops, and design sprints",
      "sz": 32,
      "ff": "Inter",
      "color": "#888888",
      "x": 120, "y": 580,
      "w": 1000
    }
  }
]
```

## OG Images

Open Graph images for link previews (1200x630):

```json
[
  {
    "name": "og",
    "command": "frame",
    "params": {
      "name": "OG Image",
      "w": 1200, "h": 630,
      "bg": "#0C1E2C",
      "layoutMode": "VERTICAL",
      "padding": 80,
      "primaryAxisAlign": "CENTER",
      "counterAxisAlign": "CENTER",
      "itemSpacing": 24
    }
  },
  {
    "name": "og_title",
    "command": "text",
    "params": {
      "pid": "${{steps.og.result.id}}",
      "text": "Give AI Full Control of Your Figma Canvas",
      "sz": 48,
      "ff": "Space Grotesk",
      "fontStyle": "Bold",
      "color": "#FFFFFF",
      "textAlign": "CENTER",
      "layoutSizingHorizontal": "FILL"
    }
  },
  {
    "name": "og_desc",
    "command": "text",
    "params": {
      "pid": "${{steps.og.result.id}}",
      "text": "Free, open source. 27 ops/sec. Local-only.",
      "sz": 24,
      "ff": "Inter",
      "color": "#888888",
      "textAlign": "CENTER",
      "layoutSizingHorizontal": "FILL"
    }
  }
]
```

## Batch Workflow for Multiple Sizes

The recommended workflow for producing a full marketing campaign:

1. **Define the content** -- headline, body, CTA text, colors, images
2. **List target formats** -- email, display ads, social, event, OG
3. **Compute tokens** for each canvas size:
   ```bash
   ai-happy-design command design.compute_tokens '{"width":1080,"height":1080}'
   ai-happy-design command design.compute_tokens '{"width":600,"height":200}'
   ai-happy-design command design.compute_tokens '{"width":1920,"height":1080}'
   ```
4. **Create one batch file per format** using the computed token values
5. **Run all batches** with lint enabled:
   ```bash
   for file in campaign/*.json; do
     ai-happy-design batch "$file" --lint
   done
   ```
6. **Review in Figma** -- each format appears as its own top-level frame on the canvas
7. **Batch export** all frames:
   ```bash
   ai-happy-design command document.find_nodes '{"nodeType":"FRAME"}'
   ai-happy-design command export.batch \
     '{"nodeIds":"1:2,1:3,1:4,1:5","format":"PNG","scale":2}'
   ```

## Brand Asset Generation

### Logo Lockups

```json
[
  {
    "name": "lockup",
    "command": "frame",
    "params": {
      "name": "Logo Lockup",
      "layoutMode": "HORIZONTAL",
      "itemSpacing": 16,
      "padding": 24,
      "noFill": true
    }
  },
  {
    "name": "brand_name",
    "command": "text",
    "params": {
      "pid": "${{steps.lockup.result.id}}",
      "text": "AI Happy Design",
      "sz": 24,
      "ff": "Space Grotesk",
      "fontStyle": "Bold",
      "color": "#FFD600"
    }
  }
]
```

## Tips for Marketing Assets

- **Compute tokens per canvas size** -- a headline that works at 1080px is unreadable at 300px
- **Use the 8px grid** for all spacing and sizing
- **Export at 2x** for retina-quality output on high-DPI screens
- **Name frames descriptively** ("Ad - Square", "Email Header - Spring Sale")
- **Batch export** all sizes at once with `export.batch`
- **Use gradient backgrounds** for hero-style marketing assets
- **Keep text punchy** -- marketing assets are glanced at, not read
- **Keep CTA buttons prominent** -- high-contrast fill, minimum 36px button height
- **Run `--lint`** on every batch to catch overflow, overlap, and sizing issues before exporting
- **Maintain brand consistency** -- same fonts, same color palette, same padding ratios across all formats

---

Made by [Ashraf Ali](https://ashrafali.net) | [GitHub](https://github.com/nerveband/ai-happy-design) | License: GPL-3.0
