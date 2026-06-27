package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/config"
	"media-manager/internal/httpapi"
	"media-manager/internal/storage"
	"media-manager/internal/web"
)

func Run(ctx context.Context, args []string) error {
	cfg := config.Load()
	if len(args) > 0 && args[0] == "reset-dev" {
		return resetDevelopment(ctx, cfg)
	}

	pool, err := openDatabase(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	server := newHTTPServer(cfg, pool)
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server listening", "addr", cfg.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return shutdownServer(server)
	case err := <-errCh:
		return err
	}
}

func resetDevelopment(ctx context.Context, cfg config.Config) error {
	if err := storage.ResetDevelopment(ctx, cfg); err != nil {
		return fmt.Errorf("development reset failed: %w", err)
	}
	slog.Info("development database reset complete")
	return nil
}

func openDatabase(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database connection setup failed: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	return pool, nil
}

func newHTTPServer(cfg config.Config, pool *pgxpool.Pool) *http.Server {
	apiRouter := chi.NewRouter()
	settingsStore := storage.NewSettingsStore(pool)
	httpapi.HandlerFromMux(httpapi.NewServer(cfg, settingsStore), apiRouter)

	router := chi.NewRouter()
	router.Mount("/api", apiRouter)
	router.Handle("/*", web.StaticHandler(cfg.WebDir))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func shutdownServer(server *http.Server) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	return nil
}
