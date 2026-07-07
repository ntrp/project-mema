package dlna

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"media-manager/internal/dlna/content"
	"media-manager/internal/dlna/soap"
	"media-manager/internal/dlna/ssdp"
	"media-manager/internal/storage"

	"github.com/google/uuid"
)

type contentFakeSource struct {
	items []storage.MediaItem
}

func (s contentFakeSource) ListMediaItems(context.Context) ([]storage.MediaItem, error) {
	return s.items, nil
}

func TestContentDirectoryBrowseSOAPReturnsDIDLAndCounts(t *testing.T) {
	path := "/media/Scenario.Movie.mkv"
	tree := content.NewTree(contentFakeSource{items: []storage.MediaItem{{
		ID:        uuid.New(),
		Type:      "movie",
		Title:     "Scenario Movie",
		FilePaths: []string{path},
	}}}).WithStat(controlFakeStat(path))
	dispatcher := soap.NewDispatcher()
	dispatcher.Register("/control", ssdp.ContentDir, contentDirectoryActions(tree, "http://127.0.0.1:18080"))
	body := browseEnvelope(content.EncodeID(content.RootContainerRef("movies")), string(content.BrowseDirectChildren), "0", "1")
	request := httptest.NewRequest("POST", "/control", strings.NewReader(body))
	request.Header.Set("SOAPACTION", `"urn:schemas-upnp-org:service:ContentDirectory:1#Browse"`)
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

func TestContentDirectoryBrowseInvalidObjectReturns701(t *testing.T) {
	tree := content.NewTree(contentFakeSource{})
	dispatcher := soap.NewDispatcher()
	dispatcher.Register("/control", ssdp.ContentDir, contentDirectoryActions(tree, "http://127.0.0.1:18080"))
	body := browseEnvelope("bad-object", string(content.BrowseDirectChildren), "0", "0")
	request := httptest.NewRequest("POST", "/control", strings.NewReader(body))
	request.Header.Set("SOAPACTION", `"urn:schemas-upnp-org:service:ContentDirectory:1#Browse"`)
	response := httptest.NewRecorder()

	dispatcher.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError || !strings.Contains(response.Body.String(), "<errorCode>701</errorCode>") {
		t.Fatalf("fault = %d %s", response.Code, response.Body.String())
	}
}

func browseEnvelope(objectID string, flag string, start string, count string) string {
	return `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body>` +
		`<u:Browse xmlns:u="urn:schemas-upnp-org:service:ContentDirectory:1">` +
		`<ObjectID>` + objectID + `</ObjectID><BrowseFlag>` + flag + `</BrowseFlag>` +
		`<Filter>*</Filter><StartingIndex>` + start + `</StartingIndex>` +
		`<RequestedCount>` + count + `</RequestedCount><SortCriteria></SortCriteria>` +
		`</u:Browse></s:Body></s:Envelope>`
}

func controlFakeStat(paths ...string) content.FileStatFunc {
	available := map[string]struct{}{}
	for _, path := range paths {
		available[path] = struct{}{}
	}
	return func(path string) (os.FileInfo, error) {
		if _, ok := available[path]; !ok {
			return nil, os.ErrNotExist
		}
		return controlFileInfo{name: filepath.Base(path)}, nil
	}
}

type controlFileInfo struct {
	name string
}

func (f controlFileInfo) Name() string       { return f.name }
func (f controlFileInfo) Size() int64        { return 1 }
func (f controlFileInfo) Mode() os.FileMode  { return 0 }
func (f controlFileInfo) ModTime() time.Time { return time.Time{} }
func (f controlFileInfo) IsDir() bool        { return false }
func (f controlFileInfo) Sys() any           { return nil }
