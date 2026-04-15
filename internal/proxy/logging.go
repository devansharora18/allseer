package proxy

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"allseer/internal/rules"
	"allseer/internal/storage/sqlite"
)

func (s *Server) writeTrafficLog(r *http.Request, decision rules.Decision, statusCode int, started time.Time, targetURL string, trafficErr error) {
	if s == nil || s.logRepo == nil || r == nil {
		return
	}

	action := string(decision.Action.Type)
	if action == "" {
		action = string(rules.ActionAllow)
	}

	errText := ""
	if trafficErr != nil {
		errText = trafficErr.Error()
	}

	entry := sqlite.TrafficLog{
		Method:     r.Method,
		URL:        resolvedRequestURL(r, targetURL),
		Host:       hostWithoutPort(requestHost(r)),
		StatusCode: statusCode,
		Duration:   time.Since(started),
		Action:     action,
		RuleName:   decision.Rule.Name,
		Error:      errText,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.logRepo.Insert(entry); err != nil {
		s.logger.Warn("failed to persist traffic log", "method", r.Method, "url", entry.URL, "error", err)
	}
}

func resolvedRequestURL(r *http.Request, targetURL string) string {
	if trimmed := strings.TrimSpace(targetURL); trimmed != "" {
		return trimmed
	}

	if r.Method == http.MethodConnect {
		host := strings.TrimSpace(r.Host)
		if host == "" {
			host = strings.TrimSpace(requestHost(r))
		}
		if host != "" {
			return "https://" + host
		}
	}

	if r.URL == nil {
		return "unknown"
	}

	if r.URL.IsAbs() {
		return r.URL.String()
	}

	host := strings.TrimSpace(requestHost(r))
	if host == "" {
		if r.URL.Path == "" {
			return "unknown"
		}
		return r.URL.String()
	}

	clone := new(url.URL)
	*clone = *r.URL
	if clone.Scheme == "" {
		clone.Scheme = "http"
	}
	clone.Host = host

	return clone.String()
}
