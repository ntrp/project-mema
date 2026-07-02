import type { MediaItem } from '$lib/settings/types';

export interface MediaFileProfileOption {
	id: string;
	targetLanguages?: string[];
	removeNonEnabledLanguages?: boolean;
}

export function fileProfileSettings(item: MediaItem, qualityProfiles: MediaFileProfileOption[]) {
	const profile = item.qualityProfileId
		? qualityProfiles.find((value) => value.id === item.qualityProfileId)
		: undefined;
	return {
		expectedLanguages: profile?.targetLanguages ?? [],
		removeNonEnabledLanguages: profile?.removeNonEnabledLanguages === true
	};
}
