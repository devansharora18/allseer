package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"allseer/internal/app"
	"allseer/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	cfgPath := os.Getenv("ALLSEER_CONFIG")
	if cfgPath == "" {
		cfgPath = "config/config.json"
	}

	cfg, err := config.Load(cfgPath)
	switch {
	case err == nil:
	case errors.Is(err, os.ErrNotExist):
		defaults := config.Default()
		cfg = &defaults
		logger.Info("config file not found, using defaults", "path", cfgPath)
	default:
		logger.Error("failed to load config", "path", cfgPath, "error", err)
		os.Exit(1)
	}

	application := app.New(*cfg, logger)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}
