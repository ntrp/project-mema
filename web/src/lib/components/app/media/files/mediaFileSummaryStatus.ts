import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';
export { audioSatisfaction } from '$lib/components/app/media/files/mediaFileAudioStatus';

export type MediaFileSummaryStatusState = 'ignored' | 'missing' | 'partial' | 'satisfied';

export interface MediaFileSummaryStatus {
	state: MediaFileSummaryStatusState;
	label: string;
	details: string[];
}

export function subtitleSatisfaction(row: MediaFileRow): MediaFileSummaryStatus {
	const satisfaction = row.subtitleSatisfaction;
	if (!satisfaction || satisfaction.state === 'ignored') {
		return { state: 'ignored', label: 'Ignored', details: ['Subtitle requirements are ignored'] };
	}
	const details = [`Mode: ${subtitleModeLabel(satisfaction.mode)}`];
	if (satisfaction.state === 'satisfied') {
		return okStatus([
			...details,
			`Subtitle requirements met: ${languageList(satisfaction.matchedLanguages)}`
		]);
	}
	const externalMissing = externallyAvailableMissingSubtitles(row);
	const missing = satisfaction.missingLanguages.filter(
		(language) => !languageKeys(externalMissing).has(languageMatchKey(language))
	);
	if (externalMissing.length > 0) {
		return partialStatus([
			...details,
			`Subtitles need to be imported: ${languageList(externalMissing)}`,
			...(missing.length > 0 ? [`Missing subtitles: ${languageList(missing)}`] : [])
		]);
	}
	details.push(`Missing subtitles: ${languageList(satisfaction.missingLanguages)}`);
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

function externallyAvailableMissingSubtitles(row: MediaFileRow) {
	const satisfaction = row.subtitleSatisfaction;
	if (!satisfaction || satisfaction.mode !== 'embedded') return [];
	const external = externalSubtitleLanguageKeys(row);
	return satisfaction.missingLanguages.filter((language) =>
		external.has(languageMatchKey(language))
	);
}

function externalSubtitleLanguageKeys(row: MediaFileRow) {
	return new Set(
		[
			...(row.externalSubtitles ?? []).map((subtitle) => languageMatchKey(subtitle.languageId)),
			...row.otherFiles
				.filter((file) => file.type === 'subtitle' && file.status === 'available')
				.map((file) => languageMatchKey(file.language))
		].filter(Boolean)
	);
}

function languageKeys(values: string[]) {
	return new Set(values.map(languageMatchKey).filter(Boolean));
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

function subtitleModeLabel(value: string) {
	switch (value) {
		case 'embedded':
			return 'Embedded';
		case 'external':
			return 'External';
		default:
			return 'Mixed';
	}
}

function languageList(values: string[]) {
	return values.map(displayLanguage).join(', ') || '-';
}
