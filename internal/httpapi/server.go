package httpapi

import (
	"net/http"
	"time"

	"media-manager/internal/config"
	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

type Server struct {
	cfg             config.Config
	settings        *storage.SettingsStore
	downloadClients *downloadclients.Service
	indexers        *indexers.Service
	sessions        *sessionStore
	now             func() time.Time
}

func NewServer(cfg config.Config, settings *storage.SettingsStore) *Server {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &Server{
		cfg:             cfg,
		settings:        settings,
		downloadClients: downloadclients.NewService(httpClient),
		indexers:        indexers.NewService(httpClient),
		sessions:        newSessionStore(),
		now:             time.Now,
	}
}
