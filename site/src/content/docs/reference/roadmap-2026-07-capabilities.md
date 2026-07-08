---
title: 2026-07 Roadmap Capabilities
description: Implemented v0.14.0 roadmap surface for agent-native Figma automation.
---

## Agent CLI

- `agent-context --json`
- `command --stdin`
- `command --payload`
- `command --payload-file`
- `command --fields`
- `command --deliver stdout|file:<path>|dir:<path>`
- `batch --stdin`
- `batch --payload`
- `batch --payload-file`
- `batch --dry-run`
- `batch --compact`
- structured stderr errors: `code`, `message`, `hint`, `retryable`

## Proof Gates

- `doctor --json`
- `verify plugin`
- `verify syntax`
- `verify live`
- `verify release`
- `document.screenshot`
- `document.screenshot_selection`
- `verify.visual`

## Reliability

- schema safety metadata enforcement
- contract drift gate: `make verify-contracts`
- MCP prompts: `prompts/list`, `prompts/get`
- central plugin `figma.commitUndo()` with `commitUndo:false`
- durable local jobs: `jobs list|get|resume|cancel`
- profile/config inspection: `profile list|use|inspect --redacted`, `config sources`

## Current Figma API Coverage

- Motion guards
- Shader guards
- Slots
- Grid row/column reorder and auto-flow readback
- Noise vectors: `noiseSize`, `noiseSizeX`, `noiseSizeY`, `noiseSizeVector`
- Slides MVP guards
- FigJam MVP guards

## Design-System Moat

- `tokens.preset_tailwind`
- `tokens.preset_shadcn`
- `tokens.preset_material`
- `tokens.setup_system`
- `design_system.health`
- `component.analyze_set`
- `component.arrange_set`
- `parity.audit_component`

## Packaging

- `packaging.generate_skills`
- npm wrapper investigation
- cloud relay proposal
- social/poster showcase artifacts
