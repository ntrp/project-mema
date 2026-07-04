package libble

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "Libble",
		DefaultBaseURL: "https://libble.me/",
		SearchPath:     "/torrents.php",
		QueryParam:     "searchstr",
		LoginPath:      "/login.php",
		ExtraLogin:     map[string]string{"code": ""},
	}, clients...)
}
