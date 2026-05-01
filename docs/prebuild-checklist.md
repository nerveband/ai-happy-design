# Prebuild Checklist

Use this before building, testing, or running MCP/CLI.

## 1. Environment
- Go `1.22+`
- Node.js `18+`
- npm
- Figma Desktop (for plugin runtime testing)

Quick checks:

```bash
go version
node -v
npm -v
```

## 2. Dependency Install
From repo root:

```bash
go mod download
```

From plugin folder:

```bash
cd plugin
npm install
cd ..
```

## 3. Build Artifacts
Build binary and plugin bundle:

```bash
make build
```

Expected outputs:
- `bin/ai-happy-design`
- `bin/ahd-figma`
- `plugin/dist/code.js`
- `plugin/dist/ui.html`

## 4. Fast Validation

```bash
go test ./...
go build ./...
cd plugin && npm run typecheck && npm run verify:syntax && cd ..
```

## 5. Runtime Bring-up
Start MCP mode:

```bash
./bin/ahd-figma mcp
```

Then in Figma run plugin **AI Happy Design**.

## 6. Smoke Commands
In another terminal:

```bash
./bin/ahd-figma tools --json
./bin/ahd-figma command document.get_info
```

If multiple channels are active, pass explicit channel or set `AHD_CHANNEL`.
