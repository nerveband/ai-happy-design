# Live Figma Acceptance

Updated: 2026-05-01

## Acceptance Scope

The live gate verifies that the rebuilt binary, managed relay, and local Figma plugin can execute real canvas reads/writes. It is intentionally separate from REST tests, which use mocked HTTP.

## Commands

```bash
make build-go
make sync-plugin
make deploy
ahd-figma relay status
ahd-figma validate docs/examples/live-acceptance-full-parity.json
ahd-figma batch -f docs/examples/live-acceptance-full-parity.json --strict-quality --live
```

Then verify representative read paths:

```bash
ahd-figma command layout.audit '{"nodeId":"<frame-id>","compact":true}'
ahd-figma command node.get_tree '{"nodeId":"<frame-id>","compact":true}'
ahd-figma command node.get_css '{"nodeId":"<frame-id>"}'
ahd-figma command export.image '{"nodeId":"<frame-id>","format":"PNG","scale":1}'
```

MCP live execution gate:

```bash
printf '{"jsonrpc":"2.0","id":1,"method":"tools/list"}\n' | ahd-figma mcp
```

## Latest Run

### HTML/CSS-Like Recreation Gate

The current installed binary also passed the dense-reference recreation gate used to verify the new measured layout primitives:

```bash
ahd-figma validate docs/examples/cair-georgia-gala-css-like-executable.json
ahd-figma batch -f docs/examples/cair-georgia-gala-css-like-executable.json --allow-overlap --no-lint --jq .files.0.summary
ahd-figma command export.image '{"nodeId":"60:661","format":"PNG","scale":1}'
```

Result:

- frame: `60:661` `CAIR-Georgia Gala Partnership Packages - CSS-like AHD Recreation`
- 14/14 batch operations succeeded
- `text.fit_box` selected a fitting 30px business-heading size
- both `layout.pricing_grid` calls succeeded
- export: `/var/folders/g4/t50hjvlj7dj9b70npwg8h3tw0000gn/T//ahd-export-CAIR-Georgia-Gala-Partnership-Packages---CSS-like-AHD-Recreation-1777666419.png`

Defects found and fixed in this loop:

- `layout.pricing_grid` rejected stroke colors because alpha was placed inside the Figma `color` object; the handler now sends alpha as paint opacity.
- `layout.pricing_grid` ignored `heading` and only rendered `tier`; it now accepts `heading`, `title`, or `tier`.
- eligibility text was double-prefixed; the handler now respects text already starting with `Eligibility:`.
- `--jq` could not descend into arrays; simple numeric path segments now work.

Visual result: materially closer than the low-level absolute-position version, with rendered tier headings, non-colliding title/rule, and measured card text. Remaining fidelity work is aesthetic, not command/runtime blocking.

### Full-Parity Smoke Gate

The current full-parity live gate created and verified a real frame through the CLI/plugin bridge:

- `57:456` `AHD Full Parity Acceptance`
- 5/5 batch operations succeeded
- strict quality gate passed with 0 lint warnings/errors
- `node.get_css` returned Figma CSS for the frame
- PNG export rendered correctly at `/var/folders/g4/t50hjvlj7dj9b70npwg8h3tw0000gn/T//ahd-export-AHD-Full-Parity-Acceptance-1777663514.png`

Two defects were found and fixed during the loop:

- `Inter SemiBold` failed to load in the live Figma runtime because available font style names are reported as spaced names such as `Semi Bold`; the schema now accepts both spaced Figma names and compact aliases, and the acceptance payload uses `Inter Bold`.
- A badge label was visually misplaced because it was a sibling in an auto-layout frame; the badge is now its own child frame containing the text.

Font availability can be checked against the live Figma runtime before creating text:

```bash
ahd-figma command text.list_fonts '{"fontFamily":"Inter"}'
```

## Defect Loop

If a live command fails:

1. record the exact command and error here
2. fix schema, routing, plugin handler, or docs
3. rerun local tests
4. rerun this live gate
