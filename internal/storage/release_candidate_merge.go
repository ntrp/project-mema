package storage

import (
	"context"
	"encoding/json"
	"strings"

	storagegen "media-manager/internal/storage/generated"

	"github.com/google/uuid"
)

func (s *SettingsStore) MergeReleaseCandidates(ctx context.Context, mediaItemID uuid.UUID, releases []ReleaseCandidateInput) error {
	if len(releases) == 0 {
		return nil
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	existing, err := listReleaseCandidates(ctx, tx, mediaItemID)
	if err != nil {
		return err
	}
	for _, release := range releases {
		if matched, ok := findMatchingReleaseCandidate(existing, release); ok {
			if err := updateReleaseCandidate(ctx, tx, matched, release); err != nil {
				return err
			}
			continue
		}
		if err := insertReleaseCandidate(ctx, tx, mediaItemID, release); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func listReleaseCandidates(ctx context.Context, q mediaItemQuerier, mediaItemID uuid.UUID) ([]ReleaseCandidate, error) {
	rows, err := storagegen.New(q).ListReleaseCandidates(ctx, mediaItemID)
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

func updateReleaseCandidate(
	ctx context.Context,
	q mediaItemQuerier,
	existing ReleaseCandidate,
	release ReleaseCandidateInput,
) error {
	sources := appendReleaseCandidateSources(
		ReleaseCandidateSourcesForStored(existing),
		ReleaseCandidateSourcesForInput(release),
	)
	payload, err := json.Marshal(sources)
	if err != nil {
		return err
	}
	_, err = q.Exec(ctx, `
		update app.media_release_candidates
		set indexer_id = $2,
			indexer_name = $3,
			indexer_protocol = $4,
			title = $5,
			download_url = $6,
			info_url = $7,
			guid = $8,
			size_bytes = $9,
			seeders = $10,
			peers = $11,
			published_at = $12,
			search_kind = $13,
			requested_season = $14,
			requested_episode = $15,
			sources = $16,
			updated_at = now()
		where id = $1
	`, existing.ID, release.IndexerID, release.IndexerName, release.IndexerProtocol, release.Title,
		release.DownloadURL, release.InfoURL, release.GUID, release.SizeBytes, release.Seeders, release.Peers,
		release.PublishedAt, release.SearchKind, release.RequestedSeason, release.RequestedEpisode, payload)
	return err
}

func findMatchingReleaseCandidate(
	candidates []ReleaseCandidate,
	release ReleaseCandidateInput,
) (ReleaseCandidate, bool) {
	for _, candidate := range candidates {
		if releaseCandidateFingerprintsMatch(candidate.GUID, release.GUID) {
			return candidate, true
		}
		if releaseCandidateFingerprint(candidate.DownloadURL) != "" &&
			releaseCandidateFingerprint(candidate.DownloadURL) == releaseCandidateFingerprint(release.DownloadURL) {
			return candidate, true
		}
	}
	return ReleaseCandidate{}, false
}

func releaseCandidateFingerprintsMatch(left *string, right *string) bool {
	if left == nil || right == nil {
		return false
	}
	leftValue := releaseCandidateFingerprint(*left)
	rightValue := releaseCandidateFingerprint(*right)
	return leftValue != "" && leftValue == rightValue
}

func releaseCandidateFingerprint(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
