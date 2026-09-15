package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MutationResult struct {
	Operation string `json:"operation"`
	Path      string `json:"path"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
}

func createFile(root, path, content string) (string, error) {
	relative, full, err := protectedPath(root, path)
	if err != nil {
		return "", err
	}
	if _, err := os.Lstat(full); err == nil {
		return "", fmt.Errorf("path already exists: %s", relative)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", fmt.Errorf("create parent directory: %w", err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("create %s: %w", relative, err)
	}
	return relative, nil
}

func updateFile(root, path, content string) (string, error) {
	relative, full, err := protectedPath(root, path)
	if err != nil {
		return "", err
	}
	info, err := os.Lstat(full)
	if err != nil {
		return "", fmt.Errorf("update %s: %w", relative, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("update target is not a regular file: %s", relative)
	}
	temporary, err := os.CreateTemp(filepath.Dir(full), ".atlas-update-*")
	if err != nil {
		return "", err
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if _, err := temporary.WriteString(content); err != nil {
		temporary.Close()
		return "", err
	}
	if err := temporary.Chmod(info.Mode().Perm()); err != nil {
		temporary.Close()
		return "", err
	}
	if err := temporary.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(temporaryName, full); err != nil {
		return "", fmt.Errorf("replace %s: %w", relative, err)
	}
	return relative, nil
}

func deleteFile(root, path string) (string, error) {
	relative, full, err := protectedPath(root, path)
	if err != nil {
		return "", err
	}
	info, err := os.Lstat(full)
	if err != nil {
		return "", fmt.Errorf("delete %s: %w", relative, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("delete target is not a regular file: %s", relative)
	}
	if err := os.Remove(full); err != nil {
		return "", fmt.Errorf("delete %s: %w", relative, err)
	}
	return relative, nil
}

func moveFile(root, source, destination string) ([]string, error) {
	sourceRelative, sourceFull, err := protectedPath(root, source)
	if err != nil {
		return nil, err
	}
	destinationRelative, destinationFull, err := protectedPath(root, destination)
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(sourceFull)
	if err != nil {
		return nil, fmt.Errorf("move %s: %w", sourceRelative, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("move source is not a regular file: %s", sourceRelative)
	}
	if _, err := os.Lstat(destinationFull); err == nil {
		return nil, fmt.Errorf("move destination already exists: %s", destinationRelative)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(destinationFull), 0o755); err != nil {
		return nil, err
	}
	if err := os.Rename(sourceFull, destinationFull); err != nil {
		return nil, fmt.Errorf("move %s to %s: %w", sourceRelative, destinationRelative, err)
	}
	return []string{sourceRelative, destinationRelative}, nil
}

func runMutation(root, operation, path, destination, content, summary, commitMessage string) (MutationResult, error) {
	result := MutationResult{Operation: operation, Path: path}
	if _, _, err := protectedPath(root, path); err != nil {
		return result, err
	}
	if destination != "" {
		if _, _, err := protectedPath(root, destination); err != nil {
			return result, err
		}
	}
	if operation == "move" && destination == "AUTOMATIZER.md" {
		return result, fmt.Errorf("AUTOMATIZER.md must be created or updated explicitly so Atlas can validate setup completion")
	}
	if err := ensureNoProtectedChanges(root); err != nil {
		return result, err
	}
	state, err := syncRepo(root)
	if err != nil {
		return result, err
	}
	if err := ensureNoProtectedChanges(root); err != nil {
		return result, err
	}
	if err := requireSetupForDataMutation(root, path, destination); err != nil {
		return result, err
	}
	if operation == "create" || operation == "update" {
		if err := validateSetupCompletion(root, path, []byte(content)); err != nil {
			return result, err
		}
	}

	var paths []string
	switch operation {
	case "create":
		changed, err := createFile(root, path, content)
		if err != nil {
			return result, err
		}
		paths = []string{changed}
	case "update":
		changed, err := updateFile(root, path, content)
		if err != nil {
			return result, err
		}
		paths = []string{changed}
	case "delete":
		changed, err := deleteFile(root, path)
		if err != nil {
			return result, err
		}
		paths = []string{changed}
	case "move":
		changed, err := moveFile(root, path, destination)
		if err != nil {
			return result, err
		}
		paths = changed
		result.Path = path + " -> " + destination
	default:
		return result, fmt.Errorf("unknown operation: %s", operation)
	}

	historyPath, err := writeHistory(root, summary, time.Now())
	if err != nil {
		return result, err
	}
	if historyPath != "" {
		paths = append(paths, historyPath)
	}
	if err := validateRepository(root); err != nil {
		return result, err
	}
	if strings.TrimSpace(commitMessage) == "" {
		commitMessage = "atlas: " + operation + " " + result.Path
	}
	commit, err := commitPaths(root, commitMessage, paths)
	if err != nil {
		return result, err
	}
	result.Commit = commit
	result.Branch = state.Branch

	if _, err := syncRepo(root); err != nil {
		return result, err
	}
	if err := validateRepository(root); err != nil {
		return result, err
	}
	if err := pushWithRetry(root, state.Branch); err != nil {
		return result, err
	}
	result.Commit, err = gitHead(root, "HEAD")
	return result, err
}
