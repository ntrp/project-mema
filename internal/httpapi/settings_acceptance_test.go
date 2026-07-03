package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"media-manager/internal/acceptance"
	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/metadata"
	"media-manager/internal/testmocks"
)

func TestScenarioSCNSettings002AdminManagesTags(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-002")
	name := "tag-" + uuid.NewString()
	renamed := name + "-renamed"

	var created Tag
	client.doJSON(t, http.MethodPost, "/settings/tags", TagRequest{Name: name}, http.StatusCreated, &created)
	if created.Name != name {
		t.Fatalf("created tag name = %q", created.Name)
	}

	var updated Tag
	client.doJSON(t, http.MethodPut, "/settings/tags/"+created.Id.String(), TagRequest{Name: renamed}, http.StatusOK, &updated)
	if updated.Name != renamed {
		t.Fatalf("updated tag name = %q", updated.Name)
	}

	var listed TagListResponse
	client.doJSON(t, http.MethodGet, "/settings/tags", nil, http.StatusOK, &listed)
	if !tagListHas(listed.Tags, updated.Id, renamed) {
		t.Fatalf("updated tag not listed: %#v", listed.Tags)
	}

	client.doJSON(t, http.MethodDelete, "/settings/tags/"+updated.Id.String(), nil, http.StatusNoContent, nil)
	var afterDelete TagListResponse
	client.doJSON(t, http.MethodGet, "/settings/tags", nil, http.StatusOK, &afterDelete)
	if tagListHas(afterDelete.Tags, updated.Id, renamed) {
		t.Fatalf("deleted tag still listed: %#v", afterDelete.Tags)
	}
}

func TestScenarioSCNSettings003AdminManagesLanguageAliases(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-003")
	code := "x" + uuid.NewString()[:7]
	displayName := "Scenario Language " + code
	updatedDisplayName := "Updated Scenario Language " + code
	createdAliases := []string{"original", "source"}
	updatedAliases := []string{"updated", "alias"}

	var created Language
	client.doJSON(t, http.MethodPost, "/settings/languages", LanguageRequest{
		Code:        code,
		DisplayName: displayName,
		Aliases:     createdAliases,
	}, http.StatusCreated, &created)
	if created.Code == "" || created.DisplayName != displayName || !containsAll(created.Aliases, createdAliases) {
		t.Fatalf("created language = %#v", created)
	}

	var updated Language
	client.doJSON(t, http.MethodPut, "/settings/languages/"+created.Code, LanguageUpdateRequest{
		DisplayName: updatedDisplayName,
		Aliases:     updatedAliases,
	}, http.StatusOK, &updated)
	if updated.DisplayName != updatedDisplayName || !containsAll(updated.Aliases, updatedAliases) {
		t.Fatalf("updated language = %#v", updated)
	}

	var listed LanguageListResponse
	client.doJSON(t, http.MethodGet, "/settings/languages", nil, http.StatusOK, &listed)
	if !languageListHas(listed.Languages, created.Code) {
		t.Fatalf("updated language not listed: %#v", listed.Languages)
	}

	client.doJSON(t, http.MethodDelete, "/settings/languages/"+created.Code, nil, http.StatusNoContent, nil)
	var afterDelete LanguageListResponse
	client.doJSON(t, http.MethodGet, "/settings/languages", nil, http.StatusOK, &afterDelete)
	if languageListHas(afterDelete.Languages, created.Code) {
		t.Fatalf("deleted language still listed: %#v", afterDelete.Languages)
	}
}

func TestScenarioSCNSettings004AdminManagesIndexerConfiguration(t *testing.T) {
	provider := testmocks.NewProviderServer()
	t.Cleanup(provider.Close)
	client := newAcceptanceClientWithProviders(t, "SCN-SETTINGS-004", provider)
	categories := []int32{2000, 2040}

	var created Indexer
	client.doJSON(t, http.MethodPost, "/settings/indexers", IndexerRequest{
		DefinitionId: "generic-torznab",
		Name:         "Scenario Torznab",
		BaseUrl:      provider.URL + "/torznab/api",
		Enabled:      true,
		Priority:     10,
		Categories:   &categories,
	}, http.StatusCreated, &created)
	if created.Name != "Scenario Torznab" || created.Protocol != IndexerProtocolTorrent || !created.Enabled {
		t.Fatalf("created indexer = %#v", created)
	}

	var updated Indexer
	client.doJSON(t, http.MethodPut, "/settings/indexers/"+created.Id.String(), IndexerRequest{
		DefinitionId: "generic-torznab",
		Name:         "Updated Torznab",
		BaseUrl:      provider.URL + "/torznab/api",
		Enabled:      false,
		Priority:     20,
		Categories:   &categories,
	}, http.StatusOK, &updated)
	if updated.Name != "Updated Torznab" || updated.Enabled {
		t.Fatalf("updated indexer = %#v", updated)
	}

	var listed IndexerListResponse
	client.doJSON(t, http.MethodGet, "/settings/indexers", nil, http.StatusOK, &listed)
	if !indexerListHas(listed.Indexers, updated.Id, "Updated Torznab") {
		t.Fatalf("updated indexer not listed: %#v", listed.Indexers)
	}

	var result IntegrationTestResponse
	client.doJSON(t, http.MethodPost, "/settings/indexers/"+updated.Id.String()+"/test", nil, http.StatusOK, &result)
	if !result.Success {
		t.Fatalf("indexer test failed: %#v", result)
	}

	client.doJSON(t, http.MethodDelete, "/settings/indexers/"+updated.Id.String(), nil, http.StatusNoContent, nil)
}

func TestScenarioSCNSettings005AdminValidatesDownloadClientConfig(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-005")
	apiKey := "scenario-key"

	var result IntegrationTestResponse
	client.doJSON(t, http.MethodPost, "/settings/download-clients/test", DownloadClientRequest{
		Name:     "Scenario SABnzbd",
		Type:     "sabnzbd",
		BaseUrl:  "http://127.0.0.1:1",
		ApiKey:   &apiKey,
		Enabled:  true,
		Priority: 1,
	}, http.StatusOK, &result)
	if result.Success {
		t.Fatalf("unreachable local client should fail validation: %#v", result)
	}
	if result.Message == "" {
		t.Fatalf("expected validation failure message: %#v", result)
	}
}

type acceptanceClient struct {
	router http.Handler
	cookie *http.Cookie
}

func newAcceptanceClient(t *testing.T, scenarioID string) acceptanceClient {
	return newAcceptanceClientWithProviders(t, scenarioID, nil)
}

func newAcceptanceClientWithProviders(t *testing.T, scenarioID string, provider *testmocks.ProviderServer) acceptanceClient {
	t.Helper()
	scenario, err := acceptance.RequireScenario("features/behavior", scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	if !scenario.HasTag("api") {
		t.Fatalf("%s missing @api tag", scenarioID)
	}
	router := authRouterWithProviders(t, provider)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, loginRequest("admin", "admin"))
	if response.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %q", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("login response did not include a session cookie")
	}
	return acceptanceClient{router: router, cookie: cookies[0]}
}

func (c acceptanceClient) doJSON(t *testing.T, method string, path string, body any, wantStatus int, output any) {
	t.Helper()
	var payload *bytes.Reader
	if body == nil {
		payload = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		payload = bytes.NewReader(data)
	}
	request := httptest.NewRequest(method, path, payload)
	request.AddCookie(c.cookie)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	c.router.ServeHTTP(response, request)
	if response.Code != wantStatus {
		t.Fatalf("%s %s status = %d, want %d, body = %q", method, path, response.Code, wantStatus, response.Body.String())
	}
	if output != nil {
		if err := json.Unmarshal(response.Body.Bytes(), output); err != nil {
			t.Fatalf("decode %s %s: %v; body = %q", method, path, err, response.Body.String())
		}
	}
}

func authRouterWithProviders(t *testing.T, provider *testmocks.ProviderServer) http.Handler {
	t.Helper()
	httpClient := http.DefaultClient
	if provider != nil {
		httpClient = provider.Client()
	}
	router := chi.NewRouter()
	HandlerFromMux(NewServer(
		testConfig(),
		testSettingsStore(t),
		downloadclients.NewService(httpClient),
		indexers.NewService(httpClient),
		metadata.NewService(httpClient, nil),
		nil,
		nil,
	), router)
	if provider == nil {
		return router
	}
	return router
}

func tagListHas(tags []Tag, id uuid.UUID, name string) bool {
	for _, tag := range tags {
		if uuid.UUID(tag.Id) == id && tag.Name == name {
			return true
		}
	}
	return false
}

func languageListHas(languages []Language, code string) bool {
	for _, language := range languages {
		if language.Code == code {
			return true
		}
	}
	return false
}

func indexerListHas(indexers []Indexer, id uuid.UUID, name string) bool {
	for _, indexer := range indexers {
		if uuid.UUID(indexer.Id) == id && indexer.Name == name {
			return true
		}
	}
	return false
}

func containsAll(values []string, expected []string) bool {
	for _, want := range expected {
		found := false
		for _, value := range values {
			if value == want {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
