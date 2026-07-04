package uniotaku

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "Uniotaku",
		DefaultBaseURL: "https://tracker.uniotaku.com/",
		SearchPath:     "/torrents_.php",
		QueryParam:     "search",
		LoginPath:      "/account-login.php",
	}, clients...)
}
