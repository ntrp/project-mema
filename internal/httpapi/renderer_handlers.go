package httpapi

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"media-manager/internal/dlna"
	"media-manager/internal/dlna/controlpoint"
)

func (s *Server) ListRenderers(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	s.listRenderers(w)
}

func (s *Server) listRenderers(w http.ResponseWriter) {
	result := RendererListResponse{Renderers: []Renderer{}}
	if s.dlna != nil {
		for _, discovered := range s.dlna.Renderers() {
			result.Renderers = append(result.Renderers, rendererResponse(discovered))
		}
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) GetRendererStatus(w http.ResponseWriter, r *http.Request, id RendererId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	s.getRendererStatus(w, r, id)
}

func (s *Server) getRendererStatus(w http.ResponseWriter, r *http.Request, id RendererId) {
	if s.dlna == nil {
		writeRendererError(w, dlna.ErrRendererOffline)
		return
	}
	discovered, err := s.dlna.Renderer(id)
	if err != nil {
		writeRendererError(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()
	state, err := s.dlna.RendererState(ctx, id)
	if err != nil {
		writeRendererError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rendererStatusResponse(discovered, state))
}

func (s *Server) PlayMediaOnRenderer(w http.ResponseWriter, r *http.Request, id RendererId) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	var body RendererPlayRequest
	if !decodeJSON(w, r, &body) {
		return
	}
	s.playMediaOnRenderer(w, r, id, body)
}

func (s *Server) playMediaOnRenderer(w http.ResponseWriter, r *http.Request, id RendererId, body RendererPlayRequest) {
	path := ""
	if body.Path != nil {
		path = *body.Path
	}
	result, err := withRendererControl(s, r, id, "play-media", func(ctx context.Context) (dlna.RendererPlayback, error) {
		if s.dlna == nil {
			return dlna.RendererPlayback{}, dlna.ErrRendererOffline
		}
		return s.dlna.PlayOnRenderer(ctx, id, uuid.UUID(body.MediaItemId), path)
	})
	if err != nil {
		writeRendererError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, RendererPlaybackResponse{ProfileId: result.ProfileID})
}

func (s *Server) ControlRendererTransport(w http.ResponseWriter, r *http.Request, id RendererId) {
	var body RendererTransportRequest
	if !rendererCommandBody(s, w, r, &body) {
		return
	}
	err := rendererControl(s, r, id, string(body.Action), func(ctx context.Context) error {
		return s.dlna.Transport(ctx, id, string(body.Action))
	})
	writeRendererCommand(w, err)
}
func (s *Server) SeekRenderer(w http.ResponseWriter, r *http.Request, id RendererId) {
	var body RendererSeekRequest
	if !rendererCommandBody(s, w, r, &body) {
		return
	}
	err := rendererControl(s, r, id, "seek", func(ctx context.Context) error {
		return s.dlna.SeekRenderer(ctx, id, int(body.PositionSeconds))
	})
	writeRendererCommand(w, err)
}
func (s *Server) SetRendererVolume(w http.ResponseWriter, r *http.Request, id RendererId) {
	var body RendererVolumeRequest
	if !rendererCommandBody(s, w, r, &body) {
		return
	}
	err := rendererControl(s, r, id, "volume", func(ctx context.Context) error {
		return s.dlna.SetRendererVolume(ctx, id, int(body.Volume))
	})
	writeRendererCommand(w, err)
}
func (s *Server) SetRendererMute(w http.ResponseWriter, r *http.Request, id RendererId) {
	var body RendererMuteRequest
	if !rendererCommandBody(s, w, r, &body) {
		return
	}
	err := rendererControl(s, r, id, "mute", func(ctx context.Context) error {
		return s.dlna.SetRendererMute(ctx, id, body.Muted)
	})
	writeRendererCommand(w, err)
}

func rendererCommandBody(s *Server, w http.ResponseWriter, r *http.Request, body any) bool {
	if _, ok := s.requireSession(w, r); !ok {
		return false
	}
	if s.dlna == nil {
		writeRendererError(w, dlna.ErrRendererOffline)
		return false
	}
	return decodeJSON(w, r, body)
}
func rendererControl(s *Server, r *http.Request, id, action string, call func(context.Context) error) error {
	_, err := withRendererControl(s, r, id, action, func(ctx context.Context) (struct{}, error) {
		return struct{}{}, call(ctx)
	})
	return err
}
func withRendererControl[T any](s *Server, r *http.Request, id, action string, call func(context.Context) (T, error)) (T, error) {
	s.events.Publish("renderer.control.attempted", map[string]any{"rendererId": id, "action": action})
	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()
	result, err := call(ctx)
	status := "succeeded"
	if err != nil {
		status = "failed"
		slog.Warn("renderer control failed", "renderer_id", id, "action", action, "error", err)
	}
	s.events.Publish("renderer.control.result", map[string]any{"rendererId": id, "action": action, "status": status})
	return result, err
}
func writeRendererCommand(w http.ResponseWriter, err error) {
	if err != nil {
		writeRendererError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func writeRendererError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, dlna.ErrRendererOffline):
		writeError(w, http.StatusNotFound, "renderer_offline", "Renderer is offline; wait for it to reappear and try again")
	case errors.Is(err, dlna.ErrUnsupportedAction):
		writeError(w, http.StatusConflict, "renderer_action_unsupported", "The renderer does not support this action")
	case errors.Is(err, dlna.ErrMediaNotPlayable):
		writeError(w, http.StatusBadRequest, "media_not_playable", "Choose an available media file for playback")
	case errors.Is(err, dlna.ErrInvalidRendererCommand):
		writeError(w, http.StatusBadRequest, "invalid_renderer_command", err.Error())
	default:
		var fault *controlpoint.FaultError
		if errors.As(err, &fault) {
			writeError(w, http.StatusBadGateway, "renderer_fault", fault.Error())
			return
		}
		writeError(w, http.StatusBadGateway, "renderer_unreachable", "The renderer did not respond; check that it is online")
	}
}

func rendererStatusResponse(renderer controlpoint.Renderer, state controlpoint.RendererState) RendererStatus {
	response := RendererStatus{
		TransportState: state.Transport.State, TransportStatus: state.Transport.Status,
		PositionSeconds: upnpSeconds(state.Position.RelTime), DurationSeconds: upnpSeconds(state.Position.TrackDuration),
	}
	if renderer.RenderingControlURL != "" && state.VolumeAvailable {
		volume := int32(state.Volume.Volume)
		response.Volume = &volume
	}
	if renderer.RenderingControlURL != "" && state.MuteAvailable {
		muted := state.Volume.Mute
		response.Muted = &muted
	}
	return response
}

func rendererResponse(value controlpoint.Renderer) Renderer {
	return Renderer{
		Id: value.UDN, Name: value.FriendlyName, Manufacturer: value.Manufacturer, ModelName: value.ModelName,
		Address: value.SourceIP.String(), LastSeen: value.LastSeen, ExpiresAt: value.ExpiresAt,
		Supports: RendererCapabilities{
			Play: controlpoint.SupportsTransport(value, "Play"), Pause: controlpoint.SupportsTransport(value, "Pause"),
			Stop: controlpoint.SupportsTransport(value, "Stop"), Seek: controlpoint.SupportsTransport(value, "Seek"),
			Next: controlpoint.SupportsTransport(value, "Next"), Previous: controlpoint.SupportsTransport(value, "Previous"),
			Volume: value.RenderingControlURL != "" && controlpoint.SupportsRendering(value, "SetVolume"),
			Mute:   value.RenderingControlURL != "" && controlpoint.SupportsRendering(value, "SetMute"),
		},
	}
}
func upnpSeconds(value string) int32 {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0
	}
	hours, e1 := strconv.Atoi(parts[0])
	minutes, e2 := strconv.Atoi(parts[1])
	seconds, e3 := strconv.ParseFloat(parts[2], 64)
	if e1 != nil || e2 != nil || e3 != nil || hours < 0 || minutes < 0 || seconds < 0 {
		return 0
	}
	return int32(hours*3600 + minutes*60 + int(seconds))
}
