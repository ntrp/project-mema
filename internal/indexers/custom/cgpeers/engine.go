package cgpeers

import (
	"media-manager/internal/indexers/custom/gazelleapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *gazelleapi.Engine {
	return gazelleapi.New(gazelleapi.Options{
		Name:           "CGPeers",
		DefaultBaseURL: "https://cgpeers.to/",
	}, clients...)
}
