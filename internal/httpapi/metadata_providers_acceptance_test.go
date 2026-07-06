package httpapi

import (
	"net/http"
	"testing"

	"github.com/google/uuid"

	"media-manager/internal/testmocks"
)

func TestScenarioSCNSettings022AdminManagesMetadataProviders(t *testing.T) {
	provider := testmocks.NewProviderServer()
	t.Cleanup(provider.Close)
	client := newAcceptanceClientWithProviders(t, "SCN-SETTINGS-022", provider)
	apiKey := "scenario-metadata-key"

	var created MetadataProvider
	client.doJSON(t, http.MethodPost, "/settings/metadata-providers", MetadataProviderRequest{
		Name:     "Scenario TMDb",
		Type:     Tmdb,
		BaseUrl:  provider.URL + "/tmdb/3",
		ApiKey:   &apiKey,
		Enabled:  true,
		Priority: 10,
	}, http.StatusCreated, &created)
	if created.Name != "Scenario TMDb" || created.Type != Tmdb || !created.Enabled || !created.ApiKeySet {
		t.Fatalf("created metadata provider = %#v", created)
	}
	if created.ApiKey == nil || *created.ApiKey != apiKey {
		t.Fatalf("created metadata provider = %#v", created)
	}

	var updated MetadataProvider
	client.doJSON(t, http.MethodPut, "/settings/metadata-providers/"+created.Id.String(), MetadataProviderRequest{
		Name:     "Updated TMDb",
		Type:     Tmdb,
		BaseUrl:  provider.URL + "/tmdb/3",
		ApiKey:   &apiKey,
		Enabled:  false,
		Priority: 20,
	}, http.StatusOK, &updated)
	if updated.Name != "Updated TMDb" || updated.Enabled || updated.Priority != 20 {
		t.Fatalf("updated metadata provider = %#v", updated)
	}
	if updated.ApiKey == nil || *updated.ApiKey != apiKey {
		t.Fatalf("updated metadata provider missing api key = %#v", updated)
	}

	var listed MetadataProviderListResponse
	client.doJSON(t, http.MethodGet, "/settings/metadata-providers", nil, http.StatusOK, &listed)
	if !metadataProviderListHas(listed.Providers, updated.Id, "Updated TMDb") {
		t.Fatalf("updated metadata provider not listed: %#v", listed.Providers)
	}

	var result IntegrationTestResponse
	client.doJSON(t, http.MethodPost, "/settings/metadata-providers/"+updated.Id.String()+"/test", nil, http.StatusOK, &result)
	if !result.Success || result.Message == "" {
		t.Fatalf("metadata provider test failed: %#v", result)
	}

	client.doJSON(t, http.MethodDelete, "/settings/metadata-providers/"+updated.Id.String(), nil, http.StatusNoContent, nil)
}

func metadataProviderListHas(providers []MetadataProvider, id uuid.UUID, name string) bool {
	for _, provider := range providers {
		if uuid.UUID(provider.Id) == id && provider.Name == name {
			return true
		}
	}
	return false
}
