package httpapi

import (
	"bytes"
	"errors"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

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
	return mediaPreviewArgsWithPlan(target, audioTrackIndex, mediaPreviewPlan(target, audioTrackIndex))
}

func mediaPreviewArgsWithPlan(target string, audioTrackIndex *int32, plan mediaPreviewTranscodePlan) []string {
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
	args = append(args,
		"-sn",
		"-c:v", plan.videoCodec,
	)
	if plan.videoCodec != "copy" {
		args = append(args,
			"-preset", "veryfast",
			"-pix_fmt", "yuv420p",
		)
	}
	args = append(args, "-c:a", plan.audioCodec)
	if plan.audioCodec != "copy" {
		args = append(args, "-ac", "2")
	}
	return append(args,
		"-movflags", "frag_keyframe+empty_moov+default_base_moof",
		"-f", "mp4",
		"pipe:1",
	)
}

type mediaPreviewTranscodePlan struct {
	videoCodec string
	audioCodec string
}

func mediaPreviewPlan(target string, audioTrackIndex *int32) mediaPreviewTranscodePlan {
	probe := mediaFileProbe(target)
	return mediaPreviewPlanFromTracks(probe.tracks, audioTrackIndex)
}

func mediaPreviewPlanFromTracks(tracks []MediaFileTrack, audioTrackIndex *int32) mediaPreviewTranscodePlan {
	plan := mediaPreviewTranscodePlan{videoCodec: "libx264", audioCodec: "aac"}
	if browserCompatiblePreviewVideo(firstTrackByType(tracks, Video, nil)) {
		plan.videoCodec = "copy"
	}
	if browserCompatiblePreviewAudio(firstTrackByType(tracks, Audio, audioTrackIndex)) {
		plan.audioCodec = "copy"
	}
	return plan
}

func firstTrackByType(tracks []MediaFileTrack, trackType MediaFileTrackType, index *int32) *MediaFileTrack {
	for i := range tracks {
		if tracks[i].Type != trackType {
			continue
		}
		if index != nil && (tracks[i].Index == nil || *tracks[i].Index != *index) {
			continue
		}
		return &tracks[i]
	}
	return nil
}

func browserCompatiblePreviewVideo(track *MediaFileTrack) bool {
	if track == nil || track.Codec == nil {
		return false
	}
	codec := strings.ToLower(strings.TrimSpace(*track.Codec))
	if codec != "h264" && codec != "avc1" {
		return false
	}
	pixelFormat := ""
	if track.PixelFormat != nil {
		pixelFormat = strings.ToLower(strings.TrimSpace(*track.PixelFormat))
	}
	return pixelFormat == "" || pixelFormat == "yuv420p" || pixelFormat == "yuvj420p"
}

func browserCompatiblePreviewAudio(track *MediaFileTrack) bool {
	if track == nil || track.Codec == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(*track.Codec), "aac")
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
