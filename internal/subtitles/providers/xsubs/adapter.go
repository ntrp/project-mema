package xsubs

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"media-manager/internal/subtitles/providercore"
	"media-manager/internal/subtitles/providers/nativehtml"
)

func init() {
	nativehtml.Register(nativehtml.Spec{Key: "xsubs", BaseURL: "https://xsubs.tv", NeedsCookie: true, ArchiveDownload: true, SearchPath: searchPath, Query: func(q url.Values, sr providercore.SearchRequest) { q.Set("q", nativehtml.QueryText(sr)) }, Candidate: candidate})
}

func searchPath(sr providercore.SearchRequest) string {
	if sr.SeasonNumber != nil && sr.EpisodeNumber != nil {
		return fmt.Sprintf("/search/%s/%d/%d", url.PathEscape(sr.Title), *sr.SeasonNumber, *sr.EpisodeNumber)
	}
	return "/search"
}

func candidate(sel *goquery.Selection, pageURL, fallbackLang string) (providercore.Candidate, bool) {
	cand, ok := nativehtml.LinkCandidate(sel, pageURL, fallbackLang, "a[href*='download'], a[href*='/xsub/'], a[href*='subtitles']")
	if !ok {
		return cand, false
	}
	lower := strings.ToLower(cand.SourceURL)
	if !strings.Contains(lower, "download") && !strings.Contains(lower, "xsub") && !strings.Contains(lower, "subtitles") {
		return providercore.Candidate{}, false
	}
	if cand.LanguageID == "" {
		cand.LanguageID = "el"
	}
	return cand, true
}
