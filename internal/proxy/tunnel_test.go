package proxy

import "testing"

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
