package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"allseer/internal/rules"
)

func TestNormalizeConnectAddressWithPort(t *testing.T) {
	t.Parallel()

	if got, want := normalizeConnectAddress("example.com:8443"), "example.com:8443"; got != want {
		t.Fatalf("unexpected normalized address: got %q want %q", got, want)
	}
}

func TestNormalizeConnectAddressWithoutPort(t *testing.T) {
	t.Parallel()

	if got, want := normalizeConnectAddress("example.com"), "example.com:443"; got != want {
		t.Fatalf("unexpected normalized address: got %q want %q", got, want)
	}
}

func TestHandleConnectBlockRule(t *testing.T) {
	t.Parallel()

	engine := rules.NewEngine([]rules.Rule{
		{
			Name:     "block-connect",
			Enabled:  true,
			Priority: 1,
			Match: rules.Matcher{
				HostPattern: "*",
				Methods:     []string{http.MethodConnect},
			},
			Action: rules.Action{Type: rules.ActionBlock},
		},
	})

	proxyServer := newTestProxyServer(t, engine)
	req := httptest.NewRequest(http.MethodConnect, "http://example.com:443", nil)
	req.Host = "example.com:443"
	rr := httptest.NewRecorder()

	proxyServer.handleConnect(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}
}
