package storage

import (
	"context"
	"encoding/json"
	"errors"

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

	if _, err := tx.Exec(ctx, `delete from app.media_release_candidates where media_item_id = $1`, mediaItemID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `delete from app.media_release_search_errors where media_item_id = $1`, mediaItemID); err != nil {
		return err
	}
	for _, release := range releases {
		if err := insertReleaseCandidate(ctx, tx, mediaItemID, release); err != nil {
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

func insertReleaseCandidate(ctx context.Context, q mediaItemQuerier, mediaItemID uuid.UUID, release ReleaseCandidateInput) error {
	sources, err := json.Marshal(ReleaseCandidateSourcesForInput(release))
	if err != nil {
		return err
	}
	_, err = q.Exec(ctx, `
		insert into app.media_release_candidates (
			id, media_item_id, indexer_id, indexer_name, indexer_protocol, title, download_url,
			info_url, guid, size_bytes, seeders, peers, published_at, search_kind,
			requested_season, requested_episode, sources
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`, uuid.New(), mediaItemID, release.IndexerID, release.IndexerName, release.IndexerProtocol, release.Title,
		release.DownloadURL, release.InfoURL, release.GUID, release.SizeBytes, release.Seeders, release.Peers,
		release.PublishedAt, release.SearchKind, release.RequestedSeason, release.RequestedEpisode, sources)
	return err
}

func (s *SettingsStore) GetReleaseCandidate(ctx context.Context, id uuid.UUID, mediaItemID uuid.UUID) (ReleaseCandidate, error) {
	release, err := scanReleaseCandidate(s.pool.QueryRow(ctx, `
		select id, media_item_id, indexer_id, indexer_name, indexer_protocol, title, download_url,
			info_url, guid, size_bytes, seeders, peers, published_at, search_kind,
			requested_season, requested_episode, sources, created_at, updated_at
		from app.media_release_candidates
		where id = $1 and media_item_id = $2
	`, id, mediaItemID))
	if errors.Is(err, pgx.ErrNoRows) {
		return ReleaseCandidate{}, ErrNotFound
	}
	return release, err
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
	rows, err := s.pool.Query(ctx, `
		select id, media_item_id, indexer_id, indexer_name, indexer_protocol, title, download_url,
			info_url, guid, size_bytes, seeders, peers, published_at, search_kind,
			requested_season, requested_episode, sources, created_at, updated_at
		from app.media_release_candidates
		where media_item_id = $1
		order by coalesce(seeders, -1) desc, size_bytes desc, created_at desc
	`, mediaItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	releases := []ReleaseCandidate{}
	for rows.Next() {
		release, err := scanReleaseCandidate(rows)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}
	return releases, rows.Err()
}

func (s *SettingsStore) listReleaseSearchErrors(ctx context.Context, mediaItemID uuid.UUID) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		select message
		from app.media_release_search_errors
		where media_item_id = $1
		order by created_at asc
	`, mediaItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []string{}
	for rows.Next() {
		var message string
		if err := rows.Scan(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func scanReleaseCandidate(row pgx.Row) (ReleaseCandidate, error) {
	var release ReleaseCandidate
	var sources []byte
	err := row.Scan(
		&release.ID,
		&release.MediaItemID,
		&release.IndexerID,
		&release.IndexerName,
		&release.IndexerProtocol,
		&release.Title,
		&release.DownloadURL,
		&release.InfoURL,
		&release.GUID,
		&release.SizeBytes,
		&release.Seeders,
		&release.Peers,
		&release.PublishedAt,
		&release.SearchKind,
		&release.RequestedSeason,
		&release.RequestedEpisode,
		&sources,
		&release.CreatedAt,
		&release.UpdatedAt,
	)
	if err != nil {
		return release, err
	}
	if len(sources) > 0 {
		err = json.Unmarshal(sources, &release.Sources)
	}
	return release, err
}
