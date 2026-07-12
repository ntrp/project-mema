package assrt

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("assrt", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "assrt", RequiredSecret: "token", SecretQueryName: "token", DefaultBaseURL: "https://api.assrt.net", SearchPath: "/v1/sub/search", TestPath: "/v1/sub/search"}})
}
