import { client } from '$lib/api/client';
import type { SessionResponse } from './types';

export async function currentSession(): Promise<SessionResponse | undefined> {
	const { data } = await client.GET('/auth/session');
	return data;
}

export async function currentSessionAuthenticated() {
	const data = await currentSession();
	return Boolean(data?.authenticated);
}

export async function login(username: string, password: string) {
	const { data, error } = await client.POST('/auth/login', {
		body: { username, password }
	});

	if (error || !data?.authenticated) {
		throw new Error(error?.message ?? 'Login failed');
	}
	return data;
}

export async function logout() {
	const { error } = await client.POST('/auth/logout');

	if (error) throw new Error(error.message);
}
