package metadata

import "testing"

func TestMergeMediaResourcesOnlyFillsMissingFields(t *testing.T) {
	primaryOverview, fallbackOverview := "AniDB overview", "TMDB overview"
	fallbackBackdrop := "/backdrop.jpg"
	primary := MediaResource{
		Title: "Primary title", Type: "serie", ExternalProvider: "anidb", ExternalID: "9541",
		Overview: &primaryOverview, Genres: []string{"Anime"},
		ProviderMappings: []ProviderMapping{{ProviderName: "tmdb", ProviderEntityType: "serie", ExternalID: "1429"}},
	}
	fallback := MediaResource{
		Title: "Fallback title", Type: "serie", ExternalProvider: "tmdb", ExternalID: "1429",
		Overview: &fallbackOverview, BackdropPath: &fallbackBackdrop,
		Genres: []string{"Animation"}, Keywords: []string{"giant monster"},
		Seasons: []Season{{Name: "Season 1", SeasonNumber: 1}},
	}

	got := MergeMediaResources(primary, fallback)
	if got.Title != primary.Title || got.ExternalProvider != "anidb" || got.ExternalID != "9541" || got.Overview != primary.Overview {
		t.Fatalf("primary fields changed: %#v", got)
	}
	if got.BackdropPath == nil || *got.BackdropPath != fallbackBackdrop || len(got.Keywords) != 1 || len(got.Seasons) != 1 {
		t.Fatalf("missing fields were not filled: %#v", got)
	}
	if len(got.Genres) != 1 || got.Genres[0] != "Anime" {
		t.Fatalf("primary genres replaced: %#v", got.Genres)
	}
	if len(got.ProviderMappings) != 1 || got.ProviderMappings[0].Canonical {
		t.Fatalf("fallback mapping was duplicated or canonical: %#v", got.ProviderMappings)
	}
}
