package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mushibana/glacier/backend/internal/app"
	"github.com/mushibana/glacier/backend/internal/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := app.NewLogger(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application, err := app.New(cfg, logger)
	if err != nil {
		return err
	}

	defer func() {
		if err := application.Close(); err != nil {
			logger.Error("failed to close application", "error", err)
		}
	}()

	return application.Run(ctx)
}
