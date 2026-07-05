package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	storagegen "media-manager/internal/storage/generated"

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
	row, err := storagegen.New(s.pool).BlockRelease(ctx, storagegen.BlockReleaseParams{
		ID:                 id,
		MediaItemID:        input.MediaItemID,
		ReleaseTitle:       title,
		IndexerName:        strings.TrimSpace(input.IndexerName),
		IndexerProtocol:    strings.TrimSpace(input.IndexerProtocol),
		DownloadClientName: downloadClientName,
		DownloadUrl:        textValue(downloadURL),
		InfoUrl:            textValue(input.InfoURL),
		Guid:               textValue(input.GUID),
		Reason:             reason,
		Source:             source,
		Temporary:          input.Temporary,
		ExpiresAt:          input.ExpiresAt,
	})
	return releaseBlocklistItemFromBlockRow(row), err
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
	rows, err := storagegen.New(s.pool).ListReleaseBlocklist(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ReleaseBlocklistItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, releaseBlocklistItemFromListRow(row))
	}
	return items, nil
}

func (s *SettingsStore) CleanupExpiredReleaseBlocks(ctx context.Context) (int32, error) {
	rows, err := storagegen.New(s.pool).CleanupExpiredReleaseBlocks(ctx)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) DeleteReleaseBlocklistItem(ctx context.Context, id uuid.UUID) error {
	rows, err := storagegen.New(s.pool).DeleteReleaseBlocklistItem(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SettingsStore) ClearReleaseBlocklist(ctx context.Context) (int32, error) {
	rows, err := storagegen.New(s.pool).ClearReleaseBlocklist(ctx)
	if err != nil {
		return 0, err
	}
	return int32(rows), nil
}

func (s *SettingsStore) findReleaseBlock(ctx context.Context, mediaItemID uuid.UUID, release releaseIdentity) (ReleaseBlocklistItem, error) {
	row, err := storagegen.New(s.pool).FindReleaseBlock(ctx, storagegen.FindReleaseBlockParams{
		MediaItemID: mediaItemID,
		Guid:        textValue(release.GUID),
		InfoUrl:     textValue(release.InfoURL),
		DownloadUrl: textValue(release.DownloadURL),
		Title:       release.Title,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return ReleaseBlocklistItem{}, ErrNotFound
	}
	return releaseBlocklistItemFromFindRow(row), err
}

type releaseIdentity struct {
	Title       string
	DownloadURL *string
	InfoURL     *string
	GUID        *string
}

func releaseBlocklistItemFromBlockRow(row storagegen.BlockReleaseRow) ReleaseBlocklistItem {
	return ReleaseBlocklistItem{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		MediaTitle:         row.MediaTitle,
		MediaType:          row.MediaType,
		ReleaseTitle:       row.ReleaseTitle,
		IndexerName:        row.IndexerName,
		IndexerProtocol:    row.IndexerProtocol,
		DownloadClientName: row.DownloadClientName,
		DownloadURL:        textPtr(row.DownloadUrl),
		InfoURL:            textPtr(row.InfoUrl),
		GUID:               textPtr(row.Guid),
		Reason:             row.Reason,
		Source:             row.Source,
		Temporary:          row.Temporary,
		ExpiresAt:          row.ExpiresAt,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func releaseBlocklistItemFromListRow(row storagegen.ListReleaseBlocklistRow) ReleaseBlocklistItem {
	return ReleaseBlocklistItem{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		MediaTitle:         row.MediaTitle,
		MediaType:          row.MediaType,
		ReleaseTitle:       row.ReleaseTitle,
		IndexerName:        row.IndexerName,
		IndexerProtocol:    row.IndexerProtocol,
		DownloadClientName: row.DownloadClientName,
		DownloadURL:        textPtr(row.DownloadUrl),
		InfoURL:            textPtr(row.InfoUrl),
		GUID:               textPtr(row.Guid),
		Reason:             row.Reason,
		Source:             row.Source,
		Temporary:          row.Temporary,
		ExpiresAt:          row.ExpiresAt,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
}

func releaseBlocklistItemFromFindRow(row storagegen.FindReleaseBlockRow) ReleaseBlocklistItem {
	return ReleaseBlocklistItem{
		ID:                 row.ID,
		MediaItemID:        row.MediaItemID,
		MediaTitle:         row.MediaTitle,
		MediaType:          row.MediaType,
		ReleaseTitle:       row.ReleaseTitle,
		IndexerName:        row.IndexerName,
		IndexerProtocol:    row.IndexerProtocol,
		DownloadClientName: row.DownloadClientName,
		DownloadURL:        textPtr(row.DownloadUrl),
		InfoURL:            textPtr(row.InfoUrl),
		GUID:               textPtr(row.Guid),
		Reason:             row.Reason,
		Source:             row.Source,
		Temporary:          row.Temporary,
		ExpiresAt:          row.ExpiresAt,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
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
