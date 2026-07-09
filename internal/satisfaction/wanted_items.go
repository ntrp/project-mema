package satisfaction

import (
	"media-manager/internal/storage"
	"media-manager/internal/targets"
)

func BuildWantedRowsForItem(item storage.MediaItem) []WantedRow {
	input := WantedRowsInput{
		Item:           item,
		HasUsableMedia: len(item.FileFacts) > 0 || len(item.FilePaths) > 0,
	}
	input.Targets = WantedTargetInputsForItem(item)
	return BuildWantedRows(input)
}

func WantedTargetInputsForItem(item storage.MediaItem) []WantedTargetInput {
	profile := mediaProfileForItem(item)
	inputs := []WantedTargetInput{}
	for _, fact := range item.FileFacts {
		video := EvaluateVideoTarget(item, &profile, fact)
		inputs = append(inputs, WantedTargetInput{Target: video.Target, FilePath: fact.FilePath})
		for _, result := range EvaluateAudioTargets(item, &profile, fact).Results {
			inputs = append(inputs, WantedTargetInput{Target: result.Target, FilePath: fact.FilePath})
		}
		for _, result := range EvaluateSubtitleTargets(item, &profile, fact).Results {
			inputs = append(inputs, WantedTargetInput{Target: result.Target, FilePath: fact.FilePath})
		}
	}
	return inputs
}

func WantedTargetsForItem(item storage.MediaItem) []targets.Target {
	result := []targets.Target{}
	for _, input := range WantedTargetInputsForItem(item) {
		result = append(result, input.Target)
	}
	return result
}

func mediaProfileForItem(item storage.MediaItem) storage.MediaProfile {
	return storage.MediaProfile{
		FinalContainer:                "",
		RemoveUnwantedAudio:           item.RemoveUnwantedAudio,
		AudioLossyTranscodePolicy:     "lossyToLossy",
		RemoveUnwantedSubtitles:       item.RemoveUnwantedSubtitles,
		SubtitleMode:                  item.SubtitleMode,
		AllowSubtitleReleaseFallback:  item.AllowSubtitleReleaseFallback,
		VideoTarget:                   item.VideoTarget,
		AudioTargets:                  item.AudioTargets,
		SubtitleTargets:               item.SubtitleTargets,
		QualityIDs:                    nil,
		UpgradeUntilCustomFormatScore: 0,
	}
}
