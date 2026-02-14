<p align="center">
  <img src="resources/AI Happy Design Icon.png" alt="AI Happy Design" width="128" height="128">
</p>

# AI Happy Design

**Give AI full control of your Figma canvas.**

A single Go binary that connects any AI (Claude, GPT, Gemini, etc.) to Figma through a local WebSocket relay. Design, edit, and export — all from natural language.

Made by [Ashraf Ali](https://ashrafali.net) | [GitHub](https://github.com/nerveband/ai-happy-design) | MIT License

---

## Architecture

```
AI / LLM              CLI (scripting)
   |                    |
   v                    v
[MCP Server] -----> [WebSocket Relay] <-----> [Figma Plugin]
                    localhost:3055
```

One binary. Three modes:
- **MCP server** — standard integration for AI editors (Claude Code, Cursor, Windsurf, etc.)
- **CLI** — direct commands and batch payloads. **10-50x faster than MCP for bulk operations** because it sends a single batch payload over one WebSocket connection instead of individual MCP tool calls. AI agents can use the CLI for heavy lifting.
- **Relay** — WebSocket bridge to the Figma plugin

## What's Included

- **14 tool domains**: paint, shape, text, layout, node, layer, component, style, variable, effect, boolean, page, document, export
- **Batch operations**: Chain 150+ commands in one payload with step interpolation (`${{steps.name.result.id}}`), per-operation timing stats, ~27 ops/sec
- **Smart canvas placement**: `document.find_free_space` scans existing frames and returns exact coordinates — AI never guesses where to place new designs
- **Design tokens**: `design.compute_tokens` auto-computes fonts, spacing, padding, and layout for any canvas size — no Figma connection required
- **Design intelligence**: Built-in design guide the AI can call to learn typography, balance, visual hierarchy, and CSS-to-Figma translation
- **Auto-registration**: `ai-happy-design register` detects installed AI editors and configures MCP automatically
- **Unicode-safe**: Full support for emojis, CJK, accented characters, em-dashes, and special characters
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

# Linux
curl -L https://github.com/nerveband/ai-happy-design/releases/latest/download/ai-happy-design_Linux_x86_64.tar.gz | tar xz
sudo mv ai-happy-design /usr/local/bin/
```

> **Building from source?** Run `codesign -s - ai-happy-design` on macOS after building. Release binaries are pre-signed.

### 2. Setup Plugin

```bash
ai-happy-design setup
```

This extracts the embedded Figma plugin and opens Finder to the manifest file.

1. Figma desktop → **Plugins > Development > Import plugin from manifest**
2. Select the revealed `manifest.json`
3. Run: **Plugins > AI Happy Design**

### 3. Register with Your AI Editor

```bash
ai-happy-design register
```

Auto-detects installed editors (Claude Code, Claude Desktop, Cursor, Windsurf, VS Code, Zed) and registers the MCP server. Or register with a specific editor:

```bash
ai-happy-design register --editor "Claude Code"
ai-happy-design register --editor "Cursor"
```

**Manual setup** (if you prefer):

<details>
<summary>Claude Code</summary>

```bash
claude mcp add ai-happy-design -- ai-happy-design mcp
```

Or add to `~/.claude.json`:
```json
{
  "mcpServers": {
    "ai-happy-design": {
      "type": "stdio",
      "command": "ai-happy-design",
      "args": ["mcp"]
    }
  }
}
```
</details>

<details>
<summary>Claude Desktop / Cursor / Windsurf</summary>

Add to your MCP config file:
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
</details>

Restart your editor after registering. The MCP server auto-detects if a relay is already running and connects as a client — no need to start the relay separately.

### 4. Verify

```bash
ai-happy-design tools --json       # List all tools
ai-happy-design command document.get_info   # Test connection
```

### 5. Upgrade

```bash
ai-happy-design upgrade
```

Auto-checks for updates daily and notifies you. Upgrade downloads and installs the new binary in-place.

## How It Works

### First Prompt

Tell the AI to discover tools first:

> "Call describe(action='catalog') to see all available Figma design tools, then design an Instagram post about [topic]."

### AI Self-Discovery

The binary teaches itself to the AI through two endpoints:

| Command | What it returns |
|---------|----------------|
| `describe(action="catalog")` | Full tool catalog with examples, batch patterns, design playbook |
| `describe(action="design_guide")` | CSS-to-Figma mappings, visual hierarchy, balance rules, typography scales |

### Design Workflow

1. **Discover** — AI calls `describe(action="catalog")` to learn all tools
2. **Compute tokens** — `design.compute_tokens` returns exact font sizes, spacing, padding for the target canvas
3. **Find space** — `document.find_free_space` returns exact x,y coordinates for new frames (never overlaps existing work)
4. **Create** — Batch create via `bulk.execute` with named steps and interpolation
5. **Verify** — `export.image` at scale=2 for crisp output

## CLI Usage

The CLI is the fastest way to drive Figma operations — **recommended for AI agents doing bulk work**.

### Single command
```bash
ai-happy-design command paint.set_solid -p '{"nodeId":"1:2","color":"#2563EB"}'
```

### Batch operations (fastest)
```bash
# From file — send 100+ operations in one shot:
ai-happy-design batch -f payload.json

# Inline:
ai-happy-design batch -o '[{"name":"card","command":"node.create_frame","params":{"width":320,"height":200}}]'
```

### List available actions
```bash
ai-happy-design actions           # All domains and actions
ai-happy-design actions document  # Just document actions
ai-happy-design tools --llm       # Full LLM-focused catalog with examples
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
4. Relay's active channel (auto-detected)

## Export Formats

Export any node as PNG, JPG, SVG, or PDF:

```bash
ai-happy-design command export.image '{"nodeId":"47:765","format":"PNG","scale":2}' -o output.png
ai-happy-design command export.image '{"nodeId":"47:765","format":"JPG","scale":1}' -o output.jpg
ai-happy-design command export.svg '{"nodeId":"47:765"}' -o output.svg
ai-happy-design command export.pdf '{"nodeId":"47:765"}' -o output.pdf
```

SVG exports are saved as raw SVG text. PNG/JPG/PDF are decoded from base64. The `-o` flag specifies the output path; without it, a filename is auto-generated from the node name.

## Performance

Benchmarked with a 5-slide carousel (108 operations):

| Metric | Value |
|--------|-------|
| Total time | ~10s |
| Operations/sec | ~11 |
| Avg per operation | ~83ms |
| Connection overhead | ~500ms (one-time) |

CLI batch sends all operations over one WebSocket connection. MCP sends each tool call individually. For designs with 20+ operations, CLI batch is significantly faster.

Batch output includes per-operation `elapsedMs` and a `timing` summary with `totalMs`, `avgMs`, and `opsPerSec`.

## Do's and Don'ts

**Do:**
- Call `describe(action="catalog")` before building commands
- Call `document.find_free_space` before creating new frames — it returns exact coordinates
- Call `design.compute_tokens` to get sizing for your canvas
- Use batch (`bulk.execute`) for multi-element designs (3+ elements)
- Use **absolute x/y positioning** for main layout + auto-layout only for badges/buttons
- Always pass `lineHeightUnit: "PERCENT"` with `text.set_line_height` (default is PIXELS)
- Export at scale=2 for crisp output
- Name every element descriptively

**Don't:**
- Guess frame positions — always use `find_free_space`
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
| Binary killed on macOS (source build) | Run `codesign -s - /path/to/ai-happy-design` |
| MCP not loading in editor | Run `ai-happy-design register --force` to re-register |

## Prerequisites

- Figma desktop app
- Go 1.21+ (only for building from source)
- Node.js 18+ (only for building from source)

## All Commands

| Command | Description |
|---------|-------------|
| `ai-happy-design setup` | Extract and install the Figma plugin |
| `ai-happy-design register` | Auto-register MCP with detected AI editors |
| `ai-happy-design mcp` | Start MCP server (stdio transport, for AI editors) |
| `ai-happy-design ws` | Start WebSocket relay only |
| `ai-happy-design command <action>` | Execute a single command |
| `ai-happy-design batch` | Execute a batch of commands |
| `ai-happy-design actions` | List all domain.action pairs |
| `ai-happy-design tools` | Print tool catalog |
| `ai-happy-design relay start/stop/status/logs` | Manage the relay process |
| `ai-happy-design upgrade` | Upgrade to the latest version |

## Documentation

See `docs/` for detailed guides:
- [Getting Started](docs/getting-started.md)
- [LLM Integration](docs/llm-integration.md)
- [CLI Batch & Payloads](docs/cli-batch-and-payloads.md)
- [Image Fill Contract](docs/image-fill-contract.md)
- [Troubleshooting](docs/troubleshooting.md)

## License

MIT — see [LICENSE](LICENSE).
