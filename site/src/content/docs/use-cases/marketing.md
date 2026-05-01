---
title: Marketing
description: Generate banners, social posts, decks, and campaign visuals in batch.
---

Marketing workflows benefit from batch-first Figma generation.

```bash
ahd-figma command design.compute_tokens '{"width":1080,"height":1350}'
ahd-figma batch campaign-carousel.json
ahd-figma command document.accessibility_audit '{"file":"campaign-carousel.json"}'
ahd-figma command export.image '{"nodeId":"FRAME_ID","scale":2}'
```
