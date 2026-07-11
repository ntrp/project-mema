import type { QueryClient } from '@tanstack/svelte-query';
import { describe, expect, it, vi } from 'vitest';
import { createSettingsCatalogCache } from './cache';
import { settingsCatalogKeys } from './queries.svelte';

describe('settings catalog cache', () => {
	it('removes catalog entries and supports refresh and logout clearing', () => {
		const client = { setQueryData: vi.fn(), invalidateQueries: vi.fn(), removeQueries: vi.fn() };
		const cache = createSettingsCatalogCache(client as unknown as QueryClient);
		cache.removeLanguage('de');
		cache.removeTag('tag-1');
		cache.removeUser('user-1');
		cache.removeDownloadClient('download-1');
		cache.removeIndexer('indexer-1');
		cache.removeMetadataProvider('metadata-1');
		cache.removeSubtitleProvider('subtitle-1');
		cache.removeLibraryFolder('folder-1');
		cache.removePathMapping('mapping-1');
		cache.removeMediaProfile('profile-1');
		cache.removeCustomFormat('format-1');
		cache.upsertLibraryFolder({ id: 'folder-2' } as never);
		cache.upsertPathMapping({ id: 'mapping-2' } as never);
		for (const [, update] of client.setQueryData.mock.calls) {
			expect(update([{ code: 'de', id: 'tag-1' }])).toEqual(expect.any(Array));
		}
		cache.refresh();
		expect(client.invalidateQueries).toHaveBeenCalledWith({ queryKey: settingsCatalogKeys.all });
		cache.clear();
		expect(client.removeQueries).toHaveBeenCalledWith({ queryKey: settingsCatalogKeys.all });
	});
});
