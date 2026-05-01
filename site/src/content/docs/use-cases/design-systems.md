---
title: Design Systems
description: Use variables, tokens, components, and parity checks for design system work.
---

Use AI Happy Design to extract tokens, inspect components, build examples, and check consistency.

```bash
ahd-figma command variable.get_all
ahd-figma command tokens.export '{"format":"css","outputPath":"tokens.css"}'
ahd-figma command component.get_local
ahd-figma command parity.compare_code '{"specPath":"design-system-contract.json"}'
```
