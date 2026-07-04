package toloka

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "Toloka",
		DefaultBaseURL: "https://toloka.to/",
		SearchPath:     "/tracker.php",
		QueryParam:     "nm",
		LoginPath:      "/login.php",
	}, clients...)
}
