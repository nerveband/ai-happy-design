---
title: Multi-Agent Workflows
description: Coordinate research, design generation, validation, and live Figma verification.
---

AI Happy Design is built for agent loops:

1. Research agent produces a design brief.
2. Layout agent creates a batch payload.
3. Validator agent runs schema validation and accessibility checks.
4. Figma agent executes the batch and exports screenshots.
5. Reviewer agent compares output against the reference and iterates.

Use `--output-format jsonl`, `--jq`, and compact fields to keep these loops predictable.
