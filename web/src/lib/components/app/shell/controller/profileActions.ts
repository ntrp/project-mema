import { updateProfile as updateProfileRequest } from '$lib/profile/profileApi';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';
import type { UserProfile } from '$lib/profile/types';
import type { UserProfileUpdateRequest } from '$lib/profile/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface ProfileDeps {
	clearNotice: () => void;
	loadProfile: () => Promise<UserProfile>;
	runMutation: RunCommandMutation;
}

export function createProfileActions(state: AppShellState, deps: ProfileDeps) {
	const clearNotice = deps.clearNotice;

	async function loadProfile() {
		state.loadingProfile = true;
		state.profileErrorMessage = '';
		try {
			state.profile = await deps.loadProfile();
		} catch (error) {
			state.profileErrorMessage = errorMessageFrom(error, 'Could not load profile');
		} finally {
			state.loadingProfile = false;
		}
	}

	async function saveProfile(request: UserProfileUpdateRequest) {
		state.savingProfile = true;
		state.profileErrorMessage = '';
		clearNotice();
		try {
			const profile = await deps.runMutation(() => updateProfileRequest(request));
			state.profile = profile;
			state.currentUser = {
				id: profile.id,
				username: profile.username,
				displayName: profile.displayName || undefined,
				pictureUrl: profile.pictureUrl || undefined,
				role: profile.role
			};
			state.message = 'Profile saved';
		} catch (error) {
			state.profileErrorMessage = errorMessageFrom(error, 'Could not save profile');
		} finally {
			state.savingProfile = false;
		}
	}

	return { loadProfile, saveProfile };
}
