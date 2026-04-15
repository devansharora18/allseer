package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"allseer/internal/rules"
)

func TestOutboundURLAbsolute(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/resource",
		},
		Host: "proxy.local",
	}

	out, err := outboundURL(req)
	if err != nil {
		t.Fatalf("outboundURL returned error: %v", err)
	}

	if got, want := out.String(), "http://example.com/resource"; got != want {
		t.Fatalf("unexpected URL: got %q want %q", got, want)
	}
}

func TestOutboundURLRelativeFallsBackToHost(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		URL: &url.URL{
			Path: "/resource",
		},
		Host: "example.com",
	}

	out, err := outboundURL(req)
	if err != nil {
		t.Fatalf("outboundURL returned error: %v", err)
	}

	if got, want := out.String(), "http://example.com/resource"; got != want {
		t.Fatalf("unexpected URL: got %q want %q", got, want)
	}
}

func TestStripHopByHopHeaders(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Set("Connection", "Keep-Alive, Upgrade")
	headers.Set("Keep-Alive", "timeout=5")
	headers.Set("Upgrade", "websocket")
	headers.Set("Transfer-Encoding", "chunked")
	headers.Set("X-Custom", "ok")

	stripHopByHopHeaders(headers)

	if got := headers.Get("Connection"); got != "" {
		t.Fatalf("expected Connection header to be removed, got %q", got)
	}

	if got := headers.Get("Keep-Alive"); got != "" {
		t.Fatalf("expected Keep-Alive header to be removed, got %q", got)
	}

	if got := headers.Get("Upgrade"); got != "" {
		t.Fatalf("expected Upgrade header to be removed, got %q", got)
	}

	if got := headers.Get("Transfer-Encoding"); got != "" {
		t.Fatalf("expected Transfer-Encoding header to be removed, got %q", got)
	}

	if got := headers.Get("X-Custom"); got != "ok" {
		t.Fatalf("expected X-Custom header to be preserved, got %q", got)
	}
}

func TestAppendForwardedFor(t *testing.T) {
	t.Parallel()

	headers := http.Header{}
	headers.Set("X-Forwarded-For", "10.0.0.2")

	appendForwardedFor(headers, "10.0.0.3:1234")

	if got, want := headers.Get("X-Forwarded-For"), "10.0.0.2, 10.0.0.3"; got != want {
		t.Fatalf("unexpected X-Forwarded-For: got %q want %q", got, want)
	}
}

func TestHandleHTTPBlockRule(t *testing.T) {
	t.Parallel()

	upstreamCalled := false
	upstream := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		upstreamCalled = true
	}))
	t.Cleanup(upstream.Close)

	engine := rules.NewEngine([]rules.Rule{
		{
			Name:     "block-all",
			Enabled:  true,
			Priority: 1,
			Match: rules.Matcher{
				HostPattern: "*",
			},
			Action: rules.Action{Type: rules.ActionBlock},
		},
	})

	proxyServer := newTestProxyServer(t, engine)
	req := httptest.NewRequest(http.MethodGet, upstream.URL+"/resource", nil)
	rr := httptest.NewRecorder()

	proxyServer.handleHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	if upstreamCalled {
		t.Fatalf("expected blocked request to avoid upstream call")
	}
}

func TestHandleHTTPRedirectRule(t *testing.T) {
	t.Parallel()

	engine := rules.NewEngine([]rules.Rule{
		{
			Name:     "redirect-all",
			Enabled:  true,
			Priority: 1,
			Match: rules.Matcher{
				HostPattern: "*",
			},
			Action: rules.Action{
				Type:   rules.ActionRedirect,
				Target: "https://example.org/landing",
			},
		},
	})

	proxyServer := newTestProxyServer(t, engine)
	req := httptest.NewRequest(http.MethodGet, "http://service.local/home", nil)
	rr := httptest.NewRecorder()

	proxyServer.handleHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	if got, want := rr.Header().Get("Location"), "https://example.org/landing"; got != want {
		t.Fatalf("unexpected Location header: got %q want %q", got, want)
	}
}

func TestHandleHTTPModifyRequestHeaders(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(r.Header.Get("X-Proxy-Injected")))
	}))
	t.Cleanup(upstream.Close)

	engine := rules.NewEngine([]rules.Rule{
		{
			Name:     "inject-header",
			Enabled:  true,
			Priority: 1,
			Match: rules.Matcher{
				HostPattern: "*",
			},
			Action: rules.Action{
				Type: rules.ActionModifyRequest,
				Headers: map[string]string{
					"X-Proxy-Injected": "allseer",
				},
			},
		},
	})

	proxyServer := newTestProxyServer(t, engine)
	req := httptest.NewRequest(http.MethodGet, upstream.URL+"/resource", nil)
	rr := httptest.NewRecorder()

	proxyServer.handleHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	body, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if got, want := strings.TrimSpace(string(body)), "allseer"; got != want {
		t.Fatalf("expected upstream to receive injected header value %q, got %q", want, got)
	}
}
