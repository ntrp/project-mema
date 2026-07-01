package metadata

import "testing"

func TestTMDBTrailerURLPrefersOfficialYouTubeTrailer(t *testing.T) {
	url := tmdbTrailerURL(tmdbVideos{Results: []tmdbVideo{
		{Key: "teaser", Site: "YouTube", Type: "Teaser", Official: true},
		{Key: "unofficial", Site: "YouTube", Type: "Trailer"},
		{Key: "official", Site: "YouTube", Type: "Trailer", Official: true},
	}})
	if url == nil || *url != "https://www.youtube.com/watch?v=official" {
		t.Fatalf("expected official YouTube trailer URL, got %v", url)
	}
}

func TestTMDBTrailerURLRequiresYouTubeTrailer(t *testing.T) {
	url := tmdbTrailerURL(tmdbVideos{Results: []tmdbVideo{
		{Key: "clip", Site: "YouTube", Type: "Clip", Official: true},
		{Key: "trailer", Site: "Vimeo", Type: "Trailer", Official: true},
	}})
	if url != nil {
		t.Fatalf("expected no trailer URL, got %v", *url)
	}
}
