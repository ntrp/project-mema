package httpapi

import (
	"net/http"

	"github.com/google/uuid"
	"media-manager/internal/delivery"
	mediatools "media-manager/internal/tools"
)

func (s *Server) PreviewMediaItemFile(w http.ResponseWriter, r *http.Request, id ResourceId, params PreviewMediaItemFileParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	target, err := s.settings.MediaItemFilePath(r.Context(), uuid.UUID(id), params.Path)
	if err != nil {
		writeSettingsError(w, err, "Could not find media file")
		return
	}
	if _, ok := statMediaFile(w, target); !ok {
		return
	}
	probe := delivery.Probe(target)
	clientProfile := mediaPreviewClientProfile(params.ClientProfile, r.UserAgent())
	decision := delivery.DecisionFromTracks(target, probe.Tracks, params.AudioTrackIndex, deliveryClientProfile(clientProfile))
	if decision.DeliveryProtocol == delivery.ProtocolFile {
		serveMediaFile(w, r, target)
		return
	}
	if !requireMediaPreviewFFmpeg(w) {
		return
	}
	serveMediaPreviewHLSPlaylist(w, r, params.Path, target, probe, params.AudioTrackIndex, clientProfile, decision)
}

func (s *Server) PreviewMediaItemFileSegment(w http.ResponseWriter, r *http.Request, id ResourceId, params PreviewMediaItemFileSegmentParams) {
	if _, ok := s.requireSession(w, r); !ok {
		return
	}
	if !delivery.ValidSegment(params.SegmentStartSeconds, params.SegmentDurationSeconds) {
		writeError(w, http.StatusBadRequest, "invalid_input", "Preview segment range is invalid")
		return
	}
	target, err := s.settings.MediaItemFilePath(r.Context(), uuid.UUID(id), params.Path)
	if err != nil {
		writeSettingsError(w, err, "Could not find media file")
		return
	}
	if _, ok := statMediaFile(w, target); !ok {
		return
	}
	if !requireMediaPreviewFFmpeg(w) {
		return
	}
	if err := mediatools.SafePathArg(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_media_path", "Media file path is invalid")
		return
	}
	probe := delivery.Probe(target)
	clientProfile := mediaPreviewClientProfile(params.ClientProfile, r.UserAgent())
	decision := delivery.DecisionFromTracks(target, probe.Tracks, params.AudioTrackIndex, deliveryClientProfile(clientProfile))
	args := delivery.SegmentArgs(target, params.AudioTrackIndex, params.SegmentStartSeconds, params.SegmentDurationSeconds, decision)
	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	wrote, err := runMediaPreviewCommand(r, w, args)
	if err != nil && !wrote && r.Context().Err() == nil {
		writeError(w, http.StatusInternalServerError, "media_preview_failed", "Could not start media preview segment")
	}
}

func requireMediaPreviewFFmpeg(w http.ResponseWriter) bool {
	if _, err := mediatools.LookPath("ffmpeg"); err == nil {
		return true
	}
	writeError(w, http.StatusInternalServerError, "ffmpeg_not_available", "ffmpeg is required for browser preview")
	return false
}

func runMediaPreviewCommand(r *http.Request, w http.ResponseWriter, args []string) (bool, error) {
	writer := &flushWriter{w: w}
	err := mediatools.RunStream(r.Context(), "ffmpeg", args, writer, 64*1024)
	return writer.wrote, err
}

type flushWriter struct {
	w     http.ResponseWriter
	wrote bool
}

func (w *flushWriter) Write(payload []byte) (int, error) {
	n, err := w.w.Write(payload)
	if flusher, ok := w.w.(http.Flusher); ok {
		flusher.Flush()
	}
	if n > 0 {
		w.wrote = true
	}
	return n, err
}
