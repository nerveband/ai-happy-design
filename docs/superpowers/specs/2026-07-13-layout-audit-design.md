# Layout Audit Design

## Goal

Give AI agents deterministic evidence for Figma layout problems so they can make bounded repairs instead of guessing, repeatedly changing tiny values, or taking unnecessary screenshots.

## Scope

Add a read-only `layout.audit` command for a required `nodeId` subtree. The first version will not mutate the document and will not include automatic repair execution.

Example:

```bash
ahd-figma command layout.audit '{"nodeId":"2260:460"}'
```

## Findings

The audit recursively reports actionable issues including:

- `overflow`: children or text extending beyond parent bounds
- `clipping`: text or titles exceeding their usable box
- `overlap`: sibling collisions, including rendered text bounds
- `gap`: suspiciously tight spacing between neighboring elements
- `fixed_text_overflow`: fixed-size text whose natural size exceeds its box
- `manual_layout_risk`: manually positioned content likely to break after edits
- `auto_layout_risk`: inconsistent sizing or children that cannot safely reflow

Each issue contains severity, code, affected node IDs, evidence, message, confidence, and a bounded repair recommendation.

Example shape:

```json
{
  "code": "TEXT_OVERFLOW",
  "severity": "error",
  "nodeIds": ["2260:628"],
  "evidence": {
    "actual": {"height": 136},
    "available": {"height": 85}
  },
  "message": "Text exceeds its container by 51px.",
  "confidence": "high",
  "fix": {
    "strategy": "resize_container_or_reduce_text",
    "commands": []
  }
}
```

## Measurement

1. Prefer Figma `absoluteBoundingBox` for rendered geometry.
2. Fall back to accumulated local coordinates when absolute bounds are unavailable.
3. Compare every child against its parent bounds.
4. Detect sibling intersections using rendered bounds.
5. Inspect text sizing mode, including `textAutoResize: 'NONE'`.
6. Temporarily measure text using a non-persistent clone with natural sizing, remove it before returning, and never leave document changes behind.
7. Use low confidence rather than guessing when exact measurement is unavailable.

## Repair recommendations

The audit remains read-only but emits concrete, minimal repair plans:

- text overflow: increase text height by the measured amount or reduce font size within a safe range
- text clipping: increase width or reduce font size
- sibling overlap: move the later node by the minimum required gap
- parent overflow: resize the parent or recommend reflow
- manual layout risk: recommend auto-layout only when children are structurally compatible
- tight gap: add only the minimum spacing required

The AI loop is:

```text
audit → apply one intentional batch → audit again → take one final screenshot
```

## CLI and discovery

Register `layout.audit` in the schema, catalog, MCP discovery, and LLM guidance. Documentation will state:

- `nodeId` is required
- the command is read-only
- screenshots are not required for diagnosis
- recommendations must be reviewed before applying
- re-audit is required after edits
- compact output contains only actionable findings

## Output modes

Default output is machine-readable JSON. Compact mode minimizes tokens while retaining actionable findings. Detailed mode may include geometry and measurement evidence.

## Testing

Add:

- unit tests for bounds, overlap, gaps, text measurement, and recommendation generation
- mocked plugin tests for Figma node traversal
- a regression fixture matching the childcare/outdoor overflow case
- tests proving temporary measurement nodes are removed
- compact output contract tests
- live verification against a real Figma subtree

No automatic repair mode is included in this version.

## Non-goals

- automatic mutation or repair
- screenshot-driven diagnosis as the primary signal
- arbitrary visual style opinions
- replacing the existing design lint system

The implementation should reuse existing layout overlap and document lint infrastructure where practical, while keeping `layout.audit` focused on current rendered geometry and actionable repair evidence.
