<p align="center">
  <img src="resource/AI Happy Design Icon.png" alt="AI Happy Design" width="128" height="128">
</p>

# AI Happy Design

**Give AI full control of your Figma canvas.**

A single Go binary that connects any MCP-compatible LLM (Claude, GPT, etc.) to Figma through a local WebSocket relay. Design, edit, and export — all from natural language or CLI commands.

Made by [Ashraf Ali](https://ashrafali.net) | [GitHub](https://github.com/nerveband/ai-happy-design) | MIT License

---

## Architecture

```
AI / LLM              CLI
   |                    |
   v                    v
[MCP Server] -----> [WebSocket Relay] <-----> [Figma Plugin]
                    localhost:3055
```

One binary. Three modes: MCP server (for AI), CLI (for humans), relay (WebSocket bridge to Figma).

> **Tip:** CLI mode is faster than MCP — it talks directly to the relay without protocol overhead. Use MCP for AI agent integration; use CLI for scripting and manual testing.

## What's Included

- **14 tool domains**: paint, shape, text, layout, node, layer, component, style, variable, effect, boolean, page, document, export
- **Batch operations**: Chain 150+ commands in one payload with step interpolation (`${{steps.name.result.id}}`), per-operation timing stats, ~27 ops/sec
- **Design tokens**: `design.compute_tokens` auto-computes fonts, spacing, padding, and layout for any canvas size — no Figma connection required
- **Design intelligence**: Built-in design guide the LLM can call to learn typography, balance, visual hierarchy, and CSS-to-Figma translation
- **Unicode-safe**: Full support for emojis, CJK, accented characters, em-dashes, and special characters in frame/element names and text
- **Local only**: Everything on localhost. Nothing leaves your machine.

## Quick Start

### 1. Install

Download the latest binary from [GitHub Releases](https://github.com/nerveband/ai-happy-design/releases):

```bash
# macOS (Apple Silicon)
curl -L https://github.com/nerveband/ai-happy-design/releases/latest/download/ai-happy-design_Darwin_arm64.tar.gz | tar xz
sudo mv ai-happy-design /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/nerveband/ai-happy-design/releases/latest/download/ai-happy-design_Darwin_x86_64.tar.gz | tar xz
sudo mv ai-happy-design /usr/local/bin/

# Or build from source:
make build-plugin   # Build Figma plugin (requires Node.js)
make build          # Sync plugin + build Go binary
```

### 2. Setup Plugin

```bash
ai-happy-design setup
```

This extracts the embedded Figma plugin and opens Finder to the manifest file.

1. Figma desktop → **Plugins > Development > Import plugin from manifest**
2. Select the revealed `manifest.json`
3. Run: **Plugins > AI Happy Design**

### 3. Start

```bash
# For AI/LLM usage (recommended):
ai-happy-design mcp

# For manual CLI testing:
ai-happy-design ws
```

### 4. Verify

```bash
ai-happy-design tools --json       # List all tools
ai-happy-design command document.get_info   # Test connection
```

### 5. Upgrade

```bash
ai-happy-design upgrade
```

The CLI checks for updates automatically (once per day) and shows a hint when a new version is available.

## Connect Your AI

Add to your MCP client config (Claude Desktop, Cursor, etc.):

```json
{
  "mcpServers": {
    "ai-happy-design": {
      "command": "ai-happy-design",
      "args": ["mcp"]
    }
  }
}
```

### First Prompt

Tell the AI to discover tools first:

> "Call describe(action='catalog') to see all available Figma design tools, then design an Instagram post about [topic]."

### How the AI Learns to Design

The binary teaches itself to the LLM through two discovery endpoints:

| Command | What it returns |
|---------|----------------|
| `describe(action="catalog")` | Full tool catalog with examples, batch patterns, design playbook |
| `describe(action="design_guide")` | Focused design guide: CSS-to-Figma mappings, visual hierarchy, balance rules, typography scales |

The LLM calls these on demand to learn how to create professional Figma designs.

## CLI Usage

### Single command
```bash
ai-happy-design command paint.set_solid -p '{"nodeId":"1:2","color":"#2563EB"}'
```

### Batch operations
```bash
ai-happy-design batch -f payload.json
# or inline:
ai-happy-design batch -o '[{"name":"card","command":"node.create_frame","params":{"width":320,"height":200}}]'
```

### Relay management
```bash
ai-happy-design relay start      # Background start
ai-happy-design relay status     # Check connection
ai-happy-design relay logs       # View logs
ai-happy-design relay stop       # Stop
```

### Channel resolution order
1. Positional argument
2. `--channel` flag
3. `AHD_CHANNEL` env var
4. Relay's active channel

## Performance

Benchmarked with a 7-card carousel (162 operations):

| Metric | Value |
|--------|-------|
| Total time | ~6s |
| Operations/sec | ~27 |
| Avg per operation | ~36ms |
| Connection overhead | ~500ms (one-time) |

Batch output includes per-operation `elapsedMs` and a `timing` summary with `totalMs`, `avgMs`, and `opsPerSec`.

**Slowest operations** (Figma API-bound): `paint.set_gradient` (62ms), `paint.set_solid` (50ms), `text.set_color` (43ms)
**Fastest operations**: `layout.set_auto_layout` (12ms), `paint.set_stroke` (24ms)

> **Tip:** The main bottleneck is Figma's plugin API execution time, not network or serialization. Batch more operations per call to amortize connection overhead.

## Do's and Don'ts

**Do:**
- Call `describe(action="catalog")` before building commands
- Call `design(action="compute_tokens", width=1080, height=1920)` to get sizing for your canvas
- Use batch (`bulk.execute`) for multi-element designs (3+ elements)
- Use **absolute x/y positioning** for main layout + auto-layout only for badges/buttons
- Always pass `lineHeightUnit: "PERCENT"` with `text.set_line_height` (default is PIXELS)
- Export at scale=2 for crisp output
- Name every element descriptively

**Don't:**
- Create cards as rectangles with floating text (use frames with children)
- Leave default names like "Frame 47"
- Use text below 16px on 1080px social media canvases
- Mix padding/spacing values across sibling cards
- Use `lineHeight: 140` without `lineHeightUnit: "PERCENT"` (140 pixels != 140%)
- Skip the balance check on sibling elements

## Troubleshooting

| Problem | Fix |
|---------|-----|
| Plugin "Disconnected" | Start relay: `ai-happy-design relay start` |
| Commands fail | Run `describe(action="catalog")` to rediscover tools |
| Wrong port | Edit Relay URL in plugin UI |
| Channel mismatch | Copy key from plugin, pass via `--channel` |
| Plugin won't load | Run `ai-happy-design setup --force` to re-extract |

## Prerequisites

- Go 1.21+ (only for building from source)
- Node.js 18+ (only for building from source)
- Figma desktop app

## Documentation

See `docs/` for detailed guides:
- [Getting Started](docs/getting-started.md)
- [LLM Integration](docs/llm-integration.md)
- [CLI Batch & Payloads](docs/cli-batch-and-payloads.md)
- [Image Fill Contract](docs/image-fill-contract.md)
- [Troubleshooting](docs/troubleshooting.md)

## License

MIT — see [LICENSE](LICENSE).
