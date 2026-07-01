import type { MediaItem, MediaSearchResult, TagForm } from '$lib/settings/types';

export function omitResult<TValue>(results: Record<string, TValue | undefined>, id: string) {
	const { [id]: _removed, ...remaining } = results;
	return remaining;
}

export function candidateKey(candidate: MediaSearchResult) {
	return `${candidate.type}:${candidate.title}:${candidate.year ?? ''}`;
}

export function mediaItemFileCount(item: MediaItem) {
	return (item.filePaths?.length ?? 0) + (item.metadataFilePaths?.length ?? 0);
}

export function errorMessageFrom(error: unknown, fallback: string) {
	return error instanceof Error ? error.message : fallback;
}

export function emptyTagForm(): TagForm {
	return { name: '' };
}
