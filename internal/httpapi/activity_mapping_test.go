package httpapi

import (
	"testing"

	"media-manager/internal/storage"
)

func TestSCNMedia009ActivityStateRulesAndManualImportInput(t *testing.T) {
	for _, status := range []string{"queued", "grabbed", "downloading"} {
		if !downloadActivityCancellable(status) {
			t.Fatalf("%s should be cancellable", status)
		}
	}
	if downloadActivityCancellable("completed") {
		t.Fatal("completed activity should not be cancellable")
	}
	if !downloadActivityDeletable("failed") || !downloadActivityDeletable("cancelled") {
		t.Fatal("failed and cancelled activity should be deletable")
	}
	if downloadActivityDeletable("queued") {
		t.Fatal("queued activity should not be deletable")
	}
	if !downloadActivityRetryable("failed") || !downloadActivityRetryable("cancelled") || downloadActivityRetryable("queued") {
		t.Fatal("only failed and cancelled activity should be retryable")
	}
	response := downloadActivityResponse(storage.DownloadActivity{ClientFiles: []storage.DownloadActivityClientFile{
		{Path: "/downloads/movie.mkv", SizeBytes: 2048, Complete: true},
	}})
	if response.ClientFiles == nil || len(*response.ClientFiles) != 1 || (*response.ClientFiles)[0].Path != "/downloads/movie.mkv" {
		t.Fatalf("download activity client files = %#v", response.ClientFiles)
	}

	failureType := "import"
	if !manualImportAllowed(storage.DownloadActivity{Status: "failed", FailureType: &failureType}) {
		t.Fatal("failed import activity should allow manual import")
	}
	if manualImportAllowed(storage.DownloadActivity{Status: "failed"}) {
		t.Fatal("failed activity without import failure type should not allow manual import")
	}

	languages := []string{" en ", "", " de "}
	mode := Move
	input := manualImportInput(ManualImportRequest{
		SourcePath:     "/downloads/movie.mkv",
		ImportMode:     &mode,
		TargetFileName: stringPtr(" Movie.mkv "),
		MovieTitle:     stringPtr(" Scenario Movie "),
		Languages:      &languages,
	})
	if input.SourcePath != "/downloads/movie.mkv" || input.TargetFileName != "Movie.mkv" {
		t.Fatalf("manual import paths = %#v", input)
	}
	if input.MovieTitle != "Scenario Movie" || len(input.Languages) != 2 || input.Languages[0] != "en" || input.Languages[1] != "de" {
		t.Fatalf("manual import metadata = %#v", input)
	}
	if input.ImportMode != "move" {
		t.Fatalf("manual import mode = %q", input.ImportMode)
	}
	if manualString(nil) != "" || manualStrings(nil) != nil {
		t.Fatal("nil manual fields should map to empty values")
	}
}
