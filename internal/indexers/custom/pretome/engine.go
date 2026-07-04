package pretome

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "PreToMe",
		DefaultBaseURL: "https://pretome.info/",
		SearchPath:     "/browse.php",
		QueryParam:     "search",
		LoginPath:      "/takelogin.php",
		ExtraLogin:     map[string]string{"login_pin": ""},
	}, clients...)
}
