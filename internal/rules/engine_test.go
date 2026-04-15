package rules

import (
	"testing"
	"time"
)

func TestEngineHostPatternsMatchesSubdomain(t *testing.T) {
	t.Parallel()

	engine := NewEngine([]Rule{
		{
			Name:     "block-ads",
			Enabled:  true,
			Priority: 1,
			Match: Matcher{
				HostPatterns: []string{"doubleclick.net"},
			},
			Action: Action{Type: ActionBlock},
		},
	})

	decision := engine.Evaluate(RequestContext{
		Method: "GET",
		Host:   "ads.doubleclick.net",
		Path:   "/",
		Time:   time.Now(),
	})

	if !decision.Matched {
		t.Fatalf("expected ad subdomain to match")
	}
}

func TestEngineHostPatternsNoFalsePositive(t *testing.T) {
	t.Parallel()

	engine := NewEngine([]Rule{
		{
			Name:     "block-ads",
			Enabled:  true,
			Priority: 1,
			Match: Matcher{
				HostPatterns: []string{"doubleclick.net"},
			},
			Action: Action{Type: ActionBlock},
		},
	})

	decision := engine.Evaluate(RequestContext{
		Method: "GET",
		Host:   "notdoubleclick.net",
		Path:   "/",
		Time:   time.Now(),
	})

	if decision.Matched {
		t.Fatalf("expected boundary-safe host matching to avoid false positive")
	}
}
