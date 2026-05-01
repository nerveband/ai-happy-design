---
title: LLM Files
description: Machine-readable docs for AI agents.
---

AI Happy Design publishes two LLM-oriented references:

- [`/llms.txt`](/llms.txt) - compact workflow and feature summary.
- [`/llms-full.txt`](/llms-full.txt) - generated full command reference from `ahd-figma schema --all`.

Regenerate after schema changes:

```bash
ahd-figma schema --all > site/public/llms-full.txt
```
