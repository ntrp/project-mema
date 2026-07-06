import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
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
	chapterSummary?: boolean;
	missing?: boolean;
	unwanted?: boolean;
}

export function fileDetailRows(row: MediaFileRow): MediaFileDetailRow[] {
	const details = [
		...row.tracks.map((track, index) => trackRow(row, track, index)),
		...fileChapterDetailRows(row)
	];
	return [...details, ...missingLanguageRows(row, details)];
}

export function fileTrackDetailRows(row: MediaFileRow): MediaFileDetailRow[] {
	const details = row.tracks.map((track, index) => trackRow(row, track, index));
	return [...details, ...missingLanguageRows(row, details)];
}

export function fileChapterDetailRows(row: MediaFileRow): MediaFileDetailRow[] {
	return row.chapters.map(chapterRow);
}

export function fileChapterSummaryRow(row: MediaFileRow): MediaFileDetailRow | undefined {
	if (row.chapters.length === 0) return undefined;
	const first = row.chapters[0];
	const last = row.chapters[row.chapters.length - 1];
	const firstTrackNumber = first.index + 1;
	const lastTrackNumber = last.index + 1;
	return {
		key: 'chapter-summary',
		trackNumber:
			firstTrackNumber === lastTrackNumber
				? String(firstTrackNumber)
				: `${firstTrackNumber}-${lastTrackNumber}`,
		type: 'chapter',
		language: '-',
		description: chapterCountLabel(row.chapters.length),
		chapterSummary: true
	};
}

function trackRow(row: MediaFileRow, track: MediaFileTrack, index: number): MediaFileDetailRow {
	return {
		key: `${track.type}-${track.index ?? index}`,
		trackNumber: track.index === undefined ? String(index + 1) : String(track.index),
		type: track.type,
		language: displayLanguage(track.language),
		description: trackDescription(row, track),
		unwanted: unwantedTrack(row, track)
	};
}

function unwantedTrack(row: MediaFileRow, track: MediaFileTrack) {
	const expectedLanguages = wantedLanguagesForTrack(row, track);
	if (expectedLanguages.length === 0) {
		return false;
	}
	const enabled = new Set(expectedLanguages.map(languageMatchKey).filter(Boolean));
	const language = languageMatchKey(track.language);
	return language !== '' && !enabled.has(language);
}

function wantedLanguagesForTrack(row: MediaFileRow, track: MediaFileTrack) {
	if (track.type === 'audio' && row.removeNonEnabledLanguages) {
		return row.expectedLanguages;
	}
	if (track.type === 'subtitle' && row.removeNonEnabledSubtitleLanguages) {
		return row.expectedSubtitleLanguages;
	}
	return [];
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
	if (row.expectedRequiredLanguages.length === 0) return [];
	const audioLanguages = new Set(
		existingRows
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

function chapterCountLabel(count: number) {
	return `${count} ${count === 1 ? 'chapter' : 'chapters'}`;
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
