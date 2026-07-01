package httpapi

import "net/http"

func (s *Server) ListDownloadActivity(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}

	activities, err := s.settings.ListDownloadActivity(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "activity_list_failed", "Could not list download activity")
		return
	}
	response := DownloadActivityListResponse{Activities: make([]DownloadActivity, 0, len(activities))}
	for _, activity := range activities {
		response.Activities = append(response.Activities, downloadActivityResponse(activity))
	}
	writeJSON(w, http.StatusOK, response)
}
