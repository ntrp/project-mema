package dlna

import (
	"net/http"

	"media-manager/internal/delivery"
	mediatools "media-manager/internal/tools"
)

func (m *Manager) resourceRemux(
	w http.ResponseWriter,
	r *http.Request,
	id string,
	target string,
	decision delivery.Decision,
	profile RendererProfile,
) {
	output := remuxOutputTarget(profile)
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
	writer := flushWriter{w: w}
	err := mediatools.RunStream(r.Context(), "ffmpeg", dlnaOutputArgs(target, "pipe:1", decision, output), &writer, 64*1024)
	if err != nil && !writer.wrote && r.Context().Err() == nil {
		http.Error(w, "could not start DLNA remux", http.StatusInternalServerError)
	}
}

func remuxDecision() delivery.Decision {
	return delivery.Decision{
		DeliveryProtocol: delivery.ProtocolFile,
		Mode:             delivery.ModeRemux,
		Plan: delivery.TranscodePlan{
			VideoCodec: "copy",
			AudioCodec: "copy",
		},
	}
}
