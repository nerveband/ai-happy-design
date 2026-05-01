# Site Update Runbook

Updated: 2026-05-01

## Source Of Truth

This repository currently contains `site/dist` and `site/node_modules`, but no obvious website source directory such as `site/src`, `site/content`, or `site/pages`.

Until the real website source is restored or identified, the checked-in generated assets are updated directly for local documentation parity:

- `site/dist/llms.txt`
- `site/dist/llms-full.txt`
- `site/dist/index.html`

## Required Website Content

The public website should mention:

- `ahd-figma` as the preferred command.
- `ai-happy-design` as compatibility alias.
- schema-backed MCP tools.
- 184 schema commands plus `ahd_describe`.
- Dev Mode focused node and `node.get_css`.
- Figma Draw guarded commands.
- extended variable commands.
- REST metadata, Dev Resources, and Webhooks.
- live Figma acceptance workflow.

## Production Update

Do not claim production website deployment from this repo alone. Once the real source is located, update source content there, rebuild, and then replace generated `site/dist` files from that source build.
