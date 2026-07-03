package httpapi

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestScenarioSCNSettings016AdminManagesLibraryFoldersAndMappings(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-016")
	root := t.TempDir()
	libraryPath := filepath.Join(root, "library")
	if err := os.MkdirAll(filepath.Join(libraryPath, "Movies"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "Scenario.Movie.2026.mkv"), []byte("video"), 0o644); err != nil {
		t.Fatal(err)
	}

	var options LibraryFolderOptionListResponse
	client.doJSON(t, http.MethodGet, "/settings/library/folder-options?path="+libraryPath, nil, http.StatusOK, &options)
	if options.CurrentPath != libraryPath || !libraryOptionHas(options.Entries, "Movies") {
		t.Fatalf("folder options = %#v", options)
	}

	var createdOption LibraryFolderOption
	client.doJSON(t, http.MethodPost, "/settings/library/folder-options", LibraryFolderOptionCreateRequest{
		ParentPath: libraryPath,
		Name:       "Imported",
	}, http.StatusCreated, &createdOption)
	if createdOption.Name != "Imported" {
		t.Fatalf("created folder option = %#v", createdOption)
	}

	var created LibraryFolderCreateResponse
	client.doJSON(t, http.MethodPost, "/settings/library/folders", LibraryFolderRequest{
		Path: libraryPath,
	}, http.StatusCreated, &created)
	if created.Folder.Path != libraryPath || created.Scan.TotalFiles == 0 {
		t.Fatalf("created library folder = %#v", created)
	}

	var listed LibraryFolderListResponse
	client.doJSON(t, http.MethodGet, "/settings/library/folders", nil, http.StatusOK, &listed)
	if !libraryFolderHas(listed.Folders, uuid.UUID(created.Folder.Id)) {
		t.Fatalf("folder not listed: %#v", listed.Folders)
	}

	var scan LibraryScan
	client.doJSON(t, http.MethodGet, "/settings/library/scans/"+created.Scan.Id.String(), nil, http.StatusOK, &scan)
	if scan.Id != created.Scan.Id || len(scan.Items) == 0 {
		t.Fatalf("scan = %#v", scan)
	}

	var profiles MediaProfileListResponse
	client.doJSON(t, http.MethodGet, "/settings/profiles", nil, http.StatusOK, &profiles)
	if len(profiles.Profiles) == 0 {
		t.Fatal("expected seeded media profile")
	}
	var matched LibraryScanItemMatchResponse
	client.doJSON(t, http.MethodPost, "/settings/library/scans/"+created.Scan.Id.String()+"/items/"+scan.Items[0].Id.String()+"/match", LibraryScanItemMatchRequest{
		MediaKind:           LibraryMediaKindMovie,
		Title:               "Scenario Movie",
		Year:                int32Ptr(2026),
		Monitored:           true,
		QualityProfileId:    profiles.Profiles[0].Id,
		MonitorMode:         OnlyMedia,
		MinimumAvailability: Released,
		ExternalProvider:    stringPtr("tmdb"),
		ExternalId:          stringPtr("scenario-movie"),
		Overview:            stringPtr("A movie imported from a scanned library folder."),
		PosterPath:          stringPtr("/poster.jpg"),
	}, http.StatusOK, &matched)
	if matched.Item.Status != LibraryScanItemStatusManuallyAdded || matched.MediaItem.Title != "Scenario Movie" {
		t.Fatalf("matched scan item = %#v", matched)
	}

	var mapping PathMapping
	client.doJSON(t, http.MethodPost, "/settings/library/path-mappings", PathMappingRequest{
		ClientPath: "/downloads",
		AppPath:    root,
	}, http.StatusCreated, &mapping)
	if mapping.ClientPath != "/downloads" || mapping.AppPath != root {
		t.Fatalf("mapping = %#v", mapping)
	}

	var mappings PathMappingListResponse
	client.doJSON(t, http.MethodGet, "/settings/library/path-mappings", nil, http.StatusOK, &mappings)
	if len(mappings.Mappings) == 0 {
		t.Fatalf("mappings = %#v", mappings)
	}
	client.doJSON(t, http.MethodDelete, "/settings/library/path-mappings/"+mapping.Id.String(), nil, http.StatusNoContent, nil)
	client.doJSON(t, http.MethodDelete, "/settings/library/folders/"+created.Folder.Id.String(), nil, http.StatusNoContent, nil)
}

func libraryOptionHas(options []LibraryFolderOption, name string) bool {
	for _, option := range options {
		if option.Name == name {
			return true
		}
	}
	return false
}

func int32Ptr(value int32) *int32 {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func libraryFolderHas(folders []LibraryFolder, id uuid.UUID) bool {
	for _, folder := range folders {
		if uuid.UUID(folder.Id) == id {
			return true
		}
	}
	return false
}
