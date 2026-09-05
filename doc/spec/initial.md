# Atlas Automata 0.1

## Purpose

Atlas Automata is a Git-native automation framework that lets an AI coding agent persist and evolve arbitrary personal or business domains without requiring a traditional application. The framework is domain-independent.

The core principle is:

> Agents read and reason. Atlas writes and synchronizes.

## Responsibility split

```text
AI AGENT
  reads repository and installed skills
  understands domain and user intent
  decides exact mutations
        |
        v
ATLAS MCP
  validates paths and repository state
  synchronizes Git using merge only
  applies the exact requested mutation
  records optional semantic history
  commits, resynchronizes, and pushes
        |
        v
GIT REPOSITORY
```

The agent owns semantic decisions: classifications, calculations, relationships, compatibility, dashboards, and how user intent maps to records. Atlas MCP has no LLM and contains no domain engine. It operates only on exact paths and content supplied by the agent.

## Child project

Atlas is installed as a pinned Git submodule at `lib/atlas-automata`:

```text
project/
├── AUTOMATIZER.md
├── AGENTS.md
├── src/
├── run/
├── bin/
├── data/
├── doc/
├── log/
├── ai/
│   ├── manifest.yaml
│   ├── skills/
│   ├── hooks/
│   └── instructions/
└── lib/atlas-automata/
```

The child repository pins an exact Atlas commit. Updating the framework is explicit. Installation is also explicit and inspectable; see [install.md](install.md).

## Example flow

For “I ran 5 km today,” the agent reads the activity skill, domain documentation, and existing records; chooses the exact record path and format; then calls `atlas_create`. Atlas synchronizes, writes only that content, validates, commits, checks the remote again, merges if needed, and pushes.

Atlas never decides whether running is exercise. That relationship belongs in project documentation or a skill and is evolved conservatively as described in [repository.md](repository.md).
