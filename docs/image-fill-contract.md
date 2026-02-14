# Image Fill Contract (MCP + CLI)

This project supports two image paths that are safe for both MCP and CLI callers:

## URL images (`set_image_fill_from_url`)

1. Caller sends `nodeId`, `url`, and optional `scaleMode`.
2. Plugin tries `figma.createImageAsync(url)` with timeout.
3. If that fails or times out, plugin falls back to `fetch(url)` and then `figma.createImage(bytes)`.
4. Plugin sets exactly one `IMAGE` fill and returns wrapped success/error envelope.

## Base64 images (`set_image_fill`)

1. Caller sends `nodeId`, `imageData` (raw base64 or data URL), optional `scaleMode`.
2. Plugin decodes with `decodeBase64ToBytes`.
3. Plugin calls `figma.createImage(bytes)`.
4. Plugin applies one `IMAGE` fill.

## Message envelope (required)

Successful command responses:

```json
{
  "type": "message",
  "channel": "<channel>",
  "message": { "id": "<requestId>", "result": { "...": "..." } }
}
```

Error responses:

```json
{
  "type": "message",
  "channel": "<channel>",
  "message": { "id": "<requestId>", "error": "..." }
}
```

Relay uses field-presence checks for `result`/`error` so falsey values (for example `result: false`) do not timeout.

## Runtime requirements

- `documentAccess` must be `"dynamic-page"` and handlers should use `await figma.getNodeByIdAsync(...)` for image operations.
- URL hosts must be listed in `manifest.json`:
  - production hosts in `networkAccess.allowedDomains`
  - local/dev hosts in `networkAccess.devAllowedDomains`
- The image server should return valid image bytes and allow browser fetch paths (CORS headers are recommended).

## Recommended caller behavior

1. Prefer `set_image_fill_from_url` for large images.
2. If URL path fails due to network/domain policies, fall back to `set_image_fill` with base64.
3. Verify via `get_fills(nodeId)` that an `IMAGE` paint with `imageHash` exists.
