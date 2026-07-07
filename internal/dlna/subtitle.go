package dlna

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"media-manager/internal/dlna/content"
	mediatools "media-manager/internal/tools"
)

func (m *Manager) subtitle(w http.ResponseWriter, r *http.Request) {
	id, index, ok := subtitleIDFromPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	object, err := m.contentTree().BrowseMetadata(r.Context(), id)
	if err != nil || index < 0 || index >= len(object.Subtitles) {
		http.NotFound(w, r)
		return
	}
	subtitle := object.Subtitles[index]
	w.Header().Set("Content-Type", subtitleContentType(subtitle))
	w.Header().Set("Cache-Control", "public, max-age=3600")
	if r.Method == http.MethodHead {
		return
	}
	if subtitle.Plan == content.SubtitleConvert {
		m.convertSubtitle(w, r, subtitle.FilePath)
		return
	}
	http.ServeFile(w, r, subtitle.FilePath)
}

func (m *Manager) convertSubtitle(w http.ResponseWriter, r *http.Request, target string) {
	if !acquireTranscodeSlot(w, r) {
		return
	}
	defer func() { <-dlnaTranscodeSlots }()
	if err := mediatools.SafePathArg(target); err != nil {
		http.Error(w, "invalid subtitle path", http.StatusBadRequest)
		return
	}
	args := []string{"-hide_banner", "-loglevel", "error", "-i", target, "-f", "webvtt", "pipe:1"}
	writer := flushWriter{w: w}
	err := mediatools.RunStream(r.Context(), "ffmpeg", args, &writer, 64*1024)
	if err != nil && !writer.wrote && r.Context().Err() == nil {
		http.Error(w, "could not convert subtitle", http.StatusInternalServerError)
	}
}

func subtitleIDFromPath(path string) (string, int, bool) {
	path = strings.TrimPrefix(path, "/dlna")
	path = strings.TrimPrefix(path, "/subtitle/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 {
		return "", 0, false
	}
	id, err := url.PathUnescape(parts[0])
	if err != nil {
		return "", 0, false
	}
	index, err := strconv.Atoi(parts[1])
	return id, index, err == nil
}

func subtitleContentType(subtitle content.Subtitle) string {
	if subtitle.Plan == content.SubtitleConvert || subtitle.Format == "vtt" {
		return "text/vtt; charset=utf-8"
	}
	return "application/x-subrip; charset=utf-8"
}
