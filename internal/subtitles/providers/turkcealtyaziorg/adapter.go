package turkcealtyaziorg

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativehtml"
)

func init() {
	nativehtml.Register(nativehtml.Spec{Key: "turkcealtyaziorg", BaseURL: "https://turkcealtyazi.org", ArchiveDownload: true, SearchPath: func(providercore.SearchRequest) string { return "/find.php" }, Query: func(q url.Values, sr providercore.SearchRequest) {
		q.Set("cat", "sub")
		q.Set("find", nativehtml.QueryText(sr))
	}, Candidate: candidate})
}

func candidate(sel *goquery.Selection, pageURL, fallbackLang string) (providercore.Candidate, bool) {
	cand, ok := nativehtml.LinkCandidate(sel, pageURL, fallbackLang, "a[href*='indirmek'], a[href*='download'], a[href*='subtitle'], a[href$='.html']")
	if !ok {
		return cand, false
	}
	lower := strings.ToLower(cand.SourceURL)
	if !strings.Contains(lower, "indirmek") && !strings.Contains(lower, "download") && !strings.Contains(lower, "subtitle") {
		return providercore.Candidate{}, false
	}
	cand.LanguageID = "tr"
	return cand, true
}
