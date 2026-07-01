import type { MediaMonitorMode, MinimumAvailability, SeriesType } from '$lib/settings/types';

export interface MediaActionSelection {
	qualityProfileId?: string;
	libraryFolderId?: string;
	tags: string[];
	monitorMode: MediaMonitorMode;
	seriesType: SeriesType;
	minimumAvailability: MinimumAvailability;
	startSearch: boolean;
}
