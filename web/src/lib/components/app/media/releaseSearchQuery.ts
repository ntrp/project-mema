import type { MediaItem } from '$lib/settings/types';

export type ReleaseSearchContext =
	| { type: 'title' }
	| { type: 'season'; seasonNumber: number }
	| { type: 'episode'; seasonNumber?: number; episodeNumber?: number };

export function releaseSearchQuery(
	item: MediaItem,
	context: ReleaseSearchContext = { type: 'title' }
) {
	const title = item.title.trim();
	if (item.type === 'movie') return [title, item.year].filter(Boolean).join(' ');
	if (context.type === 'season') return `${title} s${context.seasonNumber}`;
	if (
		context.type === 'episode' &&
		context.seasonNumber !== undefined &&
		context.episodeNumber !== undefined
	) {
		return `${title} s${context.seasonNumber}e${context.episodeNumber}`;
	}
	return [title, item.year].filter(Boolean).join(' ');
}
