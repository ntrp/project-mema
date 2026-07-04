package anidub

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "Anidub",
		DefaultBaseURL: "https://tr.anidub.com/",
		SearchPath:     "/index.php",
		QueryParam:     "story",
		LoginPath:      "/index.php",
		UsernameParam:  "login_name",
		PasswordParam:  "login_password",
		ExtraParams:    map[string]string{"do": "search", "subaction": "search"},
		ExtraLogin:     map[string]string{"login": "submit"},
	}, clients...)
}
