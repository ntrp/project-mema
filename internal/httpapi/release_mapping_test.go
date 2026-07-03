package httpapi

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/decisions"
	"media-manager/internal/storage"
)

func TestSCNMedia002ParsedReleaseMetadataMapsCatalogLanguagesOnce(t *testing.T) {
	metadata := parsedReleaseMetadataResponse(
		decisions.ParsedRelease{
			ReleaseTitle:  "Scenario.Movie.2026.German.English.1080p.WEBDL-GRP",
			MovieTitle:    "Scenario Movie",
			Year:          "2026",
			QualityID:     "webdl-1080p",
			Quality:       "WEBDL-1080p",
			Source:        "WEB-DL",
			Resolution:    "1080p",
			VideoCodec:    "H.264",
			AudioCodec:    "AC3",
			AudioChannels: "5.1",
			Version:       "2",
			Proper:        true,
			Repack:        true,
			ReleaseGroup:  "GRP",
			Languages:     []string{" ger ", "German", "English", ""},
			ReleaseType:   "movie",
		},
		[]storage.Language{
			{Code: "de", DisplayName: "German", Aliases: []string{"ger", "deu"}},
			{Code: "en", DisplayName: "English", Aliases: []string{"eng"}},
		},
	)

	if metadata.Release.MovieTitle != "Scenario Movie" || metadata.Release.Year != "2026" {
		t.Fatalf("release metadata = %#v", metadata.Release)
	}
	if metadata.Quality.QualityId != "webdl-1080p" || !metadata.Quality.Proper || !metadata.Quality.Repack {
		t.Fatalf("quality metadata = %#v", metadata.Quality)
	}
	if got := metadata.Languages; len(got) != 2 || got[0] != "German" || got[1] != "English" {
		t.Fatalf("languages = %#v, want German and English once", got)
	}
	if metadata.Details.ReleaseType != "movie" || len(metadata.Details.CustomFormatNames) != 0 {
		t.Fatalf("details = %#v", metadata.Details)
	}
}

func TestSCNMedia002ReleaseCandidateResponsePreservesIndexerAndMatchDetails(t *testing.T) {
	indexerID := uuid.New()
	releaseID := uuid.New()
	infoURL := "https://indexer.test/release/scenario"
	guid := "scenario-guid"
	seeders := int32(42)
	peers := int32(7)
	publishedAt := time.Date(2026, time.July, 3, 4, 0, 0, 0, time.UTC)

	response := releaseCandidateResponse(
		storage.MediaItem{Type: "movie", Title: "Scenario Movie"},
		storage.ReleaseCandidate{
			ID:          releaseID,
			IndexerID:   &indexerID,
			IndexerName: "Local Torznab",
			IndexerType: "torznab",
			Title:       "Scenario.Movie.2026.1080p.WEBDL.German-GRP",
			InfoURL:     &infoURL,
			GUID:        &guid,
			SizeBytes:   8 * 1024 * 1024 * 1024,
			Seeders:     &seeders,
			Peers:       &peers,
			PublishedAt: &publishedAt,
		},
		&storage.MediaProfile{QualityIDs: []string{"webdl-1080p"}, TargetLanguages: []string{"de"}},
		nil,
		[]storage.Language{{Code: "de", DisplayName: "German", Aliases: []string{"ger", "german"}}},
	)

	if response.Id != releaseID || response.IndexerId == nil || uuid.UUID(*response.IndexerId) != indexerID {
		t.Fatalf("identity = id %s indexer %v", response.Id, response.IndexerId)
	}
	if response.IndexerName != "Local Torznab" || response.IndexerType != "torznab" {
		t.Fatalf("indexer fields = %#v", response)
	}
	if response.InfoUrl == nil || *response.InfoUrl != infoURL || response.Guid == nil || *response.Guid != guid {
		t.Fatalf("urls = info %v guid %v", response.InfoUrl, response.Guid)
	}
	if response.Match.Severity != "info" {
		t.Fatalf("match severity = %q details %v", response.Match.Severity, response.Match.Details)
	}
	if response.Match.QualityId != "webdl-1080p" || response.Match.Parsed.Quality.QualityId != "webdl-1080p" {
		t.Fatalf("quality match = %#v", response.Match)
	}
	if got := response.Match.Parsed.Languages; len(got) != 1 || got[0] != "German" {
		t.Fatalf("parsed languages = %#v, want German", got)
	}
	if len(response.Match.RankContributors) == 0 {
		t.Fatalf("expected rank contributors: %#v", response.Match)
	}
}

func TestSCNMedia002ReleaseScoreContributorResponsesPreserveOrder(t *testing.T) {
	responses := releaseScoreContributorResponses([]decisions.ReleaseScoreContributor{
		{Label: "Quality", Score: 100},
		{Label: "Seeders", Score: 42},
	})

	if len(responses) != 2 {
		t.Fatalf("responses = %#v", responses)
	}
	if responses[0].Label != "Quality" || responses[0].Score != 100 {
		t.Fatalf("first response = %#v", responses[0])
	}
	if responses[1].Label != "Seeders" || responses[1].Score != 42 {
		t.Fatalf("second response = %#v", responses[1])
	}
}

func TestSCNMedia002GrabReleaseOverrideControlsMismatchBlocking(t *testing.T) {
	item := storage.MediaItem{Type: "movie", Title: "Scenario Movie", Year: int32Ptr(2026)}
	release := storage.ReleaseCandidate{Title: "Different.Movie.2025.1080p-GRP"}

	if !shouldBlockReleaseMismatch(item, release, false) {
		t.Fatal("expected mismatched release to be blocked without override")
	}
	if shouldBlockReleaseMismatch(item, release, true) {
		t.Fatal("expected override to allow mismatched release")
	}

	value := true
	if !boolValue(&value) {
		t.Fatal("expected boolValue to return true for true pointer")
	}
	value = false
	if boolValue(&value) || boolValue(nil) {
		t.Fatal("expected boolValue to return false for false or nil")
	}
}
