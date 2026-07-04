package rutracker

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "RuTracker",
		DefaultBaseURL: "https://rutracker.org/",
		SearchPath:     "/forum/tracker.php",
		QueryParam:     "nm",
		LoginPath:      "/forum/login.php",
		UsernameParam:  "login_username",
		PasswordParam:  "login_password",
	}, clients...)
}
