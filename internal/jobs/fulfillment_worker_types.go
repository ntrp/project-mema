package jobs

import (
	"github.com/riverqueue/river"

	"media-manager/internal/events"
	"media-manager/internal/storage"
	"media-manager/internal/subtitles"
)

type MediaFulfillmentWorker struct {
	river.WorkerDefaults[MediaFulfillmentArgs]
	settings           *storage.SettingsStore
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
}

type VideoTranscodeWorker struct {
	river.WorkerDefaults[VideoTranscodeArgs]
	settings           *storage.SettingsStore
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
}

type AudioTranscodeWorker struct {
	river.WorkerDefaults[AudioTranscodeArgs]
	settings           *storage.SettingsStore
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
}

type AudioSourceWorker struct {
	river.WorkerDefaults[AudioSourceArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

type ContainerRemuxWorker struct {
	river.WorkerDefaults[ContainerRemuxArgs]
	settings           *storage.SettingsStore
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
}

type SubtitleDownloadWorker struct {
	river.WorkerDefaults[SubtitleDownloadArgs]
	settings  *storage.SettingsStore
	subtitles *subtitles.Service
	events    *events.Broker
}

type SubtitleEmbedWorker struct {
	river.WorkerDefaults[SubtitleEmbedArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

type SubtitleExtractWorker struct {
	river.WorkerDefaults[SubtitleExtractArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}

type SubtitleConvertWorker struct {
	river.WorkerDefaults[SubtitleConvertArgs]
	settings *storage.SettingsStore
	events   *events.Broker
}
