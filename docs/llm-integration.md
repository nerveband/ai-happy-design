# LLM Integration

## Recommended Runtime
Run MCP mode from repo root:

```bash
./bin/ai-happy-design mcp
```

This is the mode LLM tools should target.

## MCP Config Example
Use your local binary path in your MCP client config.

```json
{
  "mcpServers": {
    "ai-happy-design": {
      "command": "/absolute/path/to/ai-happy-design/bin/ai-happy-design",
      "args": ["mcp"]
    }
  }
}
```

## LLM Workflow
1. Start MCP server.
2. Open/run AI Happy Design plugin in Figma.
3. Let plugin auto-connect to relay.
4. Ask the LLM to discover tools first: `describe(action="catalog")`.
5. The catalog includes a `workflow` section, `designPatterns`, and `playbook` — the LLM should read these before generating commands.

## Creating vs Editing (Critical)

**CREATING new designs** (3+ elements): Use `batch` (CLI) or `bulk.execute` (MCP).
- Build a JSON array of operations with named steps and interpolation.
- One connection, all steps in ~6 seconds vs minutes with individual commands.
- CLI: `ai-happy-design batch -f design.json`
- MCP: `bulk.execute` with `operations` JSON string.

**EDITING existing nodes** (1-2 changes): Use single commands.
- `ai-happy-design command paint.set_solid -p '{"nodeId":"1:2","color":"#FF0000"}'`
- Fast, precise, no batch overhead needed.

## Design Quality Guidelines

The catalog's `designPatterns` section teaches LLMs to build professional Figma designs:

- **Frames, not shapes**: Use `node.create_frame` as containers with auto-layout, not rectangles with floating text.
- **Auto-layout everything**: `layout.set_auto_layout` handles spacing, alignment, and centering — no manual x/y math.
- **8px grid**: All dimensions should be multiples of 8 (width: 320, padding: 24, spacing: 16).
- **Cards = frames**: A card is a frame with auto-layout containing text children, not a rectangle + separate text.
- **Typography scale**: 11/12/13/14/16/18/24/32/48/64/72px. Weights: 400/500/600/700/800.
- **Nesting**: Use `parentId` param when creating children, or `layer.move_to_parent` after creation.

## Discovery Commands
- MCP: `describe(action="catalog")` — full catalog with design patterns
- CLI: `ai-happy-design tools --llm --json` — same catalog as JSON
- CLI: `ai-happy-design actions [domain]` — quick domain/action listing

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
