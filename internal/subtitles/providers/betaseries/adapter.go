package betaseries

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("betaseries", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "betaseries", RequiredSecret: "token", SecretHeader: "X-BetaSeries-Key", DefaultBaseURL: "https://api.betaseries.com", SearchPath: "/subtitles/search", TestPath: "/members/infos", SeriesOnly: true}})
}
