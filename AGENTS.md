# Atlas Automata agent rules

Atlas Automata is a Git-native controlled writer for Automatizer repositories. Keep this file limited to rules that must apply on every task. Load detailed specifications from [`doc/INDEX.md`](doc/INDEX.md) only when the task needs them.

## Always preserve these invariants

1. The agent understands the domain; Atlas MCP does not contain domain intelligence.
2. The agent may read the whole repository.
3. Protected state is written only through Atlas MCP.
4. `data/` is the domain database, `doc/` is authoritative domain meaning, and `ai/` contains installed agent capabilities.
5. Agent-led configuration is mandatory before domain data; no setup phase may be silently skipped.
6. Synchronize once at session start; every protected mutation synchronizes internally.
7. Use merge only. Never rebase, force-push, or silently resolve conflicts.
8. Never silently destroy or reinterpret historical information.
9. Keep Atlas concrete, small, and readable: no interfaces and no speculative abstractions.
10. Prefer functional tests with real Git repositories and real MCP calls.

## Start each agent session

Before the first request handled in an Automatizer child project session, run `atlas_sync` once. If the MCP is not yet available, run the installed `.atlas/bin/atlas-mcp --root <repo> --sync` guard. Stop on a merge conflict and report it; never guess a resolution.

Then run `atlas_status`. Do not repeat `atlas_sync` mechanically for later requests in the same session: protected mutation tools synchronize internally. When setup is not `complete`, load `ai/skills/configure/SKILL.md` and finish its full guided workflow before accepting domain data.

## Reading and writing

Direct reads are allowed everywhere, including:

```text
AUTOMATIZER.md
data/**
doc/**
ai/**
log/**
src/**
```

The following child-project state is protected:

```text
AUTOMATIZER.md
data/**
doc/**
ai/**
log/**
```

Never modify protected state with generic file tools, shell redirection, `sed -i`, or direct Git operations. Use exactly one of:

```text
atlas_create
atlas_update
atlas_delete
atlas_move
```

Read the relevant installed skills, domain documentation, and existing data first. Provide the MCP with an exact relative path and exact content. Add a concise semantic history summary when the mutation represents a user or domain action.

The child project's `src/`, `run/`, and generated `bin/` are ordinary project areas unless that project defines stricter rules. Framework source in this repository belongs under `src/`.

## Git policy

Git is part of the runtime, not only source control.

- Fetch `origin` and merge `origin/<current-branch>` when needed.
- Never run `git rebase`, `git pull --rebase`, or `git pull -r`.
- Never run `git push --force`, `git push -f`, or `git push --force-with-lease`.
- Treat published history as immutable.
- Keep protected mutations sequential.
- If a conflict occurs before mutation, do not apply the mutation.
- If a conflict occurs after the local mutation commit, stop and report it.
- If a push race occurs, fetch, merge, revalidate, and retry without rebasing.

Load [`doc/spec/mcp.md`](doc/spec/mcp.md) when changing the Git or mutation workflow.

## Go implementation rules

The MCP server lives flat under `src/mcp/` and uses the official Go MCP SDK. Use plain functions, plain structs, explicit parameters and return values, direct control flow, and the standard library where practical.

Do not introduce Go interfaces, mocks, dependency injection, manager/service/provider/factory/adapter layers, or package hierarchies for organization alone. Some simple duplication is preferable to a premature abstraction. Avoid mutable global runtime state. Add context to ordinary Go errors. Git mutations stay synchronous.

Load [`doc/spec/mcp.md`](doc/spec/mcp.md) before changing MCP behavior or source layout.

## Developer workflow

Use the repository scripts:

```bash
./run/build
./run/build-dbg
./run/build-rel
./run/test
./run/install
./run/clean
```

Debug output belongs in `bin/dbg/`; release output belongs in `bin/rel/`. Generated child-project artifacts belong in `bin/` and must be deletable and reproducible.

## Documentation routing

Load only the specification relevant to the current task:

- [`doc/spec/initial.md`](doc/spec/initial.md): purpose, architecture, responsibility split, and child-project model.
- [`doc/spec/repository.md`](doc/spec/repository.md): repository layout, protected data, domain evolution, views, and history.
- [`doc/spec/configuration.md`](doc/spec/configuration.md): mandatory agent-led setup, setup states, indexing, required artifacts, and evolution.
- [`doc/spec/mcp.md`](doc/spec/mcp.md): MCP API, path validation, exact Git transaction, failures, and implementation constraints.
- [`doc/spec/install.md`](doc/spec/install.md): build, installation, Codex configuration, skills, and guards.
- [`doc/spec/testing.md`](doc/spec/testing.md): required end-to-end scenarios and acceptance criteria.

Keep [`doc/INDEX.md`](doc/INDEX.md) current whenever specification pages are added, moved, renamed, or materially repurposed.

## Scope of 0.2

Prioritize repository detection, MCP startup, mandatory guided configuration, synchronization, merge-only behavior, protected file mutations, commit/push workflow, guards, functional tests, minimal skill installation, and one realistic child project flow.

Do not build a SaaS backend, database server, distributed locking, remote MCP service, OAuth, marketplace, ontology engine, generic workflow engine, event bus, domain-specific GUI, custom Git/filesystem abstraction, or a hand-written MCP/JSON-RPC implementation.
