import type { MediaItem } from '$lib/settings/types';

export interface MediaFileProfileOption {
	id: string;
	qualityIds?: string[];
	upgradesAllowed?: boolean;
	upgradeUntilQualityId?: string;
	targetLanguages?: string[];
	targetLanguageScores?: { languageId: string; score?: number; required?: boolean }[];
	subtitleLanguages?: { languageId: string }[];
	removeNonEnabledLanguages?: boolean;
	removeNonEnabledSubtitleLanguages?: boolean;
}

export function fileProfileSettings(item: MediaItem, qualityProfiles: MediaFileProfileOption[]) {
	const profile = item.qualityProfileId
		? qualityProfiles.find((value) => value.id === item.qualityProfileId)
		: undefined;
	return {
		profile,
		expectedLanguages: profile?.targetLanguages ?? [],
		expectedRequiredLanguages: requiredTargetLanguages(profile),
		expectedSubtitleLanguages:
			profile?.subtitleLanguages?.map((language) => language.languageId) ?? [],
		removeNonEnabledLanguages: profile?.removeNonEnabledLanguages === true,
		removeNonEnabledSubtitleLanguages: profile?.removeNonEnabledSubtitleLanguages === true
	};
}

function requiredTargetLanguages(profile: MediaFileProfileOption | undefined) {
	if (!profile) {
		return [];
	}
	if (!profile.targetLanguageScores?.length) {
		return [];
	}
	return profile.targetLanguageScores
		.filter((language) => language.required)
		.map((language) => language.languageId);
}
