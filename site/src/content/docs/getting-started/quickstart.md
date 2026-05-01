---
title: Quick Start
description: Create, validate, audit, and export a Figma design with ahd-figma.
---

Start with schema discovery, validate before writing, run the batch, then audit and export.

## 1. Discover a Command

```bash
ahd-figma schema text.create --json
ahd-figma schema layout.pricing_grid --json
```

## 2. Validate a Payload

```bash
ahd-figma validate docs/examples/live-acceptance-full-parity.json
```

## 3. Run a Batch

```bash
ahd-figma batch docs/examples/live-acceptance-full-parity.json
```

Batch payloads may be a plain array or an object with an `operations` array.

## 4. Audit Accessibility

```bash
ahd-figma command document.accessibility_audit \
  '{"file":"docs/examples/live-acceptance-full-parity.json"}'
```

The audit checks contrast, line height, target size, non-text contrast, semantic names, color-only communication, and focus ambiguity.

## 5. Export

```bash
ahd-figma command export.image '{"nodeId":"FRAME_ID","scale":2}'
```

For very large frames, use `scale: 1` to avoid oversized exports.
