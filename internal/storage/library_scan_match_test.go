package storage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestAssignSeriesFolderPathsGroupsNestedSeasons(t *testing.T) {
	root := t.TempDir()
	items := []LibraryScanItem{
		{Path: filepath.Join(root, "Show", "Season 1", "Show.S01E01.mkv"), DetectedMediaKind: "series"},
		{Path: filepath.Join(root, "Show", "Season 2", "Show.S02E01.mkv"), DetectedMediaKind: "series"},
		{Path: filepath.Join(root, "Movie.mkv"), DetectedMediaKind: "movie"},
	}
	assignSeriesFolderPaths(items, root)
	if items[0].SeriesFolderPath != filepath.Join(root, "Show") || items[1].SeriesFolderPath != items[0].SeriesFolderPath {
		t.Fatalf("series folders = %q, %q", items[0].SeriesFolderPath, items[1].SeriesFolderPath)
	}
	if items[2].SeriesFolderPath != "" {
		t.Fatalf("movie series folder = %q", items[2].SeriesFolderPath)
	}
}

func TestResolvedImportedEpisodeNumbersUsesOverrides(t *testing.T) {
	season, episode := int32(4), int32(9)
	gotSeason, gotEpisode, ok := resolvedImportedEpisodeNumbers("Show.S01E02.mkv", &season, &episode)
	if !ok || gotSeason != 4 || gotEpisode != 9 {
		t.Fatalf("resolved episode = S%dE%d, %v; want S4E9", gotSeason, gotEpisode, ok)
	}
	if _, _, ok := resolvedImportedEpisodeNumbers("Show.S01E02.mkv", &season, nil); ok {
		t.Fatal("partial override should not resolve")
	}
}

func TestMatchLibraryScanItemPersistsEpisodeOverrides(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "series")
	if err != nil {
		t.Fatal(err)
	}
	season, episode := int32(4), int32(9)
	scan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{{
		Path: filepath.Join(folder.Path, "Show", "Season 1", "01.mkv"), DetectedTitle: "Show", DetectedMediaKind: "series",
	}})
	if err != nil {
		t.Fatal(err)
	}
	matched, _, err := store.MatchLibraryScanItem(ctx, scan.ID, scan.Items[0].ID, LibraryMatchInput{
		MediaKind: "series", Title: "Show", QualityProfileID: "any", MonitorMode: "all_episodes",
		MinimumAvailability: "released", SeasonNumber: &season, EpisodeNumber: &episode,
	})
	if err != nil {
		t.Fatal(err)
	}
	if matched.SeasonNumber == nil || *matched.SeasonNumber != season || matched.EpisodeNumber == nil || *matched.EpisodeNumber != episode {
		t.Fatalf("matched episode = %#v, want S%dE%d", matched, season, episode)
	}
}

func TestImportLibraryScanItemRecordsOriginalFileNameProvenance(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	fileName := "Linked.Movie.2026.1080p-GRP.mkv"
	path := filepath.Join(folder.Path, fileName)
	if err := os.WriteFile(path, []byte("movie"), 0o644); err != nil {
		t.Fatal(err)
	}
	title := "Linked Movie"
	kind := "movie"
	scan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{{
		Path:              path,
		FileName:          fileName,
		SizeBytes:         1024,
		DetectedTitle:     title,
		DetectedYear:      int32Ptr(2026),
		DetectedMediaKind: kind,
	}})
	if err != nil {
		t.Fatal(err)
	}

	_, item, err := store.ImportLibraryScanItem(ctx, scan.ID, scan.Items[0].ID, LibraryMatchInput{
		MediaKind:           kind,
		Title:               title,
		Year:                int32Ptr(2026),
		Monitored:           true,
		QualityProfileID:    "any",
		MonitorMode:         "movie",
		MinimumAvailability: "released",
	})
	if err != nil {
		t.Fatal(err)
	}
	provenance, err := store.ListMediaComponentProvenance(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(provenance) != 1 {
		t.Fatalf("provenance = %#v", provenance)
	}
	chain := provenance[0].TransformationChain
	if len(chain) != 1 || chain[0]["originalFileName"] != fileName {
		t.Fatalf("transformation chain = %#v, want originalFileName %q", chain, fileName)
	}
}

func TestCreateLibraryScanCanonicalizesItemPath(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(folder.Path, "Nested", "..", "Linked.Movie.2026.mkv")

	scan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{{
		Path:              path,
		FileName:          "stale-name.mkv",
		DetectedTitle:     "Linked Movie",
		DetectedMediaKind: "movie",
	}})
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(folder.Path, "Linked.Movie.2026.mkv")
	if scan.Items[0].Path != expected || scan.Items[0].FileName != filepath.Base(expected) {
		t.Fatalf("scan item path = %q, file name = %q", scan.Items[0].Path, scan.Items[0].FileName)
	}
}

func TestResetLibraryScanItemImportClearsLinkedLibraryRows(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Linked Movie",
		Year:            int32Ptr(2026),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	title := "Linked Movie"
	kind := "movie"
	source := "library"
	input := LibraryScanItemInput{
		Path:              filepath.Join(folder.Path, "Linked.Movie.2026.mkv"),
		FileName:          "Linked.Movie.2026.mkv",
		SizeBytes:         1024,
		DetectedTitle:     title,
		DetectedYear:      int32Ptr(2026),
		DetectedMediaKind: kind,
		Imported:          true,
		MatchedTitle:      &title,
		MatchedYear:       int32Ptr(2026),
		MatchedMediaKind:  &kind,
		MatchSource:       &source,
		MediaItemID:       &item.ID,
	}
	previousScan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{input})
	if err != nil {
		t.Fatal(err)
	}
	scan, err := store.CreateLibraryScan(ctx, folder, []LibraryScanItemInput{{
		Path:              input.Path,
		FileName:          "Linked.Movie.2026.mkv",
		SizeBytes:         1024,
		DetectedTitle:     title,
		DetectedYear:      int32Ptr(2026),
		DetectedMediaKind: kind,
		Imported:          true,
		MatchedTitle:      &title,
		MatchedYear:       int32Ptr(2026),
		MatchedMediaKind:  &kind,
		MatchSource:       &source,
		MediaItemID:       &item.ID,
	}})
	if err != nil {
		t.Fatal(err)
	}

	reset, err := store.ResetLibraryScanItemImport(ctx, scan.ID, scan.Items[0].ID)
	if err != nil {
		t.Fatalf("reset linked scan item: %v", err)
	}
	if reset.Item.Imported || reset.Item.MediaItemID != nil || reset.Item.MatchSource != nil {
		t.Fatalf("reset item = %#v", reset.Item)
	}
	if reset.RemovedMediaItemID == nil || *reset.RemovedMediaItemID != item.ID {
		t.Fatalf("removed media item = %v, want %s", reset.RemovedMediaItemID, item.ID)
	}
	previous, err := store.GetLibraryScan(ctx, previousScan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if previous.Items[0].Imported || previous.Items[0].MediaItemID != nil {
		t.Fatalf("previous scan item still linked = %#v", previous.Items[0])
	}
	if _, err := store.GetMediaItem(ctx, item.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("media item still exists, err=%v", err)
	}
}
