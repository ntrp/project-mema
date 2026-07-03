package testmocks

import (
	"net/http"
	"net/http/httptest"
)

type ProviderServer struct {
	*httptest.Server
}

func NewProviderServer() *ProviderServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", writeJSON(`{"status":"ok"}`))
	mux.HandleFunc("/torznab/api", handleTorznab)
	mux.HandleFunc("/tmdb/3/search/movie", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/search/tv", writeJSON(tmdbSeriesSearch))
	mux.HandleFunc("/tmdb/3/movie/936075", writeJSON(tmdbMovieDetails))
	mux.HandleFunc("/tmdb/3/movie/popular", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/movie/upcoming", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/movie/top_rated", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/tv/popular", writeJSON(tmdbSeriesSearch))
	mux.HandleFunc("/tmdb/3/tv/on_the_air", writeJSON(tmdbSeriesSearch))
	mux.HandleFunc("/tmdb/3/tv/top_rated", writeJSON(tmdbSeriesSearch))
	mux.HandleFunc("/tmdb/3/trending/all/day", writeJSON(tmdbTrendingSearch))
	mux.HandleFunc("/tmdb/3/collection/123", writeJSON(tmdbMovieCollection))
	return &ProviderServer{Server: httptest.NewServer(mux)}
}

func writeJSON(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}
}

func handleTorznab(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	if r.URL.Query().Get("t") == "caps" {
		_, _ = w.Write([]byte(torznabCaps))
		return
	}
	_, _ = w.Write([]byte(torznabSearch))
}

const torznabCaps = `<caps>
  <server title="Local Torznab Mock" version="1.0"/>
  <limits max="100" default="50"/>
  <categories>
    <category id="2000" name="Movies"><subcat id="2040" name="HD"/></category>
    <category id="5000" name="TV"><subcat id="5070" name="Anime"/></category>
  </categories>
</caps>`

const torznabSearch = `<rss xmlns:torznab="http://torznab.com/schemas/2015/feed"><channel>
  <title>Local releases</title>
  <item>
    <title>Example.Movie.2026.1080p.WEB-DL</title>
    <link>https://indexer.test/download/example-movie</link>
    <guid>https://indexer.test/release/example-movie</guid>
    <pubDate>Fri, 03 Jul 2026 04:00:00 +0200</pubDate>
    <size>8589934592</size>
    <torznab:attr name="seeders" value="42"/>
    <torznab:attr name="peers" value="7"/>
  </item>
</channel></rss>`

const tmdbMovieSearch = `{
  "page": 1,
  "results": [
    {
      "id": 936075,
      "title": "Example Movie",
      "release_date": "2026-02-14",
      "overview": "A realistic local metadata result."
    }
  ],
  "total_pages": 1,
  "total_results": 1
}`

const tmdbTrendingSearch = `{
  "page": 1,
  "results": [
    {
      "id": 936075,
      "media_type": "movie",
      "title": "Example Movie",
      "release_date": "2026-02-14",
      "overview": "A realistic local metadata result."
    }
  ],
  "total_pages": 1,
  "total_results": 1
}`

const tmdbSeriesSearch = `{
  "page": 1,
  "results": [
    {
      "id": 2026,
      "name": "Example Series",
      "first_air_date": "2026-03-01",
      "overview": "A realistic local series metadata result."
    }
  ],
  "total_pages": 1,
  "total_results": 1
}`

const tmdbMovieDetails = `{
  "id": 936075,
  "title": "Example Movie",
  "release_date": "2026-02-14",
  "overview": "A realistic local metadata detail response."
}`

const tmdbMovieCollection = `{
  "id": 123,
  "name": "Example Collection",
  "overview": "A local metadata collection.",
  "parts": [
    {
      "id": 936075,
      "title": "Example Movie",
      "release_date": "2026-02-14",
      "overview": "A realistic local metadata result."
    }
  ]
}`
