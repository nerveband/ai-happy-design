# LLM Discovery Playbook

Use this flow to reduce command failures.

## 1) Discover before acting
CLI:

```bash
ahd-figma tools --llm --json
```

MCP:
- Call `tools/list`, read `ahd://schema` or `ahd://tools`, or call `ahd_describe` with `action: "catalog"`.

This returns:
- tool/action list
- parameter hints
- generated examples (CLI/MCP/batch step)
- a short playbook for chaining/retries

## 2) Probe document context first
Run these before edits:
- `document.get_info`
- `document.get_selection`
- `document.get_focused_node` when working from Dev Mode focus
- `node.get_css` when matching an existing node's generated CSS

This reduces failures from invalid node assumptions.

## 3) For multi-step edits, use one payload
Prefer:
- CLI: `batch`
- MCP: schema-backed tools for small edits; CLI `batch` for multi-step generation

Include:
- `name` on each step
- interpolation placeholders (`${{steps.createCard.result.id}}`)
- resilient settings (`continueOnError: true`, `retries: 1`)

## 4) Image strategy
- Prefer `paint.set_image_fill_from_url` first.
- If URL fails (domain/network/CORS), fallback to `paint.set_image_fill` with base64.
- Verify with `paint.get_fills`.

## 5) Error-handling contract
Always expect wrapped envelopes:
- success: `message.result`
- error: `message.error`

If one step fails in batch/bulk:
- keep processing remaining steps (default)
- return summary with per-step attempts and errors

## LLM Prompt Template
Use this prompt pattern:

```text
1) Discover capabilities first using the catalog.
2) Build the smallest safe command plan.
3) Use CLI batch with named steps and interpolation for multi-step creation.
4) If one step fails, continue and return a per-step summary.
5) For image URL failure, fallback to base64 image fill.
```
