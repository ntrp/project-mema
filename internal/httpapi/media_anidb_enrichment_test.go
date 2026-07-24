package httpapi

import (
	"context"
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
	input := storage.MediaItemInput{
		Title: "Primary title", ExternalProvider: &primaryProvider, ExternalID: &primaryID,
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{Facts: []storage.MediaFact{{Label: "Primary", Value: "kept"}}},
	}
	absolute := int32(8)
	details := metadata.Details{
		ExternalProvider: "anidb", ExternalID: "123", Type: "serie",
		Aliases: []metadata.Alias{{Value: "Provider alias", Kind: "romaji"}},
		Facts:   []metadata.Fact{{Label: "AniDB episode 1002", Value: "type=2; number=1"}},
		EpisodeFacts: []metadata.EpisodeFact{
			{ExternalID: "1001", OriginalType: "1", OriginalNumber: 8, AbsoluteNumber: &absolute, Title: "Same title"},
			{ExternalID: "1002", OriginalType: "2", OriginalNumber: 1, Title: "Same title"},
			{ExternalID: "9999", OriginalType: "1", OriginalNumber: 9, AbsoluteNumber: &absolute, Title: "Same title"},
		},
	}

	got := appendAniDBArtifacts(item, input, details, false)
	if got.Title != "Primary title" || got.ExternalProvider == nil || *got.ExternalProvider != "tmdb" || *got.ExternalID != "77" {
		t.Fatalf("primary display identity changed: %#v", got)
	}
	if len(got.ProviderMappings) != 1 || got.ProviderMappings[0].Canonical {
		t.Fatalf("AniDB mapping = %#v", got.ProviderMappings)
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
	if len(got.Facts) != 2 {
		t.Fatalf("exceptional provider facts were not retained: %#v", got.Facts)
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
