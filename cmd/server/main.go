package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/config"
	"media-manager/internal/httpapi"
	"media-manager/internal/storage"
	"media-manager/internal/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	if len(os.Args) > 1 && os.Args[1] == "reset-dev" {
		if err := storage.ResetDevelopment(ctx, cfg); err != nil {
			if errors.Is(err, storage.ErrDevResetNotAllowed) {
				slog.Error("refusing development reset", "error", err)
				os.Exit(2)
			}
			slog.Error("development reset failed", "error", err)
			os.Exit(1)
		}
		slog.Info("development database reset complete")
		return
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection setup failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		slog.Error("database ping failed", "error", err)
		os.Exit(1)
	}

	apiRouter := chi.NewRouter()
	httpapi.HandlerFromMux(httpapi.NewServer(cfg), apiRouter)

	router := chi.NewRouter()
	router.Mount("/api", apiRouter)
	router.Handle("/*", web.StaticHandler(cfg.WebDir))

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("server listening", "addr", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "server shutdown failed: %v\n", err)
		os.Exit(1)
	}
}
