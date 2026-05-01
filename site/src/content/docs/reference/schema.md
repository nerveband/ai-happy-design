---
title: Schema System
description: How schema validation, correction, aliases, and design linting work.
---

Every command is schema-backed. Agents should inspect schemas before generating payloads.

```bash
ahd-figma schema <command> --json
ahd-figma validate payload.json
```

The validator supports:

- Required field checks.
- Type validation.
- Enums and constraints.
- Named colors.
- Common alias normalization.
- Fuzzy command correction.
- Design-lint warnings for risky layout and typography.

Validation should run before any large batch operation.
