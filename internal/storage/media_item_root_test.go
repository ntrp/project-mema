package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChangingMediaRootPreservesImportedEpisode(t *testing.T) {
	ctx, store := testDBStore(t)
	oldFolder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "series")
	if err != nil {
		t.Fatal(err)
	}
	newFolder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "series")
	if err != nil {
		t.Fatal(err)
	}
	relativeNewRoot, err := filepath.Rel(".", newFolder.Path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.pool.Exec(ctx, `update app.library_folders set path = $2 where id = $1`, newFolder.ID, relativeNewRoot); err != nil {
		t.Fatal(err)
	}
	item := createSeriesItem(t, ctx, store)
	item, err = store.UpdateMediaItemOptions(ctx, item.ID, MediaItemOptionsInput{LibraryFolderID: &oldFolder.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.UpsertMediaSeriesSeasons(ctx, item.ID, []MediaSeriesSeasonInput{{
		SeasonNumber: 1,
		Name:         "Season 1",
		Episodes:     []MediaSeriesEpisodeInput{{EpisodeNumber: 1, Name: "First"}},
	}}); err != nil {
		t.Fatal(err)
	}
	item, err = store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(*item.MediaFolderPath, "Season 1", "Show - S01E01.mkv")
	metadataSource := filepath.Join(*item.MediaFolderPath, "Season 1", "Show - S01E01.nfo")
	writePreviewFile(t, source)
	writePreviewFile(t, metadataSource)
	if err := store.RecordImportedMediaFile(ctx, item, source); err != nil {
		t.Fatal(err)
	}
	originalEpisodeID, err := store.ImportedMediaFileEpisodeID(ctx, item.ID, source)
	if err != nil || originalEpisodeID == nil {
		t.Fatalf("original episode link = %v, %v", originalEpisodeID, err)
	}

	updated, err := store.UpdateMediaItemOptions(ctx, item.ID, MediaItemOptionsInput{LibraryFolderID: &newFolder.ID})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(*updated.MediaFolderPath, "Season 1", "Show - S01E01.mkv")
	metadataTarget := filepath.Join(*updated.MediaFolderPath, "Season 1", "Show - S01E01.nfo")
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("moved file: %v", err)
	}
	if _, err := os.Stat(source); !os.IsNotExist(err) {
		t.Fatalf("source still exists or cannot be checked: %v", err)
	}
	if _, err := os.Stat(metadataTarget); err != nil {
		t.Fatalf("moved metadata file: %v", err)
	}
	var storedMetadataPath string
	if err := store.pool.QueryRow(ctx, `select file_path from app.media_item_sidecars where media_item_id = $1`, item.ID).Scan(&storedMetadataPath); err != nil {
		t.Fatal(err)
	}
	if storedMetadataPath != metadataTarget {
		t.Fatalf("metadata path = %q, want %q", storedMetadataPath, metadataTarget)
	}
	if len(updated.FilePaths) != 1 || updated.FilePaths[0] != target {
		t.Fatalf("file paths = %#v, want %q", updated.FilePaths, target)
	}
	if updated.Status != "downloaded" {
		t.Fatalf("moved item status = %q, want downloaded", updated.Status)
	}
	movedEpisodeID, err := store.ImportedMediaFileEpisodeID(ctx, item.ID, target)
	if err != nil || movedEpisodeID == nil || *movedEpisodeID != *originalEpisodeID {
		t.Fatalf("moved episode link = %v, %v; want %s", movedEpisodeID, err, *originalEpisodeID)
	}
}

func TestChangingMediaRootRollsBackWhenOptionsUpdateFails(t *testing.T) {
	ctx, store := testDBStore(t)
	oldFolder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	newFolder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type: "movie", Title: "Rollback Movie", Monitored: true, LibraryFolderID: &oldFolder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(*item.MediaFolderPath, "Rollback Movie.mkv")
	writePreviewFile(t, source)
	if err := store.RecordImportedMediaFile(ctx, item, source); err != nil {
		t.Fatal(err)
	}
	invalidProfile := "missing-profile"
	if _, err := store.UpdateMediaItemOptions(ctx, item.ID, MediaItemOptionsInput{
		LibraryFolderID: &newFolder.ID, QualityProfileID: &invalidProfile,
	}); err == nil {
		t.Fatal("expected invalid quality profile to fail")
	}
	stored, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.MediaFolderPath == nil || *stored.MediaFolderPath != *item.MediaFolderPath {
		t.Fatalf("media root = %#v, want %q", stored.MediaFolderPath, *item.MediaFolderPath)
	}
	if len(stored.FilePaths) != 1 || stored.FilePaths[0] != source {
		t.Fatalf("file paths after rollback = %#v, want %q", stored.FilePaths, source)
	}
	if _, err := os.Stat(source); err != nil {
		t.Fatalf("rolled back file: %v", err)
	}
}

func TestMediaItemUsesImportedFileRootWhenStoredRootIsTemplate(t *testing.T) {
	ctx, store := testDBStore(t)
	libraryRoot := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, libraryRoot, "movie")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.SaveFileNamingSettings(ctx, FileNamingSettingsInput{
		MovieFileFormat:      "{movie_title} ({release_year})",
		MovieFolderFormat:    "{movie_title} ({release_year})",
		SeriesEpisodeFormat:  "{series_title} - S{season:00}E{episode:00}",
		DailyEpisodeFormat:   "{series_title} - {air_date}",
		AnimeEpisodeFormat:   "{series_title} - S{season:00}E{episode:00}",
		SeriesFolderFormat:   "{series_title}",
		SeasonFolderFormat:   "Season {season}",
		SpecialsFolderFormat: "Specials",
	}); err != nil {
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
	templatedRoot := *item.MediaFolderPath
	actualRoot := filepath.Join(libraryRoot, "Existing Disk Folder")
	source := filepath.Join(actualRoot, "Bad.Name.mkv")
	writePreviewFile(t, source)
	if err := store.RecordImportedMediaFile(ctx, item, source); err != nil {
		t.Fatal(err)
	}

	updated, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.MediaFolderPath == nil || *updated.MediaFolderPath != actualRoot {
		t.Fatalf("media root = %#v, want %q", updated.MediaFolderPath, actualRoot)
	}
	if *updated.MediaFolderPath == templatedRoot {
		t.Fatalf("media root still uses template path %q", templatedRoot)
	}

	preview, err := store.PreviewMediaItemRename(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(actualRoot, "Scenario Movie (2026).mkv")
	if len(preview.Rows) != 1 || preview.Rows[0].ProposedPath != want {
		t.Fatalf("preview = %#v, want proposed path %s", preview, want)
	}
}
