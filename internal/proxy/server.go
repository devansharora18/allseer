package proxy

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"allseer/internal/rules"
	"allseer/internal/storage/sqlite"
)

type Options struct {
	ListenAddr        string
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	RuleEngine        *rules.Engine
	LogRepository     *sqlite.LogRepository
}

type Server struct {
	httpServer *http.Server
	transport  *http.Transport
	dialer     *net.Dialer
	ruleEngine *rules.Engine
	logRepo    *sqlite.LogRepository
	logger     *slog.Logger
	startTime  time.Time
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

	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy:                 nil,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	ruleEngine := opts.RuleEngine
	if ruleEngine == nil {
		ruleEngine = rules.NewEngine(nil)
	}

	s := &Server{
		logger:     logger,
		transport: transport,
		dialer:    dialer,
		ruleEngine: ruleEngine,
		logRepo:   opts.LogRepository,
		startTime: time.Now(),
	}

	mux := http.NewServeMux()
	// API endpoints
	mux.HandleFunc("/api/stats", s.handleStatistics)
	mux.HandleFunc("/api/rules", s.handleRulesList)
	mux.HandleFunc("/api/rules/delete", s.handleRuleDetail)
	// SPA and proxy
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
	if s.transport != nil {
		s.transport.CloseIdleConnections()
	}

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	if isLocalAdminLogsRequest(r) {
		s.handleAdminLogs(w, r)
		return
	}

	if r.Method == http.MethodConnect {
		s.handleConnect(w, r)
		return
	}

	s.handleHTTP(w, r)
}
