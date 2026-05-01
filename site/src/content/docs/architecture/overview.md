---
title: Architecture Overview
description: CLI, relay, plugin, schema, and MCP architecture.
---

AI Happy Design has four main layers.

## CLI

The Go binary exposes command, batch, schema, validation, MCP, guide, and relay workflows.

## Schema and Validation

Schemas live in the Go codebase and power CLI discovery, validation, batch correction, and MCP tool definitions.

## Relay

The local WebSocket relay connects the CLI and AI agents to the active Figma plugin channel.

## Figma Plugin

The plugin runs inside Figma, receives schema-routed commands, writes canvas nodes, reads document state, and returns structured response envelopes.

The plugin bundle is built for ES6 because the Figma plugin sandbox does not support newer JavaScript syntax reliably.
