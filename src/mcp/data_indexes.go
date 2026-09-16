package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const dataIndexMarker = "<!-- atlas:data-index:v1 -->"

func isDataIndexPath(path string) bool {
	clean := filepath.ToSlash(filepath.Clean(path))
	return isDataPath(clean) && filepath.Base(clean) == "index.md"
}

func validateManagedDataIndexes(root string) error {
	dataRoot := filepath.Join(root, "data")
	if _, err := os.Stat(dataRoot); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("inspect data root: %w", err)
	}

	return filepath.WalkDir(dataRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Name() != "index.md" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read data index %s: %w", path, err)
		}
		if !strings.HasPrefix(string(content), dataIndexMarker+"\n") {
			relative, _ := filepath.Rel(root, path)
			return fmt.Errorf("data index is reserved for Atlas but is not managed by Atlas: %s", filepath.ToSlash(relative))
		}
		return nil
	})
}

func rebuildDataIndexes(root string) ([]string, error) {
	dataRoot := filepath.Join(root, "data")
	if _, err := os.Stat(dataRoot); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("inspect data root: %w", err)
	}

	var directories []string
	if err := filepath.WalkDir(dataRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("data tree contains a symlink: %s", path)
		}
		if entry.IsDir() {
			directories = append(directories, path)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("inspect data directories: %w", err)
	}

	sort.Strings(directories)
	var changed []string
	for _, directory := range directories {
		content, err := renderDataIndex(root, directory)
		if err != nil {
			return nil, err
		}
		indexPath := filepath.Join(directory, "index.md")
		existing, err := os.ReadFile(indexPath)
		if err == nil {
			if !strings.HasPrefix(string(existing), dataIndexMarker+"\n") {
				relative, _ := filepath.Rel(root, indexPath)
				return nil, fmt.Errorf("data index is reserved for Atlas but is not managed by Atlas: %s", filepath.ToSlash(relative))
			}
			if string(existing) == content {
				continue
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read data index %s: %w", indexPath, err)
		}
		if err := os.WriteFile(indexPath, []byte(content), 0o644); err != nil {
			return nil, fmt.Errorf("write data index %s: %w", indexPath, err)
		}
		relative, err := filepath.Rel(root, indexPath)
		if err != nil {
			return nil, err
		}
		changed = append(changed, filepath.ToSlash(relative))
	}
	return changed, nil
}

func renderDataIndex(root, directory string) (string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("read data directory %s: %w", directory, err)
	}
	relative, err := filepath.Rel(root, directory)
	if err != nil {
		return "", err
	}

	var folders []string
	var records []string
	for _, entry := range entries {
		if entry.Name() == "index.md" {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("data directory contains a symlink: %s", filepath.Join(directory, entry.Name()))
		}
		if entry.IsDir() {
			folders = append(folders, entry.Name())
			continue
		}
		if entry.Type().IsRegular() {
			records = append(records, entry.Name())
		}
	}
	sort.Strings(folders)
	sort.Strings(records)

	var builder strings.Builder
	builder.WriteString(dataIndexMarker)
	builder.WriteString("\n# ")
	builder.WriteString(filepath.ToSlash(relative))
	builder.WriteString("\n")
	if len(folders) > 0 {
		builder.WriteString("\n## Directories\n\n")
		for _, name := range folders {
			fmt.Fprintf(&builder, "- [%s/](%s/index.md)\n", markdownLabel(name), url.PathEscape(name))
		}
	}
	if len(records) > 0 {
		builder.WriteString("\n## Records\n\n")
		for _, name := range records {
			fmt.Fprintf(&builder, "- [%s](%s)\n", markdownLabel(strings.TrimSuffix(name, filepath.Ext(name))), url.PathEscape(name))
		}
	}
	if len(folders) == 0 && len(records) == 0 {
		builder.WriteString("\n_Empty._\n")
	}
	return builder.String(), nil
}

func markdownLabel(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]")
	return replacer.Replace(value)
}
