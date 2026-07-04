package privatehd

import (
	"media-manager/internal/indexers/custom/avistazapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *avistazapi.Engine {
	return avistazapi.New(avistazapi.Options{
		Name:           "PrivateHD",
		DefaultBaseURL: "https://privatehd.to/",
	}, clients...)
}
