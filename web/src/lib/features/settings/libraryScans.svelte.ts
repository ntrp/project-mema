import { createQuery, type QueryClient } from '@tanstack/svelte-query';
import type { LibraryScan } from '$lib/settings/types';

export const libraryScanKeys = {
	all: ['library-scans'] as const,
	byFolder: (folderId: string) => [...libraryScanKeys.all, folderId] as const
};

export function createLibraryScansRuntime(client: QueryClient) {
	const scans = createQuery(() => ({
		queryKey: libraryScanKeys.all,
		queryFn: async () => ({}) as Record<string, LibraryScan>,
		staleTime: Infinity
	}));
	const upsert = (scan: LibraryScan) =>
		client.setQueryData<Record<string, LibraryScan>>(libraryScanKeys.all, (current) => ({
			...current,
			[scan.folderId]: scan
		}));
	const remove = (folderId: string) =>
		client.setQueryData<Record<string, LibraryScan>>(libraryScanKeys.all, (current) =>
			Object.fromEntries(Object.entries(current ?? {}).filter(([id]) => id !== folderId))
		);
	const clear = () => client.removeQueries({ queryKey: libraryScanKeys.all });
	return { scans, upsert, remove, clear };
}
