package httpapi

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/jobs"
	"media-manager/internal/mediafacts"
	"media-manager/internal/storage"
)

func (s *Server) EnqueueMediaFulfillmentAction(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	mediaItemID := uuid.UUID(id)
	item, err := s.settings.GetMediaItem(r.Context(), mediaItemID)
	if err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	var body MediaFulfillmentActionRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	filePath := optionalStringParam(body.FilePath)
	if filePath != "" {
		scopedFilePath, err := s.settings.MediaItemFilePath(r.Context(), mediaItemID, filePath)
		if err != nil {
			writeSettingsError(w, err, "Could not find media file")
			return
		}
		filePath = scopedFilePath
	}
	externalSubtitleID := optionalUUIDString(body.ExternalSubtitleId)
	if externalSubtitleID != "" && !mediaItemHasSubtitle(item, externalSubtitleID) {
		writeError(w, http.StatusNotFound, "media_subtitle_not_found", "Could not find media subtitle")
		return
	}
	otherFileID := optionalUUIDString(body.OtherFileId)
	if otherFileID != "" && !mediaItemHasOtherFile(item, filePath, otherFileID) {
		writeError(w, http.StatusNotFound, "media_other_file_not_found", "Could not find media other file")
		return
	}
	scopedItem := mediafacts.WithLiveFileFacts(item, filePath)
	trackID := optionalUUIDString(body.TrackId)
	var track *storage.MediaFileTrackFact
	if trackID != "" {
		track = mediaItemTrack(scopedItem, filePath, trackID)
		if track == nil {
			writeError(w, http.StatusNotFound, "media_track_not_found", "Could not find media track")
			return
		}
		if filePath == "" {
			filePath = track.FilePath
		}
	}
	if !validTrackScopedFulfillmentAction(w, string(body.Operation), track) {
		return
	}
	targetType := targetTypeString(body.TargetType)
	if targetType == "" && track != nil {
		targetType = track.TrackType
	}
	languageID := optionalStringParam(body.LanguageId)
	if languageID == "" && track != nil && track.LanguageID != nil {
		languageID = *track.LanguageID
	}
	args := jobs.FulfillmentActionArgs{
		MediaItemID:        mediaItemID.String(),
		FilePath:           filePath,
		TargetType:         targetType,
		LanguageID:         languageID,
		TrackID:            trackID,
		OtherFileID:        otherFileID,
		ExternalSubtitleID: externalSubtitleID,
		Manual:             true,
	}
	jobID, err := s.jobs.EnqueueFulfillmentAction(r.Context(), string(body.Operation), args)
	if err != nil {
		if errors.Is(err, jobs.ErrFixedScheduleNotFound) {
			writeError(w, http.StatusBadRequest, "fulfillment_operation_invalid", "Unsupported fulfillment operation")
			return
		}
		writeError(w, http.StatusInternalServerError, "fulfillment_enqueue_failed", "Could not enqueue fulfillment operation")
		s.recordEvent(r.Context(), eventSeverityError, "media", "Fulfillment enqueue failed", map[string]any{"mediaItemId": mediaItemID.String(), "operation": body.Operation, "error": err.Error()})
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Fulfillment queued", map[string]any{"mediaItemId": mediaItemID.String(), "operation": body.Operation, "jobId": jobID})
	writeJSON(w, http.StatusAccepted, JobEnqueueResponse{JobId: jobID, Message: "Fulfillment queued"})
}

func validTrackScopedFulfillmentAction(w http.ResponseWriter, operation string, track *storage.MediaFileTrackFact) bool {
	requiredType := ""
	switch operation {
	case "audio_transcode":
		requiredType = "audio"
	case "video_transcode":
		requiredType = "video"
	default:
		return true
	}
	if track == nil {
		writeError(w, http.StatusBadRequest, "media_track_required", "Fulfillment operation requires a selected media track")
		return false
	}
	if track.TrackType != requiredType {
		writeError(w, http.StatusBadRequest, "media_track_type_invalid", "Fulfillment operation does not support the selected track type")
		return false
	}
	return true
}

func targetTypeString(value *TargetSatisfactionType) string {
	if value == nil {
		return ""
	}
	return string(*value)
}

func optionalUUIDString(value *openapi_types.UUID) string {
	if value == nil {
		return ""
	}
	return value.String()
}

func mediaItemHasSubtitle(item storage.MediaItem, subtitleID string) bool {
	for _, subtitle := range item.ExternalSubtitles {
		if subtitle.ID.String() == subtitleID {
			return true
		}
	}
	return false
}

func mediaItemTrack(item storage.MediaItem, filePath string, trackID string) *storage.MediaFileTrackFact {
	for _, fact := range item.FileFacts {
		if filePath != "" && fact.FilePath != filePath {
			continue
		}
		for _, track := range fact.Tracks {
			if track.ID.String() == trackID {
				return &track
			}
		}
	}
	return nil
}

func mediaItemHasOtherFile(item storage.MediaItem, filePath string, otherFileID string) bool {
	for _, path := range item.FilePaths {
		if filePath != "" && path != filePath {
			continue
		}
		satisfaction := mediaFileSubtitleSatisfaction(
			nil,
			item.SubtitleTargets,
			item.SubtitleMode,
			externalSubtitleLanguagesForPath(item.ExternalSubtitles, item.Sidecars, path),
		)
		files := mediaFileOtherFiles(item.ID, path, item.FilePaths, item.SubtitleTargets, item.SubtitleMode, item.ExternalSubtitles, item.Sidecars, satisfaction)
		for _, file := range files {
			if file.Id != nil && file.Id.String() == otherFileID {
				return true
			}
		}
	}
	return false
}
