package subsarr

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("subsarr", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "subsarr", DefaultBaseURL: "", SearchPath: "/api/subtitles/search", TestPath: "/api/health", Local: true}})
}
