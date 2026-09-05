# Functional testing contract

Atlas 0.1 uses functional/integration tests, not unit tests, mocks, interface-based fakes, or coverage targets. Tests create real temporary directories, bare remotes, clones, branches, commits, merges, files, hooks, and MCP sessions using the official SDK.

## Required scenarios

1. **Basic mutation:** call MCP to create protected data; verify the file, commit, semantic history when supplied, and remote commit.
2. **Mutation lifecycle:** exercise create, update, move, and delete through MCP and finish clean.
3. **Sync before mutation:** advance origin from clone A; mutate from clone B; verify B incorporates A first.
4. **Divergence:** commit independently in A and B; push A; mutate from B; verify an actual merge commit and no rebase.
5. **Merge conflict:** create a genuine conflicting commit in both clones; verify Atlas reports failure and does not apply the requested mutation.
6. **Remote change during operation:** advance origin from another clone during the Atlas commit; verify the final fetch detects and merges it before push.
7. **Rejected push race:** advance origin immediately before push; verify rejection, fetch, merge, revalidation, retry, and success.
8. **Rebase rejection:** verify both `git rebase` and `git pull --rebase` fail in a configured child environment.
9. **Force-push rejection:** verify a non-fast-forward `git push --force` fails.
10. **Direct protected write:** stage a direct `data/` change and verify the commit guard rejects it.
11. **Domain evolution:** preserve historical `running` records while documenting `running is_a exercise` and compatibility; verify both running and exercise views remain possible.
12. **Generated dashboard:** generate `bin/dashboard.html`, delete `bin/`, regenerate, and verify source data is unchanged.

## Acceptance

`./run/test` must execute the real workflow locally without network services. Tests may use local filesystem remotes. A few complete scenarios are preferred to many isolated assertions. Failures must expose Git/MCP context sufficiently to diagnose the transaction stage.
