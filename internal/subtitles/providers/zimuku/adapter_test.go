package zimuku

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
)

func TestSearchPathAndCandidate(t *testing.T) {
	if got := searchPath(providercore.SearchRequest{Title: "Movie Name"}); got != "/search/Movie%20Name" {
		t.Fatalf("searchPath=%s", got)
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<article lang="zh"><h3>Movie BluRay</h3><a href="/download/abc">下载</a></article>`))
	cand, ok := candidate(doc.Find("article"), "https://zimuku.org/search/Movie", "zh")
	if !ok || cand.LanguageID != "zh" || !strings.Contains(cand.SourceURL, "/download/abc") {
		t.Fatalf("candidate=%#v ok=%v", cand, ok)
	}
}
