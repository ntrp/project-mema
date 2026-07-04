package orpheus

import (
	"media-manager/internal/indexers/custom/gazelleapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *gazelleapi.Engine {
	return gazelleapi.New(gazelleapi.Options{
		Name:                    "Orpheus",
		DefaultBaseURL:          "https://orpheus.network/",
		AuthHeader:              "Authorization",
		AuthTokenPrefix:         "token ",
		PreferMusicTitle:        true,
		DownloadPath:            "/ajax.php",
		SupportsFreeleechTokens: true,
	}, clients...)
}
