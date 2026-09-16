---
name: configure
description: Guide the mandatory first-time or evolutionary configuration of an Atlas child repository. Use when atlas_status reports setup other than complete, when AUTOMATIZER.md is missing or incomplete, before the first domain-data mutation, or when the user changes collections, schemas, identity, indexing, taxonomy, views, calculations, privacy, or workflows.
---

# Configure Atlas

Lead the user through configuration. Do not ask the user to design files, paths, schemas, or indexes unaided. Translate their goals into concrete proposals, explain consequences in plain language, and obtain confirmation at each decision point.

Never write domain records under `data/` until every phase below is complete and `atlas_status` reports `setup: complete`.

## Operating rules

- Run `atlas_sync`, then `atlas_status` before reasoning.
- Inspect `AUTOMATIZER.md`, `doc/`, `ai/skills/`, and existing `data/`. Never assume the repository is blank.
- Reuse answers already stated by the user, but summarize them for confirmation. Do not make the user repeat themselves.
- Conduct the setup as a short conversation. Ask one focused group of related questions at a time.
- Complete every phase. If a concern does not apply, record it explicitly as `Not applicable` with the reason.
- Do not treat implementation defaults as user decisions. Present a recommended choice and its tradeoff.
- Classify changes to existing domains as `SAFE`, `COMPATIBLE`, `PARTIAL`, `CONVERTIBLE`, `AMBIGUOUS`, or `BREAKING`.
- Use only Atlas MCP mutations for protected files. Mark setup complete last.

## Mandatory phases

### 1. Purpose and boundaries

Establish and confirm:

- what the repository should remember or automate;
- who uses it;
- what belongs inside and outside its scope;
- privacy, retention, or sharing restrictions;
- the user's language and terminology.

### 2. Collections and lifecycles

Identify each collection or entity. For each one, establish:

- its meaning and examples;
- whether records represent an entity, event, relationship, or snapshot;
- creation, update, archival, and deletion behavior;
- statuses and meaningful transitions;
- relationships with other collections.

### 3. Retrieval and indexing needs

Ask how the user expects to find and use the information before proposing storage. Cover:

- lookup keys and duplicate detection;
- filters, sorting, grouping, and full-text needs;
- cross-collection navigation;
- common questions and recurring reports;
- expected scale and whether generated indexes are worthwhile.

Then propose and confirm:

- stable identity rules;
- canonical file paths and slug rules;
- physical organization by entity, date, project, or a justified combination;
- the exact directory partition under `data/`, shown as a tree before approval;
- the MCP-managed `index.md` that Atlas creates in every directory under `data/`;
- any additional authoritative indexes, if the domain genuinely needs them;
- reproducible generated indexes under `bin/indexes/` when source records alone are authoritative.

Do not silently choose a flat layout, expected scale, or "no index" policy. The user must see and approve the partition. Folder `index.md` files are mandatory navigational manifests maintained automatically by Atlas MCP; neither the user nor the agent writes them directly. Never maintain the same domain fact manually in both a record and an index. Document how every additional generated index is rebuilt.

### 4. Schema, taxonomy, and validation

Define and confirm:

- required and optional fields;
- types, units, formats, controlled vocabularies, and defaults;
- taxonomy and relationship semantics;
- ambiguity and missing-data rules;
- validation and duplicate-resolution behavior;
- source, confidence, or provenance fields when relevant.

### 5. Views, calculations, and workflows

Define and confirm:

- derived views and the questions they answer;
- formulas, time boundaries, currencies, units, and rounding;
- imports, exports, dashboards, reminders, or recurring workflows;
- which artifacts are authoritative and which are reproducible output under `bin/`.

### 6. Compatibility review

When configuration or data already exists:

- map the proposal to current records and terminology;
- classify each impact;
- prefer semantic relationships over rewriting history;
- record unresolved limitations in `doc/compatibility/inventory.yaml`;
- stop for explicit user direction on `AMBIGUOUS` or `BREAKING` decisions.

For a new repository, record that compatibility is not applicable rather than omitting the review.

### 7. Final proposal and approval

Present a compact configuration summary covering every prior phase. Show the proposed collections, identities, layout, indexes, core schema, relationships, views, calculations, restrictions, and compatibility result.

Ask for explicit approval. Do not persist the configuration before approval.

### 8. Persist the configuration

Create or update these standard artifacts sequentially:

1. `AUTOMATIZER.md` with YAML front matter containing `atlas_setup: in_progress` and `atlas_setup_version: 1`, plus purpose, users, scope, collections, restrictions, and the phase checklist.
2. `doc/domain/model.md` with collections, schemas, lifecycles, taxonomy, relationships, validation, and ambiguity rules.
3. `doc/domain/indexing.md` with access patterns, identities, the approved `data/` directory tree, duplicate detection, mandatory folder indexes, any additional indexes, and rebuild rules.
4. `doc/domain/operations.md` with views, calculations, imports, exports, recurring workflows, and authoritative/generated boundaries.
5. At least one domain skill under `ai/skills/<domain-name>/SKILL.md` that tells future agents how to read, validate, query, and mutate the configured domain.
6. `doc/compatibility/inventory.yaml` when compatibility issues exist.
7. Any deterministic source or build scripts outside protected state that the approved design actually requires.
8. For an approved evolution that requires moving existing data, change only the setup state to `atlas_setup: migration`, move records exclusively with `atlas_move`, and verify the migrated tree and folder indexes. Skip this state when no data paths change.
9. Update `AUTOMATIZER.md` last, preserving the approved content and changing only the setup state to `atlas_setup: complete`.

Include a concise semantic summary with every protected mutation. Never mark setup complete if a standard artifact is missing or a decision remains unresolved.

### 9. Verify and hand off

Run `atlas_status`. Confirm all of the following before inviting data entry:

- `setup` is `complete`;
- both Git and Atlas report clean state;
- the domain skill is installed;
- the approved partition and folder-index behavior are documented;
- no unresolved decision was silently defaulted.

Tell the user what was configured in plain language, then offer the first domain action. If verification fails, remain in setup and repair the incomplete phase.

## Evolving an existing setup

Run the same nine phases for changes, focusing questions on affected decisions while explicitly confirming unaffected phases. Set `atlas_setup: in_progress` before applying an approved multi-file evolution. If existing records must move, use `atlas_setup: migration` only after the new documentation and skills agree; this state permits only `atlas_move` between data paths. Restore `complete` only after the migrated tree and folder indexes are verified.
