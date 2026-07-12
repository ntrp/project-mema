package titulky

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativehtml"
)

func init() {
	nativehtml.Register(nativehtml.Spec{Key: "titulky", BaseURL: "https://www.titulky.com", NeedsCookie: true, ArchiveDownload: true, SearchPath: func(providercore.SearchRequest) string { return "/" }, Query: func(q url.Values, sr providercore.SearchRequest) { q.Set("Fulltext", nativehtml.QueryText(sr)) }, Candidate: candidate})
}

func candidate(sel *goquery.Selection, pageURL, fallbackLang string) (providercore.Candidate, bool) {
	cand, ok := nativehtml.LinkCandidate(sel, pageURL, fallbackLang, "a[href*='id='], a[href*='download'], a[href*='titulky']")
	if !ok {
		return cand, false
	}
	lower := strings.ToLower(cand.SourceURL)
	if !strings.Contains(lower, "download") && !strings.Contains(lower, "id=") && !strings.Contains(lower, "titulky") {
		return providercore.Candidate{}, false
	}
	if cand.LanguageID == "" {
		cand.LanguageID = "cs"
	}
	return cand, true
}
