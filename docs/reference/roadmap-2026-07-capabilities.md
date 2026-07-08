# 2026-07 Roadmap Capabilities

Implemented surfaces from the competitor-beating roadmap:

## Agent CLI

- `agent-context --json`
- `command --stdin`
- `command --payload`
- `command --payload-file`
- `command --fields`
- `command --deliver stdout|file:<path>|dir:<path>`
- `batch --stdin`
- `batch --payload`
- `batch --payload-file`
- `batch --dry-run`
- `batch --compact`
- structured stderr errors: `code`, `message`, `hint`, `retryable`

## Proof Gates

- `doctor --json`
- `verify plugin`
- `verify syntax`
- `verify live`
- `verify release`
- `document.screenshot`
- `document.screenshot_selection`
- `verify.visual`

## Reliability

- schema safety metadata enforcement
- contract drift gate: `make verify-contracts`
- MCP prompts: `prompts/list`, `prompts/get`
- central plugin `figma.commitUndo()` with `commitUndo:false`
- durable local jobs: `jobs list|get|resume|cancel`
- profile/config inspection: `profile list|use|inspect --redacted`, `config sources`

## Current Figma API Coverage

- Motion guards: `motion.get_styles`, `motion.apply_style`, `motion.remove_style`, `motion.get_animations`, `motion.apply_keyframes`, `motion.remove_keyframes`, `motion.set_timeline_duration`
- Shader guards: `effect.list_shaders`, `effect.import_shader`, `effect.apply_shader_effect`, `fill.apply_shader`
- Slots: `component.create_slot`, `component.reset_slot`, `component.get_slots`
- Grid updates: `layout.reorder_grid_rows`, `layout.reorder_grid_columns`, grid auto-flow readback
- Noise vectors: `noiseSize`, `noiseSizeX`, `noiseSizeY`, `noiseSizeVector`

## Design-System Moat

- `tokens.preset_tailwind`
- `tokens.preset_shadcn`
- `tokens.preset_material`
- `tokens.setup_system`
- `design_system.health`
- `component.analyze_set`
- `component.arrange_set`
- `parity.audit_component`

## Editor Coverage

- `document.get_editor_context`
- Slides guards: `slides.get_current`, `slides.create_slide`, `slides.set_background`, `slides.add_text`, `slides.reorder`
- FigJam guards: `figjam.create_sticky`, `figjam.create_shape`, `figjam.create_connector`, `figjam.get_board`

## Packaging

- `packaging.generate_skills`
- npm wrapper investigation: `docs/research/npm-wrapper-investigation.md`
- cloud relay proposal: `docs/design/cloud-relay.md`
- showcase batch: `docs/examples/roadmap-showcase-social-posts.json`
- exported artifacts: `docs/generated/*roadmap-showcase.png`
