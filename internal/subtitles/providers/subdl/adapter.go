package subdl

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("subdl", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "subdl", RequiredSecret: "apiKey", SecretQueryName: "api_key", DefaultBaseURL: "https://api.subdl.com", SearchPath: "/api/v1/subtitles", TestPath: "/api/v1/subtitles"}})
}
