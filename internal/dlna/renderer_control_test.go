package dlna

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/dlna/controlpoint"
	"media-manager/internal/storage"
)

func TestPlayOnRendererUsesDLNAPlannerURIAndDIDL(t *testing.T) {
	var requests []string
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		requests = append(requests, r.Header.Get("SOAPACTION")+"\n"+string(body))
		_, _ = w.Write([]byte(`<Envelope><Body/></Envelope>`))
	}))
	defer fake.Close()

	mediaID := uuid.New()
	path := filepath.Join(t.TempDir(), "movie.mkv")
	if err := os.WriteFile(path, []byte("test media"), 0o600); err != nil {
		t.Fatal(err)
	}
	manager := NewManager(nil, "http://127.0.0.1:18080")
	manager.source = contentFakeSource{items: []storage.MediaItem{{
		ID: mediaID, Type: "movie", Title: "Planner Movie", FilePaths: []string{path},
	}}}
	_, network, _ := net.ParseCIDR("127.0.0.0/8")
	manager.controlPoint = controlpoint.NewService(nil, []*net.IPNet{network}, fake.Client())
	manager.controlPoint.Registry().Upsert(controlpoint.Renderer{
		UDN: "uuid:fake", FriendlyName: "Fake Renderer", AVTransportURL: fake.URL,
		SourceIP: net.ParseIP("127.0.0.1"), LocalIP: net.ParseIP("127.0.0.1"),
		AVTransportActionsKnown: true,
		AVTransportActions:      []string{"SetAVTransportURI", "Play"},
	}, net.ParseIP("127.0.0.1"), time.Minute)

	playback, err := manager.PlayOnRenderer(context.Background(), "uuid:fake", mediaID, path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(playback.URI, "http://127.0.0.1:18080/dlna/resource/") {
		t.Fatalf("playback URI = %q", playback.URI)
	}
	if !strings.Contains(playback.Metadata, "DIDL-Lite") || !strings.Contains(playback.Metadata, "movie.mkv") {
		t.Fatalf("metadata = %s", playback.Metadata)
	}
	if len(requests) != 2 || !strings.Contains(requests[0], "SetAVTransportURI") || !strings.Contains(requests[0], "CurrentURIMetaData") || !strings.Contains(requests[1], "#Play") {
		t.Fatalf("requests = %#v", requests)
	}
}

func TestRendererStateReturnsAvailableReadsWhenOptionalActionsFail(t *testing.T) {
	manager, requests := rendererManager(t, func(w http.ResponseWriter, r *http.Request) {
		action := r.Header.Get("SOAPACTION")
		switch {
		case strings.Contains(action, "GetTransportInfo"):
			_, _ = w.Write([]byte(`<Envelope><Body><Response><CurrentTransportState>PLAYING</CurrentTransportState><CurrentTransportStatus>OK</CurrentTransportStatus></Response></Body></Envelope>`))
		case strings.Contains(action, "GetVolume"):
			_, _ = w.Write([]byte(`<Envelope><Body><Response><CurrentVolume>41</CurrentVolume></Response></Body></Envelope>`))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`<Envelope><Body><Fault><detail><UPnPError><errorCode>401</errorCode><errorDescription>Invalid Action</errorDescription></UPnPError></detail></Fault></Body></Envelope>`))
		}
	})
	renderer, _ := manager.Renderer("uuid:fake")
	renderer.RenderingActionsKnown = false // Unknown SCPD: probe optional reads and retain successful ones.
	manager.controlPoint.Registry().Upsert(renderer, renderer.SourceIP, time.Minute)

	state, err := manager.RendererState(context.Background(), "uuid:fake")
	if err != nil || !state.TransportAvailable || !state.VolumeAvailable || state.PositionAvailable || state.MuteAvailable {
		t.Fatalf("partial state = %#v, %v", state, err)
	}
	if state.Transport.State != "PLAYING" || state.Volume.Volume != 41 || !containsRequest(*requests, "GetMute") {
		t.Fatalf("state/requests = %#v / %#v", state, *requests)
	}
}

func TestRendererStateTreatsInvalidActionFaultsAsUnsupported(t *testing.T) {
	manager, _ := rendererManager(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`<Envelope><Body><Fault><detail><UPnPError><errorCode>401</errorCode><errorDescription>Invalid Action</errorDescription></UPnPError></detail></Fault></Body></Envelope>`))
	})
	renderer, _ := manager.Renderer("uuid:fake")
	renderer.AVTransportActionsKnown = false
	renderer.RenderingActionsKnown = false
	manager.controlPoint.Registry().Upsert(renderer, renderer.SourceIP, time.Minute)
	state, err := manager.RendererState(context.Background(), "uuid:fake")
	if err != nil || state.TransportAvailable || state.PositionAvailable || state.VolumeAvailable || state.MuteAvailable {
		t.Fatalf("unsupported read state = %#v, %v", state, err)
	}
}

func TestRendererStateErrorsOnlyWhenEveryAttemptFails(t *testing.T) {
	manager, _ := rendererManager(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	renderer, _ := manager.controlPoint.Registry().Get("uuid:fake", time.Time{})
	renderer.AVTransportActionsKnown = false
	renderer.RenderingActionsKnown = false
	manager.controlPoint.Registry().Upsert(renderer, renderer.SourceIP, time.Minute)
	if _, err := manager.RendererState(context.Background(), "uuid:fake"); err == nil {
		t.Fatal("all failed status reads should report an error")
	}
}

func TestRendererCommandsValidateCapabilitiesAndInputs(t *testing.T) {
	manager, requests := rendererManager(t, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`<Envelope><Body/></Envelope>`))
	})
	ctx := context.Background()
	for _, action := range []string{"play", "pause", "stop", "next", "previous"} {
		if err := manager.Transport(ctx, "uuid:fake", action); err != nil {
			t.Fatalf("%s: %v", action, err)
		}
	}
	if err := manager.Transport(ctx, "uuid:fake", "rewind"); !errors.Is(err, ErrUnsupportedAction) {
		t.Fatalf("unknown transport error = %v", err)
	}
	if err := manager.SeekRenderer(ctx, "uuid:fake", -1); !errors.Is(err, ErrInvalidRendererCommand) {
		t.Fatalf("negative seek error = %v", err)
	}
	if err := manager.SetRendererVolume(ctx, "uuid:fake", 101); !errors.Is(err, ErrInvalidRendererCommand) {
		t.Fatalf("invalid volume error = %v", err)
	}
	if err := manager.SeekRenderer(ctx, "uuid:fake", 65); err != nil {
		t.Fatal(err)
	}
	if err := manager.SetRendererVolume(ctx, "uuid:fake", 25); err != nil {
		t.Fatal(err)
	}
	if err := manager.SetRendererMute(ctx, "uuid:fake", true); err != nil {
		t.Fatal(err)
	}
	if !containsRequest(*requests, "0:01:05") || !containsRequest(*requests, "DesiredVolume") || !containsRequest(*requests, "DesiredMute") {
		t.Fatalf("command requests = %#v", *requests)
	}

	renderer, _ := manager.Renderer("uuid:fake")
	renderer.AVTransportActions = []string{"Play"}
	renderer.RenderingActions = nil
	manager.controlPoint.Registry().Upsert(renderer, renderer.SourceIP, time.Minute)
	if err := manager.Transport(ctx, "uuid:fake", "pause"); !errors.Is(err, ErrUnsupportedAction) {
		t.Fatalf("unsupported pause error = %v", err)
	}
	if err := manager.SetRendererMute(ctx, "uuid:fake", false); !errors.Is(err, ErrUnsupportedAction) {
		t.Fatalf("unsupported mute error = %v", err)
	}
}

func TestRendererRegistryAccessHandlesOfflineState(t *testing.T) {
	manager := NewManager(nil, "http://127.0.0.1:18080")
	if len(manager.Renderers()) != 0 {
		t.Fatal("manager without a control point returned renderers")
	}
	if _, err := manager.Renderer("missing"); !errors.Is(err, ErrRendererOffline) {
		t.Fatalf("offline renderer error = %v", err)
	}
	manager.controlPoint = controlpoint.NewService(nil, nil, nil)
	if _, err := manager.Renderer("missing"); !errors.Is(err, ErrRendererOffline) {
		t.Fatalf("missing renderer error = %v", err)
	}
}

func rendererManager(t *testing.T, handler http.HandlerFunc) (*Manager, *[]string) {
	t.Helper()
	var mu sync.Mutex
	requests := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		requests = append(requests, r.Header.Get("SOAPACTION")+"\n"+string(body))
		mu.Unlock()
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	manager := NewManager(nil, "http://127.0.0.1:18080")
	_, network, _ := net.ParseCIDR("127.0.0.0/8")
	manager.controlPoint = controlpoint.NewService(nil, []*net.IPNet{network}, server.Client())
	manager.controlPoint.Registry().Upsert(controlpoint.Renderer{
		UDN: "uuid:fake", FriendlyName: "Fake", AVTransportURL: server.URL, RenderingControlURL: server.URL,
		SourceIP: net.ParseIP("127.0.0.1"), LocalIP: net.ParseIP("127.0.0.1"),
		AVTransportActionsKnown: true, RenderingActionsKnown: true,
		AVTransportActions: []string{"GetTransportInfo", "GetPositionInfo", "Play", "Pause", "Stop", "Next", "Previous", "Seek"},
		RenderingActions:   []string{"GetVolume", "SetVolume", "SetMute"},
	}, net.ParseIP("127.0.0.1"), time.Minute)
	return manager, &requests
}

func containsRequest(requests []string, value string) bool {
	for _, request := range requests {
		if strings.Contains(request, value) {
			return true
		}
	}
	return false
}
