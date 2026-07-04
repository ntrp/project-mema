package animedia

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "Animedia",
		DefaultBaseURL: "https://tt.animedia.tv/",
		SearchPath:     "/ajax/search_result/P0",
		QueryParam:     "search",
	}, clients...)
}
