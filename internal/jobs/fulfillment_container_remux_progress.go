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

func initializeQueuedContainerRemuxProgress(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	jobID int64,
	item storage.MediaItem,
	fact storage.MediaFileFact,
) {
	if settings == nil {
		return
	}
	zero := int32(0)
	progress := normalizedProgressData(&zero, "Waiting to remux container", containerRemuxProgressData(item, fact, mediaFactDuration(fact)))
	execution, err := settings.UpdateSystemJobExecutionProgressData(ctx, jobID, &zero, "Waiting to remux container", progress)
	if err == nil {
		publishJobExecutionUpdated(eventBroker, execution)
	}
}

func runContainerRemuxCommand(
	ctx context.Context,
	settings *storage.SettingsStore,
	eventBroker *events.Broker,
	item storage.MediaItem,
	fact storage.MediaFileFact,
	argsList []string,
) error {
	progress := containerRemuxProgress{durationMs: containerRemuxDurationMs(fact)}
	start := int32(0)
	recordJobProgressData(ctx, settings, eventBroker, &start, "Remuxing container 0%", containerRemuxProgressData(item, fact, progress.durationMs))
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
				recordJobProgressData(ctx, settings, eventBroker, &percent, "Remuxing container "+strconv.Itoa(int(percent))+"%", containerRemuxProgressData(item, fact, progress.durationMs))
			}
		},
	})
	if err != nil {
		return err
	}
	done := int32(100)
	recordJobProgressData(ctx, settings, eventBroker, &done, "Container remux complete", containerRemuxProgressData(item, fact, progress.durationMs))
	return nil
}

type containerRemuxProgress struct {
	durationMs int64
	last       int32
	lastAt     time.Time
}

func (p *containerRemuxProgress) percent(line string) (int32, bool) {
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

func containerRemuxDurationMs(fact storage.MediaFileFact) int64 {
	if duration := mediaFactDuration(fact); duration > 0 {
		return duration
	}
	probe := delivery.Probe(fact.FilePath)
	if probe.DurationSeconds == nil || *probe.DurationSeconds <= 0 {
		return 0
	}
	return int64(*probe.DurationSeconds * 1000)
}

func mediaFactDuration(fact storage.MediaFileFact) int64 {
	if fact.DurationMs == nil {
		return 0
	}
	return *fact.DurationMs
}

func containerRemuxProgressData(item storage.MediaItem, fact storage.MediaFileFact, durationMs int64) map[string]any {
	return map[string]any{
		"mediaItemId": item.ID.String(),
		"mediaTitle":  item.Title,
		"filePath":    fact.FilePath,
		"phase":       "container_remux",
		"durationMs":  durationMs,
	}
}
