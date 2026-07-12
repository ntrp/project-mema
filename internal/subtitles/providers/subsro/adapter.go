package subsro

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("subsro", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "subsro", RequiredSecret: "apiKey", SecretQueryName: "key", DefaultBaseURL: "https://api.subs.ro", SearchPath: "/subtitles", TestPath: "/subtitles", RequireIMDb: true}})
}
