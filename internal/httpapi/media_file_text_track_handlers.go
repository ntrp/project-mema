package httpapi

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/google/uuid"
)

func (s *Server) PreviewMediaItemFileSubtitle(w http.ResponseWriter, r *http.Request, id ResourceId, params PreviewMediaItemFileSubtitleParams) {
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
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		writeError(w, http.StatusInternalServerError, "ffmpeg_not_available", "ffmpeg is required for subtitle preview")
		return
	}
	w.Header().Set("Content-Type", "text/vtt; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	wrote, err := runMediaSubtitlePreview(r, w, target, params.SubtitleTrackIndex)
	if err != nil && !wrote && r.Context().Err() == nil {
		writeError(w, http.StatusInternalServerError, "media_subtitle_failed", "Could not convert subtitle track")
	}
}

func runMediaSubtitlePreview(r *http.Request, w http.ResponseWriter, target string, subtitleTrackIndex int32) (bool, error) {
	cmd := exec.CommandContext(r.Context(), "ffmpeg", mediaSubtitlePreviewArgs(target, subtitleTrackIndex)...)
	var stderr bytes.Buffer
	writer := &flushWriter{w: w}
	cmd.Stdout = writer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil && stderr.Len() > 0 {
		return writer.wrote, errors.New(stderr.String())
	}
	return writer.wrote, err
}

func mediaSubtitlePreviewArgs(target string, subtitleTrackIndex int32) []string {
	return []string{
		"-hide_banner",
		"-loglevel", "error",
		"-i", target,
		"-map", "0:" + strconv.FormatInt(int64(subtitleTrackIndex), 10),
		"-f", "webvtt",
		"pipe:1",
	}
}
