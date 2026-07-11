import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({ POST: vi.fn(), PUT: vi.fn(), DELETE: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client: clientMock }));

import { deleteTag, saveTag } from './api';

describe('tag API', () => {
	beforeEach(() => {
		clientMock.POST.mockReset();
		clientMock.PUT.mockReset();
		clientMock.DELETE.mockReset();
	});

	it('creates and updates normalized tags', async () => {
		clientMock.POST.mockResolvedValueOnce({});
		await saveTag({ name: ' Action ' });
		expect(clientMock.POST).toHaveBeenCalledWith('/settings/tags', { body: { name: 'Action' } });

		clientMock.PUT.mockResolvedValueOnce({});
		await saveTag({ id: 'tag-1', name: ' Kids ' });
		expect(clientMock.PUT).toHaveBeenCalledWith('/settings/tags/{id}', {
			params: { path: { id: 'tag-1' } },
			body: { name: 'Kids' }
		});
	});

	it('deletes tags and surfaces errors', async () => {
		clientMock.DELETE.mockResolvedValueOnce({});
		await expect(deleteTag('tag-1')).resolves.toBeUndefined();
		clientMock.DELETE.mockResolvedValueOnce({ error: { message: 'delete failed' } });
		await expect(deleteTag('tag-1')).rejects.toThrow('delete failed');
	});
});
