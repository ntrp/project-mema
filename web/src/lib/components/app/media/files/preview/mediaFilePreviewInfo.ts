import type { components } from '$lib/api/generated/schema';

export type MediaFilePreviewInfo = components['schemas']['MediaFilePreviewInfo'];
type MediaFileTrack = components['schemas']['MediaFileTrack'];

export interface MediaFileInfoRow {
	label: string;
	value: string;
}

export interface MediaFileInfoSection {
	key: string;
	title: string;
	action?: 'copy' | 'transcode';
	rows: MediaFileInfoRow[];
}

export interface MediaFilePlaybackStats {
	playing: boolean;
	variableBitRate: boolean;
	liveBitRate?: string;
	activeSubtitleLabel?: string;
}

export function mediaFileInfoSections(
	info?: MediaFilePreviewInfo,
	playbackStats?: MediaFilePlaybackStats
): MediaFileInfoSection[] {
	const sections: MediaFileInfoSection[] = [
		{
			key: 'video',
			title: 'Video',
			action: mediaTrackAction(info?.outputVideoCodec),
			rows: videoRows(info, playbackStats)
		},
		{
			key: 'audio',
			title: 'Audio',
			action: mediaTrackAction(info?.outputAudioCodec),
			rows: mediaTrackRows(info?.audioTrack, info?.outputAudioCodec, [
				['Language', (track) => track.language],
				['Title', (track) => track.title],
				['Channels', audioChannels],
				['Layout', (track) => track.channelLayout],
				['Bitrate', (track) => formatBitRate(track.bitRate)]
			])
		}
	];
	if (playbackStats?.activeSubtitleLabel) {
		sections.push({
			key: 'subtitle',
			title: 'Subtitle',
			rows: [{ label: 'Track', value: playbackStats.activeSubtitleLabel }]
		});
	}
	return sections;
}

function videoRows(info?: MediaFilePreviewInfo, playbackStats?: MediaFilePlaybackStats) {
	const rows = mediaTrackRows(info?.videoTrack, info?.outputVideoCodec, [
		['Resolution', videoResolution],
		['Profile', (track) => track.profile],
		['Pixel format', (track) => track.pixelFormat],
		['Frame rate', (track) => formatFrameRate(track.frameRate)],
		['Bitrate', (track) => videoBitRate(track, playbackStats)]
	]);
	return rows;
}

export function mediaTrackActionClass(action: NonNullable<MediaFileInfoSection['action']>) {
	if (action === 'transcode') return 'border-amber-300/35 bg-amber-400/20 text-amber-100';
	return 'border-emerald-300/35 bg-emerald-400/20 text-emerald-100';
}

export function formatBitRate(value?: string) {
	const bits = Number(value);
	if (!Number.isFinite(bits) || bits <= 0) return '-';
	if (bits >= 1_000_000) return `${(bits / 1_000_000).toFixed(bits >= 10_000_000 ? 1 : 2)} Mbps`;
	if (bits >= 1_000) return `${(bits / 1_000).toFixed(bits >= 10_000 ? 1 : 2)} Kbps`;
	return `${bits} bps`;
}

function mediaTrackRows(
	track: MediaFileTrack | undefined,
	outputCodec: string | undefined,
	fields: [string, (_track: MediaFileTrack) => string | number | undefined][]
): MediaFileInfoRow[] {
	const rows = [{ label: 'Codec', value: codecValue(track?.codec, outputCodec) }];
	for (const [label, value] of fields) {
		rows.push({
			label,
			value: formatValue(track ? value(track) : undefined)
		});
	}
	return rows;
}

function mediaTrackAction(outputCodec?: string): MediaFileInfoSection['action'] {
	return outputCodec && outputCodec !== 'copy' ? 'transcode' : 'copy';
}

function codecValue(sourceCodec?: string, outputCodec?: string) {
	const source = displayCodec(sourceCodec);
	const output = displayCodec(outputCodec);
	if (!outputCodec) return source;
	if (outputCodec === 'copy') return source === '-' ? 'copy' : source;
	if (source !== '-' && output !== '-' && source.toLowerCase() !== output.toLowerCase()) {
		return `${source} -> ${output}`;
	}
	return output;
}

function displayCodec(codec?: string) {
	const value = codec?.trim();
	if (!value) return '-';
	if (value === 'libx264') return 'h264';
	return value;
}

function videoResolution(track: MediaFileTrack) {
	if (!track.width || !track.height) return undefined;
	return `${track.width}x${track.height}`;
}

function videoBitRate(track: MediaFileTrack, playbackStats?: MediaFilePlaybackStats) {
	if (playbackStats?.playing && playbackStats.variableBitRate && playbackStats.liveBitRate) {
		return formatBitRate(playbackStats.liveBitRate);
	}
	return formatBitRate(track.bitRate);
}

function formatFrameRate(value?: string) {
	if (!value) return undefined;
	const [numerator, denominator = '1'] = value.split('/');
	const top = Number(numerator);
	const bottom = Number(denominator);
	if (!Number.isFinite(top) || !Number.isFinite(bottom) || top <= 0 || bottom <= 0) return value;
	const fps = top / bottom;
	const formatted = fps.toFixed(fps % 1 === 0 ? 0 : 3).replace(/\.?0+$/, '');
	return `${formatted} fps`;
}

function audioChannels(track: MediaFileTrack) {
	return track.channels ? `${track.channels}` : undefined;
}

function formatValue(value?: string | number) {
	if (value === undefined || value === '') return '-';
	return String(value);
}
