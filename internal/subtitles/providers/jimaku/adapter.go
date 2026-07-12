package jimaku

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("jimaku", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "jimaku", RequiredSecret: "apiKey", SecretHeader: "Authorization", DefaultBaseURL: "https://jimaku.cc", SearchPath: "/api/entries/search", TestPath: "/api/entries/search"}})
}
