---
title: Batch Operations
description: Build complete Figma designs with one validated payload.
---

Batch mode is the main path for LLM agents. It reduces round trips and keeps design generation inspectable.

## Shape

```json
[
  {
    "name": "card",
    "command": "node.create_frame",
    "params": { "name": "Card", "width": 400, "height": 300 }
  },
  {
    "name": "title",
    "command": "text.create",
    "params": {
      "parentId": "${{steps.card.result.id}}",
      "text": "Hello World",
      "fontSize": 24
    }
  }
]
```

Object form is also accepted:

```json
{ "operations": [] }
```

## Agent Rules

- Name every layer semantically.
- Use `design.compute_tokens` before laying out a design.
- Use `document.find_free_space` before placing new frames.
- Use dry-run validation before writes.
- Export and inspect after batch runs.

## Common Aliases

Compact command aliases are supported for batch payloads: `frame`, `rect`, `text`, `fill`, `stroke`, `gradient`, `shadow`, `blur`, `glass`, `modify`, `mask`, and `find`.
