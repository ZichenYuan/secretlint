package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// DiffLine represents a line added in a git diff
type DiffLine struct {
	FilePath string
	LineNum  int
	Content  string
}

// GitDiffer handles extracting added lines from git diff
type GitDiffer struct{}

// NewGitDiffer creates a new GitDiffer instance
func NewGitDiffer() *GitDiffer {
	return &GitDiffer{}
}

// GetStagedChanges returns all added lines from staged changes
func (gd *GitDiffer) GetStagedChanges() ([]DiffLine, error) {
	// Execute git diff --cached -U0 to get staged changes with no context
	cmd := exec.Command("git", "diff", "--cached", "-U0")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git diff: %w", err)
	}

	return gd.parseDiff(string(output))
}

// parseDiff parses git diff output and extracts added lines
func (gd *GitDiffer) parseDiff(diffOutput string) ([]DiffLine, error) {
	var lines []DiffLine
	scanner := bufio.NewScanner(strings.NewReader(diffOutput))
	
	var currentFile string
	var currentLineNum int
	
	// Regex to match file headers: +++ b/path/to/file
	fileHeaderRegex := regexp.MustCompile(`^\+\+\+ b/(.+)$`)
	
	// Regex to match hunk headers: @@ -old_start,old_count +new_start,new_count @@
	hunkHeaderRegex := regexp.MustCompile(`^@@ -\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@`)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Check for file header
		if matches := fileHeaderRegex.FindStringSubmatch(line); matches != nil {
			currentFile = matches[1]
			continue
		}
		
		// Check for hunk header  
		if matches := hunkHeaderRegex.FindStringSubmatch(line); matches != nil {
			lineNum, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse line number: %w", err)
			}
			currentLineNum = lineNum
			continue
		}
		
		// Check for added lines (start with +)
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			// Remove the + prefix
			content := line[1:]
			
			lines = append(lines, DiffLine{
				FilePath: currentFile,
				LineNum:  currentLineNum,
				Content:  content,
			})
			
			currentLineNum++
		} else if strings.HasPrefix(line, " ") {
			// Context line (unchanged), increment line number
			currentLineNum++
		}
		// Ignore removed lines (start with -)
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading diff output: %w", err)
	}
	
	return lines, nil
}

// IsInGitRepo checks if current directory is inside a git repository
func (gd *GitDiffer) IsInGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// HasStagedChanges checks if there are any staged changes
func (gd *GitDiffer) HasStagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	if err != nil {
		// Exit code 1 means there are differences (staged changes)
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return true, nil
		}
		return false, fmt.Errorf("failed to check for staged changes: %w", err)
	}
	// Exit code 0 means no differences (no staged changes)
	return false, nil
}