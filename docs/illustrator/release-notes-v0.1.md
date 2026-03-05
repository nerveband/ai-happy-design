# AHD Illustrator v0.1 Release Notes

## Highlights

- New `ahd-illustrator` CLI with machine-first JSON envelopes
- Shared schema registry and validator for agent discovery and safety
- `tools`, `schema`, `command`, `batch`, `host`, `doctor`, and `examples` surfaces
- `--dry-run` validation flow for every mutating command path
- Batch interpolation engine with strict step-result resolution
- AppleScript host adapter for `do javascript`, app open, and app quit flows
- JSX bridge runtime for raw script execution and plugin selector dispatch
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
- Tagged macOS integration coverage for the host adapter

## Platform Notes

- Supported platform: macOS
- Scope: CLI-first, no MCP in v0.1
- Adobe Illustrator is a trademark of Adobe. This project is unaffiliated.
