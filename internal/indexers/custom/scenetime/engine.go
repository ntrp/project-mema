package scenetime

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "SceneTime",
		DefaultBaseURL: "https://www.scenetime.com/",
		SearchPath:     "/browse.php",
		QueryParam:     "search",
	}, clients...)
}
