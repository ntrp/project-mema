package delivery

import (
	"path/filepath"
	"strconv"
	"strings"

	"media-manager/internal/playback"
)

const ReasonWebKitHLS = "webkit_native_hls_requires_video_transcode"

func DecisionFromTracks(target string, tracks []Track, audioTrackIndex *int32, clientProfile ClientProfile) Decision {
	source := playbackMediaSource(target, tracks)
	index := playbackAudioIndex(audioTrackIndex)
	options := playbackOptions(clientProfile)
	options.AudioStreamIndex = index
	decision := decisionFromStream(playback.BuildVideoStreamInfo(source, options))
	if clientProfile == ClientWebKit && decision.Mode == ModeTranscode && decision.Plan.VideoCodec != "copy" {
		decision.Reasons = appendReason(decision.Reasons, ReasonWebKitHLS)
	}
	return decision
}

func ClientProfileForRequest(explicit *ClientProfile, userAgent string) ClientProfile {
	if explicit != nil {
		return *explicit
	}
	if webkitUserAgent(userAgent) {
		return ClientWebKit
	}
	return ClientBrowser
}

func FirstTrackByType(tracks []Track, trackType TrackType, index *int32) *Track {
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

func SelectedBitRate(tracks ...*Track) *string {
	var total int64
	for _, track := range tracks {
		if track == nil || track.BitRate == nil {
			continue
		}
		value, err := strconv.ParseInt(*track.BitRate, 10, 64)
		if err != nil || value <= 0 {
			continue
		}
		total += value
	}
	if total <= 0 {
		return nil
	}
	formatted := strconv.FormatInt(total, 10)
	return &formatted
}

func playbackOptions(clientProfile ClientProfile) playback.MediaOptions {
	options := playback.MediaOptions{
		Profile:              playback.BrowserVideoProfile(),
		EnableDirectPlay:     true,
		EnableDirectStream:   true,
		EnableTranscoding:    true,
		AllowVideoStreamCopy: true,
		AllowAudioStreamCopy: true,
	}
	if clientProfile == ClientWebKit {
		options.EnableDirectStream = false
		options.AllowVideoStreamCopy = false
	}
	return options
}

func decisionFromStream(stream playback.StreamInfo) Decision {
	mode := ModeTranscode
	protocol := ProtocolHLS
	if stream.PlayMethod == playback.PlayMethodDirectPlay {
		mode = ModeDirect
		protocol = ProtocolFile
	} else if stream.PlayMethod == playback.PlayMethodDirectStream {
		mode = ModeRemux
	}
	return Decision{
		DeliveryProtocol: protocol,
		Mode:             mode,
		Plan: TranscodePlan{
			VideoCodec: ffmpegVideoCodec(stream.OutputVideoCodec),
			AudioCodec: ffmpegAudioCodec(stream.OutputAudioCodec),
		},
		Reasons: reasonStrings(stream.TranscodeReasons),
	}
}

func playbackMediaSource(target string, tracks []Track) playback.MediaSource {
	return playback.MediaSource{
		Container:            container(target),
		Video:                playbackVideoStream(FirstTrackByType(tracks, TrackVideo, nil)),
		AudioStreams:         playbackAudioStreams(tracks),
		SupportsDirectPlay:   true,
		SupportsDirectStream: true,
		SupportsTranscoding:  true,
	}
}

func playbackVideoStream(track *Track) *playback.MediaStream {
	if track == nil {
		return nil
	}
	stream := playbackTrack(track, playback.StreamVideo)
	return &stream
}

func playbackAudioStreams(tracks []Track) []playback.MediaStream {
	streams := []playback.MediaStream{}
	for i := range tracks {
		if tracks[i].Type == TrackAudio {
			streams = append(streams, playbackTrack(&tracks[i], playback.StreamAudio))
		}
	}
	return streams
}

func playbackTrack(track *Track, trackType playback.StreamType) playback.MediaStream {
	return playback.MediaStream{
		Index:       optionalTrackIndex(track.Index),
		Type:        trackType,
		Codec:       optionalTrackString(track.Codec),
		PixelFormat: optionalTrackString(track.PixelFormat),
		BitRate:     optionalTrackInt64(track.BitRate),
		Channels:    optionalTrackInt(track.Channels),
	}
}

func playbackAudioIndex(index *int32) *int {
	if index == nil {
		return nil
	}
	value := int(*index)
	return &value
}

func container(target string) string {
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

func reasonStrings(reasons []playback.TranscodeReason) []string {
	values := make([]string, 0, len(reasons))
	for _, reason := range reasons {
		values = append(values, string(reason))
	}
	return values
}

func webkitUserAgent(userAgent string) bool {
	agent := strings.ToLower(userAgent)
	if strings.Contains(agent, "iphone") || strings.Contains(agent, "ipad") || strings.Contains(agent, "ipod") {
		return true
	}
	if !strings.Contains(agent, "applewebkit") || !strings.Contains(agent, "safari") {
		return false
	}
	return !containsUserAgentToken(agent, []string{"chrome", "chromium", "crios", "fxios", "edg", "opr", "android"})
}

func containsUserAgentToken(agent string, tokens []string) bool {
	for _, token := range tokens {
		if strings.Contains(agent, token) {
			return true
		}
	}
	return false
}

func appendReason(reasons []string, reason string) []string {
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
