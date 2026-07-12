package xsubs

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
)

func TestSearchPathAndCandidate(t *testing.T) {
	season, episode := int32(1), int32(2)
	if got := searchPath(providercore.SearchRequest{Title: "My Show", SeasonNumber: &season, EpisodeNumber: &episode}); got != "/search/My%20Show/1/2" {
		t.Fatalf("searchPath=%s", got)
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<li data-language="el"><a title="Greek release" href="/download/99">download</a></li>`))
	cand, ok := candidate(doc.Find("li"), "https://xsubs.tv/search", "")
	if !ok || cand.LanguageID != "el" || cand.ReleaseName != "Greek release" {
		t.Fatalf("candidate=%#v ok=%v", cand, ok)
	}
}
