# Atlas MCP 0.2

## Scope and source layout

The server is a local stdio MCP implemented in Go with the official Go MCP SDK. Source remains flat under `src/mcp/`:

```text
main.go        process entry and one-shot sync
mcp.go         tool registration and schemas
repo.go        repository and path validation
git.go         explicit Git commands and merge-only synchronization
operations.go  concrete file mutations
history.go     semantic history append
skills.go      installed skill discovery
build.go       build-injected version
functional_test.go
```

Do not add `cmd/`, `internal/`, `pkg/`, service/manager/provider/factory/adapter layers, interfaces, mocks, or generalized Git/filesystem abstractions for 0.2.

## MCP API

The server exposes six tools:

| Tool | Behavior |
| --- | --- |
| `atlas_status` | Reports root, branch, HEAD, cleanliness, setup state, and installed skills without fetching. |
| `atlas_sync` | Fetches `origin` and merges the current remote branch when necessary. |
| `atlas_create` | Creates one new protected regular file. |
| `atlas_update` | Atomically replaces one existing protected regular file. |
| `atlas_delete` | Deletes one existing protected regular file. |
| `atlas_move` | Moves one existing protected regular file to a new protected path. |

Mutation inputs use repository-relative paths. Create/update take complete UTF-8 `content`. All mutations accept optional `summary` and `commit_message`. The summary is appended to UTC-dated Markdown under `log/` in the same commit.

## Protected paths

Only these paths are writable through mutation tools:

```text
AUTOMATIZER.md
data/**
doc/**
ai/**
log/**
```

Paths must be relative, remain inside the repository, and not traverse symlinked parents. Targets are files, not arbitrary directory trees. Atlas rejects a mutation if protected state already has staged, unstaged, or untracked changes, preventing unrelated direct edits from entering an Atlas commit.

Mutations involving `data/` additionally require `setup: complete`. Atlas derives that state from the `AUTOMATIZER.md` front-matter marker, the standard model/indexing/operations documents, and at least one installed non-base domain skill. It rejects a completion marker before those artifacts exist. See [configuration.md](configuration.md).

## Mutation transaction

Every mutation is sequential and follows this exact order:

1. Verify repository root, attached branch, target path, setup gate, and absence of existing protected changes.
2. Fetch `origin`.
3. Compare local HEAD with `origin/<branch>` and merge the remote ref when it is not already an ancestor.
4. If that merge conflicts, stop before applying the mutation and leave the conflict visible for human/agent reconciliation.
5. Recheck protected cleanliness.
6. Apply exactly one create, update, delete, or move.
7. Append semantic history when requested.
8. Validate that no unresolved conflict exists.
9. Stage only the mutation paths and generated history path, then commit with `ATLAS_MCP_COMMIT=1` for the repository guard.
10. Fetch and merge again.
11. Revalidate and push `HEAD:refs/heads/<branch>`.
12. If push is rejected because origin advanced, fetch, merge, revalidate, and retry up to three times.

Atlas never invokes rebase or force-push. A post-commit merge conflict is reported and never silently resolved.

## Git behavior

Atlas invokes the user's installed `git` executable so credential helpers, SSH configuration, user identity, remotes, and normal Git configuration continue to work. Published history is immutable. A branch must be attached; detached HEAD is rejected. A mutation requires an `origin` remote because successful completion includes push.

## Implementation style

Use concrete functions such as `findRepoRoot`, `syncRepo`, `createFile`, and `writeHistory`. Functions should expose side effects, accept explicit inputs, return normal contextual errors, and remain readable top-to-bottom. Runtime state is passed or captured explicitly. Server mutation handlers share one concrete mutex so Git mutations cannot overlap.

Concurrency is not otherwise introduced. Similar short code may remain duplicated. The standard library is preferred except for the official MCP SDK.
