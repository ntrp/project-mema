package httpapi

import (
	"net/http"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/downloadrouting"
	"media-manager/internal/jobs"
	"media-manager/internal/storage"
)

func (s *Server) GrabMediaRelease(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}

	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	var body GrabReleaseRequest
	if !decodeJSON(w, r, &body) {
		return
	}

	release, err := s.settings.GetReleaseCandidate(r.Context(), uuid.UUID(body.ReleaseId), item.ID)
	if err != nil {
		writeSettingsError(w, err, "Could not find release candidate")
		return
	}
	if shouldBlockReleaseMismatch(item, release, boolValue(body.OverrideMatch)) {
		writeError(w, http.StatusBadRequest, "release_mismatch", "Release does not match this series/movie")
		return
	}
	if blocked, err := s.settings.ReleaseCandidateBlocked(r.Context(), release); err != nil {
		writeError(w, http.StatusInternalServerError, "release_blocklist_failed", "Could not check release blocklist")
		return
	} else if blocked && !boolValue(body.OverrideMatch) {
		writeError(w, http.StatusBadRequest, "release_blocklisted", "Release is blocklisted")
		return
	}

	clients, err := s.settings.ListEnabledDownloadClients(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "download_client_list_failed", "Could not list enabled download clients")
		return
	}
	if len(clients) == 0 {
		writeError(w, http.StatusBadRequest, "no_download_client", "No enabled download client is configured")
		return
	}

	client, ok := downloadrouting.ClientForProtocol(clients, release.IndexerProtocol)
	if !ok {
		writeError(w, http.StatusBadRequest, "missing_download_client", downloadrouting.MissingClientMessage(release.IndexerProtocol))
		return
	}
	activity, err := s.settings.CreateDownloadActivity(r.Context(), storage.DownloadActivityInput{
		MediaItemID:        item.ID,
		ReleaseTitle:       release.Title,
		IndexerName:        release.IndexerName,
		DownloadClientName: client.Name,
		DownloadURL:        release.DownloadURL,
		Status:             "queued",
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "activity_create_failed", "Could not record download activity")
		return
	}
	activity.MediaTitle = item.Title
	activity.MediaType = item.Type

	jobID, err := s.jobs.EnqueueGrabRelease(r.Context(), jobs.GrabReleaseArgs{
		ActivityID:  activity.ID.String(),
		MediaItemID: item.ID.String(),
		Title:       release.Title,
		DownloadURL: release.DownloadURL,
		IndexerName: release.IndexerName,
		Protocol:    release.IndexerProtocol,
	})
	if err != nil {
		enqueueError := "Could not enqueue download job"
		_, _ = s.settings.FailDownloadActivity(r.Context(), activity.ID, &enqueueError, "download")
		s.recordEvent(r.Context(), eventSeverityError, "downloads", "Download enqueue failed", map[string]any{"mediaItemId": item.ID.String(), "activityId": activity.ID.String(), "releaseTitle": release.Title, "error": err.Error()})
		writeError(w, http.StatusInternalServerError, "download_enqueue_failed", enqueueError)
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "downloads", "Download queued", map[string]any{"mediaItemId": item.ID.String(), "activityId": activity.ID.String(), "releaseTitle": release.Title, "jobId": jobID})
	writeJSON(w, http.StatusAccepted, GrabReleaseResponse{
		JobId:    jobID,
		Message:  "Download queued",
		Activity: downloadActivityResponse(activity),
	})
}

func boolValue(value *bool) bool {
	return value != nil && *value
}

func shouldBlockReleaseMismatch(
	item storage.MediaItem,
	release storage.ReleaseCandidate,
	overrideMatch bool,
) bool {
	return !overrideMatch && decisions.EvaluateReleaseMatch(item, release).Severity == "error"
}
