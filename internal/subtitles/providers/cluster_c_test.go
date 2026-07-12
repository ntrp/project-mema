package providers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type clusterCStubService struct {
	status int
	body   string
	last   *http.Request
}

func (s *clusterCStubService) DoProviderRequest(request *http.Request, _ string, _ bool) (*http.Response, error) {
	s.last = request
	status := s.status
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(s.body)), Header: http.Header{}}, nil
}

func TestClusterCProvidersRegistered(t *testing.T) {
	for _, provider := range clusterCProviders {
		adapter, ok := AdapterFor(provider.key)
		if !ok {
			t.Fatalf("%s adapter not registered", provider.key)
		}
		if adapter == nil {
			t.Fatalf("%s adapter is nil", provider.key)
		}
	}
}

func TestClusterCPrerequisitesReturnTypedErrors(t *testing.T) {
	privateProvider := clusterCProvider{key: "avistaz", private: true, baseURL: "https://avistaz.to"}
	if err := privateProvider.Test(context.Background(), &clusterCStubService{}, providercore.Config{}); !errors.Is(err, providercore.ErrPrivateMembershipRequired) {
		t.Fatalf("private provider error = %v, want ErrPrivateMembershipRequired", err)
	}

	captchaProvider := clusterCProvider{key: "addic7ed", captcha: true, baseURL: "https://www.addic7ed.com"}
	if err := captchaProvider.Test(context.Background(), &clusterCStubService{}, providercore.Config{}); !errors.Is(err, providercore.ErrCaptchaRequired) {
		t.Fatalf("captcha provider error = %v, want ErrCaptchaRequired", err)
	}
	captchaPrivateProvider := clusterCProvider{key: "ktuvit", private: true, captcha: true, baseURL: "https://www.ktuvit.me"}
	if err := captchaPrivateProvider.Test(context.Background(), &clusterCStubService{}, providercore.Config{}); !errors.Is(err, providercore.ErrCaptchaRequired) {
		t.Fatalf("captcha-private provider error = %v, want ErrCaptchaRequired", err)
	}
}

func TestClusterCSearchParsesFixtureAndSendsCookies(t *testing.T) {
	provider := clusterCProvider{key: "titlovi", private: true, baseURL: "https://titlovi.com", searchPath: "/titlovi/"}
	service := &clusterCStubService{body: `<table><tr data-release="Example.S01E02.1080p" data-language="hr" data-format="srt"><td>Example</td><td><a href="/download/42">download</a></td></tr></table>`}
	cookie := "uid=1; pass=secret"
	config := providercore.Config{SecretSettings: map[string]string{"cookies": cookie}}
	season, episode := int32(1), int32(2)

	candidates, err := provider.Search(context.Background(), service, config, providercore.SearchRequest{Title: "Example", LanguageID: "hr", SeasonNumber: &season, EpisodeNumber: &episode})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("candidate count = %d, want 1", len(candidates))
	}
	candidate := candidates[0]
	if candidate.ProviderName != "titlovi" || candidate.ReleaseName != "Example.S01E02.1080p" || candidate.LanguageID != "hr" || candidate.Format != "srt" {
		t.Fatalf("unexpected candidate: %#v", candidate)
	}
	if candidate.SourceURL != "https://titlovi.com/download/42" {
		t.Fatalf("SourceURL = %q", candidate.SourceURL)
	}
	if got := service.last.Header.Get("Cookie"); got != cookie {
		t.Fatalf("Cookie header = %q", got)
	}
	if query := service.last.URL.Query().Get("q"); !strings.Contains(query, "Example") || !strings.Contains(query, "S01E02") {
		t.Fatalf("search query = %q", query)
	}
}

func TestClusterCDownloadReadsContent(t *testing.T) {
	provider := clusterCProvider{key: "zimuku", private: true, baseURL: "https://zimuku.org"}
	service := &clusterCStubService{body: "subtitle bytes"}
	config := providercore.Config{SecretSettings: map[string]string{"cookies": "session=ok"}}

	download, err := provider.Download(context.Background(), service, config, providercore.Candidate{SourceURL: "/download/abc"})
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if string(download.Content) != "subtitle bytes" || download.URL != "https://zimuku.org/download/abc" {
		t.Fatalf("unexpected download: %#v", download)
	}
}
