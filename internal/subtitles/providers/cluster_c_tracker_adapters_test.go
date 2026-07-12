package providers

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

func TestClusterCNativeReplacesGenericRegistrationForTargetProviders(t *testing.T) {
	for _, key := range []string{"addic7ed", "avistaz", "cinemaz", "hdbits"} {
		adapter, ok := AdapterFor(key)
		if !ok {
			t.Fatalf("%s not registered", key)
		}
		if _, generic := adapter.(clusterCProvider); generic {
			t.Fatalf("%s still uses generic cluster C adapter", key)
		}
	}
}

func TestAddic7edSearchUsesCookieCaptchaSessionAndParsesRows(t *testing.T) {
	body := `<table><tr class="epeven completed" data-language="English"><td>Show</td><td>Show.S01E02.720p-GRP</td><td>OK</td><td>English</td><td><a href="/updated/1/123/0">Download</a></td></tr></table>`
	svc := &nativeRoundTrip{body: []byte(body)}
	season, episode := int32(1), int32(2)
	got, err := addic7edAdapter().Search(context.Background(), svc, nativeCfg("PHPSESSID=abc"), providercore.SearchRequest{Title: "Show", SeasonNumber: &season, EpisodeNumber: &episode, LanguageID: "eng"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ProviderName != "addic7ed" || got[0].SourceURL != "https://www.addic7ed.com/updated/1/123/0" {
		t.Fatalf("unexpected candidates: %#v", got)
	}
	if svc.lastReq.URL.Path != "/search.php" || svc.lastReq.URL.Query().Get("Submit") != "Search" || !strings.Contains(svc.lastReq.URL.Query().Get("search"), "S01E02") {
		t.Fatalf("unexpected request: %s", svc.lastReq.URL.String())
	}
	if svc.lastReq.Header.Get("Cookie") != "PHPSESSID=abc" {
		t.Fatalf("cookie = %q", svc.lastReq.Header.Get("Cookie"))
	}
}

func TestPrivateTrackerSubtitlesRequireCookieAndReleaseProvenance(t *testing.T) {
	_, err := avistazSubtitleAdapter().Search(context.Background(), &nativeRoundTrip{}, providercore.Config{}, providercore.SearchRequest{Title: "Movie"})
	if !errors.Is(err, providercore.ErrPrivateMembershipRequired) {
		t.Fatalf("missing cookie error = %v", err)
	}
	_, err = avistazSubtitleAdapter().Search(context.Background(), &nativeRoundTrip{}, nativeCfg("session=1"), providercore.SearchRequest{Title: "Movie"})
	if !errors.Is(err, providercore.ErrReleaseProvenanceRequired) {
		t.Fatalf("missing provenance error = %v", err)
	}
}

func TestAvistaZAndCinemaZSearchUseTrackerProvenanceAndParseSubtitles(t *testing.T) {
	cases := []struct {
		name    string
		adapter nativeCProvider
		infoURL string
		base    string
	}{
		{"avistaz", avistazSubtitleAdapter(), "https://avistaz.to/torrent/987-movie", "https://avistaz.to"},
		{"cinemaz", cinemazSubtitleAdapter(), "https://cinemaz.to/torrent/654-movie", "https://cinemaz.to"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &nativeRoundTrip{body: []byte(`<table><tr data-release="Movie.2026.1080p" data-language="eng"><td class="name">Movie subtitle</td><td>English</td><td><a href="/subtitles/42/download">Download</a></td></tr></table>`)}
			got, err := tc.adapter.Search(context.Background(), svc, nativeCfg("uid=1"), trackerReq(tc.name, tc.infoURL))
			if err != nil {
				t.Fatal(err)
			}
			if svc.lastReq.URL.Path != "/subtitles" || svc.lastReq.URL.Query().Get("torrent_id") == "" {
				t.Fatalf("unexpected request: %s", svc.lastReq.URL.String())
			}
			wantURL := tc.base + "/subtitles/42/download"
			if len(got) != 1 || got[0].ProviderName != tc.name || got[0].SourceURL != wantURL || got[0].LanguageID != "eng" {
				t.Fatalf("unexpected candidates: %#v", got)
			}
		})
	}
}

func TestHDBitsSearchUsesDetailsProvenanceAndParsesSubtitleDownloads(t *testing.T) {
	svc := &nativeRoundTrip{body: []byte(`<table id="subs"><tr data-release="Movie.2026.1080p" data-language="eng"><td class="name">Movie</td><td>English</td><td><a href="downloadsubs.php?id=55">Subtitles</a></td></tr></table>`)}
	got, err := hdbitsSubtitleAdapter().Search(context.Background(), svc, nativeCfg("session=1"), trackerReq("hdbits", "https://hdbits.org/details.php?id=123"))
	if err != nil {
		t.Fatal(err)
	}
	if svc.lastReq.URL.Path != "/details.php" || svc.lastReq.URL.Query().Get("id") != "123" {
		t.Fatalf("unexpected request: %s", svc.lastReq.URL.String())
	}
	if len(got) != 1 || got[0].SourceURL != "https://hdbits.org/downloadsubs.php?id=55" || got[0].ProviderName != "hdbits" {
		t.Fatalf("unexpected candidates: %#v", got)
	}
}

func TestClusterCPrivateDownloadsExtractArchivesAndAddic7edKeepsRawSRT(t *testing.T) {
	archiveSvc := &nativeRoundTrip{body: trackerZip(t, "movie.srt", "1\n00:00:01,000 --> 00:00:02,000\nHi")}
	dl, err := avistazSubtitleAdapter().Download(context.Background(), archiveSvc, nativeCfg("uid=1"), providercore.Candidate{SourceURL: "/subtitles/42/download"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(dl.Content), "Hi") || !archiveSvc.lastDownload || !strings.Contains(archiveSvc.lastReq.Header.Get("Accept"), "application/zip") {
		t.Fatalf("archive download not extracted: %#v", dl)
	}

	rawSvc := &nativeRoundTrip{body: []byte("raw subtitle")}
	raw, err := addic7edAdapter().Download(context.Background(), rawSvc, nativeCfg("sid=1"), providercore.Candidate{SourceURL: "/updated/1/123/0"})
	if err != nil {
		t.Fatal(err)
	}
	if string(raw.Content) != "raw subtitle" || raw.URL != "https://www.addic7ed.com/updated/1/123/0" {
		t.Fatalf("raw download = %#v", raw)
	}
}

func trackerReq(source, infoURL string) providercore.SearchRequest {
	return providercore.SearchRequest{Title: "Movie", LanguageID: "eng", MediaContext: providercore.MediaContext{Provenance: []providercore.ReleaseProvenance{{Source: source, InfoURL: infoURL}}}}
}

func trackerZip(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(w, strings.NewReader(content)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
