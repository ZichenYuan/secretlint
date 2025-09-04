package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runInit() error {
	fmt.Println("🔧 Initializing secretlint...")
	
	// Check if we're in a git repository
	if err := checkGitRepository(); err != nil {
		return err
	}
	
	// Find the current secretlint binary path
	binaryPath, err := findCurrentBinary()
	if err != nil {
		return fmt.Errorf("failed to locate secretlint binary: %w", err)
	}
	
	// Create configuration files
	if err := createConfigFiles(); err != nil {
		return fmt.Errorf("failed to create config files: %w", err)
	}
	
	// Install pre-commit hook with stored binary path
	if err := installPreCommitHook(binaryPath); err != nil {
		return fmt.Errorf("failed to install pre-commit hook: %w", err)
	}
	
	fmt.Println("✅ Secretlint initialized successfully!")
	fmt.Println("")
	fmt.Println("Created files:")
	fmt.Println("  📄 .secretlintrc.yml - Configuration and rules")
	fmt.Println("  🚫 .secretignore - Files and patterns to ignore")
	fmt.Println("  🪝 .git/hooks/pre-commit - Git hook integration")
	fmt.Printf("  ⚙️  .git/hooks/secretlint-config - Binary path (%s)\n", binaryPath)
	fmt.Println("")
	fmt.Println("Try making a commit with secrets to test it:")
	fmt.Println("  echo 'API_KEY=sk-abc123' > test.txt")
	fmt.Println("  git add test.txt && git commit -m 'test'")
	
	return nil
}

func checkGitRepository() error {
	if _, err := os.Stat(".git"); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("not in a git repository - please run 'git init' first")
		}
		return fmt.Errorf("error checking git repository: %w", err)
	}
	return nil
}

func findCurrentBinary() (string, error) {
	// Try to find secretlint binary in order of preference
	
	// 1. Check if it's in PATH
	if path, err := exec.LookPath("secretlint"); err == nil {
		if absPath, err := filepath.Abs(path); err == nil {
			return absPath, nil
		}
		return path, nil
	}
	
	// 2. Check current directory
	if _, err := os.Stat("./secretlint"); err == nil {
		absPath, err := filepath.Abs("./secretlint")
		if err != nil {
			return "./secretlint", nil
		}
		return absPath, nil
	}
	
	// 3. Check parent directories
	for _, path := range []string{"../secretlint", "../../secretlint"} {
		if _, err := os.Stat(path); err == nil {
			if absPath, err := filepath.Abs(path); err == nil {
				return absPath, nil
			}
			return path, nil
		}
	}
	
	// 4. If we're running via 'go run', suggest building
	if _, err := os.Stat("cmd/secretlint/main.go"); err == nil {
		return "", fmt.Errorf("please build the binary first: go build -o secretlint cmd/secretlint/main.go")
	}
	
	// 5. Check parent for source
	if _, err := os.Stat("../cmd/secretlint/main.go"); err == nil {
		return "", fmt.Errorf("please build the binary first: cd .. && go build -o secretlint cmd/secretlint/main.go")
	}
	
	return "", fmt.Errorf("secretlint binary not found. Please build it first: go build -o secretlint cmd/secretlint/main.go")
}

func createConfigFiles() error {
	// Create .secretlintrc.yml
	configContent := `# Secretlint configuration file
# See https://github.com/ZichenYuan/secretlint for documentation

# Enable/disable specific rules
rules:
  OPENAI_API_KEY: true
  GITHUB_PAT: true
  AWS_ACCESS_KEY: true
  AWS_SECRET_KEY: true
  STRIPE_LIVE_PK: true
  STRIPE_LIVE_SK: true
  SLACK_TOKEN: true
  JWT_TOKEN: true
  GENERIC_API_KEY: true
  PRIVATE_KEY: true

# Global settings
settings:
  # Exit with code 1 when secrets are found (blocks commits)
  fail_on_detection: true
  
  # Show detailed output
  verbose: false
  
  # Minimum secret length to scan
  min_length: 10

# Custom patterns (future feature)
custom_rules: []
`

	if err := writeFileIfNotExists(".secretlintrc.yml", configContent); err != nil {
		return err
	}

	// Create .secretignore if it doesn't exist (we may already have one)
	ignoreContent := `# Secretlint ignore patterns
# Patterns use glob syntax similar to .gitignore
# Lines starting with # are comments

# Dependencies and build outputs
node_modules/
dist/
build/
target/
vendor/
.git/

# Test files
*test*
*Test*
*TEST*
**/*test*
**/*Test*
**/*TEST*

# Documentation
*.md
*.rst
*.txt
README*
CHANGELOG*
LICENSE*

# Configuration files (review these carefully!)
package.json
package-lock.json
yarn.lock
Dockerfile
docker-compose.yml

# IDE and editor files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
Thumbs.db

# Log files
*.log
**/*.log

# Temporary files
*.tmp
*.temp
.cache/
`

	if err := writeFileIfNotExists(".secretignore", ignoreContent); err != nil {
		return err
	}

	return nil
}

func installPreCommitHook(binaryPath string) error {
	hookPath := ".git/hooks/pre-commit"
	configPath := ".git/hooks/secretlint-config"
	
	// Create the hooks directory if it doesn't exist
	hookDir := filepath.Dir(hookPath)
	if err := os.MkdirAll(hookDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}
	
	// Always write/update the config with current binary path
	if err := writeSecretlintConfig(configPath, binaryPath); err != nil {
		return err
	}
	
	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		fmt.Println("⚠️  Pre-commit hook already exists")
		
		// Read existing hook content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			return fmt.Errorf("failed to read existing hook: %w", err)
		}
		
		hookContent := string(content)
		
		// Check for secretlint signature markers
		hasSecretlintMarker := strings.Contains(hookContent, "# Secretlint pre-commit hook")
		hasConfigSource := strings.Contains(hookContent, "source .git/hooks/secretlint-config") ||
			strings.Contains(hookContent, ". .git/hooks/secretlint-config")
		
		if hasSecretlintMarker && hasConfigSource {
			fmt.Println("✅ Secretlint is already integrated in pre-commit hook")
			fmt.Println("✅ Updated binary path configuration")
			return nil
		}
		
		fmt.Println("📝 Backing up existing pre-commit hook to pre-commit.backup")
		if err := os.Rename(hookPath, hookPath+".backup"); err != nil {
			return fmt.Errorf("failed to backup existing hook: %w", err)
		}
		fmt.Println("💡 You can merge your custom hook logic with the new secretlint hook if needed")
	}
	
	// Write the pre-commit hook
	if err := os.WriteFile(hookPath, []byte(getPreCommitHookContent()), 0755); err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}
	
	fmt.Println("✅ Created pre-commit hook")
	return nil
}

func writeSecretlintConfig(configPath, binaryPath string) error {
	configContent := fmt.Sprintf(`#!/bin/sh
# Secretlint configuration - stores binary path
# Generated automatically by 'secretlint init'

export SECRETLINT_BINARY="%s"
`, binaryPath)

	if err := os.WriteFile(configPath, []byte(configContent), 0755); err != nil {
		return fmt.Errorf("failed to write secretlint config: %w", err)
	}
	
	return nil
}

func writeFileIfNotExists(filename, content string) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("⚠️  %s already exists, skipping\n", filename)
		return nil
	}
	
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filename, err)
	}
	
	fmt.Printf("✅ Created %s\n", filename)
	return nil
}

func getPreCommitHookContent() string {
	return `#!/bin/sh
#
# Secretlint pre-commit hook
# Automatically scans staged changes for secrets
#

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔍 Scanning staged changes for secrets..."

# Load secretlint configuration (binary path)
if [ -f ".git/hooks/secretlint-config" ]; then
    source .git/hooks/secretlint-config
fi

# Find secretlint binary using stored path first, then fallback
SECRETLINT=""
if [ -n "$SECRETLINT_BINARY" ] && [ -f "$SECRETLINT_BINARY" ]; then
    SECRETLINT="$SECRETLINT_BINARY"
elif command -v secretlint >/dev/null 2>&1; then
    SECRETLINT="secretlint"
elif [ -f "./secretlint" ]; then
    SECRETLINT="./secretlint"
else
    echo "${RED}❌ secretlint binary not found${NC}"
    echo "Stored path: $SECRETLINT_BINARY"
    echo "Please run 'secretlint init' again or build the binary:"
    echo "  go build -o secretlint cmd/secretlint/main.go"
    exit 1
fi

# Run secretlint scan
$SECRETLINT scan

# Check exit code
if [ $? -ne 0 ]; then
    echo ""
    echo "${RED}❌ Commit blocked due to secrets detected${NC}"
    echo ""
    echo "To bypass this check (NOT recommended):"
    echo "  git commit --no-verify -m 'your message'"
    echo ""
    echo "To fix the issue:"
    echo "  1. Move secrets to environment variables"
    echo "  2. Add files to .secretignore if they're false positives"
    echo "  3. Remove secrets from the code"
    exit 1
else
    echo "${GREEN}✅ No secrets detected${NC}"
    exit 0
fi
`
}