package providers

import (
	"archive/zip"
	"bytes"
	"context"
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

func TestNativeCPrerequisitesAreProviderSpecific(t *testing.T) {
	if err := karagargaAdapter().prereq(providercore.Config{}); err != providercore.ErrPrivateMembershipRequired {
		t.Fatalf("karagarga prerequisite = %v", err)
	}
	if err := ktuvitAdapter().prereq(providercore.Config{}); err != providercore.ErrCaptchaRequired {
		t.Fatalf("ktuvit prerequisite = %v", err)
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
