package httpapi

import (
	"net/http"
	"testing"
)

func TestScenarioSCNSettings018AdminManagesPersistedDownloadClients(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-018")

	var created DownloadClient
	client.doJSON(t, http.MethodPost, "/settings/download-clients", downloadClientRequest("Scenario Client"), http.StatusCreated, &created)
	if created.Name != "Scenario Client" || created.Type != Transmission {
		t.Fatalf("created download client = %#v", created)
	}

	var updated DownloadClient
	updateRequest := downloadClientRequest("Scenario Client Updated")
	updateRequest.Enabled = false
	client.doJSON(t, http.MethodPut, "/settings/download-clients/"+created.Id.String(), updateRequest, http.StatusOK, &updated)
	if updated.Name != "Scenario Client Updated" || updated.Enabled {
		t.Fatalf("updated download client = %#v", updated)
	}

	var listed DownloadClientListResponse
	client.doJSON(t, http.MethodGet, "/settings/download-clients", nil, http.StatusOK, &listed)
	if !downloadClientListHas(listed.Clients, updated.Id.String(), "Scenario Client Updated") {
		t.Fatalf("download client not listed: %#v", listed.Clients)
	}

	var result IntegrationTestResponse
	client.doJSON(t, http.MethodPost, "/settings/download-clients/"+updated.Id.String()+"/test", nil, http.StatusOK, &result)
	if result.Success || result.Message == "" {
		t.Fatalf("unreachable local client should fail validation with a message: %#v", result)
	}

	client.doJSON(t, http.MethodDelete, "/settings/download-clients/"+updated.Id.String(), nil, http.StatusNoContent, nil)
}

func downloadClientRequest(name string) DownloadClientRequest {
	category := "movies"
	return DownloadClientRequest{
		Name:     name,
		Type:     Transmission,
		BaseUrl:  "http://127.0.0.1:1",
		Category: &category,
		Enabled:  true,
		Priority: 10,
	}
}

func downloadClientListHas(clients []DownloadClient, id string, name string) bool {
	for _, client := range clients {
		if client.Id.String() == id && client.Name == name {
			return true
		}
	}
	return false
}
