package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSCNMedia011StreamURLUsesForwardedHostAndMountedPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://internal/api/media/items/abc/files/vlc", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	request.Header.Set("X-Forwarded-Host", "media.example.test")

	got := streamURL(request, "Season 01/Episode 01.mkv", 1790000000, "signed-token")
	want := "https://media.example.test/api/media/items/abc/files/stream?path=Season+01%2FEpisode+01.mkv&streamExpires=1790000000&streamToken=signed-token"
	if got != want {
		t.Fatalf("streamURL = %q, want %q", got, want)
	}
}

func TestSCNMedia011StreamTokenBindsMediaPathAndExpiry(t *testing.T) {
	server := &Server{
		streamSecret: []byte("test-stream-secret"),
		now:          func() time.Time { return time.Unix(1000, 0) },
	}
	mediaID := uuid.New()
	path := "Season 01/Episode 01.mkv"
	expires := server.now().Add(time.Hour).Unix()
	token := server.newStreamToken(mediaID, path, expires)

	if !server.validStreamToken(mediaID, path, &expires, &token) {
		t.Fatal("expected stream token to validate for matching media, path, and expiry")
	}
	otherPath := "Season 01/Episode 02.mkv"
	if server.validStreamToken(mediaID, otherPath, &expires, &token) {
		t.Fatal("expected stream token to reject a different path")
	}
	expired := server.now().Add(-time.Second).Unix()
	if server.validStreamToken(mediaID, path, &expired, &token) {
		t.Fatal("expected stream token to reject expired links")
	}
}

func TestSCNMedia011PlaylistFilenameKeepsVLCFriendlyExtension(t *testing.T) {
	got := playlistFilename("Movie.Name.2026.mkv")
	if got != "Movie.Name.2026.m3u" {
		t.Fatalf("playlistFilename = %q", got)
	}
}

func TestSCNMedia011PlaylistDispositionOpensInline(t *testing.T) {
	got := playlistDisposition("Movie.Name.2026.mkv")
	want := "inline; filename=Movie.Name.2026.m3u"
	if got != want {
		t.Fatalf("playlistDisposition = %q, want %q", got, want)
	}
}
