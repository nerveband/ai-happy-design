# Getting Started

## 0. Prebuild Checklist
Run this first for clean setup and environment validation:
- `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/docs/prebuild-checklist.md`

## 1. Prerequisites
- Go 1.22+
- Node.js 18+
- npm
- Figma Desktop (for local plugin development)

## 2. Build the CLI/MCP Binary
From repo root:

```bash
make build
```

Binary output:
- `./bin/ai-happy-design`

## 3. Build the Figma Plugin UI/Code

```bash
cd plugin
npm install
npm run build
cd ..
```

Build outputs:
- `plugin/dist/code.js`
- `plugin/dist/ui.html`

## 4. Load Plugin in Figma
1. Open Figma Desktop.
2. Go to plugin development mode.
3. Import plugin manifest from:
- `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin/manifest.json`
4. Run the plugin: **AI Happy Design**.

## 5. Start Relay + MCP or Relay Only
### MCP mode (recommended for LLM integration)

```bash
./bin/ai-happy-design mcp
```

This starts:
- MCP stdio server
- embedded WebSocket relay (default port `3055`)

### Relay-only mode (for manual CLI testing)

```bash
./bin/ai-happy-design ws
```

### One-stop CLI mode (no manual relay start needed)
For direct CLI usage, relay auto-starts on demand for:
- `command`
- `batch`
- `connect`

Disable auto-start only when needed:

```bash
./bin/ai-happy-design --no-auto-relay command document.get_info
```

Relay lifecycle commands:

```bash
./bin/ai-happy-design relay start
./bin/ai-happy-design relay status
./bin/ai-happy-design relay logs --lines 80
./bin/ai-happy-design relay stop
```

Optional persistent relay on macOS:

```bash
./bin/ai-happy-design relay install-agent
```

## 6. Connect Plugin
Open plugin UI in Figma.
- Channel is persisted automatically.
- Auto-connect is enabled by default.
- Relay URL is editable (default `ws://localhost:3055/ws`).
- If needed, click **Connect**.

## 7. Smoke Test with CLI
In another terminal:

```bash
./bin/ai-happy-design tools --json
./bin/ai-happy-design command document.get_info
```

If only one active plugin channel exists, channel arg is optional.

You can still pass an explicit channel:

```bash
./bin/ai-happy-design command happy-unicorn-42 document.get_info
```
