package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) BlockRelease(ctx context.Context, input ReleaseBlocklistInput) (ReleaseBlocklistItem, error) {
	id := uuid.New()
	title := strings.TrimSpace(input.ReleaseTitle)
	reason := strings.TrimSpace(input.Reason)
	source := strings.TrimSpace(input.Source)
	if title == "" || reason == "" || source == "" {
		return ReleaseBlocklistItem{}, ErrInvalidInput
	}
	if input.Temporary && input.ExpiresAt == nil {
		return ReleaseBlocklistItem{}, ErrInvalidInput
	}
	if !input.Temporary {
		input.ExpiresAt = nil
	}
	downloadURL := optionalText(input.DownloadURL)
	downloadClientName := strings.TrimSpace(input.DownloadClientName)
	return scanReleaseBlocklistItem(s.pool.QueryRow(ctx, `
		insert into app.release_blocklist (
			id, media_item_id, release_title, indexer_name, indexer_protocol, download_client_name, download_url,
			info_url, guid, reason, source, temporary, expires_at
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		on conflict (id) do update set updated_at = now()
		returning id, media_item_id, '', '', release_title, indexer_name, coalesce(nullif(indexer_protocol, ''), 'torrent'),
			download_client_name, download_url, info_url, guid, reason, source, temporary, expires_at, created_at, updated_at
	`, id, input.MediaItemID, title, strings.TrimSpace(input.IndexerName), strings.TrimSpace(input.IndexerProtocol),
		downloadClientName, downloadURL, input.InfoURL, input.GUID, reason, source, input.Temporary, input.ExpiresAt))
}

func (s *SettingsStore) BlockReleaseCandidate(ctx context.Context, release ReleaseCandidateInput, reason string, source string, expiresAt *time.Time) (ReleaseBlocklistItem, error) {
	return s.BlockRelease(ctx, ReleaseBlocklistInput{
		MediaItemID:     release.MediaItemID,
		ReleaseTitle:    release.Title,
		IndexerName:     release.IndexerName,
		IndexerProtocol: release.IndexerProtocol,
		DownloadURL:     release.DownloadURL,
		InfoURL:         release.InfoURL,
		GUID:            release.GUID,
		Reason:          reason,
		Source:          source,
		Temporary:       expiresAt != nil,
		ExpiresAt:       expiresAt,
	})
}

func (s *SettingsStore) BlockReleaseActivity(ctx context.Context, activity DownloadActivity, reason string, source string, expiresAt *time.Time) (ReleaseBlocklistItem, error) {
	return s.BlockRelease(ctx, ReleaseBlocklistInput{
		MediaItemID:        activity.MediaItemID,
		ReleaseTitle:       activity.ReleaseTitle,
		IndexerName:        activity.IndexerName,
		DownloadClientName: activity.DownloadClientName,
		DownloadURL:        activity.DownloadURL,
		Reason:             reason,
		Source:             source,
		Temporary:          expiresAt != nil,
		ExpiresAt:          expiresAt,
	})
}

func (s *SettingsStore) ReleaseBlocked(ctx context.Context, mediaItemID uuid.UUID, release releaseIdentity) (bool, error) {
	item, err := s.findReleaseBlock(ctx, mediaItemID, release)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	return err == nil && item.ID != uuid.Nil, err
}

func (s *SettingsStore) ReleaseCandidateBlocked(ctx context.Context, release ReleaseCandidate) (bool, error) {
	return s.ReleaseBlocked(ctx, release.MediaItemID, releaseIdentityFromCandidate(release))
}

func (s *SettingsStore) ReleaseCandidateInputBlocked(ctx context.Context, release ReleaseCandidateInput) (bool, error) {
	return s.ReleaseBlocked(ctx, release.MediaItemID, releaseIdentityFromInput(release))
}

func (s *SettingsStore) FindReleaseBlock(ctx context.Context, release ReleaseCandidate) (ReleaseBlocklistItem, bool, error) {
	item, err := s.findReleaseBlock(ctx, release.MediaItemID, releaseIdentityFromCandidate(release))
	if errors.Is(err, ErrNotFound) {
		return ReleaseBlocklistItem{}, false, nil
	}
	return item, err == nil, err
}

func (s *SettingsStore) ListReleaseBlocklist(ctx context.Context) ([]ReleaseBlocklistItem, error) {
	rows, err := s.pool.Query(ctx, `
		select b.id, b.media_item_id, m.title, m.media_type, b.release_title, b.indexer_name,
			coalesce(nullif(b.indexer_protocol, ''), i.protocol, 'torrent'), b.download_client_name,
			b.download_url, b.info_url, b.guid, b.reason, b.source, b.temporary, b.expires_at,
			b.created_at, b.updated_at
		from app.release_blocklist b
		join app.media_items m on m.id = b.media_item_id
		left join app.indexers i on lower(i.name) = lower(b.indexer_name)
		order by b.created_at desc
		limit 200
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ReleaseBlocklistItem{}
	for rows.Next() {
		item, err := scanReleaseBlocklistItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SettingsStore) CleanupExpiredReleaseBlocks(ctx context.Context) (int32, error) {
	tag, err := s.pool.Exec(ctx, `
		delete from app.release_blocklist
		where temporary = true and expires_at <= now()
	`)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) DeleteReleaseBlocklistItem(ctx context.Context, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `
		delete from app.release_blocklist
		where id = $1
	`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ClearReleaseBlocklist(ctx context.Context) (int32, error) {
	tag, err := s.pool.Exec(ctx, `delete from app.release_blocklist`)
	if err != nil {
		return 0, err
	}
	return int32(tag.RowsAffected()), nil
}

func (s *SettingsStore) findReleaseBlock(ctx context.Context, mediaItemID uuid.UUID, release releaseIdentity) (ReleaseBlocklistItem, error) {
	item, err := scanReleaseBlocklistItem(s.pool.QueryRow(ctx, `
		select b.id, b.media_item_id, '', '', b.release_title, b.indexer_name, coalesce(nullif(b.indexer_protocol, ''), 'torrent'),
			b.download_client_name, b.download_url, b.info_url, b.guid, b.reason, b.source, b.temporary, b.expires_at,
			b.created_at, b.updated_at
		from app.release_blocklist b
		where b.media_item_id = $1
			and (b.temporary = false or b.expires_at > now())
			and (
				($2::text is not null and b.guid = $2)
				or ($3::text is not null and b.info_url = $3)
				or ($4::text is not null and b.download_url = $4)
				or lower(b.release_title) = lower($5)
			)
		order by b.created_at desc
		limit 1
	`, mediaItemID, release.GUID, release.InfoURL, release.DownloadURL, release.Title))
	if errors.Is(err, pgx.ErrNoRows) {
		return ReleaseBlocklistItem{}, ErrNotFound
	}
	return item, err
}

func scanReleaseBlocklistItem(row pgx.Row) (ReleaseBlocklistItem, error) {
	var item ReleaseBlocklistItem
	err := row.Scan(
		&item.ID,
		&item.MediaItemID,
		&item.MediaTitle,
		&item.MediaType,
		&item.ReleaseTitle,
		&item.IndexerName,
		&item.IndexerProtocol,
		&item.DownloadClientName,
		&item.DownloadURL,
		&item.InfoURL,
		&item.GUID,
		&item.Reason,
		&item.Source,
		&item.Temporary,
		&item.ExpiresAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	return item, err
}

type releaseIdentity struct {
	Title       string
	DownloadURL *string
	InfoURL     *string
	GUID        *string
}

func releaseIdentityFromCandidate(release ReleaseCandidate) releaseIdentity {
	return releaseIdentity{
		Title:       release.Title,
		DownloadURL: optionalText(release.DownloadURL),
		InfoURL:     release.InfoURL,
		GUID:        release.GUID,
	}
}

func releaseIdentityFromInput(release ReleaseCandidateInput) releaseIdentity {
	return releaseIdentity{
		Title:       release.Title,
		DownloadURL: optionalText(release.DownloadURL),
		InfoURL:     release.InfoURL,
		GUID:        release.GUID,
	}
}

func optionalText(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
