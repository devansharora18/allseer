package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAdBlockDomains(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "ads.txt")

	content := []byte(`# comment
||doubleclick.net^
0.0.0.0 ads.example.com
127.0.0.1 tracker.example.org # hosts style
*.adnxs.com
@@||allow.example^
invalid-domain
https://ads.service.net/script.js
`)

	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("failed writing ad list: %v", err)
	}

	domains, err := LoadAdBlockDomains(path)
	if err != nil {
		t.Fatalf("LoadAdBlockDomains failed: %v", err)
	}

	if len(domains) != 4 {
		t.Fatalf("expected 4 parsed domains, got %d: %#v", len(domains), domains)
	}

	expected := map[string]bool{
		"doubleclick.net":     true,
		"ads.example.com":     true,
		"tracker.example.org": true,
		"adnxs.com":           true,
	}

	for _, domain := range domains {
		if !expected[domain] {
			t.Fatalf("unexpected domain parsed: %q", domain)
		}
	}
}

func TestBuildAdBlockRule(t *testing.T) {
	t.Parallel()

	rule := BuildAdBlockRule([]string{"doubleclick.net", "adnxs.com"})

	if rule.Name != "block-ad-domains" {
		t.Fatalf("unexpected rule name %q", rule.Name)
	}

	if rule.Action.Type != ActionBlock {
		t.Fatalf("expected block action, got %q", rule.Action.Type)
	}

	if len(rule.Match.HostPatterns) != 2 {
		t.Fatalf("expected 2 host patterns, got %d", len(rule.Match.HostPatterns))
	}
}
