package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"media-manager/internal/app"
	"media-manager/internal/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, os.Args[1:]); err != nil {
		exitWithError(err)
		return
	}
}

func exitWithError(err error) {
	if errors.Is(err, storage.ErrDevResetNotAllowed) {
		slog.Error("refusing development reset", "error", err)
		os.Exit(2)
	}
	slog.Error("server failed", "error", err)
	os.Exit(1)
}
