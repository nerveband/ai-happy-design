# CLI Batch and Payload Strategy

## Single Operation
Use `command` for one action.

```bash
./bin/ai-happy-design command paint.set_solid -p '{"nodeId":"1:2","color":"#2563EB"}'
```

Equivalent with explicit channel:

```bash
./bin/ai-happy-design command happy-unicorn-42 paint.set_solid -p '{"nodeId":"1:2","color":"#2563EB"}'
```

## Multi-Operation Payload (CLI)
Use `batch` to send multiple operations over one connection.

```bash
./bin/ai-happy-design batch -o '[
  {"name":"createCard","command":"shape.create_rectangle","params":{"x":40,"y":40,"width":220,"height":120}},
  {"command":"paint.set_solid","params":{"nodeId":"${{steps.createCard.result.id}}","color":"#2563EB"}},
  {"command":"node.set_opacity","params":{"nodeId":"${{steps.0.result.id}}","opacity":0.9}}
]'
```

Use file input for larger payloads:

```bash
./bin/ai-happy-design batch -f operations.json
```

## Multi-Operation Payload (MCP)
Use tool `bulk` with action `execute` and pass `operations` as JSON string.

Example operations payload:

```json
[
  {"name": "createCard", "command": "shape.create_rectangle", "params": {"x": 40, "y": 40, "width": 220, "height": 120}},
  {"command": "paint.set_solid", "params": {"nodeId": "${{steps.createCard.result.id}}", "color": "#2563EB"}}
]
```

## Strategy: When to use what
- `command`: use for direct, isolated changes.
- `batch`/`bulk.execute`: use for related changes that should happen in one logical run.

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

## Channel Resolution
For `command` and `batch`, channel resolution order is:
1. positional channel arg
2. `--channel`
3. `AHD_CHANNEL` env var
4. relay preferred/active channel

## Discoverability for LLMs
- Use `./bin/ai-happy-design tools --json` to enumerate tool/action capabilities.
- Use `./bin/ai-happy-design tools --llm --json` for enriched examples and recommended chaining strategy.
- In MCP mode, call `describe(action="catalog")` for the same LLM-oriented catalog.
- Prefer domain actions (e.g., `paint.set_solid`) for consistency.
