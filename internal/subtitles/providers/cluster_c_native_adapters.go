package providers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"media-manager/internal/subtitles/providercore"
)

func karagargaAdapter() nativeCProvider {
	return nativeCProvider{key: "karagarga", baseURL: "https://karagarga.in", search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/browse.php", url.Values{"search": {nativeQuery(req)}, "search_type": {"title"}}, http.MethodGet
	}, parse: parseNativeHTML("karagarga", "tr", "a[href*='download.php'], a[href*='details.php']"), download: sourceDownload}
}

func ktuvitAdapter() nativeCProvider {
	return nativeCProvider{key: "ktuvit", baseURL: "https://www.ktuvit.me", captcha: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/Services/GetModuleAjax.ashx", url.Values{"moduleName": {"SubtitlesList"}, "SeriesName": {nativeQuery(req)}, "FilmName": {nativeQuery(req)}, "lang": {req.LanguageID}}, http.MethodPost
	}, parse: parseKtuvit, download: func(c providercore.Candidate) (string, url.Values, string) {
		return firstNonEmpty(c.SourceURL, "/Services/DownloadFile.ashx"), url.Values{"subtitleID": {strconv.FormatInt(c.FileID, 10)}}, http.MethodPost
	}}
}

func legendasdivxAdapter() nativeCProvider {
	return nativeCProvider{key: "legendasdivx", baseURL: "https://www.legendasdivx.pt", search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/modules.php", url.Values{"name": {"Downloads"}, "op": {"search"}, "query": {nativeQuery(req)}}, http.MethodGet
	}, parse: parseNativeHTML("legendasdivx", "tr, .download, .subtitle", "a[href*='d_op=getit'], a[href*='download']"), download: sourceDownload}
}

func legendasnetAdapter() nativeCProvider {
	return nativeCProvider{key: "legendasnet", baseURL: "https://legendas.net", captcha: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/search", url.Values{"q": {nativeQuery(req)}, "language": {req.LanguageID}}, http.MethodGet
	}, parse: parseNativeHTML("legendasnet", "tr, .subtitle, .item", "a[href*='download'], a[href*='/download/']"), download: sourceDownload}
}

func napisy24Adapter() nativeCProvider {
	return nativeCProvider{key: "napisy24", baseURL: "https://napisy24.pl", rawDownload: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/szukaj/", url.Values{"search": {nativeQuery(req)}}, http.MethodGet
	}, parse: parseNativeHTML("napisy24", ".tbl_subtitle tr, .subtitle-list .item, article", "a[href*='download'], a[href*='pobierz']"), download: sourceDownload}
}

func pipocasAdapter() nativeCProvider {
	return nativeCProvider{key: "pipocas", baseURL: "https://pipocas.tv", rawDownload: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/pesquisar", url.Values{"q": {nativeQuery(req)}}, http.MethodGet
	}, parse: parseNativeHTML("pipocas", ".subtitles-list li, .subtitle, tr", "a[href*='download'], a[href*='legenda']"), download: sourceDownload}
}

func subs4seriesAdapter() nativeCProvider {
	return nativeCProvider{key: "subs4series", baseURL: "https://www.subs4series.com", rawDownload: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/search_report.php", url.Values{"search": {nativeQuery(req)}}, http.MethodGet
	}, parse: parseNativeHTML("subs4series", "#search_results tr, .episode-row, .subtitle-row", "a[href*='download'], a[href*='subtitles']"), download: sourceDownload}
}

func subscenterAdapter() nativeCProvider {
	return nativeCProvider{key: "subscenter", baseURL: "https://www.subscenter.org", captcha: true, rawDownload: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/he/subtitle/search", url.Values{"q": {nativeQuery(req)}}, http.MethodPost
	}, parse: parseNativeHTML("subscenter", ".subtitle_result, .subtitleResult, tr", "a[href*='download'], a[href*='subtitle']"), download: sourceDownload}
}

func titloviAdapter() nativeCProvider {
	return nativeCProvider{key: "titlovi", baseURL: "https://titlovi.com", rawDownload: true, search: func(req providercore.SearchRequest) (string, url.Values, string) {
		return "/titlovi/", url.Values{"prijevod": {nativeQuery(req)}, "jezik": {req.LanguageID}}, http.MethodGet
	}, parse: parseNativeHTML("titlovi", ".subtitleContainer, .titloviRows tr, tr", "a[href*='download'], a[href*='preuzmi']"), download: sourceDownload}
}

func parseKtuvit(data []byte, pageURL, fallback string) ([]providercore.Candidate, error) {
	var payload struct {
		Subtitles []struct {
			ID                                    int64
			Name, FileName, Language, DownloadURL string
		}
	}
	if json.Unmarshal(data, &payload) == nil && len(payload.Subtitles) > 0 {
		out := make([]providercore.Candidate, 0, len(payload.Subtitles))
		for _, sub := range payload.Subtitles {
			out = append(out, providercore.Candidate{ProviderName: "ktuvit", FileID: sub.ID, LanguageID: firstNonEmpty(sub.Language, fallback), Format: "srt", ReleaseName: firstNonEmpty(sub.FileName, sub.Name), SourceURL: firstNonEmpty(sub.DownloadURL, "/Services/DownloadFile.ashx")})
		}
		return out, nil
	}
	return parseNativeHTML("ktuvit", "tr, .subtitle", "a[href*='DownloadFile'], a[href*='download']")(data, pageURL, fallback)
}
