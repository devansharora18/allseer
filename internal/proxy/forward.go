package proxy

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"allseer/internal/rules"
)

var hopByHopHeaders = []string{
	"Connection",
	"Proxy-Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	decision := s.evaluateRequestDecision(r)

	if shouldBlockDecision(decision) {
		s.logger.Info("http request blocked by rule", "rule", decision.Rule.Name, "method", r.Method, "host", requestHost(r))
		http.Error(w, "blocked by proxy rule", http.StatusForbidden)
		return
	}

	if shouldRedirectHTTP(decision) {
		target := strings.TrimSpace(decision.Action.Target)
		if target == "" {
			s.logger.Warn("redirect rule missing target", "rule", decision.Rule.Name)
			http.Error(w, "rule redirect target is missing", http.StatusInternalServerError)
			return
		}

		s.logger.Info("http request redirected by rule", "rule", decision.Rule.Name, "method", r.Method, "from", requestHost(r), "to", target)
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		return
	}

	targetURL, err := outboundURL(r)
	if err != nil {
		s.logger.Warn("invalid outbound URL", "error", err)
		http.Error(w, "invalid target URL", http.StatusBadRequest)
		return
	}

	outReq := r.Clone(r.Context())
	outReq.URL = targetURL
	outReq.RequestURI = ""
	outReq.Host = targetURL.Host

	stripHopByHopHeaders(outReq.Header)
	appendForwardedFor(outReq.Header, r.RemoteAddr)
	if decision.Matched && decision.Action.Type == rules.ActionModifyRequest {
		applyHeaderMutations(outReq.Header, decision.Action.Headers)
	}

	resp, err := s.transport.RoundTrip(outReq)
	if err != nil {
		status := http.StatusBadGateway
		message := "upstream request failed"

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			status = http.StatusGatewayTimeout
			message = "upstream request timed out"
		}

		s.logger.Error("failed to reach upstream", "method", r.Method, "target", targetURL.String(), "error", err)
		http.Error(w, message, status)
		return
	}
	defer resp.Body.Close()

	stripHopByHopHeaders(resp.Header)
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		s.logger.Warn("failed copying upstream response body", "target", targetURL.String(), "error", err)
	}

	s.logger.Info(
		"http request forwarded",
		"method", r.Method,
		"target", targetURL.String(),
		"status", resp.StatusCode,
		"rule", decision.Rule.Name,
		"rule_action", decision.Action.Type,
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func outboundURL(r *http.Request) (*url.URL, error) {
	if r.URL == nil {
		return nil, errors.New("missing request URL")
	}

	targetURL := new(url.URL)
	*targetURL = *r.URL

	if targetURL.Scheme == "" {
		targetURL.Scheme = "http"
	}

	if targetURL.Host == "" {
		targetURL.Host = r.Host
	}

	if targetURL.Host == "" {
		return nil, errors.New("missing request host")
	}

	return targetURL, nil
}

func stripHopByHopHeaders(headers http.Header) {
	if headers == nil {
		return
	}

	if connection := headers.Get("Connection"); connection != "" {
		for _, token := range strings.Split(connection, ",") {
			trimmed := strings.TrimSpace(token)
			if trimmed != "" {
				headers.Del(trimmed)
			}
		}
	}

	for _, key := range hopByHopHeaders {
		headers.Del(key)
	}
}

func appendForwardedFor(headers http.Header, remoteAddr string) {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}

	host = strings.TrimSpace(host)
	if host == "" {
		return
	}

	current := headers.Get("X-Forwarded-For")
	if current == "" {
		headers.Set("X-Forwarded-For", host)
		return
	}

	headers.Set("X-Forwarded-For", current+", "+host)
}

func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func applyHeaderMutations(headers http.Header, updates map[string]string) {
	for key, value := range updates {
		canonical := http.CanonicalHeaderKey(strings.TrimSpace(key))
		if canonical == "" {
			continue
		}

		if value == "" {
			headers.Del(canonical)
			continue
		}

		headers.Set(canonical, value)
	}
}
