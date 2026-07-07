package httpapi

import (
	"net/http"
	"strings"

	"media-manager/internal/delivery"
)

func serveMediaPreviewHLSPlaylist(
	w http.ResponseWriter,
	r *http.Request,
	filePath MediaFilePath,
	target string,
	probe delivery.ProbeResult,
	audioTrackIndex *int32,
	clientProfile MediaFilePreviewClientProfile,
	decision delivery.Decision,
) {
	if probe.DurationSeconds == nil || !delivery.ValidSeconds(*probe.DurationSeconds) || *probe.DurationSeconds <= 0 {
		writeError(w, http.StatusInternalServerError, "media_preview_duration_unavailable", "Could not determine media duration for browser preview")
		return
	}
	segments := delivery.HLSSegmentsForDecision(target, *probe.DurationSeconds, decision)
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write([]byte(delivery.HLSPlaylistText(delivery.PlaylistRequest{
		Path:          r.URL.Path,
		FilePath:      string(filePath),
		AudioTrack:    audioTrackIndex,
		ClientProfile: deliveryClientProfile(clientProfile),
		Segments:      segments,
		SegmentPath:   strings.TrimSuffix(r.URL.Path, "/preview") + "/preview-segment",
	})))
}
