import { client } from '$lib/api/client';
import type { UserProfile, UserProfileUpdateRequest } from '$lib/settings/types';

export async function getProfile(): Promise<UserProfile> {
	const { data, error } = await client.GET('/profile');
	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Profile request did not return a result');
	}
	return data;
}

export async function updateProfile(request: UserProfileUpdateRequest): Promise<UserProfile> {
	const { data, error } = await client.PUT('/profile', { body: request });
	if (error) {
		throw new Error(error.message);
	}
	if (!data) {
		throw new Error('Profile update did not return a result');
	}
	return data;
}
