package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"media-manager/internal/app"
	"media-manager/internal/logging"
)

func main() {
	logging.ConfigureDefault()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		exitWithError(err)
		return
	}
}

func exitWithError(err error) {
	slog.Error("server failed", "error", err)
	os.Exit(1)
}
