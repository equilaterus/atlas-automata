# Installation, submodule lifecycle, build, and guards

## Child installation

Atlas must be pinned as a Git submodule at `lib/atlas-automata`. From the child repository root:

```bash
git submodule add https://github.com/equilaterus/atlas-automata.git lib/atlas-automata
./lib/atlas-automata/run/install
git add .gitignore .gitmodules .agents .codex .githooks AGENTS.md ai lib/atlas-automata
ATLAS_MCP_COMMIT=1 git commit -m "Install Atlas Automata"
git push
```

Installation activates the protected-state pre-commit hook before the bootstrap files under `ai/` are committed. `ATLAS_MCP_COMMIT=1` authorizes that one trusted bootstrap commit. The installer adds `.atlas/` and `tmp/` to the child `.gitignore`; generated binaries and repository-local temporary work must not be committed. The operating system's absolute `/tmp` directory is outside the child repository and cannot be tracked by it.

The installer verifies a normal Git worktree and the expected framework location, builds a release binary, and atomically installs it at `.atlas/bin/atlas-mcp` so a running MCP can be upgraded. It configures `core.hooksPath=.githooks` and creates a project-local Codex entry in `.codex/config.toml`:

```toml
[mcp_servers.atlas]
command = "/absolute/project/path/.atlas/bin/atlas-mcp"
args = ["--root", "/absolute/project/path"]
```

The agent guards are copied into `ai/hooks/` without replacing different project-owned hooks. The installer registers the Atlas `SessionStart` and `UserPromptSubmit` lifecycle hooks in `.codex/hooks.json`; it updates an Atlas-owned hook file and stops rather than replace different project-owned Codex hooks. An existing Atlas MCP entry is preserved. Missing base skills, including mandatory `configure`, are copied into `ai/skills/`; any project-owned skill with the same name is preserved. A repository skill-discovery link at `.agents/skills` exposes that protected skill directory to Codex without duplicating it. The installer creates or refreshes only the marker-delimited Atlas block in the child `AGENTS.md`, preserving all project-owned instructions outside it and rejecting malformed markers. Installation is the trusted bootstrap step that creates initial protected capabilities before MCP enforcement is active.

After the bootstrap is committed and pushed, the user restarts the agent client so it loads the MCP configuration. The agent calls `atlas_sync` and `atlas_status`, then completes the guided workflow in `configure` whenever setup is not `complete`. Domain-data writes are technically blocked until that workflow produces the required artifacts. See [configuration.md](configuration.md).

If a target Git hook already exists with different content, installation stops instead of replacing project behavior.

## Restore, resynchronize, and update the submodule

After cloning a child repository, synchronize submodule URLs and restore its pinned revision before running the installer:

```bash
git submodule sync --recursive
git submodule update --init --recursive
./lib/atlas-automata/run/install
```

Updating Atlas is explicit. Update from the submodule remote, reinstall the generated binary, then commit and push the new submodule pointer in the child repository:

```bash
git submodule update --remote --merge lib/atlas-automata
./lib/atlas-automata/run/install
git add .gitignore .agents .codex AGENTS.md ai lib/atlas-automata
ATLAS_MCP_COMMIT=1 git commit -m "Update Atlas Automata"
git push
```

When developing Atlas inside a child checkout, the two repositories must be published in dependency order:

```bash
git -C lib/atlas-automata add <atlas-files>
git -C lib/atlas-automata commit -m "Describe the Atlas change"
git -C lib/atlas-automata push
./lib/atlas-automata/run/install
git add .gitignore .agents .codex AGENTS.md ai lib/atlas-automata
ATLAS_MCP_COMMIT=1 git commit -m "Update Atlas Automata"
git push
```

The child repository must never point to an Atlas commit that has not been pushed to the submodule remote.

## Developer commands

```bash
./run/build       # debug by default
./run/build-dbg   # bin/dbg/atlas-mcp
./run/build-rel   # bin/rel/atlas-mcp
./run/test        # functional suite
./run/clean       # remove generated binaries
```

Builds inject the version from `VERSION`. Generated binaries are not source-of-truth.

## Guards

- `pre-commit` rejects staged protected paths unless the commit was created by Atlas MCP.
- `pre-rebase` rejects every rebase.
- `pre-push` rejects branch deletion and non-fast-forward updates, including force-push.
- `session-start` runs the binary's one-shot merge-only synchronization and stops the session on failure.
- `domain-context` reinforces Atlas domain classification on every user prompt and after compaction.
- `sync-before-work` remains available as a host-neutral one-shot synchronization helper.
- `guard-command` is a small agent-hook helper that rejects obvious rebase and force-push commands before execution.

Hooks are guards, not business logic. The MCP owns the actual write transaction. The Codex lifecycle hooks synchronize only at session start; prompt and compaction hooks inject domain context without synchronizing or mutating state. `guard-command` remains available for pre-command lifecycle events on hosts that expose them.

## Operational notes

The server uses stdio and writes protocol messages only to stdout; diagnostics go to stderr. `--root` pins a server process to one child repository, `--sync` performs a one-shot sync for lifecycle hooks, and `--version` prints the build version.
