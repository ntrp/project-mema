package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestMediaComponentSourcesRetainHydrateAndRelease(t *testing.T) {
	ctx, store := testDBStore(t)
	item, finalPath := componentSourceMediaItem(t, ctx, store)
	basePath := filepath.Join(*item.MediaFolderPath, "Base.Source.mkv")
	audioPath := filepath.Join(*item.MediaFolderPath, "Audio.Source.flac")
	writeTestFile(t, basePath)
	writeTestFile(t, audioPath)

	base, err := store.RetainMediaComponentSource(ctx, item.ID, MediaComponentSourceInput{
		SourceRole:      "baseVideo",
		SourceFilePath:  basePath,
		ReleaseTitle:    stringPtr("Base Release"),
		SourceMetadata:  stringPtr("indexer=local"),
		StreamInventory: "video:h264",
		Checksum:        stringPtr("sha256:base"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if base.SourceRole != "baseVideo" || base.SizeBytes == nil || *base.SizeBytes != 5 {
		t.Fatalf("unexpected base source = %#v", base)
	}
	audio, err := store.RetainMediaComponentSource(ctx, item.ID, MediaComponentSourceInput{
		SourceRole:      "audio",
		SourceFilePath:  audioPath,
		StreamInventory: "audio:flac",
	})
	if err != nil {
		t.Fatal(err)
	}
	requireFileExists(t, base.RetainedPath)
	requireFileExists(t, audio.RetainedPath)
	requireFileExists(t, basePath)
	requireFileExists(t, finalPath)

	loaded, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.ComponentSources) != 2 {
		t.Fatalf("component sources = %#v", loaded.ComponentSources)
	}
	if !componentSourceListHas(loaded.ComponentSources, base.ID) ||
		!componentSourceListHas(loaded.ComponentSources, audio.ID) {
		t.Fatalf("component sources missing retained records: %#v", loaded.ComponentSources)
	}

	if _, err := store.UpdateFileDeleteSettings(ctx, FileDeleteSettingsInput{Mode: FileDeleteModePermanent}); err != nil {
		t.Fatal(err)
	}
	released, err := store.ReleaseMediaComponentSource(ctx, item.ID, base.ID)
	if err != nil {
		t.Fatal(err)
	}
	if released.RetentionState != "released" || released.ReleasedAt == nil {
		t.Fatalf("released source = %#v", released)
	}
	if _, err := os.Stat(base.RetainedPath); !os.IsNotExist(err) {
		t.Fatalf("retained source stat err = %v, want missing", err)
	}
	requireFileExists(t, basePath)
	requireFileExists(t, finalPath)
}

func TestMediaComponentSourceRetainRejectsUnsafePath(t *testing.T) {
	ctx, store := testDBStore(t)
	item, _ := componentSourceMediaItem(t, ctx, store)

	_, err := store.RetainMediaComponentSource(ctx, item.ID, MediaComponentSourceInput{
		SourceRole:     "baseVideo",
		SourceFilePath: "../escape.mkv",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func componentSourceMediaItem(t *testing.T, ctx context.Context, store *SettingsStore) (MediaItem, string) {
	t.Helper()
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Component Movie " + uuid.NewString(),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.MediaFolderPath == nil {
		t.Fatal("expected media folder path")
	}
	finalPath := filepath.Join(*item.MediaFolderPath, "Component.Movie.mkv")
	writeTestFile(t, finalPath)
	if err := store.RecordImportedMediaFile(ctx, item, finalPath); err != nil {
		t.Fatal(err)
	}
	return item, finalPath
}

func requireFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s: %v", path, err)
	}
}

func componentSourceListHas(sources []MediaComponentSource, id uuid.UUID) bool {
	for _, source := range sources {
		if source.ID == id {
			return true
		}
	}
	return false
}
