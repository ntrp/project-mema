import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';

export type MediaFileSummaryStatusState = 'ignored' | 'missing' | 'partial' | 'satisfied';

export interface MediaFileSummaryStatus {
	state: MediaFileSummaryStatusState;
	label: string;
	details: string[];
}

export function audioSatisfaction(row: MediaFileRow): MediaFileSummaryStatus {
	if (!row.exists) return missingStatus(['File is missing']);
	if (!hasTrack(row, 'video')) return missingStatus(['Video track is missing']);
	if (!hasTrack(row, 'audio')) return missingStatus(['Audio track is missing']);
	const missing = missingRequiredAudio(row);
	if (missing.length > 0) {
		const matched = matchedRequiredAudio(row);
		const details = [`Missing required audio: ${languageList(missing)}`];
		return matched.length > 0 ? partialStatus(details) : missingStatus(details);
	}
	const unwanted = unwantedAudio(row);
	if (unwanted.length > 0) {
		return partialStatus([`Unwanted audio tracks: ${languageList(unwanted)}`]);
	}
	return okStatus([
		row.expectedRequiredLanguages.length > 0
			? `Required audio present: ${languageList(row.expectedRequiredLanguages)}`
			: 'Audio and video are available'
	]);
}

export function subtitleSatisfaction(row: MediaFileRow): MediaFileSummaryStatus {
	const satisfaction = row.subtitleSatisfaction;
	if (!satisfaction || satisfaction.state === 'ignored') {
		return { state: 'ignored', label: 'Ignored', details: ['Subtitle requirements are ignored'] };
	}
	if (satisfaction.state === 'satisfied') {
		return okStatus([`Subtitle requirements met: ${languageList(satisfaction.matchedLanguages)}`]);
	}
	const details = [`Missing subtitles: ${languageList(satisfaction.missingLanguages)}`];
	return satisfaction.matchedLanguages.length > 0 ? partialStatus(details) : missingStatus(details);
}

export function statusBadgeClass(state: MediaFileSummaryStatusState) {
	switch (state) {
		case 'satisfied':
			return 'bg-emerald-600 text-white hover:bg-emerald-700';
		case 'missing':
			return undefined;
		case 'ignored':
			return undefined;
		default:
			return 'bg-amber-500 text-white hover:bg-amber-600';
	}
}

export function statusBadgeVariant(state: MediaFileSummaryStatusState) {
	return state === 'missing' ? 'destructive' : 'secondary';
}

function missingRequiredAudio(row: MediaFileRow) {
	const audioLanguages = new Set(
		row.tracks
			.filter((track) => track.type === 'audio')
			.map((track) => languageMatchKey(track.language))
			.filter(Boolean)
	);
	return row.expectedRequiredLanguages.filter(
		(language) => !audioLanguages.has(languageMatchKey(language))
	);
}

function matchedRequiredAudio(row: MediaFileRow) {
	const audioLanguages = audioLanguageKeys(row);
	return row.expectedRequiredLanguages.filter((language) =>
		audioLanguages.has(languageMatchKey(language))
	);
}

function unwantedAudio(row: MediaFileRow) {
	if (!row.removeNonEnabledLanguages || row.expectedLanguages.length === 0) return [];
	const expected = new Set(row.expectedLanguages.map(languageMatchKey).filter(Boolean));
	return uniqueLanguages(
		row.tracks
			.filter((track) => track.type === 'audio')
			.map((track) => track.language)
			.filter((language): language is string => Boolean(language))
			.filter((language) => {
				const key = languageMatchKey(language);
				return key !== '' && !expected.has(key);
			})
	);
}

function audioLanguageKeys(row: MediaFileRow) {
	return new Set(
		row.tracks
			.filter((track) => track.type === 'audio')
			.map((track) => languageMatchKey(track.language))
			.filter(Boolean)
	);
}

function hasTrack(row: MediaFileRow, type: 'audio' | 'video') {
	return row.tracks.some((track) => track.type === type);
}

function okStatus(details: string[]): MediaFileSummaryStatus {
	return { state: 'satisfied', label: 'Ok', details };
}

function partialStatus(details: string[]): MediaFileSummaryStatus {
	return { state: 'partial', label: 'Partial', details };
}

function missingStatus(details: string[]): MediaFileSummaryStatus {
	return { state: 'missing', label: 'Missing', details };
}

function languageList(values: string[]) {
	return values.map(displayLanguage).join(', ') || '-';
}

function uniqueLanguages(values: string[]) {
	return [...new Map(values.map((value) => [languageMatchKey(value) || value, value])).values()];
}
