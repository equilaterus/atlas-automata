# Atlas Automata repository structure

The framework repository itself should follow:

```text
atlas-automata/
├── AGENTS.md
├── README.md
│
├── src/
│   └── mcp/
│
├── run/
│
├── bin/
│   ├── dbg/
│   └── rel/
│
├── ai/
│   ├── skills/
│   ├── hooks/
│   └── instructions/
│
├── doc/
│
└── VERSION
```

Do not create additional directory layers without a concrete current need.

---

# `src/`

All framework source code belongs under:

```text
src/
```

Do not place implementation source files directly in the repository root.

---

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

# Read access

The agent should normally read the repository using its normal filesystem tools.

Do not unnecessarily proxy all reads through MCP.

Codex must be able to directly read:

```text
data/**
doc/**
ai/**
log/**
AUTOMATIZER.md
src/**
```

This keeps exploration and reasoning simple.

---

# Protected writes

The agent may read protected project state directly.

It must not write protected project state directly.

Protected areas include:

```text
data/**
doc/**
ai/**
log/**
AUTOMATIZER.md
```

Changes to these areas must be performed through Atlas MCP.

---

# `data/`

`data/` is the domain database.

Although Markdown is the primary physical storage format in 0.1, it must be treated as persistent domain data.

Examples:

```text
data/
├── clients/
├── transactions/
├── accounts/
├── activities/
├── projects/
└── people/
```

There is no universal physical layout.

The domain determines whether a collection is organized:

* by entity
* by date
* by project
* by entity and date
* another useful structure

The agent understands that structure through skills and documentation.

---

# Data writes

The agent may inspect and reason over all files in:

```text
data/
```

but must never directly modify them.

The correct flow is:

```text
agent understands desired mutation
        ↓
agent determines exact path/content/change
        ↓
agent calls Atlas MCP
        ↓
MCP performs safe write workflow
```

---

# `doc/`

`doc/` contains authoritative domain meaning.

Examples:

```text
doc/
├── domain/
├── views/
├── compatibility/
├── decisions/
└── rules/
```

This may include:

* concepts
* schemas
* taxonomy
* relationships
* compatibility information
* view definitions
* domain decisions
* rules

The agent may read it freely.

Direct writes are forbidden.

Changes go through MCP.

---

# `ai/`

`ai/` contains capabilities actually installed in a child project.

Example:

```text
ai/
├── manifest.yaml
├── skills/
├── hooks/
└── instructions/
```

The agent may read these files.

Direct writes are forbidden.

---

# Skills

Skills tell the agent how to understand and interact with the domain.

Skills are intelligence/context for the agent.

They are not the writer.

A domain skill may explain:

```text
where records live
what fields mean
which fields are required
how entities relate
which historical data is compatible
when the user must be asked
how to interpret domain concepts
```

The agent uses that knowledge to construct the exact desired mutation.

Atlas MCP then writes it.

---

# Base skills

Atlas Automata supplies reusable base skills.

Initial candidates include:

```text
configure
sync
query
mutate
history
reconcile
validate
calculate
compatibility
consolidate
migrate
view
build
dashboard
review
```

Not all skills need to exist in the first implementation.

Do not implement a skill merely because it appears in this list.

Implement features only as they become required by the 0.1 workflow.

---

# Installed skills

The framework may contain available base skills.

A child project has its actual installed skills under:

```text
ai/skills/
```

The agent should reason only from the capabilities installed in that child project plus ordinary repository documentation.

---

# `AUTOMATIZER.md`

`AUTOMATIZER.md` describes the high-level configured Automatizer.

It may include:

* purpose
* users
* domain
* collections
* restrictions
* organization strategy
* calculation requirements
* project-specific workflow

The agent may read it directly.

Updates must occur through MCP.

---

# `configure`

`configure` helps the agent establish or evolve the domain.

It must be conservative.

If configuration, documentation or data already exists, the agent must inspect it before proposing changes.

Never treat an existing Automatizer as a blank project.

Expected reasoning flow:

```text
inspect existing configuration
        ↓
inspect existing domain documentation
        ↓
inspect existing data
        ↓
understand user request
        ↓
propose model evolution
        ↓
determine compatibility impact
        ↓
ask the user about semantic ambiguity
        ↓
construct concrete file mutations
        ↓
send mutations through Atlas MCP
```

---

# Existing data must not be silently changed

A new configuration must never silently destroy or reinterpret historical information.

Possible impact categories include:

```text
SAFE
COMPATIBLE
PARTIAL
CONVERTIBLE
AMBIGUOUS
BREAKING
```

The agent determines these semantically.

Atlas MCP only performs approved concrete changes.

---

# Compatibility inventory

Known semantic compatibility issues should be persisted.

Suggested location:

```text
doc/compatibility/inventory.yaml
```

Examples:

```text
old running data does not explicitly contain exercise=true

old customer records lack new classification

historical financial category no longer maps exactly
```

Views and dashboards should be able to expose these limitations.

---

# Consolidation

`consolidate` helps the agent reconcile old data with a newer domain model.

Example:

Historical data tracks:

```text
running
```

The newer domain tracks:

```text
exercise
```

The user may decide:

```text
running is_a exercise
```

while preserving `running` independently.

The agent reasons about this.

Atlas MCP only writes the resulting taxonomy, compatibility or data changes.

---

# Prefer semantic compatibility over destructive migration

If existing data can remain untouched while a new semantic relationship makes it useful, prefer that.

Example:

```text
running → exercise
```

may be represented in domain taxonomy rather than rewriting hundreds of historical records.

Do not rewrite history unnecessarily.

---

# `log/`

`log/` contains semantic action history.

Git history tells us which lines changed.

Atlas history should explain what was done.

Example:

```text
Daniel recorded a 5 km run.
```

or:

```text
Historical running records were made compatible with the new exercise category.
```

The agent should not directly edit `log/`.

Atlas MCP/runtime generates history as part of successful protected operations where appropriate.

---

# `bin/` in a child project

`bin/` contains generated artifacts.

Examples:

```text
bin/dashboard.html
bin/data/
bin/reports/
bin/indexes/
```

Generated output is not source-of-truth.

Rule:

> If it cannot be deleted and regenerated, it does not belong in bin/.

---

# Views

Views represent useful derived interpretations of domain data.

Examples:

```text
monthly-spend
current-balance
exercise-days
active-clients
```

Domain skills and documentation explain how to derive them.

The agent understands their meaning.

Deterministic calculations may be delegated to project scripts where appropriate.

---

# Dashboard

A project may generate:

```text
bin/dashboard.html
```

The dashboard is generated output.

The agent determines what information matters by using:

```text
domain skills
view definitions
data
compatibility information
```

Rendering should remain generic where practical.

Do not put domain interpretation inside Atlas MCP.

---

# Atlas repository `run/`

All common developer workflows belong under:

```text
run/
```

Initial expected scripts:

```text
run/
├── build
├── build-dbg
├── build-rel
├── test
├── install
└── clean
```

Keep scripts simple.

The primary developer interface should look like:

```bash
./run/build
./run/test
./run/install
```

Do not create a custom build framework.

---

# Debug build

Debug output:

```text
bin/dbg/
```

Example:

```text
bin/dbg/atlas-mcp
```

Run through:

```bash
./run/build-dbg
```

---

# Release build

Release output:

```text
bin/rel/
```

Example:

```text
bin/rel/atlas-mcp
```

Run through:

```bash
./run/build-rel
```

---

# Default build

`run/build` may build the debug version by default.

Do not make the default development workflow complicated.

---

# Installation

`run/install` installs/configures Atlas Automata for the current child repository.

It should:

* identify the child repository root
* verify `lib/atlas-automata`
* build or locate the required MCP binary
* install/configure required agent hooks
* install/configure Git hooks
* install required base skills where appropriate
* preserve existing project-owned skills
* configure MCP access

Installation must be explicit and inspectable.

---

# Git is part of the runtime

Git is not merely source control.

Atlas uses Git for:

* synchronization
* versioning
* concurrent work
* history
* replication
* rollback

The MCP exists partly to enforce a safe Git workflow around protected mutations.

---

# Mandatory sync before work

Every agent interaction with an Automatizer child project must begin from synchronized repository state.

The agent environment should be configured so that each new user request causes:

```text
SYNC BEFORE WORK
```

before normal reasoning proceeds.

This may be enforced using agent lifecycle hooks.

---

# Git merge policy

Atlas never rebases.

This is a hard rule.

Forbidden:

```bash
git rebase
git pull --rebase
git pull -r
```

When local and remote histories diverge:

```bash
git fetch origin
git merge origin/<branch>
```

Use merge.

---

# Published history is not rewritten

Atlas must not use:

```bash
git push --force
git push -f
git push --force-with-lease
```

during normal operation.

Published history is treated as immutable.

---

# MCP mutation workflow

For every protected mutation requested by the agent:

```text
RECEIVE MUTATION REQUEST
        ↓
VERIFY REPOSITORY
        ↓
VERIFY TARGET PATH
        ↓
FETCH ORIGIN
        ↓
COMPARE LOCAL / REMOTE
        ↓
MERGE REMOTE CHANGES IF NECESSARY
        ↓
IF CONFLICT → FAIL WITHOUT APPLYING MUTATION
        ↓
APPLY REQUESTED FILE CHANGE
        ↓
VALIDATE REPOSITORY RULES
        ↓
WRITE SEMANTIC HISTORY IF REQUIRED
        ↓
COMMIT
        ↓
FETCH ORIGIN AGAIN
        ↓
IF REMOTE CHANGED → MERGE
        ↓
REVALIDATE
        ↓
PUSH
        ↓
RETURN RESULT
```

Never rebase.

---

# Why sync occurs inside MCP

The agent may have reasoned over a repository that was current seconds earlier.

Another user may push while the agent is preparing a mutation.

Therefore each MCP mutation independently verifies remote state.

Do not trust only the sync performed at the beginning of the user turn.

---

# Conflict before mutation

If origin changed and merging produces a conflict before the requested mutation is applied:

```text
do not apply mutation
return conflict
```

The agent must inspect the newly conflicted/current state and decide how to reconcile it.

---

# Conflict after local mutation

After committing the requested mutation, Atlas fetches origin again.

If another user changed origin:

```text
merge
revalidate
push
```

If the merge conflicts:

```text
stop
report conflict
do not guess
```

---

# Push race

Origin may advance between final fetch and push.

If push is rejected because the remote advanced:

```text
fetch
merge
revalidate
push again
```

Never solve this using rebase.

Never force push.

---

# Git executable

Use the user's installed `git` executable initially.

Do not introduce libgit2.

This preserves:

* credential helpers
* SSH configuration
* Git configuration
* remote behavior
* existing user setup

Keep Git commands explicit and inspectable.

---

# Hooks

Hooks are guards.

They are not where Atlas business logic lives.

Possible agent hooks:

```text
SessionStart
UserPromptSubmit
PreToolUse
```

Primary purposes:

```text
sync before work
prevent direct writes to protected paths
prevent forbidden Git operations
```

Possible Git hooks:

```text
pre-rebase
pre-commit
pre-push
```

Primary purposes:

```text
block rebase
detect unauthorized protected changes
verify final Atlas workflow
```

Keep hooks very small.

The MCP contains the actual workflow.

---

# Direct protected writes

Agent hooks should block generic writes into:

```text
data/
doc/
ai/
log/
AUTOMATIZER.md
```

when technically possible.

Examples to reject:

```text
Edit data/...
Write doc/...
apply_patch ai/...
rm data/...
sed -i doc/...
```

The correct alternative is an Atlas MCP mutation.

---

# Source code

The child project's:

```text
src/
```

is ordinary project source code.

Direct agent modification is allowed unless the child project defines additional rules.

Atlas protected-state restrictions are primarily about:

```text
data
domain meaning
skills
history
configuration
```

not ordinary application code.

---

# Functional testing only

Do not build a unit-test suite.

Do not write tests around isolated functions merely to increase coverage.

Do not introduce interfaces to mock dependencies.

Do not use mocking frameworks.

Atlas 0.1 uses functional/integration testing.

---

# Test philosophy

Tests should exercise real workflows using:

```text
real temporary directories
real Git repositories
real Git remotes
real commits
real branches
real merges
real files
real MCP server calls
```

A few complete functional tests are preferable to hundreds of isolated tests.

---

# Functional test environment

Tests may create:

```text
temporary bare remote
clone A
clone B
Atlas MCP process
child project files
```

and then exercise realistic workflows end-to-end.

---

# Required functional test: basic mutation

Create a real child repository.

Use MCP to create a protected data file.

Verify:

```text
file exists
commit exists
remote contains commit
```

---

# Required functional test: sync before mutation

Create:

```text
remote
clone A
clone B
```

Push a change from clone A.

Request a mutation through Atlas MCP from clone B.

Verify clone B first incorporates the remote change.

Then verify the requested mutation is applied.

---

# Required functional test: divergence

Create commits independently in clone A and clone B.

Push A.

Request mutation/push from B.

Verify Atlas performs:

```text
fetch
merge
```

and produces a merge commit where appropriate.

Verify no rebase occurs.

---

# Required functional test: merge conflict

Create a genuine Git conflict between two clones.

Verify Atlas reports failure and does not silently resolve it.

---

# Required functional test: remote changes during operation

Synchronize clone B.

Create a local Atlas mutation.

Before final push, advance origin from clone A.

Verify Atlas detects the remote change and integrates using merge.

---

# Required functional test: rejected push race

Cause origin to advance immediately before Atlas pushes.

Verify:

```text
push rejected
fetch
merge
revalidate
push succeeds
```

No rebase.

---

# Required functional test: rebase rejection

Attempt:

```bash
git rebase
```

and:

```bash
git pull --rebase
```

through an agent-controlled workflow.

Verify Atlas guards reject them.

---

# Required functional test: force push rejection

Attempt:

```bash
git push --force
```

Verify rejection.

---

# Required functional test: direct protected write

Modify:

```text
data/test.md
```

without Atlas MCP.

Attempt to commit using the configured child-project environment.

Verify the repository guard detects that the protected change did not come through Atlas workflow.

---

# Required functional test: real domain evolution

Build a small realistic activity-tracking child project.

Initial historical data tracks:

```text
running
```

Later configure:

```text
exercise
```

with:

```text
running is_a exercise
```

Verify:

* historical running data remains intact
* domain documentation evolves
* compatibility information is preserved
* derived exercise information can include running
* running remains independently queryable

This test should involve actual files and actual agent-facing structure.

---

# Required functional test: generated dashboard

Use realistic child-project data.

Generate:

```text
bin/dashboard.html
```

Delete `bin/`.

Regenerate it.

Verify source data was unaffected.

---

# No coverage targets

Do not optimize for test coverage percentages.

Tests exist to verify behavior.

Focus on important end-to-end workflows.

---

# Atlas 0.1 priorities

Build Atlas incrementally.

Initial priority:

```text
1. repository detection
2. MCP startup
3. Git synchronization
4. merge-only policy
5. protected file create/update/delete/move
6. commit/push workflow
7. agent/Git guards
8. functional tests
9. base skill installation
10. first realistic child project
```

Do not attempt to implement the entire envisioned framework before these foundations work.

---

# Non-goals for 0.1

Do not build:

```text
SaaS backend
central server
database server
distributed locking
remote MCP
OAuth
plugin marketplace
ontology engine
generic workflow engine
event bus
dependency injection framework
custom persistence abstraction
custom Git abstraction
custom filesystem abstraction
MCP implementation
domain-specific GUI
```

Use existing tools and simple code.

---

# Design rule

When faced with two implementations:

```text
clever vs obvious
```

choose obvious.

When faced with:

```text
generic vs concrete
```

choose concrete.

When faced with:

```text
abstraction vs direct function
```

choose the direct function.

When faced with:

```text
framework vs 30 lines of clear code
```

write the 30 lines.

Atlas Automata 0.1 should remain small enough that a developer can open the repository, follow the execution path, and understand what happens without learning an internal architecture.

---

# Core invariants

The following rules are mandatory:

```text
1. THE AGENT UNDERSTANDS THE DOMAIN.

2. ATLAS MCP DOES NOT CONTAIN DOMAIN INTELLIGENCE.

3. THE AGENT MAY READ THE ENTIRE REPOSITORY.

4. PROTECTED STATE IS WRITTEN ONLY THROUGH ATLAS MCP.

5. data/ IS THE DOMAIN DATABASE.

6. doc/ IS AUTHORITATIVE DOMAIN MEANING.

7. ai/ CONTAINS INSTALLED AGENT CAPABILITIES.

8. SYNC BEFORE WORK.

9. EVERY MCP MUTATION SYNCS AGAIN BEFORE MODIFYING.

10. EVERY MCP MUTATION CHECKS ORIGIN AGAIN BEFORE PUSH.

11. MERGE ONLY.

12. NEVER REBASE.

13. NEVER FORCE PUSH.

14. DO NOT SILENTLY RESOLVE GIT CONFLICTS.

15. DO NOT SILENTLY DESTROY OR REINTERPRET HISTORICAL DATA.

16. NO INTERFACES.

17. NO ABSTRACTIONS.

18. USE SIMPLE FUNCTIONS AND COMPOSITION.

19. CONCURRENCY ONLY WHERE IT HAS A CONCRETE BENEFIT.

20. FUNCTIONAL TESTS OVER UNIT TESTS.

21. USE REAL INTEGRATIONS IN TESTS WHEN PRACTICAL.

22. KEEP ATLAS SMALL, EXPLICIT AND READABLE.
```
