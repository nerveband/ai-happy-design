# Troubleshooting

## Plugin shows disconnected
1. Confirm relay is running:

```bash
PORT=3056 ai-happy-design relay status
```

or

```bash
PORT=3056 ai-happy-design relay start
```

2. Confirm plugin Relay URL matches relay port/path:
- `ws://localhost:3056/ws`
3. Click **Connect** in plugin UI.

If you want to run relay manually in foreground for debugging:

```bash
PORT=3056 ai-happy-design ws
```

## Seeing `handleSocketMessage [object Object]` in console
That message is from the legacy v1 plugin UI (`src/claude_mcp_plugin/ui.html` in the other repo), not v2.

Fix:
1. Remove the old imported plugin entry from Figma development plugins.
2. Re-import v2 manifest:
   - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin/manifest.json`
3. Open **AI Happy Design** (v2) and use the **Relay URL** field.

## Figma internal `vendor-*.min.js.br` / `statsigapi.net` / `font-files` console errors
These are typically Figma app internals or telemetry/network noise and usually not caused by this plugin.

Treat as non-blocking unless you also see plugin-specific errors from:
- `plugin/src/main.ts`
- `plugin/src/ui/app.ts`
- handler files under `plugin/src/handlers/*`

`http://127.0.0.1:<port>/figma/font-files... ERR_CONNECTION_REFUSED` specifically points to Figma's local font helper/runtime path, not AHD relay (`localhost:3055/3056`). If these persist:
1. Update/restart Figma desktop app.
2. Check Figma local permissions/helpers (especially if using browser + font installer).
3. Re-test plugin command path with relay health + channel checks below.

## CLI cannot find channel
Error usually indicates no active plugin channel.

Fix:
1. Open plugin in Figma and connect.
2. Retry command.
3. Or pass explicit channel:

```bash
./bin/ai-happy-design command <channel> document.get_info
```

4. Or set environment variable:

```bash
export AHD_CHANNEL=<channel>
```

## Auto-start disabled error
If you see:
- `relay is not running ... and auto-start is disabled (--no-auto-relay)`

Fix:
1. Remove `--no-auto-relay`
2. Or start relay manually:

```bash
./bin/ai-happy-design relay start
```

## Port occupied by non-relay process
If you see:
- `port <n> is already in use by <process>; refusing to start relay`

Fix:
1. Stop/change the conflicting process.
2. Or run AHD on another port:

```bash
PORT=3060 ./bin/ai-happy-design command document.get_info
```

3. Keep plugin UI port aligned to the same value.

## Multiple active channels error
CLI asks for explicit channel when ambiguity exists.

Fix:
- pass positional channel
- or `--channel`
- or set `AHD_CHANNEL`

## Batch interpolation error
If a step fails with an interpolation error, a placeholder path could not be resolved.

Check:
1. Referenced step ran earlier in the payload.
2. Referenced path exists in that step's result.
3. Placeholder syntax is valid:
- `${{steps.0.result.id}}`
- `${{steps.stepName.result.id}}`

If you need strict stop-on-first-failure behavior, add `--fail-fast`.

## Image URL fill fails
Check:
1. URL is reachable from plugin runtime.
2. Host is listed in `plugin/manifest.json` network access lists.
3. Server returns valid image bytes and allows cross-origin access where needed.

## MCP tool calls timeout
Check response envelope shape from plugin relay path:
- success must include wrapped `message.id` + `message.result`
- error must include wrapped `message.id` + `message.error`

See:
- `docs/image-fill-contract.md`

## Rebuild after code changes
```bash
# Go binary
make build

# Plugin assets
cd plugin
npm run build
cd ..
```
