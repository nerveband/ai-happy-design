# LLM Integration

## Recommended Runtime
Run MCP mode from repo root:

```bash
./bin/ahd-figma mcp
```

This is the mode LLM tools should target.

## MCP Config Example
Use your local binary path in your MCP client config.

```json
{
  "mcpServers": {
    "ahd-figma": {
      "command": "/absolute/path/to/ai-happy-design/bin/ahd-figma",
      "args": ["mcp"]
    }
  }
}
```

## LLM Workflow
1. Start MCP server.
2. Open/run AI Happy Design plugin in Figma.
3. Let plugin auto-connect to relay.
4. Ask the LLM to discover tools first with MCP `tools/list` or `ahd_describe {"action":"catalog"}`.
5. The catalog includes a `workflow` section, `designPatterns`, and `playbook` — the LLM should read these before generating commands.

## Creating vs Editing (Critical)

**CREATING new designs** (3+ elements): Use `batch` through the CLI.
- Build a JSON array of operations with named steps and interpolation.
- One connection, all steps in ~6 seconds vs minutes with individual commands.
- CLI: `ahd-figma batch -f design.json`
- MCP: call the schema-backed command tools for small edits; use the CLI batch path for large multi-step creation.

**EDITING existing nodes** (1-2 changes): Use single commands.
- `ahd-figma command paint.set_solid -p '{"nodeId":"1:2","color":"#FF0000"}'`
- Fast, precise, no batch overhead needed.

## Design Quality Guidelines

The catalog's `designPatterns` section teaches LLMs to build professional Figma designs:

- **Frames, not shapes**: Use `node.create_frame` as containers with auto-layout, not rectangles with floating text.
- **Auto-layout everything**: `layout.set_auto_layout` handles spacing, alignment, and centering — no manual x/y math.
- **8px grid**: All dimensions should be multiples of 8 (width: 320, padding: 24, spacing: 16).
- **Cards = frames**: A card is a frame with auto-layout containing text children, not a rectangle + separate text.
- **Typography scale**: 11/12/13/14/16/18/24/32/48/64/72px. Weights: 400/500/600/700/800.
- **Nesting**: Use `parentId` param when creating children, or `layer.move_to_parent` after creation.

## HTML/CSS-Like Recreation Workflow

For high-fidelity recreation from a screenshot, web mockup, or existing Figma node, use a measured loop instead of guessing coordinates:

1. Read the reference: use `node.get_css` for existing Figma context, or convert the visual into sections, grids, cards, and text roles.
2. Measure first: call `text.measure` for known type sizes and `text.fit_box` for labels that must fit fixed card/header boxes.
3. Create semantic text: use `text.create_rich_block` for a heading, price/kicker, bullets, and notes that should move as one content unit.
4. Use CSS-like layout primitives: `layout.pricing_grid` handles repeated pricing cards with shared columns, gaps, and measured text blocks.
5. Iterate visually: run the batch, export the frame, compare to the reference, then tune spacing/type/color in the payload and rerun.

Use `docs/examples/html-css-recreation-workflow.json` as the compact starting point.

## Discovery Commands
- MCP: `tools/list`, `resources/list`, and `ahd_describe {"action":"catalog"}`
- CLI: `ahd-figma tools --llm --json` — same catalog as JSON
- CLI: `ahd-figma actions [domain]` — quick domain/action listing

## Current Figma Context Tools
- `document.get_focused_node` reads the current Dev Mode focused node when Figma exposes it.
- `node.get_css` calls Figma `getCSSAsync()` for a selected node and returns generated CSS alongside node summary data.

## Verification Prompt Ideas
- "List available tools and summarize key actions."
- "Get current document info and selection."
- "Create a card with title and description using auto-layout."

## Command/Result Contract
Responses are correlated by request ID and wrapped using the message envelope.

Success envelope:

```json
{
  "type": "message",
  "channel": "<channel>",
  "message": { "id": "<requestId>", "result": {"...": "..."} }
}
```

Error envelope:

```json
{
  "type": "message",
  "channel": "<channel>",
  "message": { "id": "<requestId>", "error": "..." }
}
```

This prevents MCP/CLI timeouts caused by unwrapped responses.
