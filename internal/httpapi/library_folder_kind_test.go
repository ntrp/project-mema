package httpapi

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestLibraryFolderKindAllowsMediaKind(t *testing.T) {
	cases := []struct {
		folderKind string
		mediaKind  string
		want       bool
	}{
		{"movie", "movie", true},
		{"movie", "series", false},
		{"series", "series", true},
		{"series", "anime_series", true},
		{"series", "movie", false},
	}
	for _, tc := range cases {
		if got := libraryFolderKindAllowsMediaKind(tc.folderKind, tc.mediaKind); got != tc.want {
			t.Fatalf("allow %s/%s = %v, want %v", tc.folderKind, tc.mediaKind, got, tc.want)
		}
	}
}

func TestValidLibraryEpisodeNumbersRejectsInvalidValues(t *testing.T) {
	negative, episode := int32(-1), int32(2)
	for _, values := range [][2]*int32{{&episode, nil}, {&negative, &episode}} {
		response := httptest.NewRecorder()
		if validLibraryEpisodeNumbers(response, values[0], values[1]) || response.Code != http.StatusBadRequest {
			t.Fatalf("values %v accepted with status %d", values, response.Code)
		}
	}
}

func TestLibraryScanItemPathUsesFolderRoot(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "media", "movies")
	relative := filepath.ToSlash(filepath.Join("Movie (2026)", "Movie.2026.mkv"))
	want := filepath.Join(root, "Movie (2026)", "Movie.2026.mkv")
	if got := libraryScanItemPath(root, relative); got != want {
		t.Fatalf("libraryScanItemPath relative = %q, want %q", got, want)
	}

	absolute := filepath.Join(root, "Existing.mkv")
	if got := libraryScanItemPath(root, absolute); got != absolute {
		t.Fatalf("libraryScanItemPath absolute = %q, want %q", got, absolute)
	}
}

func TestSameMediaTitleStripsAccentsAndPunctuation(t *testing.T) {
	cases := [][2]string{
		{"Amelie", "Amélie"},
		{"Walle", "WALL-E"},
		{"Wall e", "Wall-e"},
	}
	for _, tc := range cases {
		if !sameMediaTitle(tc[0], tc[1]) {
			t.Fatalf("sameMediaTitle(%q, %q) = false", tc[0], tc[1])
		}
	}
}
