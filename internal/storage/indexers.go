package storage

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) ListIndexers(ctx context.Context) ([]Indexer, error) {
	rows, err := storagegen.New(s.pool).ListIndexers(ctx)
	if err != nil {
		return nil, err
	}

	indexers := make([]Indexer, 0, len(rows))
	for _, row := range rows {
		indexers = append(indexers, indexerFromRow(row))
	}
	return indexers, nil
}

func (s *SettingsStore) ListEnabledIndexers(ctx context.Context) ([]Indexer, error) {
	rows, err := storagegen.New(s.pool).ListEnabledIndexers(ctx)
	if err != nil {
		return nil, err
	}

	indexers := make([]Indexer, 0, len(rows))
	for _, row := range rows {
		indexers = append(indexers, indexerFromRow(row))
	}
	return indexers, nil
}

func (s *SettingsStore) GetIndexer(ctx context.Context, id uuid.UUID) (Indexer, error) {
	row, err := storagegen.New(s.pool).GetIndexer(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Indexer{}, ErrNotFound
	}
	return indexerFromRow(row), err
}

func (s *SettingsStore) CreateIndexer(ctx context.Context, input IndexerInput) (Indexer, error) {
	input = normalizeIndexerInput(input)
	row, err := storagegen.New(s.pool).CreateIndexer(ctx, indexerCreateParams(uuid.New(), input))
	return indexerFromRow(row), err
}

func (s *SettingsStore) UpdateIndexer(ctx context.Context, id uuid.UUID, input IndexerInput) (Indexer, error) {
	input = normalizeIndexerInput(input)
	row, err := storagegen.New(s.pool).UpdateIndexer(ctx, indexerUpdateParams(id, input))
	if errors.Is(err, pgx.ErrNoRows) {
		return Indexer{}, ErrNotFound
	}
	return indexerFromRow(row), err
}

func normalizeIndexerInput(input IndexerInput) IndexerInput {
	if input.DefinitionID == "" {
		input.DefinitionID = "generic-torznab"
	}
	if input.Implementation == "" {
		input.Implementation = "Cardigann"
	}
	if input.ImplementationName == "" {
		input.ImplementationName = input.Name
	}
	if input.Protocol == "" {
		input.Protocol = "torrent"
	}
	if input.Privacy == "" {
		input.Privacy = "private"
	}
	if input.Language == "" {
		input.Language = "en-US"
	}
	if input.IndexerURLs == nil {
		input.IndexerURLs = []string{}
	}
	if input.LegacyURLs == nil {
		input.LegacyURLs = []string{}
	}
	if input.Categories == nil {
		input.Categories = []int32{}
	}
	input.MediaTypeScopes = normalizeIndexerMediaTypeScopes(input.MediaTypeScopes)
	input.TagScopes = normalizeIndexerTagScopes(input.TagScopes)
	if len(input.Fields) == 0 {
		input.Fields = json.RawMessage("[]")
	}
	if len(input.Capabilities) == 0 {
		input.Capabilities = json.RawMessage(`{"categories":[],"supportsRawSearch":true,"searchParams":["q"],"tvSearchParams":["q","season","ep"],"movieSearchParams":["q","imdbid"]}`)
	}
	if input.AppProfileID == "" {
		input.AppProfileID = "default"
	}
	if !input.SupportsRSS && !input.SupportsSearch && !input.SupportsRedirect && !input.SupportsPagination {
		input.SupportsRSS = true
		input.SupportsSearch = true
		input.SupportsRedirect = true
		input.SupportsPagination = true
	}
	return input
}

func normalizeIndexerMediaTypeScopes(values []string) []string {
	if len(values) == 0 {
		return []string{"movie", "serie", "anime", "audio", "book"}
	}
	allowed := map[string]bool{"movie": true, "serie": true, "anime": true, "audio": true, "book": true}
	scopes := []string{}
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if !allowed[value] || seen[value] {
			continue
		}
		seen[value] = true
		scopes = append(scopes, value)
	}
	if len(scopes) == 0 {
		return []string{"movie", "serie", "anime", "audio", "book"}
	}
	return scopes
}

func normalizeIndexerTagScopes(values []string) []string {
	tags := normalizeTagNames(values)
	if tags == nil {
		return []string{}
	}
	return tags
}

func (s *SettingsStore) RecordIndexerSuccess(ctx context.Context, id uuid.UUID) (Indexer, error) {
	row, err := storagegen.New(s.pool).RecordIndexerSuccess(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Indexer{}, ErrNotFound
	}
	return indexerFromRow(row), err
}

func (s *SettingsStore) RecordIndexerFailure(
	ctx context.Context,
	id uuid.UUID,
	statusCode *int32,
	message string,
	permanent bool,
	retryUntil *time.Time,
) (Indexer, error) {
	row, err := storagegen.New(s.pool).RecordIndexerFailure(ctx, storagegen.RecordIndexerFailureParams{
		Permanent:  permanent,
		StatusCode: int4Value(statusCode),
		Message:    message,
		RetryUntil: retryUntil,
		ID:         id,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return Indexer{}, ErrNotFound
	}
	return indexerFromRow(row), err
}

func (s *SettingsStore) DeleteIndexer(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteIndexer(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
