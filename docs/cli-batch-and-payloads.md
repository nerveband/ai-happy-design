# CLI Batch and Payload Strategy

## Single Operation
Use `command` for one action.

```bash
./bin/ahd-figma command paint.set_solid -p '{"nodeId":"1:2","color":"#2563EB"}'
```

Equivalent with explicit channel:

```bash
./bin/ahd-figma command happy-unicorn-42 paint.set_solid -p '{"nodeId":"1:2","color":"#2563EB"}'
```

## Multi-Operation Payload (CLI)
Use `batch` to send multiple operations over one connection.

```bash
./bin/ahd-figma batch -o '[
  {"name":"createCard","command":"shape.create_rectangle","params":{"x":40,"y":40,"width":220,"height":120}},
  {"command":"paint.set_solid","params":{"nodeId":"${{steps.createCard.result.id}}","color":"#2563EB"}},
  {"command":"node.set_opacity","params":{"nodeId":"${{steps.0.result.id}}","opacity":0.9}}
]'
```

Use file input for larger payloads:

```bash
./bin/ahd-figma batch -f operations.json
```

## Multi-Operation Payload Shape
Use CLI `batch` for large multi-step creation. MCP exposes schema-backed command tools for small direct edits.

Example operations payload:

```json
[
  {"name": "createCard", "command": "shape.create_rectangle", "params": {"x": 40, "y": 40, "width": 220, "height": 120}},
  {"command": "paint.set_solid", "params": {"nodeId": "${{steps.createCard.result.id}}", "color": "#2563EB"}}
]
```

## Strategy: When to use what
- `command`: use for direct, isolated changes.
- `batch`: use for related changes that should happen in one logical run.

## High-Fidelity Recreation Payloads
For screenshot or HTML/CSS-like recreation, keep the batch declarative and measured:

1. Add measurement steps first: `text.measure` for real rendered text dimensions and `text.fit_box` for constrained headlines or long tier names.
2. Create the main frame and semantic sections.
3. Use `text.create_rich_block` for copy groups that combine headings, prices, bullets, and notes.
4. Use `layout.pricing_grid` for repeated pricing/package cards instead of manually positioning every text node.
5. Export the result, compare it with the reference, then update the same payload and rerun.

Example:

```bash
./bin/ahd-figma validate docs/examples/html-css-recreation-workflow.json
./bin/ahd-figma batch -f docs/examples/html-css-recreation-workflow.json
```

## Execution Behavior
- Operations run sequentially in order.
- If one fails:
  - CLI/MCP include per-step error details in `results`.
  - Remaining operations continue by default.
- There is no rollback (already-applied operations stay applied).
- Retries are enabled by default (`--retries 1` / `retries: 1`).

## Result Interpolation
A later step can reference data returned by an earlier step in the same batch.

Example idea:
- Step 1 creates a frame and returns `{ \"id\": \"10:20\" }`
- Step 2 uses that ID automatically as `nodeId` without hardcoding it

Why it is useful:
- Eliminates manual ID copy/paste between steps
- Makes LLM-generated multi-step payloads more reliable
- Reduces failures from stale or guessed node IDs
- Enables longer automated chains with less custom glue code

Supported placeholder forms:
- `${{steps.0.result.id}}` (by index)
- `${{steps.createCard.result.id}}` (by named step)
- `${{steps[0].result.id}}` (bracket index form)

If a placeholder cannot be resolved, that step returns an interpolation error.
If `--fail-fast` is not set, subsequent steps still run.

## Resilience Flags (CLI)
- `--retries` retry count per step after first attempt (default `1`)
- `--retry-delay-ms` delay between retries (default `250`)
- `--fail-fast` stop on first failed step (default `false`)
- `--interpolate` enable/disable placeholder interpolation (default `true`)

## Validation Behavior
- `ahd-figma validate` accepts either a plain JSON array or a wrapped `{ "operations": [...] }` payload.
- Interpolation references like `${{steps.createCard.result.id}}` pass pre-execution ID pattern validation.
- Figma RGB objects like `{ "r": 1, "g": 0.5, "b": 0 }` auto-convert to hex for color params.

## Channel Resolution
For `command` and `batch`, channel resolution order is:
1. positional channel arg
2. `--channel`
3. `AHD_CHANNEL` env var
4. relay preferred/active channel

## Discoverability for LLMs
- Use `./bin/ahd-figma tools --json` to enumerate tool/action capabilities.
- Use `./bin/ahd-figma tools --llm --json` for enriched examples and recommended chaining strategy.
- In MCP mode, use `tools/list`, read `ahd://schema` or `ahd://tools`, or call `ahd_describe`.
- Prefer domain actions (e.g., `paint.set_solid`) for consistency.
