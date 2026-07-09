package delivery

type TrackType string

const (
	TrackVideo    TrackType = "video"
	TrackAudio    TrackType = "audio"
	TrackSubtitle TrackType = "subtitle"
)

type ClientProfile string

const (
	ClientBrowser ClientProfile = "browser"
	ClientWebKit  ClientProfile = "webkit"
)

type Mode string

const (
	ModeDirect    Mode = "direct"
	ModeRemux     Mode = "remux"
	ModeTranscode Mode = "transcode"
)

type Protocol string

const (
	ProtocolFile Protocol = "file"
	ProtocolHLS  Protocol = "hls"
)

type Track struct {
	Index         *int32
	Type          TrackType
	Codec         *string
	Language      *string
	Title         *string
	Duration      *float64
	BitRate       *string
	ChannelLayout *string
	FrameRate     *string
	Height        *int32
	Width         *int32
	PixelFormat   *string
	Profile       *string
	Channels      *int32
}

type Chapter struct {
	Index     int32
	Title     *string
	StartTime *string
	EndTime   *string
}

type Container struct {
	BitRate    *string
	Format     *string
	FormatName *string
}

type ProbeResult struct {
	Container       Container
	Tracks          []Track
	Chapters        []Chapter
	DurationSeconds *float64
}

type Decision struct {
	DeliveryProtocol Protocol
	Mode             Mode
	Plan             TranscodePlan
	Reasons          []string
}

type TranscodePlan struct {
	VideoCodec string
	AudioCodec string
}
