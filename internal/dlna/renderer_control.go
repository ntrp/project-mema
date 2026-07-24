package dlna

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/dlna/content"
	"media-manager/internal/dlna/controlpoint"
	"media-manager/internal/dlna/ssdp"
)

var (
	ErrRendererOffline        = errors.New("renderer is offline or its discovery lease expired")
	ErrUnsupportedAction      = errors.New("renderer does not support this action")
	ErrMediaNotPlayable       = errors.New("media file is not available for DLNA playback")
	ErrInvalidRendererCommand = errors.New("invalid renderer command")
)

type RendererPlayback struct {
	URI       string
	Metadata  string
	ProfileID string
}

func (m *Manager) startControlPoint(ctx context.Context, ifaces []ssdp.Interface) {
	items := make([]controlpoint.Interface, 0, len(ifaces))
	for _, iface := range ifaces {
		items = append(items, controlpoint.Interface{Name: iface.Name, Hardware: iface.Hardware, Address: iface.Addr})
	}
	m.controlPoint = controlpoint.NewService(items, m.allowedNets, nil)
	m.controlPoint.Start(ctx)
}

func (m *Manager) stopControlPoint() {
	if m.controlPoint != nil {
		m.controlPoint.Stop()
		m.controlPoint = nil
	}
}

func (m *Manager) Renderers() []controlpoint.Renderer {
	service := m.controlPointSnapshot()
	if service == nil {
		return []controlpoint.Renderer{}
	}
	return service.Registry().List(time.Time{})
}

func (m *Manager) Renderer(udn string) (controlpoint.Renderer, error) {
	service := m.controlPointSnapshot()
	if service == nil {
		return controlpoint.Renderer{}, ErrRendererOffline
	}
	renderer, ok := service.Registry().Get(udn, time.Time{})
	if !ok {
		return controlpoint.Renderer{}, ErrRendererOffline
	}
	return renderer, nil
}

func (m *Manager) RendererState(ctx context.Context, udn string) (controlpoint.RendererState, error) {
	renderer, client, err := m.rendererClient(udn)
	if err != nil {
		return controlpoint.RendererState{}, err
	}
	type result struct {
		kind string
		data any
		err  error
	}
	calls := make([]func() result, 0, 4)
	if controlpoint.SupportsTransport(renderer, "GetTransportInfo") {
		calls = append(calls, func() result {
			value, err := client.GetTransportInfo(ctx, renderer)
			return result{"transport", value, err}
		})
	}
	if controlpoint.SupportsTransport(renderer, "GetPositionInfo") {
		calls = append(calls, func() result {
			value, err := client.GetPositionInfo(ctx, renderer)
			return result{"position", value, err}
		})
	}
	if renderer.RenderingControlURL != "" && controlpoint.SupportsRendering(renderer, "GetVolume") {
		calls = append(calls, func() result {
			value, err := client.GetVolume(ctx, renderer)
			return result{"volume", value, err}
		})
	}
	if renderer.RenderingControlURL != "" && controlpoint.SupportsRendering(renderer, "GetMute") {
		calls = append(calls, func() result {
			value, err := client.GetMute(ctx, renderer)
			return result{"mute", value, err}
		})
	}
	results := make(chan result, len(calls))
	for _, call := range calls {
		go func() { results <- call() }()
	}
	var state controlpoint.RendererState
	var readErrors []error
	successes := 0
	for range calls {
		value := <-results
		if value.err != nil {
			var fault *controlpoint.FaultError
			if !errors.As(value.err, &fault) || fault.UPnPCode != "401" {
				readErrors = append(readErrors, value.err)
			}
			continue
		}
		successes++
		switch value.kind {
		case "transport":
			state.Transport, state.TransportAvailable = value.data.(controlpoint.TransportState), true
		case "position":
			state.Position, state.PositionAvailable = value.data.(controlpoint.PositionInfo), true
		case "volume":
			state.Volume.Volume, state.VolumeAvailable = value.data.(int), true
		case "mute":
			state.Volume.Mute, state.MuteAvailable = value.data.(bool), true
		}
	}
	if successes == 0 && len(readErrors) > 0 {
		return state, errors.Join(readErrors...)
	}
	return state, nil
}

func (m *Manager) PlayOnRenderer(ctx context.Context, udn string, mediaID uuid.UUID, path string) (RendererPlayback, error) {
	renderer, client, err := m.rendererClient(udn)
	if err != nil {
		return RendererPlayback{}, err
	}
	objectID := content.EncodeID(content.MediaItemRef(mediaID))
	if strings.TrimSpace(path) != "" {
		objectID = content.EncodeID(content.FileRef(mediaID, path))
	}
	object, err := m.contentTree().BrowseMetadata(ctx, objectID)
	if err != nil || object.FilePath == "" {
		return RendererPlayback{}, ErrMediaNotPlayable
	}
	baseURL := "http://" + net.JoinHostPort(renderer.LocalIP.String(), m.httpPort)
	profile := m.ExplainRendererProfile(ctx, RendererRequest{
		FriendlyName: renderer.FriendlyName, ClientIP: renderer.SourceIP.String(), RendererUUID: renderer.UDN,
	}).Profile
	m.SetRendererProfileOverride(renderer.SourceIP.String(), profile.ID)
	resources := contentResources([]content.Object{object}, baseURL, profile)
	objects := content.ApplySubtitleURLs(baseURL, content.ApplyArtworkURLs(baseURL, []content.Object{object}))
	items := resources[object.ID]
	if len(items) == 0 || items[0].URL == "" {
		return RendererPlayback{}, ErrMediaNotPlayable
	}
	selected := items[len(items)-1]
	metadata, err := content.RenderDIDLWithOptions(objects, resources, DIDLOptionsForProfile(profile))
	if err != nil {
		return RendererPlayback{}, err
	}
	if !controlpoint.SupportsTransport(renderer, "SetAVTransportURI") || !controlpoint.SupportsTransport(renderer, "Play") {
		return RendererPlayback{}, ErrUnsupportedAction
	}
	if err := client.SetAVTransportURI(ctx, renderer, selected.URL, string(metadata)); err != nil {
		return RendererPlayback{}, err
	}
	if err := client.Play(ctx, renderer); err != nil {
		return RendererPlayback{}, err
	}
	return RendererPlayback{URI: selected.URL, Metadata: string(metadata), ProfileID: profile.ID}, nil
}

func (m *Manager) Transport(ctx context.Context, udn, action string) error {
	renderer, client, err := m.rendererClient(udn)
	if err != nil {
		return err
	}
	if !controlpoint.SupportsTransport(renderer, action) {
		return ErrUnsupportedAction
	}
	switch strings.ToLower(action) {
	case "play":
		return client.Play(ctx, renderer)
	case "pause":
		return client.Pause(ctx, renderer)
	case "stop":
		return client.Stop(ctx, renderer)
	case "next":
		return client.Next(ctx, renderer)
	case "previous":
		return client.Previous(ctx, renderer)
	default:
		return ErrUnsupportedAction
	}
}

func (m *Manager) SeekRenderer(ctx context.Context, udn string, seconds int) error {
	if seconds < 0 {
		return fmt.Errorf("%w: seek position must not be negative", ErrInvalidRendererCommand)
	}
	renderer, client, err := m.rendererClient(udn)
	if err != nil {
		return err
	}
	if !controlpoint.SupportsTransport(renderer, "Seek") {
		return ErrUnsupportedAction
	}
	target := fmt.Sprintf("%d:%02d:%02d", seconds/3600, seconds/60%60, seconds%60)
	return client.Seek(ctx, renderer, target)
}

func (m *Manager) SetRendererVolume(ctx context.Context, udn string, volume int) error {
	if volume < 0 || volume > 100 {
		return fmt.Errorf("%w: volume must be between 0 and 100", ErrInvalidRendererCommand)
	}
	renderer, client, err := m.rendererClient(udn)
	if err != nil {
		return err
	}
	if renderer.RenderingControlURL == "" || !controlpoint.SupportsRendering(renderer, "SetVolume") {
		return ErrUnsupportedAction
	}
	return client.SetVolume(ctx, renderer, volume)
}

func (m *Manager) SetRendererMute(ctx context.Context, udn string, mute bool) error {
	renderer, client, err := m.rendererClient(udn)
	if err != nil {
		return err
	}
	if renderer.RenderingControlURL == "" || !controlpoint.SupportsRendering(renderer, "SetMute") {
		return ErrUnsupportedAction
	}
	return client.SetMute(ctx, renderer, mute)
}

func (m *Manager) rendererClient(udn string) (controlpoint.Renderer, controlpoint.Client, error) {
	service := m.controlPointSnapshot()
	if service == nil {
		return controlpoint.Renderer{}, controlpoint.Client{}, ErrRendererOffline
	}
	renderer, ok := service.Registry().Get(udn, time.Time{})
	if !ok {
		return controlpoint.Renderer{}, controlpoint.Client{}, ErrRendererOffline
	}
	return renderer, service.Client(), nil
}

func (m *Manager) controlPointSnapshot() *controlpoint.Service {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.controlPoint
}
