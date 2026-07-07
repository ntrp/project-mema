package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	"media-manager/internal/delivery"
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
	probe := delivery.Probe(target)
	info := mediaPreviewInfoFromTracks(target, mediaFileTracksFromDelivery(probe.Tracks), audioTrackIndex, clientProfile)
	info.ContainerBitRate = probe.Container.BitRate
	info.ContainerFormat = probe.Container.Format
	info.ContainerFormatName = probe.Container.FormatName
	info.DurationSeconds = probe.DurationSeconds
	return info
}

func mediaPreviewInfoFromTracks(
	target string,
	tracks []MediaFileTrack,
	audioTrackIndex *int32,
	clientProfile MediaFilePreviewClientProfile,
) MediaFilePreviewInfo {
	deliveryTrackList := deliveryTracks(tracks)
	video := delivery.FirstTrackByType(deliveryTrackList, delivery.TrackVideo, nil)
	audio := delivery.FirstTrackByType(deliveryTrackList, delivery.TrackAudio, audioTrackIndex)
	decision := delivery.DecisionFromTracks(target, deliveryTrackList, audioTrackIndex, deliveryClientProfile(clientProfile))
	reasons := append([]string{}, decision.Reasons...)
	sourceBitRate := delivery.SelectedBitRate(video, audio)
	info := MediaFilePreviewInfo{
		StreamingMode:    mediaPreviewMode(decision.Mode),
		DeliveryProtocol: mediaPreviewProtocol(decision.DeliveryProtocol),
		OutputVideoCodec: decision.Plan.VideoCodec,
		OutputAudioCodec: decision.Plan.AudioCodec,
		SourceBitRate:    sourceBitRate,
		LiveBitRate:      sourceBitRate,
		TranscodeReasons: &reasons,
	}
	if video != nil {
		info.VideoTrack = mediaFileTrackFromDelivery(video)
	}
	if audio != nil {
		info.AudioTrack = mediaFileTrackFromDelivery(audio)
	}
	return info
}
