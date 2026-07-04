import type Player from 'video.js/dist/types/player';
import type { MediaFilePlaybackStats } from '$lib/components/app/media/files/preview/mediaFilePreviewInfo';

interface ByteCountVideoElement extends globalThis.HTMLVideoElement {
	webkitAudioDecodedByteCount?: number;
	webkitVideoDecodedByteCount?: number;
}

interface VideoJsTextTrack {
	kind?: string;
	label?: string;
	language?: string;
	mode?: string;
}

interface VideoJsTextTrackList {
	length: number;
	on?: (_type: string, _handler: () => void) => void;
	off?: (_type: string, _handler: () => void) => void;
	[index: number]: VideoJsTextTrack;
}

const emptyStats: MediaFilePlaybackStats = { playing: false, variableBitRate: false };

export function emptyPlaybackStats(): MediaFilePlaybackStats {
	return emptyStats;
}

export function watchPlaybackStats(
	player: Player,
	video: globalThis.HTMLVideoElement,
	source: string,
	onChange: (_stats: MediaFilePlaybackStats) => void
) {
	let latest = emptyStats;
	let previousBytes: number | undefined;
	let previousStamp = 0;
	const samples: number[] = [];
	const update = () => {
		const next = nextPlaybackStats(player, video, source, previousBytes, previousStamp, samples);
		previousBytes = next.bytes;
		previousStamp = next.stamp;
		if (!sameStats(latest, next.stats)) {
			latest = next.stats;
			onChange(next.stats);
		}
	};
	const timer = globalThis.window.setInterval(update, 1000);
	const tracks = player.textTracks() as unknown as VideoJsTextTrackList;
	const events = ['play', 'pause', 'ended', 'loadedmetadata'];
	for (const event of events) video.addEventListener(event, update);
	tracks.on?.('change', update);
	onChange(emptyStats);
	update();
	return () => {
		globalThis.window.clearInterval(timer);
		for (const event of events) video.removeEventListener(event, update);
		tracks.off?.('change', update);
	};
}

function nextPlaybackStats(
	player: Player,
	video: globalThis.HTMLVideoElement,
	source: string,
	previousBytes: number | undefined,
	previousStamp: number,
	samples: number[]
) {
	const stamp = globalThis.performance.now();
	const bytes = currentByteCount(video, source);
	const liveBitRate = liveBitRateSample(bytes, previousBytes, stamp, previousStamp, samples);
	return {
		bytes,
		stamp,
		stats: {
			playing: !video.paused && !video.ended,
			variableBitRate: variableBitRate(samples),
			liveBitRate,
			activeSubtitleLabel: activeSubtitleLabel(player)
		}
	};
}

function currentByteCount(video: globalThis.HTMLVideoElement, source: string) {
	return resourceByteCount(source) ?? decodedByteCount(video);
}

function resourceByteCount(source: string) {
	const url = new URL(source, globalThis.window.location.href).href;
	const entries = globalThis.performance.getEntriesByName(
		url,
		'resource'
	) as globalThis.PerformanceResourceTiming[];
	const entry = entries.at(-1);
	if (!entry) return undefined;
	const bytes = entry.encodedBodySize || entry.transferSize || entry.decodedBodySize;
	return bytes > 0 ? bytes : undefined;
}

function decodedByteCount(video: globalThis.HTMLVideoElement) {
	const element = video as ByteCountVideoElement;
	const bytes =
		(element.webkitVideoDecodedByteCount ?? 0) + (element.webkitAudioDecodedByteCount ?? 0);
	return bytes > 0 ? bytes : undefined;
}

function liveBitRateSample(
	bytes: number | undefined,
	previousBytes: number | undefined,
	stamp: number,
	previousStamp: number,
	samples: number[]
) {
	if (bytes === undefined || previousBytes === undefined || stamp <= previousStamp)
		return undefined;
	const byteDelta = bytes - previousBytes;
	const secondDelta = (stamp - previousStamp) / 1000;
	if (byteDelta <= 0 || secondDelta <= 0) return undefined;
	const bitRate = Math.round((byteDelta * 8) / secondDelta);
	samples.push(bitRate);
	if (samples.length > 6) samples.shift();
	return String(bitRate);
}

function variableBitRate(samples: number[]) {
	if (samples.length < 3) return false;
	const min = Math.min(...samples);
	const max = Math.max(...samples);
	return min > 0 && max / min > 1.15;
}

function activeSubtitleLabel(player: Player) {
	const tracks = player.textTracks() as unknown as VideoJsTextTrackList;
	for (let index = 0; index < tracks.length; index += 1) {
		const track = tracks[index];
		if (track?.kind === 'subtitles' && track.mode === 'showing') {
			return [track.label, track.language].filter(Boolean).join(' / ') || 'Subtitle';
		}
	}
	return undefined;
}

function sameStats(a: MediaFilePlaybackStats, b: MediaFilePlaybackStats) {
	return (
		a.playing === b.playing &&
		a.variableBitRate === b.variableBitRate &&
		a.liveBitRate === b.liveBitRate &&
		a.activeSubtitleLabel === b.activeSubtitleLabel
	);
}
