import { createMutation, useQueryClient } from '@tanstack/svelte-query';
import {
	cancelDownloadActivity,
	clearReleaseBlocklist,
	deleteDownloadActivity,
	deleteReleaseBlocklistItem,
	manualImportDownloadActivity,
	retryDownloadActivity,
	type ManualImportRequest
} from './api';
import { activityKeys } from './queries.svelte';

export function createActivityMutations() {
	const client = useQueryClient();
	const refreshDownloads = () => client.invalidateQueries({ queryKey: activityKeys.downloads() });
	const refreshBlocklist = () => client.invalidateQueries({ queryKey: activityKeys.blocklist() });

	return {
		cancel: createMutation(() => ({
			mutationFn: (id: string) => cancelDownloadActivity(id),
			onSuccess: refreshDownloads
		})),
		retry: createMutation(() => ({
			mutationFn: (id: string) => retryDownloadActivity(id),
			onSuccess: refreshDownloads
		})),
		deleteDownload: createMutation(() => ({
			mutationFn: (id: string) => deleteDownloadActivity(id),
			onSuccess: refreshDownloads
		})),
		manualImport: createMutation(() => ({
			mutationFn: ({ id, request }: { id: string; request: ManualImportRequest }) =>
				manualImportDownloadActivity(id, request),
			onSuccess: refreshDownloads
		})),
		deleteBlocklistItem: createMutation(() => ({
			mutationFn: (id: string) => deleteReleaseBlocklistItem(id),
			onSuccess: refreshBlocklist
		})),
		clearBlocklist: createMutation(() => ({
			mutationFn: () => clearReleaseBlocklist(),
			onSuccess: refreshBlocklist
		}))
	};
}
