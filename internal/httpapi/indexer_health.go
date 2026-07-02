package httpapi

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func (s *Server) recordIndexerTestResult(
	ctx context.Context,
	indexer storage.Indexer,
	result indexers.TestResult,
) {
	if result.Success {
		if _, err := s.settings.RecordIndexerSuccess(ctx, indexer.ID); err != nil {
			slog.Error("indexer status update failed", "indexerName", indexer.Name, "error", err)
		}
		s.recordEvent(ctx, eventSeverityInfo, "indexers", "Indexer status check succeeded", map[string]any{"indexerId": indexer.ID.String(), "indexerName": indexer.Name})
		return
	}

	statusCode := indexers.StatusCodeFromDetails(result.Details)
	permanent := indexers.IsPermanentFailure(statusCode)
	updated, err := s.settings.RecordIndexerFailure(ctx, indexer.ID, statusCode, result.Message, permanent, nil)
	if err != nil {
		slog.Error("indexer status update failed", "indexerName", indexer.Name, "error", err)
	}
	severity := eventSeverityWarning
	if permanent {
		severity = eventSeverityError
	}
	s.recordEvent(ctx, severity, "indexers", "Indexer status check failed", map[string]any{
		"indexerId":   indexer.ID.String(),
		"indexerName": indexer.Name,
		"statusCode":  statusCode,
		"message":     result.Message,
	})
	if err == nil && updated.HealthStatus == "disabled" {
		s.recordEvent(ctx, eventSeverityError, "indexers", "Indexer disabled", map[string]any{
			"indexerId":   indexer.ID.String(),
			"indexerName": indexer.Name,
			"statusCode":  statusCode,
			"message":     result.Message,
		})
	}
}

func (s *Server) refreshedIndexer(ctx context.Context, id uuid.UUID) (storage.Indexer, error) {
	return s.settings.GetIndexer(ctx, id)
}
