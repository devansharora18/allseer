package proxy

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Options struct {
	ListenAddr        string
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
}

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewServer(opts Options, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}

	if opts.ListenAddr == "" {
		opts.ListenAddr = "127.0.0.1:8080"
	}

	if opts.ReadHeaderTimeout <= 0 {
		opts.ReadHeaderTimeout = 10 * time.Second
	}

	if opts.IdleTimeout <= 0 {
		opts.IdleTimeout = 60 * time.Second
	}

	s := &Server{logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleProxyRequest)

	s.httpServer = &http.Server{
		Addr:              opts.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: opts.ReadHeaderTimeout,
		IdleTimeout:       opts.IdleTimeout,
	}

	return s
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		s.handleConnect(w, r)
		return
	}

	s.handleHTTP(w, r)
}
