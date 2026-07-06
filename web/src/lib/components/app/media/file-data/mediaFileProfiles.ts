import type { MediaItem } from '$lib/settings/types';

export interface MediaFileProfileOption {
	id: string;
	qualityIds?: string[];
	upgradesAllowed?: boolean;
	upgradeUntilQualityId?: string;
	audioTargets?: { languageId: string; required?: boolean }[];
	subtitleTargets?: { languageId: string }[];
	removeUnwantedAudio?: boolean;
	removeUnwantedSubtitles?: boolean;
}

export function fileProfileSettings(item: MediaItem, qualityProfiles: MediaFileProfileOption[]) {
	const profile = item.qualityProfileId
		? qualityProfiles.find((value) => value.id === item.qualityProfileId)
		: undefined;
	return {
		profile,
		expectedLanguages: profile?.audioTargets?.map((target) => target.languageId) ?? [],
		expectedRequiredLanguages: requiredTargetLanguages(profile),
		expectedSubtitleLanguages:
			profile?.subtitleTargets?.map((language) => language.languageId) ?? [],
		removeNonEnabledLanguages: profile?.removeUnwantedAudio === true,
		removeNonEnabledSubtitleLanguages: profile?.removeUnwantedSubtitles === true
	};
}

function requiredTargetLanguages(profile: MediaFileProfileOption | undefined) {
	if (!profile) {
		return [];
	}
	if (!profile.audioTargets?.length) {
		return [];
	}
	return profile.audioTargets
		.filter((language) => language.required)
		.map((language) => language.languageId);
}
