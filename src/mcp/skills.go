package main

import (
	"os"
	"path/filepath"
	"sort"
)

func loadSkills(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "ai", "skills"))
	if os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var skills []string
	for _, entry := range entries {
		if entry.IsDir() {
			skills = append(skills, entry.Name())
		}
	}
	sort.Strings(skills)
	return skills, nil
}
