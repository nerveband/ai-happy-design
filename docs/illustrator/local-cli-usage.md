# AHD Illustrator Local CLI Usage

This machine now has these binaries installed in `~/bin`:

- `ahd-illustrator`
- `ahd-figma`

Current live environment:

- Illustrator app: `/Applications/Adobe Illustrator 2026/Adobe Illustrator.app`
- Illustrator version: `30.2.1`
- Plugin bridge: not installed
- Current mode: script-first CLI, which is enough for the current Phase 1 Illustrator surface

## Quick Checks

Check that the binaries are on your path:

```bash
which ahd-illustrator
which ahd-figma
```

Check Illustrator host health:

```bash
ahd-illustrator doctor
ahd-illustrator host status
```

## Agent-First Workflow

The CLI is designed to be discoverable first, then executable.

1. List available tools:

```bash
ahd-illustrator tools --json
```

2. Inspect a command schema:

```bash
ahd-illustrator schema text.create --json
ahd-illustrator schema app.info --json
```

3. Dry-run a command before mutation:

```bash
ahd-illustrator command document.new --json '{"preset":"Web","width":1440,"height":1024}' --dry-run
```

4. Execute the real command:

```bash
ahd-illustrator command app.info --json '{}'
```

5. Use batch for multi-step work:

```bash
ahd-illustrator batch --ops ops.json --strict
```

## Commands That Work Without Any Plugin

These are script-backed in the current implementation:

- `app.*`
- `document.*`
- `artboard.*`
- `layer.*`
- `selection.*`
- `path.*`
- `text.*`
- `appearance.*`
- `action.*`
- `export.*`
- `inspect.*`
- `capture.*`
- `dataset.*`
- `matrix.*`
- `page_item.*`
- `placed.*`
- `preference.*`
- `print.*`
- `raster.*`
- `repeat.grid.*`
- `repeat.radial.*`
- `repeat.symmetry.*`
- `spot.*`
- `style.character.*`
- `style.graphic.*`
- `style.paragraph.*`
- `swatch.*`
- `symbol.*`
- `trace.preset.*`
- `variable.*`
- `view.*`
- `workspace.*`

The plugin is still optional and currently only shows up as diagnostic state in `doctor` and `host status`.

## Practical Examples

Get runtime info:

```bash
ahd-illustrator command app.info --json '{}'
```

Create a scratch document:

```bash
ahd-illustrator command document.new --json '{"preset":"Web","width":1440,"height":1024,"units":"Pixels"}'
```

Create a layer:

```bash
ahd-illustrator command layer.create --json '{"name":"Hero"}'
```

Create a rectangle:

```bash
ahd-illustrator command path.create_rect --json '{"left":80,"top":920,"width":320,"height":180,"name":"Card"}'
```

Create text:

```bash
ahd-illustrator command text.create --json '{"contents":"Hello from AHD","left":100,"top":860,"name":"Headline"}'
```

Export PNG:

```bash
ahd-illustrator command export.png --json '{"filePath":"./out/example.png","scale":2}'
```

Close the scratch document without saving:

```bash
ahd-illustrator command document.close --json '{"save":false}'
```

## Batch Example

Save this as `ops.json`:

```json
[
  {
    "name": "new_doc",
    "command": "document.new",
    "params": {
      "preset": "Web",
      "width": 1440,
      "height": 1024,
      "units": "Pixels"
    }
  },
  {
    "name": "hero_layer",
    "command": "layer.create",
    "params": {
      "name": "Hero"
    }
  },
  {
    "name": "hero_rect",
    "command": "path.create_rect",
    "params": {
      "left": 80,
      "top": 920,
      "width": 420,
      "height": 220,
      "name": "Hero Card"
    }
  },
  {
    "name": "hero_text",
    "command": "text.create",
    "params": {
      "contents": "Agent-first Illustrator",
      "left": 110,
      "top": 860,
      "name": "Hero Title"
    }
  }
]
```

Then run:

```bash
ahd-illustrator batch --ops ops.json --strict
```

## Figma CLI Notes

The Figma binary is also installed:

```bash
ahd-figma tools --llm --json
ahd-figma schema shape.create --json
ahd-figma schema effect.add --json
```

The binary is installed, but the Figma plugin itself still has to be loaded in Figma if you want live canvas operations there.

## Known Current Boundary

- `ahd-illustrator` is installed and working in script-first mode.
- The optional native Illustrator plugin is not installed yet.
- `doctor` will report `PLUGIN_REQUIRED` until the plugin bridge is built and installed.

## Files to Know

- Main Illustrator docs: `docs/illustrator/README.md`
- Command reference: `docs/illustrator/commands.md`
- Runtime learnings: `docs/illustrator/live-script-runtime-learnings.md`
- Phase 2 native plugin plan: `docs/plans/2026-03-05-ahd-illustrator-phase-2-native-plugin-bridge.md`
