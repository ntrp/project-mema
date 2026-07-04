import { displayLanguage } from '$lib/settings/languageDisplay';
import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

type MediaFileTrack = MediaFileRow['tracks'][number];

export interface AudioTrackOption {
	key: string;
	label: string;
	enabled: boolean;
	streamIndex?: number;
}

export function mediaFileStreamUrl(mediaItemId: string, filePath: string) {
	return mediaFileUrl(mediaItemId, filePath, 'stream');
}

export function mediaFilePreviewUrl(
	mediaItemId: string,
	filePath: string,
	audioTrackIndex?: number
) {
	return mediaFileUrl(mediaItemId, filePath, 'preview', audioTrackIndex);
}

export function mediaFileVlcUrl(mediaItemId: string, filePath: string) {
	return mediaFileUrl(mediaItemId, filePath, 'vlc');
}

export function metadataAudioTrackOptions(row: MediaFileRow): AudioTrackOption[] {
	return row.tracks
		.filter((track) => track.type === 'audio')
		.map((track, index) => ({
			key: `metadata-${track.index ?? index}`,
			label: audioTrackLabel(track, index),
			enabled: index === 0,
			streamIndex: track.index
		}));
}

export function chapterSeconds(value?: string) {
	const trimmed = value?.trim();
	if (!trimmed) return undefined;
	const numeric = Number(trimmed);
	if (Number.isFinite(numeric)) return numeric;
	const parts = trimmed.split(':').map(Number);
	if (parts.some((part) => !Number.isFinite(part))) return undefined;
	return parts.reduce((total, part) => total * 60 + part, 0);
}

export function formatPlaybackTime(seconds: number) {
	if (!Number.isFinite(seconds) || seconds < 0) return '0:00';
	const rounded = Math.floor(seconds);
	const hours = Math.floor(rounded / 3600);
	const minutes = Math.floor((rounded % 3600) / 60);
	const rest = String(rounded % 60).padStart(2, '0');
	return hours > 0 ? `${hours}:${String(minutes).padStart(2, '0')}:${rest}` : `${minutes}:${rest}`;
}

export function playlistDownloadName(fileName: string) {
	const base = fileName.trim().replace(/\.[^.]+$/, '') || 'media-stream';
	return `${base}.m3u`;
}

function mediaFileUrl(
	mediaItemId: string,
	filePath: string,
	action: 'preview' | 'stream' | 'vlc',
	audioTrackIndex?: number
) {
	const query = new URLSearchParams({ path: filePath });
	if (audioTrackIndex !== undefined) {
		query.set('audioTrackIndex', String(audioTrackIndex));
	}
	return `/api/media/items/${encodeURIComponent(mediaItemId)}/files/${action}?${query}`;
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

function displayLanguagePart(value?: string) {
	const label = displayLanguage(value);
	return label === '-' ? undefined : label;
}

function compactParts(values: (string | undefined)[]) {
	return (
		values
			.map((value) => value?.trim())
			.filter(Boolean)
			.join(' · ') || '-'
	);
}
