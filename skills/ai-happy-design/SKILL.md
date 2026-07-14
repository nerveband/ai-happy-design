---
name: ai-happy-design
description: Design in Figma using the ahd-figma CLI. Use when asked to create, edit, inspect, or export Figma designs — social media posts, cards, layouts, UI mockups, posters, or any visual design. Triggers on "design in Figma", "create a Figma layout", "make an Instagram post", "export from Figma", "Figma design", or any task involving AI Happy Design. Also use when asked to batch-create multi-element Figma compositions. CLI-only — all operations go through `ahd-figma command` / `ahd-figma batch`.
---

# AI Happy Design

Design in Figma via the `ahd-figma` CLI. `ai-happy-design` remains a legacy-compatible alias. The CLI is fully self-discoverable — run the bootstrap commands below to get the complete API catalog and examples.

## RULES — Read Before Doing Anything

1. **ONE batch file.** All operations in a SINGLE JSON array → one `ahd-figma batch /tmp/ops.json` call. Never write separate files per slide.
2. **Composites first.** Use `slide`/`banner` commands before building frames manually. One composite = complete design with auto-sized text.
3. **Batch handles placement.** Don't call `find_free_space` — batch does it automatically.
4. **No extra files.** No markdown copies, no summaries, no explanation files. Just the batch JSON.
5. **Execute literally.** Do exactly what was asked, nothing more. No bonus content.
6. **Write to OS temp dir.** Always write batch JSON to the OS temp directory to avoid path issues.
7. **Separate frames for separate exports.** Carousels, stories, responsive sets, A/B variants — each exportable design is its own top-level frame. Never nest slides inside a wrapper frame. Use one `slide` composite per panel.
8. **Output JSON only.** Generate the batch JSON array and run it. No prose, no explanations.

## Bootstrap — Run These First

```bash
# Schema discovery (list all commands, get exact param schemas)
ahd-figma schema                           # list all commands
ahd-figma schema text.create --json        # exact JSON schema for a command
ahd-figma validate '[...]'                 # dry-run validation (schema + design lint)
ahd-figma guide                            # design intelligence guide

# Full API catalog (commands, params, aliases, design guidance)
ahd-figma tools --llm --json
ahd-figma agent-context --json

# Example batch payloads by category
ahd-figma examples              # list categories
ahd-figma examples carousel     # 3-slide carousel
ahd-figma examples slide        # single slide composite
ahd-figma examples banner       # banner composite
ahd-figma examples effects      # shadows, glass, gradients, masks
ahd-figma examples batch        # raw batch with primitives

# Design tokens for your canvas size
ahd-figma command design.compute_tokens '{"width": 1080, "height": 1080}'
# Includes starter template + tips + alias quick reference for the next batch step
```

## Quick Reference

```bash
# Create designs (batch mode)
ahd-figma batch /tmp/ops.json          # from file (preferred)
ahd-figma batch '[...]'                # inline JSON
ahd-figma batch f1.json f2.json        # multi-file

# Edit existing nodes (single command)
ahd-figma command paint.set_solid '{"nodeId":"1:2","color":"#FF0000"}'

# Audit layout before making repairs (read-only; no screenshot needed)
ahd-figma command layout.audit '{"nodeId":"...","compact":true}'
# Apply one intentional batch, re-audit, then take one final screenshot

# Export & verify
ahd-figma command export.image '{"nodeId":"...","scale":2}'
ahd-figma command document.screenshot '{"nodeId":"...","scale":1}' --output /tmp/shot.png
ahd-figma command verify.visual --payload '{"artifactPath":"/tmp/shot.png","target":"final frame"}'

# Useful flags
--live          # print progress
--fail-fast     # stop at first failure
--parallel      # concurrent multi-file
--allow-overlap # skip auto-placement
--no-lint       # disable post-batch lint (lint is on by default)
--stdin         # read command or batch JSON from stdin
--payload-file  # read command params or batch operations from a file
--deliver       # route JSON results to stdout, file:<path>, or dir:<path>
```

Layout repair workflow:

```text
audit → apply one bounded batch → re-audit → screenshot once
```

`layout.audit` reports overflow, clipping, text-fit failures, sibling overlap, tight gaps, and manual-layout risks with measured evidence and suggested CLI fixes. Do not guess at x/y, font size, or box dimensions when an audit finding provides a measured delta. Use `compact:true` to reduce tokens. JSON output is quiet by default; parse stdout without merging stderr.

Proof gates:

```bash
ahd-figma doctor --json
ahd-figma verify plugin
ahd-figma verify syntax
ahd-figma verify live
ahd-figma verify release
```

Batch mode resolves semantic token aliases automatically:
- `sz:"hero|title|body|caption|..."`
- `padding:"side|frame|card"`
- `itemSpacing/gap:"section|card|item|tight"`
- `r:"card|button|pill"`
- `w:"content"`
