package subtis

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("subtis", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "subtis", DefaultBaseURL: "https://api.subt.is", SearchPath: "/subtitles/search", TestPath: "/subtitles/search", MovieOnly: true}})
}
