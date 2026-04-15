package proxy

import (
	"io"
	"log/slog"
	"testing"

	"allseer/internal/rules"
)

func newTestProxyServer(t *testing.T, engine *rules.Engine) *Server {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewServer(Options{
		ListenAddr: "127.0.0.1:0",
		RuleEngine: engine,
	}, logger)
}
