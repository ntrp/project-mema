package storage

import (
	"context"

	"github.com/google/uuid"
)

func (s *SettingsStore) ListRSSEnabledIndexers(ctx context.Context) ([]Indexer, error) {
	rows, err := s.pool.Query(ctx, `
		select `+indexerColumns+`
		from app.indexers
		where enabled = true
			and supports_rss = true
			and app_profile_id = any($1::text[])
			and health_status <> 'disabled'
			and (next_check_at is null or next_check_at <= now())
		order by priority asc, name asc
	`, rssEnabledProfileIDs())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexers := []Indexer{}
	for rows.Next() {
		indexer, err := scanIndexer(rows)
		if err != nil {
			return nil, err
		}
		indexers = append(indexers, indexer)
	}
	return indexers, rows.Err()
}

func (s *SettingsStore) UpdateIndexerRSSMarker(ctx context.Context, id uuid.UUID, input RSSMarkerInput) error {
	_, err := s.pool.Exec(ctx, `
		update app.indexers
		set rss_marker_published_at = $2,
			rss_marker_guid = $3,
			rss_marker_download_url = $4,
			updated_at = now()
		where id = $1
	`, id, input.PublishedAt, input.GUID, input.DownloadURL)
	return err
}

func rssEnabledProfileIDs() []string {
	profiles := DefaultIndexerAppProfiles()
	ids := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		if profile.EnableRSS {
			ids = append(ids, profile.ID)
		}
	}
	return ids
}
