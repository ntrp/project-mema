package targets

type ManualAction struct {
	ID            string
	Operation     OperationType
	Label         string
	Description   string
	Manual        bool
	Automatic     bool
	Available     bool
	BlockedReason string
	Method        string
	Path          string
	WorkerPath    string
	StateEffect   string
}

func OperationTypes() []OperationType {
	return []OperationType{
		OperationReleaseSearch,
		OperationVideoTranscode,
		OperationAudioTranscode,
		OperationAudioSourcing,
		OperationContainerRemux,
		OperationSubtitleDownload,
		OperationSubtitleEmbed,
		OperationSubtitleExtraction,
		OperationSubtitleConversion,
		OperationFileRescan,
	}
}

func ManualActions() []ManualAction {
	return []ManualAction{
		action("release_search", OperationReleaseSearch, "Search releases", "Search indexers for release candidates.", "POST", "/media/items/{id}/release-searches", "ReleaseSearchWorker"),
		action("release_grab", OperationReleaseSearch, "Grab release", "Send a selected release to a download client.", "POST", "/media/items/{id}/grab", "GrabReleaseWorker"),
		action("release_import", OperationReleaseSearch, "Import release", "Retry import for completed download activity.", "POST", "/activity/downloads/{id}/manual-import", "DownloadActivitySyncWorker"),
		action("video_upgrade", OperationReleaseSearch, "Search upgrade", "Run search and grab flow for a better video file.", "POST", "/media/items/{id}/automatic-searches", "AutoSearchDownloadWorker"),
		action("video_transcode", OperationVideoTranscode, "Transcode video", "Create a profile-matching video derivative.", "POST", "/media/items/{id}/assemblies", "ComponentMuxWorker"),
		action("audio_transcode", OperationAudioTranscode, "Transcode audio", "Create a profile-matching audio stream.", "POST", "/media/items/{id}/assemblies", "ComponentMuxWorker"),
		action("audio_source", OperationAudioSourcing, "Source audio", "Retain an alternate release as an audio source.", "POST", "/media/items/{id}/component-sources", "GrabReleaseWorker"),
		action("container_remux", OperationContainerRemux, "Remux container", "Move selected streams into the target container.", "POST", "/media/items/{id}/assemblies", "ComponentMuxWorker"),
		action("subtitle_download", OperationSubtitleDownload, "Download subtitle", "Search and download a subtitle candidate.", "POST", "/media/items/{id}/subtitle-searches", "SubtitleSearchWorker"),
		action("subtitle_grab", OperationSubtitleDownload, "Grab subtitle", "Download a selected manual subtitle result.", "POST", "/media/items/{id}/subtitle-grabs", "SubtitleSearchWorker"),
		action("subtitle_embed", OperationSubtitleEmbed, "Embed subtitle", "Merge an external subtitle into the media container.", "POST", "/media/items/{id}/assemblies", "ComponentMuxWorker"),
		action("subtitle_extraction", OperationSubtitleExtraction, "Extract subtitle", "Extract an embedded subtitle to an external artifact.", "POST", "/media/items/{id}/component-sources/{sourceId}/extractions", "ComponentExtractionWorker"),
		action("subtitle_conversion", OperationSubtitleConversion, "Convert subtitle", "Download or write subtitle content in a target format.", "POST", "/media/items/{id}/subtitle-searches", "SubtitleSearchWorker"),
		action("file_rescan", OperationFileRescan, "Rescan files", "Persist current file, track, subtitle, and sidecar facts.", "POST", "/media/items/{id}/files/rescan", "MediaFileRescan"),
	}
}

func action(id string, operation OperationType, label, description, method, path, workerPath string) ManualAction {
	return ManualAction{
		ID:          id,
		Operation:   operation,
		Label:       label,
		Description: description,
		Manual:      true,
		Automatic:   true,
		Available:   true,
		Method:      method,
		Path:        path,
		WorkerPath:  workerPath,
		StateEffect: "Uses the same persisted facts, logs, and target satisfaction recalculation as the automatic path.",
	}
}
