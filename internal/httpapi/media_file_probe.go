package httpapi

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

type ffprobeOutput struct {
	Streams  []ffprobeStream  `json:"streams"`
	Chapters []ffprobeChapter `json:"chapters"`
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
	Tags          map[string]string `json:"tags"`
}

type ffprobeChapter struct {
	ID        int32             `json:"id"`
	StartTime string            `json:"start_time"`
	EndTime   string            `json:"end_time"`
	Tags      map[string]string `json:"tags"`
}

type mediaFileProbeResult struct {
	tracks   []MediaFileTrack
	chapters []MediaFileChapter
}

func mediaFileProbe(path string) mediaFileProbeResult {
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return mediaFileProbeResult{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	output, err := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_streams",
		"-show_chapters",
		"-of", "json",
		path,
	).Output()
	if err != nil {
		return mediaFileProbeResult{}
	}
	var payload ffprobeOutput
	if err := json.Unmarshal(output, &payload); err != nil {
		return mediaFileProbeResult{}
	}
	return mediaFileProbeResult{
		tracks:   mediaFileTracks(payload.Streams),
		chapters: mediaFileChapters(payload.Chapters),
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
		BitRate:       optionalProbeString(stream.BitRate),
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
