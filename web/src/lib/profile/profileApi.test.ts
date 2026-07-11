import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({ GET: vi.fn(), PUT: vi.fn() }));
vi.mock('$lib/api/client', () => ({ client: clientMock }));

import { getProfile, updateProfile } from './profileApi';

describe('profile API', () => {
	beforeEach(() => {
		clientMock.GET.mockReset();
		clientMock.PUT.mockReset();
	});

	it('loads and updates the profile', async () => {
		clientMock.GET.mockResolvedValueOnce({ data: { id: 'user-1' } });
		await expect(getProfile()).resolves.toEqual({ id: 'user-1' });
		clientMock.PUT.mockResolvedValueOnce({ data: { id: 'user-1', displayName: 'Mema' } });
		await expect(updateProfile({ displayName: 'Mema' } as never)).resolves.toEqual({
			id: 'user-1',
			displayName: 'Mema'
		});
	});

	it('surfaces server and empty-response failures', async () => {
		clientMock.GET.mockResolvedValueOnce({ error: { message: 'profile failed' } });
		await expect(getProfile()).rejects.toThrow('profile failed');
		clientMock.PUT.mockResolvedValueOnce({ data: undefined });
		await expect(updateProfile({} as never)).rejects.toThrow(
			'Profile update did not return a result'
		);
	});
});
