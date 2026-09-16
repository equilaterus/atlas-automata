# Repository and domain model

## Framework repository

```text
atlas-automata/
├── AGENTS.md
├── README.md
├── VERSION
├── src/mcp/
├── run/
├── bin/dbg/
├── bin/rel/
├── ai/skills/
├── ai/hooks/
├── ai/instructions/
└── doc/spec/
```

Implementation source belongs under `src/`. Common workflows belong under `run/`. Do not add directory layers without a concrete current need.

## Child protected state

`data/` is the persistent domain database. Markdown is the primary 0.2 storage format, but physical layout is domain-specific: by entity, date, project, or another useful structure. The agent proposes a concrete directory tree from retrieval needs and obtains explicit user approval during configuration.

Every existing directory under `data/` contains an MCP-managed `index.md`. It lists immediate child folders and records as relative Markdown links, contains no independent domain facts, and is regenerated in the same commit as each data mutation. `index.md` is reserved; agents and users do not mutate it directly through Atlas.

`doc/` is authoritative domain meaning: concepts, schemas, taxonomy, relationships, compatibility, decisions, rules, and view definitions.

`ai/` contains capabilities installed in that child project. Skills explain record locations, required fields, meanings, relationships, ambiguity rules, and compatibility. Skills inform the agent; they never write state.

`AUTOMATIZER.md` describes the configured domain, users, collections, restrictions, organization, calculations, and project-specific workflow.

`log/` records semantic actions. Git says which lines changed; Atlas history says what the action meant. The MCP/runtime generates history as part of a successful protected mutation.

All five areas are directly readable and writable only through Atlas MCP.

## Mandatory configuration

The agent must complete the workflow in `ai/skills/configure/SKILL.md` before writing domain data. `AUTOMATIZER.md` begins with `atlas_setup: in_progress`; the agent marks it `complete` only after the user approves the model and the required domain, indexing, operations, and skill artifacts exist.

Atlas reports setup state and rejects mutations involving `data/` until setup is structurally complete. See [configuration.md](configuration.md) for the full user flow, artifacts, and indexing contract.

## Conservative domain evolution

Configuration must inspect existing configuration, documentation, skills, and data before proposing changes. Never treat an existing Automatizer as blank. Classify compatibility impact when useful as `SAFE`, `COMPATIBLE`, `PARTIAL`, `CONVERTIBLE`, `AMBIGUOUS`, or `BREAKING`.

Persist known limitations under a project-defined compatibility inventory, conventionally `doc/compatibility/inventory.yaml`. Ask the user when semantics are ambiguous.

Prefer semantic compatibility over rewriting history. If historical records use `running` and a newer model introduces `exercise`, represent `running is_a exercise` in taxonomy when that preserves both meanings. Historical running records remain intact and independently queryable.

## Views and generated output

Views are derived interpretations such as monthly spend, current balance, exercise days, or active clients. Their meaning belongs in project documentation and skills; deterministic calculation may live in project scripts.

Generated artifacts belong in child `bin/`, including `bin/dashboard.html`. If an artifact cannot be deleted and regenerated without changing source data, it does not belong in `bin/`. Atlas MCP contains no dashboard or domain interpretation logic.
