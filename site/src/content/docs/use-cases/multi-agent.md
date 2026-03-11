---
title: Multi-Agent Workflows
description: Integrate AI Happy Design with Claude Code, CI/CD pipelines, data-driven generation, and cross-tool AI workflows.
---

AI Happy Design is built for automation. It exposes two integration surfaces -- MCP (for AI editors) and CLI (for scripts and pipelines). This makes it a natural fit for multi-agent and automated design workflows.

## Claude Code + MCP Integration

AI Happy Design runs as an MCP server inside Claude Code and other MCP-capable editors.

### Supported Editors

| Editor | Integration | Setup |
|--------|-----------|-------|
| **Claude Code** | MCP (native) | `ai-happy-design register` |
| **Claude Desktop** | MCP (config file) | `ai-happy-design register` |
| **Cursor** | MCP | `ai-happy-design register` |
| **Windsurf** | MCP | `ai-happy-design register` |
| **VS Code** | MCP (via Copilot) | `ai-happy-design register` |
| **Zed** | MCP (JSONC config) | `ai-happy-design register` |
| **Any CLI agent** | CLI | Direct commands |

### Auto-Registration

Run `ai-happy-design register` and it detects all installed editors and configures MCP automatically:

```bash
ai-happy-design register
```

Or add manually to your MCP config:

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

### How It Works

When configured as an MCP server, AI Happy Design exposes its 150+ commands as MCP tools. The AI agent can call any tool directly:

```
Agent: "Create a 1080x1920 social media post with the headline 'Launch Day'"

-> Claude Code calls: design.compute_tokens {width:1080, height:1920}
-> Claude Code calls: document.find_free_space {width:1080, height:1920}
-> Claude Code calls: node.create_frame {name:"Launch Day", width:1080, height:1920, ...}
-> Claude Code calls: text.create {parentId:"42:248", text:"Launch Day", fontSize:152, ...}
-> Claude Code calls: export.image {nodeId:"42:248", format:"PNG", scale:2}
```

The agent handles the full workflow: compute tokens, find space, create elements, export.

### CLI vs MCP Performance

| Mode | Throughput | Best For |
|------|-----------|----------|
| MCP | ~3-5 ops/sec | Interactive sessions, small edits, exploration |
| CLI batch | ~27 ops/sec | Bulk creation (20+ ops), automated pipelines |
| CLI single | 1 op/call | Quick edits, scripting |

For interactive work (modifying a few elements, inspecting a design), MCP is the right choice. For creating entire layouts with 20+ operations, the agent can switch to CLI batch mode for better throughput.

### Agent CLI Fallback

An MCP-connected agent can also invoke the CLI directly when batch mode is faster:

```bash
# Agent writes a batch file, then executes it via CLI
ai-happy-design batch /tmp/design-ops.json --lint
```

This hybrid approach gives the agent the flexibility of MCP for exploration and the speed of CLI for creation.

## CI/CD Pipelines

The CLI runs anywhere Go binaries run -- macOS, Linux, Windows. This makes it straightforward to include in CI/CD pipelines.

### Automated Design Generation

Generate designs as part of a build pipeline:

```yaml
# .github/workflows/design.yml
name: Generate Marketing Assets
on:
  push:
    paths: ['content/campaigns/**']

jobs:
  generate:
    runs-on: macos-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install AI Happy Design
        run: |
          curl -sSL https://github.com/nerveband/ai-happy-design/releases/latest/download/ai-happy-design_darwin_arm64.tar.gz | tar xz
          chmod +x ai-happy-design

      - name: Start relay
        run: |
          ./ai-happy-design ws &
          sleep 2

      - name: Generate assets
        run: |
          for file in content/campaigns/*.json; do
            ./ai-happy-design batch "$file" --lint --strict-quality --fail-fast
          done

      - name: Export all frames
        run: |
          ./ai-happy-design command export.batch \
            '{"nodeIds":"1:2,1:3,1:4","format":"PNG","scale":2}'
```

### Validation in Pre-Commit Hooks

Use `validate` mode for dry-run checks without connecting to Figma:

```bash
#!/bin/bash
# .git/hooks/pre-commit
for file in $(git diff --cached --name-only -- '*.json'); do
  if head -1 "$file" | grep -q '^\['; then
    ai-happy-design validate -f "$file" || exit 1
  fi
done
```

This catches malformed batch files before they reach the repository.

### Schema Validation in PR Checks

```bash
# Validate batch JSON syntax and schema
ai-happy-design validate -f design.json

# Check schema for a specific command
ai-happy-design schema text.create --json
```

## Automated Design Generation from Data

Generate designs programmatically from structured data -- product catalogs, analytics dashboards, report cards.

### From a JSON API

Fetch data from an API and pipe it through a template generator:

```bash
# Fetch analytics data and generate a dashboard
curl -s https://api.example.com/metrics | \
  python3 generate_dashboard.py | \
  ai-happy-design batch --lint
```

The Python script transforms the API response into batch JSON, which AI Happy Design executes directly via stdin.

### Template-Based Generation

Maintain reusable templates with placeholder values:

```json
[
  {
    "name": "card",
    "command": "frame",
    "params": {
      "name": "${TITLE}",
      "w": 400, "h": 300,
      "bg": "${BG_COLOR}",
      "r": 16,
      "layoutMode": "VERTICAL",
      "padding": 24,
      "itemSpacing": 12
    }
  },
  {
    "name": "heading",
    "command": "text",
    "params": {
      "pid": "${{steps.card.result.id}}",
      "text": "${TITLE}",
      "sz": 28,
      "ff": "Inter",
      "fontStyle": "Bold",
      "color": "${TEXT_COLOR}"
    }
  }
]
```

Fill in the values with `envsubst`:

```bash
TITLE="Q1 Revenue" BG_COLOR="#0C1E2C" TEXT_COLOR="#FFFFFF" \
  envsubst < template.json | ai-happy-design batch --lint
```

### Batch File Generation from Code

An AI agent or script can generate batch JSON programmatically:

```python
import json

def create_card_batch(title, subtitle, color):
    return [
        {
            "name": "card",
            "command": "frame",
            "params": {
                "name": title,
                "w": 400, "h": 300,
                "bg": color,
                "r": 16,
                "layoutMode": "VERTICAL",
                "padding": 24,
                "itemSpacing": 12
            }
        },
        {
            "name": "card_title",
            "command": "text",
            "params": {
                "pid": "${{steps.card.result.id}}",
                "text": title,
                "sz": 24,
                "ff": "Inter",
                "fontStyle": "Bold",
                "color": "#FFFFFF"
            }
        },
        {
            "name": "card_subtitle",
            "command": "text",
            "params": {
                "pid": "${{steps.card.result.id}}",
                "text": subtitle,
                "sz": 16,
                "ff": "Inter",
                "color": "#888888"
            }
        }
    ]

ops = create_card_batch("Revenue", "$2.4M this quarter", "#0C1E2C")
with open("/tmp/card-ops.json", "w") as f:
    json.dump(ops, f)
```

Then execute:

```bash
ai-happy-design batch /tmp/card-ops.json --lint
```

## Cross-Tool Workflows

Combine AI Happy Design with other AI tools for end-to-end content production.

### Content Generation + Design

1. **Generate copy** with an LLM (ChatGPT, Claude, etc.)
2. **Generate images** with an image model (DALL-E, Midjourney, Stable Diffusion)
3. **Compose the design** with AI Happy Design

```bash
# Agent writes content, then creates the design
ai-happy-design batch - <<'EOF'
[
  {
    "name": "post",
    "command": "frame",
    "params": {"name": "Social Post", "w": 1080, "h": 1080, "bg": "#0C1E2C"}
  },
  {
    "name": "title",
    "command": "text",
    "params": {
      "pid": "${{steps.post.result.id}}",
      "text": "Ship Faster with AI",
      "sz": 84, "ff": "Space Grotesk", "fontStyle": "Bold",
      "color": "#FFFFFF", "x": 80, "y": 300, "w": 920
    }
  },
  {
    "name": "body",
    "command": "text",
    "params": {
      "pid": "${{steps.post.result.id}}",
      "text": "Automate your design workflow and ship in hours, not weeks.",
      "sz": 36, "ff": "Inter",
      "color": "#888888", "x": 80, "y": 440, "w": 920
    }
  }
]
EOF
```

### Design Review with LLMs

Use `node.get_tree` to extract a design's structure for LLM review:

```bash
# Extract the design tree (compact for token efficiency)
ai-happy-design command node.get_tree \
  '{"nodeId":"1:2","compact":true}' > /tmp/design-tree.json

# Run lint checks
ai-happy-design command document.lint '{"nodeId":"1:2"}'

# Analyze the design system
ai-happy-design command design_system.analyze
```

Feed the output to an LLM for accessibility review, visual hierarchy analysis, or design system consistency checks.

### HTML/CSS to Figma

Extract design operations from HTML/CSS files:

```bash
ai-happy-design extract component.html --width 1080 --height 1350
```

This parses the HTML/CSS structure and generates batch JSON operations.

## Architecture for Multi-Agent Systems

AI Happy Design fits into multi-agent architectures as a specialized design executor:

```
Orchestrator (Claude, GPT, etc.)
    |
    +-- Content Agent --> Copy, headlines, body text
    |
    +-- Image Agent   --> Generated images (DALL-E, etc.)
    |
    +-- Design Agent  --> Figma via AI Happy Design
         |                (batch + export)
         +-- MCP for exploration
         +-- CLI batch for bulk creation
```

The design agent receives content and images from other agents, composes them in Figma using AI Happy Design, and exports the final assets.

## Channel Isolation

When multiple editors connect to the same relay, use channel keys to isolate sessions:

```bash
# Set channel in environment
export AHD_CHANNEL=my-project

# Or per-command
ai-happy-design command document.get_info --channel my-project
```

The relay runs once on port 3055. All editors connect to the same relay but are isolated by channel.

## Tips for Multi-Agent Workflows

- **Use CLI batch for bulk operations** -- 27 ops/sec vs 3-5 ops/sec for MCP individual calls
- **Use MCP for interactive exploration** -- inspecting designs, making small edits, getting document info
- **Validate before executing** -- `ai-happy-design validate -f ops.json` catches schema errors without touching Figma
- **Use `--fail-fast`** in CI pipelines to stop on first error rather than continuing through failures
- **Use `--strict-quality`** in pipelines to enforce design quality gates
- **Export to temp dir** -- exports auto-save to `os.TempDir()`, making them easy for other tools to pick up
- **Use `compact: true`** on `node.get_tree` to minimize token usage when feeding design structure to LLMs
- **Separate concerns** -- let one agent generate content, another generate images, and AI Happy Design handle composition
- **Channel keys prevent cross-talk** between different editor sessions
- **Use `ai-happy-design register`** to auto-configure all editors at once

---

Made by [Ashraf Ali](https://ashrafali.net) | [GitHub](https://github.com/nerveband/ai-happy-design) | License: GPL-3.0
