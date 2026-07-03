import { activityDisplay, releaseGroupFromTitle } from './activityDisplay';
import type { DownloadActivity, ManualImportRequest } from '$lib/settings/types';

export interface ManualImportFormState {
	sourcePath: string;
	targetFileName: string;
	movieTitle: string;
	year?: number;
	seasonNumber?: number;
	episodeNumber?: number;
	episodeTitle: string;
	releaseGroup: string;
	edition: string;
	quality: string;
	languagesText: string;
}

export function initialManualImportForm(activity: DownloadActivity): ManualImportFormState {
	const summary = activityDisplay(activity);

	return {
		sourcePath: '',
		targetFileName: '',
		movieTitle: activity.mediaTitle,
		year: activity.mediaYear ?? undefined,
		seasonNumber: activity.mediaType === 'series' ? 1 : undefined,
		episodeNumber: activity.mediaType === 'series' ? 1 : undefined,
		episodeTitle: '',
		releaseGroup: releaseGroupFromTitle(activity.releaseTitle),
		edition: '',
		quality: summary.quality === '-' ? '' : summary.quality,
		languagesText: summary.languages.join(', ')
	};
}

export function manualImportRequestFromForm(form: ManualImportFormState): ManualImportRequest {
	return {
		sourcePath: form.sourcePath,
		targetFileName: optional(form.targetFileName),
		movieTitle: optional(form.movieTitle),
		year: form.year,
		seasonNumber: form.seasonNumber,
		episodeNumber: form.episodeNumber,
		episodeTitle: optional(form.episodeTitle),
		releaseGroup: optional(form.releaseGroup),
		edition: optional(form.edition),
		quality: optional(form.quality),
		languages: form.languagesText
			.split(',')
			.map((value) => value.trim())
			.filter(Boolean)
	};
}

function optional(value: string) {
	value = value.trim();
	return value === '' ? undefined : value;
}
