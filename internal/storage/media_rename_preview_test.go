package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMediaRenamePreviewExpandsTemplatesAndDetectsConflict(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root)
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

func TestMediaRenamePreviewDetectsMissingSource(t *testing.T) {
	item := MediaItem{Title: "Scenario", Type: "movie", Year: int32Ptr(2026), LibraryFolderPath: stringPtr(t.TempDir())}
	source := filepath.Join(*item.LibraryFolderPath, "missing.mkv")

	row := mediaRenamePreviewRow(item, DefaultFileNamingSettings(), source)

	if row.Status != "missing" {
		t.Fatalf("row = %#v", row)
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
