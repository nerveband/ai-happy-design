# Getting Started

## 0. Prebuild Checklist
Run this first for clean setup and environment validation:
- `./docs/prebuild-checklist.md`

## 1. Prerequisites
- Go 1.22+
- Node.js 18+
- npm
- Figma Desktop (for local plugin development)

## 2. Build Binary + Embedded Plugin
From repo root:

```bash
make build
```

Outputs:
- `./bin/ai-happy-design` (legacy-compatible name)
- `./bin/ahd-figma`
- `plugin/dist/code.js`
- `plugin/dist/ui.html`
- `internal/plugin/files/manifest.json`
- `internal/plugin/files/dist/code.js`
- `internal/plugin/files/dist/ui.html`

## 3. Plugin-Only Rebuild (optional during plugin iteration)

```bash
cd plugin
npm ci
npm run check
cd ..
```

## 4. Load Plugin in Figma
1. Open Figma Desktop.
2. Go to plugin development mode.
3. Import plugin manifest from:
- `./plugin/manifest.json`
4. Run the plugin: **AI Happy Design**.

## 5. Start Relay + MCP or Relay Only
### MCP mode (recommended for LLM integration)

```bash
./bin/ahd-figma mcp
```

This starts the schema-backed MCP stdio server. Tool and resource discovery come from the same schema registry used by CLI validation and discovery.

### Relay-only mode (for manual CLI testing)

```bash
./bin/ahd-figma ws
```

### One-stop CLI mode (no manual relay start needed)
For direct CLI usage, relay auto-starts on demand for:
- `command`
- `batch`
- `connect`

Disable auto-start only when needed:

```bash
./bin/ahd-figma --no-auto-relay command document.get_info
```

Relay lifecycle commands:

```bash
./bin/ahd-figma relay start
./bin/ahd-figma relay status
./bin/ahd-figma relay logs --lines 80
./bin/ahd-figma relay stop
```

Optional persistent relay on macOS:

```bash
./bin/ahd-figma relay install-agent
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
./bin/ahd-figma tools --json
./bin/ahd-figma command document.get_info
./bin/ahd-figma command document.get_focused_node
```

If only one active plugin channel exists, channel arg is optional.

You can still pass an explicit channel:

```bash
./bin/ahd-figma command happy-unicorn-42 document.get_info
```
