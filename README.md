# Secretlint Lite

🛡️ **Lightweight secret detection for Git commits** - Prevent API keys and secrets from being committed to your repository.

## 🎯 Goal

Prevent obvious API keys and secrets from being committed to Git, in a way that's lightweight, developer-friendly, and installable in <30 seconds.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/yourusername/secretlint
cd secretlint

# Build the tool
go build -o secretlint cmd/secretlint/main.go

# Test scanning
./secretlint scan

# Help
./secretlint --help
```

## ✅ Current Features (v0.1 MVP - Implemented)

### Core Functionality
- **Git Integration**: Scans only staged changes (`git diff --cached`) for performance
- **Secret Detection**: 10 curated regex patterns for high-value secrets:
  - OpenAI API keys: `sk-[A-Za-z0-9]{20,}`
  - GitHub PAT: `ghp_[A-Za-z0-9]{36}`
  - AWS Access Key: `(AKIA|ASIA)[A-Z0-9]{16}`
  - AWS Secret Key: `(?i)aws(.{0,20})?(secret|access).{0,20}['\"][A-Za-z0-9/+=]{40}['\"]`
  - Stripe keys: `pk_live_[A-Za-z0-9]{24}`, `sk_live_[A-Za-z0-9]{24}`
  - Slack tokens: `xox[baprs]-[0-9A-Za-z-]+`
  - JWTs: `eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9._-]+\.[A-Za-z0-9._-]+`
  - Generic API key patterns
  - Private keys: `-----BEGIN (RSA )?PRIVATE KEY-----`

### CLI Commands
- `secretlint scan` - Scan staged changes for secrets
- `secretlint --help` - Show usage information

### Security Features
- **Safe Display**: Masks secrets in output (shows first 4 + last 4 chars)
- **Commit Blocking**: Returns exit code 1 to prevent commits when secrets found
- **Detailed Reporting**: Shows file:line, rule matched, and remediation advice

### Performance
- Fast scanning (<200ms on typical commits)
- Only scans added lines from staged changes
- No external dependencies beyond Go standard library

## 📋 Development Plan

### ✅ Completed (Current State)
1. **Project Setup**: Go module, basic CLI structure
2. **Git Integration**: Diff parser for staged changes only
3. **Secret Detection Engine**: 10 curated regex patterns
4. **Basic Reporting**: Terminal output with masked secrets

### 🔄 Next Steps (Remaining v0.1)
1. **Entropy Detection**: Shannon entropy > 4.0 for high-entropy strings (len > 20)
2. **Ignore System**: `.secretignore` file support (glob-based patterns)
3. **Pretty Output**: Colored terminal output with ❌ and ⚠️ emojis
4. **Init Command**: `secretlint init` to setup config files and git hooks
5. **Pre-commit Hook**: Bash script integration
6. **Default Configs**: `.secretlintrc.yml` and `.secretignore` templates

### 🗺️ Future Roadmap

#### v0.2 - Enhanced Usability
- Inline `# secretlint:allow RULE_ID ttl=7d` comments
- `--explain RULE_ID` command for rule documentation
- Configuration validation and better error messages

#### v0.3 - Developer Experience  
- Lightweight editor bridge (VS Code extension)
- `--fix` helpers to move secrets to `.env` files
- Auto-replacement with `process.env` references

#### v0.4 - Enterprise Features
- Pluggable rules via YAML/JSON configuration  
- Organization presets without requiring a server
- Custom rule definitions

#### v1.0 - Production Ready
- Pre-commit framework integration
- Cross-platform builds (Windows/macOS/Linux)
- Simple telemetry opt-in
- Package manager distribution (Homebrew, etc.)

## 🏗️ Architecture

```
secretlint/
├── cmd/secretlint/        # CLI entry point
├── internal/
│   ├── cli/              # Command handlers (scan, init)
│   ├── scanner/          # Core scanning logic
│   │   ├── differ.go     # Git diff parser ✅
│   │   ├── rules.go      # Regex rules & findings ✅
│   │   └── entropy.go    # Shannon entropy detection (TODO)
│   ├── reporter/         # Output formatting (TODO)
│   └── config/          # Configuration management (TODO)
├── templates/           # Default config templates (TODO)
└── README.md
```

## 🧪 Testing

Current test coverage includes:
- Git diff parsing with various change types
- Secret detection across all 10 rule categories
- Clean file scanning (no false positives)
- Staged vs unstaged change handling

```bash
# Test with sample secrets
echo 'API_KEY=sk-abc123def456ghi789jklmno' > test.txt
git add test.txt
go run cmd/secretlint/main.go scan
```

## 🤝 Contributing

This is an MVP implementation focused on core functionality. Key areas for contribution:

1. **Performance optimization** - Make scanning even faster
2. **Rule improvements** - Reduce false positives, add new secret types  
3. **Integration testing** - More comprehensive test scenarios
4. **Documentation** - Usage examples and best practices

## 📄 License

MIT License - See LICENSE file for details.

## 🔮 Design Principles

- **Lightweight**: Fast startup, minimal dependencies
- **Developer-friendly**: Clear output, helpful error messages
- **Git-native**: Works seamlessly with existing Git workflows
- **Secure by default**: Conservative detection, safe secret masking
- **Extensible**: Clean architecture for adding new rules and features

---

**Status**: 🚧 **Active Development** - Core v0.1 features implemented, working toward full MVP completion.