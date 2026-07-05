package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMediaRenamePreviewExpandsTemplatesAndDetectsConflict(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
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
	source := filepath.Join(*item.MediaFolderPath, "Bad.Name.mkv")
	writePreviewFile(t, source)
	if err := store.RecordImportedMediaFile(ctx, item, source); err != nil {
		t.Fatal(err)
	}

	preview, err := store.PreviewMediaItemRename(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "Scenario Movie (2026)", "Scenario Movie (2026).mkv")
	if len(preview.Rows) != 1 || preview.Rows[0].Status != "safe" || preview.Rows[0].ProposedPath != want {
		t.Fatalf("preview = %#v, want %s", preview, want)
	}
	if _, err := os.Stat(source); err != nil {
		t.Fatalf("preview mutated source: %v", err)
	}

	writePreviewFile(t, want)
	preview, err = store.PreviewMediaItemRename(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if preview.Rows[0].Status != "conflict" {
		t.Fatalf("conflict preview = %#v", preview.Rows[0])
	}
}

func TestApplyMediaItemRenameMovesSafeRowsAndRecordsHistory(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.SaveFileNamingSettings(ctx, defaultFileNamingSettingsInput()); err != nil {
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
	source := filepath.Join(*item.MediaFolderPath, "Bad.Name.mkv")
	writePreviewFile(t, source)
	if err := store.RecordImportedMediaFile(ctx, item, source); err != nil {
		t.Fatal(err)
	}

	result, err := store.ApplyMediaItemRename(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.AppliedCount != 1 || result.FailedCount != 0 {
		t.Fatalf("result = %#v", result)
	}
	destination := filepath.Join(root, "Scenario Movie (2026)", "Scenario Movie (2026).mkv")
	if _, err := os.Stat(source); !os.IsNotExist(err) {
		t.Fatalf("source stat err = %v, want missing", err)
	}
	if _, err := os.Stat(destination); err != nil {
		t.Fatalf("destination missing: %v", err)
	}
	updated, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	expectStrings(t, updated.FilePaths, []string{destination})
	history, err := store.ListMediaFileHistory(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !historyHasOperation(history, "renamed", "succeeded") {
		t.Fatalf("history = %#v", history)
	}
}

func TestApplyMediaItemRenameSkipsConflictRows(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.SaveFileNamingSettings(ctx, defaultFileNamingSettingsInput()); err != nil {
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
	source := filepath.Join(*item.MediaFolderPath, "Bad.Name.mkv")
	destination := filepath.Join(root, "Scenario Movie (2026)", "Scenario Movie (2026).mkv")
	writePreviewFile(t, source)
	writePreviewFile(t, destination)
	if err := store.RecordImportedMediaFile(ctx, item, source); err != nil {
		t.Fatal(err)
	}

	result, err := store.ApplyMediaItemRename(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if result.SkippedCount != 1 || result.AppliedCount != 0 {
		t.Fatalf("result = %#v", result)
	}
	if got := result.Rows[0].Status; got != "skipped" {
		t.Fatalf("row status = %q", got)
	}
}

func TestApplyMediaItemRenameRollsBackWhenRecordIsStale(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
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
	source := filepath.Join(*item.MediaFolderPath, "orphan.mkv")
	destination := filepath.Join(root, "Scenario Movie (2026)", "Scenario Movie (2026).mkv")
	writePreviewFile(t, source)

	row := store.applyMediaRenameRow(ctx, item, MediaRenamePreviewRow{
		CurrentPath:  source,
		ProposedPath: destination,
		Status:       "safe",
		Messages:     []string{},
	})

	if row.Status != "failed" {
		t.Fatalf("row = %#v", row)
	}
	if _, err := os.Stat(source); err != nil {
		t.Fatalf("source should have been restored: %v", err)
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("destination stat err = %v, want missing", err)
	}
}

func TestMediaRenamePreviewDetectsMissingSource(t *testing.T) {
	item := MediaItem{Title: "Scenario", Type: "movie", Year: int32Ptr(2026), LibraryFolderPath: stringPtr(t.TempDir())}
	source := filepath.Join(*item.LibraryFolderPath, "missing.mkv")

	row := mediaRenamePreviewRow(item, DefaultFileNamingSettings(), source)

	if row.Status != "missing" {
		t.Fatalf("row = %#v", row)
	}
}

func TestCheckedRenamePathRejectsRootAndOutsideDestinations(t *testing.T) {
	root := t.TempDir()

	for _, proposed := range []struct {
		folder string
		file   string
	}{
		{folder: root, file: ""},
		{folder: filepath.Join(root, ".."), file: "outside.mkv"},
		{folder: root, file: "%2e%2e/outside.mkv"},
	} {
		_, _, ok := checkedRenamePath(root, proposed.folder, proposed.file)
		if ok {
			t.Fatalf("expected %q/%q to be rejected", proposed.folder, proposed.file)
		}
	}
}

func writePreviewFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("movie"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func historyHasOperation(history []MediaFileHistoryEntry, operation string, status string) bool {
	for _, entry := range history {
		if entry.Operation == operation && entry.Status == status {
			return true
		}
	}
	return false
}

func defaultFileNamingSettingsInput() FileNamingSettingsInput {
	defaults := DefaultFileNamingSettings()
	return FileNamingSettingsInput{
		MovieFileFormat:      defaults.MovieFileFormat,
		MovieFolderFormat:    defaults.MovieFolderFormat,
		SeriesEpisodeFormat:  defaults.SeriesEpisodeFormat,
		DailyEpisodeFormat:   defaults.DailyEpisodeFormat,
		AnimeEpisodeFormat:   defaults.AnimeEpisodeFormat,
		SeriesFolderFormat:   defaults.SeriesFolderFormat,
		SeasonFolderFormat:   defaults.SeasonFolderFormat,
		SpecialsFolderFormat: defaults.SpecialsFolderFormat,
	}
}
