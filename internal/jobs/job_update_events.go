package jobs

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/riverqueue/river/rivertype"

	"media-manager/internal/events"
	"media-manager/internal/storage"
)

type jobUpdateEvent struct {
	ID          int64      `json:"id"`
	Status      string     `json:"status"`
	Kind        string     `json:"kind"`
	Queue       string     `json:"queue"`
	Attempt     int32      `json:"attempt"`
	MaxAttempts int32      `json:"maxAttempts"`
	Priority    int32      `json:"priority"`
	Args        string     `json:"args"`
	Metadata    string     `json:"metadata"`
	Errors      string     `json:"errors"`
	InfoMessage string     `json:"infoMessage"`
	ScheduledAt time.Time  `json:"scheduledAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	AttemptedAt *time.Time `json:"attemptedAt,omitempty"`
	FinalizedAt *time.Time `json:"finalizedAt,omitempty"`
}

type jobExecutionEvent struct {
	RiverJobID      int64      `json:"riverJobId"`
	ScheduleID      string     `json:"scheduleId,omitempty"`
	Classification  string     `json:"classification"`
	HistoryPolicy   string     `json:"historyPolicy"`
	Status          string     `json:"status"`
	Kind            string     `json:"kind"`
	Queue           string     `json:"queue"`
	Attempt         int32      `json:"attempt"`
	MaxAttempts     int32      `json:"maxAttempts"`
	Priority        int32      `json:"priority"`
	ProgressPercent *int32     `json:"progressPercent,omitempty"`
	ProgressLabel   string     `json:"progressLabel"`
	Args            string     `json:"args"`
	Metadata        string     `json:"metadata"`
	Errors          string     `json:"errors"`
	InfoMessage     string     `json:"infoMessage"`
	ScheduledAt     time.Time  `json:"scheduledAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	AttemptedAt     *time.Time `json:"attemptedAt,omitempty"`
	FinalizedAt     *time.Time `json:"finalizedAt,omitempty"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type jobExecutionContextKey struct{}

func withJobExecution(ctx context.Context, riverJobID int64) context.Context {
	return context.WithValue(ctx, jobExecutionContextKey{}, riverJobID)
}

func jobExecutionID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(jobExecutionContextKey{}).(int64)
	return id, ok && id > 0
}

func publishJobUpdated(broker *events.Broker, row *rivertype.JobRow, status string) {
	if broker == nil || row == nil {
		return
	}
	broker.Publish("system.job.updated", jobUpdateFromRow(row, status))
}

func recordJobUpdated(ctx context.Context, store *storage.SettingsStore, broker *events.Broker, row *rivertype.JobRow, status string) {
	publishJobUpdated(broker, row, status)
	if store == nil || row == nil {
		return
	}
	execution, err := store.UpsertSystemJobExecution(ctx, jobExecutionInputFromRow(row, status))
	if err != nil {
		return
	}
	_, _ = store.CreateSystemJobExecutionLog(ctx, execution.RiverJobID, severityForJobStatus(execution.Status), messageForJobStatus(execution.Status), map[string]any{
		"kind":  execution.Kind,
		"queue": execution.Queue,
	})
	publishJobExecutionUpdated(broker, execution)
}

func recordJobFinished(ctx context.Context, store *storage.SettingsStore, broker *events.Broker, row *rivertype.JobRow, err error) {
	if err != nil {
		recordJobUpdated(ctx, store, broker, row, "retryable")
		return
	}
	recordJobUpdated(ctx, store, broker, row, "completed")
}

func recordJobProgress(ctx context.Context, store *storage.SettingsStore, broker *events.Broker, percent *int32, label string) {
	riverJobID, ok := jobExecutionID(ctx)
	if !ok || store == nil {
		return
	}
	execution, err := store.UpdateSystemJobExecutionProgress(ctx, riverJobID, percent, label)
	if err != nil {
		return
	}
	_, _ = store.CreateSystemJobExecutionLog(ctx, riverJobID, "info", label, nil)
	publishJobExecutionUpdated(broker, execution)
}

func jobUpdateFromRow(row *rivertype.JobRow, status string) jobUpdateEvent {
	if strings.TrimSpace(status) == "" {
		status = string(row.State)
	}
	finalizedAt := row.FinalizedAt
	if finalizedAt == nil && isFinalJobStatus(status) {
		now := time.Now().UTC()
		finalizedAt = &now
	}
	return jobUpdateEvent{
		ID:          row.ID,
		Status:      status,
		Kind:        row.Kind,
		Queue:       row.Queue,
		Attempt:     int32(row.Attempt),
		MaxAttempts: int32(row.MaxAttempts),
		Priority:    int32(row.Priority),
		Args:        jsonText(row.EncodedArgs, "{}"),
		Metadata:    jsonText(row.Metadata, "{}"),
		Errors:      errorsJSON(row.Errors),
		InfoMessage: jobInfoMessage(row.Errors, status),
		ScheduledAt: row.ScheduledAt,
		CreatedAt:   row.CreatedAt,
		AttemptedAt: row.AttemptedAt,
		FinalizedAt: finalizedAt,
	}
}

func publishJobFinished(broker *events.Broker, row *rivertype.JobRow, err error) {
	if err != nil {
		publishJobUpdated(broker, row, "retryable")
		return
	}
	publishJobUpdated(broker, row, "completed")
}

func publishJobExecutionUpdated(broker *events.Broker, execution storage.SystemJobExecution) {
	if broker == nil {
		return
	}
	broker.Publish("system.job.execution.updated", jobExecutionEventFromStorage(execution))
}

func jobExecutionInputFromRow(row *rivertype.JobRow, status string) storage.SystemJobExecutionInput {
	if strings.TrimSpace(status) == "" {
		status = string(row.State)
	}
	finalizedAt := row.FinalizedAt
	if finalizedAt == nil && isFinalJobStatus(status) {
		now := time.Now().UTC()
		finalizedAt = &now
	}
	scheduleID := periodicScheduleID(row.Metadata)
	classification := "one_shot"
	if scheduleID != "" {
		classification = "fixed"
	}
	return storage.SystemJobExecutionInput{
		RiverJobID:     row.ID,
		ScheduleID:     scheduleID,
		Classification: classification,
		Status:         status,
		Kind:           row.Kind,
		Queue:          row.Queue,
		Attempt:        int32(row.Attempt),
		MaxAttempts:    int32(row.MaxAttempts),
		Priority:       int32(row.Priority),
		Args:           []byte(jsonText(row.EncodedArgs, "{}")),
		Metadata:       []byte(jsonText(row.Metadata, "{}")),
		Errors:         []byte(errorsJSON(row.Errors)),
		InfoMessage:    jobInfoMessage(row.Errors, status),
		ScheduledAt:    row.ScheduledAt,
		CreatedAt:      row.CreatedAt,
		AttemptedAt:    row.AttemptedAt,
		FinalizedAt:    finalizedAt,
	}
}

func periodicScheduleID(metadata []byte) string {
	var value map[string]any
	if err := json.Unmarshal(metadata, &value); err != nil {
		return ""
	}
	id, _ := value["river:periodic_job_id"].(string)
	return strings.TrimSpace(id)
}

func severityForJobStatus(status string) string {
	if status == "retryable" || status == "cancelled" || status == "discarded" {
		return "warning"
	}
	return "info"
}

func messageForJobStatus(status string) string {
	switch status {
	case "running":
		return "Job started"
	case "completed":
		return "Job completed"
	case "retryable":
		return "Job will retry"
	case "cancelled":
		return "Job cancelled"
	case "discarded":
		return "Job discarded"
	default:
		return "Job updated"
	}
}

func jobExecutionEventFromStorage(execution storage.SystemJobExecution) jobExecutionEvent {
	return jobExecutionEvent{
		RiverJobID:      execution.RiverJobID,
		ScheduleID:      execution.ScheduleID,
		Classification:  execution.Classification,
		HistoryPolicy:   execution.HistoryPolicy,
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

func jobInfoMessage(errors []rivertype.AttemptError, status string) string {
	if len(errors) == 0 {
		return status
	}
	message := strings.TrimSpace(errors[len(errors)-1].Error)
	if message == "" {
		return status
	}
	return message
}

func errorsJSON(errors []rivertype.AttemptError) string {
	if len(errors) == 0 {
		return "[]"
	}
	data, err := json.Marshal(errors)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func jsonText(data []byte, fallback string) string {
	text := strings.TrimSpace(string(data))
	if text == "" {
		return fallback
	}
	return text
}

func isFinalJobStatus(status string) bool {
	return status == "completed" || status == "cancelled" || status == "discarded"
}
