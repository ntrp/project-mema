import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { audioTrackMatchesTarget } from '$lib/components/app/media/files/mediaFileAudioTargetMatching';
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
	const audioTracks = row.tracks.filter((track) => track.type === 'audio');
	const targets =
		row.expectedAudioTargets.length > 0
			? row.expectedAudioTargets
			: row.expectedRequiredLanguages.map((languageId) => ({ languageId }));
	if (targets.length === 0) return [];
	return targets
		.filter(
			(target) =>
				!audioTracks.some(
					(track) =>
						languageMatchKey(track.language) === languageMatchKey(target.languageId) &&
						audioTrackMatchesTarget(track, target)
				)
		)
		.map((target) => ({
			key: `missing-audio-${targetKey(target)}`,
			trackNumber: '-',
			type: 'audio' as const,
			language: displayLanguage(target.languageId),
			description: 'Missing expected audio track',
			missing: true
		}));
}

function targetKey(target: MediaFileRow['expectedAudioTargets'][number]) {
	return [
		languageMatchKey(target.languageId) || target.languageId,
		target.targetCodec,
		target.targetChannels?.join('-'),
		target.minimumBitrateKbps
	]
		.filter(Boolean)
		.join('-');
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
