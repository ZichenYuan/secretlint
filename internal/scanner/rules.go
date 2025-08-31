package scanner

import (
	"fmt"
	"regexp"
	"strings"
)

// SecretRule represents a regex-based rule for detecting secrets
type SecretRule struct {
	ID          string
	Name        string
	Pattern     *regexp.Regexp
	Description string
	Advice      string
}

// Finding represents a detected secret
type Finding struct {
	RuleID      string
	RuleName    string
	FilePath    string
	LineNum     int
	Content     string
	Match       string
	StartPos    int
	EndPos      int
	Description string
	Advice      string
}

// SecretScanner handles secret detection using regex rules
type SecretScanner struct {
	rules []SecretRule
}

// NewSecretScanner creates a new SecretScanner with default rules
func NewSecretScanner() *SecretScanner {
	scanner := &SecretScanner{}
	scanner.loadDefaultRules()
	return scanner
}

// loadDefaultRules loads the curated regex patterns from the specification
func (s *SecretScanner) loadDefaultRules() {
	rules := []struct {
		id          string
		name        string
		pattern     string
		description string
		advice      string
	}{
		{
			id:          "OPENAI_API_KEY",
			name:        "OpenAI API Key",
			pattern:     `sk-[A-Za-z0-9]{20,}`,
			description: "OpenAI API key detected",
			advice:      "Move this to an environment variable (.env file) and add .env to .gitignore",
		},
		{
			id:          "GITHUB_PAT",
			name:        "GitHub Personal Access Token",
			pattern:     `ghp_[A-Za-z0-9]{36}`,
			description: "GitHub Personal Access Token detected",
			advice:      "Store in environment variables or GitHub Secrets for CI/CD",
		},
		{
			id:          "AWS_ACCESS_KEY",
			name:        "AWS Access Key ID",
			pattern:     `(AKIA|ASIA)[A-Z0-9]{16}`,
			description: "AWS Access Key ID detected",
			advice:      "Use AWS IAM roles or store in AWS credentials file/environment variables",
		},
		{
			id:          "AWS_SECRET_KEY",
			name:        "AWS Secret Access Key",
			pattern:     `(?i)aws(.{0,20})?(secret|access).{0,20}['\"][A-Za-z0-9/+=]{40}['\"]`,
			description: "AWS Secret Access Key detected",
			advice:      "Use AWS IAM roles or store in AWS credentials file/environment variables",
		},
		{
			id:          "STRIPE_LIVE_PK",
			name:        "Stripe Live Publishable Key",
			pattern:     `pk_live_[A-Za-z0-9]{24}`,
			description: "Stripe Live Publishable Key detected",
			advice:      "Move to environment variables and ensure it's not exposed in client-side code",
		},
		{
			id:          "STRIPE_LIVE_SK",
			name:        "Stripe Live Secret Key", 
			pattern:     `sk_live_[A-Za-z0-9]{24}`,
			description: "Stripe Live Secret Key detected",
			advice:      "Move to environment variables and never expose in client-side code",
		},
		{
			id:          "SLACK_TOKEN",
			name:        "Slack Token",
			pattern:     `xox[baprs]-[0-9A-Za-z\-]+`,
			description: "Slack API token detected",
			advice:      "Store in environment variables or secure configuration management",
		},
		{
			id:          "JWT_TOKEN",
			name:        "JSON Web Token",
			pattern:     `eyJ[A-Za-z0-9_\-]+\.[A-Za-z0-9._\-]+\.[A-Za-z0-9._\-]+`,
			description: "JWT token detected",
			advice:      "Avoid committing JWTs; use secure token storage and short expiration times",
		},
		{
			id:          "GENERIC_API_KEY",
			name:        "Generic API Key Pattern",
			pattern:     `(?i)(api[_\-]?key|apikey|secret[_\-]?key|secretkey|access[_\-]?token|accesstoken)\s*[=:]\s*['\"]?[A-Za-z0-9\+/]{32,}['\"]?`,
			description: "Generic API key pattern detected",
			advice:      "Move sensitive keys to environment variables or secure configuration",
		},
		{
			id:          "PRIVATE_KEY",
			name:        "Private Key",
			pattern:     `-----BEGIN\s+(RSA\s+)?PRIVATE\s+KEY-----`,
			description: "Private key detected",
			advice:      "Store private keys securely, never commit to version control",
		},
	}

	for _, rule := range rules {
		compiled, err := regexp.Compile(rule.pattern)
		if err != nil {
			// Skip invalid patterns
			continue
		}
		
		s.rules = append(s.rules, SecretRule{
			ID:          rule.id,
			Name:        rule.name,
			Pattern:     compiled,
			Description: rule.description,
			Advice:      rule.advice,
		})
	}
}

// ScanLine scans a single line for secrets using all loaded rules
func (s *SecretScanner) ScanLine(filePath string, lineNum int, content string) []Finding {
	var findings []Finding
	
	for _, rule := range s.rules {
		matches := rule.Pattern.FindAllStringSubmatch(content, -1)
		if matches == nil {
			continue
		}
		
		// Find all match positions
		indexes := rule.Pattern.FindAllStringIndex(content, -1)
		
		for i, match := range matches {
			var matchText string
			if len(match) > 0 {
				matchText = match[0]
			}
			
			// Get position info
			var startPos, endPos int
			if i < len(indexes) {
				startPos = indexes[i][0]
				endPos = indexes[i][1]
			}
			
			findings = append(findings, Finding{
				RuleID:      rule.ID,
				RuleName:    rule.Name,
				FilePath:    filePath,
				LineNum:     lineNum,
				Content:     content,
				Match:       matchText,
				StartPos:    startPos,
				EndPos:      endPos,
				Description: rule.Description,
				Advice:      rule.Advice,
			})
		}
	}
	
	return findings
}

// ScanLines scans multiple lines for secrets
func (s *SecretScanner) ScanLines(lines []DiffLine) []Finding {
	var allFindings []Finding
	
	for _, line := range lines {
		findings := s.ScanLine(line.FilePath, line.LineNum, line.Content)
		allFindings = append(allFindings, findings...)
	}
	
	return allFindings
}

// MaskSecret returns a masked version of the secret for safe display
func (f *Finding) MaskSecret() string {
	if len(f.Match) <= 8 {
		return strings.Repeat("*", len(f.Match))
	}
	
	// Show first 4 and last 4 characters, mask the middle
	prefix := f.Match[:4]
	suffix := f.Match[len(f.Match)-4:]
	middle := strings.Repeat("*", len(f.Match)-8)
	
	return fmt.Sprintf("%s%s%s", prefix, middle, suffix)
}