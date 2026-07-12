package whisperai

import (
	"media-manager/internal/subtitles/providers"
	"media-manager/internal/subtitles/providers/clusterapi"
)

func init() {
	providers.Register("whisperai", clusterapi.Adapter{Spec: clusterapi.Spec{Key: "whisperai", DefaultBaseURL: "http://localhost:9000", SearchPath: "/transcribe", TestPath: "/health", Local: true, CommandName: "ffmpeg", CommandArgs: []string{"-version"}}})
}
