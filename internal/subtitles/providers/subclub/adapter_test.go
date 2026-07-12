package subclub

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"media-manager/internal/subtitles/providercore"
)

type providerStub struct{}

func (providerStub) DoProviderRequest(req *http.Request, providerType string, isDownload bool) (*http.Response, error) {
	if isDownload {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("1\nfixture\n"))}, nil
	}
	if strings.Contains(req.URL.Path, "subtitles_archivecontent") {
		return resp(`<a href="../down.php?id=5&filename=abc">Fixture.S01E02.srt</a>`), nil
	}
	row := `<table id="tale_list"><tbody><tr><td></td><td><a class="sc_link" href="down.php?id=5">Fixture Show (2020) [1x2]</a></td><td></td><td><a href="https://imdb.com/title/tt1/">imdb</a></td><td></td><td></td><td>23.976</td><td><span>9.1</span></td><td>uploader</td></tr></tbody></table>`
	return resp(row), nil
}

func TestSearchAndDownloadFixture(t *testing.T) {
	s, e := int32(1), int32(2)
	candidates, err := Adapter.Search(context.Background(), providerStub{}, providercore.Config{BaseURL: "https://example.test"}, providercore.SearchRequest{MediaType: "serie", Title: "Fixture Show", LanguageID: "est", SeasonNumber: &s, EpisodeNumber: &e})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(candidates) != 1 || candidates[0].ProviderName != "subclub" || candidates[0].LanguageID != "est" {
		t.Fatalf("unexpected candidates: %#v", candidates)
	}
	dl, err := Adapter.Download(context.Background(), providerStub{}, providercore.Config{}, candidates[0])
	if err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if !bytes.Contains(dl.Content, []byte("fixture")) {
		t.Fatalf("unexpected content: %q", dl.Content)
	}
}

func resp(s string) *http.Response {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(s))}
}
