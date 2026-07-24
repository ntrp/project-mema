package httpapi

import (
	"reflect"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func TestSCNMedia008ApplyMediaDetailsCopiesProviderSnapshot(t *testing.T) {
	year := int32(2026)
	runtime := int32(120)
	seasonCount := int32(1)
	episodeCount := int32(2)
	voteAverage := 8.4
	episodeTotal := int32(2)
	details := metadata.Details{
		Title:            "Scenario Series",
		Type:             "serie",
		Year:             &year,
		ExternalProvider: "tmdb",
		ExternalID:       "series-1",
		ExternalURL:      stringPointer("https://www.themoviedb.org/tv/series-1"),
		Overview:         stringPointer("Overview"),
		PosterPath:       stringPointer("/poster.jpg"),
		CollectionID:     stringPointer("collection-1"),
		CollectionName:   stringPointer("Scenario Collection"),
		BackdropPath:     stringPointer("/backdrop.jpg"),
		Status:           stringPointer("returning"),
		OriginalLanguage: stringPointer("de"),
		FirstAirDate:     stringPointer("2026-01-02"),
		RuntimeMinutes:   &runtime,
		SeasonCount:      &seasonCount,
		EpisodeCount:     &episodeCount,
		VoteAverage:      &voteAverage,
		Genres:           []string{"Drama"},
		Keywords:         []string{"scenario"},
		Facts:            []metadata.Fact{{Label: "Network", Value: "Local"}},
		Seasons: []metadata.Season{{
			Name:         "Season 1",
			EpisodeCount: &episodeTotal,
			AirDate:      stringPointer("2026-01-02"),
			PosterPath:   stringPointer("/season.jpg"),
			Episodes: []metadata.Episode{{
				Name:          "Pilot",
				EpisodeNumber: 1,
				Overview:      stringPointer("Episode overview"),
				AirDate:       stringPointer("2026-01-02"),
				StillPath:     stringPointer("/still.jpg"),
			}},
		}},
		Cast: []metadata.Person{{
			Name:        "Actor",
			Role:        stringPointer("Lead"),
			ProfilePath: stringPointer("/actor.jpg"),
		}},
		Crew: []metadata.Person{{
			ExternalProvider: stringPointer("tmdb"),
			ExternalID:       stringPointer("crew-1"),
			Name:             "Director",
			Role:             stringPointer("Director"),
			ProfilePath:      stringPointer("/director.jpg"),
		}},
		Recommendations: []metadata.SearchResult{{
			Title:            "Recommended",
			Type:             "movie",
			Year:             &year,
			ExternalProvider: "tmdb",
			ExternalID:       "rec-1",
			ExternalURL:      stringPointer("https://www.themoviedb.org/movie/rec-1"),
			Overview:         stringPointer("Recommendation"),
			PosterPath:       stringPointer("/rec.jpg"),
		}},
		Similar: []metadata.SearchResult{{
			Title:            "Similar",
			Type:             "serie",
			ExternalProvider: "tmdb",
			ExternalID:       "sim-1",
		}},
	}

	input := applyMediaDetails(storage.MediaItemInput{Title: "Old"}, details)

	if input.Title != "Scenario Series" || input.Type != "serie" || input.Year == nil || *input.Year != year {
		t.Fatalf("core fields = %#v", input)
	}
	if input.ExternalProvider == nil || *input.ExternalProvider != "tmdb" || input.ExternalID == nil || *input.ExternalID != "series-1" {
		t.Fatalf("external ids = provider %v id %v", input.ExternalProvider, input.ExternalID)
	}
	if len(input.ProviderMappings) != 1 || input.ProviderMappings[0].Source["externalUrl"] != "https://www.themoviedb.org/tv/series-1" {
		t.Fatalf("provider mappings = %#v", input.ProviderMappings)
	}
	if len(input.Genres) != 1 || input.Genres[0] != "Drama" || len(input.Keywords) != 1 {
		t.Fatalf("genres/keywords = %#v %#v", input.Genres, input.Keywords)
	}
	if len(input.Facts) != 1 || input.Facts[0].Label != "Network" {
		t.Fatalf("facts = %#v", input.Facts)
	}
	if len(input.Seasons) != 1 || len(input.Seasons[0].Episodes) != 1 {
		t.Fatalf("seasons = %#v", input.Seasons)
	}
	if input.Seasons[0].Episodes[0].Name != "Pilot" || input.Seasons[0].Episodes[0].EpisodeNumber != 1 {
		t.Fatalf("episode = %#v", input.Seasons[0].Episodes[0])
	}
	if len(input.Cast) != 1 || input.Cast[0].Name != "Actor" {
		t.Fatalf("cast = %#v", input.Cast)
	}
	if len(input.Crew) != 1 || input.Crew[0].ExternalID == nil || *input.Crew[0].ExternalID != "crew-1" {
		t.Fatalf("crew = %#v", input.Crew)
	}
	if len(input.Recommendations) != 1 || input.Recommendations[0].ExternalID != "rec-1" {
		t.Fatalf("recommendations = %#v", input.Recommendations)
	}
	if input.Recommendations[0].ExternalURL == nil || *input.Recommendations[0].ExternalURL != "https://www.themoviedb.org/movie/rec-1" {
		t.Fatalf("recommendation external url = %#v", input.Recommendations[0].ExternalURL)
	}
	if len(input.Similar) != 1 || input.Similar[0].ExternalID != "sim-1" {
		t.Fatalf("similar = %#v", input.Similar)
	}
}

func TestMetadataSearchResultResponseCopiesAliases(t *testing.T) {
	language := "ja-Latn"
	response := metadataSearchResultResponse(metadata.SearchResult{Title: "Frieren", Type: "serie", Aliases: []metadata.Alias{{
		Value: "Sousou no Frieren", Language: &language, Kind: "romaji", ProviderName: "anilist",
	}}})

	if response.Aliases == nil || len(*response.Aliases) != 1 {
		t.Fatalf("aliases = %#v", response.Aliases)
	}
	alias := (*response.Aliases)[0]
	if alias.Value != "Sousou no Frieren" || alias.Language == nil || *alias.Language != language || alias.Kind != "romaji" || alias.ProviderName == nil || *alias.ProviderName != "anilist" {
		t.Fatalf("alias = %#v", alias)
	}
}

func TestMergeMediaProviderMappingsPreservesExistingPrimaryProvenance(t *testing.T) {
	source := map[string]any{"source": "tmdb_import", "raw": map[string]any{"id": float64(209867)}}
	existing := []storage.MediaProviderMappingInput{{EntityType: "media_item", ProviderName: "tmdb", ExternalID: "209867", Canonical: true, Source: source}}
	next := []storage.MediaProviderMappingInput{{EntityType: "media_item", ProviderName: "anilist", ExternalID: "154587", Canonical: true}}

	mappings := mergeMediaProviderMappings(existing, next, "tmdb", "209867", "serie")

	if len(mappings) != 2 || !reflect.DeepEqual(mappings[0].Source, source) || mappings[1].Canonical {
		t.Fatalf("provider mappings = %#v", mappings)
	}
}

func TestApplyAniListDetailsPreservesExistingTMDBPrimary(t *testing.T) {
	tmdb := "tmdb"
	tmdbID := "209867"
	details := metadata.Details{
		Title: "Frieren: Beyond Journey's End", Type: "serie", ExternalProvider: "anilist", ExternalID: "154587",
		Aliases: []metadata.Alias{{Value: "Sousou no Frieren", Kind: "romaji", ProviderName: "anilist", Source: map[string]any{"field": "title.romaji"}}},
		ProviderMappings: []metadata.ProviderMapping{{ProviderName: "anilist", ProviderEntityType: "serie", ExternalID: "154587", Canonical: true,
			Source: map[string]any{"source": "anilist", "raw": map[string]any{"id": float64(154587)}}}},
	}

	input := applyMediaDetails(storage.MediaItemInput{Type: "serie", ExternalProvider: &tmdb, ExternalID: &tmdbID}, details)

	if input.ExternalProvider == nil || *input.ExternalProvider != "tmdb" || input.ExternalID == nil || *input.ExternalID != "209867" {
		t.Fatalf("primary provider changed: %#v %#v", input.ExternalProvider, input.ExternalID)
	}
	if input.ContentKind != "anime" || len(input.ProviderMappings) != 2 || !input.ProviderMappings[0].Canonical || input.ProviderMappings[1].Canonical {
		t.Fatalf("provider mappings = %#v", input.ProviderMappings)
	}
	if input.ProviderMappings[1].Source["raw"] == nil || len(input.Aliases) != 1 || input.Aliases[0].ProviderName == nil || *input.Aliases[0].ProviderName != "anilist" {
		t.Fatalf("provenance/aliases = %#v %#v", input.ProviderMappings, input.Aliases)
	}
}

func TestMediaRefreshUsesAniListMappingWithoutChangingPrimaryLink(t *testing.T) {
	item := storage.MediaItem{
		ExternalProvider: stringPointer("tmdb"), ExternalID: stringPointer("209867"),
		ProviderMappings: []storage.MediaProviderMapping{{EntityType: "media_item", ProviderName: "anilist", ExternalID: "154587"}},
	}
	provider, externalID := mediaRefreshProviderLink(item)
	if provider != "anilist" || externalID != "154587" || *item.ExternalProvider != "tmdb" || *item.ExternalID != "209867" {
		t.Fatalf("refresh link = %q %q, item = %#v", provider, externalID, item)
	}
}

func TestSCNMedia006MediaInputFromRequestCarriesApprovalChoices(t *testing.T) {
	year := int32(2026)
	folderID := uuid.New()
	request := storage.MediaRequest{
		Type:                "movie",
		Title:               "Requested Movie",
		Year:                &year,
		ExternalProvider:    stringPointer("tmdb"),
		ExternalID:          stringPointer("movie-1"),
		Overview:            stringPointer("Overview"),
		PosterPath:          stringPointer("/poster.jpg"),
		MonitorMode:         "movie",
		MinimumAvailability: "released",
		Tags:                []string{"favorite"},
	}

	input := mediaInputFromRequest(request, storage.MediaRequestApprovalInput{
		QualityProfileID:    "profile-1",
		LibraryFolderID:     folderID,
		MonitorMode:         "collection",
		MinimumAvailability: "in_cinema",
		Tags:                []string{"favorite"},
	})

	if input.Type != "movie" || input.Title != "Requested Movie" || !input.Monitored {
		t.Fatalf("core input = %#v", input)
	}
	if input.QualityProfileID == nil || *input.QualityProfileID != "profile-1" {
		t.Fatalf("quality profile = %v", input.QualityProfileID)
	}
	if input.LibraryFolderID == nil || *input.LibraryFolderID != folderID {
		t.Fatalf("library folder = %v", input.LibraryFolderID)
	}
	if input.MonitorMode != "collection" || input.MinimumAvailability != "in_cinema" {
		t.Fatalf("approval options = %#v", input)
	}
	if len(input.Tags) != 1 || input.Tags[0] != "favorite" {
		t.Fatalf("tags = %#v", input.Tags)
	}
}

func stringPointer(value string) *string {
	return &value
}
