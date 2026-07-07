import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';

type MediaFileAudioTrack = MediaFileRow['tracks'][number];
type MediaFileAudioTarget = MediaFileRow['expectedAudioTargets'][number];

export function audioTrackMatchesTarget(track: MediaFileAudioTrack, target: MediaFileAudioTarget) {
	return audioTargetMismatchDetails(track, target).length === 0;
}

export function audioTargetMismatchDetails(
	track: MediaFileAudioTrack,
	target: MediaFileAudioTarget
) {
	const details: string[] = [];
	if (
		target.targetCodec &&
		normalizeAudioCodec(track.codec) !== normalizeAudioCodec(target.targetCodec)
	) {
		details.push(`codec ${track.codec ?? 'unknown'} != ${target.targetCodec}`);
	}
	if (
		target.targetChannels?.length &&
		!target.targetChannels.some((value) => channelMatches(track, value))
	) {
		details.push(`channels ${channelLabel(track)} not in ${target.targetChannels.join(', ')}`);
	}
	const bitrate = bitrateKbps(track.bitRate);
	if (target.minimumBitrateKbps && (!bitrate || bitrate < target.minimumBitrateKbps)) {
		details.push(
			`bitrate ${bitrate ? `${bitrate} kbps` : 'unknown'} below ${target.minimumBitrateKbps} kbps`
		);
	}
	return details;
}

function channelMatches(track: MediaFileAudioTrack, target: string) {
	const normalized = normalizeChannel(target);
	if (!normalized) return false;
	if (normalized === 'atmos') return channelLabel(track).includes('atmos');
	return trackChannelValues(track).includes(normalized);
}

function trackChannelValues(track: MediaFileAudioTrack) {
	return [normalizeChannel(track.channelLayout), channelCountLabel(track.channels)].filter(Boolean);
}

function channelLabel(track: MediaFileAudioTrack) {
	return track.channelLayout ?? channelCountLabel(track.channels) ?? 'unknown';
}

function channelCountLabel(value: number | undefined) {
	switch (value) {
		case 1:
			return '1.0';
		case 2:
			return '2.0';
		case 6:
			return '5.1';
		case 7:
			return '6.1';
		case 8:
			return '7.1';
		default:
			return value ? `${value}ch` : undefined;
	}
}

function normalizeChannel(value: string | undefined) {
	const normalized = value?.trim().toLowerCase();
	if (!normalized) return undefined;
	if (normalized.includes('atmos')) return 'atmos';
	if (normalized.includes('5.1')) return '5.1';
	if (normalized.includes('6.1')) return '6.1';
	if (normalized.includes('7.1')) return '7.1';
	if (normalized === 'stereo') return '2.0';
	if (normalized === 'mono') return '1.0';
	return normalized;
}

function bitrateKbps(value: string | undefined) {
	const bitrate = Number(value);
	return Number.isFinite(bitrate) && bitrate > 0 ? Math.round(bitrate / 1000) : undefined;
}

function normalizeAudioCodec(value: string | undefined) {
	const trimmed = value?.trim() ?? '';
	if (trimmed.toLowerCase() === 'dd+') return 'eac3';
	switch (trimmed.toLowerCase().replace(/[^a-z0-9]/g, '')) {
		case 'ddp':
		case 'ddplus':
		case 'eac3':
			return 'eac3';
		case 'dd':
		case 'ac3':
		case 'dolbydigital':
			return 'ac3';
		case 'truehdatmos':
			return 'truehd';
		default:
			return trimmed.toLowerCase();
	}
}
