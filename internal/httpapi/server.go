package httpapi

import (
	"time"

	"media-manager/internal/config"
	"media-manager/internal/dlna"
	"media-manager/internal/downloadclients"
	"media-manager/internal/events"
	"media-manager/internal/indexers"
	"media-manager/internal/jobs"
	"media-manager/internal/metadata"
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

type Server struct {
	cfg             config.Config
	settings        *storage.SettingsStore
	downloadClients *downloadclients.Service
	indexers        *indexers.Service
	metadata        *metadata.Service
	subtitles       *subtitles.Service
	jobs            *jobs.Client
	events          *events.Broker
	dlna            *dlna.Manager
	streamSecret    []byte
	now             func() time.Time
}

func NewServer(cfg config.Config, settings *storage.SettingsStore, downloadClients *downloadclients.Service, indexerService *indexers.Service, metadataService *metadata.Service, jobs *jobs.Client, eventBroker *events.Broker) *Server {
	if eventBroker == nil {
		eventBroker = events.NewBroker()
	}
	return &Server{
		cfg:             cfg,
		settings:        settings,
		downloadClients: downloadClients,
		indexers:        indexerService,
		metadata:        metadataService,
		subtitles:       subtitles.NewService(nil),
		jobs:            jobs,
		events:          eventBroker,
		streamSecret:    newStreamTokenSecret(),
		now:             time.Now,
	}
}

func (s *Server) SetDLNAManager(dlnaManager *dlna.Manager) {
	s.dlna = dlnaManager
}
