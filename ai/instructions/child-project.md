<!-- atlas-automata:start -->
## Atlas Automata

Start every user request by calling `atlas_sync` and `atlas_status`.

If `atlas_status` reports `setup` other than `complete`, load `ai/skills/configure/SKILL.md` and conduct its complete agent-led configuration before accepting or writing domain data. Do not skip a phase, silently choose a semantic default, or ask the user to design the repository unaided.

Read `AUTOMATIZER.md`, relevant documentation, installed skills, and existing data before reasoning about a mutation. `AUTOMATIZER.md` and files under `data/`, `doc/`, `ai/`, and `log/` are protected: read them directly, but write them only with `atlas_create`, `atlas_update`, `atlas_delete`, or `atlas_move`.

Use merge-only Git synchronization. Never rebase, force-push, silently resolve a conflict, or silently reinterpret historical information.
<!-- atlas-automata:end -->
