package jobs

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"media-manager/internal/delivery"
	"media-manager/internal/storage"
	"media-manager/internal/subtitleformats"
)

func persistLiveMediaFileFact(ctx context.Context, settings *storage.SettingsStore, item storage.MediaItem, filePath string) error {
	probe := delivery.Probe(filePath)
	if len(probe.Tracks) == 0 {
		return fmt.Errorf("probe returned no tracks for %s", filePath)
	}
	input := liveMediaFileFactInput(item, filePath, probe)
	if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
		size := info.Size()
		input.SizeBytes = &size
	}
	_, err := settings.UpsertMediaFileFact(ctx, input)
	return err
}

func liveMediaFileFactInput(item storage.MediaItem, filePath string, probe delivery.ProbeResult) storage.MediaFileFactInput {
	input := storage.MediaFileFactInput{
		MediaItemID:         item.ID,
		FilePath:            filePath,
		ContainerFormat:     probe.Container.Format,
		ContainerFormatName: probe.Container.FormatName,
		ContainerBitrate:    int64FromString(probe.Container.BitRate),
		DurationMs:          durationMsFromSeconds(probe.DurationSeconds),
		SourceKind:          "probe",
		ProbedAt:            time.Now().UTC(),
	}
	for _, fact := range item.FileFacts {
		if fact.FilePath == filePath {
			input.SeasonID = fact.SeasonID
			input.EpisodeID = fact.EpisodeID
			input.QualityID = fact.QualityID
			break
		}
	}
	for _, track := range probe.Tracks {
		input.Tracks = append(input.Tracks, liveMediaFileTrackInput(track))
	}
	return input
}

func liveMediaFileTrackInput(track delivery.Track) storage.MediaFileTrackFactInput {
	return storage.MediaFileTrackFactInput{
		StreamIndex: int32Value(track.Index),
		TrackType:   string(track.Type),
		LanguageID:  track.Language,
		Codec:       track.Codec,
		Channels:    track.ChannelLayout,
		DurationMs:  durationMsFromSeconds(track.Duration),
		BitrateKbps: int32KbpsFromString(track.BitRate),
		Width:       track.Width,
		Height:      track.Height,
		PixelFormat: track.PixelFormat,
		Format:      subtitleFormatFromTrack(track),
		Title:       track.Title,
	}
}

func subtitleFormatFromTrack(track delivery.Track) *string {
	if track.Type != delivery.TrackSubtitle || track.Codec == nil {
		return nil
	}
	normalized := subtitleformats.Normalize(*track.Codec)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func int32Value(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}

func durationMsFromSeconds(value *float64) *int64 {
	if value == nil || *value <= 0 {
		return nil
	}
	ms := int64(*value * 1000)
	return &ms
}

func int64FromString(value *string) *int64 {
	if value == nil {
		return nil
	}
	parsed, err := strconv.ParseInt(*value, 10, 64)
	if err != nil || parsed < 0 {
		return nil
	}
	return &parsed
}

func int32KbpsFromString(value *string) *int32 {
	bits := int64FromString(value)
	if bits == nil || *bits <= 0 {
		return nil
	}
	kbps := int32(*bits / 1000)
	return &kbps
}
