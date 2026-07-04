package httpapi

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (s *Server) GetMediaItemFilePreviewInfo(w http.ResponseWriter, r *http.Request, id ResourceId, params GetMediaItemFilePreviewInfoParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	target, err := s.settings.MediaItemFilePath(r.Context(), uuid.UUID(id), params.Path)
	if err != nil {
		writeSettingsError(w, err, "Could not find media file")
		return
	}
	if _, ok := statMediaFile(w, target); !ok {
		return
	}
	writeJSON(w, http.StatusOK, mediaPreviewInfo(target, params.AudioTrackIndex))
}

func mediaPreviewInfo(target string, audioTrackIndex *int32) MediaFilePreviewInfo {
	probe := mediaFileProbe(target)
	info := mediaPreviewInfoFromTracks(target, probe.tracks, audioTrackIndex)
	info.DurationSeconds = probe.durationSeconds
	return info
}

func mediaPreviewInfoFromTracks(target string, tracks []MediaFileTrack, audioTrackIndex *int32) MediaFilePreviewInfo {
	video := firstTrackByType(tracks, Video, nil)
	audio := firstTrackByType(tracks, Audio, audioTrackIndex)
	plan := mediaPreviewPlanFromTracks(tracks, audioTrackIndex)
	sourceBitRate := selectedBitRate(video, audio)
	info := MediaFilePreviewInfo{
		StreamingMode:    mediaPreviewMode(target, tracks, audioTrackIndex, plan),
		OutputVideoCodec: plan.videoCodec,
		OutputAudioCodec: plan.audioCodec,
		SourceBitRate:    sourceBitRate,
		LiveBitRate:      sourceBitRate,
	}
	if video != nil {
		info.VideoTrack = video
	}
	if audio != nil {
		info.AudioTrack = audio
	}
	return info
}

func mediaPreviewMode(target string, tracks []MediaFileTrack, audioTrackIndex *int32, plan mediaPreviewTranscodePlan) MediaFilePreviewMode {
	if mediaPreviewDirectFromTracks(target, tracks, audioTrackIndex) {
		return Direct
	}
	if plan.videoCodec == "copy" && plan.audioCodec == "copy" {
		return Remux
	}
	return Transcode
}

func selectedBitRate(tracks ...*MediaFileTrack) *string {
	var total int64
	for _, track := range tracks {
		if track == nil || track.BitRate == nil {
			continue
		}
		value, err := strconv.ParseInt(*track.BitRate, 10, 64)
		if err != nil || value <= 0 {
			continue
		}
		total += value
	}
	if total <= 0 {
		return nil
	}
	formatted := strconv.FormatInt(total, 10)
	return &formatted
}
