package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IgnoreChecker handles .secretignore file parsing and matching
type IgnoreChecker struct {
	patterns []string
	regexes  []*regexp.Regexp
}

// NewIgnoreChecker creates a new ignore checker
func NewIgnoreChecker() *IgnoreChecker {
	return &IgnoreChecker{
		patterns: make([]string, 0),
		regexes:  make([]*regexp.Regexp, 0),
	}
}

// LoadIgnoreFile loads patterns from .secretignore file
func (ic *IgnoreChecker) LoadIgnoreFile(ignoreFilePath string) error {
	file, err := os.Open(ignoreFilePath)
	if err != nil {
		// If .secretignore doesn't exist, that's OK
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open .secretignore: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Add pattern
		if err := ic.AddPattern(line); err != nil {
			// Log warning but continue processing
			fmt.Fprintf(os.Stderr, "Warning: invalid pattern '%s': %v\n", line, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .secretignore: %w", err)
	}

	return nil
}

// AddPattern adds a glob pattern to the ignore list
func (ic *IgnoreChecker) AddPattern(pattern string) error {
	// Convert glob pattern to regex
	regex, err := ic.globToRegex(pattern)
	if err != nil {
		return fmt.Errorf("invalid glob pattern: %w", err)
	}
	
	compiledRegex, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("failed to compile pattern: %w", err)
	}
	
	ic.patterns = append(ic.patterns, pattern)
	ic.regexes = append(ic.regexes, compiledRegex)
	
	return nil
}

// ShouldIgnore checks if a file path should be ignored
func (ic *IgnoreChecker) ShouldIgnore(filePath string) bool {
	// Normalize path separators for cross-platform compatibility
	normalizedPath := filepath.ToSlash(filePath)
	
	for _, regex := range ic.regexes {
		if regex.MatchString(normalizedPath) {
			return true
		}
		
		// Also check just the filename
		filename := filepath.Base(normalizedPath)
		if regex.MatchString(filename) {
			return true
		}
	}
	
	return false
}

// globToRegex converts a glob pattern to a regular expression
func (ic *IgnoreChecker) globToRegex(glob string) (string, error) {
	// Start with anchored regex
	regex := "^"
	
	i := 0
	for i < len(glob) {
		switch glob[i] {
		case '*':
			if i+1 < len(glob) && glob[i+1] == '*' {
				// Double star (**) - matches any number of directories
				if i+2 < len(glob) && glob[i+2] == '/' {
					// **/
					regex += "(?:.*/|^)"
					i += 3
				} else if i+2 == len(glob) {
					// ** at end
					regex += ".*"
					i += 2
				} else {
					// **something - treat as single *
					regex += "[^/]*"
					i++
				}
			} else {
				// Single star (*) - matches anything except /
				regex += "[^/]*"
				i++
			}
		case '?':
			// Question mark matches any single character except /
			regex += "[^/]"
			i++
		case '[':
			// Character class
			j := i + 1
			if j < len(glob) && glob[j] == '!' {
				j++
			}
			if j < len(glob) && glob[j] == ']' {
				j++
			}
			for j < len(glob) && glob[j] != ']' {
				j++
			}
			if j >= len(glob) {
				return "", fmt.Errorf("unterminated character class")
			}
			// Convert [!...] to [^...]
			class := glob[i:j+1]
			if strings.HasPrefix(class, "[!") {
				class = "[^" + class[2:]
			}
			regex += class
			i = j + 1
		case '.', '^', '$', '+', '{', '}', '(', ')', '|', '\\':
			// Escape regex special characters
			regex += "\\" + string(glob[i])
			i++
		default:
			regex += string(glob[i])
			i++
		}
	}
	
	regex += "$"
	return regex, nil
}

// GetPatterns returns the loaded patterns for debugging
func (ic *IgnoreChecker) GetPatterns() []string {
	return ic.patterns
}