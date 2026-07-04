package animez

import (
	"media-manager/internal/indexers/custom/avistazapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *avistazapi.Engine {
	return avistazapi.New(avistazapi.Options{
		Name:            "AnimeZ",
		DefaultBaseURL:  "https://animez.to/",
		AnimeCategories: true,
		PreferRelease:   true,
	}, clients...)
}
