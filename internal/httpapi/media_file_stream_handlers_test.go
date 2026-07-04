package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSCNMedia011StreamURLUsesForwardedHostAndMountedPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://internal/api/media/items/abc/files/vlc", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	request.Header.Set("X-Forwarded-Host", "media.example.test")

	got := streamURL(request, "Season 01/Episode 01.mkv")
	want := "https://media.example.test/api/media/items/abc/files/stream?path=Season+01%2FEpisode+01.mkv"
	if got != want {
		t.Fatalf("streamURL = %q, want %q", got, want)
	}
}

func TestSCNMedia011PlaylistFilenameKeepsVLCFriendlyExtension(t *testing.T) {
	got := playlistFilename("Movie.Name.2026.mkv")
	if got != "Movie.Name.2026.m3u" {
		t.Fatalf("playlistFilename = %q", got)
	}
}
