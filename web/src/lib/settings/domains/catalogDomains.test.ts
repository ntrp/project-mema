import { beforeEach, describe, expect, it, vi } from 'vitest';

const client = vi.hoisted(() => ({
	POST: vi.fn(),
	PUT: vi.fn(),
	DELETE: vi.fn()
}));

vi.mock('$lib/api/client', () => ({ client }));

import { deleteCustomFormat, saveCustomFormat } from './customFormats';
import { deleteLanguage, saveLanguage } from './languages';
import { deleteUser, saveUser } from './users';

describe('focused settings catalog domains', () => {
	beforeEach(() => vi.clearAllMocks());

	it('creates and updates languages', async () => {
		client.POST.mockResolvedValueOnce({});
		client.PUT.mockResolvedValueOnce({});

		await saveLanguage({ code: ' en ', displayName: ' English ', aliasesText: 'eng' });
		await saveLanguage({
			originalCode: 'en',
			code: 'en',
			displayName: 'English',
			aliasesText: 'eng'
		});

		expect(client.POST).toHaveBeenCalledWith('/settings/languages', {
			body: { code: 'EN', displayName: 'English', aliases: ['eng'] }
		});
		expect(client.PUT).toHaveBeenCalledWith('/settings/languages/{code}', {
			params: { path: { code: 'en' } },
			body: { displayName: 'English', aliases: ['eng'] }
		});
	});

	it('creates and updates users', async () => {
		client.POST.mockResolvedValueOnce({});
		client.PUT.mockResolvedValueOnce({});

		await saveUser({ username: ' admin ', password: 'secret', role: 'admin' });
		await saveUser({ id: 'user-1', username: 'admin', password: '', role: 'user' });

		expect(client.POST).toHaveBeenCalledWith('/settings/users', {
			body: { username: 'admin', password: 'secret', role: 'admin' }
		});
		expect(client.PUT).toHaveBeenCalledWith('/settings/users/{id}', {
			params: { path: { id: 'user-1' } },
			body: { username: 'admin', role: 'user' }
		});
	});

	it('creates and updates custom formats', async () => {
		client.POST.mockResolvedValueOnce({});
		client.PUT.mockResolvedValueOnce({});
		const format = {
			name: ' WEB ',
			includeInRenameTemplate: true,
			includeSpecs: [],
			excludeSpecs: []
		};

		await saveCustomFormat(format);
		await saveCustomFormat({ ...format, id: 'format-1' });

		expect(client.POST).toHaveBeenCalledWith('/settings/custom-formats', {
			body: expect.objectContaining({ name: 'WEB' })
		});
		expect(client.PUT).toHaveBeenCalledWith(
			'/settings/custom-formats/{id}',
			expect.objectContaining({ params: { path: { id: 'format-1' } } })
		);
	});

	it.each([
		['language', deleteLanguage, '/settings/languages/{code}', 'en'],
		['user', deleteUser, '/settings/users/{id}', 'user-1'],
		['custom format', deleteCustomFormat, '/settings/custom-formats/{id}', 'format-1']
	] as const)('deletes a %s and surfaces failures', async (_label, remove, path, id) => {
		client.DELETE.mockResolvedValueOnce({});
		await expect(remove(id)).resolves.toBeUndefined();
		expect(client.DELETE).toHaveBeenCalledWith(path, { params: { path: expect.any(Object) } });

		client.DELETE.mockResolvedValueOnce({ error: { message: 'delete failed' } });
		await expect(remove(id)).rejects.toThrow('delete failed');
	});

	it.each([
		['language', () => saveLanguage({ code: 'en', displayName: 'English', aliasesText: '' })],
		['user', () => saveUser({ username: 'admin', password: 'secret', role: 'admin' })],
		[
			'custom format',
			() =>
				saveCustomFormat({
					name: 'WEB',
					includeInRenameTemplate: true,
					includeSpecs: [],
					excludeSpecs: []
				})
		]
	] as const)('surfaces %s save failures', async (_label, save) => {
		client.POST.mockResolvedValueOnce({ error: { message: 'save failed' } });
		await expect(save()).rejects.toThrow('save failed');
	});
});
