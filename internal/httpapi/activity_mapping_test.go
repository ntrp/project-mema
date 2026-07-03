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

	failureType := "import"
	if !manualImportAllowed(storage.DownloadActivity{Status: "failed", FailureType: &failureType}) {
		t.Fatal("failed import activity should allow manual import")
	}
	if manualImportAllowed(storage.DownloadActivity{Status: "failed"}) {
		t.Fatal("failed activity without import failure type should not allow manual import")
	}

	languages := []string{" en ", "", " de "}
	input := manualImportInput(ManualImportRequest{
		SourcePath:     "/downloads/movie.mkv",
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
	if manualString(nil) != "" || manualStrings(nil) != nil {
		t.Fatal("nil manual fields should map to empty values")
	}
}
