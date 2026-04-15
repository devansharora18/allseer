package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"allseer/internal/config"
	"allseer/internal/proxy"
	"allseer/internal/rules"
)

type App struct {
	cfg    config.Config
	logger *slog.Logger
	server *proxy.Server
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	if logger == nil {
		logger = slog.Default()
	}

	engine := rules.NewEngine(nil)
	if cfg.Rules.File != "" {
		loadedRules, err := rules.LoadFromFile(cfg.Rules.File)
		switch {
		case err == nil:
			engine.ReplaceRules(loadedRules)
			logger.Info("loaded rules", "count", len(loadedRules), "file", cfg.Rules.File)
		case errors.Is(err, os.ErrNotExist):
			logger.Warn("rules file not found, continuing without rules", "file", cfg.Rules.File)
		default:
			return nil, fmt.Errorf("load rules: %w", err)
		}
	}

	return &App{
		cfg:    cfg,
		logger: logger,
		server: proxy.NewServer(proxy.Options{
			ListenAddr:        cfg.Proxy.ListenAddr,
			ReadHeaderTimeout: cfg.Proxy.ReadHeaderTimeout.Value(),
			IdleTimeout:       cfg.Proxy.IdleTimeout.Value(),
			RuleEngine:        engine,
		}, logger),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		a.logger.Info("starting proxy server", "addr", a.cfg.Proxy.ListenAddr)
		errCh <- a.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownTimeout := a.cfg.Proxy.ShutdownTimeout.Value()
		if shutdownTimeout <= 0 {
			shutdownTimeout = 10 * time.Second
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		a.logger.Info("shutting down proxy server")
		return a.server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
