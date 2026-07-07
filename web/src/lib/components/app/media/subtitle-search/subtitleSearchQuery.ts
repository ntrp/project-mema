import type { MediaFileRow } from '$lib/components/app/media/files/mediaFiles';
import type { MediaItem } from '$lib/settings/types';

export function subtitleSearchQuery(item: MediaItem) {
	return item.title.trim();
}

export function subtitleSearchQueryVariants(item: MediaItem, row: MediaFileRow) {
	const values: string[] = [];
	const add = (value: string | undefined) => {
		value = value?.trim();
		if (!value || values.some((item) => item.toLowerCase() === value.toLowerCase())) return;
		values.push(value);
	};
	add(item.title);
	add([item.title, item.year].filter(Boolean).join(' '));
	add(
		row.path
			?.split('/')
			.pop()
			?.replace(/\.[^.]+$/, '')
	);
	return values;
}
