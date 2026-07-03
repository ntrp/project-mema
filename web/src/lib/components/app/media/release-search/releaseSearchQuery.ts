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

export function releaseSearchQueryVariants(
	item: MediaItem,
	context: ReleaseSearchContext = { type: 'title' }
) {
	const title = item.title.trim();
	const variants: string[] = [];
	const addVariant = (value: string) => {
		value = value.trim();
		if (!value || variants.some((variant) => variant.toLowerCase() === value.toLowerCase())) return;
		variants.push(value);
	};

	addVariant(releaseSearchQuery(item, context));
	if (item.type === 'series' && context.type === 'season') {
		addVariant(`${title} s${context.seasonNumber}`);
		addVariant(`${title} S${padded(context.seasonNumber, 2)}`);
	}
	if (
		item.type === 'series' &&
		context.type === 'episode' &&
		context.seasonNumber !== undefined &&
		context.episodeNumber !== undefined
	) {
		addVariant(`${title} s${context.seasonNumber}e${context.episodeNumber}`);
		addVariant(`${title} S${padded(context.seasonNumber, 2)}E${padded(context.episodeNumber, 2)}`);
	}
	return variants;
}

function padded(value: number, width: number) {
	return String(value).padStart(width, '0');
}
