package main

import (
	"context"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type EmptyInput struct{}

type StatusOutput struct {
	Root   string   `json:"root"`
	Branch string   `json:"branch"`
	Head   string   `json:"head"`
	Clean  bool     `json:"clean"`
	Setup  string   `json:"setup"`
	Skills []string `json:"skills"`
}

type FileInput struct {
	Path          string `json:"path" jsonschema:"repository-relative path in protected Automatizer state"`
	Content       string `json:"content" jsonschema:"complete UTF-8 file content"`
	Summary       string `json:"summary,omitempty" jsonschema:"optional human-readable semantic history entry"`
	CommitMessage string `json:"commit_message,omitempty" jsonschema:"optional Git commit message"`
}

type DeleteInput struct {
	Path          string `json:"path" jsonschema:"repository-relative path in protected Automatizer state"`
	Summary       string `json:"summary,omitempty" jsonschema:"optional human-readable semantic history entry"`
	CommitMessage string `json:"commit_message,omitempty" jsonschema:"optional Git commit message"`
}

type MoveInput struct {
	Source        string `json:"source" jsonschema:"existing repository-relative protected file"`
	Destination   string `json:"destination" jsonschema:"new repository-relative protected file"`
	Summary       string `json:"summary,omitempty" jsonschema:"optional human-readable semantic history entry"`
	CommitMessage string `json:"commit_message,omitempty" jsonschema:"optional Git commit message"`
}

func newMCPServer(root string) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "atlas-automata", Version: atlasVersion()}, nil)
	var mutationLock sync.Mutex

	mcp.AddTool(server, &mcp.Tool{
		Name:        "atlas_status",
		Description: "Inspect the configured Automatizer repository without changing it.",
		Annotations: toolAnnotations(true, false, true, false),
	},
		func(_ context.Context, _ *mcp.CallToolRequest, _ EmptyInput) (*mcp.CallToolResult, StatusOutput, error) {
			branch, err := currentBranch(root)
			if err != nil {
				return nil, StatusOutput{}, err
			}
			head, err := gitHead(root, "HEAD")
			if err != nil {
				return nil, StatusOutput{}, err
			}
			status, err := gitOutput(root, "status", "--porcelain")
			if err != nil {
				return nil, StatusOutput{}, err
			}
			skills, err := loadSkills(root)
			if err != nil {
				return nil, StatusOutput{}, err
			}
			setup, err := setupState(root)
			if err != nil {
				return nil, StatusOutput{}, err
			}
			return nil, StatusOutput{Root: root, Branch: branch, Head: head, Clean: status == "", Setup: setup, Skills: skills}, nil
		})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "atlas_sync",
		Description: "Fetch origin and merge the current remote branch without rebasing.",
		Annotations: toolAnnotations(false, false, true, true),
	},
		func(_ context.Context, _ *mcp.CallToolRequest, _ EmptyInput) (*mcp.CallToolResult, SyncState, error) {
			mutationLock.Lock()
			defer mutationLock.Unlock()
			state, err := syncRepo(root)
			return nil, state, err
		})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "atlas_create",
		Description: "Create one protected file, refresh managed data folder indexes, record history, commit, synchronize, and push.",
		Annotations: toolAnnotations(false, false, false, true),
	},
		func(_ context.Context, _ *mcp.CallToolRequest, input FileInput) (*mcp.CallToolResult, MutationResult, error) {
			mutationLock.Lock()
			defer mutationLock.Unlock()
			result, err := runMutation(root, "create", input.Path, "", input.Content, input.Summary, input.CommitMessage)
			return nil, result, err
		})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "atlas_update",
		Description: "Replace one protected file, refresh managed data folder indexes, record history, commit, synchronize, and push.",
		Annotations: toolAnnotations(false, true, false, true),
	},
		func(_ context.Context, _ *mcp.CallToolRequest, input FileInput) (*mcp.CallToolResult, MutationResult, error) {
			mutationLock.Lock()
			defer mutationLock.Unlock()
			result, err := runMutation(root, "update", input.Path, "", input.Content, input.Summary, input.CommitMessage)
			return nil, result, err
		})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "atlas_delete",
		Description: "Delete one protected file, refresh managed data folder indexes, record history, commit, synchronize, and push.",
		Annotations: toolAnnotations(false, true, false, true),
	},
		func(_ context.Context, _ *mcp.CallToolRequest, input DeleteInput) (*mcp.CallToolResult, MutationResult, error) {
			mutationLock.Lock()
			defer mutationLock.Unlock()
			result, err := runMutation(root, "delete", input.Path, "", "", input.Summary, input.CommitMessage)
			return nil, result, err
		})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "atlas_move",
		Description: "Move one protected file, refresh managed data folder indexes, record history, commit, synchronize, and push.",
		Annotations: toolAnnotations(false, true, false, true),
	},
		func(_ context.Context, _ *mcp.CallToolRequest, input MoveInput) (*mcp.CallToolResult, MutationResult, error) {
			mutationLock.Lock()
			defer mutationLock.Unlock()
			result, err := runMutation(root, "move", input.Source, input.Destination, "", input.Summary, input.CommitMessage)
			return nil, result, err
		})

	return server
}

func toolAnnotations(readOnly, destructive, idempotent, openWorld bool) *mcp.ToolAnnotations {
	return &mcp.ToolAnnotations{
		ReadOnlyHint:    readOnly,
		DestructiveHint: boolPointer(destructive),
		IdempotentHint:  idempotent,
		OpenWorldHint:   boolPointer(openWorld),
	}
}

func boolPointer(value bool) *bool {
	return &value
}
