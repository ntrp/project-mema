package speedcd

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "SpeedCD",
		DefaultBaseURL: "https://speed.cd/",
		SearchPath:     "/browse.php",
		QueryParam:     "search",
		LoginPath:      "/checkpoint/",
		PasswordParam:  "pwd",
	}, clients...)
}
