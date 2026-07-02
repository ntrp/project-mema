import type { MediaFileRow } from './mediaFiles';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';

type MediaFileTrack = MediaFileRow['tracks'][number];
type MediaFileChapter = MediaFileRow['chapters'][number];
type TrackType = MediaFileTrack['type'] | 'chapter';

export interface MediaFileDetailRow {
	key: string;
	trackNumber: string;
	type: TrackType;
	language: string;
	description: string;
	missing?: boolean;
}

export function fileDetailRows(row: MediaFileRow): MediaFileDetailRow[] {
	const details = [
		...row.tracks.map((track, index) => trackRow(row, track, index)),
		...row.chapters.map(chapterRow)
	];
	return [...details, ...missingLanguageRows(row, details)];
}

function trackRow(row: MediaFileRow, track: MediaFileTrack, index: number): MediaFileDetailRow {
	return {
		key: `${track.type}-${track.index ?? index}`,
		trackNumber: track.index === undefined ? String(index + 1) : String(track.index),
		type: track.type,
		language: displayLanguage(track.language),
		description: trackDescription(row, track)
	};
}

function chapterRow(chapter: MediaFileChapter): MediaFileDetailRow {
	return {
		key: `chapter-${chapter.index}`,
		trackNumber: String(chapter.index + 1),
		type: 'chapter',
		language: '-',
		description: compactParts([
			valueOrDash(chapter.title),
			rangeLabel(chapter.startTime, chapter.endTime)
		])
	};
}

function missingLanguageRows(
	row: MediaFileRow,
	existingRows: MediaFileDetailRow[]
): MediaFileDetailRow[] {
	if (row.expectedLanguages.length < 2) return [];
	const audioLanguages = new Set(
		existingRows
			.filter((track) => track.type === 'audio')
			.map((track) => languageMatchKey(track.language))
			.filter(Boolean)
	);
	return row.expectedLanguages
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

function trackDescription(row: MediaFileRow, track: MediaFileTrack) {
	if (track.type === 'video') {
		return compactParts([
			track.codec,
			resolution(track, row.quality),
			track.profile,
			track.pixelFormat,
			track.frameRate,
			bitRate(track.bitRate)
		]);
	}
	if (track.type === 'audio') {
		return compactParts([track.codec, channels(track), bitRate(track.bitRate), track.title]);
	}
	return compactParts([track.codec, track.title]);
}

function resolution(track: MediaFileTrack, fallback: string) {
	return track.width && track.height ? `${track.width}x${track.height}` : fallback;
}

function channels(track: MediaFileTrack) {
	if (track.channelLayout) return track.channelLayout;
	return track.channels ? `${track.channels}ch` : undefined;
}

function bitRate(value?: string) {
	const numeric = Number(value);
	return Number.isFinite(numeric) && numeric > 0 ? `${Math.round(numeric / 1000)} kbps` : undefined;
}

function rangeLabel(start?: string, end?: string) {
	if (!start && !end) return undefined;
	return compactParts([start, end], ' - ');
}

function compactParts(values: (string | undefined)[], separator = ' · ') {
	return (
		values
			.map((value) => value?.trim())
			.filter(Boolean)
			.join(separator) || '-'
	);
}

function valueOrDash(value?: string) {
	return value?.trim() || '-';
}
