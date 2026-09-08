package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SyncState struct {
	Branch     string `json:"branch"`
	LocalHead  string `json:"local_head"`
	RemoteHead string `json:"remote_head"`
	Merged     bool   `json:"merged"`
}

func gitCommand(root string, extraEnv []string, args ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	command.Env = append(os.Environ(), extraEnv...)
	output, err := command.CombinedOutput()
	text := strings.TrimSpace(string(output))
	if err != nil {
		if text == "" {
			return "", fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
		}
		return text, fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), text, err)
	}
	return text, nil
}

func gitOutput(root string, args ...string) (string, error) {
	return gitCommand(root, nil, args...)
}

func currentBranch(root string) (string, error) {
	branch, err := gitOutput(root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("Atlas requires an attached branch: %w", err)
	}
	return strings.TrimSpace(branch), nil
}

func gitHead(root, revision string) (string, error) {
	head, err := gitOutput(root, "rev-parse", revision)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(head), nil
}

func gitFetch(root string) error {
	if _, err := gitOutput(root, "fetch", "--prune", "origin"); err != nil {
		return fmt.Errorf("fetch origin: %w", err)
	}
	return nil
}

func revisionExists(root, revision string) bool {
	_, err := gitOutput(root, "rev-parse", "--verify", "--quiet", revision)
	return err == nil
}

func isAncestor(root, older, newer string) (bool, error) {
	command := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", older, newer)
	err := command.Run()
	if err == nil {
		return true, nil
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("compare revisions %s and %s: %w", older, newer, err)
}

func syncRepo(root string) (SyncState, error) {
	state := SyncState{}
	if err := validateRepository(root); err != nil {
		return state, err
	}
	branch, err := currentBranch(root)
	if err != nil {
		return state, err
	}
	state.Branch = branch
	if err := gitFetch(root); err != nil {
		return state, err
	}

	remoteRevision := "refs/remotes/origin/" + branch
	if !revisionExists(root, remoteRevision) {
		state.LocalHead, err = gitHead(root, "HEAD")
		return state, err
	}
	state.RemoteHead, err = gitHead(root, remoteRevision)
	if err != nil {
		return state, err
	}
	localBefore, err := gitHead(root, "HEAD")
	if err != nil {
		return state, err
	}
	remoteIsAncestor, err := isAncestor(root, state.RemoteHead, localBefore)
	if err != nil {
		return state, err
	}
	if !remoteIsAncestor {
		if _, err := gitOutput(root, "merge", "--no-edit", remoteRevision); err != nil {
			return state, fmt.Errorf("merge origin/%s; resolve the conflict without rebasing: %w", branch, err)
		}
		state.Merged = true
	}
	state.LocalHead, err = gitHead(root, "HEAD")
	return state, err
}

func commitPaths(root, message string, paths []string) (string, error) {
	args := append([]string{"add", "--"}, paths...)
	if _, err := gitOutput(root, args...); err != nil {
		return "", fmt.Errorf("stage Atlas mutation: %w", err)
	}
	diffArgs := append([]string{"diff", "--cached", "--quiet", "--"}, paths...)
	command := exec.Command("git", append([]string{"-C", root}, diffArgs...)...)
	err := command.Run()
	if err == nil {
		return "", fmt.Errorf("mutation produced no changes")
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() != 1 {
		return "", fmt.Errorf("inspect staged mutation: %w", err)
	}
	commitArgs := append([]string{"commit", "-m", message, "--"}, paths...)
	if _, err := gitCommand(root, []string{"ATLAS_MCP_COMMIT=1"}, commitArgs...); err != nil {
		return "", fmt.Errorf("commit Atlas mutation: %w", err)
	}
	return gitHead(root, "HEAD")
}

func pushWithRetry(root, branch string) error {
	for attempt := 1; attempt <= 3; attempt++ {
		if _, err := gitOutput(root, "push", "origin", "HEAD:refs/heads/"+branch); err == nil {
			return nil
		} else if attempt == 3 {
			return fmt.Errorf("push rejected after %d merge-only attempts: %w", attempt, err)
		}
		if _, err := syncRepo(root); err != nil {
			return fmt.Errorf("recover from rejected push: %w", err)
		}
		if err := validateRepository(root); err != nil {
			return err
		}
	}
	return fmt.Errorf("push failed")
}
