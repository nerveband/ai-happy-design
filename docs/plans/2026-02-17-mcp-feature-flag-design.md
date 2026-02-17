# MCP Feature Flag Design

**Date**: 2026-02-17
**Status**: Approved

## Goal

Put the MCP server behind a feature flag so it's disabled by default. Users should use the CLI. MCP requires explicit opt-in via a config file.

## Config File

**Location**: `~/.config/ai-happy-design/config.toml`

```toml
# AI Happy Design configuration

[mcp]
# MCP server is disabled by default. Enable with:
#   ai-happy-design config set mcp.enabled true
# Or run 'ai-happy-design register' which enables it automatically.
enabled = false

[server]
port = 3055
host = "localhost"
idle_timeout = "15m"   # "0" or "off" to disable
```

## New `config` Subcommand

```
ai-happy-design config set <key> <value>   # Set a config value
ai-happy-design config get <key>            # Get a config value
ai-happy-design config path                 # Print config file path
ai-happy-design config init                 # Create default config file
```

Examples:
```bash
ai-happy-design config set mcp.enabled true
ai-happy-design config get mcp.enabled
# → true
```

## Behavior Changes

### `mcp` command
Checks `mcp.enabled`. If `false`, prints error and exits non-zero:
```
MCP server is disabled by default. To enable it:
  ai-happy-design config set mcp.enabled true
Or register with an editor (auto-enables MCP):
  ai-happy-design register
```

### `register` command
Before registering with editors, sets `mcp.enabled = true` in config automatically. Then proceeds as normal.

### `config.Load()`
Reads TOML file first, then env vars override. Existing env vars (`PORT`, `SERVER_HOST`, `AHD_IDLE_TIMEOUT`) continue to work and take precedence over the file.

### Config file auto-creation
On first `config set` or `register`, the file is created with defaults if it doesn't exist.

## Migration

- **Existing users with env vars**: Everything keeps working, no breakage.
- **Existing users who already ran `register`**: Their editor configs still point to `ai-happy-design mcp`, which will now error. They need to run `ai-happy-design config set mcp.enabled true` once. The error message tells them exactly what to do.
- **New users**: CLI works out of the box. MCP requires explicit opt-in.

## Dependencies

- `github.com/BurntSushi/toml` (standard Go TOML library)

## Key Decisions

- **TOML format**: Go-idiomatic, supports comments, self-documenting.
- **Error and exit** when MCP disabled: Clear, actionable message.
- **Register auto-enables**: Running `register` implies intent to use MCP. One step, not two.
- **Env vars override file**: Backwards compatible with existing deployments.
