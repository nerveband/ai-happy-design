# AI Happy Design v2

A single Go binary that serves as both an MCP (Model Context Protocol) server for AI code editors and a CLI for direct Figma manipulation. Communicates with a Figma plugin through an embedded WebSocket relay.

## Features

- **MCP Server** (stdio) for Claude Code, Cursor, Windsurf, and other AI editors
- **WebSocket relay** for real-time communication with the Figma plugin
- **CLI mode** for direct interaction with connected Figma sessions
- **17 domain tools** covering paint, shape, text, layout, node, layer, component, style, variable, effect, boolean, page, document, export, bulk operations, self-documentation, and connection management

## Quick Start

```bash
# Build
make build

# Start MCP server (stdio + WebSocket relay)
./bin/ai-happy-design mcp

# Start WebSocket relay only
./bin/ai-happy-design ws

# Connect to a plugin channel
./bin/ai-happy-design connect happy-unicorn-42
```

## MCP Configuration

Add to your editor's MCP config:

```json
{
  "mcpServers": {
    "ai-happy-design": {
      "command": "/path/to/ai-happy-design",
      "args": ["mcp"]
    }
  }
}
```

## License

MIT License - see LICENSE file.
