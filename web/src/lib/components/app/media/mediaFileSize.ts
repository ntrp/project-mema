import type { MediaItem } from '$lib/settings/types';

export function mediaFileInfo(item: MediaItem, path: string) {
	return item.files?.find((file) => file.path === path);
}

export function mediaFileSize(item: MediaItem, path: string) {
	const size = mediaFileInfo(item, path)?.sizeBytes;
	return typeof size === 'number' ? formatBytes(size) : '-';
}

export function formatBytes(bytes: number) {
	if (!Number.isFinite(bytes) || bytes < 0) {
		return '-';
	}
	if (bytes < 1024) {
		return `${bytes} B`;
	}
	const units = ['KiB', 'MiB', 'GiB', 'TiB'];
	let value = bytes / 1024;
	let unit = units[0];
	for (const nextUnit of units.slice(1)) {
		if (value < 1024) break;
		value /= 1024;
		unit = nextUnit;
	}
	return `${value.toFixed(value >= 10 ? 1 : 2)} ${unit}`;
}
