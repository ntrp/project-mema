package storage

import "testing"

func TestIndexerScopesDefaultAllMediaTypes(t *testing.T) {
	indexers := []Indexer{
		{Name: "Default scopes"},
		{Name: "Series only", MediaTypeScopes: []string{"serie"}},
	}

	got := EligibleIndexers(indexers, MediaItem{Type: "movie", Title: "Scenario Movie"})

	if len(got) != 1 || got[0].Name != "Default scopes" {
		t.Fatalf("eligible indexers = %#v", got)
	}
}

func TestIndexerScopesMatchTaggedMediaOnlyToMatchingTaggedIndexers(t *testing.T) {
	indexers := []Indexer{
		{Name: "Movies all tags", MediaTypeScopes: []string{"movie"}},
		{Name: "Movies kids", MediaTypeScopes: []string{"movie"}, TagScopes: []string{"Kids"}},
		{Name: "Movies anime", MediaTypeScopes: []string{"movie"}, TagScopes: []string{"Anime"}},
		{Name: "Series kids", MediaTypeScopes: []string{"serie"}, TagScopes: []string{"Kids"}},
	}

	got := EligibleIndexers(indexers, MediaItem{
		Type:  "movie",
		Title: "Tagged Movie",
		Tags:  []string{"kids"},
	})

	if len(got) != 1 || got[0].Name != "Movies kids" {
		t.Fatalf("eligible indexers = %#v", got)
	}
}

func TestIndexerScopesMatchUntaggedMediaToTaggedAndUntaggedIndexers(t *testing.T) {
	indexers := []Indexer{
		{Name: "Movies all tags", MediaTypeScopes: []string{"movie"}},
		{Name: "Movies kids", MediaTypeScopes: []string{"movie"}, TagScopes: []string{"Kids"}},
		{Name: "Series kids", MediaTypeScopes: []string{"serie"}, TagScopes: []string{"Kids"}},
	}

	got := EligibleIndexers(indexers, MediaItem{Type: "movie", Title: "Untagged Movie"})

	if len(got) != 2 || got[0].Name != "Movies all tags" || got[1].Name != "Movies kids" {
		t.Fatalf("eligible indexers = %#v", got)
	}
}
