package dlna

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"media-manager/internal/dlna/content"
	"media-manager/internal/dlna/soap"
	"media-manager/internal/dlna/ssdp"
	"media-manager/internal/storage"

	"github.com/google/uuid"
)

func TestContentDirectorySearchSOAPReturnsDIDLAndCounts(t *testing.T) {
	path := "/media/Scenario.Movie.mkv"
	tree := content.NewTree(contentFakeSource{items: []storage.MediaItem{{
		ID:        uuid.New(),
		Type:      "movie",
		Title:     "Scenario Movie",
		FilePaths: []string{path},
	}}}).WithStat(controlFakeStat(path))
	dispatcher := soapDispatcherWithContent(tree)
	body := searchEnvelope(content.RootID, `dc:title contains "Scenario"`, "0", "10")
	request := soapRequest(
		"/dlna/control/content-directory",
		"urn:schemas-upnp-org:service:ContentDirectory:1#Search",
		body,
	)
	response := httptest.NewRecorder()

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
	for _, want := range []string{
		"<NumberReturned>1</NumberReturned>",
		"<TotalMatches>1</TotalMatches>",
		"&lt;DIDL-Lite",
		"Scenario Movie",
	} {
		if !strings.Contains(response.Body.String(), want) {
			t.Fatalf("SOAP response missing %q:\n%s", want, response.Body.String())
		}
	}
}

func TestContentDirectorySearchUnsupportedCriteriaReturns402(t *testing.T) {
	dispatcher := soapDispatcherWithContent(content.NewTree(contentFakeSource{}))
	body := searchEnvelope(content.RootID, `res@duration > "0:01:00"`, "0", "10")
	request := soapRequest(
		"/dlna/control/content-directory",
		"urn:schemas-upnp-org:service:ContentDirectory:1#Search",
		body,
	)
	response := httptest.NewRecorder()

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>402</errorCode>") {
		t.Fatalf("fault = %d %s", response.Code, response.Body.String())
	}
}

func searchEnvelope(containerID string, criteria string, start string, count string) string {
	return `<u:Search xmlns:u="urn:schemas-upnp-org:service:ContentDirectory:1">` +
		`<ContainerID>` + containerID + `</ContainerID><SearchCriteria>` + criteria + `</SearchCriteria>` +
		`<Filter>*</Filter><StartingIndex>` + start + `</StartingIndex>` +
		`<RequestedCount>` + count + `</RequestedCount><SortCriteria></SortCriteria></u:Search>`
}

func soapDispatcherWithContent(tree *content.Tree) *soap.Dispatcher {
	dispatcher := soap.NewDispatcher()
	dispatcher.Register("/dlna/control/content-directory", ssdp.ContentDir, contentDirectoryActions(tree, "http://127.0.0.1:18080", func() int { return 0 }))
	return dispatcher
}
