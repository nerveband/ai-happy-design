# Illustrator CLI Surface

## Core Commands

```text
ahd-illustrator tools [--json] [--llm]
ahd-illustrator schema [domain.action] [--json] [--all --llms-txt]
ahd-illustrator command <domain.action> --json '<payload>' [--dry-run] [--fields '<mask>'] [--output json|ndjson|text]
ahd-illustrator batch --ops <file|json> [--dry-run] [--strict] [--output json|ndjson]
ahd-illustrator host status|open|quit
ahd-illustrator doctor
ahd-illustrator examples [category]
```

`host status` and `doctor` include the resolved Illustrator app path, resolved version, and the current plugin probe result.

## Domain Coverage

| Domain | Commands |
| --- | --- |
| `app.*` | `info`, `version`, `select_tool`, `execute_menu`, `user_interaction_level` |
| `document.*` | `new`, `open`, `save`, `save_as`, `close`, `list`, `info` |
| `artboard.*` | `list`, `create`, `resize`, `set_active`, `fit_to_artwork` |
| `layer.*` | `list`, `create`, `rename`, `visibility`, `lock`, `reorder` |
| `selection.*` | `get`, `clear`, `set_by_ids`, `select_by_name` |
| `path.*` | `create_rect`, `create_ellipse`, `create_path`, `transform`, `duplicate` |
| `text.*` | `create`, `set_contents`, `set_style`, `outline` |
| `appearance.*` | `set_fill`, `set_stroke`, `set_gradient`, `apply_graphic_style` |
| `action.*` | `load`, `run`, `unload` |
| `export.*` | `png`, `jpg`, `svg`, `pdf`, `ai` |
| `inspect.*` | `tree`, `styles`, `bounds`, `fonts`, `summary` |

## Output Contract

Single-command success:

```json
{
  "ok": true,
  "requestId": "uuid",
  "command": "text.create",
  "result": {},
  "warnings": [],
  "timingMs": 0
}
```

Single-command error:

```json
{
  "ok": false,
  "requestId": "uuid",
  "command": "export.png",
  "timingMs": 0,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "human-readable",
    "details": {}
  },
  "retryable": false
}
```

Batch success/error summary:

```json
{
  "ok": true,
  "requestId": "uuid",
  "summary": {
    "total": 2,
    "succeeded": 2,
    "failed": 0
  },
  "steps": [],
  "timingMs": 0,
  "retryable": false
}
```

## Discovery and Safety

- `tools --json` is the fastest way for an agent to discover the command surface.
- `schema <command> --json` is canonical for payload generation.
- Use `--dry-run` on all mutating commands first.
- `--fields` trims the top-level response envelope to reduce context bloat.
- `batch --strict` stops on the first invalid step.
- Live validation should use a scratch document and disposable export paths.

## Batch Interpolation

Step references use the form `${{steps.step_name.result.path}}`.

Example:

```json
[
  {
    "name": "new_doc",
    "command": "document.new",
    "params": { "width": 1440, "height": 1024 }
  },
  {
    "name": "cover_board",
    "command": "artboard.create",
    "params": {
      "name": "Cover",
      "left": 0,
      "top": 1024,
      "right": "${{steps.new_doc.result.params.width}}",
      "bottom": 0
    }
  }
]
```

Interpolation is intentionally strict:

- missing step names fail validation
- missing result paths fail validation
- non-scalar values cannot be embedded into larger strings

## Plugin-Required Commands

These currently route through the `sendScriptMessage` bridge:

- `inspect.tree`
- `inspect.styles`
- `inspect.bounds`
- `inspect.fonts`
- `inspect.summary`
- `appearance.set_gradient`
- `appearance.apply_graphic_style`

If the plugin bridge is unavailable, they fail with `PLUGIN_REQUIRED`.

The plugin probe currently uses `ahd.version`. A healthy installed bridge should return a JSON payload for:

```text
app.sendScriptMessage("AHDIllustrator", "ahd.version", "{}")
```
