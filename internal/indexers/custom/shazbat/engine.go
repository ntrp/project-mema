package shazbat

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "Shazbat",
		DefaultBaseURL: "https://www.shazbat.tube/",
		SearchPath:     "/search",
		QueryParam:     "search",
		LoginPath:      "/login",
		ExtraParams:    map[string]string{"portlet": "true"},
	}, clients...)
}
