import { displayLanguage, languageMatchKey } from '$lib/settings/languageDisplay';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import { mediaFileChapterTrack } from '$lib/components/app/media/files/preview/mediaFileChapters';

type MediaFileTrack = MediaFileRow['tracks'][number];

export interface AudioTrackOption {
	key: string;
	label: string;
	enabled: boolean;
	language?: string;
	streamIndex?: number;
}

export interface MediaFileTextTrack {
	key: string;
	kind: 'subtitles' | 'chapters';
	label: string;
	src: string;
	default?: boolean;
	srclang?: string;
}

export function mediaFilePreviewUrl(
	mediaItemId: string,
	filePath: string,
	audioTrackIndex?: number,
	startTimeSeconds?: number
) {
	return mediaFileUrl(mediaItemId, filePath, 'preview', {
		audioTrackIndex,
		startTimeSeconds: validStartTime(startTimeSeconds) ? startTimeSeconds.toFixed(3) : undefined
	});
}
export function mediaFilePreviewInfoUrl(
	mediaItemId: string,
	filePath: string,
	audioTrackIndex?: number
) {
	return mediaFileUrl(mediaItemId, filePath, 'preview-info', { audioTrackIndex });
}
export function mediaFileVlcUrl(mediaItemId: string, filePath: string) {
	return mediaFileUrl(mediaItemId, filePath, 'vlc');
}
export function mediaFileTextTracks(mediaItemId: string, row: MediaFileRow): MediaFileTextTrack[] {
	const subtitles = row.path ? mediaFileSubtitleTracks(mediaItemId, row.path, row.tracks) : [];
	const chapters = mediaFileChapterTrack(row);
	return chapters ? [...subtitles, chapters] : subtitles;
}
export function metadataAudioTrackOptions(row: MediaFileRow): AudioTrackOption[] {
	return row.tracks
		.filter((track) => track.type === 'audio')
		.map((track, index) => ({
			key: `metadata-${track.index ?? index}`,
			label: audioTrackLabel(track, index),
			enabled: index === 0,
			language: languageTag(track.language),
			streamIndex: track.index
		}));
}

function mediaFileSubtitleTracks(
	mediaItemId: string,
	filePath: string,
	tracks: MediaFileTrack[]
): MediaFileTextTrack[] {
	return tracks
		.filter((track) => track.type === 'subtitle' && track.index !== undefined)
		.map((track, index) => ({
			key: `subtitle-${track.index}`,
			kind: 'subtitles',
			label: subtitleTrackLabel(track, index),
			src: mediaFileUrl(mediaItemId, filePath, 'subtitle', { subtitleTrackIndex: track.index }),
			srclang: languageTag(track.language)
		}));
}

function mediaFileUrl(
	mediaItemId: string,
	filePath: string,
	action: 'preview' | 'preview-info' | 'vlc' | 'subtitle',
	params: Record<string, string | number | undefined> = {}
) {
	const query = new URLSearchParams({ path: filePath });
	for (const [name, value] of Object.entries(params)) {
		if (value !== undefined) query.set(name, String(value));
	}
	return `/api/media/items/${encodeURIComponent(mediaItemId)}/files/${action}?${query}`;
}

function validStartTime(value: number | undefined): value is number {
	return typeof value === 'number' && Number.isFinite(value) && value > 0;
}

function audioTrackLabel(track: MediaFileTrack, index: number) {
	return compactParts([
		track.title,
		displayLanguagePart(track.language),
		track.codec,
		track.channelLayout,
		track.channels ? `${track.channels}ch` : undefined,
		`Track ${track.index ?? index + 1}`
	]);
}

function subtitleTrackLabel(track: MediaFileTrack, index: number) {
	return compactParts([
		track.title,
		displayLanguagePart(track.language),
		track.codec,
		`Track ${track.index ?? index + 1}`
	]);
}

function displayLanguagePart(value?: string) {
	const label = displayLanguage(value);
	return label === '-' ? undefined : label;
}

function languageTag(value?: string) {
	return languageMatchKey(value) || 'und';
}

function compactParts(values: (string | undefined)[]) {
	return (
		values
			.map((value) => value?.trim())
			.filter(Boolean)
			.join(' · ') || '-'
	);
}
