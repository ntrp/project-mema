package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MediaItem struct {
	ID                 uuid.UUID
	Type               string
	Title              string
	Year               *int32
	Monitored          bool
	ExternalProvider   *string
	ExternalID         *string
	Overview           *string
	PosterPath         *string
	QualityProfileID   *string
	QualityProfileName *string
	Status             string
	LibraryFolderID    *uuid.UUID
	LibraryFolderPath  *string
	MediaFolderPath    *string
	FilePaths          []string
	MetadataFilePaths  []string
	Tags               []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type MediaItemInput struct {
	Type             string
	Title            string
	Year             *int32
	Monitored        bool
	ExternalProvider *string
	ExternalID       *string
	Overview         *string
	PosterPath       *string
	QualityProfileID *string
	LibraryFolderID  *uuid.UUID
	Tags             []string
}

type DownloadActivity struct {
	ID                 uuid.UUID
	MediaItemID        uuid.UUID
	MediaTitle         string
	MediaType          string
	ReleaseTitle       string
	IndexerName        string
	DownloadClientName string
	DownloadID         *string
	DownloadURL        string
	Status             string
	Error              *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type DownloadActivityInput struct {
	MediaItemID        uuid.UUID
	ReleaseTitle       string
	IndexerName        string
	DownloadClientName string
	DownloadID         *string
	DownloadURL        string
	Status             string
	Error              *string
}

type ReleaseCandidate struct {
	ID          uuid.UUID
	MediaItemID uuid.UUID
	IndexerID   *uuid.UUID
	IndexerName string
	IndexerType string
	Title       string
	DownloadURL string
	InfoURL     *string
	GUID        *string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReleaseCandidateInput struct {
	MediaItemID uuid.UUID
	IndexerID   *uuid.UUID
	IndexerName string
	IndexerType string
	Title       string
	DownloadURL string
	InfoURL     *string
	GUID        *string
	SizeBytes   int64
	Seeders     *int32
	Peers       *int32
}

type ReleaseSearchSnapshot struct {
	Releases []ReleaseCandidate
	Errors   []string
}

const mediaItemSelectFields = `
	m.id, m.media_type, m.title, m.year, m.monitored, m.external_provider, m.external_id, m.overview, m.poster_path,
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
			id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, quality_profile_id, library_folder_id, media_folder_path
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		returning id
	`, id, input.Type, input.Title, input.Year, input.Monitored, input.ExternalProvider, input.ExternalID, input.Overview, input.PosterPath, input.QualityProfileID, input.LibraryFolderID, mediaFolderPath).Scan(&itemID); err != nil {
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

func (s *SettingsStore) DeleteMediaItem(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.media_items where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ListMissingMediaItems(ctx context.Context) ([]MediaItem, error) {
	rows, err := s.pool.Query(ctx, `
		select `+mediaItemSelectFields+`
		from app.media_items m
		`+mediaItemJoins+`
		where m.monitored = true
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

func (s *SettingsStore) ListEnabledIndexers(ctx context.Context) ([]Indexer, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, type, base_url, api_key, categories, enabled, priority, created_at, updated_at
		from app.indexers
		where enabled = true
		order by priority asc, name asc
	`)
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

func (s *SettingsStore) ListEnabledDownloadClients(ctx context.Context) ([]DownloadClient, error) {
	rows, err := s.pool.Query(ctx, `
		select id, name, type, base_url, username, password, api_key, category, enabled, priority, created_at, updated_at
		from app.download_clients
		where enabled = true
		order by priority asc, name asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []DownloadClient{}
	for rows.Next() {
		client, err := scanDownloadClient(rows)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, rows.Err()
}

func (s *SettingsStore) ReplaceReleaseSearchResults(ctx context.Context, mediaItemID uuid.UUID, releases []ReleaseCandidateInput, searchErrors []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, `delete from app.media_release_candidates where media_item_id = $1`, mediaItemID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `delete from app.media_release_search_errors where media_item_id = $1`, mediaItemID); err != nil {
		return err
	}
	for _, release := range releases {
		if _, err := tx.Exec(ctx, `
			insert into app.media_release_candidates (
				id, media_item_id, indexer_id, indexer_name, indexer_type, title, download_url,
				info_url, guid, size_bytes, seeders, peers
			)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, uuid.New(), mediaItemID, release.IndexerID, release.IndexerName, release.IndexerType, release.Title,
			release.DownloadURL, release.InfoURL, release.GUID, release.SizeBytes, release.Seeders, release.Peers); err != nil {
			return err
		}
	}
	for _, message := range searchErrors {
		if _, err := tx.Exec(ctx, `
			insert into app.media_release_search_errors (id, media_item_id, message)
			values ($1, $2, $3)
		`, uuid.New(), mediaItemID, message); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (s *SettingsStore) GetReleaseCandidate(ctx context.Context, id uuid.UUID, mediaItemID uuid.UUID) (ReleaseCandidate, error) {
	release, err := scanReleaseCandidate(s.pool.QueryRow(ctx, `
		select id, media_item_id, indexer_id, indexer_name, indexer_type, title, download_url,
			info_url, guid, size_bytes, seeders, peers, created_at, updated_at
		from app.media_release_candidates
		where id = $1 and media_item_id = $2
	`, id, mediaItemID))
	if errors.Is(err, pgx.ErrNoRows) {
		return ReleaseCandidate{}, ErrNotFound
	}
	return release, err
}

func (s *SettingsStore) ListReleaseSearchResults(ctx context.Context, mediaItemID uuid.UUID) (ReleaseSearchSnapshot, error) {
	releaseRows, err := s.pool.Query(ctx, `
		select id, media_item_id, indexer_id, indexer_name, indexer_type, title, download_url,
			info_url, guid, size_bytes, seeders, peers, created_at, updated_at
		from app.media_release_candidates
		where media_item_id = $1
		order by coalesce(seeders, -1) desc, size_bytes desc, created_at desc
	`, mediaItemID)
	if err != nil {
		return ReleaseSearchSnapshot{}, err
	}
	defer releaseRows.Close()

	snapshot := ReleaseSearchSnapshot{Releases: []ReleaseCandidate{}, Errors: []string{}}
	for releaseRows.Next() {
		release, err := scanReleaseCandidate(releaseRows)
		if err != nil {
			return ReleaseSearchSnapshot{}, err
		}
		snapshot.Releases = append(snapshot.Releases, release)
	}
	if err := releaseRows.Err(); err != nil {
		return ReleaseSearchSnapshot{}, err
	}

	errorRows, err := s.pool.Query(ctx, `
		select message
		from app.media_release_search_errors
		where media_item_id = $1
		order by created_at asc
	`, mediaItemID)
	if err != nil {
		return ReleaseSearchSnapshot{}, err
	}
	defer errorRows.Close()

	for errorRows.Next() {
		var message string
		if err := errorRows.Scan(&message); err != nil {
			return ReleaseSearchSnapshot{}, err
		}
		snapshot.Errors = append(snapshot.Errors, message)
	}
	return snapshot, errorRows.Err()
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

var metadataFileExtensions = map[string]struct{}{
	".ass":  {},
	".idx":  {},
	".jpeg": {},
	".jpg":  {},
	".nfo":  {},
	".png":  {},
	".srt":  {},
	".ssa":  {},
	".sub":  {},
	".tbn":  {},
	".txt":  {},
	".webp": {},
}

func collectMetadataFilePaths(mediaPaths []string) []string {
	paths := map[string]struct{}{}
	for _, mediaPath := range mediaPaths {
		dir := filepath.Dir(mediaPath)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		mediaBase := strings.TrimSuffix(filepath.Base(mediaPath), filepath.Ext(mediaPath))
		mediaBase = strings.ToLower(mediaBase)
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			ext := strings.ToLower(filepath.Ext(name))
			if _, ok := metadataFileExtensions[ext]; !ok {
				continue
			}
			base := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))
			if !isRelatedMetadataBase(base, mediaBase) {
				continue
			}
			fullPath := filepath.Join(dir, name)
			if fullPath != mediaPath {
				paths[fullPath] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(paths))
	for path := range paths {
		result = append(result, path)
	}
	sort.Strings(result)
	return result
}

func isRelatedMetadataBase(base string, mediaBase string) bool {
	if base == mediaBase || strings.HasPrefix(base, mediaBase+".") || strings.HasPrefix(base, mediaBase+"-") {
		return true
	}
	switch base {
	case "banner", "clearlogo", "cover", "fanart", "folder", "landscape", "movie", "poster":
		return true
	default:
		return false
	}
}

func scanReleaseCandidate(row pgx.Row) (ReleaseCandidate, error) {
	var release ReleaseCandidate
	err := row.Scan(
		&release.ID,
		&release.MediaItemID,
		&release.IndexerID,
		&release.IndexerName,
		&release.IndexerType,
		&release.Title,
		&release.DownloadURL,
		&release.InfoURL,
		&release.GUID,
		&release.SizeBytes,
		&release.Seeders,
		&release.Peers,
		&release.CreatedAt,
		&release.UpdatedAt,
	)
	return release, err
}
