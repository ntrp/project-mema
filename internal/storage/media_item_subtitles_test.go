package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestMediaItemSubtitlesStoreProvenanceAndDeleteManagedFile(t *testing.T) {
	ctx, store := testDBStore(t)
	root := t.TempDir()
	folder, err := store.CreateLibraryFolder(ctx, root, "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Subtitle Movie " + uuid.NewString(),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.MediaFolderPath == nil {
		t.Fatal("expected media folder path")
	}
	subtitlePath := filepath.Join(*item.MediaFolderPath, "Subtitle Movie.english.srt")
	if err := os.MkdirAll(filepath.Dir(subtitlePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(subtitlePath, []byte("subtitle"), 0o644); err != nil {
		t.Fatal(err)
	}
	size := int64(8)
	record, err := store.UpsertMediaItemSubtitle(ctx, MediaItemSubtitleInput{
		MediaItemID:        item.ID,
		ProviderName:       "OpenSubtitles",
		LanguageID:         "english",
		Format:             "srt",
		FilePath:           subtitlePath,
		SourceURL:          stringPtr("https://example.test/subtitle"),
		SourceRef:          stringPtr("feature-1"),
		ReleaseName:        stringPtr("Subtitle.Movie.2026"),
		ProviderSubtitleID: stringPtr("123"),
		Checksum:           stringPtr("sha256:abc"),
		SizeBytes:          &size,
	})
	if err != nil {
		t.Fatal(err)
	}
	if record.ProviderSubtitleID == nil || *record.ProviderSubtitleID != "123" || record.SizeBytes == nil || *record.SizeBytes != size {
		t.Fatalf("expected provenance fields, got %#v", record)
	}
	item, err = store.DeleteMediaItemSubtitle(ctx, item.ID, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(subtitlePath); !os.IsNotExist(err) {
		t.Fatalf("expected subtitle file deleted, err=%v", err)
	}
	if len(item.ExternalSubtitles) != 0 {
		t.Fatalf("expected subtitle record deleted, got %#v", item.ExternalSubtitles)
	}
}

func TestMediaItemSubtitleDeleteRejectsTraversal(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Traversal Movie " + uuid.NewString(),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	record, err := store.UpsertMediaItemSubtitle(ctx, MediaItemSubtitleInput{
		MediaItemID:  item.ID,
		ProviderName: "OpenSubtitles",
		LanguageID:   "english",
		FilePath:     "../escape.srt",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.DeleteMediaItemSubtitle(ctx, item.ID, record.ID); err != ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestMediaItemSubtitleSelectionAndAssemblyArtifacts(t *testing.T) {
	ctx, store := testDBStore(t)
	folder, err := store.CreateLibraryFolder(ctx, t.TempDir(), "movie")
	if err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateMediaItem(ctx, MediaItemInput{
		Type:            "movie",
		Title:           "Selection Movie " + uuid.NewString(),
		Monitored:       true,
		LibraryFolderID: &folder.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.MediaFolderPath == nil {
		t.Fatal("expected media folder path")
	}
	size := int64(12)
	selected, err := store.UpsertMediaItemSubtitle(ctx, MediaItemSubtitleInput{
		MediaItemID:  item.ID,
		ProviderName: "OpenSubtitles",
		LanguageID:   "english",
		FilePath:     filepath.Join(*item.MediaFolderPath, "Selection.Movie.english.srt"),
		Checksum:     stringPtr("sha256:selected"),
		SizeBytes:    &size,
	})
	if err != nil {
		t.Fatal(err)
	}
	ignored, err := store.UpsertMediaItemSubtitle(ctx, MediaItemSubtitleInput{
		MediaItemID:  item.ID,
		ProviderName: "OpenSubtitles",
		LanguageID:   "german",
		FilePath:     filepath.Join(*item.MediaFolderPath, "Selection.Movie.german.srt"),
	})
	if err != nil {
		t.Fatal(err)
	}
	external, err := store.UpsertMediaItemSubtitle(ctx, MediaItemSubtitleInput{
		MediaItemID:  item.ID,
		ProviderName: "OpenSubtitles",
		LanguageID:   "french",
		FilePath:     filepath.Join(*item.MediaFolderPath, "Selection.Movie.french.srt"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !selected.Selected || selected.RetentionMode != SubtitleRetentionExternal {
		t.Fatalf("expected default active external subtitle, got %#v", selected)
	}
	item, err = store.UpdateMediaItemSubtitleSelection(ctx, item.ID, selected.ID, MediaItemSubtitleSelectionInput{
		Selected:      true,
		RetentionMode: SubtitleRetentionMux,
	})
	if err != nil {
		t.Fatal(err)
	}
	item, err = store.UpdateMediaItemSubtitleSelection(ctx, item.ID, ignored.ID, MediaItemSubtitleSelectionInput{
		Selected:      false,
		RetentionMode: SubtitleRetentionIgnore,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.UpdateMediaItemSubtitleSelection(ctx, item.ID, ignored.ID, MediaItemSubtitleSelectionInput{
		Selected:      true,
		RetentionMode: "invalid",
	}); err != ErrInvalidInput {
		t.Fatalf("expected invalid retention mode, got %v", err)
	}
	artifacts := SelectedSubtitleArtifacts(item)
	if len(artifacts) != 1 || artifacts[0].ID != selected.ID || artifacts[0].ID == external.ID {
		t.Fatalf("expected only selected subtitle artifact, got %#v", artifacts)
	}
	if artifacts[0].RetentionMode != SubtitleRetentionMux || artifacts[0].Checksum == nil {
		t.Fatalf("expected mux artifact with provenance, got %#v", artifacts[0])
	}
}
