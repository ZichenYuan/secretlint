Product Design Plan: Secretlint Lite (MVP)
🎯 Goal

Prevent obvious API keys / secrets from being committed to Git, in a way that’s lightweight, developer-friendly, and installable in <30s.

1. Core Requirements

Scope: Only scan staged changes (added lines from git diff --cached).

Rules: Start with ~10 regexes for high-value secrets:

OpenAI API keys: sk-[A-Za-z0-9]{20,}

GitHub PAT: ghp_[A-Za-z0-9]{20,}

AWS Access Key: (AKIA|ASIA)[A-Z0-9]{16}

AWS Secret Key: (?i)aws(.{0,20})?(secret|access).{0,20}['"][A-Za-z0-9/+=]{40}['"]

Stripe keys: rk_live_[A-Za-z0-9]{24}, sk_live_[A-Za-z0-9]{24}

Slack tokens: xox[baprs]-[0-9A-Za-z-]+

JWTs: eyJ[A-Za-z0-9_-]+?\.[A-Za-z0-9._-]+?\.[A-Za-z0-9._-]+

Generic high-entropy strings (Shannon entropy > 4.0, len > 20)

Performance: Must run in <200ms on most commits.

UX: Print file:line, the offending token snippet, and a short suggestion. Block the commit by default.

Ignore Mechanism:

.secretignore (glob-based paths, like .gitignore).

Inline # secretlint:allow comment (optional in v0.2).

2. System Architecture

CLI Tool (secretlint)

Input: staged diff (git diff --cached -U0).

Diff Parser: extract only added lines (+ ...).

Scanner:

Run regex rules line-by-line.

Run entropy check if no regex hit but line “looks random”.

Reporter:

Pretty-print results (with red ❌ and yellow ⚠️).

Exit 1 if hits found.

Config loader:

YAML or JSON config file .secretlintrc.

.secretignore for file globs.

Pre-commit Hook

Thin wrapper calling secretlint scan --staged.

If exit code != 0 → abort commit.

3. User Flow
Install
brew install secretlint-lite
# or
go install github.com/yourname/secretlint-lite@latest

Setup
secretlint init


Creates:

.secretlintrc.yml (default rules)

.secretignore (ignore node_modules, dist, etc.)

.git/hooks/pre-commit (auto-added, calls secretlint scan --staged)

Everyday Use
git add src/config.ts
git commit -m "oops add api key"


Output:

⛔ Secret detected in staged changes

Rule     : OPENAI_API_KEY
File     : src/config.ts:12
Snippet  : sk-abc123****************
Advice   : Move this to an env var (.env file) and add .env to .gitignore.

Commit aborted.
Sensible Roadmap

v0.1: staged-only, curated rules, .secretignore, baseline snapshot.

v0.2: inline # secretlint:allow RULE_ID ttl=7d, --explain RULE_ID.

v0.3: lightweight editor bridge (optional), --fix helpers (move value to .env, replace with process.env).

v0.4: pluggable rules (YAML/JSON), org presets without a server.

v1.0: pre-commit framework support, Windows/macOS/Linux builds, simple telemetry opt-in.
