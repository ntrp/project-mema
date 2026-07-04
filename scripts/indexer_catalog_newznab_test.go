package main

import "testing"

func TestNewznabEntriesFromSource(t *testing.T) {
	source := `
yield return GetDefinition("NZBgeek", GetSettings("https://api.nzbgeek.info"), categories: new[] { 1000, 2000, 5000 });
yield return GetDefinition("Tabula Rasa", GetSettings("https://www.tabula-rasa.pw", apiPath: @"/api/v1/api"), categories: new[] { 1000, 7000 });
yield return GetDefinition("Generic Newznab", GetSettings(""));
`
	entries := newznabEntriesFromSource(source)
	if len(entries) != 2 {
		t.Fatalf("entries = %d", len(entries))
	}
	if entries[0].DefinitionID != "newznab-nzbgeek" || entries[0].IndexerURLs[0] != "https://api.nzbgeek.info/api" {
		t.Fatalf("first entry = %#v", entries[0])
	}
	if entries[1].DefinitionID != "newznab-tabula-rasa" || entries[1].IndexerURLs[0] != "https://www.tabula-rasa.pw/api/v1/api" {
		t.Fatalf("second entry = %#v", entries[1])
	}
}
