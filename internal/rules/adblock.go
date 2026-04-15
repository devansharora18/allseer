package rules

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const adBlockRuleName = "block-ad-domains"

func LoadAdBlockDomains(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	seen := make(map[string]struct{})
	domains := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		domain := normalizeAdBlockDomain(line)
		if domain == "" {
			continue
		}

		if _, exists := seen[domain]; exists {
			continue
		}

		seen[domain] = struct{}{}
		domains = append(domains, domain)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan ad block domains file %q: %w", path, err)
	}

	return domains, nil
}

func BuildAdBlockRule(domains []string) Rule {
	hostPatterns := make([]string, 0, len(domains))
	for _, domain := range domains {
		domain = strings.TrimSpace(strings.ToLower(domain))
		if domain == "" {
			continue
		}
		hostPatterns = append(hostPatterns, domain)
	}

	return Rule{
		Name:     adBlockRuleName,
		Enabled:  true,
		Priority: 5,
		Match: Matcher{
			HostPatterns: hostPatterns,
		},
		Action: Action{Type: ActionBlock},
	}
}

func normalizeAdBlockDomain(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}

	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "@@") {
		return ""
	}

	if hash := strings.Index(line, "#"); hash >= 0 {
		line = strings.TrimSpace(line[:hash])
	}

	if line == "" {
		return ""
	}

	fields := strings.Fields(line)
	if len(fields) >= 2 && looksLikeIP(fields[0]) {
		line = fields[1]
	}

	line = strings.TrimPrefix(line, "||")
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimPrefix(line, "*.")
	line = strings.TrimPrefix(line, ".")

	if cut := strings.IndexAny(line, "^/$?"); cut >= 0 {
		line = line[:cut]
	}

	line = strings.TrimSpace(strings.ToLower(line))
	line = strings.TrimSuffix(line, ".")

	if host, _, err := net.SplitHostPort(line); err == nil {
		line = host
	}

	if !isDomainLike(line) {
		return ""
	}

	return line
}

func looksLikeIP(value string) bool {
	return net.ParseIP(strings.TrimSpace(value)) != nil
}

func isDomainLike(value string) bool {
	if value == "" || strings.Contains(value, " ") {
		return false
	}

	if net.ParseIP(value) != nil {
		return false
	}

	if !strings.Contains(value, ".") {
		return false
	}

	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.' {
			continue
		}
		return false
	}

	return true
}
