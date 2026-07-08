package httpapi

import (
	"encoding/json"
	"strings"

	"media-manager/internal/storage"
)

func systemJobScheduleResponses(schedules []storage.SystemJobSchedule) []SystemJobSchedule {
	response := make([]SystemJobSchedule, 0, len(schedules))
	for _, schedule := range schedules {
		response = append(response, systemJobScheduleResponse(schedule))
	}
	return response
}

func systemJobScheduleResponse(schedule storage.SystemJobSchedule) SystemJobSchedule {
	return SystemJobSchedule{
		Id:                    schedule.ID,
		Name:                  schedule.Name,
		Kind:                  schedule.Kind,
		Queue:                 schedule.Queue,
		IntervalSeconds:       schedule.IntervalSeconds,
		Paused:                schedule.Paused,
		NextRunAt:             schedule.NextRunAt,
		ActiveRiverJobId:      schedule.ActiveRiverJobID,
		ActiveStatus:          schedule.ActiveStatus,
		ActiveProgressPercent: schedule.ActiveProgressPercent,
		ActiveProgressLabel:   schedule.ActiveProgressLabel,
		ActiveInfoMessage:     schedule.ActiveInfoMessage,
		LastRiverJobId:        schedule.LastRiverJobID,
		LastStatus:            schedule.LastStatus,
		LastCreatedAt:         schedule.LastCreatedAt,
		LastFinalizedAt:       schedule.LastFinalizedAt,
		CreatedAt:             schedule.CreatedAt,
		UpdatedAt:             schedule.UpdatedAt,
	}
}

func systemJobExecutionResponses(executions []storage.SystemJobExecution) []SystemJobExecution {
	response := make([]SystemJobExecution, 0, len(executions))
	for _, execution := range executions {
		response = append(response, systemJobExecutionResponse(execution))
	}
	return response
}

func systemJobExecutionResponse(execution storage.SystemJobExecution) SystemJobExecution {
	return SystemJobExecution{
		RiverJobId:      execution.RiverJobID,
		ScheduleId:      optionalStringPointer(execution.ScheduleID),
		Classification:  SystemJobExecutionClassification(execution.Classification),
		Status:          execution.Status,
		Kind:            execution.Kind,
		Queue:           execution.Queue,
		Attempt:         execution.Attempt,
		MaxAttempts:     execution.MaxAttempts,
		Priority:        execution.Priority,
		ProgressPercent: execution.ProgressPercent,
		ProgressLabel:   execution.ProgressLabel,
		Args:            execution.Args,
		Metadata:        execution.Metadata,
		Errors:          execution.Errors,
		InfoMessage:     execution.InfoMessage,
		ScheduledAt:     execution.ScheduledAt,
		CreatedAt:       execution.CreatedAt,
		AttemptedAt:     execution.AttemptedAt,
		FinalizedAt:     execution.FinalizedAt,
		UpdatedAt:       execution.UpdatedAt,
	}
}

func systemJobExecutionLogResponses(logs []storage.SystemJobExecutionLog) []SystemJobExecutionLog {
	response := make([]SystemJobExecutionLog, 0, len(logs))
	for _, log := range logs {
		response = append(response, SystemJobExecutionLog{
			Id:         log.ID,
			RiverJobId: log.RiverJobID,
			Severity:   SystemEventSeverity(log.Severity),
			Message:    log.Message,
			Data:       log.Data,
			CreatedAt:  log.CreatedAt,
		})
	}
	return response
}

func systemJobHistorySettingsResponse(settings storage.SystemJobHistorySettings) SystemJobHistorySettings {
	return SystemJobHistorySettings{RetentionDays: settings.RetentionDays}
}

func systemJobExecutionInputFromJob(job storage.SystemJob, status string) storage.SystemJobExecutionInput {
	if strings.TrimSpace(status) == "" {
		status = job.State
	}
	scheduleID := scheduleIDFromMetadata(job.Metadata)
	classification := "one_shot"
	if scheduleID != "" {
		classification = "fixed"
	}
	return storage.SystemJobExecutionInput{
		RiverJobID:     job.ID,
		ScheduleID:     scheduleID,
		Classification: classification,
		Status:         status,
		Kind:           job.Kind,
		Queue:          job.Queue,
		Attempt:        job.Attempt,
		MaxAttempts:    job.MaxAttempts,
		Priority:       job.Priority,
		Args:           []byte(job.Args),
		Metadata:       []byte(job.Metadata),
		Errors:         []byte(job.Errors),
		InfoMessage:    job.InfoMessage,
		ScheduledAt:    job.ScheduledAt,
		CreatedAt:      job.CreatedAt,
		AttemptedAt:    job.AttemptedAt,
		FinalizedAt:    job.FinalizedAt,
	}
}

func optionalStringPointer(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func scheduleIDFromMetadata(metadata string) string {
	var value map[string]any
	if err := json.Unmarshal([]byte(metadata), &value); err != nil {
		return ""
	}
	id, _ := value["river:periodic_job_id"].(string)
	return strings.TrimSpace(id)
}

func historyLimit(value *int32) int32 {
	limit := optionalInt32(value)
	if limit <= 0 {
		return 20
	}
	if limit >= 500 {
		return 499
	}
	return limit
}

func schedulePauseMessage(paused bool) string {
	if paused {
		return "Scheduled job paused"
	}
	return "Scheduled job resumed"
}
