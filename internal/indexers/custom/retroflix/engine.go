package retroflix

import (
	"media-manager/internal/indexers/custom/speedappapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *speedappapi.Engine {
	return speedappapi.New(speedappapi.Options{
		Name:           "RetroFlix",
		DefaultBaseURL: "https://retroflix.club/",
	}, clients...)
}
