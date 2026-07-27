package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSCNMedia011StreamURLUsesFileIdentifier(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://internal/api/media/items/abc/files/vlc", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	request.Header.Set("X-Forwarded-Host", "media.example.test")
	fileID := uuid.MustParse("10000000-0000-4000-8000-000000000001")

	got := streamURL(request, fileID, 1790000000, "signed-token")
	want := "https://media.example.test/api/media/items/abc/files/stream?fileId=10000000-0000-4000-8000-000000000001&streamExpires=1790000000&streamToken=signed-token"
	if got != want {
		t.Fatalf("streamURL = %q, want %q", got, want)
	}
}

func TestSCNMedia011StreamTokenBindsMediaFileIdentifierAndExpiry(t *testing.T) {
	server := &Server{
		streamSecret: []byte("test-stream-secret"),
		now:          func() time.Time { return time.Unix(1000, 0) },
	}
	mediaID, fileID := uuid.New(), uuid.New()
	expires := server.now().Add(time.Hour).Unix()
	token := server.newStreamToken(mediaID, fileID.String(), expires)

	if !server.validStreamToken(mediaID, fileID.String(), &expires, &token) {
		t.Fatal("expected stream token to validate for matching media, file, and expiry")
	}
	if server.validStreamToken(mediaID, uuid.New().String(), &expires, &token) {
		t.Fatal("expected stream token to reject a different file")
	}
	expired := server.now().Add(-time.Second).Unix()
	if server.validStreamToken(mediaID, fileID.String(), &expired, &token) {
		t.Fatal("expected stream token to reject expired links")
	}
}

func TestSCNMedia011PlaylistFilenameKeepsVLCFriendlyExtension(t *testing.T) {
	got := playlistFilename("Movie.Name.2026.mkv")
	if got != "Movie.Name.2026.m3u" {
		t.Fatalf("playlistFilename = %q", got)
	}
}

func TestSCNMedia011PlaylistDispositionDownloadsAuthenticatedPlaylist(t *testing.T) {
	got := playlistDisposition("Movie.Name.2026.mkv")
	want := "attachment; filename=Movie.Name.2026.m3u"
	if got != want {
		t.Fatalf("playlistDisposition = %q, want %q", got, want)
	}
}
