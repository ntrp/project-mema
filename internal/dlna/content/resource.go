package content

import (
	"fmt"
	"mime"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/delivery"
)

type Resource struct {
	URL             string
	ProtocolInfo    string
	SizeBytes       *int64
	Duration        *string
	BitRate         *int64
	Resolution      *string
	AudioChannels   *int32
	SampleFrequency *int64
}

type ResourceInput struct {
	URL       string
	SizeBytes *int64
	Probe     delivery.ProbeResult
	Decision  delivery.Decision
}

func ResourceFromDelivery(input ResourceInput) Resource {
	video := delivery.FirstTrackByType(input.Probe.Tracks, delivery.TrackVideo, nil)
	audio := delivery.FirstTrackByType(input.Probe.Tracks, delivery.TrackAudio, nil)
	return Resource{
		URL:             input.URL,
		ProtocolInfo:    ProtocolInfo(input.URL, input.Probe.Container, input.Decision),
		SizeBytes:       resourceSize(input.SizeBytes, input.Decision),
		Duration:        resourceDuration(input.Probe.DurationSeconds),
		BitRate:         parseInt64Ptr(delivery.SelectedBitRate(video, audio)),
		Resolution:      resourceResolution(video),
		AudioChannels:   resourceChannels(audio),
		SampleFrequency: nil,
	}
}

func ResourceURL(baseURL string, object Object) string {
	base, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return "/dlna/resource/" + url.PathEscape(object.ID)
	}
	base.Path = path.Join(base.Path, "/dlna/resource", object.ID)
	return base.String()
}

func ProtocolInfo(resourceURL string, container delivery.Container, decision delivery.Decision) string {
	mimeType := resourceMIME(resourceURL, container, decision)
	conversion := "0"
	if decision.Mode == delivery.ModeTranscode {
		conversion = "1"
	}
	return fmt.Sprintf("http-get:*:%s:DLNA.ORG_OP=01;DLNA.ORG_CI=%s", mimeType, conversion)
}

func resourceMIME(resourceURL string, container delivery.Container, decision delivery.Decision) string {
	if decision.DeliveryProtocol == delivery.ProtocolHLS {
		return "application/vnd.apple.mpegurl"
	}
	if container.FormatName != nil {
		if value := containerFormatMIME(*container.FormatName); value != "" {
			return value
		}
	}
	extension := strings.ToLower(filepath.Ext(resourceURL))
	if value := mime.TypeByExtension(extension); value != "" {
		return strings.Split(value, ";")[0]
	}
	if extension == ".mkv" {
		return "video/x-matroska"
	}
	return "application/octet-stream"
}

func containerFormatMIME(format string) string {
	format = strings.ToLower(format)
	switch {
	case strings.Contains(format, "matroska"):
		return "video/x-matroska"
	case strings.Contains(format, "mp4") || strings.Contains(format, "quicktime"):
		return "video/mp4"
	case strings.Contains(format, "mpegts"):
		return "video/mp2t"
	default:
		return ""
	}
}

func resourceSize(size *int64, decision delivery.Decision) *int64 {
	if decision.DeliveryProtocol == delivery.ProtocolHLS {
		return nil
	}
	return size
}

func resourceDuration(value *float64) *string {
	if value == nil || *value <= 0 {
		return nil
	}
	duration := time.Duration(*value * float64(time.Second))
	hours := int(duration / time.Hour)
	duration -= time.Duration(hours) * time.Hour
	minutes := int(duration / time.Minute)
	duration -= time.Duration(minutes) * time.Minute
	seconds := int(duration / time.Second)
	millis := int((duration - time.Duration(seconds)*time.Second) / time.Millisecond)
	formatted := fmt.Sprintf("%d:%02d:%02d.%03d", hours, minutes, seconds, millis)
	return &formatted
}

func resourceResolution(track *delivery.Track) *string {
	if track == nil || track.Width == nil || track.Height == nil {
		return nil
	}
	value := fmt.Sprintf("%dx%d", *track.Width, *track.Height)
	return &value
}

func resourceChannels(track *delivery.Track) *int32 {
	if track == nil || track.Channels == nil {
		return nil
	}
	return track.Channels
}

func parseInt64Ptr(value *string) *int64 {
	if value == nil {
		return nil
	}
	parsed, err := strconv.ParseInt(*value, 10, 64)
	if err != nil || parsed <= 0 {
		return nil
	}
	return &parsed
}
