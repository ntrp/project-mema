package imports

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"media-manager/internal/storage"
	"media-manager/internal/testdb"
)

func TestSCNActivity002ManualMovieImportLinksSanitizedTarget(t *testing.T) {
	store := importTestStore(t)
	ctx := context.Background()
	root := t.TempDir()
	sourceRoot := filepath.Join(root, "client")
	libraryPath := filepath.Join(root, "library")
	sourcePath := filepath.Join(sourceRoot, "downloads", "Scenario.Movie.Raw.MP4")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sourcePath, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}
	folder := createImportFolder(t, store, libraryPath)
	createImportPathMapping(t, store, "/client", sourceRoot)
	year := int32(2026)
	item := createImportMediaItem(t, store, storage.MediaItemInput{
		Type:            "movie",
		Title:           "Scenario Movie",
		Year:            &year,
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})

	activityID := uuid.New()
	err := NewService(store).ImportManualDownload(ctx, storage.DownloadActivity{
		ID:          activityID,
		MediaItemID: item.ID,
	}, ManualImportInput{
		SourcePath:   "/client/downloads/Scenario.Movie.Raw.MP4",
		MovieTitle:   "Scenario: Movie",
		Year:         &year,
		Edition:      "Director/Cut",
		Quality:      "WEBDL-1080p",
		Languages:    []string{" English ", "", "Japanese"},
		ReleaseGroup: "Group*Name",
	})
	if err != nil {
		t.Fatal(err)
	}

	imported, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if imported.Status != "downloaded" {
		t.Fatalf("status = %q", imported.Status)
	}
	if len(imported.FilePaths) != 1 {
		t.Fatalf("file paths = %#v", imported.FilePaths)
	}
	wantName := "Scenario Movie (2026) - Director Cut - English Japanese - WEBDL-1080p - Group Name.mp4"
	if filepath.Base(imported.FilePaths[0]) != wantName {
		t.Fatalf("imported file = %q, want %q", filepath.Base(imported.FilePaths[0]), wantName)
	}
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	targetInfo, err := os.Stat(imported.FilePaths[0])
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(sourceInfo, targetInfo) {
		t.Fatalf("imported file is not linked to source")
	}
	attempts, err := store.ListImportAttemptsForActivity(ctx, activityID)
	if err != nil {
		t.Fatal(err)
	}
	if len(attempts) != 1 || attempts[0].Status != "succeeded" || len(attempts[0].InsertedMediaFilePaths) != 1 {
		t.Fatalf("attempts = %#v", attempts)
	}
}

func TestManualSeriesImportAssociatesPersistedEpisode(t *testing.T) {
	store := importTestStore(t)
	ctx := context.Background()
	root := t.TempDir()
	sourcePath := filepath.Join(root, "episode.mkv")
	if err := os.WriteFile(sourcePath, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}
	folder := createImportFolder(t, store, filepath.Join(root, "library"))
	item := createImportMediaItem(t, store, storage.MediaItemInput{
		Type:            "serie",
		Title:           "Scenario Series",
		Monitored:       true,
		LibraryFolderID: &folder.ID,
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{{
				Name:         "Season 1",
				SeasonNumber: 1,
				Monitored:    false,
				Episodes: []storage.MediaEpisode{{
					Name:          "Pilot",
					EpisodeNumber: 1,
					Monitored:     true,
				}},
			}},
		},
	})
	seasons, err := store.ListMediaSeriesSeasons(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	wantEpisodeID := seasons[0].Episodes[0].ID

	err = NewService(store).ImportManualDownload(ctx, storage.DownloadActivity{
		ID:          uuid.New(),
		MediaItemID: item.ID,
	}, ManualImportInput{
		SourcePath:    sourcePath,
		SeasonNumber:  int32Ptr(1),
		EpisodeNumber: int32Ptr(1),
		EpisodeTitle:  "Pilot",
	})
	if err != nil {
		t.Fatal(err)
	}
	imported, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	gotEpisodeID, err := store.ImportedMediaFileEpisodeID(ctx, item.ID, imported.FilePaths[0])
	if err != nil {
		t.Fatal(err)
	}
	if gotEpisodeID == nil || *gotEpisodeID != wantEpisodeID {
		t.Fatalf("episode id = %v, want %s", gotEpisodeID, wantEpisodeID)
	}
}

func TestSCNActivity002ManualSeriesImportRequiresEpisodeCoordinates(t *testing.T) {
	store := importTestStore(t)
	ctx := context.Background()
	root := t.TempDir()
	sourcePath := filepath.Join(root, "episode.mkv")
	if err := os.WriteFile(sourcePath, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}
	folder := createImportFolder(t, store, filepath.Join(root, "library"))
	item := createImportMediaItem(t, store, storage.MediaItemInput{
		Type:            "serie",
		Title:           "Scenario Series",
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})

	activityID := uuid.New()
	err := NewService(store).ImportManualDownload(ctx, storage.DownloadActivity{
		ID:          activityID,
		MediaItemID: item.ID,
	}, ManualImportInput{SourcePath: sourcePath})

	if err == nil || err.Error() != "season and episode are required for series imports" {
		t.Fatalf("error = %v", err)
	}
	attempts, err := store.ListImportAttemptsForActivity(ctx, activityID)
	if err != nil {
		t.Fatal(err)
	}
	if len(attempts) != 1 || attempts[0].Status != "failed" || *attempts[0].FailureStage != "select_source" {
		t.Fatalf("attempts = %#v", attempts)
	}
}

func TestManualImportRestoresMovedSourceWhenMediaRecordFails(t *testing.T) {
	store := importTestStore(t)
	ctx := context.Background()
	root := t.TempDir()
	sourcePath := filepath.Join(root, "episode.mkv")
	if err := os.WriteFile(sourcePath, []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}
	libraryPath := filepath.Join(root, "library")
	folder := createImportFolder(t, store, libraryPath)
	item := createImportMediaItem(t, store, storage.MediaItemInput{
		Type:            "serie",
		Title:           "Scenario Series",
		Monitored:       true,
		LibraryFolderID: &folder.ID,
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			Seasons: []storage.MediaSeason{{
				Name:         "Season 1",
				SeasonNumber: 1,
				Monitored:    false,
				Episodes: []storage.MediaEpisode{{
					Name:          "Pilot",
					EpisodeNumber: 1,
					Monitored:     true,
				}},
			}},
		},
	})

	activityID := uuid.New()
	err := NewService(store).ImportManualDownload(ctx, storage.DownloadActivity{
		ID:          activityID,
		MediaItemID: item.ID,
	}, ManualImportInput{
		SourcePath:     sourcePath,
		TargetFileName: "Scenario.Series.S01E99.mkv",
		ImportMode:     ImportModeMove,
	})

	if err == nil || !strings.Contains(err.Error(), "episode import target S01E99 is not known") {
		t.Fatalf("error = %v", err)
	}
	if _, err := os.Stat(sourcePath); err != nil {
		t.Fatalf("source was not restored: %v", err)
	}
	targetPath := filepath.Join(libraryPath, "Scenario.Series.S01E99.mkv")
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		t.Fatalf("target stat err = %v, want missing", err)
	}
	attempts, err := store.ListImportAttemptsForActivity(ctx, activityID)
	if err != nil {
		t.Fatal(err)
	}
	if len(attempts) != 1 || attempts[0].Status != "failed" || *attempts[0].FailureStage != "record_media_file" {
		t.Fatalf("attempts = %#v", attempts)
	}
	if len(attempts[0].CreatedTargets) != 1 || attempts[0].CreatedTargets[0] != targetPath {
		t.Fatalf("created targets = %#v", attempts[0].CreatedTargets)
	}
}

func importTestStore(t *testing.T) *storage.SettingsStore {
	t.Helper()
	databaseURL := testdb.Create(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := storage.EnsureSchema(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	store := storage.NewSettingsStore(pool)
	if err := store.EnsureDefaultFileNamingSettings(ctx); err != nil {
		t.Fatal(err)
	}
	return store
}

func createImportFolder(t *testing.T, store *storage.SettingsStore, path string) storage.LibraryFolder {
	t.Helper()
	folder, err := store.CreateLibraryFolder(context.Background(), path, "movie")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.DeleteLibraryFolder(context.Background(), folder.ID) })
	return folder
}

func createImportPathMapping(t *testing.T, store *storage.SettingsStore, clientPath string, appPath string) {
	t.Helper()
	mapping, err := store.CreatePathMapping(context.Background(), storage.PathMappingInput{
		ClientPath: clientPath,
		AppPath:    appPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.DeletePathMapping(context.Background(), mapping.ID) })
}

func int32Ptr(value int32) *int32 {
	return &value
}

func createImportMediaItem(t *testing.T, store *storage.SettingsStore, input storage.MediaItemInput) storage.MediaItem {
	t.Helper()
	item, err := store.CreateMediaItem(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.DeleteMediaItem(context.Background(), item.ID, false) })
	return item
}
