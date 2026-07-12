package subx

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("subx", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "subx", RequiredSecret: "apiKey", SecretHeader: "X-API-Key", DefaultBaseURL: "https://subx-api.duckdns.org", SearchPath: "/search", TestPath: "/status"}})
}
