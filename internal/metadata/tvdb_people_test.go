package metadata

import (
	"context"
	"net/http"
	"testing"
)

func TestTVDBPersonDetailsMapsExtendedRecord(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := ""
		status := http.StatusOK
		switch r.URL.Path {
		case "/login":
			body = `{"data":{"token":"tvdb-token"}}`
		case "/people/4001/extended":
			if r.URL.Query().Get("meta") != "translations" {
				t.Fatalf("meta query = %q, want translations", r.URL.Query().Get("meta"))
			}
			body = `{
				"data": {
					"id": 4001,
					"name": "Example Actor",
					"image": "https://artworks.thetvdb.com/actor.jpg",
					"birth": "1970-01-02",
					"birthPlace": "Scenario City",
					"death": "",
					"aliases": [{"name": "E. Actor"}],
					"biographies": [
						{"language": "fra", "biography": "Biographie francaise."},
						{"language": "eng", "biography": "English biography."}
					],
					"characters": [
						{
							"name": "Robot",
							"movieId": 516,
							"movie": {"name": "WALL-E", "year": "2008", "image": "https://artworks.thetvdb.com/walle.jpg"}
						},
						{
							"name": "Commander",
							"seriesId": 121361,
							"series": {"name": "Game of Thrones", "year": "2011", "image": "https://artworks.thetvdb.com/got.jpg"}
						}
					]
				}
			}`
		default:
			status = http.StatusNotFound
		}
		return tvdbTestResponse(status, body), nil
	})}

	apiKey := "test-key"
	service := NewService(client, nil)
	details, err := service.PersonDetails(context.Background(), Config{
		Type:    "tvdb",
		BaseURL: "https://metadata.test",
		APIKey:  &apiKey,
	}, "4001")
	if err != nil {
		t.Fatalf("PersonDetails() error = %v", err)
	}
	if details.ID != "4001" || details.Name != "Example Actor" {
		t.Fatalf("details identity = %#v", details)
	}
	if details.Biography == nil || *details.Biography != "English biography." {
		t.Fatalf("biography = %#v", details.Biography)
	}
	if details.Birthday == nil || *details.Birthday != "1970-01-02" || details.PlaceOfBirth == nil || *details.PlaceOfBirth != "Scenario City" {
		t.Fatalf("birth fields = %#v %#v", details.Birthday, details.PlaceOfBirth)
	}
	if details.ProfilePath == nil || *details.ProfilePath != "https://artworks.thetvdb.com/actor.jpg" {
		t.Fatalf("profile = %#v", details.ProfilePath)
	}
	if len(details.AlsoKnownAs) != 1 || details.AlsoKnownAs[0] != "E. Actor" {
		t.Fatalf("aliases = %#v", details.AlsoKnownAs)
	}
	if len(details.Appearances) != 2 {
		t.Fatalf("appearances = %#v", details.Appearances)
	}
	if details.Appearances[0].Type != "movie" || details.Appearances[0].ExternalID != "516" || details.Appearances[0].Title != "WALL-E" {
		t.Fatalf("movie appearance = %#v", details.Appearances[0])
	}
	if details.Appearances[1].Type != "serie" || details.Appearances[1].ExternalID != "121361" || details.Appearances[1].Title != "Game of Thrones" {
		t.Fatalf("series appearance = %#v", details.Appearances[1])
	}
}
