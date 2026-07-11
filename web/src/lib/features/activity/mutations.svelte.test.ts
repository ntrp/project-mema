import { beforeEach, describe, expect, it, vi } from 'vitest';

const mocks = vi.hoisted(() => ({
	client: { invalidateQueries: vi.fn() },
	createMutation: vi.fn((options: () => unknown) => options()),
	cancel: vi.fn(),
	clear: vi.fn(),
	deleteDownload: vi.fn(),
	deleteBlocklistItem: vi.fn(),
	manualImport: vi.fn()
}));

vi.mock('@tanstack/svelte-query', () => ({
	createMutation: mocks.createMutation,
	useQueryClient: () => mocks.client
}));
vi.mock('./api', () => ({
	cancelDownloadActivity: mocks.cancel,
	clearReleaseBlocklist: mocks.clear,
	deleteDownloadActivity: mocks.deleteDownload,
	deleteReleaseBlocklistItem: mocks.deleteBlocklistItem,
	manualImportDownloadActivity: mocks.manualImport
}));

import { createActivityMutations } from './mutations.svelte';
import { activityKeys } from './queries.svelte';

describe('activity mutations', () => {
	beforeEach(() => vi.clearAllMocks());

	it('forwards command arguments to the generated API wrappers', () => {
		const mutations = createActivityMutations() as unknown as MutationOptions;
		const request = { path: '/downloads/item.mkv' };

		mutations.cancel.mutationFn('one');
		mutations.deleteDownload.mutationFn('two');
		mutations.manualImport.mutationFn({ id: 'three', request });
		mutations.deleteBlocklistItem.mutationFn('four');
		mutations.clearBlocklist.mutationFn();

		expect(mocks.cancel).toHaveBeenCalledWith('one');
		expect(mocks.deleteDownload).toHaveBeenCalledWith('two');
		expect(mocks.manualImport).toHaveBeenCalledWith('three', request);
		expect(mocks.deleteBlocklistItem).toHaveBeenCalledWith('four');
		expect(mocks.clear).toHaveBeenCalledOnce();
	});

	it('invalidates only the collection affected by a successful command', () => {
		const mutations = createActivityMutations() as unknown as MutationOptions;

		mutations.cancel.onSuccess();
		mutations.deleteDownload.onSuccess();
		mutations.manualImport.onSuccess();
		mutations.deleteBlocklistItem.onSuccess();
		mutations.clearBlocklist.onSuccess();

		expect(mocks.client.invalidateQueries).toHaveBeenCalledTimes(5);
		expect(mocks.client.invalidateQueries).toHaveBeenNthCalledWith(1, {
			queryKey: activityKeys.downloads()
		});
		expect(mocks.client.invalidateQueries).toHaveBeenNthCalledWith(4, {
			queryKey: activityKeys.blocklist()
		});
	});
});

type Command<T> = { mutationFn: (value: T) => unknown; onSuccess: () => unknown };
type MutationOptions = {
	cancel: Command<string>;
	deleteDownload: Command<string>;
	manualImport: Command<{ id: string; request: { path: string } }>;
	deleteBlocklistItem: Command<string>;
	clearBlocklist: { mutationFn: () => unknown; onSuccess: () => unknown };
};
