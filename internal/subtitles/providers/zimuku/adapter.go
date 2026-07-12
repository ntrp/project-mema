package zimuku

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativehtml"
)

func init() {
	nativehtml.Register(nativehtml.Spec{Key: "zimuku", BaseURL: "https://zimuku.org", ArchiveDownload: true, SearchPath: searchPath, Query: func(q url.Values, sr providercore.SearchRequest) { q.Set("q", nativehtml.QueryText(sr)) }, Candidate: candidate})
}

func searchPath(sr providercore.SearchRequest) string {
	if strings.TrimSpace(sr.Title) == "" {
		return "/search"
	}
	return fmt.Sprintf("/search/%s", url.PathEscape(nativehtml.QueryText(sr)))
}

func candidate(sel *goquery.Selection, pageURL, fallbackLang string) (providercore.Candidate, bool) {
	cand, ok := nativehtml.LinkCandidate(sel, pageURL, fallbackLang, "a[href*='/download/'], a[href*='/detail/'], a[href*='subtitle']")
	if !ok {
		return cand, false
	}
	lower := strings.ToLower(cand.SourceURL)
	if !strings.Contains(lower, "/download/") && !strings.Contains(lower, "/detail/") && !strings.Contains(lower, "subtitle") {
		return providercore.Candidate{}, false
	}
	if cand.LanguageID == "" {
		cand.LanguageID = fallbackLang
	}
	return cand, true
}
