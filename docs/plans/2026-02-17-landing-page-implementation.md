# AI Happy Design Landing Page — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a bold, playful single-page marketing site for AI Happy Design with scroll-triggered animations, a live terminal demo, and clear CTAs.

**Architecture:** Single HTML file with separate CSS and JS files. No frameworks, no build tools. Google Fonts via CDN. IntersectionObserver for scroll animations. Deployed to Netlify via CLI from a private GitHub repo.

**Tech Stack:** HTML5, CSS3 (custom properties, grid, flexbox), vanilla JS (IntersectionObserver, typed animation), Google Fonts (Space Grotesk, Inter, JetBrains Mono)

---

## Pre-Implementation Setup

### Task 0: Create site directory and private GitHub repo

**Files:**
- Create: `site/index.html`
- Create: `site/styles.css`
- Create: `site/script.js`
- Create: `site/netlify.toml`

**Step 1: Create directory structure**

```bash
mkdir -p site
```

**Step 2: Create netlify.toml**

```toml
[build]
  publish = "."

[[headers]]
  for = "/*"
  [headers.values]
    X-Frame-Options = "DENY"
    X-Content-Type-Options = "nosniff"
    Referrer-Policy = "strict-origin-when-cross-origin"
```

**Step 3: Create private GitHub repo**

```bash
cd site
git init
gh repo create ai-happy-design-site --private --source=. --push
```

**Step 4: Commit skeleton**

```bash
git add -A && git commit -m "chore: init landing page skeleton"
```

---

## Page Construction (Section by Section)

### Task 1: HTML skeleton + CSS reset + custom properties

**Files:**
- Create: `site/index.html` — full HTML skeleton with all section placeholders
- Create: `site/styles.css` — CSS reset, custom properties, base typography, responsive breakpoints

**Step 1: Write index.html**

Complete HTML document with:
- `<head>`: meta tags (charset, viewport, description, OG tags), Google Fonts preconnect + link, CSS link
- `<body>`: semantic sections (`<nav>`, `<header>`, `<section>` x7, `<footer>`), JS link at bottom
- Each section has `id` and `class` for styling/animation targeting
- All text content populated (headlines, descriptions, code snippets)
- SVG icons inline (lightning bolt, calculator, grid, layers, zap, lock, puzzle, chat bubble, diamond)
- Terminal window markup with typed-text spans
- Feature cards markup
- Speed comparison bar markup
- Tabbed code blocks markup
- Integration logos (text-based, no external images)
- CTA buttons and social links

OG tags:
- `og:title`: "AI Happy Design — Give AI full control of your Figma canvas"
- `og:description`: "Free, open source Figma automation. 9x faster than MCP. One binary."
- `og:type`: "website"

**Step 2: Write styles.css**

CSS organized in sections:
1. **Reset** — box-sizing, margin, font smoothing
2. **Custom properties** — all color tokens, spacing scale, transition defaults
3. **Base** — html/body, typography defaults, link styles
4. **Nav** — sticky, transparent→solid transition, layout
5. **Hero** — split layout, responsive stack on mobile
6. **Terminal** — window chrome, typing cursor animation, glow
7. **Narrative** — centered text, merge animation keyframes
8. **How It Works** — 3-column, dashed line animation
9. **Features** — 2x3 grid, card hover glow
10. **Speed** — bar chart layout, animated width
11. **DX** — tab buttons, code block styling
12. **Integrations** — horizontal flex, hover color transition
13. **CTA** — centered, button styles
14. **Footer** — simple flex layout
15. **Animations** — fade-up keyframes, intersection states
16. **Responsive** — mobile breakpoints (768px, 480px)

Key CSS patterns:
- `.fade-up` class: `opacity:0; transform:translateY(30px)` → `.visible`: `opacity:1; transform:translateY(0)`
- `.terminal-glow`: `box-shadow: 0 0 60px rgba(255,214,0,0.15)`
- `.accent-pill`: yellow bg, dark text, rounded, small text
- `.btn-primary`: yellow bg `#FFD600`, dark text `#0a0a0a`, hover darken
- `.btn-outline`: transparent bg, white border, white text, hover yellow

**Step 3: Open in browser and verify structure renders**

```bash
open site/index.html
```

**Step 4: Commit**

```bash
cd site && git add -A && git commit -m "feat: HTML skeleton + CSS foundation with all sections"
```

---

### Task 2: Terminal typing animation (hero)

**Files:**
- Create: `site/script.js` — typing animation, scroll observer, tab switching

**Step 1: Write script.js**

Three systems:

**A) Typing Animation:**
- Array of demo sequences, each with `command` string and `result` string
- Types command character-by-character (40ms per char)
- After command completes, fades in result block
- Pauses 2s, clears, types next command
- Loops forever
- Blinking cursor via CSS `@keyframes blink`

**B) Scroll-triggered animations:**
- `IntersectionObserver` with `threshold: 0.15`
- Observes all `.fade-up` elements
- Adds `.visible` class when element enters viewport
- `once: true` — don't re-hide on scroll up

**C) Tab switching (DX section):**
- Click handler on tab buttons
- Shows/hides corresponding code block
- Updates active tab button state

**D) Nav background:**
- Scroll listener (throttled via `requestAnimationFrame`)
- Adds `.nav-solid` class when `scrollY > 50`

**E) Speed bar animation:**
- IntersectionObserver on speed section
- Triggers CSS width transition on bars when visible

**Step 2: Verify terminal animation works in browser**

```bash
open site/index.html
```

**Step 3: Commit**

```bash
cd site && git add script.js && git commit -m "feat: terminal typing animation + scroll reveals + tab switching"
```

---

### Task 3: Narrative merge animation

**Files:**
- Modify: `site/styles.css` — add merge animation keyframes

**Step 1: Add CSS keyframes for the "Design + Code merge" animation**

- Two boxes start 200px apart
- On scroll-trigger, they slide toward center
- At center they overlap and a glow expands
- Lightning bolt icon fades in at the merge point
- Pure CSS animation triggered by `.visible` class

**Step 2: Verify in browser**

**Step 3: Commit**

```bash
cd site && git add styles.css && git commit -m "feat: design-code merge animation"
```

---

### Task 4: How It Works animated connections

**Files:**
- Modify: `site/styles.css` — dashed line animation between steps

**Step 1: Add animated dashed lines**

- CSS `border-top: 2px dashed #FFD600` between the three step circles
- `@keyframes dash-flow` using `background-position` shift to simulate flow
- Triggered on scroll into view

**Step 2: Verify animation flows left-to-right**

**Step 3: Commit**

```bash
cd site && git add styles.css && git commit -m "feat: animated connection lines in How It Works"
```

---

### Task 5: Speed bars animation

**Files:**
- Modify: `site/styles.css` — animated bar widths

**Step 1: Add bar animation**

- MCP bar: animates from 0% to 25% width (gray)
- CLI bar: animates from 0% to 90% width (yellow)
- Number counters (optional CSS approach or JS)
- "9x faster" badge pulses once

**Step 2: Verify bars animate on scroll**

**Step 3: Commit**

```bash
cd site && git add -A && git commit -m "feat: animated speed comparison bars"
```

---

### Task 6: Responsive design pass

**Files:**
- Modify: `site/styles.css` — mobile breakpoints

**Step 1: Add responsive rules**

At `max-width: 768px`:
- Hero: stack vertically (text above terminal)
- Features: 1-column stack
- How It Works: vertical stack
- Speed bars: full width
- Nav: smaller text, hamburger optional (or just shrink)
- Terminal: full width, smaller font

At `max-width: 480px`:
- Reduce headline sizes
- Increase padding for touch targets
- Stack CTA buttons vertically

**Step 2: Test at 375px, 768px, 1024px, 1440px widths**

**Step 3: Commit**

```bash
cd site && git add styles.css && git commit -m "feat: responsive design for mobile and tablet"
```

---

### Task 7: Copy lightning bolt SVG icon from resources

**Files:**
- Read: `resources/icon-optimized.svg` for the lightning bolt path data
- Modify: `site/index.html` — inline SVG in nav and merge animation

**Step 1: Extract SVG path data from project resources**

Use the optimized SVG icon. Inline it in:
- Nav logo (28x28)
- Narrative merge point (48x48)
- How It Works relay step (32x32)

**Step 2: Verify icons render correctly**

**Step 3: Commit**

```bash
cd site && git add index.html && git commit -m "feat: inline brand SVG icons from resources"
```

---

## Quality Assurance (All 11 Skills)

### Task 8: Run all web quality skills

Each sub-step uses one of the 11 web design skills listed in requirements. Run them sequentially, fix issues found by each before moving to the next.

**Step 1: frontend-design skill** — Review overall design quality and polish

**Step 2: web-design-guidelines skill** — Check against Web Interface Guidelines

**Step 3: accessibility skill** — WCAG 2.1 audit (contrast, aria labels, keyboard nav, focus states)

**Step 4: seo skill** — Meta tags, structured data, semantic HTML, heading hierarchy

**Step 5: performance skill** — Optimize loading (font loading strategy, image optimization, CSS/JS size)

**Step 6: core-web-vitals skill** — Check LCP, INP, CLS scores

**Step 7: best-practices skill** — Security headers, code quality, modern patterns

**Step 8: web-quality-audit skill** — Comprehensive audit combining all above

**Step 9: interface-design skill** — Review interactive elements (tabs, buttons, animations)

**Step 10: Fix all issues found, commit**

```bash
cd site && git add -A && git commit -m "fix: address quality audit findings"
```

---

## Branding Guidelines Document

### Task 9: Create branding guidelines

**Files:**
- Create: `site/brand-guidelines.md` (or HTML version)

**Step 1: Write branding guidelines covering:**

- Logo usage (lightning bolt icon, spacing rules, minimum sizes)
- Color palette (primary, accent, backgrounds, text, semantic colors)
- Typography (font families, weights, sizes, hierarchy)
- Tone of voice (bold, developer-friendly, no-nonsense, playful)
- Do's and don'ts (examples of correct/incorrect usage)
- Button styles and interactive states
- Code block styling conventions
- Animation principles

**Step 2: Commit**

```bash
cd site && git add brand-guidelines.md && git commit -m "docs: add branding guidelines"
```

---

## Deployment

### Task 10: Deploy to Netlify

**Step 1: Install Netlify CLI if needed**

```bash
npm install -g netlify-cli
```

**Step 2: Deploy**

```bash
cd site && netlify deploy --prod --dir=.
```

**Step 3: Push to private repo**

```bash
cd site && git push origin main
```

**Step 4: Verify live site loads correctly**

---

## Summary

| Task | Description | Est. |
|------|-------------|------|
| 0 | Setup: directory, repo, skeleton | 2 min |
| 1 | HTML + CSS foundation (all sections) | 15 min |
| 2 | Terminal typing animation + scroll reveals | 10 min |
| 3 | Narrative merge animation | 5 min |
| 4 | How It Works animated connections | 5 min |
| 5 | Speed bars animation | 5 min |
| 6 | Responsive design pass | 10 min |
| 7 | SVG icons from resources | 5 min |
| 8 | Quality audit (11 skills) | 15 min |
| 9 | Branding guidelines doc | 10 min |
| 10 | Deploy to Netlify | 5 min |
