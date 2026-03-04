---
name: ai-happy-design
description: Design in Figma using the ai-happy-design CLI. Use when asked to create, edit, or export Figma designs — social media posts, cards, layouts, UI mockups, posters, or any visual design. Triggers on "design in Figma", "create a Figma layout", "make an Instagram post", "export from Figma", "Figma design", or any task involving ai-happy-design. Also use when asked to batch-create multi-element Figma compositions. CLI-only — all operations go through `ai-happy-design command` / `ai-happy-design batch`.
---

# AI Happy Design

Design in Figma via the `ai-happy-design` CLI. The CLI is fully self-discoverable — run the bootstrap commands below to get the complete API catalog and examples.

## RULES — Read Before Doing Anything

1. **ONE batch file.** All operations in a SINGLE JSON array → one `ai-happy-design batch /tmp/ops.json` call. Never write separate files per slide.
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
ai-happy-design schema                           # list all commands
ai-happy-design schema text.create --json        # exact JSON schema for a command
ai-happy-design validate '[...]'                 # dry-run validation (schema + design lint)
ai-happy-design guide                            # design intelligence guide

# Full API catalog (commands, params, aliases, design guidance)
ai-happy-design tools --llm --json

# Example batch payloads by category
ai-happy-design examples              # list categories
ai-happy-design examples carousel     # 3-slide carousel
ai-happy-design examples slide        # single slide composite
ai-happy-design examples banner       # banner composite
ai-happy-design examples effects      # shadows, glass, gradients, masks
ai-happy-design examples batch        # raw batch with primitives

# Design tokens for your canvas size
ai-happy-design command design.compute_tokens '{"width": 1080, "height": 1080}'
# Includes starter template + tips + alias quick reference for the next batch step
```

## Quick Reference

```bash
# Create designs (batch mode)
ai-happy-design batch /tmp/ops.json          # from file (preferred)
ai-happy-design batch '[...]'                # inline JSON
ai-happy-design batch f1.json f2.json        # multi-file

# Edit existing nodes (single command)
ai-happy-design command paint.set_solid '{"nodeId":"1:2","color":"#FF0000"}'

# Export & verify
ai-happy-design command export.image '{"nodeId":"...","scale":2}'

# Useful flags
--live          # print progress
--fail-fast     # stop at first failure
--parallel      # concurrent multi-file
--allow-overlap # skip auto-placement
--no-lint       # disable post-batch lint (lint is on by default)
```

Batch mode resolves semantic token aliases automatically:
- `sz:"hero|title|body|caption|..."`
- `padding:"side|frame|card"`
- `itemSpacing/gap:"section|card|item|tight"`
- `r:"card|button|pill"`
- `w:"content"`
