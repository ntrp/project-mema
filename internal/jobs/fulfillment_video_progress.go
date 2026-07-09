package jobs

import (
	"context"
	"strconv"
	"time"

	"media-manager/internal/delivery"
	"media-manager/internal/events"
	"media-manager/internal/storage"
	mediatools "media-manager/internal/tools"
)

func initializeQueuedVideoTranscodeProgress(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	jobID int64,
	item storage.MediaItem,
	track storage.MediaFileTrackFact,
) {
	if settings == nil {
		return
	}
	zero := int32(0)
	progress := normalizedProgressData(&zero, "Waiting to transcode video", videoTranscodeProgressData(item, track, mediaFactDurationMs(item, track)))
	execution, err := settings.UpdateSystemJobExecutionProgressData(ctx, jobID, &zero, "Waiting to transcode video", progress)
	if err == nil {
		publishJobExecutionUpdated(eventBroker, execution)
	}
}

func runVideoTranscodeCommand(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	track storage.MediaFileTrackFact,
	argsList []string,
) error {
	progress := videoTranscodeProgress{durationMs: videoTranscodeDurationMs(item, track)}
	start := int32(0)
	recordJobProgressData(ctx, settings, eventBroker, &start, "Transcoding video 0%", videoTranscodeProgressData(item, track, progress.durationMs))
	_, err := mediatools.RunOutputProgress(ctx, mediatools.ProgressCommandSpec{
		CommandSpec: mediatools.CommandSpec{
			Name:           "ffmpeg",
			Args:           ffmpegProgressArgs(argsList),
			Timeout:        6 * time.Hour,
			MaxOutputBytes: 0,
			MaxStderrBytes: 128 * 1024,
		},
		Progress: func(line string) {
			if percent, ok := progress.percent(line); ok {
				recordJobProgressData(ctx, settings, eventBroker, &percent, "Transcoding video "+strconv.Itoa(int(percent))+"%", videoTranscodeProgressData(item, track, progress.durationMs))
			}
		},
	})
	if err != nil {
		return err
	}
	done := int32(100)
	recordJobProgressData(ctx, settings, eventBroker, &done, "Video transcode complete", videoTranscodeProgressData(item, track, progress.durationMs))
	return nil
}

type videoTranscodeProgress struct {
	durationMs int64
	last       int32
	lastAt     time.Time
}

func (p *videoTranscodeProgress) percent(line string) (int32, bool) {
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

func videoTranscodeDurationMs(item storage.MediaItem, track storage.MediaFileTrackFact) int64 {
	if duration := mediaFactDurationMs(item, track); duration > 0 {
		return duration
	}
	probe := delivery.Probe(track.FilePath)
	if probe.DurationSeconds == nil || *probe.DurationSeconds <= 0 {
		return 0
	}
	return int64(*probe.DurationSeconds * 1000)
}

func videoTranscodeProgressData(item storage.MediaItem, track storage.MediaFileTrackFact, durationMs int64) map[string]any {
	return map[string]any{
		"mediaItemId": item.ID.String(),
		"mediaTitle":  item.Title,
		"filePath":    track.FilePath,
		"trackId":     track.ID.String(),
		"phase":       "video_transcode",
		"durationMs":  durationMs,
	}
}
