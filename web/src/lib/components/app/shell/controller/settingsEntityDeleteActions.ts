import {
	deleteMediaProfile as deleteMediaProfileRequest,
	deletePathMapping as deletePathMappingRequest
} from '$lib/settings/api';
import { deleteTag as deleteTagRequest } from '$lib/components/settings/tags/api';
import { deleteCustomFormat as deleteCustomFormatRequest } from '$lib/settings/domains/customFormats';
import { deleteLanguage as deleteLanguageRequest } from '$lib/settings/domains/languages';
import { deleteUser as deleteUserRequest } from '$lib/settings/domains/users';
import {
	emptyCustomFormatForm,
	emptyLanguageForm,
	emptyMediaProfileForm,
	emptyUserForm
} from '$lib/settings/forms';
import { emptyTagForm, errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';

interface SettingsEntityDeleteDeps {
	clearNotice: () => void;
	runMutation?: RunCommandMutation;
	removeLanguage?: (_code: string) => void;
	removeTag?: (_id: string) => void;
	removeUser?: (_id: string) => void;
	removePathMapping?: (_id: string) => void;
	removeMediaProfile?: (_id: string) => void;
	removeCustomFormat?: (_id: string) => void;
}

export function createSettingsEntityDeleteActions(
	state: AppShellState,
	deps: SettingsEntityDeleteDeps
) {
	const clearNotice = deps.clearNotice;
	const runMutation = deps.runMutation ?? ((command) => command());

	async function deletePathMapping(id: string) {
		state.deletingPathMappingId = id;
		clearNotice();

		try {
			await runMutation(() => deletePathMappingRequest(id));
			deps.removePathMapping?.(id);
			state.message = 'Path mapping deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete path mapping');
		} finally {
			state.deletingPathMappingId = undefined;
		}
	}

	async function deleteUser(id: string) {
		clearNotice();

		try {
			await runMutation(() => deleteUserRequest(id));
			if (state.userForm.id === id) {
				state.userForm = emptyUserForm();
			}
			deps.removeUser?.(id);
			state.message = 'User deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete user');
		}
	}

	async function deleteTag(id: string) {
		state.deletingTagId = id;
		clearNotice();

		try {
			await runMutation(() => deleteTagRequest(id));
			if (state.tagForm.id === id) {
				state.tagForm = emptyTagForm();
			}
			deps.removeTag?.(id);
			state.message = 'Tag deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete tag');
		} finally {
			state.deletingTagId = undefined;
		}
	}

	async function deleteLanguage(code: string) {
		state.deletingLanguageCode = code;
		clearNotice();

		try {
			await runMutation(() => deleteLanguageRequest(code));
			if (state.languageForm.originalCode === code) {
				state.languageForm = emptyLanguageForm();
			}
			deps.removeLanguage?.(code);
			state.message = 'Language deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete language');
		} finally {
			state.deletingLanguageCode = undefined;
		}
	}

	async function deleteMediaProfile(id: string) {
		state.deletingMediaProfileId = id;
		clearNotice();

		try {
			await runMutation(() => deleteMediaProfileRequest(id));
			if (state.mediaProfileForm.id === id) {
				state.mediaProfileForm = emptyMediaProfileForm();
			}
			deps.removeMediaProfile?.(id);
			state.message = 'Profile deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete profile');
		} finally {
			state.deletingMediaProfileId = undefined;
		}
	}

	async function deleteCustomFormat(id: string) {
		state.deletingCustomFormatId = id;
		clearNotice();

		try {
			await runMutation(() => deleteCustomFormatRequest(id));
			if (state.customFormatForm.id === id) {
				state.customFormatForm = emptyCustomFormatForm();
			}
			deps.removeCustomFormat?.(id);
			state.message = 'Custom format deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete custom format');
		} finally {
			state.deletingCustomFormatId = undefined;
		}
	}

	return {
		deletePathMapping,
		deleteUser,
		deleteTag,
		deleteLanguage,
		deleteMediaProfile,
		deleteCustomFormat
	};
}
