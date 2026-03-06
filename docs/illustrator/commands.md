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
| `app.*` | `info`, `version`, `select_tool`, `execute_menu`, `user_interaction_level`, `beep`, `copy`, `cut`, `paste`, `undo`, `redo`, `redraw`, `convert_sample_color`, `translate_placeholder_text`, `preset_lists`, `get_preset_file`, `get_preset_settings`, `load_color_settings`, `show_presets` |
| `document.*` | `new`, `open`, `save`, `save_as`, `close`, `list`, `info`, `activate`, `arrange`, `write_as_library` (experimental), `export_pdf_preset`, `import_pdf_preset`, `export_print_preset`, `import_print_preset` |
| `artboard.*` | `list`, `create`, `resize`, `set_active`, `fit_to_artwork`, `rearrange` |
| `layer.*` | `list`, `create`, `rename`, `visibility`, `lock`, `reorder` |
| `selection.*` | `get`, `clear`, `set_by_ids`, `select_by_name`, `select_active_artboard_objects` |
| `page_item.*` | `remove`, `duplicate`, `move`, `resize`, `rotate`, `transform`, `translate`, `z_order`, `bring_in_perspective` (experimental) |
| `path.*` | `create_rect`, `create_ellipse`, `create_path`, `transform`, `duplicate`, `create_polygon`, `create_star`, `create_rounded_rect`, `set_entire_path` |
| `text.*` | `create`, `set_contents`, `set_style`, `outline`, `create_area`, `create_on_path`, `thread`, `convert_to_area`, `convert_to_point`, `change_case` |
| `appearance.*` | `set_fill`, `set_stroke`, `set_gradient`, `apply_graphic_style` |
| `action.*` | `load`, `run`, `unload` |
| `export.*` | `png`, `jpg`, `svg`, `pdf`, `ai`, `gif`, `png8`, `tiff`, `webp`, `photoshop`, `autocad`, `eps`, `fxg` |
| `inspect.*` | `tree`, `styles`, `bounds`, `fonts`, `summary` |
| `workspace.*` | `save`, `switch`, `reset`, `delete` |
| `preference.*` | `get`, `set`, `delete` |
| `view.*` | `info`, `set_screen_mode`, `set_zoom`, `set_ruler_visibility`, `set_transparency_grid_visibility`, `set_center`, `rotate` |
| `matrix.*` | `identity`, `rotation`, `scale`, `translation`, `concatenate`, `concatenate_rotation`, `concatenate_scale`, `concatenate_translation`, `invert`, `compare`, `singular` |
| `perspective.*` | `show`, `hide`, `get_active_plane`, `set_active_plane`, `select_preset`, `import_preset`, `export_preset` |
| `style.character.*` | `list`, `apply`, `import` |
| `style.paragraph.*` | `list`, `apply`, `import` |
| `style.graphic.*` | `list`, `apply`, `merge`, `remove` |
| `swatch.*` | `list`, `create`, `delete` |
| `spot.*` | `list`, `create`, `delete` |
| `symbol.*` | `list`, `create`, `place`, `break_link` |
| `placed.*` | `list`, `place`, `embed`, `relink`, `trace` |
| `raster.*` | `list`, `rasterize`, `trace`, `release_tracing`, `colorize` |
| `repeat.grid.*` | `list`, `create`, `update` |
| `repeat.radial.*` | `list`, `create`, `update` |
| `repeat.symmetry.*` | `list`, `create`, `update` |
| `dataset.*` | `list`, `create`, `apply`, `update`, `delete`, `import`, `export` |
| `variable.*` | `list`, `create`, `delete`, `bind_visibility`, `bind_text`, `bind_content`, `import`, `export` |
| `trace.preset.*` | `list`, `store` (experimental) |
| `capture.*` | `image`, `window` |
| `print.*` | `presets`, `devices`, `run` |

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
- Opaque identifiers reject URL, query, fragment, and encoded traversal syntax.
- Nested payloads such as `appearance.set_gradient.stops` and `path.create_path.points` are validated structurally before execution.
- `document.new` enforces cross-field rules for `artboards`, `artboardLayout`, and `artboardRowsOrCols` so impossible payloads fail validation instead of hanging Illustrator.

## Live-Validated Script Notes

- `app.info` returns `name`, `version`, `scriptingVersion`, document counts, selection count, and the installed startup preset list.
- `app.preset_lists` still returns preset, tracing, and color-settings lists, but printer and PPD enumeration degrade to `unavailable` on this host instead of hanging Illustrator.
- `document.new` supports `width`, `height`, `artboards`, `artboardLayout`, `artboardSpacing`, `artboardRowsOrCols`, `colorSpace`, and optional `preset`.
- `document.save_as` honors `format` and now supports `ai` and `pdf` without relying on extension guessing alone.
- `document.arrange` is normalized through the stable `cascade` and `tile` menu commands because the documented `doc.arrange(...)` surface is not live on Illustrator 30.2.1.
- `document.write_as_library` is shipped as an experimental command. The CLI surfaces it with a warning because Illustrator 30.2.1 can reject documented library types at runtime.
- `path.create_rect` uses Illustrator's real `roundedRectangle(...)` API when `cornerRadius` is greater than zero.
- `view.set_ruler_visibility` and `view.set_transparency_grid_visibility` are normalized setter commands built on top of Illustrator menu toggles rather than leaking toggle semantics to the caller.
- `symbol.break_link` is script-backed through the live `SymbolItem.breakLink()` runtime method.
- `export.png` and `export.jpg` honor `artboardId` by activating the requested artboard and enabling artboard clipping.
- `export.svg` honors `artboardId`, but Illustrator 30.2.1 emits the file as `basename_<Artboard Name>.svg` when multi-artboard export is enabled. The command result returns that actual output path.
- `export.tiff` may materialize as either `basename.tif`, `basename.tiff`, or `basename-01.tif`. The CLI resolves the actual created artifact path before returning.
- `capture.window` can emit `basename.tif.tiff` on this host. The CLI normalizes the returned `outputPath` to the actual file that appeared on disk.
- `document.export_pdf_preset` and `document.export_print_preset` return the real exported artifact path because Illustrator can replace the requested filename with the preset name.
- `print.devices` intentionally reports a deterministic degraded result because live `printerList` and `PPDFileList` enumeration on this host was unreliable enough to hang Illustrator.
- `trace.preset.list` can legitimately return an empty array when the host exposes no script-visible tracing presets beyond the defaults.
- `page_item.bring_in_perspective` and `trace.preset.store` are shipped as experimental commands. Both return structured warnings because Illustrator 30.2.1 can reject the documented runtime path.

## No Plugin / No SDK

- The current command surface is script-backed and works without the native plugin.
- End users do not need the Illustrator SDK to use the CLI.
- Maintainers only need the SDK if they want to build the optional native plugin from source.
- `host status` and `doctor` may still show the native plugin probe as unavailable; that does not block the current CLI surface.
- The current script-first exclusions are:
  - `workspace.list`
  - `export.for_screens`
- The current script-first experimental commands are:
  - `document.write_as_library`
  - `page_item.bring_in_perspective`
  - `trace.preset.store`
  These are documented with runtime notes in [live-script-runtime-learnings.md](live-script-runtime-learnings.md) and emit `EXPERIMENTAL_COMMAND` warnings in the response envelope.

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

## Native Bridge Probe

The current CLI surface is script-backed for inspection, gradients, and graphic style application.

`host status` and `doctor` still probe the optional native bridge with `ahd.version`. A healthy installed bridge should return a JSON payload for:

```text
app.sendScriptMessage("AHDIllustrator", "ahd.version", "{}")
```

This probe is diagnostic only. It is reserved for future native-only capabilities rather than the current Phase 1 script-first command surface.
