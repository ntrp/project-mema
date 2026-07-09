package jobs

import (
	"github.com/riverqueue/river"

	"media-manager/internal/events"
	"media-manager/internal/storage"
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

type ContainerRemuxWorker struct {
	river.WorkerDefaults[ContainerRemuxArgs]
	settings           *storage.SettingsStore
	events             *events.Broker
	enqueueFulfillment fulfillmentEnqueueFunc
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
