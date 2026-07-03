package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"media-manager/internal/acceptance"
	"media-manager/internal/downloadclients"
	"media-manager/internal/indexers"
	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func TestScenarioSCNMedia009AdminManagesDownloadActivityLifecycle(t *testing.T) {
	store := testSettingsStore(t)
	client := newAcceptanceClientWithStore(t, "SCN-MEDIA-009", store)
	activity := createScenarioDownloadActivity(t, store)

	var listed DownloadActivityListResponse
	client.doJSON(t, http.MethodGet, "/activity/downloads", nil, http.StatusOK, &listed)
	if !downloadActivityListHas(listed.Activities, activity.ID.String(), DownloadActivityStatusQueued) {
		t.Fatalf("download activity not listed: %#v", listed.Activities)
	}

	var cancelled DownloadActivity
	client.doJSON(t, http.MethodPost, "/activity/downloads/"+activity.ID.String()+"/cancel", nil, http.StatusOK, &cancelled)
	if cancelled.Status != DownloadActivityStatusCancelled || cancelled.ReleaseTitle != activity.ReleaseTitle {
		t.Fatalf("cancelled download activity = %#v", cancelled)
	}

	client.doJSON(t, http.MethodDelete, "/activity/downloads/"+activity.ID.String(), nil, http.StatusNoContent, nil)
	var afterDelete DownloadActivityListResponse
	client.doJSON(t, http.MethodGet, "/activity/downloads", nil, http.StatusOK, &afterDelete)
	if downloadActivityListHas(afterDelete.Activities, activity.ID.String(), DownloadActivityStatusCancelled) {
		t.Fatalf("deleted download activity still listed: %#v", afterDelete.Activities)
	}
}

func newAcceptanceClientWithStore(t *testing.T, scenarioID string, store *storage.SettingsStore) acceptanceClient {
	t.Helper()
	scenario, err := acceptance.RequireScenario("features/behavior", scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	if !scenario.HasTag("api") {
		t.Fatalf("%s missing @api tag", scenarioID)
	}
	router := chi.NewRouter()
	HandlerFromMux(NewServer(
		testConfig(),
		store,
		downloadclients.NewService(http.DefaultClient),
		indexers.NewService(http.DefaultClient),
		metadata.NewService(http.DefaultClient, nil),
		nil,
		nil,
	), router)
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

func createScenarioDownloadActivity(t *testing.T, store *storage.SettingsStore) storage.DownloadActivity {
	t.Helper()
	item, err := store.CreateMediaItem(t.Context(), storage.MediaItemInput{
		Type:                "movie",
		Title:               "Activity Scenario Movie",
		Year:                int32Ptr(2026),
		Monitored:           true,
		MonitorMode:         "onlyMedia",
		MinimumAvailability: "released",
	})
	if err != nil {
		t.Fatal(err)
	}
	activity, err := store.CreateDownloadActivity(t.Context(), storage.DownloadActivityInput{
		MediaItemID:        item.ID,
		ReleaseTitle:       "Activity.Scenario.Movie.2026.1080p",
		IndexerName:        "Scenario Indexer",
		DownloadClientName: "Scenario Client",
		DownloadURL:        "https://example.test/download/activity",
		Status:             "queued",
	})
	if err != nil {
		t.Fatal(err)
	}
	return activity
}

func downloadActivityListHas(activities []DownloadActivity, id string, status DownloadActivityStatus) bool {
	for _, activity := range activities {
		if activity.Id.String() == id && activity.Status == status {
			return true
		}
	}
	return false
}
