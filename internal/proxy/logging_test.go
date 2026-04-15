package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolvedRequestURLConnect(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodConnect, "http://example.com:443", nil)
	req.Host = "example.com:443"

	if got, want := resolvedRequestURL(req, ""), "https://example.com:443"; got != want {
		t.Fatalf("unexpected resolved CONNECT URL: got %q want %q", got, want)
	}
}

func TestResolvedRequestURLRelativeRequest(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://example.com/path?q=1", nil)
	req.URL.Scheme = ""
	req.URL.Host = ""
	req.Host = "example.com"

	if got, want := resolvedRequestURL(req, ""), "http://example.com/path?q=1"; got != want {
		t.Fatalf("unexpected resolved URL: got %q want %q", got, want)
	}
}
