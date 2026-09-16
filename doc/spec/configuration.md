# Agent-led configuration contract

Atlas installation is incomplete until the child domain has passed the guided configuration. The end user describes goals and approves semantic decisions; the agent owns repository design and conducts every phase in `ai/skills/configure/SKILL.md`.

## User experience

The agent asks focused questions in ordinary language, reuses information the user already provided, recommends concrete choices, and explains tradeoffs. It never asks the user to invent directory layouts, schemas, indexes, or taxonomies without guidance.

All nine phases are mandatory: purpose and boundaries; collections and lifecycles; retrieval and indexing; schema and taxonomy; views and workflows; compatibility; final approval; persistence; and verification. A phase may be recorded as not applicable with a reason, but may not disappear from the review.

The agent presents one complete proposal before protected state is written. `AMBIGUOUS` and `BREAKING` choices require explicit user direction.

## Setup states

`atlas_status` reports one of:

| State | Meaning |
| --- | --- |
| `not_started` | `AUTOMATIZER.md` does not exist. |
| `in_progress` | The file exists but its front matter does not contain `atlas_setup: complete`. |
| `incomplete` | The completion marker exists but one or more structural artifacts are missing. |
| `complete` | The marker, standard domain documents, and at least one non-base domain skill exist. |

Atlas rejects every mutation involving `data/` unless setup is `complete`. Configuration mutations under `AUTOMATIZER.md`, `doc/`, and `ai/` remain available so the agent can finish or repair setup.

## Required artifacts

The guided setup persists:

- `AUTOMATIZER.md`, including `atlas_setup` and `atlas_setup_version` in YAML front matter;
- `doc/domain/model.md` for collections, schemas, lifecycles, taxonomy, relationships, and validation;
- `doc/domain/indexing.md` for access patterns, identity, layout, indexes, duplicate detection, and rebuild rules;
- `doc/domain/operations.md` for views, calculations, integrations, workflows, and authoritative/generated boundaries;
- at least one domain skill other than the base `configure` and `mutate` skills;
- `doc/compatibility/inventory.yaml` when known compatibility issues exist.

The agent first writes `atlas_setup: in_progress`. It changes that marker to `complete` only after the approved supporting artifacts exist. The MCP rejects a premature completion marker.

## Indexing rule

Design indexes from user retrieval needs, not from storage convenience. Establish lookup keys, filters, sorting, grouping, cross-collection navigation, duplicate detection, and expected scale before selecting physical layout.

The final proposal must show the exact directory tree under `data/`; a flat layout, expected scale, and optional indexes are decisions, not defaults. The user explicitly approves that partition before persistence.

Atlas MCP reserves and maintains `index.md` in every existing directory under `data/`. These files are deterministic navigational manifests listing immediate child folders and records. A data mutation rebuilds affected manifests in the same commit, and agents never edit them directly.

Additional authoritative indexes may live under `data/` when they contain unique domain state. Other reproducible indexes may live under `bin/indexes/` only when the approved domain needs them and must document their rebuild process. Never require manual synchronization of the same domain fact in both a record and an index.

## Evolution

Reconfiguration follows the same phases. The agent inspects existing state, classifies compatibility, obtains approval, marks setup `in_progress`, applies sequential Atlas mutations, and restores `complete` only after documentation, domain skills, and indexing rules agree. Historical information is never silently reinterpreted.
