package cli

import (
	"fmt"

	"secretlint/internal/scanner"
)

func scanStagedChanges() error {
	differ := scanner.NewGitDiffer()
	
	// Check if we're in a git repository
	if !differ.IsInGitRepo() {
		return fmt.Errorf("not in a git repository")
	}
	
	// Check if there are staged changes
	hasChanges, err := differ.HasStagedChanges()
	if err != nil {
		return fmt.Errorf("failed to check for staged changes: %w", err)
	}
	
	if !hasChanges {
		fmt.Println("✅ No staged changes to scan")
		return nil
	}
	
	// Get the staged changes
	lines, err := differ.GetStagedChanges()
	if err != nil {
		return fmt.Errorf("failed to get staged changes: %w", err)
	}
	
	if len(lines) == 0 {
		fmt.Println("✅ No new lines to scan")
		return nil
	}
	
	fmt.Printf("📄 Found %d added lines to scan\n", len(lines))
	
	// Initialize the secret scanner
	secretScanner := scanner.NewSecretScanner()
	
	// Scan all lines for secrets
	findings := secretScanner.ScanLines(lines)
	
	if len(findings) == 0 {
		fmt.Println("✅ No secrets detected in staged changes")
		return nil
	}
	
	// Report findings
	fmt.Printf("\n⛔ %d secret(s) detected in staged changes:\n\n", len(findings))
	
	for _, finding := range findings {
		fmt.Printf("Rule     : %s\n", finding.RuleID)
		fmt.Printf("File     : %s:%d\n", finding.FilePath, finding.LineNum)
		fmt.Printf("Snippet  : %s\n", finding.MaskSecret())
		fmt.Printf("Advice   : %s\n\n", finding.Advice)
	}
	
	fmt.Println("Commit aborted.")
	return fmt.Errorf("secrets detected - commit blocked")
}