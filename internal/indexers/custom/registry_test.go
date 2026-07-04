package custom

import (
	"context"
	"strings"
	"testing"

	"media-manager/internal/indexers/engine"
)

func TestRegistryRoutesGenericNabDefinitions(t *testing.T) {
	registry := NewRegistry(nil)
	for _, definitionID := range []string{"generic-newznab", "generic-torznab"} {
		if _, ok := registry.EngineFor(engine.Config{DefinitionID: definitionID}); !ok {
			t.Fatalf("missing engine for %s", definitionID)
		}
	}
}

func TestRegistryIncludesProwlarrSourceImplementations(t *testing.T) {
	registry := NewRegistry(nil)
	implemented := implementedProwlarrSourceImplementations()
	for _, implementation := range prowlarrSourceImplementations() {
		indexer, ok := registry.EngineFor(engine.Config{Implementation: implementation})
		if !ok {
			t.Fatalf("missing engine for %s", implementation)
		}
		if implemented[implementation] {
			continue
		}
		result := indexer.Test(context.Background(), engine.Config{Implementation: implementation})
		if !strings.Contains(result.Message, "not implemented") {
			t.Fatalf("%s returned unexpected result %#v", implementation, result)
		}
	}
}

func implementedProwlarrSourceImplementations() map[string]bool {
	implemented := map[string]bool{}
	for _, implementation := range prowlarrSourceImplementations() {
		implemented[implementation] = true
	}
	return implemented
}

func prowlarrSourceImplementations() []string {
	return []string{
		"AlphaRatio", "Anidex", "Anidub", "AnimeBytes", "AnimeZ", "Animedia",
		"AvistaZ", "BakaBT", "BeyondHD", "BinSearch", "BitHDTV", "BroadcastheNet",
		"BrokenStones", "CGPeers", "CinemaZ", "DICMusic", "ExoticaZ", "FileList",
		"FunFile", "Gazelle", "GazelleGames", "GreatPosterWall", "HDBits", "HDSpace",
		"HDTorrents", "Headphones", "IPTorrents", "ImmortalSeed", "Knaben", "Libble",
		"MTeamTp", "MyAnonamouse", "Nebulance", "NorBits", "NzbIndex", "Orpheus",
		"PassThePopcorn", "PixelHD", "PornoLab", "PreToMe", "PrivateHD", "Rarbg",
		"Redacted", "RetroFlix", "RevolutionTT", "RuTracker", "SceneHD", "SceneTime",
		"SecretCinema", "Shazbat", "Shizaproject", "SpeedApp", "SpeedCD", "SubsPlease",
		"Toloka", "TorrentBytes", "TorrentDay", "TorrentPotato", "TorrentRss",
		"TorrentSyndikat", "TorrentsCSV", "UNIT3D", "Uniotaku", "XSpeeds", "Xthor",
		"ZonaQ",
	}
}
