package pixelhd

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "PixelHD",
		DefaultBaseURL: "https://pixelhd.me/",
		SearchPath:     "/torrents.php",
		QueryParam:     "searchstr",
	}, clients...)
}
