package storage

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func scanIndexerRow(row pgx.Row) (Indexer, error) {
	indexer, err := scanIndexer(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Indexer{}, ErrNotFound
	}
	return indexer, err
}

func scanIndexer(row pgx.Row) (Indexer, error) {
	var indexer Indexer
	err := row.Scan(
		&indexer.ID,
		&indexer.DefinitionID,
		&indexer.Name,
		&indexer.Implementation,
		&indexer.ImplementationName,
		&indexer.Protocol,
		&indexer.Privacy,
		&indexer.Language,
		&indexer.Encoding,
		&indexer.Description,
		&indexer.IndexerURLs,
		&indexer.LegacyURLs,
		&indexer.BaseURL,
		&indexer.APIKey,
		&indexer.Categories,
		&indexer.Fields,
		&indexer.Capabilities,
		&indexer.Redirect,
		&indexer.AppProfileID,
		&indexer.MinimumSeeders,
		&indexer.SeedRatio,
		&indexer.SeedTime,
		&indexer.PackSeedTime,
		&indexer.PreferMagnetURL,
		&indexer.SupportsRSS,
		&indexer.SupportsSearch,
		&indexer.SupportsRedirect,
		&indexer.SupportsPagination,
		&indexer.Enabled,
		&indexer.Priority,
		&indexer.HealthStatus,
		&indexer.LastQueryAt,
		&indexer.LastSuccessAt,
		&indexer.LastFailureAt,
		&indexer.NextCheckAt,
		&indexer.LastStatusCode,
		&indexer.LastError,
		&indexer.FailureCount,
		&indexer.RSSMarkerPublishedAt,
		&indexer.RSSMarkerGUID,
		&indexer.RSSMarkerDownloadURL,
		&indexer.CreatedAt,
		&indexer.UpdatedAt,
	)
	return indexer, err
}

const indexerColumns = `
	id, definition_id, name, implementation, implementation_name, protocol, privacy,
	language, encoding, description, indexer_urls, legacy_urls, base_url, api_key,
	categories, fields, capabilities, redirect, app_profile_id, minimum_seeders,
	seed_ratio, seed_time, pack_seed_time, prefer_magnet_url, supports_rss,
	supports_search, supports_redirect, supports_pagination, enabled, priority,
	health_status, last_query_at, last_success_at, last_failure_at, next_check_at,
	last_status_code, last_error, failure_count, rss_marker_published_at,
	rss_marker_guid, rss_marker_download_url, created_at, updated_at
`
