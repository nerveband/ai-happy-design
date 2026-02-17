# AI Happy Design — Landing Page Design

**Date:** 2026-02-17
**Status:** Approved
**Hosting:** Static files (HTML/CSS/JS) → Netlify CLI → Private GitHub repo

---

## Visual System

### Vibe
Bold + playful. Dark background with bright yellow pops. Animated, energetic. Mix of dark sections with glowing accents.

### Color Palette
| Token | Value | Usage |
|-------|-------|-------|
| `--bg-dark` | `#0a0a0a` | Primary background |
| `--bg-alt` | `#111111` | Alternate section background |
| `--surface` | `#161616` | Card surfaces |
| `--border` | `#2a2a2a` | Subtle card borders |
| `--accent` | `#FFD600` | Primary brand yellow |
| `--accent-glow` | `rgba(255,214,0,0.2)` | Glow effects |
| `--text` | `#FFFFFF` | Primary text |
| `--text-muted` | `#888888` | Secondary text |
| `--terminal-bg` | `#1a1a1a` | Code/terminal background |

### Typography
- **Headlines:** Space Grotesk (bold, geometric)
- **Body:** Inter (clean, readable)
- **Code:** JetBrains Mono (monospace)

### Animation
- Scroll-triggered fade-up reveals via IntersectionObserver (no library)
- Terminal auto-typing with blinking cursor
- Yellow glow pulses on key elements
- CSS transitions only — no heavy animation libraries
- Animated dashed lines for data flow visualization

---

## Page Structure

### 1. Nav Bar (sticky)
- Transparent → solid `#0a0a0a` on scroll
- Left: Lightning bolt icon + "AI Happy Design"
- Right: "GitHub" link + star badge, "Get Started" yellow pill button
- Height: 64px

### 2. Hero (100vh)
- **Left (55%):**
  - Yellow pill badge: "Free & Open Source — MIT License"
  - Headline: "Give AI full control of your Figma canvas" (Space Grotesk, ~56px)
  - Subline: "Design, edit, and export — all from natural language. One binary. Local-only. Blazing fast." (Inter, 20px, muted)
  - Buttons: "Get Started" (yellow fill) + "View on GitHub" (outline)
- **Right (45%):**
  - Terminal window with fake chrome (three dots, title bar)
  - Auto-typing animation cycling through 3 demos:
    1. `text.create` command → result JSON
    2. `batch` command → "27 operations in 0.9s"
    3. `export.image` → "Exported to design.png"
  - Yellow glow behind terminal, subtle pulse

### 3. Narrative: "The Gap"
- Centered headline: "Design and code don't have to live in separate worlds"
- Animation: "Design" box + "Code" box slide together and merge with lightning bolt glow
- Below: "AI Happy Design bridges the gap. Type a command. Get a design. It's that simple."

### 4. How It Works (3-step)
- Headline: "How it works"
- Three steps connected by animated dashed yellow lines:
  1. **Prompt** — chat bubble — "Tell your AI what to design"
  2. **Relay** — lightning bolt — "Commands flow through a local WebSocket"
  3. **Canvas** — Figma diamond — "Designs appear instantly in Figma"

### 5. Features (2x3 card grid)
Cards with icon + title + 1-line description:

1. **Design Token Calculator** — "Perfect font sizes, spacing, and padding for any canvas — computed instantly"
2. **Free Space Detector** — "Never overlap frames again. Auto-finds the next open spot on your canvas"
3. **14 Tool Domains** — "Paint, text, shape, layout, effects, components, variables — everything Figma can do"
4. **27 ops/sec Batch** — "Chain 150+ commands in one payload. 10-50x faster than individual MCP calls"
5. **Local-Only** — "Nothing leaves your machine. No cloud. No API keys to Figma. Pure privacy"
6. **Chain With Any Tool** — "Pipe into scripts, CI/CD, or other CLIs. It's just a binary that speaks JSON"

### 6. Speed Comparison
- Headline: "Faster than you think"
- Two animated horizontal bars:
  - "Standard MCP" — gray bar, ~3 ops/sec
  - "AI Happy Design CLI" — yellow bar, ~27 ops/sec, "9x faster" badge
- "150 Figma operations in under 6 seconds. One WebSocket. Zero overhead."

### 7. Developer Experience
- Headline: "Dead simple to use"
- Three tabbed code blocks: CLI / MCP Config / Batch
- "Works with Claude Code, Cursor, Windsurf, VS Code, Zed — and any AI that speaks MCP"

### 8. Integrations
- Horizontal row of editor/AI monochrome icons → yellow on hover
- Claude, ChatGPT, Cursor, VS Code, Windsurf, Zed
- "One binary. Every AI editor."

### 9. CTA Section
- Headline: "Start designing with AI"
- "Free forever. Open source. MIT licensed."
- Buttons: "Get Started on GitHub" (yellow) + "Read the Docs" (outline)
- MIT License badge, GitHub stars badge
- "Follow @workandpromise on X" + "info@workandpromise.com"

### 10. Footer
- Left: "Made by Ashraf Ali"
- Right: "MIT License · GitHub · X"

---

## Key Messaging
- **Primary:** Free, open source, MIT licensed
- **Speed:** 9x faster than standard MCP (27 ops/sec batch)
- **Simplicity:** One binary, no setup, local-only
- **Capability:** 14 tool domains covering everything Figma can do
- **Workflow:** Chain with any tool, pipe into scripts, plug into CI/CD
- **Intelligence:** Built-in design tokens, free space detection, design guidance
- **Privacy:** Nothing leaves your machine

## Social & Contact
- X: @workandpromise
- Email: info@workandpromise.com
- GitHub: nerveband/ai-happy-design
- Creator: Ashraf Ali

## Technical Decisions
- Single HTML file with inline CSS and JS (or separate files, Netlify-ready)
- No build tools, no frameworks, no dependencies
- Google Fonts loaded from CDN (Space Grotesk, Inter, JetBrains Mono)
- IntersectionObserver for scroll animations
- CSS custom properties for theming
- Responsive: mobile-first, works on all screen sizes
- Private GitHub repo for hosting source
- Deploy via Netlify CLI
