package httpapi

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"

	"media-manager/internal/storage"
)

func (s *Server) GetSystemJobsOverview(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	schedules, err := s.settings.ListSystemJobSchedules(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_jobs_overview_failed", "Could not list fixed scheduled jobs")
		return
	}
	oneShotJobs, err := s.settings.ListCurrentOneShotJobExecutions(r.Context(), 100)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_jobs_overview_failed", "Could not list current one-shot jobs")
		return
	}
	settings, err := s.settings.GetSystemJobHistorySettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_jobs_overview_failed", "Could not load job history settings")
		return
	}
	writeJSON(w, http.StatusOK, SystemJobsOverviewResponse{
		Schedules:       systemJobScheduleResponses(schedules),
		OneShotJobs:     systemJobExecutionResponses(oneShotJobs),
		HistorySettings: systemJobHistorySettingsResponse(settings),
	})
}

func (s *Server) PauseSystemJobSchedule(w http.ResponseWriter, r *http.Request, id string) {
	s.updateSystemJobSchedulePaused(w, r, id, true)
}

func (s *Server) ResumeSystemJobSchedule(w http.ResponseWriter, r *http.Request, id string) {
	s.updateSystemJobSchedulePaused(w, r, id, false)
}

func (s *Server) ListSystemJobExecutions(w http.ResponseWriter, r *http.Request, params ListSystemJobExecutionsParams) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	limit := historyLimit(params.Limit)
	executions, err := s.settings.ListSystemJobExecutions(r.Context(), storage.SystemJobExecutionFilters{
		States:     stringList(params.Status),
		ScheduleID: optionalStringParam(params.ScheduleId),
		Kind:       optionalStringParam(params.Kind),
		Queue:      optionalStringParam(params.Queue),
		Query:      optionalStringParam(params.Query),
		Before:     params.Before,
		Limit:      limit + 1,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_job_executions_failed", "Could not list job execution history")
		return
	}
	hasMore := len(executions) > int(limit)
	if hasMore {
		executions = executions[:limit]
	}
	writeJSON(w, http.StatusOK, SystemJobExecutionListResponse{
		Executions: systemJobExecutionResponses(executions),
		HasMore:    hasMore,
	})
}

func (s *Server) ListSystemJobExecutionLogs(w http.ResponseWriter, r *http.Request, riverJobId int64) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	if _, err := s.settings.GetSystemJobExecution(r.Context(), riverJobId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "system_job_execution_not_found", "Could not find job execution")
			return
		}
		writeError(w, http.StatusInternalServerError, "system_job_execution_load_failed", "Could not load job execution")
		return
	}
	logs, err := s.settings.ListSystemJobExecutionLogs(r.Context(), riverJobId, 500)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "system_job_execution_logs_failed", "Could not list job execution logs")
		return
	}
	writeJSON(w, http.StatusOK, SystemJobExecutionLogListResponse{Logs: systemJobExecutionLogResponses(logs)})
}

func (s *Server) UpdateSystemJobHistorySettings(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	var body SystemJobHistorySettings
	if !decodeJSON(w, r, &body) {
		return
	}
	settings, err := s.settings.UpdateSystemJobHistorySettings(r.Context(), body.RetentionDays)
	if err != nil {
		writeError(w, http.StatusBadRequest, "system_job_history_settings_invalid", "Could not update job history settings")
		return
	}
	writeJSON(w, http.StatusOK, systemJobHistorySettingsResponse(settings))
}

func (s *Server) updateSystemJobSchedulePaused(w http.ResponseWriter, r *http.Request, id string, paused bool) {
	if _, ok := s.requireAdmin(w, r); !ok {
		return
	}
	schedule, err := s.settings.SetSystemJobSchedulePaused(r.Context(), id, paused)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || errors.Is(err, storage.ErrNotFound) {
			writeError(w, http.StatusNotFound, "system_job_schedule_not_found", "Could not find fixed scheduled job")
			return
		}
		writeError(w, http.StatusInternalServerError, "system_job_schedule_update_failed", "Could not update fixed scheduled job")
		return
	}
	s.recordEvent(r.Context(), eventSeverityWarning, "jobs", schedulePauseMessage(paused), map[string]any{"scheduleId": schedule.ID, "kind": schedule.Kind})
	writeJSON(w, http.StatusOK, systemJobScheduleResponse(schedule))
}
