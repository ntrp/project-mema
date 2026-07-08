import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import type {
	MediaFileDetailRow,
	TrackDeleteRequest
} from '$lib/components/app/media/files/mediaFileDetailRows';
import { relativePath } from '$lib/components/app/media/files/mediaFilePath';
import {
	rowsWithMissingAudio,
	rowsWithMissingSubtitles
} from '$lib/components/app/media/files/details/mediaFileMissingRows';
import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';

type MediaFileTrack = MediaFileRow['tracks'][number];
type MediaFileChapter = MediaFileRow['chapters'][number];

export type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetailRows';

export function fileDetailRows(row: MediaFileRow): MediaFileDetailRow[] {
	return [...trackRowsWithMissingTargets(row), ...fileChapterDetailRows(row)];
}

export function fileTrackDetailRows(row: MediaFileRow): MediaFileDetailRow[] {
	return trackRowsWithMissingTargets(row);
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
		chapterSummary: true,
		deleteRequest: { targetType: 'chapters' }
	};
}

function trackRow(row: MediaFileRow, track: MediaFileTrack, index: number): MediaFileDetailRow {
	return {
		key: `${track.type}-${track.index ?? index}`,
		trackNumber: track.index === undefined ? String(index + 1) : String(track.index),
		type: track.type,
		language: displayLanguage(track.language),
		description: trackDescription(row, track),
		provenance: track.provenance,
		unwanted: unwantedTrack(row, track),
		deleteRequest: trackDeleteRequest(track)
	};
}

function trackRowsWithMissingTargets(row: MediaFileRow): MediaFileDetailRow[] {
	const rows = [
		...row.tracks.map((track, index) => trackRow(row, track, index)),
		...externalSubtitleTrackRows(row)
	];
	return rowsWithMissingSubtitles(row, rowsWithMissingAudio(row, rows));
}

function externalSubtitleTrackRows(row: MediaFileRow): MediaFileDetailRow[] {
	if (row.subtitleSatisfaction?.mode !== 'embedded') return [];
	return (row.externalSubtitles ?? [])
		.filter((subtitle) => subtitle.retentionMode === 'mux')
		.map((subtitle) => ({
			key: `external-subtitle-${subtitle.id}`,
			trackNumber: '-',
			type: 'subtitle' as const,
			language: displayLanguage(subtitle.languageId),
			description: compactParts([
				'External subtitle',
				subtitle.format.toUpperCase(),
				relativePath(row.path ? row.path.replace(/[^/]+$/, '') : undefined, subtitle.filePath)
			]),
			unwanted: unwantedSubtitleLanguage(row, subtitle.languageId)
		}));
}

function unwantedTrack(row: MediaFileRow, track: MediaFileTrack) {
	if (track.type === 'subtitle' && row.subtitleSatisfaction?.mode === 'external') {
		return true;
	}
	const expectedLanguages = wantedLanguagesForTrack(row, track);
	if (expectedLanguages.length === 0) {
		return false;
	}
	const enabled = new Set(expectedLanguages.map(languageMatchKey).filter(Boolean));
	const language = languageMatchKey(track.language);
	return language !== '' && !enabled.has(language);
}

function trackDeleteRequest(track: MediaFileTrack): TrackDeleteRequest | undefined {
	if ((track.type !== 'audio' && track.type !== 'subtitle') || track.index === undefined) {
		return undefined;
	}
	return { targetType: track.type, trackIndex: track.index };
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

function unwantedSubtitleLanguage(row: MediaFileRow, languageId: string) {
	if (!row.removeNonEnabledSubtitleLanguages || row.expectedSubtitleLanguages.length === 0) {
		return false;
	}
	const enabled = new Set(row.expectedSubtitleLanguages.map(languageMatchKey).filter(Boolean));
	const language = languageMatchKey(languageId);
	return language !== '' && !enabled.has(language);
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
		]),
		deleteRequest: { targetType: 'chapter', chapterIndex: chapter.index }
	};
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
