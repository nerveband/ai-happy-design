# Turnkey Pairing Design (v2)

## Scope
Repository: `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design`

Problem to solve: plugin reconnect friction (new/random channel IDs, manual reconnect steps) and CLI friction (channel arg required repeatedly).

## Goals
1. Plugin should keep the same identity/channel between opens.
2. Plugin should auto-connect when possible.
3. MCP/CLI should route commands without repeated manual channel management.
4. Error/result envelope behavior must remain stable for CLI/MCP.

## Implemented in this session

### 1) Persistent plugin pairing state
- Implemented `figma.clientStorage` settings in plugin runtime:
  - `channelKey`
  - `port`
  - `autoConnect`
- Main/runtime now supports UI messages:
  - `request-connection-settings`
  - `save-connection-settings`
- Files:
  - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin/src/main.ts`
  - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin/src/ui/app.ts`

### 2) Auto-connect startup flow
- UI requests settings at startup and applies them.
- If `autoConnect` is enabled (default), UI initiates connection automatically.
- Reconnect policy changed to retry indefinitely (with backoff), so the plugin eventually recovers when relay starts later.
- Files:
  - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin/src/ws/client.ts`

### 3) Channel lifecycle UX improvements
- Added channel regeneration button in plugin UI (`New`) for explicit re-pairing.
- Keeps copy flow and persisted update path.
- Files:
  - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin/src/ui/index.html`
  - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/plugin/src/ui/styles.css`

### 4) Relay preferred-channel selection
- Relay now tracks a preferred active channel and returns deterministic fallback instead of map-order randomness.
- `/status` endpoint now includes `preferredChannel`.
- File:
  - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/internal/ws/server.go`

### 5) CLI auto channel resolution
- `command` and `batch` now support running without explicit positional channel.
- Resolution order:
  1. positional channel argument
  2. `--channel`
  3. `AHD_CHANNEL` env
  4. relay `/status` preferred/active channel
- File:
  - `/Users/nerveband/wavedepth Dropbox/Ashraf Ali/Mac (2)/Documents/GitHub/ai-happy-design/cmd/ai-happy-design/main.go`

## Resulting user flow
1. Open plugin first time -> generated channel persists.
2. Open plugin later -> same channel appears and auto-connect attempts begin.
3. Run CLI command without channel (if one preferred/active channel exists) -> routes automatically.
4. MCP uses relay active/preferred channel selection with deterministic behavior.

## Recommended next phase
1. Add explicit `pair` command group:
- `pair status`
- `pair set <channel>`
- `pair clear`
- `pair list`

2. Add plugin identity handshake (optional, stronger pairing):
- persist `clientId` in plugin storage
- include `clientId` in `join` message
- relay tracks `clientId -> preferred channel`

3. Add authentication hardening for non-local relays:
- pre-shared token in join/command envelopes
- allowlist host/port settings

4. Add focused tests:
- CLI channel auto-resolution behavior
- relay preferred channel fallback
- plugin settings request/save protocol

## Verification commands used
- `go test ./...`
- `go build ./...`
- `cd plugin && npm run build`

All passed in this session.
