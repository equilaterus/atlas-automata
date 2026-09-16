package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type testRepository struct {
	remote string
	seed   string
	cloneA string
	cloneB string
}

func runTestGit(t *testing.T, directory string, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), directory, err, output)
	}
	return strings.TrimSpace(string(output))
}

func expectTestGitFailure(t *testing.T, directory string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", directory}, args...)...)
	if output, err := command.CombinedOutput(); err == nil {
		t.Fatalf("git %s unexpectedly succeeded in %s\n%s", strings.Join(args, " "), directory, output)
	}
}

func writeTestFile(t *testing.T, root, path, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func commitTestFile(t *testing.T, root, path, content, message string) {
	t.Helper()
	writeTestFile(t, root, path, content)
	runTestGit(t, root, "add", "--", path)
	runTestGit(t, root, "commit", "-m", message)
}

func makeTestRepository(t *testing.T) testRepository {
	return makeTestRepositoryWithSetup(t, true)
}

func makeTestRepositoryWithSetup(t *testing.T, configured bool) testRepository {
	t.Helper()
	base := t.TempDir()
	repository := testRepository{
		remote: filepath.Join(base, "remote.git"),
		seed:   filepath.Join(base, "seed"),
		cloneA: filepath.Join(base, "a"),
		cloneB: filepath.Join(base, "b"),
	}
	runTestGit(t, base, "init", "--bare", "--initial-branch=dev", repository.remote)
	runTestGit(t, base, "clone", repository.remote, repository.seed)
	configureTestGit(t, repository.seed)
	writeTestFile(t, repository.seed, "README.md", "child repository\n")
	if configured {
		writeTestFile(t, repository.seed, "AUTOMATIZER.md", "---\natlas_setup: complete\natlas_setup_version: 1\n---\n# Test domain\n")
		writeTestFile(t, repository.seed, "doc/domain/model.md", "# Model\n")
		writeTestFile(t, repository.seed, "doc/domain/indexing.md", "# Indexing\n")
		writeTestFile(t, repository.seed, "doc/domain/operations.md", "# Operations\n")
		writeTestFile(t, repository.seed, "ai/skills/test-domain/SKILL.md", "---\nname: test-domain\ndescription: Test domain.\n---\n")
	}
	runTestGit(t, repository.seed, "add", ".")
	runTestGit(t, repository.seed, "commit", "-m", "initial")
	runTestGit(t, repository.seed, "push", "-u", "origin", "dev")
	runTestGit(t, base, "clone", repository.remote, repository.cloneA)
	runTestGit(t, base, "clone", repository.remote, repository.cloneB)
	configureTestGit(t, repository.cloneA)
	configureTestGit(t, repository.cloneB)
	return repository
}

func TestSetupGateRequiresCompleteGuidedConfiguration(t *testing.T) {
	repository := makeTestRepositoryWithSetup(t, false)
	if state, err := setupState(repository.cloneB); err != nil || state != "not_started" {
		t.Fatalf("initial setup state = %q, %v; want not_started", state, err)
	}
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/item.md", "content": "blocked\n",
	}, true)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "AUTOMATIZER.md", "content": "---\natlas_setup: complete\natlas_setup_version: 1\n---\n",
	}, true)
	if _, err := os.Stat(filepath.Join(repository.cloneB, "AUTOMATIZER.md")); !os.IsNotExist(err) {
		t.Fatalf("premature completed setup was written: %v", err)
	}

	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "AUTOMATIZER.md", "content": "---\natlas_setup: in_progress\natlas_setup_version: 1\n---\n# Test domain\n",
	}, false)
	for path, content := range map[string]string{
		"doc/domain/model.md":            "# Model\n",
		"doc/domain/indexing.md":         "# Indexing\n",
		"doc/domain/operations.md":       "# Operations\n",
		"ai/skills/test-domain/SKILL.md": "---\nname: test-domain\ndescription: Test domain.\n---\n",
	} {
		callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
			"path": path, "content": content,
		}, false)
	}
	callTestTool(t, repository.cloneB, "atlas_update", map[string]any{
		"path": "AUTOMATIZER.md", "content": "---\natlas_setup: complete\natlas_setup_version: 1\n---\n# Test domain\n",
	}, false)
	if state, err := setupState(repository.cloneB); err != nil || state != "complete" {
		t.Fatalf("final setup state = %q, %v; want complete", state, err)
	}
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/item.md", "content": "allowed\n",
	}, false)
}

func TestDataMutationsMaintainFolderIndexes(t *testing.T) {
	repository := makeTestRepository(t)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/games/positive/alpha.md", "content": "title: Alpha\n",
	}, false)

	for path, expected := range map[string][]string{
		"data/index.md":                {dataIndexMarker, "[games/](games/index.md)"},
		"data/games/index.md":          {dataIndexMarker, "[positive/](positive/index.md)"},
		"data/games/positive/index.md": {dataIndexMarker, "[alpha](alpha.md)"},
	} {
		content, err := os.ReadFile(filepath.Join(repository.cloneB, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("read generated index %s: %v", path, err)
		}
		for _, fragment := range expected {
			if !strings.Contains(string(content), fragment) {
				t.Fatalf("generated index %s does not contain %q: %s", path, fragment, content)
			}
		}
	}

	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/games/negative/beta.md", "content": "title: Beta\n",
	}, false)
	callTestTool(t, repository.cloneB, "atlas_move", map[string]any{
		"source": "data/games/positive/alpha.md", "destination": "data/games/negative/alpha.md",
	}, false)

	positive, err := os.ReadFile(filepath.Join(repository.cloneB, "data/games/positive/index.md"))
	if err != nil || strings.Contains(string(positive), "alpha.md") || !strings.Contains(string(positive), "_Empty._") {
		t.Fatalf("source index was not refreshed after move: %v %s", err, positive)
	}
	negative, err := os.ReadFile(filepath.Join(repository.cloneB, "data/games/negative/index.md"))
	if err != nil || !strings.Contains(string(negative), "alpha.md") || !strings.Contains(string(negative), "beta.md") {
		t.Fatalf("destination index was not refreshed after move: %v %s", err, negative)
	}

	callTestTool(t, repository.cloneB, "atlas_delete", map[string]any{
		"path": "data/games/negative/beta.md",
	}, false)
	negative, err = os.ReadFile(filepath.Join(repository.cloneB, "data/games/negative/index.md"))
	if err != nil || !strings.Contains(string(negative), "alpha.md") || strings.Contains(string(negative), "beta.md") {
		t.Fatalf("index was not refreshed after delete: %v %s", err, negative)
	}

	callTestTool(t, repository.cloneB, "atlas_update", map[string]any{
		"path": "data/games/index.md", "content": "manual\n",
	}, true)
	if got := runTestGit(t, repository.cloneB, "status", "--porcelain"); got != "" {
		t.Fatalf("repository is dirty after rejected index mutation: %s", got)
	}
}

func TestMigrationStateOnlyAllowsDataMoves(t *testing.T) {
	repository := makeTestRepository(t)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/games/item.md", "content": "title: Item\n",
	}, false)
	callTestTool(t, repository.cloneB, "atlas_update", map[string]any{
		"path": "AUTOMATIZER.md", "content": "---\natlas_setup: migration\natlas_setup_version: 1\n---\n# Test domain\n",
	}, false)
	if state, err := setupState(repository.cloneB); err != nil || state != "migration" {
		t.Fatalf("migration setup state = %q, %v; want migration", state, err)
	}

	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/games/new.md", "content": "blocked\n",
	}, true)
	callTestTool(t, repository.cloneB, "atlas_update", map[string]any{
		"path": "data/games/item.md", "content": "blocked\n",
	}, true)
	callTestTool(t, repository.cloneB, "atlas_delete", map[string]any{
		"path": "data/games/item.md",
	}, true)
	callTestTool(t, repository.cloneB, "atlas_move", map[string]any{
		"source": "data/games/item.md", "destination": "data/games/positive/item.md",
	}, false)
	if _, err := os.Stat(filepath.Join(repository.cloneB, "data/games/positive/item.md")); err != nil {
		t.Fatalf("approved migration move missing: %v", err)
	}

	callTestTool(t, repository.cloneB, "atlas_update", map[string]any{
		"path": "AUTOMATIZER.md", "content": "---\natlas_setup: complete\natlas_setup_version: 1\n---\n# Test domain\n",
	}, false)
	if state, err := setupState(repository.cloneB); err != nil || state != "complete" {
		t.Fatalf("final setup state = %q, %v; want complete", state, err)
	}
}

func configureTestGit(t *testing.T, root string) {
	t.Helper()
	runTestGit(t, root, "config", "user.name", "Atlas Test")
	runTestGit(t, root, "config", "user.email", "atlas@example.test")
}

func callTestTool(t *testing.T, root, name string, arguments map[string]any, wantError bool) *mcp.CallToolResult {
	t.Helper()
	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := newMCPServer(root).Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "atlas-functional-test", Version: "0.2.0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()
	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
	failed := err != nil || (result != nil && result.IsError)
	if failed != wantError {
		var content []string
		if result != nil {
			for _, item := range result.Content {
				if text, ok := item.(*mcp.TextContent); ok {
					content = append(content, text.Text)
				}
			}
		}
		t.Fatalf("tool %s failure=%v, want %v; error=%v content=%q", name, failed, wantError, err, content)
	}
	return result
}

func TestMCPMutationLifecycleAndSync(t *testing.T) {
	repository := makeTestRepository(t)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "../escape.md", "content": "forbidden\n",
	}, true)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "src/not-protected.md", "content": "forbidden\n",
	}, true)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/activities/run.md", "content": "distance_km: 5\n", "summary": "Recorded a 5 km run.",
	}, false)
	runTestGit(t, repository.cloneA, "pull", "--ff-only")
	if content, err := os.ReadFile(filepath.Join(repository.cloneA, "data/activities/run.md")); err != nil || string(content) != "distance_km: 5\n" {
		t.Fatalf("remote mutation missing: %q, %v", content, err)
	}

	callTestTool(t, repository.cloneB, "atlas_update", map[string]any{
		"path": "data/activities/run.md", "content": "distance_km: 6\n",
	}, false)
	callTestTool(t, repository.cloneB, "atlas_move", map[string]any{
		"source": "data/activities/run.md", "destination": "data/activities/2026-09-04.md",
	}, false)
	callTestTool(t, repository.cloneB, "atlas_delete", map[string]any{
		"path": "data/activities/2026-09-04.md",
	}, false)

	runTestGit(t, repository.cloneA, "pull", "--ff-only")
	commitTestFile(t, repository.cloneA, "remote-change.txt", "remote\n", "remote change")
	runTestGit(t, repository.cloneA, "push", "origin", "dev")
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "doc/domain/state.md", "content": "# State\n",
	}, false)
	if _, err := os.Stat(filepath.Join(repository.cloneB, "remote-change.txt")); err != nil {
		t.Fatalf("mutation did not synchronize first: %v", err)
	}
	if got := runTestGit(t, repository.cloneB, "status", "--porcelain"); got != "" {
		t.Fatalf("repository is dirty after lifecycle: %s", got)
	}
}

func TestDivergenceUsesMergeAndConflictStopsMutation(t *testing.T) {
	repository := makeTestRepository(t)
	commitTestFile(t, repository.cloneA, "from-a.txt", "a\n", "change from A")
	runTestGit(t, repository.cloneA, "push", "origin", "dev")
	commitTestFile(t, repository.cloneB, "from-b.txt", "b\n", "change from B")
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/item.md", "content": "value\n",
	}, false)
	parents := strings.Fields(runTestGit(t, repository.cloneB, "show", "--format=%P", "--no-patch", "HEAD~1"))
	if len(parents) != 2 {
		t.Fatalf("expected a merge commit before the mutation, got %d parents", len(parents))
	}

	runTestGit(t, repository.cloneA, "pull", "--ff-only")
	commitTestFile(t, repository.cloneA, "doc/conflict.md", "remote\n", "remote conflict base")
	runTestGit(t, repository.cloneA, "push", "origin", "dev")
	runTestGit(t, repository.cloneB, "pull", "--ff-only")
	commitTestFile(t, repository.cloneA, "doc/conflict.md", "remote changed\n", "remote side")
	runTestGit(t, repository.cloneA, "push", "origin", "dev")
	commitTestFile(t, repository.cloneB, "doc/conflict.md", "local changed\n", "local side")
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/must-not-exist.md", "content": "no\n",
	}, true)
	if _, err := os.Stat(filepath.Join(repository.cloneB, "data/must-not-exist.md")); !os.IsNotExist(err) {
		t.Fatalf("mutation was applied after a merge conflict: %v", err)
	}
}

func TestRejectedPushRaceRetriesWithMerge(t *testing.T) {
	repository := makeTestRepository(t)
	hook := fmt.Sprintf(`#!/usr/bin/env bash
set -e
marker=%q
if [[ ! -e "$marker" ]]; then
  touch "$marker"
  unset $(git rev-parse --local-env-vars)
  printf 'race\n' > %q
  git -C %q add race.txt
  git -C %q commit -m 'advance during push'
  git -C %q push origin dev
fi
`, filepath.Join(repository.cloneB, ".race-fired"), filepath.Join(repository.cloneA, "race.txt"), repository.cloneA, repository.cloneA, repository.cloneA)
	writeTestFile(t, repository.cloneB, ".git/hooks/pre-push", hook)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/raced.md", "content": "survived\n",
	}, false)
	if _, err := os.Stat(filepath.Join(repository.cloneB, "race.txt")); err != nil {
		t.Fatalf("push-race change was not merged: %v", err)
	}
	if parents := strings.Fields(runTestGit(t, repository.cloneB, "show", "--format=%P", "--no-patch", "HEAD")); len(parents) != 2 {
		t.Fatalf("expected retry merge commit, got %d parents", len(parents))
	}
}

func TestRemoteChangeDuringMutationIsMergedBeforePush(t *testing.T) {
	repository := makeTestRepository(t)
	hook := fmt.Sprintf(`#!/usr/bin/env bash
set -e
marker=%q
if [[ ! -e "$marker" ]]; then
  touch "$marker"
  unset $(git rev-parse --local-env-vars)
  printf 'during commit\n' > %q
  git -C %q add during-commit.txt
  git -C %q commit -m 'advance during mutation'
  git -C %q push origin dev
fi
`, filepath.Join(repository.cloneB, ".commit-race-fired"), filepath.Join(repository.cloneA, "during-commit.txt"), repository.cloneA, repository.cloneA, repository.cloneA)
	writeTestFile(t, repository.cloneB, ".git/hooks/commit-msg", hook)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/during-race.md", "content": "survived\n",
	}, false)
	if _, err := os.Stat(filepath.Join(repository.cloneB, "during-commit.txt")); err != nil {
		t.Fatalf("change made during mutation was not merged: %v", err)
	}
	if parents := strings.Fields(runTestGit(t, repository.cloneB, "show", "--format=%P", "--no-patch", "HEAD")); len(parents) != 2 {
		t.Fatalf("expected final synchronization merge, got %d parents", len(parents))
	}
}

func TestRepositoryGuards(t *testing.T) {
	repository := makeTestRepository(t)
	frameworkRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	hooks := filepath.Join(repository.cloneB, ".githooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"pre-commit", "pre-rebase", "pre-push"} {
		content, err := os.ReadFile(filepath.Join(frameworkRoot, "ai", "hooks", name))
		if err != nil {
			t.Fatal(err)
		}
		writeTestFile(t, repository.cloneB, ".githooks/"+name, string(content))
	}
	runTestGit(t, repository.cloneB, "config", "core.hooksPath", ".githooks")
	writeTestFile(t, repository.cloneB, "data/direct.md", "forbidden\n")
	runTestGit(t, repository.cloneB, "add", "data/direct.md")
	command := exec.Command("git", "-C", repository.cloneB, "commit", "-m", "direct")
	if err := command.Run(); err == nil {
		t.Fatal("direct protected commit unexpectedly succeeded")
	}
	runTestGit(t, repository.cloneB, "reset", "HEAD", "--", "data/direct.md")
	if err := os.Remove(filepath.Join(repository.cloneB, "data/direct.md")); err != nil {
		t.Fatal(err)
	}
	command = exec.Command(filepath.Join(hooks, "pre-rebase"))
	command.Dir = repository.cloneB
	if err := command.Run(); err == nil {
		t.Fatal("rebase guard unexpectedly succeeded")
	}

	commitTestFile(t, repository.cloneA, "remote-only.txt", "remote\n", "remote only")
	runTestGit(t, repository.cloneA, "push", "origin", "dev")
	commitTestFile(t, repository.cloneB, "local-only.txt", "local\n", "local only")
	runTestGit(t, repository.cloneB, "fetch", "origin")
	expectTestGitFailure(t, repository.cloneB, "rebase", "origin/dev")
	expectTestGitFailure(t, repository.cloneB, "pull", "--rebase", "origin", "dev")
	expectTestGitFailure(t, repository.cloneB, "push", "--force", "origin", "dev")
}

func TestDomainEvolutionAndRegenerableDashboard(t *testing.T) {
	repository := makeTestRepository(t)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "data/activities/2026-08-01.md", "content": "type: running\ndistance_km: 5\n",
	}, false)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "doc/domain/activities.md", "content": "# Activities\n\n`running is_a exercise`; running remains queryable as its own type.\n",
	}, false)
	callTestTool(t, repository.cloneB, "atlas_create", map[string]any{
		"path": "doc/compatibility/inventory.yaml", "content": "- scope: historical-running\n  status: COMPATIBLE\n  rule: running is_a exercise\n",
	}, false)
	original, err := os.ReadFile(filepath.Join(repository.cloneB, "data/activities/2026-08-01.md"))
	if err != nil || !strings.Contains(string(original), "type: running") {
		t.Fatalf("historical running data changed: %v %q", err, original)
	}

	build := `#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
mkdir -p "$root/bin"
count="$(grep -R -l '^type: running$' "$root/data/activities" | wc -l)"
printf '<html><body>exercise-days: %s; running-days: %s</body></html>\n' "$count" "$count" > "$root/bin/dashboard.html"
`
	writeTestFile(t, repository.cloneB, "src/build-dashboard", build)
	command := exec.Command(filepath.Join(repository.cloneB, "src/build-dashboard"))
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build dashboard: %v %s", err, output)
	}
	if err := os.RemoveAll(filepath.Join(repository.cloneB, "bin")); err != nil {
		t.Fatal(err)
	}
	command = exec.Command(filepath.Join(repository.cloneB, "src/build-dashboard"))
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("regenerate dashboard: %v %s", err, output)
	}
	dashboard, err := os.ReadFile(filepath.Join(repository.cloneB, "bin/dashboard.html"))
	if err != nil || !strings.Contains(string(dashboard), "exercise-days: 1; running-days: 1") {
		t.Fatalf("unexpected dashboard: %v %q", err, dashboard)
	}
	after, _ := os.ReadFile(filepath.Join(repository.cloneB, "data/activities/2026-08-01.md"))
	if string(after) != string(original) {
		t.Fatal("dashboard regeneration changed source data")
	}
}
