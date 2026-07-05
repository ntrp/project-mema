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

type IndexerProxy struct {
	ID                    uuid.UUID
	Name                  string
	Implementation        string
	Link                  string
	Enabled               bool
	OnHealthIssue         bool
	SupportsOnHealthIssue bool
	IncludeHealthWarnings bool
	TestCommand           string
	Fields                json.RawMessage
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type IndexerProxyInput struct {
	Name                  string
	Implementation        string
	Link                  string
	Enabled               bool
	OnHealthIssue         bool
	SupportsOnHealthIssue bool
	IncludeHealthWarnings bool
	TestCommand           string
	Fields                json.RawMessage
}

func (s *SettingsStore) ListIndexerProxies(ctx context.Context) ([]IndexerProxy, error) {
	rows, err := storagegen.New(s.pool).ListIndexerProxies(ctx)
	if err != nil {
		return nil, err
	}
	proxies := make([]IndexerProxy, 0, len(rows))
	for _, row := range rows {
		proxies = append(proxies, indexerProxyFromRow(row))
	}
	return proxies, nil
}

func (s *SettingsStore) CreateIndexerProxy(ctx context.Context, input IndexerProxyInput) (IndexerProxy, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Link) == "" {
		return IndexerProxy{}, ErrInvalidInput
	}
	row, err := storagegen.New(s.pool).CreateIndexerProxy(ctx, indexerProxyParams(uuid.New(), input))
	return indexerProxyFromRow(row), err
}

func (s *SettingsStore) UpdateIndexerProxy(ctx context.Context, id uuid.UUID, input IndexerProxyInput) (IndexerProxy, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Link) == "" {
		return IndexerProxy{}, ErrInvalidInput
	}
	row, err := storagegen.New(s.pool).UpdateIndexerProxy(ctx, storagegen.UpdateIndexerProxyParams(indexerProxyParams(id, input)))
	return indexerProxyResult(row, err)
}

func (s *SettingsStore) GetIndexerProxy(ctx context.Context, id uuid.UUID) (IndexerProxy, error) {
	row, err := storagegen.New(s.pool).GetIndexerProxy(ctx, id)
	return indexerProxyResult(row, err)
}

func (s *SettingsStore) DeleteIndexerProxy(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteIndexerProxy(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func indexerProxyResult(row storagegen.AppIndexerProxy, err error) (IndexerProxy, error) {
	if errors.Is(err, pgx.ErrNoRows) {
		return IndexerProxy{}, ErrNotFound
	}
	return indexerProxyFromRow(row), err
}

func indexerProxyParams(id uuid.UUID, input IndexerProxyInput) storagegen.CreateIndexerProxyParams {
	return storagegen.CreateIndexerProxyParams{
		ID:                    id,
		Name:                  strings.TrimSpace(input.Name),
		Implementation:        input.Implementation,
		Link:                  strings.TrimSpace(input.Link),
		Enabled:               input.Enabled,
		OnHealthIssue:         input.OnHealthIssue,
		SupportsOnHealthIssue: input.SupportsOnHealthIssue,
		IncludeHealthWarnings: input.IncludeHealthWarnings,
		TestCommand:           input.TestCommand,
		Fields:                input.Fields,
	}
}

func indexerProxyFromRow(row storagegen.AppIndexerProxy) IndexerProxy {
	return IndexerProxy{
		ID:                    row.ID,
		Name:                  row.Name,
		Implementation:        row.Implementation,
		Link:                  row.Link,
		Enabled:               row.Enabled,
		OnHealthIssue:         row.OnHealthIssue,
		SupportsOnHealthIssue: row.SupportsOnHealthIssue,
		IncludeHealthWarnings: row.IncludeHealthWarnings,
		TestCommand:           row.TestCommand,
		Fields:                json.RawMessage(row.Fields),
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}
