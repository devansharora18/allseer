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
	loadedRules := make([]rules.Rule, 0)

	if cfg.Rules.File != "" {
		fileRules, err := rules.LoadFromFile(cfg.Rules.File)
		switch {
		case err == nil:
			loadedRules = append(loadedRules, fileRules...)
			logger.Info("loaded rules", "count", len(fileRules), "file", cfg.Rules.File)
		case errors.Is(err, os.ErrNotExist):
			logger.Warn("rules file not found, continuing without rules", "file", cfg.Rules.File)
		default:
			return nil, fmt.Errorf("load rules: %w", err)
		}
	}

	if cfg.Rules.AdBlockFile != "" {
		adDomains, err := rules.LoadAdBlockDomains(cfg.Rules.AdBlockFile)
		switch {
		case err == nil:
			if len(adDomains) > 0 {
				adRule := rules.BuildAdBlockRule(adDomains)
				loadedRules = append(loadedRules, adRule)
				logger.Info("loaded ad block domains", "count", len(adDomains), "file", cfg.Rules.AdBlockFile)
			}
		case errors.Is(err, os.ErrNotExist):
			logger.Warn("ad block list file not found, continuing without ad presets", "file", cfg.Rules.AdBlockFile)
		default:
			return nil, fmt.Errorf("load ad block domains: %w", err)
		}
	}

	engine.ReplaceRules(loadedRules)

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
