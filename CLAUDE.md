@AGENTS.md

# AI Happy Design - Agent Instructions

## Build & Deploy (REQUIRED)

Every time you build the Go binary, you MUST sign it for macOS and restart the relay. Unsigned binaries will be blocked by Gatekeeper. The relay is a long-running daemon — replacing the binary on disk does NOT update the running process.

**Preferred: One-command deploy** (builds plugin + Go, signs, installs, restarts relay):
```bash
make deploy
```

**Manual steps** (if `make deploy` isn't suitable):
```bash
# 1. Build plugin (from plugin/ dir)
cd plugin && npm run build && cd ..

# 2. Sync plugin files into Go embed
make sync-plugin

# 3. Build Go binary
go build -o ai-happy-design ./cmd/ai-happy-design/

# 4. Sign (REQUIRED on macOS) — use /tmp to avoid Dropbox path issues
cp ai-happy-design /tmp/ai-happy-design
codesign -f -s - /tmp/ai-happy-design
cp /tmp/ai-happy-design ~/bin/ai-happy-design

# 5. Restart relay (REQUIRED — old process uses stale binary)
pkill -f "ai-happy-design ws" 2>/dev/null; sleep 1
nohup ~/bin/ai-happy-design ws > /tmp/ahd-relay.log 2>&1 &
```

After deploying, **reopen the Figma plugin** to load the new `code.js`.

The Dropbox path with spaces breaks `codesign -v` verification. Always sign from a simple path like `/tmp`.

## esbuild HTML Inlining

The plugin build (`plugin/esbuild.config.mjs`) inlines CSS and JS into `index.html` using `.replace()`. The replacement MUST use function callbacks:

```js
.replace('/* STYLES */', () => css)
.replace('/* SCRIPT */', () => js)
```

**Never** use plain string replacement (`.replace('/* SCRIPT */', js)`) — JavaScript's `$&` in minified output gets interpreted as a back-reference, injecting the matched pattern into the JS and breaking the plugin.

## Plugin Icon

The plugin icon is an inline SVG in `plugin/src/ui/index.html`. Figma's CSP blocks:
- Dynamic `data:` URIs via JS (`img.src`, `backgroundImage`)
- CSS `background-image` with base64

Only inline SVG works reliably in Figma plugin sandbox.
