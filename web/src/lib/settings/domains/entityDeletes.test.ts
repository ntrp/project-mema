import { beforeEach, describe, expect, it, vi } from 'vitest';

const remove = vi.hoisted(() => vi.fn());
vi.mock('$lib/api/client', () => ({ client: { DELETE: remove } }));

import {
	deleteDownloadClient,
	deleteIndexer,
	deleteMediaProfile,
	deleteMetadataProvider,
	deleteSubtitleProvider
} from './entityDeletes';

const commands = [
	[deleteDownloadClient, '/settings/download-clients/{id}'],
	[deleteIndexer, '/settings/indexers/{id}'],
	[deleteMetadataProvider, '/settings/metadata-providers/{id}'],
	[deleteSubtitleProvider, '/settings/subtitle-providers/{id}'],
	[deleteMediaProfile, '/settings/profiles/{id}']
] as const;

describe('settings entity deletes', () => {
	beforeEach(() => vi.clearAllMocks());

	it.each(commands)('deletes through %s', async (command, path) => {
		remove.mockResolvedValueOnce({});
		await expect(command('entity-1')).resolves.toBeUndefined();
		expect(remove).toHaveBeenCalledWith(path, { params: { path: { id: 'entity-1' } } });
	});

	it.each(commands)('surfaces errors from %s', async (command) => {
		remove.mockResolvedValueOnce({ error: { message: 'delete failed' } });
		await expect(command('entity-1')).rejects.toThrow('delete failed');
	});
});
