import type { DLNAMediaProbeRequest, MediaFilePreviewInfo } from '$lib/settings/types';

function traceValue(value?: string | null) {
	const trimmed = value?.trim() ?? '';
	return trimmed === '' ? undefined : trimmed;
}

export function probeFromMediaPreview(
	info?: MediaFilePreviewInfo | null
): DLNAMediaProbeRequest | undefined {
	const container = traceValue(info?.containerFormatName) ?? traceValue(info?.containerFormat);
	const videoCodec = traceValue(info?.videoTrack?.codec);
	const audioCodec = traceValue(info?.audioTrack?.codec);
	const height = info?.videoTrack?.height;

	if (!container && !videoCodec && !audioCodec && height == null) {
		return undefined;
	}

	return {
		container,
		videoCodec,
		audioCodec,
		height
	};
}
