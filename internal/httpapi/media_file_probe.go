package httpapi

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	mediatools "media-manager/internal/tools"
)

type ffprobeOutput struct {
	Streams  []ffprobeStream  `json:"streams"`
	Chapters []ffprobeChapter `json:"chapters"`
	Format   ffprobeFormat    `json:"format"`
}

type ffprobeFormat struct {
	BitRate    string `json:"bit_rate"`
	Duration   string `json:"duration"`
	Format     string `json:"format_name"`
	FormatName string `json:"format_long_name"`
}

type ffprobeStream struct {
	Index         int32             `json:"index"`
	CodecName     string            `json:"codec_name"`
	CodecType     string            `json:"codec_type"`
	Profile       string            `json:"profile"`
	Width         int32             `json:"width"`
	Height        int32             `json:"height"`
	PixelFormat   string            `json:"pix_fmt"`
	FrameRate     string            `json:"avg_frame_rate"`
	Channels      int32             `json:"channels"`
	ChannelLayout string            `json:"channel_layout"`
	BitRate       string            `json:"bit_rate"`
	Duration      string            `json:"duration"`
	Tags          map[string]string `json:"tags"`
}

type ffprobeChapter struct {
	ID        int32             `json:"id"`
	StartTime string            `json:"start_time"`
	EndTime   string            `json:"end_time"`
	Tags      map[string]string `json:"tags"`
}

type mediaFileProbeResult struct {
	container       mediaFileContainer
	tracks          []MediaFileTrack
	chapters        []MediaFileChapter
	durationSeconds *float64
}

type mediaFileContainer struct {
	bitRate    *string
	format     *string
	formatName *string
}

func mediaFileProbe(path string) mediaFileProbeResult {
	if _, err := mediatools.LookPath("ffprobe"); err != nil {
		return mediaFileProbeResult{}
	}
	if err := mediatools.SafePathArg(path); err != nil {
		return mediaFileProbeResult{}
	}
	output, err := mediatools.RunOutput(context.Background(), mediatools.CommandSpec{
		Name: "ffprobe",
		Args: []string{
			"-v", "error",
			"-show_streams",
			"-show_chapters",
			"-show_format",
			"-of", "json",
			path,
		},
		Timeout:        3 * time.Second,
		MaxOutputBytes: 4 * 1024 * 1024,
		MaxStderrBytes: 64 * 1024,
	})
	if err != nil {
		return mediaFileProbeResult{}
	}
	var payload ffprobeOutput
	if err := json.Unmarshal(output, &payload); err != nil {
		return mediaFileProbeResult{}
	}
	return mediaFileProbeResult{
		container:       mediaFileContainerInfo(payload.Format),
		tracks:          mediaFileTracks(payload.Streams),
		chapters:        mediaFileChapters(payload.Chapters),
		durationSeconds: optionalProbeDuration(payload.Format.Duration),
	}
}

func mediaFileContainerInfo(format ffprobeFormat) mediaFileContainer {
	return mediaFileContainer{
		bitRate:    optionalProbeString(format.BitRate),
		format:     optionalProbeString(format.Format),
		formatName: optionalProbeString(format.FormatName),
	}
}

func mediaFileTracks(streams []ffprobeStream) []MediaFileTrack {
	tracks := []MediaFileTrack{}
	for _, stream := range streams {
		track, ok := mediaFileTrack(stream)
		if ok {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

func mediaFileChapters(chapters []ffprobeChapter) []MediaFileChapter {
	results := make([]MediaFileChapter, 0, len(chapters))
	for index, chapter := range chapters {
		number := chapter.ID
		if number <= 0 {
			number = int32(index)
		}
		results = append(results, MediaFileChapter{
			Index:     number,
			Title:     optionalProbeString(chapter.Tags["title"]),
			StartTime: optionalProbeString(chapter.StartTime),
			EndTime:   optionalProbeString(chapter.EndTime),
		})
	}
	return results
}

func mediaFileTrack(stream ffprobeStream) (MediaFileTrack, bool) {
	trackType, ok := mediaFileTrackType(stream.CodecType)
	if !ok {
		return MediaFileTrack{}, false
	}
	track := MediaFileTrack{
		Type:          trackType,
		Index:         optionalProbeIndex(stream.Index),
		Codec:         optionalProbeString(stream.CodecName),
		Language:      optionalProbeString(languageTag(stream.Tags)),
		Title:         optionalProbeString(stream.Tags["title"]),
		BitRate:       streamBitRate(stream),
		ChannelLayout: optionalProbeString(stream.ChannelLayout),
		FrameRate:     optionalProbeString(normalFrameRate(stream.FrameRate)),
		Height:        optionalProbeInt(stream.Height),
		Width:         optionalProbeInt(stream.Width),
		PixelFormat:   optionalProbeString(stream.PixelFormat),
		Profile:       optionalProbeString(stream.Profile),
		Channels:      optionalProbeInt(stream.Channels),
	}
	return track, true
}

func streamBitRate(stream ffprobeStream) *string {
	if bitRate := optionalProbeString(stream.BitRate); bitRate != nil {
		return bitRate
	}
	if bitRate := optionalProbeString(probeTag(stream.Tags, "BPS")); bitRate != nil {
		return bitRate
	}
	bytes, err := strconv.ParseFloat(probeTag(stream.Tags, "NUMBER_OF_BYTES"), 64)
	duration := streamDurationSeconds(stream)
	if err != nil || bytes <= 0 || duration <= 0 {
		return nil
	}
	value := strconv.FormatInt(int64(bytes*8/duration), 10)
	return &value
}

func streamDurationSeconds(stream ffprobeStream) float64 {
	if duration := optionalProbeDuration(stream.Duration); duration != nil {
		return *duration
	}
	return probeDurationTagSeconds(probeTag(stream.Tags, "DURATION"))
}

func probeTag(tags map[string]string, name string) string {
	name = strings.ToLower(strings.ReplaceAll(name, "_", "-"))
	for key, value := range tags {
		key = strings.ToLower(strings.ReplaceAll(key, "_", "-"))
		if key == name || strings.HasPrefix(key, name+"-") {
			return value
		}
	}
	return ""
}

func probeDurationTagSeconds(value string) float64 {
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 3 {
		return 0
	}
	hours, errHours := strconv.ParseFloat(parts[0], 64)
	minutes, errMinutes := strconv.ParseFloat(parts[1], 64)
	seconds, errSeconds := strconv.ParseFloat(parts[2], 64)
	if errHours != nil || errMinutes != nil || errSeconds != nil {
		return 0
	}
	return hours*3600 + minutes*60 + seconds
}

func mediaFileTrackType(value string) (MediaFileTrackType, bool) {
	switch strings.ToLower(value) {
	case "video":
		return Video, true
	case "audio":
		return Audio, true
	case "subtitle":
		return Subtitle, true
	default:
		return "", false
	}
}
