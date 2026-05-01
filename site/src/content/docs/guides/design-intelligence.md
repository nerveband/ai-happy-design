---
title: Design Intelligence
description: Built-in rules for design tokens, accessibility, layout quality, and visual parity.
---

AI Happy Design gives agents guidance they can query at runtime.

## Design Tokens

```bash
ahd-figma command design.compute_tokens '{"width":1440,"height":900}'
```

Use the returned scale for typography, spacing, margins, cards, and content density.

## Token Export

```bash
ahd-figma command tokens.export \
  '{"variablesFile":"figma-vars.json","outputs":{"css":"tokens.css","tailwind":"tokens.tailwind.json","swift":"FigmaTokens.swift"}}'
```

Supported outputs:

- JSON
- CSS custom properties
- Tailwind config fragments
- Swift
- Android resources

## Accessibility Audit

```bash
ahd-figma command document.accessibility_audit '{"file":"ops.json"}'
```

Checks include text contrast, non-text contrast, line height, target size, semantic naming, color-only communication, and focus ambiguity.

## Parity Compare

```bash
ahd-figma command parity.compare_code '{"specPath":"code-spec.json","threshold":0.8}'
```

Use parity checks to catch missing states, token drift, semantic mismatches, and design-code gaps before handoff.
