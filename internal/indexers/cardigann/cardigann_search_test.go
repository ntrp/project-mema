package cardigann

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCardigannHTMLSearchUsesYAMLDefinitionAndDownloadSelector(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/search/Example":
			body := `<html><body><table><tr class="release">
				<td class="name"><a class="title" href="/details/1">Example.Movie.2026.1080p</a></td>
				<td class="size">1.5 GiB</td><td class="seeders">12</td><td class="leechers">3</td>
			</tr></table></body></html>`
			return response(http.StatusOK, body), nil
		case "/details/1":
			return response(http.StatusOK, `<a href="magnet:?xt=urn:btih:abc123">Magnet</a>`), nil
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
			return response(http.StatusNotFound, ""), nil
		}
	})
	service := cardigannTestService(client, "mockhtml", htmlCardigannDefinition)

	releases, err := service.Search(context.Background(), Config{
		ID:             "idx-1",
		DefinitionID:   "mockhtml",
		Name:           "Mock HTML",
		Implementation: "Cardigann",
		Protocol:       "torrent",
		BaseURL:        "http://indexer.local/",
	}, "Example", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Movie.2026.1080p" {
		t.Fatalf("title = %q", release.Title)
	}
	if !strings.HasPrefix(release.DownloadURL, "magnet:?xt=urn:btih:abc123") {
		t.Fatalf("download url = %q", release.DownloadURL)
	}
	if release.SizeBytes != 1610612736 {
		t.Fatalf("size = %d", release.SizeBytes)
	}
	if release.Seeders == nil || *release.Seeders != 12 || release.Peers == nil || *release.Peers != 3 {
		t.Fatalf("peers = %#v seeders = %#v", release.Peers, release.Seeders)
	}
}

func TestCardigannJSONSearchUsesYAMLDefinition(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api" || r.URL.Query().Get("q") != "Example" {
			t.Fatalf("unexpected request %s", r.URL.String())
		}
		body := `[{"name":"Example.Series.S01E01","info_hash":"abc123","size":"2048","seeders":5,"leechers":2,"added":1783036800}]`
		return response(http.StatusOK, body), nil
	})
	service := cardigannTestService(client, "mockjson", jsonCardigannDefinition)

	releases, err := service.Search(context.Background(), Config{
		ID:             "idx-2",
		DefinitionID:   "mockjson",
		Name:           "Mock JSON",
		Implementation: "Cardigann",
		Protocol:       "torrent",
		BaseURL:        "http://indexer.local/",
	}, "Example", "series")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	release := releases[0]
	if release.Title != "Example.Series.S01E01" {
		t.Fatalf("title = %q", release.Title)
	}
	if release.DownloadURL != "magnet:?xt=urn:btih:abc123" {
		t.Fatalf("download = %q", release.DownloadURL)
	}
	if release.PublishedAt == nil {
		t.Fatal("published date was not parsed")
	}
}

func TestCardigannYTSRowsUseTorrentAttributeAndMovieParent(t *testing.T) {
	client := fakeHTTPDoer(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/api/v2/list_movies.json" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		body := `{"data":{"movie_count":1,"movies":[{"year":2026,"title":"Example Movie","title_long":"Example Movie","url":"https://yts.gg/movies/example","large_cover_image":"https://yts.gg/cover.jpg","imdb_code":"tt1234567","torrents":[{"quality":"720p","audio_channels":"5.1","bit_depth":"10","type":"web","video_codec":"x264","url":"https://yts.gg/torrent/download/abc","hash":"abc123","date_uploaded_unix":1783036800,"size_bytes":2048,"seeds":11,"peers":2}]}]}}`
		return response(http.StatusOK, body), nil
	})
	service := New(client)
	service.loader.remote = ""
	fields, _ := json.Marshal([]struct {
		Name  string `json:"name"`
		Value any    `json:"value"`
	}{{Name: "apiurl", Value: "api.yts.test"}})

	releases, err := service.Search(context.Background(), Config{
		ID:             "idx-yts",
		DefinitionID:   "yts",
		Name:           "YTS",
		Implementation: "Cardigann",
		Protocol:       "torrent",
		BaseURL:        "https://yts.gg/",
		Fields:         fields,
	}, "Example Movie", "movie")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 {
		t.Fatalf("release count = %d", len(releases))
	}
	if !strings.Contains(releases[0].Title, "720p WEBRip 5.1 10Bit x264 -YTS") {
		t.Fatalf("title = %q", releases[0].Title)
	}
	if releases[0].DownloadURL != "https://yts.gg/torrent/download/abc" {
		t.Fatalf("download = %q", releases[0].DownloadURL)
	}
}

func TestGeneratedCardigannDefinitionsDecode(t *testing.T) {
	loader := newCardigannLoader(nil)
	if len(loader.local) < 500 {
		t.Fatalf("embedded definitions = %d, want full v11 catalog", len(loader.local))
	}
	for id, body := range loader.local {
		if _, err := decodeCardigannDefinition(id, []byte(body)); err != nil {
			t.Fatalf("decode %s: %v", id, err)
		}
	}
}

func cardigannTestService(client HTTPDoer, id string, definition string) *Engine {
	service := New(client)
	service.UseLocalDefinitions(map[string]string{id: definition})
	return service
}

type fakeHTTPDoer func(req *http.Request) (*http.Response, error)

func (f fakeHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func response(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

const htmlCardigannDefinition = `
id: mockhtml
name: Mock HTML
type: public
links: [http://indexer.local/]
settings:
  - name: downloadlink
    type: select
    default: "magnet:"
caps:
  modes:
    search: [q]
    movie-search: [q]
search:
  paths:
    - path: "/search/{{ .Keywords }}"
  rows:
    selector: "tr.release"
  fields:
    title_default:
      selector: "a.title"
    title:
      text: "{{ .Result.title_default }}"
    details:
      selector: "a.title"
      attribute: href
    download:
      selector: "a.title"
      attribute: href
    size:
      selector: ".size"
    seeders:
      selector: ".seeders"
    leechers:
      selector: ".leechers"
download:
  selectors:
    - selector: "a[href^=\"{{ .Config.downloadlink }}\"]"
      attribute: href
`

const jsonCardigannDefinition = `
id: mockjson
name: Mock JSON
type: public
links: [http://indexer.local/]
caps:
  modes:
    search: [q]
    tv-search: [q, season, ep]
search:
  paths:
    - path: "/api"
      inputs:
        q: "{{ .Keywords }}"
      response:
        type: json
  rows:
    selector: "$"
  fields:
    title:
      selector: name
    infohash:
      selector: info_hash
    size:
      selector: size
    seeders:
      selector: seeders
    leechers:
      selector: leechers
    date:
      selector: added
`
