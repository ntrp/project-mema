package gestdown

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("gestdown", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "gestdown", DefaultBaseURL: "https://api.gestdown.info", SearchPath: "/subtitles/search", TestPath: "/status", SeriesOnly: true}})
}
