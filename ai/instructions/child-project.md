<!-- atlas-automata:start -->
## Atlas Automata

At the start of an agent session, call `atlas_sync` once and then call `atlas_status`. Do not repeat `atlas_sync` for every user request in the same session; protected mutation tools synchronize internally.

Every user request concerning the configured domain is domain work, even when phrased as casual conversation. Before answering, load the relevant discovered skill and existing records. Treat new facts, preferences, ratings, corrections, and decisions supplied by the user as candidate domain mutations; preserve them through Atlas MCP when they belong in the configured domain.

If `atlas_status` reports `setup` other than `complete`, load `ai/skills/configure/SKILL.md` and conduct its complete agent-led configuration before accepting or writing domain data. Do not skip a phase, silently choose a semantic default, or ask the user to design the repository unaided.

Read `AUTOMATIZER.md`, relevant documentation, installed skills, and existing data before reasoning about a mutation. `AUTOMATIZER.md` and files under `data/`, `doc/`, `ai/`, and `log/` are protected: read them directly, but write them only with `atlas_create`, `atlas_update`, `atlas_delete`, or `atlas_move`.

Use merge-only Git synchronization. Never rebase, force-push, silently resolve a conflict, or silently reinterpret historical information.
<!-- atlas-automata:end -->
