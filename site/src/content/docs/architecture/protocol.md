---
title: WebSocket Protocol
description: Message format, response envelopes, channel management, error handling, and image fill flow for the AI Happy Design relay protocol.
---

AI Happy Design uses a WebSocket-based protocol to communicate between the CLI/MCP server and the Figma plugin. The relay server on port 3055 routes messages between clients and plugins via named channels.

## Connection Flow

```
1. Relay starts:       ai-happy-design ws  (listens on ws://127.0.0.1:3055)
2. Plugin connects:    WebSocket open -> sends join message
3. CLI connects:       WebSocket open -> sends command
4. Relay routes:       CLI command -> plugin (same channel)
5. Plugin responds:    Result -> relay -> CLI
6. CLI disconnects:    WebSocket close (per-command lifecycle)
```

The plugin maintains a persistent connection. CLI clients connect per command (or per batch session) and disconnect when done.

## WebSocket Endpoint

The relay listens on `ws://127.0.0.1:3055` by default. All communication is local -- nothing leaves your machine.

### Connection Limits

| Property | Value |
|----------|-------|
| Max message size | 64 MB |
| Read timeout | 300 seconds |
| Write deadline | 120 seconds |
| Default port | 3055 |

The 64 MB read limit was increased from 1 MB to support large image exports (e.g., 2160x3840 frames at 2x scale produce multi-megabyte PNG payloads).

## Message Format

### Join Message

When the Figma plugin connects, it sends a join message to register on a channel:

```json
{
  "id": "uuid-1",
  "type": "join",
  "channel": "channel-key"
}
```

The relay confirms with a `joined` response:

```json
{
  "type": "joined",
  "channel": "channel-key"
}
```

The plugin UI shows "Connected" only after receiving the `joined` message, not on WebSocket open. If you see "Connecting..." persistently, the relay may not be running.

### Command Message (CLI/MCP to Plugin)

```json
{
  "id": "cmd_1710000000_001",
  "type": "message",
  "channel": "channel-key",
  "message": {
    "id": "cmd_1710000000_001",
    "command": "text.create",
    "params": {
      "parentId": "42:248",
      "text": "Hello World",
      "fontSize": 28,
      "fontFamily": "Inter",
      "fontStyle": "Bold",
      "color": "#FFFFFF",
      "x": 32,
      "y": 32,
      "commandId": "cmd_1710000000_001"
    }
  }
}
```

Fields:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique message ID for correlation |
| `type` | string | Always `"message"` for commands |
| `channel` | string | Channel name for routing |
| `message.id` | string | Same as outer `id` (for plugin-side correlation) |
| `message.command` | string | The `domain.action` command name |
| `message.params` | object | Command parameters |
| `message.params.commandId` | string | Injected by CLI for response matching |

### Response Message -- Success

```json
{
  "type": "message",
  "channel": "channel-key",
  "message": {
    "id": "cmd_1710000000_001",
    "result": {
      "id": "42:249",
      "name": "Hello World",
      "type": "TEXT",
      "width": 200,
      "height": 28
    }
  }
}
```

### Response Message -- Error

```json
{
  "type": "message",
  "channel": "channel-key",
  "message": {
    "id": "cmd_1710000000_001",
    "error": "Node not found: 99:999"
  }
}
```

Responses are always wrapped in the full envelope `{type:"message", channel, message:{id, result/error}}`. The relay never sends bare `{id, error}` objects. This is a protocol invariant.

A response has either `result` or `error`, never both.

## Channel System

Channels route messages to the correct Figma plugin instance. This matters when you have multiple Figma files open, each running the plugin.

### How Channels Work

1. The Figma plugin generates a channel key on startup
2. The plugin joins that channel on the relay
3. CLI commands target a channel to reach the right plugin instance
4. The relay forwards messages only to clients on the matching channel

### Channel Resolution Order

When a CLI command does not specify a channel explicitly:

1. **Positional argument** -- `ai-happy-design command text.create '{}' my-channel`
2. **`--channel` flag** -- `ai-happy-design command text.create '{}' --channel my-channel`
3. **`AHD_CHANNEL` env var** -- `export AHD_CHANNEL=my-channel`
4. **Relay preferred/active channel** -- the relay tracks which channel was most recently joined and uses it as the default

For most single-file workflows, the channel is resolved automatically (option 4). Explicit channel selection matters when multiple Figma files are open simultaneously.

### Multiple Channels

```
Relay (port 3055)
    +-- Channel "project-a" <-> Figma File A (plugin)
    +-- Channel "project-b" <-> Figma File B (plugin)
    +-- Channel "project-c" <-> Figma File C (plugin)
```

### Checking Connection Status

```bash
ai-happy-design command connect.status
```

Returns the current connection state and active channel.

## Plugin Connection Flow

1. Plugin starts in Figma
2. Plugin opens WebSocket to `ws://127.0.0.1:3055`
3. Plugin sends a `join` message with its channel key
4. Relay acknowledges with a `joined` response
5. Plugin UI changes from "Connecting..." to "Connected"
6. The relay is now ready to route commands

## Dynamic Page Access

The Figma Plugin API requires using asynchronous node access:

```typescript
// Correct -- async access works across pages
const node = await figma.getNodeByIdAsync(nodeId);

// Deprecated -- only works on current page
const node = figma.getNodeById(nodeId);  // DON'T USE
```

All AI Happy Design handlers use `getNodeByIdAsync` via the shared `getNode()` utility, which also calls `loadAsync()` before lookup to handle recently-created nodes.

## Stable Node IDs

Figma's `create*` APIs (createFrame, createText, etc.) return transient session IDs that may expire after the handler completes. The plugin resolves stable IDs by reading the committed ID from the parent's `.children` list.

This is handled automatically by `resolveStableId()` in all creation handlers (node, text, shape, component).

## Image Fill Flow

Setting image fills involves a multi-step process because Figma's plugin sandbox has restrictions on image data handling.

### URL-Based (`paint.set_image_url`)

```
1. Plugin receives: {nodeId, url, scaleMode}
2. Try: hash = await figma.createImageAsync(url)
3. If that fails (CORS, network):
   a. Fallback: response = await fetch(url)
   b. bytes = new Uint8Array(await response.arrayBuffer())
   c. hash = figma.createImage(bytes)
4. Apply: node.fills = [{type:"IMAGE", imageHash:hash, scaleMode}]
```

The `figma.createImageAsync(url)` path is preferred because it handles Figma's internal caching. The fetch fallback handles URLs that Figma cannot reach directly.

URL-based image fills require the domain to be listed in the plugin manifest's `allowedDomains` (production) or `devAllowedDomains` (development).

### Base64/Data URL (`paint.set_image`)

```
1. Plugin receives: {nodeId, imageData, scaleMode}
2. Decode base64 or data URL to Uint8Array
3. hash = figma.createImage(bytes)
4. Apply: node.fills = [{type:"IMAGE", imageHash:hash, scaleMode}]
```

### One-Step Image (`shape.create_image`)

```
1. Plugin receives: {imageData, x, y, width, height, parentId, scaleMode}
2. Create rectangle: figma.createRectangle()
3. Set position and size
4. Decode image data -> figma.createImage(bytes)
5. Apply image fill to rectangle
6. Return node info
```

### Export Flow

```
CLI -> export command -> relay -> plugin -> node.exportAsync()
    -> base64 PNG/JPG/SVG -> relay -> CLI -> save to OS temp dir
```

The export response can be up to 64MB (the relay read limit), supporting high-resolution exports. Exports auto-save to `os.TempDir()` and the response includes the file path.

## Error Handling

### Transport Errors

If the WebSocket connection drops, the CLI retries with exponential backoff (up to 3 attempts by default in batch mode).

### Command Errors

Command errors are returned in the response envelope:

```json
{
  "type": "message",
  "channel": "channel-key",
  "message": {
    "id": "cmd_001",
    "error": "Cannot set layoutPositioning:ABSOLUTE on child of non-auto-layout parent"
  }
}
```

Error messages are passed through from the plugin without modification. They are designed to be useful to both humans and LLM agents.

### Common Error Types

| Error | Cause |
|-------|-------|
| `Node not found: 99:999` | Invalid node ID or node was deleted |
| `Font not available: Nonexistent Font Bold` | Font family or style not installed |
| `Cannot set layoutPositioning:ABSOLUTE...` | Trying to absolute-position a child of a non-auto-layout parent |
| `Fill index 0 out of range` | Frame has no fills (safe to ignore) |

### Timeout

If the plugin does not respond within 300 seconds, the CLI treats the command as failed and returns a timeout error.

### Read-Modify-Write Safety

When modifying effects or fills, the plugin:

1. Reads the current effects/fills with `node.effects.slice()` or `node.fills.slice()`
2. **Sanitizes** the copy (removes internal Figma properties like `boundVariables`)
3. Modifies the sanitized copy
4. Writes back the clean array

This prevents "Cannot set properties of read-only" errors that occur when passing internal Figma objects back. The sanitization utilities (`sanitizeEffects.ts`, `sanitizeFills.ts`) allowlist only the properties that Figma accepts on write-back.

## Batch Protocol

In batch mode, the CLI sends operations sequentially over a single WebSocket connection:

```
CLI                    Relay                  Plugin
 |                       |                       |
 |-- command (step 0) -->|-- command (step 0) -->|
 |                       |                       |-- execute
 |<-- response (step 0) -|<-- response (step 0) -|
 |                       |                       |
 |-- command (step 1) -->|-- command (step 1) -->|
 |                       |                       |-- execute
 |<-- response (step 1) -|<-- response (step 1) -|
```

The CLI waits for each response before sending the next command. This ensures step interpolation (`${{steps.NAME.result.id}}`) has the data it needs.

With `--parallel`, the CLI identifies independent operations (no interpolation dependencies) and sends them concurrently.

### Batch Error Modes

| Flag | Behavior |
|------|----------|
| (default) | Continue on error, report failures in summary |
| `--fail-fast` | Stop batch on first error |

## Protocol Invariants

These rules must never be violated:

1. **All responses use the full envelope** -- `{type:"message", channel, message:{id, result/error}}`. Never bare `{id, result}`.
2. **Responses have either `result` or `error`** -- never both, never neither.
3. **The `id` field correlates request to response** -- the plugin must echo back the same `id` it received.
4. **Join must precede commands** -- a plugin must send a `join` message before it can receive commands on a channel.
5. **Node access uses async API** -- all node lookups use `figma.getNodeByIdAsync()`, never the deprecated sync getter.
6. **Stable IDs** -- creation handlers resolve transient session IDs to stable committed IDs before returning results.

---

Made by [Ashraf Ali](https://ashrafali.net) | [GitHub](https://github.com/nerveband/ai-happy-design) | License: GPL-3.0
