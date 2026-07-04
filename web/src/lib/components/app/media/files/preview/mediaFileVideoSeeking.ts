import type Player from 'video.js/dist/types/player';
import type { TimeRange } from 'video.js/dist/types/utils/time';

export function addSourceTimeline(
	instance: Player,
	enabled: boolean,
	sourceStartTime: number,
	sourceDuration: number | undefined,
	onSeekRequest: (_timeSeconds: number) => void
) {
	if (!enabled) return undefined;
	const offset = Math.max(0, sourceStartTime);
	const duration = validDuration(sourceDuration) ? sourceDuration : undefined;
	const nativeCurrentTime = instance.currentTime;
	const nativeDuration = instance.duration;
	const nativeBuffered = instance.buffered;
	const nativeSeekable = instance.seekable;
	let lastRequest: number | undefined;
	let ready = false;
	const markReady = () => {
		ready = true;
		instance.trigger('timeupdate');
	};
	instance.one('loadedmetadata', markReady);
	instance.currentTime = ((seconds?: number | string) => {
		if (seconds === undefined) return sourceCurrentTime(nativeCurrentTime.call(instance), offset);
		const target = Number(seconds);
		if (!Number.isFinite(target)) return nativeCurrentTime.call(instance, seconds);
		const localTarget = Math.max(0, target - offset);
		if (!ready && Math.abs(target - offset) <= 0.5) return nativeCurrentTime.call(instance, 0);
		if (target < offset - 0.5) {
			if (target <= 0.5 && (nativeCurrentTime.call(instance) ?? 0) <= 0.5) {
				return nativeCurrentTime.call(instance, 0);
			}
			requestSourceSeek(target);
			return target;
		}
		if (!bufferedAround(nativeBuffered.call(instance), localTarget)) {
			requestSourceSeek(target);
			return target;
		}
		return nativeCurrentTime.call(instance, localTarget);
	}) as Player['currentTime'];
	instance.duration = ((seconds?: number) => {
		if (duration === undefined) return nativeDuration.call(instance, seconds);
		return duration;
	}) as Player['duration'];
	instance.buffered = (() =>
		shiftedRanges(nativeBuffered.call(instance), offset)) as Player['buffered'];
	instance.seekable = (() =>
		duration === undefined
			? shiftedRanges(nativeSeekable.call(instance), offset)
			: timeRanges([[0, duration]])) as Player['seekable'];
	return () => {
		instance.off('loadedmetadata', markReady);
		instance.currentTime = nativeCurrentTime;
		instance.duration = nativeDuration;
		instance.buffered = nativeBuffered;
		instance.seekable = nativeSeekable;
	};

	function requestSourceSeek(target: number) {
		if (!validSeekTarget(target) || Math.abs(target - (lastRequest ?? Number.NaN)) < 0.5) return;
		lastRequest = target;
		onSeekRequest(target);
	}
}

function validSeekTarget(value: number | undefined): value is number {
	return typeof value === 'number' && Number.isFinite(value) && value >= 0;
}

function validDuration(value: number | undefined): value is number {
	return typeof value === 'number' && Number.isFinite(value) && value > 0;
}

function sourceCurrentTime(value: number | undefined, offset: number) {
	return typeof value === 'number' && Number.isFinite(value) ? value + offset : value;
}

function bufferedAround(ranges: TimeRange, time: number) {
	for (let index = 0; index < ranges.length; index += 1) {
		if (time >= ranges.start(index) - 1 && time <= ranges.end(index) + 1) {
			return true;
		}
	}
	return false;
}

function shiftedRanges(ranges: TimeRange, offset: number): TimeRange {
	if (offset <= 0) return ranges;
	return timeRanges(
		Array.from({ length: ranges.length }, (_, index) => [
			ranges.start(index) + offset,
			ranges.end(index) + offset
		])
	);
}

function timeRanges(ranges: [number, number][]): TimeRange {
	return {
		length: ranges.length,
		start: (index: number) => ranges[index][0],
		end: (index: number) => ranges[index][1]
	};
}
