# `src/mcp/`

The Atlas MCP server is implemented in Go.

Keep the project initially simple and flat.

Expected initial structure:

```text
src/mcp/
├── main.go
├── mcp.go
├── git.go
├── repo.go
├── skills.go
├── operations.go
├── history.go
├── build.go
├── go.mod
└── go.sum
```

Do not create package hierarchies merely to organize a small number of files.

Do not introduce:

```text
cmd/
internal/
pkg/
services/
repositories/
providers/
managers/
factories/
adapters/
```

unless actual code complexity eventually makes one of them clearly necessary.

For 0.1, keep `src/mcp/` simple.

---

# Go programming rules

Use plain, straightforward Go.

Favor:

```text
functions
plain structs
explicit parameters
explicit return values
composition
direct control flow
```

---

# NO INTERFACES

Do not introduce Go interfaces.

Atlas Automata 0.1 does not use interfaces.

Do not create interfaces for:

* testing
* mocking
* Git
* repositories
* MCP operations
* filesystem access
* future implementations
* dependency injection
* architectural cleanliness

Use concrete functions and concrete data.

---

# NO ABSTRACTIONS

Do not introduce abstractions.

This is intentional.

Do not build generalized frameworks around code that is currently simple.

Avoid:

```text
RepositoryManager
GitService
OperationExecutor
DataProvider
StorageBackend
FileRepository
SkillRegistryService
BuildPipeline
EventBus
PluginManager
Factory
Provider
Adapter hierarchy
```

unless a later real requirement forces the design to change.

For Atlas 0.1:

> Solve the concrete problem directly.

Prefer:

```go
findRepoRoot()
gitFetch()
gitMerge()
gitStatus()
createFile()
updateFile()
deleteFile()
moveFile()
loadSkills()
writeHistory()
```

If two functions contain some similar code, that is acceptable.

Do not generalize code merely because duplication exists.

Duplicated simple code is preferable to a premature abstraction.

---

# Composition

Use simple composition.

Example:

```text
syncRepo()
    ↓
checkPath()
    ↓
applyMutation()
    ↓
validateRepo()
    ↓
commitChanges()
    ↓
syncRepo()
    ↓
push()
```

Each function should perform a concrete recognizable task.

Avoid hidden behavior.

---

# Functions

Functions should normally:

* have one clear purpose
* accept explicit inputs
* return explicit results or errors
* expose side effects clearly
* remain easy to read top-to-bottom

Do not artificially split coherent logic into dozens of tiny functions.

Optimize for comprehension, not line-count metrics.

---

# Data structures

Use plain Go structs where structured state is useful.

Example:

```go
type SyncState struct {
    Branch     string
    LocalHead  string
    RemoteHead string
}
```

Do not attach methods merely to simulate classes.

A function is usually preferable:

```go
syncRepo(root string) error
```

instead of:

```go
repo.Sync()
```

when no meaningful object behavior exists.

---

# Global state

Avoid mutable global state.

Constants are fine.

Runtime information should normally be passed explicitly or loaded when needed.

---

# Errors

Use normal Go errors.

Add useful context.

Example:

```go
if err != nil {
    return fmt.Errorf("fetch origin: %w", err)
}
```

Do not build custom error hierarchies unless an actual requirement appears.

---

# Concurrency

Default to synchronous execution.

Do not use goroutines just because Go provides them.

Concurrency is appropriate only when:

* operations are independent
* they are meaningfully slow
* concurrent execution clearly improves behavior

Git mutations must remain sequential.

Protected repository mutations must remain sequential.

Do not build complex channel-based architectures.

Use async/concurrency only where it demonstrably helps.

---

# Dependencies

Keep dependencies extremely small.

Prefer the Go standard library.

Use the official Go MCP SDK.

Do not implement MCP or JSON-RPC manually.

Do not add dependencies for trivial functionality that can be implemented clearly using the standard library.

---

# MCP responsibility

Atlas MCP is the controlled write gateway for an Automatizer repository.

Its responsibilities are intentionally narrow:

```text
synchronize Git
inspect repository state
validate target paths
create files
update files
delete files
move files
validate repository workflow
record operation information
commit
synchronize again
push
```

It is not a domain engine.

It is not an AI.

It is not a query engine.

It is not an ontology engine.

---

# MCP write operations

The initial MCP API should remain small.

Conceptually:

```text
atlas.sync

atlas.create
atlas.update
atlas.delete
atlas.move

atlas.status
```

Exact MCP tool names may change during implementation.

Do not expose an unnecessarily large API.

---