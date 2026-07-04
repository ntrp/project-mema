package httpapi

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const mediaPreviewHLSSegmentSeconds = 6.0

type mediaPreviewHLSSegment struct {
	start    float64
	duration float64
}

func serveMediaPreviewHLSPlaylist(
	w http.ResponseWriter,
	r *http.Request,
	filePath MediaFilePath,
	target string,
	probe mediaFileProbeResult,
	audioTrackIndex *int32,
	clientProfile MediaFilePreviewClientProfile,
	decision mediaPreviewDecision,
) {
	if probe.durationSeconds == nil || !validPreviewSeconds(*probe.durationSeconds) || *probe.durationSeconds <= 0 {
		writeError(w, http.StatusInternalServerError, "media_preview_duration_unavailable", "Could not determine media duration for browser preview")
		return
	}
	segments := mediaPreviewHLSSegmentsForDecision(target, *probe.durationSeconds, decision)
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write([]byte(mediaPreviewHLSPlaylistText(r, filePath, audioTrackIndex, clientProfile, segments)))
}

func mediaPreviewHLSPlaylistText(
	r *http.Request,
	filePath MediaFilePath,
	audioTrackIndex *int32,
	clientProfile MediaFilePreviewClientProfile,
	segments []mediaPreviewHLSSegment,
) string {
	targetDuration := mediaPreviewHLSTargetDuration(segments)
	var builder strings.Builder
	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:3\n")
	builder.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")
	builder.WriteString("#EXT-X-TARGETDURATION:")
	builder.WriteString(strconv.Itoa(targetDuration))
	builder.WriteString("\n#EXT-X-MEDIA-SEQUENCE:0\n")
	for _, segment := range segments {
		builder.WriteString("#EXTINF:")
		builder.WriteString(formatPreviewSeconds(segment.duration))
		builder.WriteString(",\n")
		builder.WriteString(mediaPreviewHLSSegmentURL(r, filePath, audioTrackIndex, clientProfile, segment.start, segment.duration))
		builder.WriteByte('\n')
	}
	builder.WriteString("#EXT-X-ENDLIST\n")
	return builder.String()
}

func mediaPreviewHLSSegmentURL(
	r *http.Request,
	filePath MediaFilePath,
	audioTrackIndex *int32,
	clientProfile MediaFilePreviewClientProfile,
	start float64,
	duration float64,
) string {
	query := url.Values{
		"path":                   []string{string(filePath)},
		"segmentStartSeconds":    []string{formatPreviewSeconds(start)},
		"segmentDurationSeconds": []string{formatPreviewSeconds(duration)},
	}
	if audioTrackIndex != nil {
		query.Set("audioTrackIndex", strconv.FormatInt(int64(*audioTrackIndex), 10))
	}
	if clientProfile != "" && clientProfile != Browser {
		query.Set("clientProfile", string(clientProfile))
	}
	return (&url.URL{
		Path:     strings.TrimSuffix(r.URL.Path, "/preview") + "/preview-segment",
		RawQuery: query.Encode(),
	}).String()
}

func mediaPreviewHLSSegmentsForDecision(target string, duration float64, decision mediaPreviewDecision) []mediaPreviewHLSSegment {
	if decision.plan.videoCodec == "copy" {
		return mediaPreviewHLSSegments(duration, mediaPreviewVideoKeyframes(target))
	}
	return mediaPreviewFixedHLSSegments(duration)
}

func mediaPreviewHLSSegments(duration float64, keyframes []float64) []mediaPreviewHLSSegment {
	if len(keyframes) < 2 {
		return mediaPreviewFixedHLSSegments(duration)
	}
	boundaries := []float64{0}
	current := 0.0
	for current+mediaPreviewHLSSegmentSeconds < duration {
		next, ok := nextPreviewKeyframe(keyframes, current+mediaPreviewHLSSegmentSeconds, duration)
		if !ok {
			break
		}
		boundaries = append(boundaries, next)
		current = next
	}
	boundaries = append(boundaries, duration)
	return mediaPreviewSegmentsFromBoundaries(boundaries)
}

func mediaPreviewFixedHLSSegments(duration float64) []mediaPreviewHLSSegment {
	segments := []mediaPreviewHLSSegment{}
	for start := 0.0; start < duration; start += mediaPreviewHLSSegmentSeconds {
		segments = append(segments, mediaPreviewHLSSegment{
			start:    start,
			duration: math.Min(mediaPreviewHLSSegmentSeconds, duration-start),
		})
	}
	return segments
}

func nextPreviewKeyframe(keyframes []float64, minimum float64, duration float64) (float64, bool) {
	for _, value := range keyframes {
		if value >= duration || value < minimum {
			continue
		}
		if value > 0 {
			return value, true
		}
	}
	return 0, false
}

func mediaPreviewSegmentsFromBoundaries(boundaries []float64) []mediaPreviewHLSSegment {
	segments := []mediaPreviewHLSSegment{}
	for index := 0; index < len(boundaries)-1; index++ {
		start := boundaries[index]
		duration := boundaries[index+1] - start
		if duration > 0.05 {
			segments = append(segments, mediaPreviewHLSSegment{start: start, duration: duration})
		}
	}
	return segments
}

func mediaPreviewHLSTargetDuration(segments []mediaPreviewHLSSegment) int {
	maxDuration := mediaPreviewHLSSegmentSeconds
	for _, segment := range segments {
		maxDuration = math.Max(maxDuration, segment.duration)
	}
	return int(math.Ceil(maxDuration))
}

func mediaPreviewHLSSegmentArgs(target string, audioTrackIndex *int32, start float64, duration float64, decision mediaPreviewDecision) []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-ss", formatPreviewSeconds(start),
		"-t", formatPreviewSeconds(duration),
		"-fflags", "+genpts",
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
		"-dn",
		"-c:v", decision.plan.videoCodec,
	)
	if decision.plan.videoCodec == "copy" {
		args = append(args, "-bsf:v", "h264_mp4toannexb")
	} else {
		args = append(args,
			"-preset", "veryfast",
			"-pix_fmt", "yuv420p",
			"-profile:v", "high",
			"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%.3f)", mediaPreviewHLSSegmentSeconds),
		)
	}
	args = append(args, "-c:a", decision.plan.audioCodec)
	if decision.plan.audioCodec != "copy" {
		args = append(args, "-ac", "2")
	}
	return append(args,
		"-avoid_negative_ts", "make_zero",
		"-muxdelay", "0",
		"-muxpreload", "0",
		"-f", "mpegts",
		"pipe:1",
	)
}

func formatPreviewSeconds(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
}
