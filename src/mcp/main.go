package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	rootFlag := flag.String("root", "", "Automatizer repository root (defaults to the current repository)")
	versionFlag := flag.Bool("version", false, "print the Atlas version")
	syncFlag := flag.Bool("sync", false, "synchronize the repository once and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(atlasVersion())
		return
	}

	root, err := findRepoRoot(*rootFlag)
	if err != nil {
		log.Fatal(err)
	}
	if *syncFlag {
		state, err := syncRepo(root)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("synchronized %s at %s\n", state.Branch, state.LocalHead)
		return
	}

	server := newMCPServer(root)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
