package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"media-manager/internal/metadata"
	"media-manager/internal/storage"
)

func TestAniDBPersonDetailsUseStoredLibraryAppearances(t *testing.T) {
	ctx := context.Background()
	store := testSettingsStore(t)
	clientName, clientVersion := "testclient", "1"
	provider, err := store.CreateMetadataProvider(ctx, storage.MetadataProviderInput{
		Name: "AniDB person test", Type: "anidb", BaseURL: "http://api.anidb.invalid:9001",
		APIKey: &clientName, PIN: &clientVersion, Enabled: true, Priority: 120,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := store.DeleteMetadataProvider(ctx, provider.ID); err != nil {
			t.Errorf("delete metadata provider: %v", err)
		}
	})

	providerType, animeID := "anidb", "9541"
	personID, role := "1678", "Eren Yeager"
	_, err = store.CreateMediaItem(ctx, storage.MediaItemInput{
		Type: "serie", Title: "Shingeki no Kyojin", ExternalProvider: &providerType, ExternalID: &animeID,
		MediaMetadataSnapshot: storage.MediaMetadataSnapshot{Cast: []storage.MediaPerson{{
			ExternalProvider: &providerType, ExternalID: &personID, Name: "Kaji Yuuki", Role: &role,
		}}},
	})
	if err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	HandlerFromMux(NewServer(testConfig(), store, nil, nil, metadata.NewService(http.DefaultClient, nil), nil, nil), router)
	request := httptest.NewRequest(http.MethodGet, "/people/anidb/1678", nil)
	request.AddCookie(authenticatedCookie(t, router))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%q", response.Code, response.Body.String())
	}
	var details PersonDetails
	if err := json.Unmarshal(response.Body.Bytes(), &details); err != nil {
		t.Fatal(err)
	}
	if details.Name != "Kaji Yuuki" || len(details.Appearances) != 1 || details.Appearances[0].Role == nil || *details.Appearances[0].Role != role {
		t.Fatalf("details = %#v", details)
	}
}
