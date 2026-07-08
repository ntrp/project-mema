import type { MediaItem, MediaProfileAudioTarget } from '$lib/settings/types';

export type MediaFileAudioTargetOption = Omit<MediaProfileAudioTarget, 'score'> & {
	score?: number;
	required?: boolean;
};
export interface MediaFileSubtitleTargetOption {
	languageId: string;
	formats?: string[];
}

export interface MediaFileProfileOption {
	id: string;
	qualityIds?: string[];
	upgradesAllowed?: boolean;
	upgradeUntilQualityId?: string;
	audioTargets?: MediaFileAudioTargetOption[];
	subtitleTargets?: MediaFileSubtitleTargetOption[];
	removeUnwantedAudio?: boolean;
	removeUnwantedSubtitles?: boolean;
}

export function fileProfileSettings(item: MediaItem, qualityProfiles: MediaFileProfileOption[]) {
	const profile = item.qualityProfileId
		? qualityProfiles.find((value) => value.id === item.qualityProfileId)
		: undefined;
	return { profile };
}
