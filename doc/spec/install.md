# Build, installation, and guards

## Developer commands

```bash
./run/build       # debug by default
./run/build-dbg   # bin/dbg/atlas-mcp
./run/build-rel   # bin/rel/atlas-mcp
./run/test        # functional suite
./run/clean       # remove generated binaries
```

Builds inject the version from `VERSION`. Generated binaries are not source-of-truth.

## Child installation

With Atlas pinned at `lib/atlas-automata`, run from the child root:

```bash
./lib/atlas-automata/run/install
```

The installer verifies a normal Git worktree and expected framework location, builds a release binary, and installs it at `.atlas/bin/atlas-mcp`. It configures `core.hooksPath=.githooks` and creates a project-local Codex entry in `.codex/config.toml`:

```toml
[mcp_servers.atlas]
command = "/absolute/project/path/.atlas/bin/atlas-mcp"
args = ["--root", "/absolute/project/path"]
```

The agent guards `sync-before-work` and `guard-command` are copied into `ai/hooks/` without replacing different project-owned hooks. An existing Atlas MCP entry is preserved. Missing base skills are copied into `ai/skills/`; any project-owned skill with the same name is preserved. Installation is the trusted bootstrap step that creates initial protected capabilities before MCP enforcement is active.

If a target Git hook already exists with different content, installation stops instead of replacing project behavior.

## Guards

- `pre-commit` rejects staged protected paths unless the commit was created by Atlas MCP.
- `pre-rebase` rejects every rebase.
- `pre-push` rejects branch deletion and non-fast-forward updates, including force-push.
- `sync-before-work` runs the binary's one-shot merge-only synchronization.
- `guard-command` is a small agent-hook helper that rejects obvious rebase and force-push commands before execution.

Hooks are guards, not business logic. The MCP owns the actual write transaction. Agent hosts should wire `sync-before-work` to session/prompt startup and `guard-command` to pre-command lifecycle events when those surfaces are available.

## Operational notes

The server uses stdio and writes protocol messages only to stdout; diagnostics go to stderr. `--root` pins a server process to one child repository, `--sync` performs a one-shot sync for lifecycle hooks, and `--version` prints the build version.
