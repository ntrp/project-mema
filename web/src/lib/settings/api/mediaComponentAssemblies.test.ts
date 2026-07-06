import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({
	POST: vi.fn()
}));

vi.mock('$lib/api/client', () => ({ client: clientMock }));

import { enqueueMediaComponentAssembly } from './mediaComponentAssemblies';

describe('media component assembly API helpers', () => {
	beforeEach(() => {
		clientMock.POST.mockReset();
	});

	it('queues an assembly job', async () => {
		clientMock.POST.mockResolvedValueOnce({
			data: { jobId: 42, message: 'queued', run: { id: 'run-1' } }
		});

		await expect(
			enqueueMediaComponentAssembly('media-1', {
				baseSourceId: 'base-1',
				artifactIds: ['artifact-1']
			})
		).resolves.toEqual({ jobId: 42, message: 'queued', run: { id: 'run-1' } });
		expect(clientMock.POST).toHaveBeenLastCalledWith('/media/items/{id}/assemblies', {
			params: { path: { id: 'media-1' } },
			body: { baseSourceId: 'base-1', artifactIds: ['artifact-1'] }
		});
	});
});
