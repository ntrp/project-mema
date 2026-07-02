import type { MediaItem } from '$lib/settings/types';

export interface MediaFileProfileOption {
	id: string;
	targetLanguages?: string[];
}

export function expectedProfileLanguages(
	item: MediaItem,
	qualityProfiles: MediaFileProfileOption[]
) {
	if (!item.qualityProfileId) return [];
	return (
		qualityProfiles.find((profile) => profile.id === item.qualityProfileId)?.targetLanguages ?? []
	);
}
