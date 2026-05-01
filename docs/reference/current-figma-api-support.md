# Current Figma API Support

Updated: 2026-05-01

`ahd-figma` now exposes 184 schema commands plus the MCP-only `ahd_describe` helper. The preferred path remains CLI and batch execution through the local relay/plugin bridge, with REST reserved for metadata and Dev Mode integration work.

## Plugin API

Supported through the Figma plugin bridge:

- document reads, selection, focused Dev Mode node, styles, search, lint, and zoom
- node create/read/update/delete, compact tree reads, CSS extraction through `getCSSAsync()`
- frames, text, shapes, images, layer ordering, grouping, masks, boolean operations
- auto layout, layout grids, guarded grid-layout container and child placement commands
- paints, gradients, image fills, strokes, effects, beta glass/noise/texture effects
- local variables, modes, binding, unbinding, extended collection guards, and override helpers
- components, instances, component properties, and reset through `removeOverrides()` with fallback
- guarded Figma Draw commands for text paths, transform groups, brushes, variable-width strokes, dynamic strokes, and pattern fills/strokes

Beta or runtime-dependent APIs are guarded in the plugin and return explicit errors when unavailable in the active Figma runtime.

## REST API

Supported as local commands that do not require the relay/plugin:

- `figma.oembed`
- `figma.file_metadata`
- `figma.dev_resources_list`
- `figma.dev_resource_create`
- `figma.dev_resource_update`
- `figma.dev_resource_delete`
- `figma.webhooks_list`
- `figma.webhook_create`
- `figma.webhook_get`
- `figma.webhook_update`
- `figma.webhook_delete`

Token resolution order is explicit token param, `FIGMA_ACCESS_TOKEN`, then `FIGMA_TOKEN`. Tests use mocked HTTP only.

## Local Agent Commands

Local commands do not touch Figma:

- `design.compute_tokens`
- `tokens.export`
- `parity.compare_code`
- `document.accessibility_audit` when passed a batch JSON file

## Verification Gates

- `go test ./...`
- `cd plugin && npm run check`
- `make build-go && make sync-plugin`
- every `docs/examples/*.json` validates
- MCP `tools/list` reports 185 tools
- live Figma acceptance batch runs through the relay/plugin
