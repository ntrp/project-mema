import { relativePath } from '$lib/components/app/media/files/mediaFilePath';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';
import type { MediaItemSubtitle } from '$lib/settings/types';

type MediaFileTrack = MediaFileRow['tracks'][number];

export interface SubtitleStateRow {
	key: string;
	languageId: string;
	language: string;
	state: 'embedded' | 'external' | 'missing' | 'downloading' | 'satisfied';
	label: string;
}

export interface EmbeddedSubtitleRow {
	key: string;
	languageId: string;
	language: string;
	description: string;
}

export function subtitleStateRows(row: MediaFileRow, downloading = false): SubtitleStateRow[] {
	const satisfaction = row.subtitleSatisfaction;
	if (!satisfaction || satisfaction.state === 'ignored') return [];
	return satisfaction.wantedLanguages.map((languageId) => {
		const match = subtitleMatch(row, languageId);
		const missing = satisfaction.missingLanguages.some((value) => sameLanguage(value, languageId));
		const state =
			downloading && missing ? 'downloading' : (match ?? (missing ? 'missing' : 'satisfied'));
		return {
			key: languageMatchKey(languageId) || languageId,
			languageId,
			language: displayLanguage(languageId),
			state,
			label: stateLabel(state)
		};
	});
}

export function embeddedSubtitleRows(row: MediaFileRow): EmbeddedSubtitleRow[] {
	return row.tracks.filter(isSubtitleTrack).map((track, index) => ({
		key: `embedded-${track.index ?? index}`,
		languageId: track.language ?? '',
		language: displayLanguage(track.language),
		description: compactParts([track.codec, track.title])
	}));
}

export function externalSubtitlesForRow(row: MediaFileRow): MediaItemSubtitle[] {
	return row.externalSubtitles ?? [];
}

export function subtitleSourceLabel(subtitle: MediaItemSubtitle) {
	return compactParts([
		subtitle.providerName,
		subtitle.sourceReference,
		subtitle.providerSubtitleId
	]);
}

export function subtitleFileLabel(row: MediaFileRow, subtitle: MediaItemSubtitle) {
	return relativePath(row.path ? row.path.replace(/[^/]+$/, '') : undefined, subtitle.filePath);
}

function subtitleMatch(
	row: MediaFileRow,
	languageId: string
): SubtitleStateRow['state'] | undefined {
	if (
		externalSubtitlesForRow(row).some((subtitle) => sameLanguage(subtitle.languageId, languageId))
	) {
		return 'external';
	}
	if (
		row.tracks.some((track) => isSubtitleTrack(track) && sameLanguage(track.language, languageId))
	) {
		return 'embedded';
	}
	return undefined;
}

function isSubtitleTrack(track: MediaFileTrack) {
	return track.type === 'subtitle';
}

function sameLanguage(left?: string, right?: string) {
	const leftKey = languageMatchKey(left);
	const rightKey = languageMatchKey(right);
	return leftKey !== '' && leftKey === rightKey;
}

function stateLabel(state: SubtitleStateRow['state']) {
	switch (state) {
		case 'embedded':
			return 'Embedded';
		case 'external':
			return 'External';
		case 'downloading':
			return 'Downloading';
		case 'missing':
			return 'Missing';
		default:
			return 'Satisfied';
	}
}

function compactParts(values: (string | undefined)[]) {
	return (
		values
			.map((value) => value?.trim())
			.filter(Boolean)
			.join(' · ') || '-'
	);
}
