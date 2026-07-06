package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecordImportedMediaFileRecordsParsedReleaseGroupProvenance(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	target := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.1080p-ARR.mkv")
	source := filepath.Join(t.TempDir(), "download.mkv")

	if err := store.RecordImportedMediaFileWithHistory(ctx, item, source, target, "hardlink"); err != nil {
		t.Fatal(err)
	}

	provenance, err := store.ListMediaComponentProvenance(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(provenance) != 1 {
		t.Fatalf("provenance = %#v", provenance)
	}
	record := provenance[0]
	if record.ComponentType != "container" || record.ReleaseGroup != "ARR" || record.ReleaseName != "Scenario.Movie.2026.1080p" {
		t.Fatalf("provenance = %#v", record)
	}
	if record.SourceFilePath == nil || *record.SourceFilePath != target {
		t.Fatalf("source file path = %#v", record.SourceFilePath)
	}
	if len(record.TransformationChain) != 1 || record.TransformationChain[0]["originalPath"] != source {
		t.Fatalf("transformation chain = %#v", record.TransformationChain)
	}
}

func TestRecordImportedMediaFileRecordsSidecars(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	target := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.1080p-ARR.mkv")
	poster := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.1080p-ARR-poster.jpg")
	subtitle := filepath.Join(*item.MediaFolderPath, "English.srt")
	for _, path := range []string{poster, subtitle} {
		if err := os.WriteFile(path, []byte("sidecar"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	if err := store.RecordImportedMediaFileWithHistory(ctx, item, "", target, "hardlink"); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !mediaItemHasSidecar(loaded.Sidecars, MediaSidecarMetadata, poster) {
		t.Fatalf("metadata sidecar not recorded: %#v", loaded.Sidecars)
	}
	if !mediaItemHasSidecar(loaded.Sidecars, MediaSidecarSubtitle, subtitle) {
		t.Fatalf("subtitle sidecar not recorded: %#v", loaded.Sidecars)
	}
	sidecar := mediaItemSidecar(loaded.Sidecars, MediaSidecarSubtitle, subtitle)
	if sidecar == nil || sidecar.LanguageID == nil || *sidecar.LanguageID != "english" {
		t.Fatalf("subtitle sidecar language not recorded: %#v", sidecar)
	}
	if len(loaded.ExternalSubtitles) != 1 || loaded.ExternalSubtitles[0].LanguageID != "english" || loaded.ExternalSubtitles[0].FilePath != subtitle {
		t.Fatalf("external subtitles = %#v", loaded.ExternalSubtitles)
	}
	if !stringSliceHas(loaded.MetadataFilePaths, poster) {
		t.Fatalf("metadata paths = %#v", loaded.MetadataFilePaths)
	}
}

func TestRecordImportedMediaFileMuxesSidecarSubtitlesForEmbeddedProfile(t *testing.T) {
	ctx, store := testDBStore(t)
	item := rescanMediaItem(t, ctx, store)
	item.SubtitlePreferredMode = "embedded"
	target := filepath.Join(*item.MediaFolderPath, "Scenario.Movie.2026.1080p-ARR.mkv")
	subtitle := filepath.Join(*item.MediaFolderPath, "English.srt")
	if err := os.WriteFile(subtitle, []byte("sidecar"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := store.RecordImportedMediaFileWithHistory(ctx, item, "", target, "hardlink"); err != nil {
		t.Fatal(err)
	}

	loaded, err := store.GetMediaItem(ctx, item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.ExternalSubtitles) != 1 {
		t.Fatalf("external subtitles = %#v", loaded.ExternalSubtitles)
	}
	if loaded.ExternalSubtitles[0].RetentionMode != SubtitleRetentionMux {
		t.Fatalf("expected mux retention for embedded profile, got %#v", loaded.ExternalSubtitles[0])
	}
}

func mediaItemHasSidecar(sidecars []MediaItemSidecar, sidecarType MediaSidecarType, path string) bool {
	return mediaItemSidecar(sidecars, sidecarType, path) != nil
}

func mediaItemSidecar(sidecars []MediaItemSidecar, sidecarType MediaSidecarType, path string) *MediaItemSidecar {
	for _, sidecar := range sidecars {
		if sidecar.SidecarType == sidecarType && sidecar.FilePath == path {
			return &sidecar
		}
	}
	return nil
}

func stringSliceHas(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
