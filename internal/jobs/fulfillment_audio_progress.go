package jobs

import (
	"context"
	"strconv"
	"strings"
	"time"

	"media-manager/internal/delivery"
	"media-manager/internal/events"
	"media-manager/internal/storage"
	mediatools "media-manager/internal/tools"
)

func runAudioTranscodeCommand(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	track storage.MediaFileTrackFact,
	argsList []string,
) error {
	progress := audioTranscodeProgress{durationMs: audioTranscodeDurationMs(item, track)}
	start := int32(0)
	recordJobProgressData(ctx, settings, eventBroker, &start, "Transcoding audio 0%", audioTranscodeProgressData(item, track, progress.durationMs))
	_, err := mediatools.RunOutputProgress(ctx, mediatools.ProgressCommandSpec{
		CommandSpec: mediatools.CommandSpec{
			Name:           "ffmpeg",
			Args:           ffmpegProgressArgs(argsList),
			Timeout:        2 * time.Hour,
			MaxOutputBytes: 0,
			MaxStderrBytes: 128 * 1024,
		},
		Progress: func(line string) {
			if percent, ok := progress.percent(line); ok {
				recordJobProgressData(ctx, settings, eventBroker, &percent, "Transcoding audio "+strconv.Itoa(int(percent))+"%", audioTranscodeProgressData(item, track, progress.durationMs))
			}
		},
	})
	if err != nil {
		return err
	}
	done := int32(100)
	recordJobProgressData(ctx, settings, eventBroker, &done, "Audio transcode complete", audioTranscodeProgressData(item, track, progress.durationMs))
	return nil
}

type audioTranscodeProgress struct {
	durationMs int64
	last       int32
	lastAt     time.Time
}

func (p *audioTranscodeProgress) percent(line string) (int32, bool) {
	if p.durationMs <= 0 {
		return 0, false
	}
	value, ok := ffmpegOutTimeMicroseconds(line)
	if !ok {
		return 0, false
	}
	percent := int32((value * 100) / (p.durationMs * 1000))
	if percent < 0 {
		percent = 0
	}
	if percent > 99 {
		percent = 99
	}
	now := time.Now()
	if percent <= p.last && now.Sub(p.lastAt) < 2*time.Second {
		return 0, false
	}
	p.last = percent
	p.lastAt = now
	return percent, true
}

func ffmpegOutTimeMicroseconds(line string) (int64, bool) {
	key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
	if !ok {
		return 0, false
	}
	switch key {
	case "out_time_us", "out_time_ms":
		parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		return parsed, err == nil
	case "out_time":
		return ffmpegClockMicroseconds(strings.TrimSpace(value))
	default:
		return 0, false
	}
}

func ffmpegClockMicroseconds(value string) (int64, bool) {
	parts := strings.Split(value, ":")
	if len(parts) != 3 {
		return 0, false
	}
	hours, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, false
	}
	minutes, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, false
	}
	seconds, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, false
	}
	totalSeconds := float64(hours*3600+minutes*60) + seconds
	return int64(totalSeconds * 1_000_000), true
}

func ffmpegProgressArgs(args []string) []string {
	next := append([]string{}, args...)
	output := next[len(next)-1]
	next = next[:len(next)-1]
	next = append(next, "-nostats", "-progress", "pipe:1", output)
	return next
}

func mediaFactDurationMs(item storage.MediaItem, track storage.MediaFileTrackFact) int64 {
	for _, fact := range item.FileFacts {
		if fact.FilePath == track.FilePath && fact.DurationMs != nil {
			return *fact.DurationMs
		}
	}
	return 0
}

func audioTranscodeDurationMs(item storage.MediaItem, track storage.MediaFileTrackFact) int64 {
	if duration := mediaFactDurationMs(item, track); duration > 0 {
		return duration
	}
	probe := delivery.Probe(track.FilePath)
	if probe.DurationSeconds == nil || *probe.DurationSeconds <= 0 {
		return 0
	}
	return int64(*probe.DurationSeconds * 1000)
}

func audioTranscodeProgressData(item storage.MediaItem, track storage.MediaFileTrackFact, durationMs int64) map[string]any {
	return map[string]any{
		"mediaItemId": item.ID.String(),
		"mediaTitle":  item.Title,
		"filePath":    track.FilePath,
		"trackId":     track.ID.String(),
		"phase":       "audio_transcode",
		"durationMs":  durationMs,
	}
}
