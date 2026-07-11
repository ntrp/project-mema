import { describe, expect, it, vi } from 'vitest';

vi.mock('@tanstack/svelte-query', () => ({
	createQuery: vi.fn((options: () => { queryFn: () => unknown }) => ({ data: options().queryFn() }))
}));

import { createLibraryScansRuntime, libraryScanKeys } from './libraryScans.svelte';

describe('library scan runtime', () => {
	it('upserts, removes, and clears folder scan cache data', () => {
		const client = { setQueryData: vi.fn(), removeQueries: vi.fn() };
		const runtime = createLibraryScansRuntime(client as never);
		runtime.upsert({ folderId: 'folder-1' } as never);
		expect(client.setQueryData).toHaveBeenCalledWith(libraryScanKeys.all, expect.any(Function));
		runtime.remove('folder-1');
		runtime.clear();
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: libraryScanKeys.all });
	});
});
