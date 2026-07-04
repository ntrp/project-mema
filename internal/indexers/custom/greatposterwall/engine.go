package greatposterwall

import (
	"media-manager/internal/indexers/custom/gazelleapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *gazelleapi.Engine {
	return gazelleapi.New(gazelleapi.Options{
		Name:           "GreatPosterWall",
		DefaultBaseURL: "https://greatposterwall.com/",
		FreeleechParam: "freetorrent",
	}, clients...)
}
