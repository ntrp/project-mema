package httpapi

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (s *Server) EnqueueMediaComponentAssembly(w http.ResponseWriter, r *http.Request, id ResourceId) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body MediaComponentAssemblyRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	run, err := s.settings.CreateMediaComponentAssemblyRun(
		r.Context(),
		uuid.UUID(id),
		mediaComponentAssemblyInput(body),
	)
	if err != nil {
		writeSettingsError(w, err, "Could not create component assembly")
		return
	}
	jobID, err := s.jobs.EnqueueMediaComponentMux(r.Context(), run.ID)
	if err != nil {
		_, _ = s.settings.FailMediaComponentAssemblyRun(r.Context(), run.ID, "", err.Error())
		writeError(w, http.StatusInternalServerError, "component_assembly_enqueue_failed", "Could not enqueue component assembly")
		return
	}
	run, err = s.settings.AssignMediaComponentAssemblyJob(r.Context(), run.ID, strconv.FormatInt(jobID, 10))
	if err != nil {
		writeSettingsError(w, err, "Could not update component assembly")
		return
	}
	s.recordEvent(r.Context(), eventSeverityInfo, "media", "Component assembly queued", map[string]any{
		"mediaItemId": uuid.UUID(id).String(),
		"runId":       run.ID.String(),
		"jobId":       jobID,
	})
	writeJSON(w, http.StatusAccepted, MediaComponentAssemblyEnqueueResponse{
		JobId:   jobID,
		Message: "Component assembly queued",
		Run:     mediaComponentAssemblyRunResponse(run),
	})
}
