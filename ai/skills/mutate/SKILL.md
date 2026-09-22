---
name: atlas-mutate
description: Read and persist configured-domain knowledge in an Atlas child project. Use for every request about the project's domain, including recommendations and casual discussion, and whenever the user supplies facts, preferences, ratings, corrections, decisions, or other information that may update protected Automatizer state.
---

# Atlas mutation

Treat every request about the configured domain as domain work. Read before answering even when the request appears conversational. User-provided facts, preferences, ratings, corrections, and decisions are candidate mutations; preserve them when they fit the configured model instead of merely acknowledging them in chat.

Call `atlas_status` first. If `setup` is not `complete`, stop normal domain work and run the complete `configure` skill. Never attempt to write under `data/` before configuration is verified.

Read `AUTOMATIZER.md`, the relevant files under `doc/`, installed domain skills, and existing data before deciding a mutation.

Use only the Atlas MCP tools for files under `data/`, `doc/`, `ai/`, `log/`, or `AUTOMATIZER.md`:

- `atlas_create` for a new file.
- `atlas_update` to replace an existing file with complete content.
- `atlas_delete` for an existing file.
- `atlas_move` to rename or relocate an existing file.

Always supply a concise semantic `summary` when the change represents a user or domain action. Do not reinterpret historical data silently. If meaning is ambiguous, ask the user before calling a mutation tool.
