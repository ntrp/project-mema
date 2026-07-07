package delivery

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

func Probe(path string) ProbeResult {
	if _, err := mediatools.LookPath("ffprobe"); err != nil {
		return ProbeResult{}
	}
	if err := mediatools.SafePathArg(path); err != nil {
		return ProbeResult{}
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
		return ProbeResult{}
	}
	var payload ffprobeOutput
	if err := json.Unmarshal(output, &payload); err != nil {
		return ProbeResult{}
	}
	return ProbeResult{
		Container:       containerInfo(payload.Format),
		Tracks:          tracks(payload.Streams),
		Chapters:        chapters(payload.Chapters),
		DurationSeconds: OptionalDuration(payload.Format.Duration),
	}
}

func containerInfo(format ffprobeFormat) Container {
	return Container{
		BitRate:    optionalString(format.BitRate),
		Format:     optionalString(format.Format),
		FormatName: optionalString(format.FormatName),
	}
}

func tracks(streams []ffprobeStream) []Track {
	tracks := []Track{}
	for _, stream := range streams {
		track, ok := trackFromStream(stream)
		if ok {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

func chapters(chapters []ffprobeChapter) []Chapter {
	results := make([]Chapter, 0, len(chapters))
	for index, chapter := range chapters {
		number := chapter.ID
		if number <= 0 {
			number = int32(index)
		}
		results = append(results, Chapter{
			Index:     number,
			Title:     optionalString(chapter.Tags["title"]),
			StartTime: optionalString(chapter.StartTime),
			EndTime:   optionalString(chapter.EndTime),
		})
	}
	return results
}

func trackFromStream(stream ffprobeStream) (Track, bool) {
	trackType, ok := trackTypeFromCodec(stream.CodecType)
	if !ok {
		return Track{}, false
	}
	return Track{
		Type:          trackType,
		Index:         optionalIndex(stream.Index),
		Codec:         optionalString(stream.CodecName),
		Language:      optionalString(languageTag(stream.Tags)),
		Title:         optionalString(stream.Tags["title"]),
		BitRate:       streamBitRate(stream),
		ChannelLayout: optionalString(stream.ChannelLayout),
		FrameRate:     optionalString(normalFrameRate(stream.FrameRate)),
		Height:        optionalInt(stream.Height),
		Width:         optionalInt(stream.Width),
		PixelFormat:   optionalString(stream.PixelFormat),
		Profile:       optionalString(stream.Profile),
		Channels:      optionalInt(stream.Channels),
	}, true
}

func streamBitRate(stream ffprobeStream) *string {
	if bitRate := optionalString(stream.BitRate); bitRate != nil {
		return bitRate
	}
	if bitRate := optionalString(probeTag(stream.Tags, "BPS")); bitRate != nil {
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
	if duration := OptionalDuration(stream.Duration); duration != nil {
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

func trackTypeFromCodec(value string) (TrackType, bool) {
	switch strings.ToLower(value) {
	case "video":
		return TrackVideo, true
	case "audio":
		return TrackAudio, true
	case "subtitle":
		return TrackSubtitle, true
	default:
		return "", false
	}
}
