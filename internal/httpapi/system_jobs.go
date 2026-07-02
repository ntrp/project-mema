package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"

	"media-manager/internal/storage"
)

func (s *Server) ListSystemJobs(w http.ResponseWriter, r *http.Request, params ListSystemJobsParams) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	jobs, err := s.settings.ListSystemJobs(r.Context(), storage.SystemJobFilters{
		States: stringList(params.Status),
		Queue:  optionalStringParam(params.Queue),
		Kind:   optionalStringParam(params.Kind),
		Query:  optionalStringParam(params.Query),
		Limit:  optionalInt32(params.Limit),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_jobs_failed", "Could not list system jobs")
		return
	}
	response := SystemJobListResponse{Jobs: make([]SystemJob, 0, len(jobs))}
	for _, job := range jobs {
		response.Jobs = append(response.Jobs, systemJobResponse(job))
	}
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) AbortSystemJob(w http.ResponseWriter, r *http.Request, id int64) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if err := s.jobs.AbortJob(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, "system_job_abort_failed", "Could not abort system job")
		s.recordEvent(r.Context(), eventSeverityError, "jobs", "Job abort failed", map[string]any{"jobId": id, "error": err.Error()})
		return
	}
	job, err := s.settings.GetSystemJob(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "system_job_not_found", "Could not find system job")
			return
		}
		writeError(w, http.StatusInternalServerError, "system_job_load_failed", "Could not load system job")
		return
	}
	s.recordEvent(r.Context(), eventSeverityWarning, "jobs", "Job aborted", map[string]any{"jobId": id, "kind": job.Kind, "queue": job.Queue})
	s.events.Publish("system.job.updated", systemJobResponse(job))
	writeJSON(w, http.StatusOK, systemJobResponse(job))
}

func systemJobResponse(job storage.SystemJob) SystemJob {
	return SystemJob{
		Id:          job.ID,
		Status:      job.State,
		Kind:        job.Kind,
		Queue:       job.Queue,
		Attempt:     job.Attempt,
		MaxAttempts: job.MaxAttempts,
		Priority:    job.Priority,
		Args:        job.Args,
		Metadata:    job.Metadata,
		Errors:      job.Errors,
		InfoMessage: job.InfoMessage,
		ScheduledAt: job.ScheduledAt,
		CreatedAt:   job.CreatedAt,
		AttemptedAt: job.AttemptedAt,
		FinalizedAt: job.FinalizedAt,
	}
}

func stringList(value *[]string) []string {
	if value == nil {
		return nil
	}
	values := []string{}
	for _, item := range *value {
		item = strings.TrimSpace(item)
		if item != "" {
			values = append(values, item)
		}
	}
	return values
}

func optionalStringParam(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func optionalInt32(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}
