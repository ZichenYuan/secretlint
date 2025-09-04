# Secretlint - Git Secret Detection Tool

🛡️ **Prevent API keys and secrets from being committed to your Git repository**

Secretlint is a lightweight, fast tool that automatically scans your Git commits for secrets like API keys, tokens, and credentials. It integrates seamlessly with your Git workflow to catch secrets before they reach your repository.

## 🚀 Quick Start

### 1. Build Secretlint
```bash
# Clone the repository
git clone https://github.com/ZichenYuan/secretlint
cd secretlint

# Build the binary
go build -o secretlint cmd/secretlint/main.go
```

### 2. Initialize in Your Project
```bash
# Navigate to your Git repository
cd /path/to/your/project

# Initialize secretlint (creates config files and pre-commit hook)
/path/to/secretlint init
```

### 3. Test It Works
```bash
# Create a file with a fake secret
echo 'const apiKey = "sk-1234567890abcdefghijklmnopqrstuvwxyz";' > api.js

# Try to commit it
git add api.js
git commit -m "test commit"

# Secretlint should block the commit and show the detected secret
```

## 📖 Detailed User Guide

### Installation & Setup

#### Step 1: Build the Tool
```bash
# Ensure you have Go installed (1.19+)
go version

# Clone and build
git clone https://github.com/ZichenYuan/secretlint
cd secretlint
go build -o secretlint cmd/secretlint/main.go

# Optional: Move to PATH for global access
sudo mv secretlint /usr/local/bin/
```

#### Step 2: Initialize in Your Repository
```bash
# Navigate to your Git repository
cd your-project

# Run init command
secretlint init
```

**What `init` does:**
- ✅ Creates `.secretlintrc.yml` - Configuration file with rules
- ✅ Creates `.secretignore` - Files/patterns to ignore  
- ✅ Installs Git pre-commit hook that automatically scans commits
- ✅ Stores the secretlint binary path for the hook to use

#### Step 3: Verify Installation
```bash
# Check that files were created
ls -la .secretlint* .git/hooks/pre-commit

# Test manual scanning
secretlint scan
```

### Daily Usage

#### Automatic Protection (Recommended)
Once initialized, secretlint works automatically:

```bash
# Normal Git workflow - secretlint runs automatically
git add .
git commit -m "your changes"

# If secrets detected:
# ❌ Commit blocked with details
# ✅ If no secrets: commit proceeds normally
```

#### Manual Scanning
Scan staged changes without committing:

```bash
# Scan currently staged changes
secretlint scan

# Example output:
# 🔍 Scanning for secrets...
# 📄 Found 15 added lines to scan
# ⛔ 1 secret(s) detected in staged changes:
# 
# Rule     : OPENAI_API_KEY
# File     : src/config.js:12
# Snippet  : sk-proj**********************xyz123
# Advice   : Move this to an environment variable (.env file) and add .env to .gitignore
```

#### Bypassing Protection (Not Recommended)
```bash
# Skip secretlint check (emergency use only)
git commit --no-verify -m "emergency commit"
```

### Configuration

#### `.secretlintrc.yml` - Rules Configuration
```yaml
# Enable/disable specific rules
rules:
  OPENAI_API_KEY: true      # sk-[A-Za-z0-9]{20,}
  GITHUB_PAT: true          # ghp_[A-Za-z0-9]{36}
  AWS_ACCESS_KEY: true      # (AKIA|ASIA)[A-Z0-9]{16}
  AWS_SECRET_KEY: true      # AWS secret patterns
  STRIPE_LIVE_SK: true      # sk_live_[A-Za-z0-9]{24}
  SLACK_TOKEN: true         # xox[baprs]-[0-9A-Za-z-]+
  JWT_TOKEN: true           # eyJ... patterns
  GENERIC_API_KEY: true     # Common API key patterns
  PRIVATE_KEY: true         # -----BEGIN PRIVATE KEY-----

# Global settings
settings:
  fail_on_detection: true   # Block commits when secrets found
  verbose: false           # Show detailed output
  min_length: 10          # Minimum secret length to check
```

#### `.secretignore` - Ignore Patterns
Use glob patterns to exclude files from scanning:

```bash
# Dependencies and build outputs
node_modules/
dist/
build/
*.log

# Test files (often contain fake secrets)
*test*
*Test*
**/*test*

# Documentation
*.md
*.rst
README*

# Configuration files (review carefully!)
package.json
Dockerfile

# Temporary files
*.tmp
.cache/
```

### Understanding Output

#### Clean Scan (No Secrets)
```
🔍 Scanning for secrets...
📄 Found 8 added lines to scan
✅ No secrets detected in staged changes
```

#### Secrets Detected
```
🔍 Scanning for secrets...
📄 Found 12 added lines to scan
🚫 Ignored files:
   README.md (3 lines)
   test-config.js (2 lines)

⛔ 2 secret(s) detected in staged changes:

Rule     : OPENAI_API_KEY
File     : src/api.js:15
Snippet  : sk-proj**********************abc123
Advice   : Move this to an environment variable (.env file) and add .env to .gitignore

Rule     : AWS_ACCESS_KEY
File     : config/aws.js:8
Snippet  : AKIA****************EXAM
Advice   : Use AWS IAM roles or store in AWS credentials file/environment variables

Commit aborted.
```

### What Secretlint Detects

| Secret Type | Pattern | Example |
|-------------|---------|---------|
| **OpenAI API Keys** | `sk-[A-Za-z0-9]{20,}` | `sk-proj-abc123...` |
| **GitHub PAT** | `ghp_[A-Za-z0-9]{36}` | `ghp_1234567890abcdef...` |
| **AWS Access Key** | `(AKIA\|ASIA)[A-Z0-9]{16}` | `AKIAIOSFODNN7EXAMPLE` |
| **AWS Secret Key** | Context-aware patterns | `aws_secret_access_key = "..."` |
| **Stripe Keys** | `sk_live_[A-Za-z0-9]{24}` | `sk_live_abc123def456...` |
| **Slack Tokens** | `xox[baprs]-[0-9A-Za-z-]+` | `xoxb-1234567890-...` |
| **JWT Tokens** | `eyJ[A-Za-z0-9_-]+\.\.\.` | `eyJhbGciOiJIUzI1NiI...` |
| **Private Keys** | `-----BEGIN.*PRIVATE KEY-----` | RSA/SSH private keys |
| **Generic API Keys** | Common patterns | `api_key = "abc123..."` |

### Troubleshooting

#### "secretlint binary not found" Error
```bash
# Check if binary exists
which secretlint
ls -la ./secretlint

# Rebuild and reinitialize
go build -o secretlint cmd/secretlint/main.go
./secretlint init
```

#### Pre-commit Hook Not Working
```bash
# Check hook exists and is executable
ls -la .git/hooks/pre-commit
cat .git/hooks/pre-commit

# Reinitialize
secretlint init
```

#### False Positives
```bash
# Add file to .secretignore
echo "false-positive-file.js" >> .secretignore

# Or add pattern for file type
echo "*.example" >> .secretignore
```

#### Hook Conflicts with Other Tools
The init command detects existing pre-commit hooks and backs them up:
```bash
# Your original hook is saved as:
ls -la .git/hooks/pre-commit.backup

# You can merge them manually if needed
```

### Best Practices

#### 1. Initialize Early
```bash
# Add secretlint to new projects immediately
git init
secretlint init
```

#### 2. Team Setup
```bash
# Commit secretlint config files to share with team
git add .secretlintrc.yml .secretignore
git commit -m "Add secretlint configuration"

# Team members then run:
secretlint init  # Uses existing config files
```

#### 3. CI/CD Integration
```bash
# Add to CI pipeline
secretlint scan
if [ $? -ne 0 ]; then
  echo "Secrets detected in codebase"
  exit 1
fi
```

#### 4. Regular Maintenance
```bash
# Periodically review and update ignore patterns
vim .secretignore

# Update rules as needed
vim .secretlintrc.yml
```

### Advanced Usage

#### Custom Ignore Patterns
```bash
# Ignore specific directories
echo "vendor/" >> .secretignore
echo "third-party/**" >> .secretignore

# Ignore file types
echo "*.min.js" >> .secretignore
echo "*.bundle.*" >> .secretignore

# Ignore by filename pattern
echo "*-generated.py" >> .secretignore
echo "test-*" >> .secretignore
```

#### Temporary Bypass
```bash
# For emergency commits (use sparingly)
git commit --no-verify -m "emergency fix"

# Better: fix the issue properly
git reset --soft HEAD~1  # Undo last commit
# Remove secrets, then commit normally
```

## 🔧 Commands Reference

| Command | Description | Example |
|---------|-------------|---------|
| `secretlint init` | Setup config files and pre-commit hook | `secretlint init` |
| `secretlint scan` | Scan staged changes for secrets | `secretlint scan` |
| `secretlint --help` | Show help and usage information | `secretlint --help` |

## 🚨 What to Do When Secrets Are Detected

### 1. **Don't Panic** - The secret hasn't been committed yet

### 2. **Fix the Issue**
```bash
# Option A: Move to environment variables
# Before:
const apiKey = "sk-1234567890abcdef";

# After:
const apiKey = process.env.OPENAI_API_KEY;

# Add to .env file (and .env to .gitignore)
echo "OPENAI_API_KEY=sk-1234567890abcdef" >> .env
echo ".env" >> .gitignore
```

```bash
# Option B: Remove the secret entirely
# If it's test/example code, use fake values:
const apiKey = "sk-fake-key-for-testing";
```

```bash
# Option C: Add to .secretignore if it's a false positive
echo "path/to/false-positive-file.js" >> .secretignore
```

### 3. **Commit Again**
```bash
git add .
git commit -m "your changes"  # Should now pass
```

### 4. **If Secret Was Already Committed**
```bash
# Remove from Git history (use with caution)
git filter-branch --force --index-filter 'git rm --cached --ignore-unmatch path/to/file' --prune-empty --tag-name-filter cat -- --all

# Or use BFG Repo Cleaner for large repositories
# https://rtyley.github.io/bfg-repo-cleaner/
```

## 🎯 Why Use Secretlint?

- **🚀 Fast**: Scans only changed lines, < 200ms typically
- **🎯 Accurate**: 10 carefully tuned patterns, low false positives  
- **🔧 Simple**: One command setup, zero configuration needed
- **🛡️ Automatic**: Pre-commit hook catches secrets before they're committed
- **📁 Smart**: Ignores test files, docs, and dependencies by default
- **🔍 Transparent**: Shows exactly what's being ignored and why
- **⚙️ Reliable**: Stores binary path, works even if you move the tool

---

**Need help?** Open an issue at https://github.com/ZichenYuan/secretlint/issues