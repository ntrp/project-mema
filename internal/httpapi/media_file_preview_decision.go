package httpapi

import (
	"path/filepath"
	"strconv"
	"strings"

	"media-manager/internal/playback"
)

const (
	mediaPreviewDeliveryFile MediaFilePreviewDeliveryProtocol = "file"
	mediaPreviewDeliveryHLS  MediaFilePreviewDeliveryProtocol = "hls"

	mediaPreviewReasonWebKitHLS = "webkit_native_hls_requires_video_transcode"
)

type mediaPreviewDecision struct {
	deliveryProtocol MediaFilePreviewDeliveryProtocol
	mode             MediaFilePreviewMode
	plan             mediaPreviewTranscodePlan
	reasons          []string
}

type mediaPreviewTranscodePlan struct {
	videoCodec string
	audioCodec string
}

func mediaPreviewDecisionFromTracks(
	target string,
	tracks []MediaFileTrack,
	audioTrackIndex *int32,
	clientProfile MediaFilePreviewClientProfile,
) mediaPreviewDecision {
	source := playbackMediaSource(target, tracks)
	index := playbackAudioIndex(audioTrackIndex)
	options := mediaPreviewPlaybackOptions(clientProfile)
	options.AudioStreamIndex = index
	decision := mediaPreviewDecisionFromStream(playback.BuildVideoStreamInfo(source, options))
	if clientProfile == Webkit && decision.mode == Transcode && decision.plan.videoCodec != "copy" {
		decision.reasons = appendPreviewReason(decision.reasons, mediaPreviewReasonWebKitHLS)
	}
	return decision
}

func mediaPreviewPlaybackOptions(clientProfile MediaFilePreviewClientProfile) playback.MediaOptions {
	options := playback.MediaOptions{
		Profile:              playback.BrowserVideoProfile(),
		EnableDirectPlay:     true,
		EnableDirectStream:   true,
		EnableTranscoding:    true,
		AllowVideoStreamCopy: true,
		AllowAudioStreamCopy: true,
	}
	if clientProfile == Webkit {
		options.EnableDirectStream = false
		options.AllowVideoStreamCopy = false
	}
	return options
}

func mediaPreviewDecisionFromStream(stream playback.StreamInfo) mediaPreviewDecision {
	mode := Transcode
	delivery := mediaPreviewDeliveryHLS
	if stream.PlayMethod == playback.PlayMethodDirectPlay {
		mode = Direct
		delivery = mediaPreviewDeliveryFile
	} else if stream.PlayMethod == playback.PlayMethodDirectStream {
		mode = Remux
	}
	return mediaPreviewDecision{
		deliveryProtocol: delivery,
		mode:             mode,
		plan: mediaPreviewTranscodePlan{
			videoCodec: ffmpegVideoCodec(stream.OutputVideoCodec),
			audioCodec: ffmpegAudioCodec(stream.OutputAudioCodec),
		},
		reasons: mediaPreviewReasonStrings(stream.TranscodeReasons),
	}
}

func playbackMediaSource(target string, tracks []MediaFileTrack) playback.MediaSource {
	return playback.MediaSource{
		Container:            mediaPreviewContainer(target),
		Video:                playbackVideoStream(firstTrackByType(tracks, Video, nil)),
		AudioStreams:         playbackAudioStreams(tracks),
		SupportsDirectPlay:   true,
		SupportsDirectStream: true,
		SupportsTranscoding:  true,
	}
}

func playbackVideoStream(track *MediaFileTrack) *playback.MediaStream {
	if track == nil {
		return nil
	}
	stream := playbackTrack(track, playback.StreamVideo)
	return &stream
}

func playbackAudioStreams(tracks []MediaFileTrack) []playback.MediaStream {
	streams := []playback.MediaStream{}
	for i := range tracks {
		if tracks[i].Type == Audio {
			streams = append(streams, playbackTrack(&tracks[i], playback.StreamAudio))
		}
	}
	return streams
}

func playbackTrack(track *MediaFileTrack, trackType playback.StreamType) playback.MediaStream {
	return playback.MediaStream{
		Index:       optionalTrackIndex(track.Index),
		Type:        trackType,
		Codec:       optionalTrackString(track.Codec),
		PixelFormat: optionalTrackString(track.PixelFormat),
		BitRate:     optionalTrackInt64(track.BitRate),
		Channels:    optionalTrackInt(track.Channels),
	}
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

func playbackAudioIndex(index *int32) *int {
	if index == nil {
		return nil
	}
	value := int(*index)
	return &value
}

func mediaPreviewContainer(target string) string {
	return strings.TrimPrefix(strings.ToLower(filepath.Ext(target)), ".")
}

func ffmpegVideoCodec(codec string) string {
	switch strings.ToLower(strings.TrimSpace(codec)) {
	case "", "h264", "avc1":
		return "libx264"
	case "copy":
		return "copy"
	default:
		return codec
	}
}

func ffmpegAudioCodec(codec string) string {
	if strings.TrimSpace(codec) == "" {
		return "aac"
	}
	return codec
}

func mediaPreviewReasonStrings(reasons []playback.TranscodeReason) []string {
	values := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		values = append(values, string(reason))
	}
	return values
}

func mediaPreviewClientProfile(explicit *MediaFilePreviewClientProfile, userAgent string) MediaFilePreviewClientProfile {
	if explicit != nil {
		return *explicit
	}
	if webkitPreviewUserAgent(userAgent) {
		return Webkit
	}
	return Browser
}

func webkitPreviewUserAgent(userAgent string) bool {
	agent := strings.ToLower(userAgent)
	if strings.Contains(agent, "iphone") || strings.Contains(agent, "ipad") || strings.Contains(agent, "ipod") {
		return true
	}
	if !strings.Contains(agent, "applewebkit") || !strings.Contains(agent, "safari") {
		return false
	}
	return !containsPreviewUserAgentToken(agent, []string{"chrome", "chromium", "crios", "fxios", "edg", "opr", "android"})
}

func containsPreviewUserAgentToken(agent string, tokens []string) bool {
	for _, token := range tokens {
		if strings.Contains(agent, token) {
			return true
		}
	}
	return false
}

func appendPreviewReason(reasons []string, reason string) []string {
	for _, existing := range reasons {
		if existing == reason {
			return reasons
		}
	}
	return append(reasons, reason)
}

func optionalTrackIndex(value *int32) int {
	if value == nil {
		return -1
	}
	return int(*value)
}

func optionalTrackString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func optionalTrackInt(value *int32) int {
	if value == nil {
		return 0
	}
	return int(*value)
}

func optionalTrackInt64(value *string) int64 {
	if value == nil {
		return 0
	}
	parsed, err := strconv.ParseInt(*value, 10, 64)
	if err != nil || parsed < 0 {
		return 0
	}
	return parsed
}
