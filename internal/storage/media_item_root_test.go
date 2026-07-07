package storage

import (
	"path/filepath"
	"testing"
)

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
