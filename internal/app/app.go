package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/config"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/httpapi"
	"media-manager/internal/indexers"
	"media-manager/internal/jobs"
	"media-manager/internal/logging"
	"media-manager/internal/metadata"
	"media-manager/internal/storage"
	"media-manager/internal/web"
)

func Run(ctx context.Context, args []string) error {
	cfg := config.Load()
	if len(args) > 0 && args[0] == "reset-dev" {
		return resetDevelopment(ctx, cfg)
	}
	if err := ensureMediaDataDir(cfg.MediaDataDir); err != nil {
		return err
	}

	pool, err := openDatabase(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()
	settingsStore := storage.NewSettingsStore(pool)
	if err := configureFileLogging(ctx, settingsStore); err != nil {
		return err
	}

	server, jobClient, err := newHTTPServer(cfg, pool)
	if err != nil {
		return err
	}
	if err := jobClient.Start(ctx); err != nil {
		return fmt.Errorf("job client start failed: %w", err)
	}
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
		return shutdown(server, jobClient)
	case err := <-errCh:
		if stopErr := stopJobs(jobClient); stopErr != nil {
			return stopErr
		}
		return err
	}
}

func configureFileLogging(ctx context.Context, store *storage.SettingsStore) error {
	settings, err := store.GetLogFileSettings(ctx)
	if err != nil {
		return fmt.Errorf("log file settings load failed: %w", err)
	}
	if err := logging.Default.ConfigureFile(logging.FileSettings{
		Enabled:       settings.Enabled,
		Directory:     settings.Directory,
		RetentionDays: int(settings.RetentionDays),
	}); err != nil {
		return fmt.Errorf("log file setup failed: %w", err)
	}
	return nil
}

func ensureMediaDataDir(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("media data directory setup failed: %w", err)
	}
	return nil
}

func resetDevelopment(ctx context.Context, cfg config.Config) error {
	if err := storage.ResetDevelopment(ctx, cfg); err != nil {
		return fmt.Errorf("development reset failed: %w", err)
	}
	slog.Info("development database reset complete")
	return nil
}

func openDatabase(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("database connection setup failed: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	if err := storage.EnsureSchema(ctx, cfg.DatabaseURL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database schema setup failed: %w", err)
	}
	if err := storage.NewSettingsStore(pool).EnsureDefaultAdminUser(ctx, cfg.AdminUsername, cfg.AdminPassword); err != nil {
		pool.Close()
		return nil, fmt.Errorf("default admin user setup failed: %w", err)
	}
	return pool, nil
}

func newHTTPServer(cfg config.Config, pool *pgxpool.Pool) (*http.Server, *jobs.Client, error) {
	apiRouter := chi.NewRouter()
	settingsStore := storage.NewSettingsStore(pool)
	httpClient := &http.Client{Timeout: 10 * time.Second}
	downloadClientService := downloadclients.NewService(httpClient)
	indexerService := indexers.NewService(httpClient)
	metadataService := metadata.NewService(httpClient, settingsStore)
	eventBroker := events.NewBroker()
	jobClient, err := jobs.NewClient(pool, settingsStore, indexerService, downloadClientService, eventBroker)
	if err != nil {
		return nil, nil, fmt.Errorf("job client setup failed: %w", err)
	}
	httpapi.HandlerFromMux(httpapi.NewServer(cfg, settingsStore, downloadClientService, indexerService, metadataService, jobClient, eventBroker), apiRouter)

	router := chi.NewRouter()
	router.Mount("/api", apiRouter)
	router.Handle("/*", web.StaticHandler(cfg.WebDir))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}, jobClient, nil
}

func shutdown(server *http.Server, jobClient *jobs.Client) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	if err := jobClient.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("job client shutdown failed: %w", err)
	}
	return nil
}

func stopJobs(jobClient *jobs.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := jobClient.Stop(ctx); err != nil {
		return fmt.Errorf("job client shutdown failed: %w", err)
	}
	return nil
}
