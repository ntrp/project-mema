package zonaq

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "ZonaQ",
		DefaultBaseURL: "https://www.zonaq.pw/",
		SearchPath:     "/retorno/2/index.php",
		QueryParam:     "search",
		LoginPath:      "/paDentro.php",
		UsernameParam:  "user",
		PasswordParam:  "passwrd",
	}, clients...)
}
