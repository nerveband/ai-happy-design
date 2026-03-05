# AHD Illustrator v0.1 Release Notes

## Highlights

- New `ahd-illustrator` CLI with machine-first JSON envelopes
- Shared schema registry and validator for agent discovery and safety
- `tools`, `schema`, `command`, `batch`, `host`, `doctor`, and `examples` surfaces
- `--dry-run` validation flow for every mutating command path
- Batch interpolation engine with strict step-result resolution
- AppleScript host adapter for `do javascript`, app open, and app quit flows
- JSX bridge runtime for raw script execution and plugin selector dispatch
- Live hardening for Illustrator 30.2.1: bundle-id app discovery, version reporting, ES3-safe bridge serialization, and deterministic plugin probing
- Plugin dependency reduction: `inspect.*`, gradient application, and graphic style application now run through pure scripting
- Script-depth improvements: multi-artboard `document.new`, format-aware `document.save_as`, rounded rectangles, artboard-targeted PNG/JPG/SVG export, and `app.info` runtime preset discovery
- Buildable C++ plugin bridge skeleton with `sendScriptMessage` handlers

## Command Coverage

v0.1 ships command domains for:

- `app.*`
- `document.*`
- `artboard.*`
- `layer.*`
- `selection.*`
- `path.*`
- `text.*`
- `appearance.*`
- `action.*`
- `export.*`
- `inspect.*`

## Safety and Agent UX

- Deterministic JSON by default
- Schema introspection via `ahd-illustrator schema`
- Low-risk fuzzy correction for enum-like fields only
- Working-directory path hardening for output paths
- Opaque identifier hardening against URL/query/fragment/encoded traversal syntax
- Command-specific validation for impossible multi-artboard layout payloads
- Tagged macOS integration coverage for the host adapter

## Platform Notes

- Supported platform: macOS
- Scope: CLI-first, no MCP in v0.1
- The native plugin is now optional for the current CLI surface and reserved for future native-only capabilities
- Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.
