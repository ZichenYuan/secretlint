package cli

import (
	"fmt"
	"os"
)

func Execute() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: secretlint <command>\n\nCommands:\n  init    Setup secretlint in current repository\n  scan    Scan staged changes for secrets")
	}

	command := os.Args[1]
	
	switch command {
	case "init":
		return runInit()
	case "scan":
		return runScan(os.Args[2:])
	case "--help", "-h":
		fmt.Println("secretlint - Lightweight secret detection for Git")
		fmt.Println("\nUsage: secretlint <command>")
		fmt.Println("\nCommands:")
		fmt.Println("  init    Setup secretlint in current repository")
		fmt.Println("  scan    Scan staged changes for secrets")
		fmt.Println("\nOptions:")
		fmt.Println("  --staged    Scan only staged changes (default for scan)")
		return nil
	default:
		return fmt.Errorf("unknown command: %s\n\nRun 'secretlint --help' for usage", command)
	}
}

func runInit() error {
	fmt.Println("🔧 Initializing secretlint...")
	return fmt.Errorf("init command not implemented yet")
}

func runScan(args []string) error {
	fmt.Println("🔍 Scanning for secrets...")
	
	// Import and use the git differ
	return scanStagedChanges()
}