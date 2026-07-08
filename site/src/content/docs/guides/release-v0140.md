---
title: Release v0.14.0
description: AI Happy Design v0.14.0 release notes.
---

`v0.14.0` is the roadmap release for agent-native Figma automation.

## Highlights

- 222 schema-backed command surfaces.
- MCP `prompts/list` and `prompts/get`, backed by the design catalog.
- Screenshot proof loops: `document.screenshot`, `document.screenshot_selection`, and `verify.visual`.
- Structured stderr errors with `code`, `message`, `hint`, and `retryable`.
- `agent-context --json` for compact machine-readable onboarding.
- Durable local jobs and artifact routing with `--deliver`.
- `doctor` and `verify` proof gates for plugin, syntax, live, and release checks.
- Token presets, design-system health, component-set analysis, and component parity audits.
- Guarded current Figma API coverage for Motion, Shaders, Slots, grid updates, noise vectors, Slides, and FigJam.
- Generated showcase artifacts for social square, story poster, and event poster formats.

## Install

```bash
ahd-figma upgrade
```

For a clean local build from source:

```bash
make deploy
```

Reload the Figma development plugin after upgrading so it picks up the bundled `code.js`.
