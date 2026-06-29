package httpapi

import (
	"time"

	"media-manager/internal/config"
	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/jobs"
	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

type Server struct {
	cfg             config.Config
	settings        *storage.SettingsStore
	downloadClients *downloadclients.Service
	indexers        *indexers.Service
	metadata        *metadata.Service
	jobs            *jobs.Client
	sessions        *sessionStore
	now             func() time.Time
}

func NewServer(cfg config.Config, settings *storage.SettingsStore, downloadClients *downloadclients.Service, indexerService *indexers.Service, metadataService *metadata.Service, jobs *jobs.Client) *Server {
	return &Server{
		cfg:             cfg,
		settings:        settings,
		downloadClients: downloadClients,
		indexers:        indexerService,
		metadata:        metadataService,
		jobs:            jobs,
		sessions:        newSessionStore(),
		now:             time.Now,
	}
}
