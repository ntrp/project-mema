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
	firstIndexerID := uuid.New()
	secondIndexerID := uuid.New()
	if err := store.ReplaceReleaseSearchResults(ctx, item.ID, []ReleaseCandidateInput{
		{
			IndexerID:   &firstIndexerID,
			IndexerName: "Low Seed Indexer",
			IndexerType: "torznab",
			Title:       "Release.Snapshot.2026.1080p.WEB-DL",
			DownloadURL: "http://indexer.test/download/low-" + suffix,
			InfoURL:     stringPtr("http://indexer.test/details/low-" + suffix),
			GUID:        stringPtr("low-" + suffix),
			SizeBytes:   7_000_000_000,
			Seeders:     int32Ptr(12),
			Peers:       int32Ptr(18),
			PublishedAt: &published,
			SearchKind:  "manual",
		},
		{
			IndexerID:        &secondIndexerID,
			IndexerName:      "High Seed Indexer",
			IndexerType:      "newznab",
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

	if err := store.ReplaceReleaseSearchResults(ctx, item.ID, []ReleaseCandidateInput{{
		IndexerName: "Replacement Indexer",
		IndexerType: "torznab",
		Title:       "Release.Snapshot.2026.720p.WEB-DL",
		DownloadURL: "http://indexer.test/download/replacement-" + suffix,
		SizeBytes:   4_000_000_000,
		SearchKind:  "manual",
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
	if _, err := store.GetReleaseCandidate(ctx, release.ID, item.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected stale release lookup to be not found, got %v", err)
	}
}
