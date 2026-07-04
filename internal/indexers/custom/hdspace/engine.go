package hdspace

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "HDSpace",
		DefaultBaseURL: "https://hd-space.org/",
		SearchPath:     "/index.php",
		QueryParam:     "search",
		LoginPath:      "/index.php",
		UsernameParam:  "uid",
		PasswordParam:  "pwd",
		ExtraParams:    map[string]string{"page": "torrents", "options": "0"},
		ExtraLogin:     map[string]string{"page": "login"},
	}, clients...)
}
