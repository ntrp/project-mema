package subsource

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("subsource", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "subsource", RequiredSecret: "apiKey", SecretHeader: "Authorization", DefaultBaseURL: "https://subsource.net", SearchPath: "/api/search", TestPath: "/api/search"}})
}
