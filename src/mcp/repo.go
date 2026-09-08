package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var protectedDirectories = []string{"data", "doc", "ai", "log"}

func findRepoRoot(start string) (string, error) {
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get current directory: %w", err)
		}
	}

	abs, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolve repository path: %w", err)
	}
	out, err := gitOutput(abs, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("find Git repository from %s: %w", abs, err)
	}
	return filepath.Clean(strings.TrimSpace(out)), nil
}

func cleanRelativePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("path must be relative: %s", path)
	}
	clean := filepath.Clean(path)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes the repository: %s", path)
	}
	return filepath.ToSlash(clean), nil
}

func isProtectedPath(path string) bool {
	path = filepath.ToSlash(filepath.Clean(path))
	if path == "AUTOMATIZER.md" {
		return true
	}
	for _, directory := range protectedDirectories {
		if strings.HasPrefix(path, directory+"/") {
			return true
		}
	}
	return false
}

func protectedPath(root, path string) (string, string, error) {
	relative, err := cleanRelativePath(path)
	if err != nil {
		return "", "", err
	}
	if !isProtectedPath(relative) {
		return "", "", fmt.Errorf("path is not protected Automatizer state: %s", relative)
	}

	full := filepath.Join(root, filepath.FromSlash(relative))
	parent := filepath.Dir(full)
	for {
		info, statErr := os.Lstat(parent)
		if statErr == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				return "", "", fmt.Errorf("path traverses a symlink: %s", relative)
			}
			break
		}
		if !os.IsNotExist(statErr) {
			return "", "", fmt.Errorf("inspect parent for %s: %w", relative, statErr)
		}
		if parent == root {
			break
		}
		parent = filepath.Dir(parent)
	}
	return relative, full, nil
}

func validateRepository(root string) error {
	gitRoot, err := findRepoRoot(root)
	if err != nil {
		return err
	}
	if filepath.Clean(gitRoot) != filepath.Clean(root) {
		return fmt.Errorf("configured root %s is not the repository root %s", root, gitRoot)
	}
	if out, err := gitOutput(root, "diff", "--name-only", "--diff-filter=U"); err != nil {
		return err
	} else if strings.TrimSpace(out) != "" {
		return fmt.Errorf("repository has unresolved conflicts: %s", strings.TrimSpace(out))
	}
	return nil
}

func ensureNoProtectedChanges(root string) error {
	out, err := gitOutput(root, "status", "--porcelain=v1", "-z", "--", "data", "doc", "ai", "log", "AUTOMATIZER.md")
	if err != nil {
		return fmt.Errorf("inspect protected state: %w", err)
	}
	if out != "" {
		return fmt.Errorf("protected state already has changes outside this Atlas operation")
	}
	return nil
}
