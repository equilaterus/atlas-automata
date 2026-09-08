package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func writeHistory(root, summary string, now time.Time) (string, error) {
	if strings.TrimSpace(summary) == "" {
		return "", nil
	}
	summary = strings.Join(strings.Fields(summary), " ")
	relative := filepath.ToSlash(filepath.Join("log", now.UTC().Format("2006-01-02")+".md"))
	full := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", fmt.Errorf("create history directory: %w", err)
	}
	file, err := os.OpenFile(full, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return "", fmt.Errorf("open semantic history: %w", err)
	}
	defer file.Close()
	line := fmt.Sprintf("- %s — %s\n", now.UTC().Format(time.RFC3339), summary)
	if _, err := file.WriteString(line); err != nil {
		return "", fmt.Errorf("write semantic history: %w", err)
	}
	return relative, nil
}
