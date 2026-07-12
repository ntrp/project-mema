package providercore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/storage"
)

func TestBuildSearchRequestUsesMediaContext(t *testing.T) {
	item := storage.MediaItem{Type: "movie", Title: "Request Movie"}
	request := BuildSearchRequest(item, " english ", "/media/movie.mkv")
	if request.Title != "Request Movie" || request.LanguageID != "english" || request.MediaContext.File.Name != "movie.mkv" {
		t.Fatalf("request = %#v", request)
	}
}

func TestBuildMediaContextUsesStoredIDsAliasesNumberingAndFileFacts(t *testing.T) {
	provider := "tmdb"
	externalID := "123"
	size := int64(99)
	infoURL := "https://tracker.example/release"
	item := storage.MediaItem{
		ExternalProvider: &provider,
		ExternalID:       &externalID,
		ProviderMappings: []storage.MediaProviderMapping{
			{ProviderName: "imdb", ExternalID: "tt0000123", EntityType: "media"},
			{ProviderName: "tvdb", ExternalID: "456", EntityType: "episode"},
		},
		Aliases: []storage.MediaItemAlias{{Alias: "Alt Title", Kind: "localized"}},
		EpisodeNumbering: []storage.MediaEpisodeNumbering{{
			EpisodeID: uuid.New(), ProviderName: "anidb", NumberingScheme: "absolute", AbsoluteNumber: int32Ptr(7),
		}},
		FileFacts: []storage.MediaFileFact{{FilePath: "/media/show.mkv", SizeBytes: &size}},
		ComponentSources: []storage.MediaComponentSource{{
			SourceRole: "download", ReleaseID: &infoURL,
		}},
	}
	ctx := BuildMediaContext(item, "/media/show.mkv")
	if ctx.ExternalIDs["tmdb"] != "123" || ctx.ExternalIDs["imdb"] != "tt0000123" {
		t.Fatalf("external IDs = %#v", ctx.ExternalIDs)
	}
	if ctx.EpisodeExternalIDs["tvdb"] != "456" {
		t.Fatalf("episode IDs = %#v", ctx.EpisodeExternalIDs)
	}
	if len(ctx.Aliases) != 1 || ctx.Aliases[0].Value != "Alt Title" {
		t.Fatalf("aliases = %#v", ctx.Aliases)
	}
	if len(ctx.EpisodeNumbering) != 1 || *ctx.EpisodeNumbering[0].AbsoluteNumber != 7 {
		t.Fatalf("numbering = %#v", ctx.EpisodeNumbering)
	}
	if ctx.File.Name != "show.mkv" || ctx.File.Extension != "mkv" || ctx.File.SizeBytes != 99 {
		t.Fatalf("file context = %#v", ctx.File)
	}
	if len(ctx.Provenance) != 1 || ctx.Provenance[0].InfoURL != infoURL {
		t.Fatalf("provenance = %#v", ctx.Provenance)
	}
}

func TestBuildMediaContextDoesNotHashExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "movie.mkv")
	data := make([]byte, 140*1024)
	for index := range data {
		data[index] = byte(index)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := BuildMediaContext(storage.MediaItem{}, path)
	if len(ctx.File.Hashes) != 0 {
		t.Fatalf("generic context should not hash files, got %#v", ctx.File.Hashes)
	}
	if ctx.File.SizeBytes != int64(len(data)) {
		t.Fatalf("size = %d", ctx.File.SizeBytes)
	}
}

func TestComputeFileHashesIsExplicitAndReportsErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "movie.mkv")
	data := make([]byte, 140*1024)
	for index := range data {
		data[index] = byte(index)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	hashes, err := ComputeFileHashes(path, "sha256", "opensubtitles")
	if err != nil {
		t.Fatalf("ComputeFileHashes failed: %v", err)
	}
	if hashes["sha256"] == "" || hashes["opensubtitles"] == "" {
		t.Fatalf("hashes = %#v", hashes)
	}
	if _, err := ComputeFileHashes(filepath.Join(t.TempDir(), "missing.mkv"), "sha256"); err == nil {
		t.Fatal("expected missing file error")
	}
	shortPath := filepath.Join(t.TempDir(), "short.mkv")
	if err := os.WriteFile(shortPath, []byte("short"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ComputeFileHashes(shortPath, "opensubtitles"); err == nil {
		t.Fatal("expected short file error")
	}
}

func TestFileSHA256ReportsReadError(t *testing.T) {
	if _, err := fileSHA256(t.TempDir()); err == nil {
		t.Fatal("expected directory read error")
	}
}

func int32Ptr(value int32) *int32 { return &value }
