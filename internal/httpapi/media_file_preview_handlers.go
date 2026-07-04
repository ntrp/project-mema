package httpapi

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/google/uuid"
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
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		writeError(w, http.StatusInternalServerError, "ffmpeg_not_available", "ffmpeg is required for browser preview")
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Accel-Buffering", "no")
	wrote, err := runMediaPreview(r, w, target, params.AudioTrackIndex)
	if err != nil && !wrote && r.Context().Err() == nil {
		writeError(w, http.StatusInternalServerError, "media_preview_failed", "Could not start media preview")
	}
}

func runMediaPreview(r *http.Request, w http.ResponseWriter, target string, audioTrackIndex *int32) (bool, error) {
	args := mediaPreviewArgs(target, audioTrackIndex)
	cmd := exec.CommandContext(r.Context(), "ffmpeg", args...)
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

func mediaPreviewArgs(target string, audioTrackIndex *int32) []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-i", target,
		"-map", "0:v:0",
	}
	if audioTrackIndex != nil {
		args = append(args, "-map", "0:"+strconv.FormatInt(int64(*audioTrackIndex), 10))
	} else {
		args = append(args, "-map", "0:a:0?")
	}
	return append(args,
		"-sn",
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-ac", "2",
		"-movflags", "frag_keyframe+empty_moov+default_base_moof",
		"-f", "mp4",
		"pipe:1",
	)
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
