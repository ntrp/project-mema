package indexers

import (
	"reflect"
	"testing"
)

func TestDefaultMediaTypeScopesFromCatalogCategories(t *testing.T) {
	entry := CatalogEntry{Capabilities: Capabilities{Categories: []Category{
		{Name: "Movies"},
		{Name: "TV Shows"},
		{Name: "Audio", Children: []Category{{Name: "Music"}}},
		{Name: "Books", Children: []Category{{Name: "Audiobook"}}},
		{Name: "Anime"},
	}}}

	got := DefaultMediaTypeScopes(entry)
	want := []string{"movie", "serie", "anime", "audio", "book"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scopes = %#v, want %#v", got, want)
	}
}

func TestDefaultMediaTypeScopesFallbackToAllTypes(t *testing.T) {
	got := DefaultMediaTypeScopes(CatalogEntry{})
	want := []string{"movie", "serie", "anime", "audio", "book"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scopes = %#v, want %#v", got, want)
	}
}

func TestDefaultMediaTypeScopesClassifiesAudiobooksAsBooks(t *testing.T) {
	got := DefaultMediaTypeScopes(CatalogEntry{Capabilities: Capabilities{Categories: []Category{
		{Name: "Audiobook"},
	}}})
	want := []string{"book"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("scopes = %#v, want %#v", got, want)
	}
}
