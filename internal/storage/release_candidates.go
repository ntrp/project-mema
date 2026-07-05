package storage

import (
	"context"
	"encoding/json"
	"errors"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *SettingsStore) ReplaceReleaseSearchResults(ctx context.Context, mediaItemID uuid.UUID, releases []ReleaseCandidateInput, searchErrors []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	q := storagegen.New(tx)
	if err := q.ClearReleaseCandidatesForMedia(ctx, mediaItemID); err != nil {
		return err
	}
	if err := q.ClearReleaseSearchErrorsForMedia(ctx, mediaItemID); err != nil {
		return err
	}
	for _, release := range releases {
		if err := insertReleaseCandidate(ctx, tx, mediaItemID, release); err != nil {
			return err
		}
	}
	for _, message := range searchErrors {
		if err := q.AddReleaseSearchError(ctx, storagegen.AddReleaseSearchErrorParams{
			ID:          uuid.New(),
			MediaItemID: mediaItemID,
			Message:     message,
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func insertReleaseCandidate(ctx context.Context, q mediaItemQuerier, mediaItemID uuid.UUID, release ReleaseCandidateInput) error {
	sources, err := json.Marshal(ReleaseCandidateSourcesForInput(release))
	if err != nil {
		return err
	}
	return storagegen.New(q).AddReleaseCandidate(ctx, storagegen.AddReleaseCandidateParams{
		ID:               uuid.New(),
		MediaItemID:      mediaItemID,
		IndexerID:        release.IndexerID,
		IndexerName:      release.IndexerName,
		IndexerProtocol:  release.IndexerProtocol,
		Title:            release.Title,
		DownloadUrl:      release.DownloadURL,
		InfoUrl:          textValue(release.InfoURL),
		Guid:             textValue(release.GUID),
		SizeBytes:        release.SizeBytes,
		Seeders:          int4Value(release.Seeders),
		Peers:            int4Value(release.Peers),
		PublishedAt:      release.PublishedAt,
		SearchKind:       release.SearchKind,
		RequestedSeason:  int4Value(release.RequestedSeason),
		RequestedEpisode: int4Value(release.RequestedEpisode),
		Sources:          sources,
	})
}

func (s *SettingsStore) GetReleaseCandidate(ctx context.Context, id uuid.UUID, mediaItemID uuid.UUID) (ReleaseCandidate, error) {
	row, err := storagegen.New(s.pool).GetReleaseCandidate(ctx, storagegen.GetReleaseCandidateParams{
		ID:          id,
		MediaItemID: mediaItemID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return ReleaseCandidate{}, ErrNotFound
	}
	if err != nil {
		return ReleaseCandidate{}, err
	}
	return releaseCandidateFromRow(row)
}

func (s *SettingsStore) ListReleaseSearchResults(ctx context.Context, mediaItemID uuid.UUID) (ReleaseSearchSnapshot, error) {
	releases, err := s.listReleaseCandidates(ctx, mediaItemID)
	if err != nil {
		return ReleaseSearchSnapshot{}, err
	}
	errors, err := s.listReleaseSearchErrors(ctx, mediaItemID)
	if err != nil {
		return ReleaseSearchSnapshot{}, err
	}
	return ReleaseSearchSnapshot{Releases: releases, Errors: errors}, nil
}

func (s *SettingsStore) listReleaseCandidates(ctx context.Context, mediaItemID uuid.UUID) ([]ReleaseCandidate, error) {
	rows, err := storagegen.New(s.pool).ListReleaseCandidates(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}

	releases := make([]ReleaseCandidate, 0, len(rows))
	for _, row := range rows {
		release, err := releaseCandidateFromRow(row)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}
	return releases, nil
}

func (s *SettingsStore) listReleaseSearchErrors(ctx context.Context, mediaItemID uuid.UUID) ([]string, error) {
	return storagegen.New(s.pool).ListReleaseSearchErrors(ctx, mediaItemID)
}

func releaseCandidateFromRow(row storagegen.AppMediaReleaseCandidate) (ReleaseCandidate, error) {
	release := ReleaseCandidate{
		ID:               row.ID,
		MediaItemID:      row.MediaItemID,
		IndexerID:        row.IndexerID,
		IndexerName:      row.IndexerName,
		IndexerProtocol:  row.IndexerProtocol,
		Title:            row.Title,
		DownloadURL:      row.DownloadUrl,
		InfoURL:          textPtr(row.InfoUrl),
		GUID:             textPtr(row.Guid),
		SizeBytes:        row.SizeBytes,
		Seeders:          int4Ptr(row.Seeders),
		Peers:            int4Ptr(row.Peers),
		PublishedAt:      row.PublishedAt,
		SearchKind:       row.SearchKind,
		RequestedSeason:  int4Ptr(row.RequestedSeason),
		RequestedEpisode: int4Ptr(row.RequestedEpisode),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
	if len(row.Sources) > 0 {
		if err := json.Unmarshal(row.Sources, &release.Sources); err != nil {
			return ReleaseCandidate{}, err
		}
	}
	return release, nil
}
