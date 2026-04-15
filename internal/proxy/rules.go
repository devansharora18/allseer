package proxy

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"allseer/internal/rules"
)

func (s *Server) evaluateRequestDecision(r *http.Request) rules.Decision {
	host := requestHost(r)
	path := "/"
	if r.URL != nil && r.URL.Path != "" {
		path = r.URL.Path
	}

	ctx := rules.RequestContext{
		Method: r.Method,
		Host:   hostWithoutPort(host),
		Path:   path,
		Time:   time.Now(),
	}

	return s.ruleEngine.Evaluate(ctx)
}

func requestHost(r *http.Request) string {
	if r == nil {
		return ""
	}

	if r.URL != nil && r.URL.Host != "" {
		return r.URL.Host
	}

	return r.Host
}

func hostWithoutPort(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}

	if strings.Contains(host, "://") {
		if parsed, err := url.Parse(host); err == nil {
			host = parsed.Host
		}
	}

	if stripped, _, err := net.SplitHostPort(host); err == nil {
		return stripped
	}

	if strings.HasPrefix(host, "[") && strings.Contains(host, "]") {
		host = strings.TrimPrefix(host, "[")
		host = strings.SplitN(host, "]", 2)[0]
	}

	host = strings.TrimSuffix(host, ".")
	return strings.ToLower(host)
}

func shouldBlockDecision(decision rules.Decision) bool {
	return decision.Matched && decision.Action.Type == rules.ActionBlock
}

func shouldRedirectHTTP(decision rules.Decision) bool {
	return decision.Matched && decision.Action.Type == rules.ActionRedirect
}
