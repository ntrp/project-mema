package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/downloadclients"
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
	writeJSON(w, http.StatusOK, downloadActivityResponse(updated))
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
