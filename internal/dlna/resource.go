package dlna

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"media-manager/internal/delivery"
	mediatools "media-manager/internal/tools"
)

var dlnaTranscodeSlots = make(chan struct{}, 2)

func (m *Manager) resource(w http.ResponseWriter, r *http.Request) {
	id, segment := resourceIDFromPath(r.URL.Path)
	object, err := m.contentTree().BrowseMetadata(r.Context(), id)
	if err != nil || object.FilePath == "" {
		http.NotFound(w, r)
		return
	}
	if segment {
		m.resourceSegment(w, r, object.FilePath)
		return
	}
	profile := m.RendererProfileFromRequest(r)
	if rejectUnsupportedSeek(w, r, profile) {
		return
	}
	if r.URL.Query().Get("mode") == "hls" {
		m.resourcePlaylist(w, r, id, object.FilePath)
		return
	}
	probe := probeWithPathContainer(delivery.Probe(object.FilePath), object.FilePath)
	capability := EvaluateRendererCapability(profile, probe)
	mode := r.URL.Query().Get("mode")
	if mode == "remux" || capability.Decision.Mode == delivery.ModeRemux {
		m.resourceRemux(w, r, id, object.FilePath, remuxDecision(), profile)
		return
	}
	decision := capability.Decision
	if mode == "transcode" && decision.Mode != delivery.ModeTranscode {
		decision = matroskaAudioTranscodeDecision(probe)
	}
	if mode == "transcode" ||
		(capability.Decision.Mode == delivery.ModeTranscode &&
			capability.Decision.DeliveryProtocol == delivery.ProtocolFile) {
		m.resourceTranscode(w, r, id, object.FilePath, decision, profile)
		return
	}
	if capability.Decision.DeliveryProtocol == delivery.ProtocolHLS {
		m.resourcePlaylist(w, r, id, object.FilePath)
		return
	}
	done, ok := m.beginStream(r, id, "direct", false)
	if !ok {
		http.Error(w, "too many DLNA streams", http.StatusTooManyRequests)
		return
	}
	defer done()
	writeFileError(w, delivery.ServeFile(w, r, object.FilePath))
}

func (m *Manager) resourceTranscode(
	w http.ResponseWriter,
	r *http.Request,
	id string,
	target string,
	decision delivery.Decision,
	profile RendererProfile,
) {
	output := transcodeOutputTarget(profile)
	done, ok := m.beginStream(r, id, output.StreamMode, true)
	if !ok {
		http.Error(w, "too many DLNA streams", http.StatusTooManyRequests)
		return
	}
	defer done()
	cachePath, cached := m.existingDLNAOutput(target, output)
	if isSeekRange(r.Header.Get("Range")) && !cached {
		var err error
		cachePath, err = m.cachedDLNAOutput(r, target, decision, output)
		if err != nil {
			http.Error(w, "could not prepare DLNA remux", http.StatusInternalServerError)
			return
		}
		cached = true
	}
	if cached {
		setDLNAOutputHeaders(w, output)
		writeFileError(w, delivery.ServeFile(w, r, cachePath))
		return
	}
	setDLNAOutputHeaders(w, output)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	if r.Method == http.MethodHead {
		return
	}
	if !acquireTranscodeSlot(w, r) {
		return
	}
	defer func() { <-dlnaTranscodeSlots }()
	if err := mediatools.SafePathArg(target); err != nil {
		http.Error(w, "invalid media path", http.StatusBadRequest)
		return
	}
	ctx, cancel := m.commandContext(r)
	defer cancel()
	writer := flushWriter{w: w}
	err := mediatools.RunStream(ctx, "ffmpeg", dlnaOutputArgs(target, "pipe:1", decision, output), &writer, 64*1024)
	if err != nil {
		if !writer.wrote && ctx.Err() == nil {
			http.Error(w, "could not start DLNA remux", http.StatusInternalServerError)
		}
	}
}

func (m *Manager) resourcePlaylist(w http.ResponseWriter, r *http.Request, id string, target string) {
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("TransferMode.DLNA.ORG", "Streaming")
	w.Header().Set("ContentFeatures.DLNA.ORG", "DLNA.ORG_OP=01;DLNA.ORG_CI=1")
	if r.Method == http.MethodHead {
		return
	}
	done, ok := m.beginStream(r, id, "hls_playlist", true)
	if !ok {
		http.Error(w, "too many DLNA streams", http.StatusTooManyRequests)
		return
	}
	defer done()
	probe := delivery.Probe(target)
	duration := 0.0
	if probe.DurationSeconds != nil {
		duration = *probe.DurationSeconds
	}
	decision := delivery.DecisionFromTracks(target, probe.Tracks, nil, DeliveryClientProfile(m.RendererProfileFromRequest(r)))
	request := delivery.PlaylistRequest{
		Path:        r.URL.Path,
		FilePath:    target,
		Segments:    delivery.HLSSegmentsForDecision(target, duration, decision),
		SegmentPath: "/dlna/resource/" + url.PathEscape(id) + "/segment",
	}
	_, _ = w.Write([]byte(delivery.HLSPlaylistText(request)))
}

func (m *Manager) resourceSegment(w http.ResponseWriter, r *http.Request, target string) {
	start, duration, ok := segmentRange(r)
	if !ok {
		http.Error(w, "invalid segment range", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("TransferMode.DLNA.ORG", "Streaming")
	if r.Method == http.MethodHead {
		return
	}
	done, ok := m.beginStream(r, r.URL.Path, "hls_segment", true)
	if !ok {
		http.Error(w, "too many DLNA streams", http.StatusTooManyRequests)
		return
	}
	defer done()
	if !acquireTranscodeSlot(w, r) {
		return
	}
	defer func() { <-dlnaTranscodeSlots }()
	if err := mediatools.SafePathArg(target); err != nil {
		http.Error(w, "invalid media path", http.StatusBadRequest)
		return
	}
	probe := delivery.Probe(target)
	decision := delivery.DecisionFromTracks(target, probe.Tracks, nil, DeliveryClientProfile(m.RendererProfileFromRequest(r)))
	args := delivery.SegmentArgs(target, nil, start, duration, decision)
	writer := flushWriter{w: w}
	ctx, cancel := m.commandContext(r)
	defer cancel()
	err := mediatools.RunStream(ctx, "ffmpeg", args, &writer, 64*1024)
	if err != nil && !writer.wrote && ctx.Err() == nil {
		http.Error(w, "could not start DLNA segment", http.StatusInternalServerError)
	}
}

func resourceIDFromPath(path string) (string, bool) {
	path = strings.TrimPrefix(path, "/dlna")
	path = strings.TrimPrefix(path, "/resource/")
	segment := strings.HasSuffix(path, "/segment")
	path = strings.TrimSuffix(path, "/segment")
	id, err := url.PathUnescape(strings.Trim(path, "/"))
	if err != nil {
		return "", segment
	}
	return id, segment
}

func isSeekRange(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	if !strings.HasPrefix(value, "bytes=") {
		return false
	}
	spec := strings.TrimSpace(strings.TrimPrefix(value, "bytes="))
	if index := strings.Index(spec, ","); index >= 0 {
		spec = strings.TrimSpace(spec[:index])
	}
	start, _, ok := strings.Cut(spec, "-")
	if !ok || strings.TrimSpace(start) == "" {
		return false
	}
	offset, err := strconv.ParseInt(strings.TrimSpace(start), 10, 64)
	return err == nil && offset > 0
}

func segmentRange(r *http.Request) (float64, float64, bool) {
	start, err := strconv.ParseFloat(r.URL.Query().Get("segmentStartSeconds"), 64)
	if err != nil {
		return 0, 0, false
	}
	duration, err := strconv.ParseFloat(r.URL.Query().Get("segmentDurationSeconds"), 64)
	if err != nil || !delivery.ValidSegment(start, duration) {
		return 0, 0, false
	}
	return start, duration, true
}

func acquireTranscodeSlot(w http.ResponseWriter, r *http.Request) bool {
	select {
	case dlnaTranscodeSlots <- struct{}{}:
		return true
	case <-r.Context().Done():
		return false
	default:
		http.Error(w, "too many DLNA transcodes", http.StatusTooManyRequests)
		return false
	}
}

func writeFileError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, delivery.ErrDirectory) {
		http.Error(w, "media file path points to a directory", http.StatusBadRequest)
		return
	}
	if os.IsNotExist(err) {
		http.Error(w, "could not find media file", http.StatusNotFound)
		return
	}
	http.Error(w, "could not open media file", http.StatusInternalServerError)
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
