package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRecordImportedMediaFilePersistsDBBackedFileFact(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	target := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.1080p-ARR.mkv")
	writeSizedFile(t, target, 17)

	if err := store.RecordImportedMediaFileWithHistory(ctx, item, "", target, "hardlink"); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.FileFacts) != 1 {
		t.Fatalf("file facts = %#v", loaded.FileFacts)
	}
	fact := loaded.FileFacts[0]
	if fact.FilePath != target || fact.SizeBytes == nil || *fact.SizeBytes != 17 || fact.SourceKind != "import" {
		t.Fatalf("file fact = %#v", fact)
	}
}

func TestUpsertMediaFileFactPersistsProbeTracksForSatisfaction(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	path := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.1080p-ARR.mkv")
	bitrate := int32(768)
	width := int32(1920)

	_, err := store.UpsertMediaFileFact(ctx, MediaFileFactInput{
		MediaItemID:         item.ID,
		FilePath:            path,
		ContainerFormat:     stringPtr("matroska"),
		ContainerFormatName: stringPtr("Matroska / WebM"),
		SourceKind:          "probe",
		Tracks: []MediaFileTrackFactInput{
			{StreamIndex: 0, TrackType: "video", Codec: stringPtr("h264"), Width: &width},
			{StreamIndex: 1, TrackType: "audio", LanguageID: stringPtr("english"), Codec: stringPtr("aac"), BitrateKbps: &bitrate},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	facts, err := store.ListMediaFileFacts(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(facts) != 1 || len(facts[0].Tracks) != 2 {
		t.Fatalf("facts = %#v", facts)
	}
	if facts[0].Tracks[1].LanguageID == nil || *facts[0].Tracks[1].LanguageID != "english" {
		t.Fatalf("audio track facts = %#v", facts[0].Tracks)
	}
}

func TestRescanMediaItemFilesUpdatesPersistedFileFact(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	path := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.mkv")
	writeSizedFile(t, path, 23)

	if _, err := store.RescanMediaItemFiles(ctx, item.ID); err != nil {
		t.Fatal(err)
	}

	facts, err := store.ListMediaFileFacts(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(facts) != 1 || facts[0].SizeBytes == nil || *facts[0].SizeBytes != 23 || facts[0].SourceKind != "rescan" {
		t.Fatalf("rescan facts = %#v", facts)
	}
}

func writeSizedFile(t *testing.T, path string, size int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, bytes.Repeat([]byte("x"), size), 0o644); err != nil {
		t.Fatal(err)
	}
}
