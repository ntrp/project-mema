package httpapi

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func TestAppendAniDBArtifactsPreservesPrimaryAndAttachesOnlyResolvedRegularEpisode(t *testing.T) {
	seasonID, episodeID := uuid.New(), uuid.New()
	item := storage.MediaItem{ProviderMappings: []storage.MediaProviderMapping{
		{ProviderName: "anidb", EntityType: "episode", ExternalID: "1001", SeasonID: &seasonID, EpisodeID: &episodeID},
	}}
	primaryProvider, primaryID := "tmdb", "77"
	primaryPoster, primaryDate := "tmdb.jpg", "2024-01-01"
	input := storage.MediaItemInput{
		Title: "Primary title", ExternalProvider: &primaryProvider, ExternalID: &primaryID, PosterPath: &primaryPoster,
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{
			ReleaseDate: &primaryDate,
			Facts:       []storage.MediaFact{{Label: "Primary", Value: "kept"}},
		},
	}
	absolute := int32(8)
	auxPoster, auxDate := "https://cdn-eu.anidb.net/images/main/456.jpg", "2024-02-03"
	details := metadata.Details{
		ExternalProvider: "anidb", ExternalID: "123", Type: "serie", PosterPath: &auxPoster, ReleaseDate: &auxDate,
		Aliases: []metadata.Alias{{Value: "Provider alias", Kind: "romaji"}},
		Facts:   []metadata.Fact{{Label: "AniDB episode 1002", Value: "type=2; number=1"}},
		EpisodeFacts: []metadata.EpisodeFact{
			{ExternalID: "1001", OriginalType: "1", OriginalNumber: 8, AbsoluteNumber: &absolute, Title: "Same title"},
			{ExternalID: "1002", OriginalType: "2", OriginalNumber: 1, Title: "Same title"},
			{ExternalID: "9999", OriginalType: "1", OriginalNumber: 9, AbsoluteNumber: &absolute, Title: "Same title"},
		},
	}

	got := appendAniDBArtifacts(item, input, details, false)
	if got.Title != "Primary title" || got.ExternalProvider == nil || *got.ExternalProvider != "tmdb" || *got.ExternalID != "77" || *got.PosterPath != primaryPoster || *got.ReleaseDate != primaryDate {
		t.Fatalf("primary display identity changed: %#v", got)
	}
	if len(got.ProviderMappings) != 1 || got.ProviderMappings[0].Canonical || got.ProviderMappings[0].Source["artwork"] != auxPoster || got.ProviderMappings[0].Source["startDate"] != auxDate {
		t.Fatalf("AniDB mapping provenance = %#v", got.ProviderMappings)
	}
	if len(got.Aliases) != 1 || got.Aliases[0].ProviderName == nil || *got.Aliases[0].ProviderName != "anidb" {
		t.Fatalf("aliases = %#v", got.Aliases)
	}
	if len(got.EpisodeNumbering) != 1 || got.EpisodeNumbering[0].EpisodeID != episodeID || *got.EpisodeNumbering[0].AbsoluteNumber != 8 {
		t.Fatalf("resolved numbering = %#v", got.EpisodeNumbering)
	}
	if got.EpisodeNumbering[0].Source["originalType"] != "1" || got.EpisodeNumbering[0].Source["resolvedBy"] != "provider_mapping" {
		t.Fatalf("numbering provenance = %#v", got.EpisodeNumbering[0].Source)
	}
	if len(got.Facts) != 4 || got.Facts[2].Label != "AniDB artwork" || got.Facts[3].Label != "AniDB start date" || got.Facts[3].Source["externalId"] != "123" {
		t.Fatalf("auxiliary provider facts were not retained with provenance: %#v", got.Facts)
	}
}

func TestDirectAniDBDetailsMarkMediaAsAnime(t *testing.T) {
	provider, externalID := "anidb", "123"
	got := applyMediaResource(storage.MediaItemInput{Type: "serie", ExternalProvider: &provider, ExternalID: &externalID}, metadata.Details{
		Title: "Anime", Type: "serie", ExternalProvider: "anidb", ExternalID: "123",
	})
	if got.ContentKind != "anime" {
		t.Fatalf("content kind = %q", got.ContentKind)
	}
}

func TestEnrichAniDBFactsUsesPrimaryDetailsWithoutChangingDisplay(t *testing.T) {
	provider, externalID := "anidb", "123"
	input := storage.MediaItemInput{Title: "Display", ExternalProvider: &provider, ExternalID: &externalID}
	got, err := (&Server{}).enrichAniDBFacts(context.Background(), storage.MediaItem{}, nil, input, metadata.Details{
		ExternalProvider: "anidb", ExternalID: "123", Type: "serie",
	})
	if err != nil || got.Title != "Display" || len(got.ProviderMappings) != 1 || !got.ProviderMappings[0].Canonical {
		t.Fatalf("input=%#v err=%v", got, err)
	}
}

func TestMetadataFallbackMustBeEnabledAndOnlyFillsMissingFields(t *testing.T) {
	apiKey, poster, backdrop := "tmdb-key", "/anidb.jpg", "/tmdb-backdrop.jpg"
	year := int32(2013)
	primary := metadata.MediaResource{
		Title: "Shingeki no Kyojin", Type: "serie", Year: &year, ExternalProvider: "anidb", ExternalID: "9541",
		PosterPath:       &poster,
		ProviderMappings: []metadata.ProviderMapping{{ProviderName: "tmdb", ExternalID: "1429"}},
	}
	loads := 0
	load := func(_ context.Context, _ storage.MetadataProvider, request metadata.DetailsRequest) (metadata.MediaResource, error) {
		loads++
		if request.ExternalID != "1429" {
			t.Fatalf("TMDB request = %#v", request)
		}
		return metadata.MediaResource{
			Title: "Attack on Titan", Type: "serie", ExternalProvider: "tmdb", ExternalID: "1429",
			BackdropPath: &backdrop, Genres: []string{"Animation"},
		}, nil
	}
	search := func(_ context.Context, _ storage.MetadataProvider, request metadata.SearchRequest) ([]metadata.SearchResult, error) {
		return []metadata.SearchResult{{Title: request.Query, Year: request.Year, ExternalID: "1429"}}, nil
	}
	provider := storage.MetadataProvider{Type: "tmdb", Enabled: true, APIKey: &apiKey}
	server := &Server{}
	if got := server.resourceWithMetadataFallback(context.Background(), []storage.MetadataProvider{provider}, primary, load, search); got.BackdropPath != nil || loads != 0 {
		t.Fatalf("disabled fallback changed resource: %#v loads=%d", got, loads)
	}
	provider.FallbackEnabled = true
	got := server.resourceWithMetadataFallback(context.Background(), []storage.MetadataProvider{provider}, primary, load, search)
	if loads != 1 || got.Title != primary.Title || got.PosterPath != primary.PosterPath || got.BackdropPath == nil || *got.BackdropPath != backdrop {
		t.Fatalf("fallback resource = %#v loads=%d", got, loads)
	}
	unmapped := primary
	unmapped.ProviderMappings = nil
	got = server.resourceWithMetadataFallback(context.Background(), []storage.MetadataProvider{provider}, unmapped, load, search)
	if loads != 2 || got.BackdropPath == nil || *got.BackdropPath != backdrop {
		t.Fatalf("searched fallback resource = %#v loads=%d", got, loads)
	}
	failed := server.resourceWithMetadataFallback(
		context.Background(), []storage.MetadataProvider{provider}, primary,
		func(context.Context, storage.MetadataProvider, metadata.DetailsRequest) (metadata.MediaResource, error) {
			return metadata.MediaResource{}, errors.New("fallback unavailable")
		}, search,
	)
	if failed.Title != primary.Title || failed.BackdropPath != nil {
		t.Fatalf("failed fallback changed primary: %#v", failed)
	}
}

func TestAniDBMediaMappingRequiresResolvedItemMapping(t *testing.T) {
	if _, ok := anidbMediaMapping(nil); ok {
		t.Fatal("empty mappings resolved")
	}
	input := storage.MediaItemInput{Title: "unchanged"}
	got, err := (&Server{}).enrichAniDBFacts(context.Background(), storage.MediaItem{}, nil, input, metadata.Details{ExternalProvider: "tmdb"})
	if err != nil || got.Title != "unchanged" {
		t.Fatalf("unmapped enrichment input=%#v err=%v", got, err)
	}
	id, ok := anidbMediaMapping([]storage.MediaProviderMapping{{ProviderName: "anidb", EntityType: "media_item", ExternalID: " 123 "}})
	if !ok || id != "123" {
		t.Fatalf("id=%q ok=%v", id, ok)
	}
}

func TestAniDBProviderCredentialsRequireClientNameAndVersion(t *testing.T) {
	name, version := "client", "1"
	if providerHasCredentials(storage.MetadataProvider{Type: "anidb", APIKey: &name}) {
		t.Fatal("AniDB provider enabled without client version")
	}
	if !providerHasCredentials(storage.MetadataProvider{Type: "anidb", APIKey: &name, PIN: &version}) {
		t.Fatal("configured AniDB identity was not enabled")
	}
}

func TestResolvedAniDBNumberingDoesNotGuessByTitle(t *testing.T) {
	absolute := int32(1)
	facts := []metadata.EpisodeFact{{ExternalID: "1001", OriginalType: "1", OriginalNumber: 1, AbsoluteNumber: &absolute, Title: "Matching title"}}
	if got := resolvedAniDBNumbering(nil, facts); len(got) != 0 {
		t.Fatalf("title-only facts attached: %#v", got)
	}
}
