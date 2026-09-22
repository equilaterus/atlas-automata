# Functional testing contract

Atlas 0.2 uses functional/integration tests, not unit tests, mocks, interface-based fakes, or coverage targets. Tests create real temporary directories, bare remotes, clones, branches, commits, merges, files, hooks, and MCP sessions using the official SDK.

## Required scenarios

1. **Guided setup gate:** reject data before setup, reject a premature completion marker, permit configuration artifacts, and permit data only after setup becomes complete.
2. **Basic mutation:** call MCP to create protected data; verify the file, commit, semantic history when supplied, and remote commit.
3. **Mutation lifecycle:** exercise create, update, move, and delete through MCP and finish clean.
4. **Sync before mutation:** advance origin from clone A; mutate from clone B; verify B incorporates A first.
5. **Divergence:** commit independently in A and B; push A; mutate from B; verify an actual merge commit and no rebase.
6. **Merge conflict:** create a genuine conflicting commit in both clones; verify Atlas reports failure and does not apply the requested mutation.
7. **Remote change during operation:** advance origin from another clone during the Atlas commit; verify the final fetch detects and merges it before push.
8. **Rejected push race:** advance origin immediately before push; verify rejection, fetch, merge, revalidation, retry, and success.
9. **Rebase rejection:** verify both `git rebase` and `git pull --rebase` fail in a configured child environment.
10. **Force-push rejection:** verify a non-fast-forward `git push --force` fails.
11. **Direct protected write:** stage a direct `data/` change and verify the commit guard rejects it.
12. **Domain evolution:** preserve historical `running` records while documenting `running is_a exercise` and compatibility; verify both running and exercise views remain possible.
13. **Generated dashboard:** generate `bin/dashboard.html`, delete `bin/`, regenerate, and verify source data is unchanged.
14. **Folder indexes:** create, move, and delete nested data records; verify Atlas maintains linked `index.md` files in every data directory and rejects direct index mutation.
15. **Safe migration:** verify `atlas_setup: migration` blocks data create, update, and delete while allowing data-to-data moves, then returns cleanly to complete.
16. **MCP safety metadata:** list the six tools through an MCP client and verify their read-only, destructive, idempotent, and open-world annotations exactly.
17. **Codex activation:** install Atlas into a real temporary Git repository; verify MCP configuration, lifecycle hooks, executable hook scripts, discovered skill link, and replacement of stale managed instructions without changing project-owned instructions.

## Acceptance

`./run/test` must execute the real workflow locally without network services. Tests may use local filesystem remotes. A few complete scenarios are preferred to many isolated assertions. Failures must expose Git/MCP context sufficiently to diagnose the transaction stage.
