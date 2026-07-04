package exoticaz

import (
	"media-manager/internal/indexers/custom/avistazapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *avistazapi.Engine {
	return avistazapi.New(avistazapi.Options{
		Name:           "ExoticaZ",
		DefaultBaseURL: "https://exoticaz.to/",
	}, clients...)
}
