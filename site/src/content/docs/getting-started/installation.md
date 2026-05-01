---
title: Installation
description: Install AI Happy Design v0.13.2, including ahd-figma, the Figma plugin, and editor integration.
---

AI Happy Design ships as one Go binary with two command names:

- `ahd-figma` - preferred CLI for Figma automation
- `ai-happy-design` - compatibility command

The current release is `v0.13.2`.

## Install the Binary

Download the latest release for your platform from GitHub:

```bash
curl -LO https://github.com/nerveband/ai-happy-design/releases/latest/download/ai-happy-design_Darwin_arm64.tar.gz
tar xzf ai-happy-design_Darwin_arm64.tar.gz
mkdir -p ~/bin
mv ai-happy-design ~/bin/ai-happy-design
cp ~/bin/ai-happy-design ~/bin/ahd-figma
codesign -s - -f ~/bin/ai-happy-design ~/bin/ahd-figma
```

Verify:

```bash
ahd-figma --version
ahd-figma schema --json | jq length
```

Expected command count: `184`.

## Upgrade

After `v0.13.2`, the built-in upgrader updates the canonical binary and mirrors the `ahd-figma` alias:

```bash
ahd-figma upgrade
```

## Build from Source

```bash
git clone https://github.com/nerveband/ai-happy-design.git
cd ai-happy-design
make build
make deploy
```

`make deploy` rebuilds the plugin bundle, installs `~/bin/ai-happy-design` and `~/bin/ahd-figma`, signs the macOS binary, and restarts the local relay.

## Load the Figma Plugin

1. Open Figma.
2. Go to Plugins -> Development -> Import plugin from manifest.
3. Select `plugin/manifest.json` from the repo checkout.
4. Start or restart the relay with `ahd-figma ws` or `make deploy`.
5. Reopen the plugin after each plugin rebuild.

## Register With AI Editors

Use MCP mode for editors and agents that support MCP:

```bash
ahd-figma mcp
```

Manual MCP config:

```json
{
  "mcpServers": {
    "ahd-figma": {
      "command": "ahd-figma",
      "args": ["mcp"]
    }
  }
}
```

## Verify the Installed Surface

```bash
ahd-figma schema --json | jq length
ahd-figma tools --json | jq '[to_entries[] | .value | length] | add'
```

Expected:

- CLI schemas: `184`
- CLI tools: `184`
- MCP tools/list: `185`, including `ahd_describe`
