package content

import "testing"

func TestApplyArtworkURLsAddsStableThumbnailRoutes(t *testing.T) {
	objects := []Object{{ID: "item-1", Kind: ObjectItem, Title: "Scenario", FilePath: "/media/movie.mkv"}}

	updated := ApplyArtworkURLs("http://127.0.0.1:18080", objects)

	if updated[0].Artwork == nil || *updated[0].Artwork != "http://127.0.0.1:18080/dlna/artwork/item-1?kind=thumbnail" {
		t.Fatalf("artwork = %#v", updated[0].Artwork)
	}
	if objects[0].Artwork != nil {
		t.Fatalf("ApplyArtworkURLs mutated input: %#v", objects[0])
	}
}
