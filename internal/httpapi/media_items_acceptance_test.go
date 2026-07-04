package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestScenarioSCNMedia007AdminManagesMediaItemMonitoringOptions(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-MEDIA-007")

	var profiles MediaProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/profiles", nil, http.StatusOK, &profiles)
	if len(profiles.Profiles) == 0 {
		t.Fatal("expected seeded media profile")
	}
	libraryPath := filepath.Join(t.TempDir(), "library")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	var folder LibraryFolderCreateResponse
	client.doJSON(t, http.MethodPost, "/settings/library/folders", LibraryFolderRequest{
		Path: libraryPath,
	}, http.StatusCreated, &folder)

	var created MediaItem
	client.doJSON(t, http.MethodPost, "/media/items", mediaItemCreateRequest(profiles.Profiles[0].Id, folder.Folder.Id), http.StatusCreated, &created)
	if created.Title != "Managed Scenario Movie" || !created.Monitored || created.Status != Missing {
		t.Fatalf("created media item = %#v", created)
	}

	var listed MediaItemListResponse
	client.doJSON(t, http.MethodGet, "/media/items", nil, http.StatusOK, &listed)
	if !mediaItemListHas(listed.Items, created.Id.String(), "Managed Scenario Movie") {
		t.Fatalf("media item not listed: %#v", listed.Items)
	}

	notMonitored := false
	monitorMode := None
	var updated MediaItem
	client.doJSON(t, http.MethodPut, "/media/items/"+created.Id.String(), MediaItemUpdateRequest{
		Monitored:   &notMonitored,
		MonitorMode: &monitorMode,
	}, http.StatusOK, &updated)
	if updated.Monitored || updated.MonitorMode != None {
		t.Fatalf("updated media item = %#v", updated)
	}

	client.doJSON(t, http.MethodDelete, "/media/items/"+created.Id.String()+"?keepFiles=true", nil, http.StatusNoContent, nil)
	var afterDelete MediaItemListResponse
	client.doJSON(t, http.MethodGet, "/media/items", nil, http.StatusOK, &afterDelete)
	if mediaItemListHas(afterDelete.Items, created.Id.String(), "Managed Scenario Movie") {
		t.Fatalf("deleted media item still listed: %#v", afterDelete.Items)
	}
}

func mediaItemCreateRequest(qualityProfileID string, libraryFolderID ResourceId) MediaItemCreateRequest {
	tags := []string{"managed", "scenario"}
	overview := "A movie managed directly by the admin."
	return MediaItemCreateRequest{
		Type:                MediaTypeMovie,
		Title:               "Managed Scenario Movie",
		Year:                int32Ptr(2026),
		Overview:            &overview,
		Monitored:           true,
		MonitorMode:         OnlyMedia,
		MinimumAvailability: Released,
		QualityProfileId:    &qualityProfileID,
		LibraryFolderId:     &libraryFolderID,
		Tags:                &tags,
	}
}

func mediaItemListHas(items []MediaItem, id string, title string) bool {
	for _, item := range items {
		if item.Id.String() == id && item.Title == title {
			return true
		}
	}
	return false
}
