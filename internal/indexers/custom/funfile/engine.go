package funfile

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "FunFile",
		DefaultBaseURL: "https://www.funfile.org/",
		SearchPath:     "/browse.php",
		QueryParam:     "search",
		LoginPath:      "/takelogin.php",
	}, clients...)
}
