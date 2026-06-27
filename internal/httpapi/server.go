package httpapi

import (
	"time"

	"media-manager/internal/config"
	"media-manager/internal/storage"
)

type Server struct {
	cfg      config.Config
	settings *storage.SettingsStore
	sessions *sessionStore
	now      func() time.Time
}

func NewServer(cfg config.Config, settings *storage.SettingsStore) *Server {
	return &Server{
		cfg:      cfg,
		settings: settings,
		sessions: newSessionStore(),
		now:      time.Now,
	}
}
