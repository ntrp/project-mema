package newznab

import (
	"media-manager/internal/indexers/custom/nab"
	"media-manager/internal/indexers/engine"
)

func New(client engine.HTTPDoer) engine.Engine {
	return nab.New("Newznab", client)
}
