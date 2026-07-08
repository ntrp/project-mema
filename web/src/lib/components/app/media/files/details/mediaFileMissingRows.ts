import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';

type TrackType = MediaFileDetailRow['type'];

export function rowsWithMissingSubtitles(
	row: MediaFileRow,
	rows: MediaFileDetailRow[]
): MediaFileDetailRow[] {
	const missing = missingSubtitleRows(row);
	if (missing.length === 0) return rows;
	const insertAt = missingSubtitleInsertIndex(rows);
	return [...rows.slice(0, insertAt), ...missing, ...rows.slice(insertAt)];
}

export function rowsWithMissingAudio(
	row: MediaFileRow,
	rows: MediaFileDetailRow[]
): MediaFileDetailRow[] {
	const missing = missingAudioRows(row, rows);
	if (missing.length === 0) return rows;
	const insertAt = missingAudioInsertIndex(rows);
	return [...rows.slice(0, insertAt), ...missing, ...rows.slice(insertAt)];
}

function missingAudioRows(row: MediaFileRow, rows: MediaFileDetailRow[]): MediaFileDetailRow[] {
	if (row.expectedRequiredLanguages.length === 0) return [];
	const audioLanguages = new Set(
		rows
			.filter((track) => track.type === 'audio')
			.map((track) => languageMatchKey(track.language))
			.filter(Boolean)
	);
	return row.expectedRequiredLanguages
		.filter((language) => !audioLanguages.has(languageMatchKey(language)))
		.map((language) => ({
			key: `missing-audio-${languageMatchKey(language)}`,
			trackNumber: '-',
			type: 'audio' as const,
			language: displayLanguage(language),
			description: 'Missing expected audio track',
			missing: true
		}));
}

function missingAudioInsertIndex(rows: MediaFileDetailRow[]) {
	const lastAudio = lastTrackIndex(rows, 'audio');
	if (lastAudio >= 0) return lastAudio + 1;
	const lastVideo = lastTrackIndex(rows, 'video');
	return lastVideo >= 0 ? lastVideo + 1 : rows.length;
}

function missingSubtitleRows(row: MediaFileRow): MediaFileDetailRow[] {
	const satisfaction = row.subtitleSatisfaction;
	if (
		!satisfaction ||
		satisfaction.mode === 'external' ||
		satisfaction.missingLanguages.length === 0
	) {
		return [];
	}
	return satisfaction.missingLanguages.map((language) => ({
		key: `missing-subtitle-${languageMatchKey(language) || language}`,
		trackNumber: '',
		type: 'subtitle' as const,
		language: displayLanguage(language),
		description: 'Missing expected subtitle track',
		missing: true
	}));
}

function missingSubtitleInsertIndex(rows: MediaFileDetailRow[]) {
	const lastAudio = lastTrackIndex(rows, 'audio');
	if (lastAudio >= 0) return lastAudio + 1;
	const lastVideo = lastTrackIndex(rows, 'video');
	return lastVideo >= 0 ? lastVideo + 1 : rows.length;
}

function lastTrackIndex(rows: MediaFileDetailRow[], type: TrackType) {
	for (let index = rows.length - 1; index >= 0; index -= 1) {
		if (rows[index].type === type) return index;
	}
	return -1;
}
