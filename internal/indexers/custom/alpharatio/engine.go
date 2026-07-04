package alpharatio

import (
	"media-manager/internal/indexers/custom/gazelleapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *gazelleapi.Engine {
	return gazelleapi.New(gazelleapi.Options{
		Name:           "AlphaRatio",
		DefaultBaseURL: "https://alpharatio.cc/",
		FreeleechParam: "freetorrent",
		ExcludeScene:   true,
	}, clients...)
}
