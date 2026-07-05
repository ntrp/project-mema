package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestScenarioSCNMedia011StorageReleaseSearchSnapshot(t *testing.T) {
	requireStorageScenario(t, "SCN-MEDIA-011")
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "movie",
		Title:     "Release Snapshot " + suffix,
		Year:      int32Ptr(2026),
		Monitored: true,
	})
	if err != nil {
		t.Fatalf("create media item: %v", err)
	}

	published := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	firstIndexer := releaseCandidateTestIndexer(t, store, "Low Seed Indexer "+suffix)
	secondIndexer := releaseCandidateTestIndexer(t, store, "High Seed Indexer "+suffix)
	firstIndexerID := firstIndexer.ID
	secondIndexerID := secondIndexer.ID
	if err := store.ReplaceReleaseSearchResults(ctx, item.ID, []ReleaseCandidateInput{
		{
			IndexerID:       &firstIndexerID,
			IndexerName:     "Low Seed Indexer",
			IndexerProtocol: "torrent",
			Title:           "Release.Snapshot.2026.1080p.WEB-DL",
			DownloadURL:     "http://indexer.test/download/low-" + suffix,
			InfoURL:         stringPtr("http://indexer.test/details/low-" + suffix),
			GUID:            stringPtr("low-" + suffix),
			SizeBytes:       7_000_000_000,
			Seeders:         int32Ptr(12),
			Peers:           int32Ptr(18),
			PublishedAt:     &published,
			SearchKind:      "manual",
		},
		{
			IndexerID:        &secondIndexerID,
			IndexerName:      "High Seed Indexer",
			IndexerProtocol:  "usenet",
			Title:            "Release.Snapshot.2026.2160p.WEB-DL",
			DownloadURL:      "http://indexer.test/download/high-" + suffix,
			GUID:             stringPtr("high-" + suffix),
			SizeBytes:        15_000_000_000,
			Seeders:          int32Ptr(44),
			Peers:            int32Ptr(50),
			PublishedAt:      &published,
			SearchKind:       "season",
			RequestedSeason:  int32Ptr(1),
			RequestedEpisode: int32Ptr(2),
			Sources: []ReleaseCandidateSource{
				{
					IndexerID:       &secondIndexerID,
					IndexerName:     "High Seed Indexer",
					IndexerProtocol: "usenet",
					Title:           "Release.Snapshot.2026.2160p.WEB-DL",
					DownloadURL:     "http://indexer.test/download/high-" + suffix,
					GUID:            stringPtr("high-" + suffix),
				},
				{
					IndexerID:       &firstIndexerID,
					IndexerName:     "Mirror Indexer",
					IndexerProtocol: "torrent",
					Title:           "Release.Snapshot.2026.2160p.WEB-DL",
					DownloadURL:     "http://mirror.test/download/high-" + suffix,
					InfoURL:         stringPtr("http://mirror.test/details/high-" + suffix),
				},
			},
		},
	}, []string{"torznab timeout", "newznab rejected query"}); err != nil {
		t.Fatalf("replace release search results: %v", err)
	}

	snapshot, err := store.ListReleaseSearchResults(ctx, item.ID)
	if err != nil {
		t.Fatalf("list release search results: %v", err)
	}
	if len(snapshot.Releases) != 2 || snapshot.Releases[0].Title != "Release.Snapshot.2026.2160p.WEB-DL" {
		t.Fatalf("unexpected release ordering: %#v", snapshot.Releases)
	}
	expectStrings(t, snapshot.Errors, []string{"torznab timeout", "newznab rejected query"})

	release, err := store.GetReleaseCandidate(ctx, snapshot.Releases[0].ID, item.ID)
	if err != nil {
		t.Fatalf("get release candidate: %v", err)
	}
	if release.RequestedSeason == nil || *release.RequestedSeason != 1 || release.SearchKind != "season" {
		t.Fatalf("release candidate details = %#v", release)
	}
	if len(release.Sources) != 2 || release.Sources[0].IndexerName != "High Seed Indexer" || release.Sources[1].IndexerName != "Mirror Indexer" {
		t.Fatalf("release candidate sources = %#v", release.Sources)
	}

	if err := store.ReplaceReleaseSearchResults(ctx, item.ID, []ReleaseCandidateInput{{
		IndexerName:     "Replacement Indexer",
		IndexerProtocol: "torrent",
		Title:           "Release.Snapshot.2026.720p.WEB-DL",
		DownloadURL:     "http://indexer.test/download/replacement-" + suffix,
		SizeBytes:       4_000_000_000,
		SearchKind:      "manual",
	}}, nil); err != nil {
		t.Fatalf("replace release search results again: %v", err)
	}

	snapshot, err = store.ListReleaseSearchResults(ctx, item.ID)
	if err != nil {
		t.Fatalf("list replacement release search results: %v", err)
	}
	if len(snapshot.Releases) != 1 || snapshot.Releases[0].Title != "Release.Snapshot.2026.720p.WEB-DL" || len(snapshot.Errors) != 0 {
		t.Fatalf("replacement snapshot = %#v", snapshot)
	}
	if len(snapshot.Releases[0].Sources) != 1 || snapshot.Releases[0].Sources[0].IndexerName != "Replacement Indexer" {
		t.Fatalf("replacement sources = %#v", snapshot.Releases[0].Sources)
	}
	if _, err := store.GetReleaseCandidate(ctx, release.ID, item.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected stale release lookup to be not found, got %v", err)
	}
}

func TestReleaseCandidateIndexerIDClearsWhenIndexerDeleted(t *testing.T) {
	ctx, store := testDBStore(t)
	suffix := uuid.NewString()
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "movie",
		Title:     "Deleted Indexer Candidate " + suffix,
		Monitored: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	indexer := releaseCandidateTestIndexer(t, store, "Deleted Candidate Indexer "+suffix)
	if err := store.ReplaceReleaseSearchResults(ctx, item.ID, []ReleaseCandidateInput{{
		IndexerID:       &indexer.ID,
		IndexerName:     indexer.Name,
		IndexerProtocol: indexer.Protocol,
		Title:           "Deleted.Indexer.Candidate.2026.1080p.WEB-DL",
		DownloadURL:     "http://indexer.test/download/deleted-" + suffix,
		SizeBytes:       4_000_000_000,
		SearchKind:      "manual",
	}}, nil); err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteIndexer(ctx, indexer.ID); err != nil {
		t.Fatal(err)
	}
	snapshot, err := store.ListReleaseSearchResults(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Releases) != 1 || snapshot.Releases[0].IndexerID != nil {
		t.Fatalf("release candidate after indexer delete = %#v", snapshot.Releases)
	}
	if snapshot.Releases[0].IndexerName != indexer.Name {
		t.Fatalf("release candidate should retain indexer name, got %#v", snapshot.Releases[0])
	}
}

func TestReleaseCandidateStoresPersistedEpisodeReference(t *testing.T) {
	ctx, store := testDBStore(t)
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "serie",
		Title:     "Episode Link " + uuid.NewString(),
		Monitored: true,
		MediaMetadataSnapshot: MediaMetadataSnapshot{
			Seasons: []MediaSeason{{
				Name:         "Season 1",
				SeasonNumber: 1,
				Monitored:    false,
				Episodes: []MediaEpisode{{
					Name:          "Pilot",
					EpisodeNumber: 1,
					Monitored:     true,
				}},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	seasons, err := store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	seasonID := seasons[0].ID
	episodeID := seasons[0].Episodes[0].ID
	if err := store.ReplaceReleaseSearchResults(ctx, item.ID, []ReleaseCandidateInput{{
		SeasonID:         &seasonID,
		EpisodeID:        &episodeID,
		IndexerName:      "Episode Indexer",
		IndexerProtocol:  "torrent",
		Title:            "Episode.Link.S01E01.1080p.WEB-DL",
		DownloadURL:      "http://indexer.test/download/episode-link",
		SizeBytes:        4_000_000_000,
		SearchKind:       "episode",
		RequestedSeason:  int32Ptr(1),
		RequestedEpisode: int32Ptr(1),
	}}, nil); err != nil {
		t.Fatal(err)
	}
	snapshot, err := store.ListReleaseSearchResults(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	release := snapshot.Releases[0]
	if release.SeasonID == nil || *release.SeasonID != seasonID || release.EpisodeID == nil || *release.EpisodeID != episodeID {
		t.Fatalf("expected persisted episode reference, got %#v", release)
	}
}

func releaseCandidateTestIndexer(t *testing.T, store *SettingsStore, name string) Indexer {
	t.Helper()
	indexer, err := store.CreateIndexer(t.Context(), IndexerInput{
		Name:       name,
		Protocol:   "torrent",
		BaseURL:    "http://indexer.test/" + uuid.NewString(),
		Categories: []int32{2000},
		Enabled:    true,
		Priority:   100,
	})
	if err != nil {
		t.Fatalf("create indexer: %v", err)
	}
	return indexer
}

func TestScenarioSCNMedia011ReleaseBlocklistMatchesAndExpires(t *testing.T) {
	ctx, store := testDBStore(t)
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:      "movie",
		Title:     "Blocked Release Movie " + uuid.NewString(),
		Monitored: true,
	})
	if err != nil {
		t.Fatalf("create media item: %v", err)
	}
	guid := "blocked-guid-" + uuid.NewString()
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	release := ReleaseCandidateInput{
		MediaItemID:     item.ID,
		IndexerName:     "Scenario Indexer",
		IndexerProtocol: "usenet",
		Title:           "Blocked.Release.2026.1080p",
		DownloadURL:     "https://indexer.test/download/blocked",
		GUID:            &guid,
		SizeBytes:       1,
		SearchKind:      "manual",
	}
	if _, err := store.BlockReleaseCandidate(ctx, release, "missing pieces", "download_failed", &expiresAt); err != nil {
		t.Fatalf("block release: %v", err)
	}
	blocked, err := store.ReleaseCandidateInputBlocked(ctx, ReleaseCandidateInput{
		MediaItemID: item.ID,
		Title:       "Different title",
		GUID:        &guid,
	})
	if err != nil || !blocked {
		t.Fatalf("blocked by guid = %v, %v", blocked, err)
	}
	items, err := store.ListReleaseBlocklist(ctx)
	if err != nil || len(items) != 1 || items[0].MediaTitle != item.Title {
		t.Fatalf("list blocklist = %#v, %v", items, err)
	}
	expired := time.Now().Add(-time.Hour)
	if _, err := store.BlockReleaseCandidate(ctx, ReleaseCandidateInput{
		MediaItemID:     item.ID,
		IndexerName:     "Scenario Indexer",
		IndexerProtocol: "usenet",
		Title:           "Expired.Release.2026.1080p",
		DownloadURL:     "https://indexer.test/download/expired",
		SizeBytes:       1,
		SearchKind:      "manual",
	}, "server unavailable", "download_status_unavailable", &expired); err != nil {
		t.Fatalf("block expired release: %v", err)
	}
	deleted, err := store.CleanupExpiredReleaseBlocks(ctx)
	if err != nil || deleted != 1 {
		t.Fatalf("cleanup expired = %d, %v", deleted, err)
	}
}
