# AI Happy Design — Brand Guidelines

> **"Give AI full control of your Figma canvas"**

---

## Table of Contents

1. [Brand Overview](#1-brand-overview)
2. [Logo & Icon](#2-logo--icon)
3. [Color Palette](#3-color-palette)
4. [Typography](#4-typography)
5. [Tone of Voice](#5-tone-of-voice)
6. [UI Components](#6-ui-components)
7. [Animation Principles](#7-animation-principles)
8. [Do's and Don'ts](#8-dos-and-donts)
9. [Social & Marketing](#9-social--marketing)

---

## 1. Brand Overview

| Attribute   | Value                                      |
|-------------|--------------------------------------------|
| Product     | AI Happy Design                            |
| Tagline     | "Give AI full control of your Figma canvas"|
| Category    | Open source developer tool (Figma automation) |
| License     | MIT                                        |
| Creator     | Ashraf Ali                                 |
| Social      | @workandpromise on X                       |
| Email       | info@workandpromise.com                    |
| GitHub      | github.com/nerveband/ai-happy-design       |

AI Happy Design is a developer tool that connects AI coding agents directly to Figma. It ships as a single binary, works with every major AI editor, and speaks the language developers already know: CLIs, JSON, and raw throughput.

---

## 2. Logo & Icon

### Primary Mark

The primary mark is a **lightning bolt on a rounded square**. It represents speed, energy, and the "spark" of AI-powered design.

### Specifications

| Property       | Value                            |
|----------------|----------------------------------|
| Symbol         | Lightning bolt                   |
| Bolt color     | #FFD600 (Brand Yellow)           |
| Background     | #0a0a0a (Dark Base)              |
| Corner radius  | Rounded square (proportional)    |
| Minimum size   | 16 x 16 px                       |
| Clear space    | At least 50% of icon width on all sides |

### Icon SVG Reference

```html
<svg viewBox="0 0 28 28" xmlns="http://www.w3.org/2000/svg">
  <rect width="28" height="28" rx="6" fill="#0a0a0a"/>
  <path d="M15.5 4L8 15h4.5l-1 9L20 13h-4.5l1-9z" fill="#FFD600"/>
</svg>
```

### Usage Rules

- Always render the bolt on a dark background (#0a0a0a or darker).
- Never stretch, rotate, skew, or add effects (drop shadows, outlines, glows) to the icon.
- Use inline SVG for web rendering. Figma's Content Security Policy blocks `data:` URIs set dynamically via JavaScript and CSS `background-image` with base64.
- At sizes below 16px, the bolt loses legibility. Do not use the icon smaller than 16 x 16 px.

---

## 3. Color Palette

### Primary Colors

| Name         | Hex       | RGB              | Usage                                  |
|--------------|-----------|------------------|----------------------------------------|
| Brand Yellow | `#FFD600` | 255, 214, 0      | Primary accent, CTAs, highlights, icon |
| Dark Base    | `#0a0a0a` | 10, 10, 10       | Page background, primary surfaces      |
| Surface      | `#161616` | 22, 22, 22       | Cards, elevated surfaces               |
| Dark Alt     | `#111111` | 17, 17, 17       | Alternate section backgrounds          |

### Supporting Colors

| Name         | Hex       | RGB              | Usage                                  |
|--------------|-----------|------------------|----------------------------------------|
| Border       | `#2a2a2a` | 42, 42, 42       | Subtle borders, dividers               |
| Text Primary | `#FFFFFF` | 255, 255, 255    | Headlines, primary text                |
| Text Muted   | `#888888` | 136, 136, 136    | Secondary text, descriptions           |
| Terminal BG  | `#1a1a1a` | 26, 26, 26       | Code blocks, terminal windows          |
| Code Text    | `#e0e0e0` | 224, 224, 224    | Code text color                        |

### Accent Variations

| Name               | Value                        | Usage                          |
|--------------------|------------------------------|--------------------------------|
| Accent Hover       | `#e6c200`                    | Button hover states            |
| Accent Glow        | `rgba(255, 214, 0, 0.15)`   | Subtle glows, halos            |
| Accent Glow Strong | `rgba(255, 214, 0, 0.3)`    | Hover glows, emphasis          |
| Accent Tint        | `rgba(255, 214, 0, 0.1)`    | Badge/pill backgrounds         |

### Status Colors

| Name    | Hex                   | Usage                            |
|---------|-----------------------|----------------------------------|
| Success | `#22a355` / `#28c840` | Connected state, success messages|
| Warning | `#febc2e`             | Connecting state, cautions       |
| Error   | `#ff5f57`             | Disconnected, errors             |

### CSS Custom Properties

```css
:root {
  /* Primary */
  --color-brand-yellow: #FFD600;
  --color-dark-base: #0a0a0a;
  --color-surface: #161616;
  --color-dark-alt: #111111;

  /* Supporting */
  --color-border: #2a2a2a;
  --color-text-primary: #FFFFFF;
  --color-text-muted: #888888;
  --color-terminal-bg: #1a1a1a;
  --color-code-text: #e0e0e0;

  /* Accent */
  --color-accent-hover: #e6c200;
  --color-accent-glow: rgba(255, 214, 0, 0.15);
  --color-accent-glow-strong: rgba(255, 214, 0, 0.3);
  --color-accent-tint: rgba(255, 214, 0, 0.1);

  /* Status */
  --color-success: #22a355;
  --color-success-bright: #28c840;
  --color-warning: #febc2e;
  --color-error: #ff5f57;
}
```

---

## 4. Typography

### Font Stack

| Role      | Family          | Weights       | Purpose                          |
|-----------|-----------------|---------------|----------------------------------|
| Headlines | Space Grotesk   | 700           | Bold, geometric, modern          |
| Body      | Inter           | 400, 500, 600 | Clean, highly readable           |
| Code      | JetBrains Mono  | 400           | Designed for code readability    |

### Type Scale

| Level           | Size   | Family         | Weight | Line Height |
|-----------------|--------|----------------|--------|-------------|
| Display / Hero  | 56px   | Space Grotesk  | 700    | 1.1         |
| Section Header  | 42px   | Space Grotesk  | 700    | 1.1         |
| Card Title      | 18px   | Inter          | 600    | 1.1         |
| Body            | 16px   | Inter          | 400    | 1.6         |
| Small / Muted   | 14-15px| Inter          | 400-500| 1.6         |
| Code            | 13-14px| JetBrains Mono | 400    | 1.5         |

### CSS Font Definitions

```css
/* Google Fonts import */
@import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@700&family=Inter:wght@400;500;600&family=JetBrains+Mono:wght@400&display=swap');

:root {
  --font-headline: 'Space Grotesk', system-ui, sans-serif;
  --font-body: 'Inter', system-ui, sans-serif;
  --font-code: 'JetBrains Mono', 'Fira Code', monospace;
}

h1, h2, h3 {
  font-family: var(--font-headline);
  font-weight: 700;
  line-height: 1.1;
}

body {
  font-family: var(--font-body);
  font-weight: 400;
  line-height: 1.6;
  font-size: 16px;
}

code, pre {
  font-family: var(--font-code);
  font-weight: 400;
  font-size: 14px;
  line-height: 1.5;
}
```

---

## 5. Tone of Voice

### Principles

| Principle                 | Description                                                                 |
|---------------------------|-----------------------------------------------------------------------------|
| **Bold**                  | State features directly. No hedging or "might" language.                    |
| **Developer-first**       | Speak the developer's language. CLIs, binaries, JSON, ops/sec.              |
| **No-nonsense**           | Get to the point. Short sentences. Active voice.                            |
| **Playful when appropriate** | The product name is "Happy Design" — there is room for energy and personality. |
| **Never condescending**   | Assume the reader is smart. Do not over-explain.                            |

### Examples

**Do say:**

- "One binary. Every AI editor."
- "27 ops/sec. Zero overhead."
- "Type a command. Get a design."
- "Open source. MIT licensed. Fork it, ship it, break it."

**Do not say:**

- "Simply leverage our innovative solution..."
- "AI Happy Design is the best tool for..."
- "As a revolutionary breakthrough..."
- "We are excited to announce..."
- "Our cutting-edge AI-powered platform..."

### Writing Checklist

- [ ] Can you cut any words? Cut them.
- [ ] Is every sentence in active voice?
- [ ] Would a developer nod reading this, or roll their eyes?
- [ ] Does it sound like a person talking, not a press release?

---

## 6. UI Components

### Buttons

#### Primary Button

```css
.btn-primary {
  background: #FFD600;
  color: #0a0a0a;
  border: none;
  border-radius: 10px;
  padding: 14px 28px;
  font-family: 'Inter', sans-serif;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.2s ease, transform 0.2s ease;
}

.btn-primary:hover {
  background: #e6c200;
  transform: translateY(-1px);
}
```

#### Outline Button

```css
.btn-outline {
  background: transparent;
  color: #FFFFFF;
  border: 1px solid #FFFFFF;
  border-radius: 10px;
  padding: 14px 28px;
  font-family: 'Inter', sans-serif;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  transition: border-color 0.2s ease, color 0.2s ease, transform 0.2s ease;
}

.btn-outline:hover {
  border-color: #FFD600;
  color: #FFD600;
  transform: translateY(-1px);
}
```

#### Small Button

```css
.btn-small {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 13px;
}
```

### Pills / Badges

```css
.pill {
  display: inline-block;
  background: rgba(255, 214, 0, 0.1);
  border: 1px solid rgba(255, 214, 0, 0.3);
  color: #FFD600;
  padding: 6px 16px;
  border-radius: 100px;
  font-family: 'Inter', sans-serif;
  font-weight: 500;
  font-size: 14px;
}
```

### Cards

```css
.card {
  background: #161616;
  border: 1px solid #2a2a2a;
  border-radius: 16px;
  padding: 32px;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

.card:hover {
  border-color: #FFD600;
  box-shadow: 0 0 40px rgba(255, 214, 0, 0.08);
}
```

### Terminal Windows

```css
.terminal {
  background: #1a1a1a;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 0 80px rgba(255, 214, 0, 0.15);
}

.terminal-chrome {
  background: #2a2a2a;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.terminal-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.terminal-dot--red    { background: #ff5f57; }
.terminal-dot--yellow { background: #febc2e; }
.terminal-dot--green  { background: #28c840; }

.terminal-body {
  padding: 20px 24px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  color: #e0e0e0;
}
```

---

## 7. Animation Principles

### Core Rules

1. **Subtle and purposeful.** Every animation guides attention. Nothing decorates.
2. **No animation exceeds 1.5 seconds.**
3. **Always respect `prefers-reduced-motion`.** Disable all motion when the user requests it.

### Standard Animations

| Animation         | Duration | Easing   | Properties                       | Use Case                        |
|-------------------|----------|----------|----------------------------------|---------------------------------|
| Fade-up reveal    | 0.8s     | ease     | opacity 0 to 1, translateY 30px to 0 | Section entrances          |
| Sibling stagger   | +0.1s    | ease     | Same as fade-up, delayed         | Multiple cards, list items      |
| Glow pulse        | 1.5s     | ease     | box-shadow opacity oscillation   | Key moments (merge, speed demo) |
| Terminal typing   | 30ms/char| linear   | Character-by-character append    | Code demo sequences             |
| Button hover      | 0.2s     | ease     | background, translateY(-1px)     | Interactive feedback            |
| Card hover        | 0.3s     | ease     | border-color, box-shadow         | Hover emphasis                  |

### CSS Reference

```css
/* Fade-up reveal */
.reveal {
  opacity: 0;
  transform: translateY(30px);
  transition: opacity 0.8s ease, transform 0.8s ease;
}

.reveal.visible {
  opacity: 1;
  transform: translateY(0);
}

/* Stagger children */
.reveal:nth-child(2) { transition-delay: 0.1s; }
.reveal:nth-child(3) { transition-delay: 0.2s; }
.reveal:nth-child(4) { transition-delay: 0.3s; }

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 8. Do's and Don'ts

### Do

- Use the yellow accent **sparingly** -- it should pop, not overwhelm.
- Keep dark backgrounds consistent. `#0a0a0a` is the base everywhere.
- Use Space Grotesk for **all** headlines. No exceptions.
- Maintain generous whitespace between sections (120px vertical padding).
- Use inline SVGs for icons (web rendering, Figma compatibility).
- Keep code examples **real and functional**. Every snippet should work if pasted.
- Include the MIT License badge on all marketing materials.

### Don't

- Use Brand Yellow (`#FFD600`) on light backgrounds. The contrast is too low.
- Introduce additional accent colors. Yellow is the only accent.
- Use more than two font families on a single page (headline + body, or body + code).
- Add unnecessary decoration or ornamentation. Let content and whitespace do the work.
- Use stock photos. The product speaks through its UI, code, and terminal output.
- Over-animate. One reveal per scroll position. No parallax, no bouncing, no spinning.

---

## 9. Social & Marketing

### Channels

| Channel  | Handle / URL                                  |
|----------|-----------------------------------------------|
| X        | [@workandpromise](https://x.com/workandpromise) (always lowercase) |
| Email    | info@workandpromise.com                       |
| GitHub   | [github.com/nerveband/ai-happy-design](https://github.com/nerveband/ai-happy-design) |

### Naming Conventions

- **First mention:** Always use the full name "AI Happy Design".
- **Subsequent mentions:** "AHD" is acceptable after the full name has been introduced.
- **Never:** "AI-Happy-Design", "aihappydesign", "Happy Design" (without "AI"), or "AIHD".

### Required Badges

The MIT License badge should be visible on all marketing materials and landing pages:

```markdown
[![MIT License](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)
```

### Social Copy Guidelines

Follow the [Tone of Voice](#5-tone-of-voice) rules. Keep social posts:

- Under 280 characters when possible (X-friendly).
- Focused on one feature or metric per post.
- Ending with a link or CTA, not a question.

**Good social post:**
> AI Happy Design v0.5 is out. 27 ops/sec, single binary, works with Claude, Cursor, Copilot, and Windsurf. MIT licensed.
> github.com/nerveband/ai-happy-design

**Bad social post:**
> We are SO excited to announce the latest version of our amazing tool! It does so many things and we can't wait for you to try it!!!

---

## Quick Reference Card

| Element          | Value                              |
|------------------|------------------------------------|
| Brand Yellow     | `#FFD600`                          |
| Dark Base        | `#0a0a0a`                          |
| Surface          | `#161616`                          |
| Border           | `#2a2a2a`                          |
| Headline font    | Space Grotesk 700                  |
| Body font        | Inter 400                          |
| Code font        | JetBrains Mono 400                 |
| Button radius    | 10px                               |
| Card radius      | 16px                               |
| Section spacing  | 120px vertical padding             |
| Min icon size    | 16 x 16 px                         |
| Icon clear space | 50% of icon width, all sides       |

---

*AI Happy Design is open source software released under the MIT License.*
*Created by Ashraf Ali. Learn more at [github.com/nerveband/ai-happy-design](https://github.com/nerveband/ai-happy-design).*
