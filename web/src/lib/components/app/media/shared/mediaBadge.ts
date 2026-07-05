import type { MediaType } from '$lib/settings/types';

export type MediaBadgeTone = MediaType | 'series';

export function mediaBadgeToneClass(type: MediaBadgeTone) {
	return type === 'movie'
		? 'border-yellow-500 bg-yellow-400 text-yellow-950'
		: 'border-blue-600 bg-blue-600 text-white';
}
