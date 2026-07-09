package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/config"
	"media-manager/internal/dlna"
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

func Run(ctx context.Context) error {
	cfg := config.Load()
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

	server, jobClient, dlnaManager, err := newHTTPServer(cfg, pool)
	if err != nil {
		return err
	}
	if err := jobClient.Start(ctx); err != nil {
		return fmt.Errorf("job client start failed: %w", err)
	}
	if err := dlnaManager.Start(ctx); err != nil {
		return fmt.Errorf("dlna start failed: %w", err)
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
		return shutdown(server, jobClient, dlnaManager)
	case err := <-errCh:
		_ = dlnaManager.Stop(context.Background())
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

func newHTTPServer(cfg config.Config, pool *pgxpool.Pool) (*http.Server, *jobs.Client, *dlna.Manager, error) {
	apiRouter := chi.NewRouter()
	apiRouter.Use(middleware.Recoverer)
	settingsStore := storage.NewSettingsStore(pool)
	httpClient := &http.Client{Timeout: 10 * time.Second}
	downloadClientService := downloadclients.NewService(httpClient)
	indexerService := indexers.NewService(httpClient)
	metadataService := metadata.NewService(httpClient, settingsStore)
	eventBroker := events.NewBroker()
	jobClient, err := jobs.NewClient(pool, settingsStore, indexerService, downloadClientService, eventBroker)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("job client setup failed: %w", err)
	}
	dlnaManager := dlna.NewManager(settingsStore, "http://"+cfg.Addr)
	apiServer := httpapi.NewServer(cfg, settingsStore, downloadClientService, indexerService, metadataService, jobClient, eventBroker)
	apiServer.SetDLNAManager(dlnaManager)
	apiRouter.Get("/docs", httpapi.SwaggerUIHandler)
	apiRouter.Get("/openapi.yaml", httpapi.OpenAPISpecHandler)
	httpapi.HandlerFromMux(apiServer, apiRouter)

	handler := routeDLNA(appRouter(cfg, apiRouter), dlnaManager.Handler())

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}, jobClient, dlnaManager, nil
}

func appRouter(cfg config.Config, apiRouter http.Handler) http.Handler {
	router := chi.NewRouter()
	router.Mount("/api", apiRouter)
	if !cfg.IsDevelopment() {
		router.Handle("/*", web.StaticHandler(cfg.WebDir))
	}
	return router
}

func routeDLNA(next http.Handler, dlnaHandler http.Handler) http.Handler {
	stripped := http.StripPrefix("/dlna", dlnaHandler)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/dlna" || strings.HasPrefix(r.URL.Path, "/dlna/") {
			stripped.ServeHTTP(w, r)
			return
		}
		if isRootDLNAPath(r.URL.Path) {
			dlnaHandler.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isRootDLNAPath(path string) bool {
	switch path {
	case "/rootDesc.xml", "/contentDirectory.xml", "/connectionManager.xml", "/mediaReceiverRegistrar.xml",
		"/icon-256.png", "/icon-128.png", "/icon-120.png", "/icon-48.png":
		return true
	default:
		return strings.HasPrefix(path, "/control/") ||
			strings.HasPrefix(path, "/resource/") ||
			strings.HasPrefix(path, "/artwork/") ||
			strings.HasPrefix(path, "/subtitle/") ||
			strings.HasPrefix(path, "/events/")
	}
}

func shutdown(server *http.Server, jobClient *jobs.Client, dlnaManager *dlna.Manager) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if dlnaManager != nil {
		if err := dlnaManager.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("dlna shutdown failed: %w", err)
		}
	}
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
