package playback

import "strings"

type ProfileType string

const (
	ProfileVideo ProfileType = "video"
	ProfileAudio ProfileType = "audio"
)

type PlayMethod string

const (
	PlayMethodDirectPlay   PlayMethod = "direct_play"
	PlayMethodDirectStream PlayMethod = "direct_stream"
	PlayMethodTranscode    PlayMethod = "transcode"
)

type StreamProtocol string

const (
	ProtocolHTTP StreamProtocol = "http"
	ProtocolHLS  StreamProtocol = "hls"
)

type StreamType string

const (
	StreamVideo    StreamType = "video"
	StreamAudio    StreamType = "audio"
	StreamSubtitle StreamType = "subtitle"
)

type CodecType string

const (
	CodecVideo CodecType = "video"
	CodecAudio CodecType = "audio"
)

type ProfileConditionProperty string

const (
	ConditionPixelFormat   ProfileConditionProperty = "pixel_format"
	ConditionAudioChannels ProfileConditionProperty = "audio_channels"
)

type TranscodeReason string

const (
	ReasonContainerNotSupported       TranscodeReason = "container_not_supported"
	ReasonVideoCodecNotSupported      TranscodeReason = "video_codec_not_supported"
	ReasonAudioCodecNotSupported      TranscodeReason = "audio_codec_not_supported"
	ReasonVideoPixelFormatUnsupported TranscodeReason = "video_pixel_format_not_supported"
	ReasonSecondaryAudioNotSupported  TranscodeReason = "secondary_audio_not_supported"
	ReasonContainerBitrateExceeded    TranscodeReason = "container_bitrate_exceeded"
)

type DeviceProfile struct {
	Name                string
	MaxStreamingBitrate int64
	DirectPlayProfiles  []DirectPlayProfile
	TranscodingProfiles []TranscodingProfile
	CodecProfiles       []CodecProfile
}

type DirectPlayProfile struct {
	Type        ProfileType
	Containers  []string
	VideoCodecs []string
	AudioCodecs []string
}

type TranscodingProfile struct {
	Type          ProfileType
	Container     string
	Protocol      StreamProtocol
	VideoCodecs   []string
	AudioCodecs   []string
	SegmentLength int
}

type CodecProfile struct {
	Type       CodecType
	Codecs     []string
	Conditions []ProfileCondition
}

type ProfileCondition struct {
	Property ProfileConditionProperty
	Allowed  []string
	Maximum  int
}

type MediaSource struct {
	Container            string
	BitRate              int64
	Video                *MediaStream
	AudioStreams         []MediaStream
	SupportsDirectPlay   bool
	SupportsDirectStream bool
	SupportsTranscoding  bool
}

type MediaStream struct {
	Index       int
	Type        StreamType
	Codec       string
	PixelFormat string
	BitRate     int64
	Channels    int
}

type MediaOptions struct {
	Profile              DeviceProfile
	AudioStreamIndex     *int
	MaxStreamingBitrate  int64
	EnableDirectPlay     bool
	EnableDirectStream   bool
	EnableTranscoding    bool
	AllowVideoStreamCopy bool
	AllowAudioStreamCopy bool
}

type StreamInfo struct {
	PlayMethod       PlayMethod
	Container        string
	Protocol         StreamProtocol
	OutputVideoCodec string
	OutputAudioCodec string
	AudioStreamIndex *int
	SegmentLength    int
	TranscodeReasons []TranscodeReason
}

func (p DirectPlayProfile) SupportsContainer(container string) bool {
	return containsToken(p.Containers, container)
}

func (p DirectPlayProfile) SupportsVideoCodec(codec string) bool {
	return len(p.VideoCodecs) == 0 || containsToken(p.VideoCodecs, codec)
}

func (p DirectPlayProfile) SupportsAudioCodec(codec string) bool {
	return codec == "" || len(p.AudioCodecs) == 0 || containsToken(p.AudioCodecs, codec)
}

func (p TranscodingProfile) SupportsVideoCodec(codec string) bool {
	return codec != "" && containsToken(p.VideoCodecs, codec)
}

func (p TranscodingProfile) SupportsAudioCodec(codec string) bool {
	return codec == "" || containsToken(p.AudioCodecs, codec)
}

func containsToken(tokens []string, value string) bool {
	value = normalizeToken(value)
	for _, token := range tokens {
		if normalizeToken(token) == value {
			return true
		}
	}
	return false
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
