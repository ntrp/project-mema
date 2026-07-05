import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';

export type MediaFileSummaryStatusState = 'ignored' | 'missing' | 'satisfied';

export interface MediaFileSummaryStatus {
	state: MediaFileSummaryStatusState;
	label: string;
}

export function audioSatisfaction(row: MediaFileRow): MediaFileSummaryStatus {
	if (!row.exists) return { state: 'missing', label: 'Missing: -' };
	const missing = missingRequiredAudio(row);
	if (missing.length > 0) return { state: 'missing', label: `Missing: ${languageList(missing)}` };
	return { state: 'satisfied', label: `Satisfied: ${languageList(row.expectedRequiredLanguages)}` };
}

export function subtitleSatisfaction(row: MediaFileRow): MediaFileSummaryStatus {
	const satisfaction = row.subtitleSatisfaction;
	if (!satisfaction || satisfaction.state === 'ignored') {
		return { state: 'ignored', label: 'Ignored' };
	}
	if (satisfaction.state === 'satisfied') {
		return { state: 'satisfied', label: `Satisfied: ${languageList(satisfaction.matchedLanguages)}` };
	}
	return { state: 'missing', label: `Missing: ${languageList(satisfaction.missingLanguages)}` };
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

function languageList(values: string[]) {
	return values.map(displayLanguage).join(', ') || '-';
}
