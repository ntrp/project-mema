package jobs

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/riverqueue/river/rivertype"

	"media-manager/internal/events"
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

func publishJobUpdated(broker *events.Broker, row *rivertype.JobRow, status string) {
	if broker == nil || row == nil {
		return
	}
	broker.Publish("system.job.updated", jobUpdateFromRow(row, status))
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
