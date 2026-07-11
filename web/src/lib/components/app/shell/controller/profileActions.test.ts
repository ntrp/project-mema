import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { AppShellState } from './state.svelte';

const profileApi = vi.hoisted(() => ({ get: vi.fn(), update: vi.fn() }));
vi.mock('$lib/profile/profileApi', () => ({
	getProfile: profileApi.get,
	updateProfile: profileApi.update
}));
import { createProfileActions } from './profileActions';

describe('profile actions', () => {
	let state: AppShellState;

	beforeEach(() => {
		Object.values(profileApi).forEach((request) => request.mockReset());
		state = {
			loadingProfile: false,
			savingProfile: false,
			profileErrorMessage: '',
			message: ''
		} as AppShellState;
	});

	it('loads and saves a profile while synchronizing the current user', async () => {
		const loaded = { id: 'user-1', username: 'ada' };
		const saved = { ...loaded, displayName: '', pictureUrl: '', role: 'admin' };
		profileApi.get.mockResolvedValue(loaded);
		profileApi.update.mockResolvedValue(saved);
		const actions = createProfileActions(state, dependencies());
		await actions.loadProfile();
		expect(state.profile).toEqual(loaded);
		expect(state.loadingProfile).toBe(false);
		await actions.saveProfile({} as never);
		expect(state.profile).toEqual(saved);
		expect(state.currentUser).toEqual({
			id: 'user-1',
			username: 'ada',
			displayName: undefined,
			pictureUrl: undefined,
			role: 'admin'
		});
		expect(state.message).toBe('Profile saved');
		expect(state.savingProfile).toBe(false);
	});

	it('records load and save errors and resets busy state', async () => {
		profileApi.get.mockRejectedValue(new Error('load failed'));
		profileApi.update.mockRejectedValue(new Error('save failed'));
		const actions = createProfileActions(state, dependencies());
		await actions.loadProfile();
		expect(state.profileErrorMessage).toBe('load failed');
		expect(state.loadingProfile).toBe(false);
		await actions.saveProfile({} as never);
		expect(state.profileErrorMessage).toBe('save failed');
		expect(state.savingProfile).toBe(false);
	});

	function dependencies() {
		return {
			clearNotice: vi.fn(),
			loadProfile: profileApi.get,
			runMutation: <T>(command: () => Promise<T>) => command()
		};
	}
});
