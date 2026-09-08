# Atlas Automata

Atlas Automata is a small Git-native persistence framework for AI agents. The agent reads and understands a domain; the Atlas MCP server is the exclusive writer for protected project state and applies a fetch/merge/commit/fetch/merge/push workflow.

## Build and test

```bash
./run/build
./run/test
./run/build-rel
```

The debug and release binaries are written to `bin/dbg/atlas-mcp` and `bin/rel/atlas-mcp`.

## Install in a child repository

Add Atlas at `lib/atlas-automata`, then run the installer from the child repository root:

```bash
./lib/atlas-automata/run/install
```

The installer builds `.atlas/bin/atlas-mcp`, configures project-local Codex MCP access, installs Git guards under `.githooks`, and copies missing base skills without replacing project-owned skills.

## MCP tools

- `atlas_status` inspects repository state.
- `atlas_sync` fetches and merges the current `origin` branch.
- `atlas_create`, `atlas_update`, `atlas_delete`, and `atlas_move` mutate protected files, optionally append semantic history, commit, resynchronize, and push.

Protected state consists of `AUTOMATIZER.md` and files under `data/`, `doc/`, `ai/`, and `log/`. See [the documentation index](doc/INDEX.md) for the architecture and detailed contracts.
