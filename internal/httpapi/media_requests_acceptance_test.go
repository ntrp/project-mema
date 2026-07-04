package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestScenarioSCNMedia006SignedInUsersCreateAndInspectMediaRequests(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-MEDIA-006")

	var created MediaRequest
	client.doJSON(t, http.MethodPost, "/media/requests", mediaRequestCreateRequest(), http.StatusCreated, &created)
	if created.Title != "Requested Scenario Movie" || created.Status != MediaRequestStatusPending {
		t.Fatalf("created media request = %#v", created)
	}

	var listed MediaRequestListResponse
	client.doJSON(t, http.MethodGet, "/media/requests", nil, http.StatusOK, &listed)
	if !mediaRequestListHas(listed.Requests, created.Id.String(), "Requested Scenario Movie") {
		t.Fatalf("media request not listed: %#v", listed.Requests)
	}

	var fetched MediaRequest
	client.doJSON(t, http.MethodGet, "/media/requests/"+created.Id.String(), nil, http.StatusOK, &fetched)
	if fetched.Id != created.Id || fetched.RequestedByUsername != "admin" {
		t.Fatalf("fetched media request = %#v", fetched)
	}

	var profiles MediaProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/profiles", nil, http.StatusOK, &profiles)
	if len(profiles.Profiles) == 0 {
		t.Fatal("expected seeded media profile")
	}
	libraryPath := filepath.Join(t.TempDir(), "requests-library")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	var folder LibraryFolderCreateResponse
	client.doJSON(t, http.MethodPost, "/settings/library/folders", LibraryFolderRequest{
		Path: libraryPath,
	}, http.StatusCreated, &folder)

	var approved MediaRequestApproveResponse
	client.doJSON(t, http.MethodPost, "/media/requests/"+created.Id.String()+"/approve", MediaRequestApproveRequest{
		QualityProfileId: profiles.Profiles[0].Id,
		LibraryFolderId:  folder.Folder.Id,
	}, http.StatusOK, &approved)
	if approved.Request.Status != MediaRequestStatusApproved {
		t.Fatalf("approved request = %#v", approved.Request)
	}
	if approved.MediaItem.Title != created.Title || !approved.MediaItem.Monitored {
		t.Fatalf("approved media item = %#v", approved.MediaItem)
	}
}

func mediaRequestCreateRequest() MediaRequestCreateRequest {
	tags := []string{"family", "uhd"}
	overview := "A requested movie waiting for approval."
	return MediaRequestCreateRequest{
		Type:                MediaTypeMovie,
		Title:               "Requested Scenario Movie",
		Year:                int32Ptr(2026),
		Overview:            &overview,
		MonitorMode:         OnlyMedia,
		MinimumAvailability: Released,
		Tags:                &tags,
	}
}

func mediaRequestListHas(requests []MediaRequest, id string, title string) bool {
	for _, request := range requests {
		if request.Id.String() == id && request.Title == title {
			return true
		}
	}
	return false
}
