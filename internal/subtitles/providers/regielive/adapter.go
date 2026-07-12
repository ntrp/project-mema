package regielive

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("regielive", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "regielive", DefaultBaseURL: "https://api.regielive.ro", SearchPath: "/subtitles/search", TestPath: "/subtitles/search"}})
}
