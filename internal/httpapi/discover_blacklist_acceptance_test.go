package httpapi

import (
	"net/http"
	"testing"

	"media-manager/internal/storage"
)

func TestScenarioSCNMedia005AdminManagesDiscoveryBlacklist(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-MEDIA-005")

	var created DiscoverBlacklistItem
	client.doJSON(t, http.MethodPost, "/media/discover/blacklist", discoverBlacklistRequest(), http.StatusCreated, &created)
	if created.Title != "Scenario Movie" || created.Type != MediaTypeMovie {
		t.Fatalf("created discover blacklist item = %#v", created)
	}

	var listed DiscoverBlacklistResponse
	client.doJSON(t, http.MethodGet, "/media/discover/blacklist", nil, http.StatusOK, &listed)
	if !discoverBlacklistListHas(listed.Items, created.Id.String()) {
		t.Fatalf("blacklist item not listed: %#v", listed.Items)
	}

	filtered := filterDiscoverBlacklist([]MediaSearchResult{
		{Type: MediaTypeMovie, Title: "Scenario Movie", Year: int32Ptr(2026)},
		{Type: MediaTypeMovie, Title: "Visible Movie", Year: int32Ptr(2027)},
	}, []storage.DiscoverBlacklistItem{{
		Type:  string(MediaTypeMovie),
		Title: "Scenario Movie",
		Year:  int32Ptr(2026),
	}})
	if len(filtered) != 1 || filtered[0].Title != "Visible Movie" {
		t.Fatalf("filtered results = %#v", filtered)
	}

	client.doJSON(t, http.MethodDelete, "/media/discover/blacklist/"+created.Id.String(), nil, http.StatusNoContent, nil)
}

func discoverBlacklistRequest() DiscoverBlacklistRequest {
	provider := "tmdb"
	externalID := "scenario-movie"
	return DiscoverBlacklistRequest{
		Type:             MediaTypeMovie,
		Title:            "Scenario Movie",
		Year:             int32Ptr(2026),
		ExternalProvider: &provider,
		ExternalId:       &externalID,
	}
}

func discoverBlacklistListHas(items []DiscoverBlacklistItem, id string) bool {
	for _, item := range items {
		if item.Id.String() == id {
			return true
		}
	}
	return false
}
