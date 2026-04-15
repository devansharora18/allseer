package proxy

import (
	"net/http"
	"net/url"
	"testing"
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
