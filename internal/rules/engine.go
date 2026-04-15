package rules

import (
	"path"
	"sort"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	mu    sync.RWMutex
	rules []Rule
}

func NewEngine(rules []Rule) *Engine {
	e := &Engine{}
	e.ReplaceRules(rules)
	return e
}

func (e *Engine) ReplaceRules(in []Rule) {
	rules := append([]Rule(nil), in...)

	sort.SliceStable(rules, func(i, j int) bool {
		if rules[i].Priority == rules[j].Priority {
			return rules[i].Name < rules[j].Name
		}

		return rules[i].Priority < rules[j].Priority
	})

	e.mu.Lock()
	defer e.mu.Unlock()

	e.rules = rules
}

func (e *Engine) Evaluate(ctx RequestContext) Decision {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}

		if !matchMethod(rule.Match.Methods, ctx.Method) {
			continue
		}

		if !matchHost(rule.Match.HostPattern, rule.Match.HostPatterns, ctx.Host) {
			continue
		}

		if !matchPattern(rule.Match.PathPattern, ctx.Path) {
			continue
		}

		if !matchTimeWindow(rule.Match.TimeWindow, ctx.Time) {
			continue
		}

		return Decision{
			Matched: true,
			Rule:    rule,
			Action:  rule.Action,
		}
	}

	return Decision{
		Matched: false,
		Action: Action{
			Type: ActionAllow,
		},
	}
}

func matchHost(hostPattern string, hostPatterns []string, host string) bool {
	if hostPattern == "" && len(hostPatterns) == 0 {
		return true
	}

	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" {
		return false
	}

	if hostPattern != "" && matchSingleHostPattern(strings.ToLower(hostPattern), host) {
		return true
	}

	for _, pattern := range hostPatterns {
		if matchSingleHostPattern(strings.ToLower(pattern), host) {
			return true
		}
	}

	return false
}

func matchSingleHostPattern(pattern, host string) bool {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return false
	}

	if strings.ContainsAny(pattern, "*?[") {
		ok, err := path.Match(pattern, host)
		if err != nil {
			return false
		}
		return ok
	}

	return host == pattern || strings.HasSuffix(host, "."+pattern)
}

func matchMethod(methods []string, reqMethod string) bool {
	if len(methods) == 0 {
		return true
	}

	for _, method := range methods {
		if method == "*" || strings.EqualFold(method, reqMethod) {
			return true
		}
	}

	return false
}

func matchPattern(pattern, value string) bool {
	if pattern == "" {
		return true
	}

	ok, err := path.Match(pattern, value)
	if err != nil {
		return false
	}

	return ok
}

func matchTimeWindow(window *TimeWindow, now time.Time) bool {
	if window == nil {
		return true
	}

	if len(window.Weekdays) > 0 {
		match := false
		for _, day := range window.Weekdays {
			if day == now.Weekday() {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	hour := now.Hour()

	if window.StartHour == window.EndHour {
		return true
	}

	if window.StartHour < window.EndHour {
		return hour >= window.StartHour && hour < window.EndHour
	}

	return hour >= window.StartHour || hour < window.EndHour
}
