# Atlas Automata

Atlas Automata is a small Git-native persistence framework for AI agents. The agent understands and configures the user's domain; the Atlas MCP server is the exclusive writer for protected project state.

## Start here

Work with Atlas through your coding agent. In a Git repository, tell the agent:

> Install Atlas Automata from https://github.com/equilaterus/atlas-automata and guide me through the complete setup for what I want to manage.

The agent must perform the technical installation, commit and push it, tell you when an agent restart is required, and then conduct the configuration as a guided conversation. You describe your goals in ordinary language; the agent proposes the structure and explains meaningful choices.

The guided setup always covers, in order:

1. purpose, users, boundaries, privacy, and terminology;
2. collections, relationships, statuses, and lifecycles;
3. how information must be found, filtered, grouped, and indexed;
4. record identity, paths, schemas, taxonomy, and validation;
5. views, calculations, imports, exports, and recurring workflows;
6. compatibility with existing information;
7. a complete proposal for your approval;
8. creation of the configuration, domain documentation, indexing rules, and domain skill;
9. final verification and handoff.

No phase may be silently skipped. A non-applicable phase is recorded with its reason. Atlas rejects writes under `data/` until the setup is structurally complete, so adding the first record is never the first step.

An approved evolution that changes physical data paths uses the restricted `atlas_setup: migration` state. In that state Atlas permits only `atlas_move` between data paths, updates folder indexes, and keeps create, update, and delete blocked until the migration is verified.

## What the agent installs

Atlas is pinned at `lib/atlas-automata` as a Git submodule. The agent runs this bootstrap from the child repository root:

```bash
git submodule add https://github.com/equilaterus/atlas-automata.git lib/atlas-automata
./lib/atlas-automata/run/install
git add .gitignore .gitmodules .agents .codex .githooks AGENTS.md ai lib/atlas-automata
ATLAS_MCP_COMMIT=1 git commit -m "Install Atlas Automata"
git push
```

The installer builds `.atlas/bin/atlas-mcp`, configures project-local Codex MCP access and lifecycle hooks, installs Git and agent guards, installs the base `configure` and `mutate` skills, exposes protected skills through `.agents/skills`, and adds the mandatory Atlas block to the child `AGENTS.md` without replacing existing project instructions. It ignores `.atlas/` and repository-local `tmp/` output.

After bootstrap, restart the agent client so it loads the project MCP configuration. At the start of that agent session, call `atlas_sync` once and then call `atlas_status`; do not repeat `atlas_sync` for each request. Any setup state other than `complete` requires the full `configure` workflow before domain data can be written.

## Clone or restore a child project

Tell the agent to restore Atlas, or run:

```bash
git submodule sync --recursive
git submodule update --init --recursive
./lib/atlas-automata/run/install
```

The child repository pins an exact Atlas commit. Restoring does not silently upgrade it.

## Update Atlas

Tell the agent to update and resynchronize Atlas, or run:

```bash
git submodule update --remote --merge lib/atlas-automata
./lib/atlas-automata/run/install
git add .gitignore .agents .codex AGENTS.md ai lib/atlas-automata
ATLAS_MCP_COMMIT=1 git commit -m "Update Atlas Automata"
git push
```

When editing Atlas itself inside a child checkout, commit and push inside `lib/atlas-automata` first. Then reinstall, commit the published submodule pointer in the child, and push the child repository. Never leave a child pointing to an unpublished Atlas commit.

## MCP tools

- `atlas_status` reports repository and setup state.
- `atlas_sync` fetches and merges the current `origin` branch.
- `atlas_create`, `atlas_update`, `atlas_delete`, and `atlas_move` mutate protected files, append semantic history when requested, commit, resynchronize, and push.
- Every data mutation also creates or refreshes the reserved `index.md` in each directory under `data/`, so the approved partition remains directly navigable without a separate generated index tree.

Protected state consists of `AUTOMATIZER.md` and files under `data/`, `doc/`, `ai/`, and `log/`. Domain-data mutations require a completed setup. See [the documentation index](doc/INDEX.md) for detailed contracts.

## Develop Atlas

```bash
./run/build
./run/test
./run/build-rel
```

Debug and release binaries are written to `bin/dbg/atlas-mcp` and `bin/rel/atlas-mcp`.
