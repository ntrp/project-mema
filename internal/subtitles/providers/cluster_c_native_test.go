package providers

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type nativeRoundTrip struct {
	lastReq      *http.Request
	lastProvider string
	lastDownload bool
	body         []byte
	status       int
}

func (f *nativeRoundTrip) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	f.lastReq = req
	f.lastProvider = providerType
	f.lastDownload = isDownload
	status := f.status
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{StatusCode: status, Body: io.NopCloser(bytes.NewReader(f.body)), Header: http.Header{}}, nil
}

func TestNativeCProvidersRegistered(t *testing.T) {
	for _, provider := range nativeCProviders {
		adapter, ok := AdapterFor(provider.key)
		if !ok || adapter == nil {
			t.Fatalf("%s adapter not registered", provider.key)
		}
	}
}

func TestNativeCPrerequisitesAreProviderSpecific(t *testing.T) {
	if err := karagargaAdapter().prereq(providercore.Config{}); !errors.Is(err, providercore.ErrPrivateMembershipRequired) {
		t.Fatalf("karagarga prerequisite = %v", err)
	}
	if err := ktuvitAdapter().prereq(providercore.Config{}); !errors.Is(err, providercore.ErrCaptchaRequired) {
		t.Fatalf("ktuvit prerequisite = %v", err)
	}
	if err := subscenterAdapter().prereq(providercore.Config{}); !errors.Is(err, providercore.ErrCaptchaRequired) {
		t.Fatalf("subscenter prerequisite = %v", err)
	}
}

func TestKaragargaSearchUsesBazarrCookieSessionShape(t *testing.T) {
	svc := &nativeRoundTrip{body: []byte(`<table><tr data-language="eng"><td>Film 1999</td><td><a href="download.php?id=42">kg.release</a></td></tr></table>`)}
	cfg := nativeCfg("sid=abc")
	year := int32(1999)
	got, err := karagargaAdapter().Search(context.Background(), svc, cfg, providercore.SearchRequest{Title: "Film", Year: &year, LanguageID: "eng"})
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if len(got) != 1 || got[0].ProviderName != "karagarga" || !strings.Contains(got[0].SourceURL, "download.php?id=42") {
		t.Fatalf("unexpected candidates: %#v", got)
	}
	if svc.lastReq.Method != http.MethodGet || svc.lastReq.URL.Path != "/browse.php" {
		t.Fatalf("unexpected request: %s %s", svc.lastReq.Method, svc.lastReq.URL.String())
	}
	if svc.lastReq.URL.Query().Get("search_type") != "title" || !strings.Contains(svc.lastReq.URL.Query().Get("search"), "Film") {
		t.Fatalf("unexpected search query: %s", svc.lastReq.URL.RawQuery)
	}
	if svc.lastReq.Header.Get("Cookie") != "sid=abc" || svc.lastReq.Header.Get("User-Agent") == "" {
		t.Fatalf("missing session headers: %#v", svc.lastReq.Header)
	}
}

func TestKtuvitSearchPostsAjaxModuleAndParsesJSON(t *testing.T) {
	svc := &nativeRoundTrip{body: []byte(`{"Subtitles":[{"ID":77,"Name":"Episode","FileName":"Episode.S01E02","Language":"heb","DownloadURL":"/Services/DownloadFile.ashx"}]}`)}
	season, episode := int32(1), int32(2)
	got, err := ktuvitAdapter().Search(context.Background(), svc, nativeCfg("ASP.NET_SessionId=x"), providercore.SearchRequest{Title: "Episode", SeasonNumber: &season, EpisodeNumber: &episode, LanguageID: "heb"})
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	if svc.lastReq.Method != http.MethodPost || svc.lastReq.URL.Path != "/Services/GetModuleAjax.ashx" {
		t.Fatalf("unexpected request: %s %s", svc.lastReq.Method, svc.lastReq.URL.String())
	}
	body, _ := io.ReadAll(svc.lastReq.Body)
	if !strings.Contains(string(body), "moduleName=SubtitlesList") || !strings.Contains(string(body), "S01E02") {
		t.Fatalf("unexpected ktuvit form: %s", body)
	}
	if len(got) != 1 || got[0].FileID != 77 || got[0].ReleaseName != "Episode.S01E02" {
		t.Fatalf("unexpected candidates: %#v", got)
	}
}

func TestLegendasDivxAndNetSearchRoutesAreNative(t *testing.T) {
	cases := []struct {
		name      string
		adapter   nativeCProvider
		wantPath  string
		wantParam string
	}{
		{"legendasdivx", legendasdivxAdapter(), "/modules.php", "op=search"},
		{"legendasnet", legendasnetAdapter(), "/search", "q=Movie"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &nativeRoundTrip{body: []byte(`<div class="subtitle" data-language="por"><a href="/download/9">Movie rip</a></div>`)}
			got, err := tc.adapter.Search(context.Background(), svc, nativeCfg("member=1"), providercore.SearchRequest{Title: "Movie", LanguageID: "por"})
			if err != nil {
				t.Fatalf("Search error: %v", err)
			}
			if svc.lastReq.URL.Path != tc.wantPath || !strings.Contains(svc.lastReq.URL.RawQuery, tc.wantParam) {
				t.Fatalf("unexpected route: %s", svc.lastReq.URL.String())
			}
			if len(got) != 1 || got[0].ProviderName != tc.name || got[0].LanguageID != "por" {
				t.Fatalf("unexpected candidates: %#v", got)
			}
		})
	}
}

func TestNativeCSearchParsesRealShapedFixtures(t *testing.T) {
	cases := []struct{ key, body, wantURL, wantTitle, wantLang string }{
		{"napisy24", `<table class="tbl_subtitle"><tr data-release="Example.2020.PL.1080p" data-language="pol"><td class="title">Example release</td><td><a href="/download/123/example.zip">Pobierz</a></td></tr></table>`, "https://napisy24.pl/download/123/example.zip", "Example.2020.PL.1080p", "pol"},
		{"pipocas", `<ul class="subtitles-list"><li><span class="release">Example.S01E02.HDTV</span><span class="language">por</span><a href="/legenda/55">download</a></li></ul>`, "https://pipocas.tv/legenda/55", "Example.S01E02.HDTV", "por"},
		{"subs4series", `<table id="search_results"><tr><td class="episode">Example - 1x02</td><td class="lang">ell</td><td><a href="/subtitles/example-1x02/download">Download</a></td></tr></table>`, "https://www.subs4series.com/subtitles/example-1x02/download", "Example - 1x02", "ell"},
		{"subscenter", `<div class="subtitle_result"><a class="title" href="/he/subtitle/download/77">Example S01E02</a><span class="language">heb</span></div>`, "https://www.subscenter.org/he/subtitle/download/77", "Example S01E02", "heb"},
		{"titlovi", `<div class="subtitleContainer" data-release="Example.S01E02.Balkan"><span class="jezik">hrv</span><a href="/download/?type=1&mediaid=42">preuzmi</a></div>`, "https://titlovi.com/download/?type=1&mediaid=42", "Example.S01E02.Balkan", "hrv"},
	}
	season, episode := int32(1), int32(2)
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			service := &nativeRoundTrip{body: []byte(tc.body)}
			candidates, err := findNativeProvider(t, tc.key).Search(context.Background(), service, nativeCfg("sid=ok"), providercore.SearchRequest{Title: "Example", LanguageID: "eng", SeasonNumber: &season, EpisodeNumber: &episode})
			if err != nil {
				t.Fatalf("Search error: %v", err)
			}
			if len(candidates) != 1 {
				t.Fatalf("candidate count = %d", len(candidates))
			}
			candidate := candidates[0]
			if candidate.ProviderName != tc.key || candidate.ReleaseName != tc.wantTitle || candidate.LanguageID != tc.wantLang || candidate.SourceURL != tc.wantURL {
				t.Fatalf("unexpected candidate: %#v", candidate)
			}
			if got := service.lastReq.Header.Get("Cookie"); got != "sid=ok" {
				t.Fatalf("Cookie = %q", got)
			}
			if service.lastProvider != tc.key {
				t.Fatalf("provider = %q", service.lastProvider)
			}
		})
	}
}

func TestNativeCSearchUsesProviderSpecificRequestShape(t *testing.T) {
	for _, tc := range []struct{ key, method, field string }{{"napisy24", http.MethodGet, "search"}, {"pipocas", http.MethodGet, "q"}, {"subs4series", http.MethodGet, "search"}, {"subscenter", http.MethodPost, ""}, {"titlovi", http.MethodGet, "prijevod"}} {
		t.Run(tc.key, func(t *testing.T) {
			service := &nativeRoundTrip{body: []byte(`<div></div>`)}
			_, err := findNativeProvider(t, tc.key).Search(context.Background(), service, nativeCfg("sid=ok"), providercore.SearchRequest{Title: "Example"})
			if err != nil {
				t.Fatalf("Search error: %v", err)
			}
			if service.lastReq.Method != tc.method {
				t.Fatalf("method = %s", service.lastReq.Method)
			}
			if tc.field != "" && service.lastReq.URL.Query().Get(tc.field) != "Example" {
				t.Fatalf("%s query = %q", tc.field, service.lastReq.URL.RawQuery)
			}
			if tc.key == "subscenter" && service.lastReq.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
				t.Fatalf("missing form content type")
			}
		})
	}
}

func TestNativeCDownloadExtractsArchiveMember(t *testing.T) {
	svc := &nativeRoundTrip{body: zipSubtitle(t, "movie.srt", "1\n00:00:01,000 --> 00:00:02,000\nOlá")}
	dl, err := legendasnetAdapter().Download(context.Background(), svc, nativeCfg("member=1"), providercore.Candidate{SourceURL: "/download/9/movie.zip"})
	if err != nil {
		t.Fatalf("Download error: %v", err)
	}
	if !strings.Contains(string(dl.Content), "Olá") || !strings.Contains(dl.URL, "/download/9/movie.zip") {
		t.Fatalf("unexpected download: %#v", dl)
	}
	if !svc.lastDownload || !strings.Contains(svc.lastReq.Header.Get("Accept"), "application/zip") {
		t.Fatalf("download request missing archive semantics: %v %#v", svc.lastDownload, svc.lastReq.Header)
	}
}

func TestNativeCRawDownloadSendsCookieAndDownloadFlag(t *testing.T) {
	service := &nativeRoundTrip{body: []byte("subtitle bytes")}
	download, err := findNativeProvider(t, "titlovi").Download(context.Background(), service, nativeCfg("sid=ok"), providercore.Candidate{SourceURL: "/download/42"})
	if err != nil {
		t.Fatalf("Download error: %v", err)
	}
	if string(download.Content) != "subtitle bytes" || download.URL != "https://titlovi.com/download/42" {
		t.Fatalf("download = %#v", download)
	}
	if !service.lastDownload || service.lastReq.Header.Get("Accept") == "" || service.lastReq.Header.Get("Cookie") != "sid=ok" {
		t.Fatalf("download request not shaped correctly")
	}
}

func findNativeProvider(t *testing.T, key string) nativeCProvider {
	t.Helper()
	for _, provider := range nativeCProviders {
		if provider.key == key {
			return provider
		}
	}
	t.Fatalf("native provider %s not found", key)
	return nativeCProvider{}
}

func nativeCfg(cookie string) providercore.Config {
	return providercore.Config{SecretSettings: map[string]string{"cookies": cookie}}
}

func zipSubtitle(t *testing.T, name, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
