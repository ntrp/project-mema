package satisfaction

import "media-manager/internal/targets"

type MediaFileState string

const (
	MediaFileMissing     MediaFileState = "missing"
	MediaFileDownloading MediaFileState = "downloading"
	MediaFilePartial     MediaFileState = "partial"
	MediaFileDownloaded  MediaFileState = "downloaded"
	MediaFileUpgradeable MediaFileState = "upgradeable"
)

type MediaItemState string

const (
	MediaItemMissing     MediaItemState = "missing"
	MediaItemDownloading MediaItemState = "downloading"
	MediaItemPartial     MediaItemState = "partial"
	MediaItemDownloaded  MediaItemState = "downloaded"
	MediaItemUpgradeable MediaItemState = "upgradeable"
)

func RollupMediaFileState(hasUsableFile bool, targetStates []targets.State, activeWork bool) MediaFileState {
	if !hasUsableFile {
		if activeWork {
			return MediaFileDownloading
		}
		return MediaFileMissing
	}
	if hasAnyTargetState(targetStates, targets.StateBlocked, targets.StateFailed, targets.StateMissing, targets.StatePartial, targets.StatePending) {
		return MediaFilePartial
	}
	if hasAnyTargetState(targetStates, targets.StateUpgradeable) {
		return MediaFileUpgradeable
	}
	if activeWork {
		return MediaFileDownloading
	}
	return MediaFileDownloaded
}

func RollupMediaItemState(requiredFileStates []MediaFileState, activeWork bool) MediaItemState {
	if len(requiredFileStates) == 0 {
		if activeWork {
			return MediaItemDownloading
		}
		return MediaItemMissing
	}
	if hasAnyFileState(requiredFileStates, MediaFilePartial) {
		return MediaItemPartial
	}
	if hasAnyFileState(requiredFileStates, MediaFileMissing) {
		if activeWork {
			return MediaItemDownloading
		}
		return MediaItemMissing
	}
	if hasAnyFileState(requiredFileStates, MediaFileDownloading) || activeWork {
		return MediaItemDownloading
	}
	if hasAnyFileState(requiredFileStates, MediaFileUpgradeable) {
		return MediaItemUpgradeable
	}
	return MediaItemDownloaded
}

func hasAnyTargetState(states []targets.State, candidates ...targets.State) bool {
	for _, state := range states {
		for _, candidate := range candidates {
			if state == candidate {
				return true
			}
		}
	}
	return false
}

func hasAnyFileState(states []MediaFileState, candidates ...MediaFileState) bool {
	for _, state := range states {
		for _, candidate := range candidates {
			if state == candidate {
				return true
			}
		}
	}
	return false
}
