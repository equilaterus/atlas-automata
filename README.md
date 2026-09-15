# Atlas Automata

Atlas Automata is a small Git-native persistence framework for AI agents. The agent reads and understands a domain; the Atlas MCP server is the exclusive writer for protected project state and applies a fetch/merge/commit/fetch/merge/push workflow.

## Install first

Atlas is not copied into a child project. Add it as a pinned Git submodule at the exact `lib/atlas-automata` path, run the installer from the child repository root, then commit and push the bootstrap:

```bash
git submodule add https://github.com/equilaterus/atlas-automata.git lib/atlas-automata
./lib/atlas-automata/run/install
git add .gitignore .gitmodules .codex .githooks ai lib/atlas-automata
ATLAS_MCP_COMMIT=1 git commit -m "Install Atlas Automata"
git push
```

The environment override on the bootstrap commit is required because installation activates the protected-state pre-commit guard. The installer adds `.atlas/` and `tmp/` to the child `.gitignore`: generated runtime files and repository-local temporary work are not committed. The operating system's absolute `/tmp` directory is already outside the repository and is never tracked by the child project.

After cloning a child repository, restore the exact pinned framework revision before installing:

```bash
git submodule sync --recursive
git submodule update --init --recursive
./lib/atlas-automata/run/install
```

To update Atlas deliberately, update the submodule, reinstall the generated binary, then commit and push the new pointer from the child repository:

```bash
git submodule update --remote --merge lib/atlas-automata
./lib/atlas-automata/run/install
git add lib/atlas-automata
git commit -m "Update Atlas Automata"
git push
```

The installer builds `.atlas/bin/atlas-mcp`, configures project-local Codex MCP access, installs Git guards under `.githooks`, and copies missing base skills without replacing project-owned skills.

If Atlas itself is edited from inside a child checkout, commit and push inside `lib/atlas-automata` first. Then return to the child repository, reinstall, commit the changed submodule pointer, and push the child repository. Never leave the child pointing to an unpublished Atlas commit.

## Build and test

```bash
./run/build
./run/test
./run/build-rel
```

The debug and release binaries are written to `bin/dbg/atlas-mcp` and `bin/rel/atlas-mcp`.

## MCP tools

- `atlas_status` inspects repository state.
- `atlas_sync` fetches and merges the current `origin` branch.
- `atlas_create`, `atlas_update`, `atlas_delete`, and `atlas_move` mutate protected files, optionally append semantic history, commit, resynchronize, and push.

Protected state consists of `AUTOMATIZER.md` and files under `data/`, `doc/`, `ai/`, and `log/`. See [the documentation index](doc/INDEX.md) for the architecture and detailed contracts.
