package storage

import (
	"context"
	"testing"
)

func TestReleaseCandidatesPersistCustomFormatFacts(t *testing.T) {
	ctx, store := testDBStore(t)
	item := releaseCandidateMediaItem(t, ctx, store)

	err := store.ReplaceReleaseSearchResults(ctx, item.ID, []ReleaseCandidateInput{{
		Title:             "Scenario.Movie.2026.1080p.WEBDL.Preferred",
		IndexerName:       "Indexer",
		IndexerProtocol:   "torrent",
		DownloadURL:       "https://indexer.test/download/1",
		CustomFormatScore: 100,
		MatchedCustomFormats: []ReleaseCandidateCustomFormatMatch{{
			Name:  "Preferred group",
			Score: 100,
		}},
	}}, nil)
	if err != nil {
		t.Fatal(err)
	}

	snapshot, err := store.ListReleaseSearchResults(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Releases) != 1 {
		t.Fatalf("releases = %#v", snapshot.Releases)
	}
	release := snapshot.Releases[0]
	if release.CustomFormatScore != 100 || len(release.MatchedCustomFormats) != 1 {
		t.Fatalf("custom format facts = %#v", release)
	}
}

func releaseCandidateMediaItem(t *testing.T, ctx context.Context, store *SettingsStore) MediaItem {
	t.Helper()
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Scenario Movie",
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	return item
}
