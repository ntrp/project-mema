package storage

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	storagegen "media-manager/internal/storage/generated"
)

func systemJobScheduleFromRow(row storagegen.ListSystemJobSchedulesRow) SystemJobSchedule {
	schedule := SystemJobSchedule{
		ID:                   row.ID,
		Name:                 row.Name,
		Kind:                 row.Kind,
		Queue:                row.Queue,
		IntervalSeconds:      row.IntervalSeconds,
		IntervalConfigurable: row.IntervalConfigurable,
		HistoryPolicy:        row.HistoryPolicy,
		Paused:               row.Paused,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
	if row.ActiveRiverJobID > 0 {
		schedule.ActiveRiverJobID = &row.ActiveRiverJobID
		schedule.ActiveStatus = row.ActiveStatus
		schedule.ActiveProgressPercent = int4Ptr(row.ActiveProgressPercent)
		schedule.ActiveProgressLabel = row.ActiveProgressLabel
		schedule.ActiveInfoMessage = row.ActiveInfoMessage
	}
	if row.LastRiverJobID > 0 {
		schedule.LastRiverJobID = &row.LastRiverJobID
		schedule.LastStatus = row.LastStatus
		schedule.LastCreatedAt = &row.LastCreatedAt
		schedule.LastFinalizedAt = row.LastFinalizedAt
		next := row.LastCreatedAt.Add(time.Duration(row.IntervalSeconds) * time.Second)
		schedule.NextRunAt = &next
	}
	if schedule.NextRunAt == nil && !schedule.Paused {
		next := schedule.UpdatedAt.Add(time.Duration(row.IntervalSeconds) * time.Second)
		schedule.NextRunAt = &next
	}
	if schedule.Paused {
		schedule.NextRunAt = nil
	}
	return schedule
}

func systemJobExecutionsFromRows(rows []storagegen.AppSystemJobExecution, err error) ([]SystemJobExecution, error) {
	if err != nil {
		return nil, err
	}
	executions := make([]SystemJobExecution, 0, len(rows))
	for _, row := range rows {
		executions = append(executions, systemJobExecutionFromRow(row))
	}
	return executions, nil
}

func systemJobExecutionFromRow(row storagegen.AppSystemJobExecution) SystemJobExecution {
	return SystemJobExecution{
		RiverJobID:      row.RiverJobID,
		ScheduleID:      textString(row.ScheduleID),
		Classification:  row.Classification,
		HistoryPolicy:   row.HistoryPolicy,
		Status:          row.Status,
		Kind:            row.Kind,
		Queue:           row.Queue,
		Attempt:         row.Attempt,
		MaxAttempts:     row.MaxAttempts,
		Priority:        row.Priority,
		ProgressPercent: int4Ptr(row.ProgressPercent),
		ProgressLabel:   row.ProgressLabel,
		Args:            string(row.Args),
		Metadata:        string(row.Metadata),
		Errors:          string(row.Errors),
		InfoMessage:     row.InfoMessage,
		ScheduledAt:     row.ScheduledAt,
		CreatedAt:       row.CreatedAt,
		AttemptedAt:     row.AttemptedAt,
		FinalizedAt:     row.FinalizedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func systemJobExecutionLogFromRow(row storagegen.AppSystemJobExecutionLog) SystemJobExecutionLog {
	var data map[string]any
	_ = json.Unmarshal(row.Data, &data)
	return SystemJobExecutionLog{
		ID:         row.ID,
		RiverJobID: row.RiverJobID,
		Severity:   row.Severity,
		Message:    row.Message,
		Data:       nonNilMap(data),
		CreatedAt:  row.CreatedAt,
	}
}

func nullableText(value string) pgtype.Text {
	value = strings.TrimSpace(value)
	return pgtype.Text{String: value, Valid: value != ""}
}

func nullableInt4(value *int32) pgtype.Int4 {
	return int4Value(value)
}

func textString(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
