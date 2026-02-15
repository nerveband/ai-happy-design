# Figma MCP Image Fill - Debug Investigation

## Problem
Image fill operations (`set_image_fill` and `set_image_fill_from_url`) time out when called via the AI Happy Design MCP plugin. All other commands (create_frame, create_text, set_font_name, etc.) work fine.

## Root Cause Found
**Bug in WebSocket relay** (`src/talk_to_figma_mcp/utils/websocket.ts` line 107-110):
```typescript
// OLD - drops error responses because they have no `result` field
if (myResponse.id && pendingRequests.has(myResponse.id) && myResponse.result)

// FIXED - also handles error-only responses
if (myResponse.id && pendingRequests.has(myResponse.id) && (myResponse.result || myResponse.error))
```
Error responses from the plugin have `{ id, error }` but NO `result` field. The old code required `myResponse.result` to be truthy, so errors were silently dropped and requests hung until timeout.

**This fix has been applied and rebuilt** but needs testing.

## Additional Issues Found & Fixed

### 1. `figma.base64Decode` doesn't exist
- **File**: `src/claude_mcp_plugin/code.js` line 4132
- **Old**: `figma.base64Decode(imageData)` - not a real Figma Plugin API method
- **Fixed**: Manual `atob()` + Uint8Array conversion
- **Note**: `atob` may also not exist in the Figma plugin main thread sandbox

### 2. `createImageAsync` can hang forever
- **File**: `src/claude_mcp_plugin/code.js` line 4163
- **Fixed**: Wrapped in `Promise.race` with 30s timeout
- Now returns actual error instead of hanging

### 3. Plugin network access too restrictive
- **File**: `src/claude_mcp_plugin/manifest.json`
- `devAllowedDomains` only had `localhost:3055`
- Added `localhost:8765` for local image server
- Added `picsum.photos`, `fastly.picsum.photos`, `i.imgur.com` to `allowedDomains`

### 4. Local HTTP server needs CORS headers
- Figma plugin iframes have `null` origin
- Python HTTP server at `/tmp/cors_server.py` serves images with `Access-Control-Allow-Origin: *`

## Key Resources
- **Figma Plugin Network Requests**: https://developers.figma.com/docs/plugins/making-network-requests/
- **Figma Working with Images**: https://developers.figma.com/docs/plugins/working-with-images/
- **Figma createImageAsync API**: https://developers.figma.com/docs/plugins/api/properties/figma-createimageasync
- **Figma createImage API**: https://developers.figma.com/docs/plugins/api/properties/figma-createimage

## Key Figma API Facts
- `figma.createImageAsync(url)` - fetches image from URL, needs domain in `networkAccess`
- `figma.createImage(bytes: Uint8Array)` - creates image from raw bytes (no network needed)
- Plugin main thread is NOT a browser - may not have `atob`, `fetch`, etc.
- Plugin UI iframe IS a browser - has `fetch`, `canvas`, `atob`, etc.
- Images must be PNG, JPG, or GIF, max 4096x4096
- CORS applies: server needs `Access-Control-Allow-Origin: *`

## Architecture: Message Flow
```
MCP Server (TypeScript/Bun)
  → WebSocket relay (localhost:3055)
    → ui.html (iframe, receives WS messages)
      → code.js (main plugin thread via figma.ui.postMessage)
        → executes command (e.g., createImageAsync)
      ← code.js sends result back via figma.ui.postMessage
    ← ui.html forwards result via WebSocket
  ← WebSocket relay matches response to pending request
← MCP Server returns result
```

## Setup Instructions for New AI Instance

### Prerequisites
- Bun runtime installed (`~/.bun/bin/bun`)
- Figma Desktop app
- AI Happy Design plugin loaded from manifest

### 1. Start WebSocket Relay
```bash
cd /Users/nerveband/Documents/GitHub/ai-happy-design
bun run src/socket.ts
# Runs on port 3055
```

### 2. Start Local Image Server (for testing)
```bash
python3 /tmp/cors_server.py
# Serves images from /tmp/carousel-images/ on port 8765 with CORS headers
```

### 3. Configure MCP in codex
```AIHappyDesign -- /Users/nerveband/.bun/bin/bun run /Users/nerveband/Documents/GitHub/ai-happy-design/dist/talk_to_figma_mcp/server.js
```

### 4. Load Plugin in Figma
- Plugins > Development > Import plugin from manifest
- Select: `~/Documents/GitHub/ai-happy-design/src/claude_mcp_plugin/manifest.json`
- Click "Run" on the plugin
- Note the channel ID shown in the plugin UI

### 5. Connect in codex
```
join_channel with channel: "<channel_id>"
```

### 6. Test Image Fill
```
set_image_fill_from_url with nodeId: "<rect_id>", url: "http://localhost:8765/style2-03-hospital-waiting.png"
```

## Key Files
- **MCP Server**: `src/talk_to_figma_mcp/` (TypeScript)
- **WebSocket relay**: `src/socket.ts`
- **Plugin code**: `src/claude_mcp_plugin/code.js` (plain JS, loaded directly)
- **Plugin UI**: `src/claude_mcp_plugin/ui.html`
- **Plugin manifest**: `src/claude_mcp_plugin/manifest.json`
- **WebSocket utils (where bug was)**: `src/talk_to_figma_mcp/utils/websocket.ts`
- **Build**: `bun run build` (uses tsup, outputs to `dist/`)

## Next Steps
1. **Test with rebuilt MCP server** - the WebSocket error-handling fix should now surface actual error messages instead of generic timeouts
2. **If `atob` doesn't work** in plugin main thread, alternatives:
   - Route image fetch through UI iframe (has full browser APIs)
   - Use `figma.ui.postMessage` to send base64 to UI, decode there, send bytes back
3. **If `createImageAsync` can't reach localhost**, alternatives:
   - Upload images to public host (imgur, imgbb)
   - Proxy through UI iframe's `fetch()`
4. **Still uses `figma.base64Decode`** in `addFill` function (line 4200) - needs same fix

## Test Images on Desktop
```
/Users/nerveband/Desktop/style1-*.png  (10 images, sketch/line-art style)
/Users/nerveband/Desktop/style2-*.png  (10 images, full-color flat illustration)
/Users/nerveband/Desktop/chaplaincy-*.png  (6 images)
```

## Current Carousel Status
7 slide frames created with all text content and styling (Poppins + DM Sans) on page "AMC Carousel - The Care Behind the Care". Images not yet applied. Design uses green (#0F4D2E), teal (#2E8FA6), cream (#F5F2EB) palette.
