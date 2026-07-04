package secretcinema

import (
	"media-manager/internal/indexers/custom/gazelleapi"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *gazelleapi.Engine {
	return gazelleapi.New(gazelleapi.Options{
		Name:             "SecretCinema",
		DefaultBaseURL:   "https://secret-cinema.pw/",
		PreferGroupTitle: true,
	}, clients...)
}
