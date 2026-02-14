# AI Happy Design - Agent Instructions

## Build & Sign (REQUIRED)

Every time you build the Go binary, you MUST sign it for macOS. Unsigned binaries will be blocked by Gatekeeper.

```bash
# 1. Build plugin (from plugin/ dir)
cd plugin && npm run build && cd ..

# 2. Sync plugin files into Go embed
make sync-plugin

# 3. Build Go binary
go build -o ai-happy-design ./cmd/ai-happy-design/

# 4. Sign (REQUIRED on macOS) — use /tmp to avoid Dropbox path issues
cp ai-happy-design /tmp/ai-happy-design
codesign -s - /tmp/ai-happy-design
cp /tmp/ai-happy-design ~/bin/ai-happy-design
```

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
