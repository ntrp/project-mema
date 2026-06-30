package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type LibraryFolder struct {
	ID        uuid.UUID
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LibraryScan struct {
	ID               uuid.UUID
	FolderID         uuid.UUID
	FolderPath       string
	Status           string
	TotalFiles       int32
	AutoMatchedCount int32
	ManualCount      int32
	Items            []LibraryScanItem
	CreatedAt        time.Time
	CompletedAt      *time.Time
}

type LibraryScanItem struct {
	ID                uuid.UUID
	ScanID            uuid.UUID
	Path              string
	FileName          string
	DetectedTitle     string
	DetectedYear      *int32
	DetectedMediaKind string
	Status            string
	MatchedTitle      *string
	MatchedYear       *int32
	MatchedMediaKind  *string
	MediaItemID       *uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type LibraryScanItemInput struct {
	Path              string
	FileName          string
	DetectedTitle     string
	DetectedYear      *int32
	DetectedMediaKind string
	SafeMatch         bool
}

type LibraryMatchInput struct {
	MediaKind           string
	Title               string
	Year                *int32
	Monitored           bool
	QualityProfileID    string
	MonitorMode         string
	MinimumAvailability string
	ExternalProvider    *string
	ExternalID          *string
	Overview            *string
	PosterPath          *string
}

func (s *SettingsStore) ListLibraryFolders(ctx context.Context) ([]LibraryFolder, error) {
	rows, err := s.pool.Query(ctx, `
		select id, path, created_at, updated_at
		from app.library_folders
		order by path asc
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	folders := []LibraryFolder{}
	for rows.Next() {
		folder, err := scanLibraryFolder(rows)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	return folders, rows.Err()
}

func (s *SettingsStore) CreateLibraryFolder(ctx context.Context, path string) (LibraryFolder, error) {
	id := uuid.New()
	return scanLibraryFolder(s.pool.QueryRow(ctx, `
		insert into app.library_folders (id, path)
		values ($1, $2)
		on conflict (path) do update set updated_at = now()
		returning id, path, created_at, updated_at
	`, id, path))
}

func (s *SettingsStore) DeleteLibraryFolder(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `delete from app.library_folders where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) LibraryFolderExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	if err := s.pool.QueryRow(ctx, `select exists(select 1 from app.library_folders where id = $1)`, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *SettingsStore) CreateLibraryScan(ctx context.Context, folder LibraryFolder, inputs []LibraryScanItemInput) (LibraryScan, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return LibraryScan{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	scanID := uuid.New()
	manualCount := int32(0)
	if _, err := tx.Exec(ctx, `
		insert into app.library_scans (
			id, library_folder_id, status, total_files, auto_matched_count, manual_count, completed_at
		)
		values ($1, $2, 'completed', $3, 0, 0, now())
	`, scanID, folder.ID, int32(len(inputs))); err != nil {
		return LibraryScan{}, err
	}

	for _, input := range inputs {
		manualCount++
		if _, err := tx.Exec(ctx, `
			insert into app.library_scan_items (
				id, scan_id, path, file_name, detected_title, detected_year, detected_media_kind,
				status, matched_title, matched_year, matched_media_kind, media_item_id
			)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, uuid.New(), scanID, input.Path, input.FileName, input.DetectedTitle, input.DetectedYear,
			input.DetectedMediaKind, "pending", nil, nil, nil, nil); err != nil {
			return LibraryScan{}, err
		}
	}

	if _, err := tx.Exec(ctx, `
		update app.library_scans
		set auto_matched_count = $2, manual_count = $3
		where id = $1
	`, scanID, int32(0), manualCount); err != nil {
		return LibraryScan{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return LibraryScan{}, err
	}
	return s.GetLibraryScan(ctx, scanID)
}

func (s *SettingsStore) GetLibraryScan(ctx context.Context, id uuid.UUID) (LibraryScan, error) {
	scan, err := scanLibraryScan(s.pool.QueryRow(ctx, `
		select s.id, s.library_folder_id, f.path, s.status, s.total_files, s.auto_matched_count,
			s.manual_count, s.created_at, s.completed_at
		from app.library_scans s
		join app.library_folders f on f.id = s.library_folder_id
		where s.id = $1
	`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return LibraryScan{}, ErrNotFound
	}
	if err != nil {
		return LibraryScan{}, err
	}
	items, err := s.listLibraryScanItems(ctx, id)
	if err != nil {
		return LibraryScan{}, err
	}
	scan.Items = items
	return scan, nil
}

func (s *SettingsStore) MatchLibraryScanItem(ctx context.Context, scanID uuid.UUID, itemID uuid.UUID, input LibraryMatchInput) (LibraryScanItem, MediaItem, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	mediaType, ok := mediaKindToMediaType(input.MediaKind)
	if !ok {
		return LibraryScanItem{}, MediaItem{}, ErrNotFound
	}
	var folderID uuid.UUID
	if err := tx.QueryRow(ctx, `
		select library_folder_id
		from app.library_scans
		where id = $1
	`, scanID).Scan(&folderID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LibraryScanItem{}, MediaItem{}, ErrNotFound
		}
		return LibraryScanItem{}, MediaItem{}, err
	}
	item, err := createMediaItemIfMissing(ctx, tx, MediaItemInput{
		Type:                mediaType,
		Title:               input.Title,
		Year:                input.Year,
		Monitored:           input.Monitored,
		ExternalProvider:    input.ExternalProvider,
		ExternalID:          input.ExternalID,
		Overview:            input.Overview,
		PosterPath:          input.PosterPath,
		MonitorMode:         input.MonitorMode,
		MinimumAvailability: input.MinimumAvailability,
		QualityProfileID:    &input.QualityProfileID,
		LibraryFolderID:     &folderID,
	})
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	updated, err := scanLibraryScanItem(tx.QueryRow(ctx, `
		update app.library_scan_items
		set status = 'manually_added',
			matched_title = $3,
			matched_year = $4,
			matched_media_kind = $5,
			media_item_id = $6,
			updated_at = now()
		where scan_id = $1 and id = $2 and status = 'pending'
		returning id, scan_id, path, file_name, detected_title, detected_year, detected_media_kind,
			status, matched_title, matched_year, matched_media_kind, media_item_id, created_at, updated_at
	`, scanID, itemID, input.Title, input.Year, input.MediaKind, item.ID))
	if errors.Is(err, pgx.ErrNoRows) {
		return LibraryScanItem{}, MediaItem{}, ErrNotFound
	}
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	if _, err := tx.Exec(ctx, `
		update app.library_scans
		set manual_count = (
			select count(*)::integer
			from app.library_scan_items
			where scan_id = $1 and status = 'pending'
		)
		where id = $1
	`, scanID); err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	item, err = getMediaItem(ctx, tx, item.ID)
	if err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return LibraryScanItem{}, MediaItem{}, err
	}
	return updated, item, nil
}

func (s *SettingsStore) listLibraryScanItems(ctx context.Context, scanID uuid.UUID) ([]LibraryScanItem, error) {
	rows, err := s.pool.Query(ctx, `
		select id, scan_id, path, file_name, detected_title, detected_year, detected_media_kind,
			status, matched_title, matched_year, matched_media_kind, media_item_id, created_at, updated_at
		from app.library_scan_items
		where scan_id = $1
		order by status desc, path asc
	`, scanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []LibraryScanItem{}
	for rows.Next() {
		item, err := scanLibraryScanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type mediaItemQuerier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func createMediaItemIfMissing(ctx context.Context, q mediaItemQuerier, input MediaItemInput) (MediaItem, error) {
	input = normalizeMediaItemOptions(input)
	var existingID uuid.UUID
	err := q.QueryRow(ctx, `
		select id
		from app.media_items
		where lower(media_type) = lower($1) and lower(title) = lower($2)
			and (($3::integer is null and year is null) or year = $3)
		order by created_at asc
		limit 1
	`, input.Type, input.Title, input.Year).Scan(&existingID)
	if err == nil {
		mediaFolderPath, err := ensureMediaMainFolder(ctx, q, input)
		if err != nil {
			return MediaItem{}, err
		}
		if _, err := q.Exec(ctx, `
			update app.media_items
			set quality_profile_id = coalesce(quality_profile_id, $2::text),
				library_folder_id = coalesce(library_folder_id, $3::uuid),
				media_folder_path = coalesce(media_folder_path, $4::text),
				monitor_mode = $5,
				minimum_availability = $6,
				manual = $7,
				updated_at = case
					when (quality_profile_id is null and $2::text is not null)
						or (library_folder_id is null and $3::uuid is not null)
						or (media_folder_path is null and $4::text is not null)
						or monitor_mode <> $5
						or minimum_availability <> $6
						or manual <> $7
					then now()
					else updated_at
				end
			where id = $1
		`, existingID, input.QualityProfileID, input.LibraryFolderID, mediaFolderPath, input.MonitorMode, input.MinimumAvailability, input.Manual); err != nil {
			return MediaItem{}, err
		}
		if len(input.Tags) > 0 {
			if err := assignMediaItemTags(ctx, q, existingID, input.Tags); err != nil {
				return MediaItem{}, err
			}
		}
		return getMediaItem(ctx, q, existingID)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return MediaItem{}, err
	}
	id := uuid.New()
	var itemID uuid.UUID
	mediaFolderPath, err := ensureMediaMainFolder(ctx, q, input)
	if err != nil {
		return MediaItem{}, err
	}
	if err := q.QueryRow(ctx, `
		insert into app.media_items (
			id, media_type, title, year, monitored, external_provider, external_id, overview, poster_path, monitor_mode, minimum_availability, manual, quality_profile_id, library_folder_id, media_folder_path
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		returning id
	`, id, input.Type, input.Title, input.Year, input.Monitored, input.ExternalProvider, input.ExternalID, input.Overview, input.PosterPath, input.MonitorMode, input.MinimumAvailability, input.Manual, input.QualityProfileID, input.LibraryFolderID, mediaFolderPath).Scan(&itemID); err != nil {
		return MediaItem{}, err
	}
	if err := assignMediaItemTags(ctx, q, itemID, input.Tags); err != nil {
		return MediaItem{}, err
	}
	return getMediaItem(ctx, q, itemID)
}

func mediaKindToMediaType(kind string) (string, bool) {
	switch kind {
	case "movie", "anime_movie":
		return "movie", true
	case "series", "anime_series":
		return "series", true
	default:
		return "", false
	}
}

func scanLibraryFolder(row pgx.Row) (LibraryFolder, error) {
	var folder LibraryFolder
	err := row.Scan(&folder.ID, &folder.Path, &folder.CreatedAt, &folder.UpdatedAt)
	return folder, err
}

func scanLibraryScan(row pgx.Row) (LibraryScan, error) {
	var scan LibraryScan
	err := row.Scan(
		&scan.ID,
		&scan.FolderID,
		&scan.FolderPath,
		&scan.Status,
		&scan.TotalFiles,
		&scan.AutoMatchedCount,
		&scan.ManualCount,
		&scan.CreatedAt,
		&scan.CompletedAt,
	)
	return scan, err
}

func scanLibraryScanItem(row pgx.Row) (LibraryScanItem, error) {
	var item LibraryScanItem
	err := row.Scan(
		&item.ID,
		&item.ScanID,
		&item.Path,
		&item.FileName,
		&item.DetectedTitle,
		&item.DetectedYear,
		&item.DetectedMediaKind,
		&item.Status,
		&item.MatchedTitle,
		&item.MatchedYear,
		&item.MatchedMediaKind,
		&item.MediaItemID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}
