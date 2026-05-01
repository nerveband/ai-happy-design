# Figma Tool Implementation Notes

Updated: 2026-05-01

## Research Sources

The implementation plan was based on current Figma Plugin API, REST API, MCP guide research through Context7 and official Figma docs, plus local source inspection of representative Figma CLI/MCP tools.

## Implementation Lessons Applied

- From official Figma MCP: keep `get_design_context`-style compact reads and CSS extraction available, but preserve AHD's deterministic CLI/batch path.
- From Figma Context MCP: separate fetching, simplifying, serialization, and metrics. AHD applies this through schema discovery, compact output, and batch summaries.
- From figma-mcp-go: executable MCP tools can be thin wrappers over an existing bridge. AHD now generates MCP tools from `internal/schema`.
- From figma-console-mcp: broad tool coverage needs schema compatibility tests, accessibility checks, and plugin/runtime version awareness. AHD now has parity tests, local accessibility audit, and guarded plugin APIs.
- From TalkToFigma-style bridges: plugin write access is useful but needs stronger contracts. AHD keeps the WebSocket bridge behind validation and routing.
- From figma-cli/design-system CLIs: token export and design-system automation are valuable. AHD now includes local `tokens.export` and keeps variables/token commands schema-visible.
- From RedMadRobot figma-export: production export should be deterministic and config-friendly. AHD's token/export direction follows deterministic snapshots instead of ad-hoc strings.

## Features Added From The Research

- Schema-backed MCP tools and resources.
- `node.get_css` via `getCSSAsync()`.
- Dev Mode focused node read.
- Guarded Figma Draw commands.
- Guarded grid-layout commands.
- Extended variable collection and override helpers.
- Figma REST metadata, oEmbed, Dev Resources, and Webhooks V2 commands.
- Local accessibility audit and parity/token commands.
- CLI `--output-format` and simple `--jq`.

## Guardrails Kept

- The plugin still builds to ES6 for the known desktop runtime compatibility rule in this repo.
- REST is not the primary write path.
- Beta APIs are runtime-guarded.
- No release, tag, or push is part of implementation without explicit user permission.
