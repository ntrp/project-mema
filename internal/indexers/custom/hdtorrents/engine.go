package hdtorrents

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "HDTorrents",
		DefaultBaseURL: "https://hdts.ru/",
		SearchPath:     "/torrents.php",
		QueryParam:     "search",
		LoginPath:      "/login.php",
		UsernameParam:  "uid",
		PasswordParam:  "pwd",
	}, clients...)
}
