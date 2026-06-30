import type { MediaMonitorMode, MinimumAvailability } from '$lib/settings/types';

export interface MediaActionSelection {
	qualityProfileId?: string;
	libraryFolderId?: string;
	tags: string[];
	monitorMode: MediaMonitorMode;
	minimumAvailability: MinimumAvailability;
	startSearch: boolean;
}
