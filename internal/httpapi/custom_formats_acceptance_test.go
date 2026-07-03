package httpapi

import (
	"net/http"
	"testing"
)

func TestScenarioSCNSettings017AdminManagesCustomFormatsAndParsing(t *testing.T) {
	client := newAcceptanceClient(t, "SCN-SETTINGS-017")
	request := customFormatRequest("Scenario Source", true)

	var created CustomFormat
	client.doJSON(t, http.MethodPost, "/settings/custom-formats", request, http.StatusCreated, &created)
	if created.Name != "Scenario Source" || len(created.IncludeSpecs) != 1 {
		t.Fatalf("created custom format = %#v", created)
	}

	updatedRequest := customFormatRequest("Scenario Source Updated", false)
	var updated CustomFormat
	client.doJSON(t, http.MethodPut, "/settings/custom-formats/"+created.Id.String(), updatedRequest, http.StatusOK, &updated)
	if updated.Name != "Scenario Source Updated" || updated.IncludeInRenameTemplate {
		t.Fatalf("updated custom format = %#v", updated)
	}

	var listed CustomFormatListResponse
	client.doJSON(t, http.MethodGet, "/settings/custom-formats", nil, http.StatusOK, &listed)
	if !customFormatListHas(listed.Formats, updated.Id.String(), "Scenario Source Updated") {
		t.Fatalf("custom format not listed: %#v", listed.Formats)
	}

	var parsed CustomFormatParsingResponse
	client.doJSON(t, http.MethodPost, "/settings/custom-formats/test-parsing", CustomFormatParsingRequest{
		FileName: "Scenario.Movie.2026.1080p.WEB-DL.x265-GRP.mkv",
	}, http.StatusOK, &parsed)
	if parsed.Details.MatchedSpecCount == 0 || len(parsed.MatchedCustomFormats) == 0 {
		t.Fatalf("expected custom format match: %#v", parsed)
	}

	client.doJSON(t, http.MethodDelete, "/settings/custom-formats/"+updated.Id.String(), nil, http.StatusNoContent, nil)
}

func customFormatRequest(name string, includeInRenameTemplate bool) CustomFormatRequest {
	return CustomFormatRequest{
		Name:                    name,
		IncludeInRenameTemplate: includeInRenameTemplate,
		IncludeSpecs: []CustomFormatSpec{{
			Id:       "source-webdl",
			Name:     "Source WEBDL",
			Type:     CustomFormatSpecTypeSource,
			Value:    "WEB.?DL",
			Required: true,
		}},
		ExcludeSpecs: []CustomFormatSpec{},
	}
}

func customFormatListHas(formats []CustomFormat, id string, name string) bool {
	for _, format := range formats {
		if format.Id.String() == id && format.Name == name {
			return true
		}
	}
	return false
}
