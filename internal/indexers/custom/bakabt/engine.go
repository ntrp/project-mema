package bakabt

import (
	"media-manager/internal/indexers/custom/htmltable"
	"media-manager/internal/indexers/engine"
)

func New(clients ...engine.HTTPDoer) *htmltable.Engine {
	return htmltable.New(htmltable.Options{
		Name:           "BakaBT",
		DefaultBaseURL: "https://bakabt.me/",
		SearchPath:     "/browse.php",
		QueryParam:     "q",
		LoginPath:      "/login.php",
		ExtraParams: map[string]string{
			"only": "0", "incomplete": "1", "lossless": "1", "hd": "1", "multiaudio": "1", "bonus": "1", "reorder": "1",
		},
	}, clients...)
}
