import { beforeEach, describe, expect, it, vi } from 'vitest';

const clientMock = vi.hoisted(() => ({ GET: vi.fn(), POST: vi.fn() }));

vi.mock('$lib/api/client', () => ({ client: clientMock }));

import { currentSession, currentSessionAuthenticated, login, logout } from './api';

describe('session API', () => {
	beforeEach(() => {
		clientMock.GET.mockReset();
		clientMock.POST.mockReset();
	});

	it('reads the current session and authenticated state', async () => {
		clientMock.GET.mockResolvedValueOnce({ data: { authenticated: true } });
		await expect(currentSession()).resolves.toEqual({ authenticated: true });
		expect(clientMock.GET).toHaveBeenCalledWith('/auth/session');

		clientMock.GET.mockResolvedValueOnce({ data: undefined });
		await expect(currentSessionAuthenticated()).resolves.toBe(false);
	});

	it('logs in authenticated users and rejects failed authentication', async () => {
		clientMock.POST.mockResolvedValueOnce({ data: { authenticated: true } });
		await expect(login('admin', 'secret')).resolves.toEqual({ authenticated: true });
		expect(clientMock.POST).toHaveBeenCalledWith('/auth/login', {
			body: { username: 'admin', password: 'secret' }
		});

		clientMock.POST.mockResolvedValueOnce({ error: { message: 'Invalid credentials' } });
		await expect(login('admin', 'bad')).rejects.toThrow('Invalid credentials');
		clientMock.POST.mockResolvedValueOnce({ data: { authenticated: false } });
		await expect(login('admin', 'bad')).rejects.toThrow('Login failed');
	});

	it('logs out and surfaces server errors', async () => {
		clientMock.POST.mockResolvedValueOnce({});
		await expect(logout()).resolves.toBeUndefined();
		clientMock.POST.mockResolvedValueOnce({ error: { message: 'logout failed' } });
		await expect(logout()).rejects.toThrow('logout failed');
	});
});
