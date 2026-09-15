---
name: atlas-mutate
description: Persist concrete changes to protected Automatizer state through Atlas MCP.
---

# Atlas mutation

Call `atlas_status` first. If `setup` is not `complete`, stop normal domain work and run the complete `configure` skill. Never attempt to write under `data/` before configuration is verified.

Read `AUTOMATIZER.md`, the relevant files under `doc/`, installed domain skills, and existing data before deciding a mutation.

Use only the Atlas MCP tools for files under `data/`, `doc/`, `ai/`, `log/`, or `AUTOMATIZER.md`:

- `atlas_create` for a new file.
- `atlas_update` to replace an existing file with complete content.
- `atlas_delete` for an existing file.
- `atlas_move` to rename or relocate an existing file.

Always supply a concise semantic `summary` when the change represents a user or domain action. Do not reinterpret historical data silently. If meaning is ambiguous, ask the user before calling a mutation tool.
