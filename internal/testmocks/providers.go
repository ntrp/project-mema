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
	mux.HandleFunc("/tmdb/3/search/person", writeJSON(tmdbPersonSearch))
	mux.HandleFunc("/tmdb/3/discover/movie", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/genre/movie/list", writeJSON(tmdbMovieGenres))
	mux.HandleFunc("/tmdb/3/movie/936075", writeJSON(tmdbMovieDetails))
	mux.HandleFunc("/tmdb/3/movie/popular", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/movie/upcoming", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/movie/top_rated", writeJSON(tmdbMovieSearch))
	mux.HandleFunc("/tmdb/3/tv/popular", writeJSON(tmdbSeriesSearch))
	mux.HandleFunc("/tmdb/3/tv/on_the_air", writeJSON(tmdbSeriesSearch))
	mux.HandleFunc("/tmdb/3/tv/top_rated", writeJSON(tmdbSeriesSearch))
	mux.HandleFunc("/tmdb/3/trending/all/day", writeJSON(tmdbTrendingSearch))
	mux.HandleFunc("/tmdb/3/collection/123", writeJSON(tmdbMovieCollection))
	mux.HandleFunc("/tvdb/v4/login", writeJSON(tvdbLogin))
	mux.HandleFunc("/tvdb/v4/search", writeJSON(tvdbMovieSearch))
	mux.HandleFunc("/tvdb/v4/movies/900/extended", writeJSON(tvdbMovieDetails))
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

const tmdbMovieGenres = `{
  "genres": [
    { "id": 18, "name": "Drama" },
    { "id": 28, "name": "Action" }
  ]
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

const tmdbPersonSearch = `{
  "page": 1,
  "results": [
    {
      "id": 1001,
      "name": "Example Actor",
      "profile_path": "/actor.jpg",
      "popularity": 12.5,
      "known_for": [
        {
          "id": 936075,
          "media_type": "movie",
          "title": "Example Movie",
          "release_date": "2026-02-14"
        }
      ]
    }
  ],
  "total_pages": 1,
  "total_results": 1
}`

const tmdbMovieDetails = `{
  "id": 936075,
  "title": "Example Movie",
  "release_date": "2026-02-14",
  "overview": "A realistic local metadata detail response.",
  "credits": {
    "cast": [
      {
        "id": 1001,
        "name": "Example Actor",
        "character": "Lead",
        "profile_path": "/actor.jpg"
      }
    ],
    "crew": [
      {
        "id": 2001,
        "name": "Example Director",
        "job": "Director",
        "department": "Directing",
        "profile_path": "/director.jpg"
      },
      {
        "id": 2002,
        "name": "Example Writer",
        "job": "Screenplay",
        "department": "Writing",
        "profile_path": "/writer.jpg"
      }
    ]
  }
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

const tvdbLogin = `{
  "status": "success",
  "data": {
    "token": "scenario-tvdb-token"
  }
}`

const tvdbMovieSearch = `{
  "status": "success",
  "data": [
    {
      "id": "search-900",
      "tvdb_id": "900",
      "type": "movie",
      "name": "Example TVDB Movie",
      "year": "2026",
      "overview": "A realistic TVDB metadata result.",
      "image_url": "/tvdb-movie.jpg"
    }
  ]
}`

const tvdbMovieDetails = `{
  "status": "success",
  "data": {
    "id": 900,
    "name": "Example TVDB Movie",
    "year": "2026",
    "overview": "A realistic TVDB metadata detail response.",
    "image": "/tvdb-movie.jpg",
    "first_release": "2026-04-12",
    "runtime": 101,
    "status": { "name": "Released" },
    "originalLanguage": "eng",
    "score": 7.4,
    "boxOffice": "533300000.00",
    "budget": "120000000.00",
    "artworks": [
      { "type": 15, "image": "/tvdb-backdrop.jpg", "width": 1920, "height": 1080, "score": 50 }
    ],
    "production_countries": [
      { "country": "usa", "name": "United States" }
    ],
    "genres": [
      { "name": "Adventure" }
    ],
    "characters": [
      {
        "peopleId": 4001,
        "personName": "Example TVDB Actor",
        "name": "Lead",
        "peopleType": "Actor",
        "image": "/tvdb-actor.jpg"
      }
    ]
  }
}`
