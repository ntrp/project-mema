package httpapi

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"media-manager/internal/jobs"
)

func (s *Server) SearchMediaSubtitles(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	var body ManualSubtitleSearchRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if !validManualSubtitleSearch(w, body) {
		return
	}
	results, logs, err := jobs.SearchManualSubtitles(
		r.Context(),
		s.settings,
		s.subtitles,
		item,
		body.Query,
		body.LanguageId,
		body.FilePath,
	)
	if err != nil {
		writeSettingsError(w, err, "Could not search subtitles")
		return
	}
	writeJSON(w, http.StatusOK, ManualSubtitleSearchResponse{
		Candidates: subtitleCandidateResponses(results),
		Errors:     []string{},
		Logs:       logs,
	})
}

func (s *Server) GrabMediaSubtitle(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	item, err := s.settings.GetMediaItem(r.Context(), uuid.UUID(id))
	if err != nil {
		writeSettingsError(w, err, "Could not find media item")
		return
	}
	var body GrabSubtitleRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	if !validGrabSubtitle(w, body) {
		return
	}
	err = jobs.GrabManualSubtitle(r.Context(), s.settings, s.subtitles, item, jobs.ManualSubtitleGrabInput{
		ProviderID: uuid.UUID(body.ProviderId),
		FilePath:   body.FilePath,
		LanguageID: body.LanguageId,
		Title:      body.Title,
		Format:     body.Format,
		FileID:     body.FileId,
		SourceURL:  body.SourceUrl,
		SourceRef:  body.SourceReference,
	})
	if err != nil {
		writeSettingsError(w, err, "Could not grab subtitle")
		return
	}
	updated, err := s.settings.GetMediaItem(r.Context(), item.ID)
	if err != nil {
		writeSettingsError(w, err, "Could not load media item")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "subtitles", "Subtitle grabbed", map[string]any{
		"mediaItemId": item.ID.String(),
		"providerId":  uuid.UUID(body.ProviderId).String(),
		"languageId":  body.LanguageId,
	})
	writeJSON(w, http.StatusOK, mediaItemResponse(updated))
}

func validManualSubtitleSearch(w http.ResponseWriter, body ManualSubtitleSearchRequest) bool {
	if strings.TrimSpace(body.FilePath) == "" {
		writeError(w, http.StatusBadRequest, "subtitle_file_missing", "Subtitle search file path is required")
		return false
	}
	if strings.TrimSpace(body.LanguageId) == "" {
		writeError(w, http.StatusBadRequest, "subtitle_language_missing", "Subtitle search language is required")
		return false
	}
	return true
}

func validGrabSubtitle(w http.ResponseWriter, body GrabSubtitleRequest) bool {
	if strings.TrimSpace(body.FilePath) == "" ||
		strings.TrimSpace(body.LanguageId) == "" ||
		strings.TrimSpace(body.Title) == "" ||
		strings.TrimSpace(body.Format) == "" {
		writeError(w, http.StatusBadRequest, "subtitle_grab_invalid", "Subtitle grab candidate is incomplete")
		return false
	}
	return true
}

func subtitleCandidateResponses(results []jobs.ManualSubtitleCandidate) []SubtitleCandidate {
	items := make([]SubtitleCandidate, 0, len(results))
	for _, result := range results {
		candidate := result.Candidate
		var fileID *int64
		if candidate.FileID != 0 {
			fileID = &candidate.FileID
		}
		items = append(items, SubtitleCandidate{
			Id:              result.ID,
			Protocol:        result.Protocol,
			ProviderId:      openapi_types.UUID(result.Provider.ID),
			ProviderName:    result.Provider.Name,
			Title:           firstSubtitleCandidateText(candidate.ReleaseName, candidate.SourceRef, candidate.SourceURL),
			LanguageId:      candidate.LanguageID,
			Format:          candidate.Format,
			FileId:          fileID,
			SourceUrl:       optionalStringPtr(candidate.SourceURL),
			SourceReference: optionalStringPtr(candidate.SourceRef),
			Match: SubtitleCandidateMatch{
				Severity: SubtitleCandidateMatchSeverity(result.Match.Severity),
				Label:    result.Match.Label,
				Details:  result.Match.Details,
			},
		})
	}
	return items
}

func firstSubtitleCandidateText(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return "-"
}

func optionalStringPtr(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
