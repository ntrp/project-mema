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
	clientProfile := mediaPreviewClientProfile(params.ClientProfile, r.UserAgent())
	writeJSON(w, http.StatusOK, mediaPreviewInfo(target, params.AudioTrackIndex, clientProfile))
}

func mediaPreviewInfo(target string, audioTrackIndex *int32, clientProfile MediaFilePreviewClientProfile) MediaFilePreviewInfo {
	probe := mediaFileProbe(target)
	info := mediaPreviewInfoFromTracks(target, probe.tracks, audioTrackIndex, clientProfile)
	info.ContainerBitRate = probe.container.bitRate
	info.ContainerFormat = probe.container.format
	info.ContainerFormatName = probe.container.formatName
	info.DurationSeconds = probe.durationSeconds
	return info
}

func mediaPreviewInfoFromTracks(
	target string,
	tracks []MediaFileTrack,
	audioTrackIndex *int32,
	clientProfile MediaFilePreviewClientProfile,
) MediaFilePreviewInfo {
	video := firstTrackByType(tracks, Video, nil)
	audio := firstTrackByType(tracks, Audio, audioTrackIndex)
	decision := mediaPreviewDecisionFromTracks(target, tracks, audioTrackIndex, clientProfile)
	reasons := append([]string{}, decision.reasons...)
	sourceBitRate := selectedBitRate(video, audio)
	info := MediaFilePreviewInfo{
		StreamingMode:    decision.mode,
		DeliveryProtocol: decision.deliveryProtocol,
		OutputVideoCodec: decision.plan.videoCodec,
		OutputAudioCodec: decision.plan.audioCodec,
		SourceBitRate:    sourceBitRate,
		LiveBitRate:      sourceBitRate,
		TranscodeReasons: &reasons,
	}
	if video != nil {
		info.VideoTrack = video
	}
	if audio != nil {
		info.AudioTrack = audio
	}
	return info
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
