package httpapi

import (
	"testing"

	"media-manager/internal/indexers"
	"media-manager/internal/storage"
)

func TestIndexerCatalogLanguagesUseConfiguredLanguageCodes(t *testing.T) {
	languages := []storage.Language{
		{Code: "EN", DisplayName: "English", Aliases: []string{"ENG"}},
		{Code: "DE", DisplayName: "German", Aliases: []string{"Deutsch", "DEU"}},
	}
	entries := []indexers.CatalogEntry{
		{DefinitionID: "english", Name: "English", Language: "en-US"},
		{DefinitionID: "german", Name: "German", Language: "de-DE"},
		{DefinitionID: "alias", Name: "Alias", Language: "Deutsch"},
	}

	response := indexerCatalogResponse(entries, languages)

	if got := response.Entries[0].Language; got != "EN" {
		t.Fatalf("english entry language = %q", got)
	}
	if got := response.Entries[1].Language; got != "DE" {
		t.Fatalf("german entry language = %q", got)
	}
	if got := response.Entries[2].Language; got != "DE" {
		t.Fatalf("alias entry language = %q", got)
	}
	assertStringsEqual(t, response.Languages, []string{"EN", "DE"})
}

func assertStringsEqual(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("values = %#v, want %#v", got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("values = %#v, want %#v", got, want)
		}
	}
}
