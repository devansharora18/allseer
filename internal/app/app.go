package app

import (
	"context"
	"log/slog"
	"time"

	"allseer/internal/config"
	"allseer/internal/proxy"
)

type App struct {
	cfg    config.Config
	logger *slog.Logger
	server *proxy.Server
}

func New(cfg config.Config, logger *slog.Logger) *App {
	if logger == nil {
		logger = slog.Default()
	}

	return &App{
		cfg:    cfg,
		logger: logger,
		server: proxy.NewServer(proxy.Options{
			ListenAddr:        cfg.Proxy.ListenAddr,
			ReadHeaderTimeout: cfg.Proxy.ReadHeaderTimeout.Value(),
			IdleTimeout:       cfg.Proxy.IdleTimeout.Value(),
		}, logger),
	}
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
