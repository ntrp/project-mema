package httpapi

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-manager/internal/downloadclients"
	"media-manager/internal/imports"
	"media-manager/internal/storage"
)

func (s *Server) CancelDownloadActivity(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	activity, err := s.settings.GetDownloadActivity(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find download activity")
		return
	}
	if !downloadActivityCancellable(activity.Status) {
		writeError(w, http.StatusBadRequest, "activity_not_cancellable", "Download activity cannot be cancelled")
		return
	}
	if activity.DownloadID != nil {
		if err := s.cancelClientDownload(r.Context(), activity); err != nil {
			writeError(w, http.StatusBadGateway, "download_cancel_failed", err.Error())
			return
		}
	}

	updated, err := s.settings.CancelDownloadActivity(r.Context(), activity.ID)
	if err != nil {
		writeSettingsError(w, err, "Could not cancel download activity")
		return
	}
	updated.MediaTitle = activity.MediaTitle
	updated.MediaType = activity.MediaType
	s.publishDownloadActivity(updated)
	s.recordEvent(r.Context(), eventSeverityWarning, "downloads", "Download activity cancelled", map[string]any{"activityId": activity.ID.String(), "mediaItemId": activity.MediaItemID.String(), "releaseTitle": activity.ReleaseTitle})
	writeJSON(w, http.StatusOK, downloadActivityResponse(updated))
}

func (s *Server) ManualImportDownloadActivity(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var request ManualImportRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	activity, err := s.settings.GetDownloadActivity(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find download activity")
		return
	}
	if !manualImportAllowed(activity) {
		writeError(w, http.StatusBadRequest, "activity_not_import_failed", "Only failed import activity can be manually imported")
		return
	}
	if err := imports.NewService(s.settings).ImportManualDownload(r.Context(), activity, manualImportInput(request)); err != nil {
		writeError(w, http.StatusBadRequest, "manual_import_failed", err.Error())
		return
	}
	progress := 100
	updated, err := s.settings.UpdateDownloadActivityProgress(r.Context(), activity.ID, "completed", &progress, nil)
	if err != nil {
		writeSettingsError(w, err, "Could not update download activity")
		return
	}
	updated.MediaTitle = activity.MediaTitle
	updated.MediaType = activity.MediaType
	updated.MediaYear = activity.MediaYear
	s.publishDownloadActivity(updated)
	s.recordEvent(r.Context(), eventSeverityInfo, "downloads", "Download activity manually imported", map[string]any{"activityId": activity.ID.String(), "mediaItemId": activity.MediaItemID.String(), "releaseTitle": activity.ReleaseTitle})
	writeJSON(w, http.StatusOK, downloadActivityResponse(updated))
}

func (s *Server) DeleteDownloadActivity(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	activity, err := s.settings.GetDownloadActivity(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find download activity")
		return
	}
	if !downloadActivityDeletable(activity.Status) {
		writeError(w, http.StatusBadRequest, "activity_not_deletable", "Only failed or cancelled download activity can be deleted")
		return
	}
	if err := s.settings.DeleteDownloadActivity(r.Context(), activity.ID); err != nil {
		writeSettingsError(w, err, "Could not delete download activity")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "downloads", "Download activity deleted", map[string]any{"activityId": activity.ID.String(), "mediaItemId": activity.MediaItemID.String(), "releaseTitle": activity.ReleaseTitle})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) cancelClientDownload(ctx context.Context, activity storage.DownloadActivity) error {
	clients, err := s.settings.ListEnabledDownloadClients(ctx)
	if err != nil {
		return err
	}
	for _, client := range clients {
		if client.Name != activity.DownloadClientName {
			continue
		}
		result := s.downloadClients.Cancel(ctx, downloadClientConfig(client), downloadclients.CancelRequest{
			DownloadID: *activity.DownloadID,
		})
		if !result.Success {
			return errors.New(result.Message)
		}
		return nil
	}
	return errors.New("Download client is not enabled")
}

func (s *Server) publishDownloadActivity(activity storage.DownloadActivity) {
	s.events.Publish("activity.download.updated", downloadActivityResponse(activity))
}

func downloadActivityCancellable(status string) bool {
	return status == "queued" || status == "grabbed" || status == "downloading"
}

func downloadActivityDeletable(status string) bool {
	return status == "failed" || status == "cancelled"
}

func manualImportAllowed(activity storage.DownloadActivity) bool {
	return activity.Status == "failed" && activity.FailureType != nil && *activity.FailureType == "import"
}

func manualImportInput(request ManualImportRequest) imports.ManualImportInput {
	return imports.ManualImportInput{
		SourcePath:     request.SourcePath,
		TargetFileName: manualString(request.TargetFileName),
		MovieTitle:     manualString(request.MovieTitle),
		Year:           request.Year,
		SeasonNumber:   request.SeasonNumber,
		EpisodeNumber:  request.EpisodeNumber,
		EpisodeTitle:   manualString(request.EpisodeTitle),
		ReleaseGroup:   manualString(request.ReleaseGroup),
		Edition:        manualString(request.Edition),
		Quality:        manualString(request.Quality),
		Languages:      manualStrings(request.Languages),
	}
}

func manualString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func manualStrings(value *[]string) []string {
	if value == nil {
		return nil
	}
	values := make([]string, 0, len(*value))
	for _, item := range *value {
		if item = strings.TrimSpace(item); item != "" {
			values = append(values, item)
		}
	}
	return values
}
