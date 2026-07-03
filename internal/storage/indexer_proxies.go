package storage

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

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
	rows, err := s.pool.Query(ctx, `select `+indexerProxyColumns+` from app.indexer_proxies order by name asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	proxies := []IndexerProxy{}
	for rows.Next() {
		proxy, err := scanIndexerProxy(rows)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	return proxies, rows.Err()
}

func (s *SettingsStore) CreateIndexerProxy(ctx context.Context, input IndexerProxyInput) (IndexerProxy, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Link) == "" {
		return IndexerProxy{}, ErrInvalidInput
	}
	return scanIndexerProxyRow(s.pool.QueryRow(ctx, `
		insert into app.indexer_proxies (
			id, name, implementation, link, enabled, on_health_issue, supports_on_health_issue,
			include_health_warnings, test_command, fields
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning `+indexerProxyColumns+`
	`, uuid.New(), strings.TrimSpace(input.Name), input.Implementation, strings.TrimSpace(input.Link),
		input.Enabled, input.OnHealthIssue, input.SupportsOnHealthIssue, input.IncludeHealthWarnings,
		input.TestCommand, input.Fields))
}

func (s *SettingsStore) UpdateIndexerProxy(ctx context.Context, id uuid.UUID, input IndexerProxyInput) (IndexerProxy, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.Link) == "" {
		return IndexerProxy{}, ErrInvalidInput
	}
	return scanIndexerProxyRow(s.pool.QueryRow(ctx, `
		update app.indexer_proxies
		set name = $2,
			implementation = $3,
			link = $4,
			enabled = $5,
			on_health_issue = $6,
			supports_on_health_issue = $7,
			include_health_warnings = $8,
			test_command = $9,
			fields = $10,
			updated_at = now()
		where id = $1
		returning `+indexerProxyColumns+`
	`, id, strings.TrimSpace(input.Name), input.Implementation, strings.TrimSpace(input.Link),
		input.Enabled, input.OnHealthIssue, input.SupportsOnHealthIssue, input.IncludeHealthWarnings,
		input.TestCommand, input.Fields))
}

func (s *SettingsStore) GetIndexerProxy(ctx context.Context, id uuid.UUID) (IndexerProxy, error) {
	return scanIndexerProxyRow(s.pool.QueryRow(ctx, `select `+indexerProxyColumns+` from app.indexer_proxies where id = $1`, id))
}

func (s *SettingsStore) DeleteIndexerProxy(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.indexer_proxies where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func scanIndexerProxyRow(row pgx.Row) (IndexerProxy, error) {
	proxy, err := scanIndexerProxy(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return IndexerProxy{}, ErrNotFound
	}
	return proxy, err
}

func scanIndexerProxy(row pgx.Row) (IndexerProxy, error) {
	var proxy IndexerProxy
	err := row.Scan(
		&proxy.ID,
		&proxy.Name,
		&proxy.Implementation,
		&proxy.Link,
		&proxy.Enabled,
		&proxy.OnHealthIssue,
		&proxy.SupportsOnHealthIssue,
		&proxy.IncludeHealthWarnings,
		&proxy.TestCommand,
		&proxy.Fields,
		&proxy.CreatedAt,
		&proxy.UpdatedAt,
	)
	return proxy, err
}

const indexerProxyColumns = `
	id, name, implementation, link, enabled, on_health_issue, supports_on_health_issue,
	include_health_warnings, test_command, fields, created_at, updated_at
`
