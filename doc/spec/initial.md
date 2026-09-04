# Atlas Automata 0.1

## Purpose

Atlas Automata is a Git-native automation framework for AI coding agents.

The goal is to allow an agent such as Codex or Claude Code to automate arbitrary domains without requiring a traditional application.

Examples:

* personal finance
* company operations
* customer management
* calls and purchases
* personal assistant
* coaching
* activity tracking
* project management
* inventories
* research
* arbitrary personal or business workflows

The framework itself must remain domain-independent.

The AI agent understands the domain.

Atlas Automata provides the structure and workflow that allows the agent to safely persist and evolve that domain.

---

# Core architecture

The fundamental division of responsibility is:

```text
AI AGENT
    │
    │ reads the repository
    │ loads installed skills
    │ understands the domain
    │ reasons about user intent
    │ decides what needs to change
    │
    ▼
ATLAS MCP
    │
    │ exclusive controlled writer
    │ synchronizes Git
    │ validates repository rules
    │ performs requested file mutation
    │ records workflow state
    │ synchronizes Git again
    │
    ▼
GIT REPOSITORY
```

The agent is the intelligence.

The MCP is not.

---

# Fundamental principle

> Agents read and reason. Atlas writes and synchronizes.

The agent must be able to inspect the entire repository.

Protected project state must only be modified through Atlas MCP.

---

# Domain intelligence belongs to the agent

Atlas MCP does not contain an LLM.

Atlas MCP does not attempt to understand domain semantics.

It does not decide:

* what "exercise" means
* whether running counts as exercise
* what a financial transfer means
* what a customer status means
* which metric should appear in a dashboard
* how historical information should be interpreted

The host agent handles those decisions using:

```text
ai/skills/
doc/
data/
AUTOMATIZER.md
```

The agent may freely read those sources and reason over them.

Once the agent decides a concrete repository mutation is required, it requests that mutation through Atlas MCP.

Example:

```text
User:
"I ran 5 km today."

Agent:
- reads activity skill
- reads current domain structure
- determines required record format
- determines destination path
- produces exact content

Atlas MCP:
- synchronizes repository
- creates requested file
- validates workflow
- commits
- synchronizes again
- pushes
```

The MCP does not need to know what running means.

---

# Child project installation

Atlas Automata is installed in another repository as a Git submodule.

Example:

```bash
git init
git submodule add <atlas-automata-repository> lib/atlas-automata
```

A child project should conceptually look like:

```text
project/
├── AUTOMATIZER.md
├── AGENTS.md
│
├── src/
├── run/
├── bin/
│
├── data/
├── doc/
├── log/
│
├── ai/
│   ├── manifest.yaml
│   ├── skills/
│   ├── hooks/
│   └── instructions/
│
└── lib/
    └── atlas-automata/
```

The child project pins the exact Atlas Automata commit through the Git submodule.

Updating Atlas Automata must therefore be explicit.

---