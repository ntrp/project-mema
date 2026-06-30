package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const mediaItemSelectFields = `
	m.id, m.media_type, m.title, m.year, m.monitored, m.external_provider, m.external_id, m.overview, m.poster_path,
	m.monitor_mode, m.minimum_availability, m.manual,
	m.quality_profile_id, mp.name as quality_profile_name,
	case
		when exists (
			select 1
			from app.library_scan_items status_lsi
			where status_lsi.media_item_id = m.id
		) then 'downloaded'
		when exists (
			select 1
			from app.download_activity status_activity
			where status_activity.media_item_id = m.id
				and status_activity.status in ('queued', 'grabbed', 'downloading')
			) then 'downloading'
		when exists (
			select 1
			from app.download_activity status_activity
			where status_activity.media_item_id = m.id
				and status_activity.status = 'completed'
		) then 'downloaded'
		else 'missing'
	end as status,
	m.library_folder_id,
	m.media_folder_path,
	coalesce(lf.path, (
		select lf2.path
		from app.library_scan_items lsi2
		join app.library_scans ls2 on ls2.id = lsi2.scan_id
		join app.library_folders lf2 on lf2.id = ls2.library_folder_id
		where lsi2.media_item_id = m.id
		order by lsi2.updated_at desc
		limit 1
	)) as library_folder_path,
	array(
		select distinct lsi.path
		from app.library_scan_items lsi
		where lsi.media_item_id = m.id
		order by lsi.path
	) as file_paths,
	coalesce(array(
		select t.name
		from app.media_item_tags mit
		join app.tags t on t.id = mit.tag_id
		where mit.media_item_id = m.id
		order by lower(t.name)
	), '{}') as tags,
	m.created_at, m.updated_at
`

const mediaItemJoins = `
	left join app.media_profiles mp on mp.id = m.quality_profile_id
	left join app.library_folders lf on lf.id = m.library_folder_id
`

func (s *SettingsStore) ListMediaItems(ctx context.Context) ([]MediaItem, error) {
	rows, err := s.pool.Query(ctx, `
		select `+mediaItemSelectFields+`
		from app.media_items m
		`+mediaItemJoins+`
		order by created_at desc, title asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []MediaItem{}
	for rows.Next() {
		item, err := scanMediaItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SettingsStore) SearchMediaItems(ctx context.Context, query string, mediaType *string, limit int) ([]MediaItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	rows, err := s.pool.Query(ctx, `
		select `+mediaItemSelectFields+`
		from app.media_items m
		`+mediaItemJoins+`
		where title ilike '%' || $1 || '%'
			and ($2::text is null or media_type = $2)
		order by
			case when lower(title) = lower($1) then 0 else 1 end,
			title asc
		limit $3
	`, query, mediaType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []MediaItem{}
	for rows.Next() {
		item, err := scanMediaItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SettingsStore) GetMediaItem(ctx context.Context, id uuid.UUID) (MediaItem, error) {
	return getMediaItem(ctx, s.pool, id)
}

func getMediaItem(ctx context.Context, q mediaItemQuerier, id uuid.UUID) (MediaItem, error) {
	return scanMediaItemRow(q.QueryRow(ctx, `
		select `+mediaItemSelectFields+`
		from app.media_items m
		`+mediaItemJoins+`
		where m.id = $1
	`, id))
}

func (s *SettingsStore) CreateMediaItem(ctx context.Context, input MediaItemInput) (MediaItem, error) {
	input = normalizeMediaItemOptions(input)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	id := uuid.New()
	var itemID uuid.UUID
	mediaFolderPath, err := ensureMediaMainFolder(ctx, tx, input)
	if err != nil {
		return MediaItem{}, err
	}
	if err := tx.QueryRow(ctx, `
		insert into app.media_items (
			id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, monitor_mode, minimum_availability, manual, quality_profile_id, library_folder_id, media_folder_path
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		returning id
	`, id, input.Type, input.Title, input.Year, input.Monitored, input.ExternalProvider, input.ExternalID, input.Overview, input.PosterPath, input.MonitorMode, input.MinimumAvailability, input.Manual, input.QualityProfileID, input.LibraryFolderID, mediaFolderPath).Scan(&itemID); err != nil {
		return MediaItem{}, err
	}
	if err := assignMediaItemTags(ctx, tx, itemID, input.Tags); err != nil {
		return MediaItem{}, err
	}
	item, err := getMediaItem(ctx, tx, itemID)
	if err != nil {
		return MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return MediaItem{}, err
	}
	return item, nil
}

func (s *SettingsStore) ListMissingMediaItems(ctx context.Context) ([]MediaItem, error) {
	rows, err := s.pool.Query(ctx, `
		select `+mediaItemSelectFields+`
		from app.media_items m
		`+mediaItemJoins+`
		where m.monitored = true
			and m.manual = false
			and not exists (
				select 1
				from app.library_scan_items lsi
				where lsi.media_item_id = m.id
			)
			and not exists (
				select 1
				from app.download_activity activity
				where activity.media_item_id = m.id
					and activity.status in ('queued', 'grabbed', 'downloading')
			)
		order by m.created_at asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []MediaItem{}
	for rows.Next() {
		item, err := scanMediaItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanMediaItemRow(row pgx.Row) (MediaItem, error) {
	item, err := scanMediaItem(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, ErrNotFound
	}
	return item, err
}

func scanMediaItem(row pgx.Row) (MediaItem, error) {
	var item MediaItem
	err := row.Scan(
		&item.ID,
		&item.Type,
		&item.Title,
		&item.Year,
		&item.Monitored,
		&item.ExternalProvider,
		&item.ExternalID,
		&item.Overview,
		&item.PosterPath,
		&item.MonitorMode,
		&item.MinimumAvailability,
		&item.Manual,
		&item.QualityProfileID,
		&item.QualityProfileName,
		&item.Status,
		&item.LibraryFolderID,
		&item.MediaFolderPath,
		&item.LibraryFolderPath,
		&item.FilePaths,
		&item.Tags,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	item.MetadataFilePaths = collectMetadataFilePaths(item.FilePaths)
	return item, err
}
