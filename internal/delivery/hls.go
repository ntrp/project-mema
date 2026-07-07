package delivery

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

const HLSSegmentSeconds = 6.0

type HLSSegment struct {
	Start    float64
	Duration float64
}

type PlaylistRequest struct {
	Path          string
	FilePath      string
	AudioTrack    *int32
	ClientProfile ClientProfile
	Segments      []HLSSegment
	SegmentPath   string
}

func HLSPlaylistText(request PlaylistRequest) string {
	targetDuration := HLSTargetDuration(request.Segments)
	var builder strings.Builder
	builder.WriteString("#EXTM3U\n")
	builder.WriteString("#EXT-X-VERSION:3\n")
	builder.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")
	builder.WriteString("#EXT-X-TARGETDURATION:")
	builder.WriteString(strconv.Itoa(targetDuration))
	builder.WriteString("\n#EXT-X-MEDIA-SEQUENCE:0\n")
	for _, segment := range request.Segments {
		builder.WriteString("#EXTINF:")
		builder.WriteString(FormatSeconds(segment.Duration))
		builder.WriteString(",\n")
		builder.WriteString(HLSSegmentURL(request, segment.Start, segment.Duration))
		builder.WriteByte('\n')
	}
	builder.WriteString("#EXT-X-ENDLIST\n")
	return builder.String()
}

func HLSSegmentURL(request PlaylistRequest, start float64, duration float64) string {
	query := url.Values{
		"path":                   []string{request.FilePath},
		"segmentStartSeconds":    []string{FormatSeconds(start)},
		"segmentDurationSeconds": []string{FormatSeconds(duration)},
	}
	if request.AudioTrack != nil {
		query.Set("audioTrackIndex", strconv.FormatInt(int64(*request.AudioTrack), 10))
	}
	if request.ClientProfile != "" && request.ClientProfile != ClientBrowser {
		query.Set("clientProfile", string(request.ClientProfile))
	}
	segmentPath := request.SegmentPath
	if segmentPath == "" {
		segmentPath = strings.TrimSuffix(request.Path, "/preview") + "/preview-segment"
	}
	return (&url.URL{Path: segmentPath, RawQuery: query.Encode()}).String()
}

func HLSSegmentsForDecision(target string, duration float64, decision Decision) []HLSSegment {
	if decision.Plan.VideoCodec == "copy" {
		return HLSSegments(duration, VideoKeyframes(target))
	}
	return FixedHLSSegments(duration)
}

func HLSSegments(duration float64, keyframes []float64) []HLSSegment {
	if len(keyframes) < 2 {
		return FixedHLSSegments(duration)
	}
	boundaries := []float64{0}
	current := 0.0
	for current+HLSSegmentSeconds < duration {
		next, ok := nextKeyframe(keyframes, current+HLSSegmentSeconds, duration)
		if !ok {
			break
		}
		boundaries = append(boundaries, next)
		current = next
	}
	boundaries = append(boundaries, duration)
	return segmentsFromBoundaries(boundaries)
}

func FixedHLSSegments(duration float64) []HLSSegment {
	segments := []HLSSegment{}
	for start := 0.0; start < duration; start += HLSSegmentSeconds {
		segments = append(segments, HLSSegment{
			Start:    start,
			Duration: math.Min(HLSSegmentSeconds, duration-start),
		})
	}
	return segments
}

func SegmentArgs(target string, audioTrackIndex *int32, start float64, duration float64, decision Decision) []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-ss", FormatSeconds(start),
		"-t", FormatSeconds(duration),
		"-fflags", "+genpts",
		"-i", target,
		"-map", "0:v:0",
	}
	if audioTrackIndex != nil {
		args = append(args, "-map", "0:"+strconv.FormatInt(int64(*audioTrackIndex), 10))
	} else {
		args = append(args, "-map", "0:a:0?")
	}
	args = append(args, "-sn", "-dn", "-c:v", decision.Plan.VideoCodec)
	if decision.Plan.VideoCodec == "copy" {
		args = append(args, "-bsf:v", "h264_mp4toannexb")
	} else {
		args = append(args,
			"-preset", "veryfast",
			"-pix_fmt", "yuv420p",
			"-profile:v", "high",
			"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%.3f)", HLSSegmentSeconds),
		)
	}
	args = append(args, "-c:a", decision.Plan.AudioCodec)
	if decision.Plan.AudioCodec != "copy" {
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

func ValidSegment(start, duration float64) bool {
	return ValidSeconds(start) && ValidSeconds(duration) && duration > 0 && duration <= 60
}

func ValidSeconds(value float64) bool {
	return value >= 0 && !math.IsInf(value, 0) && !math.IsNaN(value)
}

func HLSTargetDuration(segments []HLSSegment) int {
	maxDuration := HLSSegmentSeconds
	for _, segment := range segments {
		maxDuration = math.Max(maxDuration, segment.Duration)
	}
	return int(math.Ceil(maxDuration))
}

func FormatSeconds(value float64) string {
	return strconv.FormatFloat(value, 'f', 3, 64)
}

func nextKeyframe(keyframes []float64, minimum float64, duration float64) (float64, bool) {
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

func segmentsFromBoundaries(boundaries []float64) []HLSSegment {
	segments := []HLSSegment{}
	for index := 0; index < len(boundaries)-1; index++ {
		start := boundaries[index]
		duration := boundaries[index+1] - start
		if duration > 0.05 {
			segments = append(segments, HLSSegment{Start: start, Duration: duration})
		}
	}
	return segments
}
