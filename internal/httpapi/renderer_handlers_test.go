package httpapi

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"media-manager/internal/config"
	"media-manager/internal/dlna"
	"media-manager/internal/dlna/controlpoint"
)

func TestRendererRoutesRequireAuthentication(t *testing.T) {
	router := chi.NewRouter()
	HandlerFromMux(NewServer(config.Config{}, nil, nil, nil, nil, nil, nil), router)
	request := httptest.NewRequest(http.MethodGet, "/renderers", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d body=%s", response.Code, response.Body.String())
	}
}

func TestRendererResponseMapsCapabilitiesAndLease(t *testing.T) {
	now := time.Now().UTC()
	response := rendererResponse(controlpoint.Renderer{
		UDN: "uuid:tv", FriendlyName: "Living room", SourceIP: net.ParseIP("192.168.1.2"),
		LastSeen: now, ExpiresAt: now.Add(time.Minute), RenderingControlURL: "http://192.168.1.2/render",
		AVTransportActionsKnown: true, RenderingActionsKnown: true,
		AVTransportActions: []string{"Play", "Pause", "Seek"}, RenderingActions: []string{"SetVolume"},
	})
	if !response.Supports.Play || !response.Supports.Seek || !response.Supports.Volume || response.Supports.Stop || response.Supports.Mute {
		t.Fatalf("response = %#v", response)
	}
	if upnpSeconds("1:02:03.500") != 3723 || upnpSeconds("NOT_IMPLEMENTED") != 0 || upnpSeconds("-1:00:00") != 0 {
		t.Fatal("UPnP time conversion failed")
	}

	unknownRenderer := controlpoint.Renderer{RenderingControlURL: "http://192.168.1.2/render"}
	unknown := rendererResponse(unknownRenderer)
	if !unknown.Supports.Play || !unknown.Supports.Volume {
		t.Fatalf("unknown SCPD should use optimistic fallback: %#v", unknown.Supports)
	}
	partial := rendererStatusResponse(unknownRenderer, controlpoint.RendererState{
		Transport: controlpoint.TransportState{State: "PLAYING", Status: "OK"},
		Position:  controlpoint.PositionInfo{RelTime: "0:00:12", TrackDuration: "0:02:00"},
		Volume:    controlpoint.VolumeState{Volume: 41}, VolumeAvailable: true,
	})
	if partial.Volume == nil || *partial.Volume != 41 || partial.Muted != nil || partial.PositionSeconds != 12 {
		t.Fatalf("partial status response = %#v", partial)
	}
}

func TestEveryRendererRouteRejectsUnauthenticatedRequests(t *testing.T) {
	server := NewServer(config.Config{}, nil, nil, nil, nil, nil, nil)
	router := chi.NewRouter()
	HandlerFromMux(server, router)
	requests := []struct {
		method, path, body string
	}{
		{http.MethodGet, "/renderers", ""},
		{http.MethodGet, "/renderers/uuid:tv/status", ""},
		{http.MethodPost, "/renderers/uuid:tv/play", `{}`},
		{http.MethodPost, "/renderers/uuid:tv/transport", `{}`},
		{http.MethodPost, "/renderers/uuid:tv/seek", `{}`},
		{http.MethodPost, "/renderers/uuid:tv/volume", `{}`},
		{http.MethodPost, "/renderers/uuid:tv/mute", `{}`},
	}
	for _, item := range requests {
		request := httptest.NewRequest(item.method, item.path, bytes.NewBufferString(item.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s status = %d", item.method, item.path, response.Code)
		}
	}
}

func TestAuthenticatedRendererOperationsHandleEmptyAndOfflineManagers(t *testing.T) {
	server := NewServer(config.Config{}, nil, nil, nil, nil, nil, nil)
	response := httptest.NewRecorder()
	server.listRenderers(response)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"renderers":[]`) {
		t.Fatalf("empty list = %d %s", response.Code, response.Body.String())
	}

	request := httptest.NewRequest(http.MethodGet, "/renderers/uuid:missing/status", nil)
	response = httptest.NewRecorder()
	server.getRendererStatus(response, request, "uuid:missing")
	if response.Code != http.StatusNotFound {
		t.Fatalf("nil status manager = %d", response.Code)
	}

	server.SetDLNAManager(dlna.NewManager(nil, "http://127.0.0.1:18080"))
	response = httptest.NewRecorder()
	server.listRenderers(response)
	if response.Code != http.StatusOK {
		t.Fatalf("manager list = %d", response.Code)
	}
	response = httptest.NewRecorder()
	server.getRendererStatus(response, request, "uuid:missing")
	if response.Code != http.StatusNotFound {
		t.Fatalf("offline status = %d", response.Code)
	}

	server.SetDLNAManager(nil)
	path := "/movie.mkv"
	response = httptest.NewRecorder()
	server.playMediaOnRenderer(response, request, "uuid:missing", RendererPlayRequest{MediaItemId: uuid.Nil, Path: &path})
	if response.Code != http.StatusNotFound {
		t.Fatalf("offline play = %d %s", response.Code, response.Body.String())
	}
}

func TestRendererControlReportsSuccessFailureAndMappedErrors(t *testing.T) {
	server := NewServer(config.Config{}, nil, nil, nil, nil, nil, nil)
	request := httptest.NewRequest(http.MethodPost, "/renderers/uuid:tv/transport", nil)
	if err := rendererControl(server, request, "uuid:tv", "play", func(context.Context) error { return nil }); err != nil {
		t.Fatal(err)
	}
	failed := errors.New("network failed")
	if err := rendererControl(server, request, "uuid:tv", "play", func(context.Context) error { return failed }); !errors.Is(err, failed) {
		t.Fatalf("control failure = %v", err)
	}

	cases := []struct {
		err  error
		code int
	}{
		{dlna.ErrRendererOffline, http.StatusNotFound},
		{dlna.ErrUnsupportedAction, http.StatusConflict},
		{dlna.ErrMediaNotPlayable, http.StatusBadRequest},
		{errors.Join(dlna.ErrInvalidRendererCommand, errors.New("bad value")), http.StatusBadRequest},
		{&controlpoint.FaultError{Action: "Play", StatusCode: 500, Description: "Transition unavailable"}, http.StatusBadGateway},
		{errors.New("timeout"), http.StatusBadGateway},
	}
	for _, item := range cases {
		response := httptest.NewRecorder()
		writeRendererCommand(response, item.err)
		if response.Code != item.code || !strings.Contains(response.Body.String(), "message") {
			t.Fatalf("%v => %d %s", item.err, response.Code, response.Body.String())
		}
	}
	success := httptest.NewRecorder()
	writeRendererCommand(success, nil)
	if success.Code != http.StatusNoContent {
		t.Fatalf("success status = %d", success.Code)
	}
}
