package providers

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"media-manager/internal/subtitles/providercore"

	"github.com/PuerkitoBio/goquery"
)

func addic7edAdapter() nativeCProvider {
	return nativeCProvider{key: "addic7ed", baseURL: "https://www.addic7ed.com", captcha: true, rawDownload: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/search.php", url.Values{"search": {nativeQuery(req)}, "Submit": {"Search"}}, http.MethodGet
	}, parse: parseAddic7ed, download: sourceDownload}
}

func avistazSubtitleAdapter() nativeCProvider {
	return avistaZLikeSubtitleAdapter("avistaz", "https://avistaz.to")
}

func cinemazSubtitleAdapter() nativeCProvider {
	return avistaZLikeSubtitleAdapter("cinemaz", "https://cinemaz.to")
}

func hdbitsSubtitleAdapter() nativeCProvider {
	return nativeCProvider{key: "hdbits", baseURL: "https://hdbits.org", provenanceSource: "hdbits", search: func(req providercore.SearchRequest) (string, url.Values, string) {
		infoURL := releaseProvenanceURL(req, "hdbits")
		return hdbitsDetailsPath(infoURL), nil, http.MethodGet
	}, parse: parseHDBitsSubtitles, download: sourceDownload}
}

func avistaZLikeSubtitleAdapter(key, baseURL string) nativeCProvider {
	return nativeCProvider{key: key, baseURL: baseURL, provenanceSource: key, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		form := url.Values{"search": {nativeQuery(req)}}
		if id := trailingNumericID(releaseProvenanceURL(req, key)); id != "" {
			form.Set("torrent_id", id)
		}
		return "/subtitles", form, http.MethodGet
	}, parse: parseAvistaZLikeSubtitles(key), download: sourceDownload}
}

func parseAddic7ed(data []byte, pageURL, fallback string) ([]providercore.Candidate, error) {
	doc, err := providercore.ParseHTML(data)
	if err != nil {
		return nil, err
	}
	out := []providercore.Candidate{}
	seen := map[string]bool{}
	doc.Find("tr.epeven, tr.epodd, tr").Each(func(_ int, row *goquery.Selection) {
		link := row.Find("a[href*='/updated/'], a[href*='/original/'], a[href*='download']").First()
		if link.Length() == 0 {
			return
		}
		href, _ := link.Attr("href")
		cells := row.Find("td")
		title := firstNonEmpty(attr(row, "data-release"), strings.TrimSpace(cells.Eq(1).Text()), strings.TrimSpace(cells.Eq(2).Text()), strings.TrimSpace(link.Text()))
		lang := firstNonEmpty(attr(row, "data-language"), strings.TrimSpace(row.Find("td.language, .language").First().Text()), strings.TrimSpace(cells.Eq(3).Text()), fallback)
		appendCandidate(&out, seen, "addic7ed", pageURL, href, title, lang)
	})
	return out, nil
}

func parseAvistaZLikeSubtitles(provider string) func([]byte, string, string) ([]providercore.Candidate, error) {
	return func(data []byte, pageURL, fallback string) ([]providercore.Candidate, error) {
		doc, err := providercore.ParseHTML(data)
		if err != nil {
			return nil, err
		}
		out := []providercore.Candidate{}
		seen := map[string]bool{}
		doc.Find("tr, .subtitle, .subtitle-row, .card").Each(func(_ int, row *goquery.Selection) {
			link := row.Find("a[href*='subtitle'][href*='download'], a[href*='download'], a[href*='/subtitles/']").First()
			if link.Length() == 0 {
				return
			}
			href, _ := link.Attr("href")
			title := firstNonEmpty(attr(row, "data-release"), strings.TrimSpace(row.Find(".release, .name, .title, td:first-child").First().Text()), strings.TrimSpace(link.Text()))
			lang := firstNonEmpty(attr(row, "data-language"), strings.TrimSpace(row.Find(".language, .lang, td:nth-child(2)").First().Text()), fallback)
			appendCandidate(&out, seen, provider, pageURL, href, title, lang)
		})
		return out, nil
	}
}

func parseHDBitsSubtitles(data []byte, pageURL, fallback string) ([]providercore.Candidate, error) {
	doc, err := providercore.ParseHTML(data)
	if err != nil {
		return nil, err
	}
	out := []providercore.Candidate{}
	seen := map[string]bool{}
	doc.Find("tr, .subtitle, .subtitles li").Each(func(_ int, row *goquery.Selection) {
		link := row.Find("a[href*='downloadsubs'], a[href*='subtitle'], a[href*='download']").First()
		if link.Length() == 0 {
			return
		}
		href, _ := link.Attr("href")
		title := firstNonEmpty(attr(row, "data-release"), strings.TrimSpace(row.Find(".release, .name, td:first-child").First().Text()), strings.TrimSpace(link.Text()))
		lang := firstNonEmpty(attr(row, "data-language"), strings.TrimSpace(row.Find(".language, .lang, td:nth-child(2)").First().Text()), fallback)
		appendCandidate(&out, seen, "hdbits", pageURL, href, title, lang)
	})
	return out, nil
}

func appendCandidate(out *[]providercore.Candidate, seen map[string]bool, provider, pageURL, href, title, lang string) {
	if strings.TrimSpace(href) == "" || strings.TrimSpace(title) == "" {
		return
	}
	abs := resolveAgainst(pageURL, href)
	key := title + "\x00" + abs
	if seen[key] {
		return
	}
	seen[key] = true
	*out = append(*out, providercore.Candidate{ProviderName: provider, LanguageID: lang, Format: formatFrom(href), ReleaseName: title, SourceURL: abs, SourceRef: href})
}

func releaseProvenanceURL(req providercore.SearchRequest, source string) string {
	want := strings.ToLower(strings.TrimSpace(source))
	for _, provenance := range req.MediaContext.Provenance {
		got := strings.ToLower(strings.TrimSpace(provenance.Source))
		if got == want || strings.Contains(got, want) || strings.Contains(strings.ToLower(provenance.InfoURL), want) {
			if strings.TrimSpace(provenance.InfoURL) != "" {
				return strings.TrimSpace(provenance.InfoURL)
			}
		}
	}
	return ""
}

func hdbitsDetailsPath(infoURL string) string {
	if parsed, err := url.Parse(infoURL); err == nil {
		if id := parsed.Query().Get("id"); id != "" {
			return "/details.php?id=" + url.QueryEscape(id)
		}
	}
	return infoURL
}

func trailingNumericID(raw string) string {
	matches := regexp.MustCompile(`\d+`).FindAllString(raw, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[len(matches)-1]
}
