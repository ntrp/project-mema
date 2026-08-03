import type { MediaItemStatus, MediaMetadataDetails } from '$lib/settings/types';
import { mediaImageUrl } from '$lib/components/app/media/shared/mediaImage';

export function imageUrl(path?: string, size = 'w780') {
	return mediaImageUrl(path, size);
}

export function runtimeText(minutes?: number) {
	if (!minutes || minutes <= 0) return undefined;
	const hours = Math.floor(minutes / 60);
	const remainingMinutes = minutes % 60;
	if (hours > 0 && remainingMinutes > 0) return `${hours}h ${remainingMinutes}m`;
	if (hours > 0) return `${hours}h`;
	return `${remainingMinutes}m`;
}

export function mediaHeroTopInfo(details: MediaMetadataDetails) {
	return [
		details.seasonCount ? ['Seasons', `${details.seasonCount}`] : undefined,
		details.episodeCount ? ['Episodes', `${details.episodeCount}`] : undefined
	].filter((item): item is [string, string] => Boolean(item));
}

export function statusLabel(status: MediaItemStatus) {
	switch (status) {
		case 'downloaded':
			return 'Downloaded';
		case 'downloading':
			return 'Downloading';
		default:
			return 'Missing';
	}
}

export function statusBadgeClass(status: MediaItemStatus) {
	switch (status) {
		case 'downloaded':
			return 'border-emerald-500 text-emerald-300';
		case 'downloading':
			return 'border-primary text-primary';
		default:
			return 'border-destructive text-destructive';
	}
}

export function monitorStatus(details: MediaMetadataDetails) {
	return details.monitored ? 'Monitored' : 'Not monitored';
}

export function monitorHint(details: MediaMetadataDetails) {
	if (details.type === 'serie') {
		return details.monitored
			? 'Click to stop monitoring future episodes'
			: 'Click to monitor future episodes';
	}
	return details.monitored ? 'Click to stop monitoring this movie' : 'Click to monitor this movie';
}
