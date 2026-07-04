package redacted

import (
	"media-manager/internal/indexers/custom/gazelleapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *gazelleapi.Engine {
	return gazelleapi.New(gazelleapi.Options{
		Name:                    "Redacted",
		DefaultBaseURL:          "https://redacted.sh/",
		AuthHeader:              "Authorization",
		FreeleechParam:          "freetorrent",
		FreeleechValue:          "4",
		PreferMusicTitle:        true,
		DownloadPath:            "/ajax.php",
		SupportsFreeleechTokens: true,
	}, clients...)
}
