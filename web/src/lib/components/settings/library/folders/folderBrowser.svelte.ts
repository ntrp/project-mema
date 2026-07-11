import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import { createLibraryFolderOption, listLibraryFolderOptions } from '$lib/settings/api';

export const folderBrowserKeys = {
	all: ['settings', 'folder-browser'] as const,
	path: (path?: string) => [...folderBrowserKeys.all, path ?? ''] as const
};

export function createFolderBrowserResources() {
	const client = useQueryClient();
	return {
		load: (path?: string) =>
			client.fetchQuery({
				queryKey: folderBrowserKeys.path(path),
				queryFn: () => listLibraryFolderOptions(path),
				staleTime: 30_000
			}),
		create: createMutation(() => ({
			mutationFn: ({ parentPath, name }: { parentPath: string; name: string }) =>
				createLibraryFolderOption(parentPath, name),
			onSuccess: (_created, variables) =>
				client.invalidateQueries({ queryKey: folderBrowserKeys.path(variables.parentPath) })
		}))
	};
}
