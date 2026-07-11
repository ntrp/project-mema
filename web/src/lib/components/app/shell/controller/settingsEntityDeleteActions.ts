import {
	deleteCustomFormat as deleteCustomFormatRequest,
	deleteLanguage as deleteLanguageRequest,
	deleteMediaProfile as deleteMediaProfileRequest,
	deletePathMapping as deletePathMappingRequest,
	deleteUser as deleteUserRequest
} from '$lib/settings/api';
import { deleteTag as deleteTagRequest } from '$lib/components/settings/tags/api';
import {
	emptyCustomFormatForm,
	emptyLanguageForm,
	emptyMediaProfileForm,
	emptyUserForm
} from '$lib/settings/forms';
import { emptyTagForm, errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface SettingsEntityDeleteDeps {
	clearNotice: () => void;
}

export function createSettingsEntityDeleteActions(
	state: AppShellState,
	deps: SettingsEntityDeleteDeps
) {
	const clearNotice = deps.clearNotice;

	async function deletePathMapping(id: string) {
		state.deletingPathMappingId = id;
		clearNotice();

		try {
			await deletePathMappingRequest(id);
			state.pathMappings = state.pathMappings.filter((mapping) => mapping.id !== id);
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
			await deleteUserRequest(id);
			if (state.userForm.id === id) {
				state.userForm = emptyUserForm();
			}
			state.users = state.users.filter((user) => user.id !== id);
			state.message = 'User deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete user');
		}
	}

	async function deleteTag(id: string) {
		state.deletingTagId = id;
		clearNotice();

		try {
			await deleteTagRequest(id);
			if (state.tagForm.id === id) {
				state.tagForm = emptyTagForm();
			}
			state.tags = state.tags.filter((tag) => tag.id !== id);
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
			await deleteLanguageRequest(code);
			if (state.languageForm.originalCode === code) {
				state.languageForm = emptyLanguageForm();
			}
			state.languages = state.languages.filter((language) => language.code !== code);
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
			await deleteMediaProfileRequest(id);
			if (state.mediaProfileForm.id === id) {
				state.mediaProfileForm = emptyMediaProfileForm();
			}
			state.mediaProfiles = state.mediaProfiles.filter((profile) => profile.id !== id);
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
			await deleteCustomFormatRequest(id);
			if (state.customFormatForm.id === id) {
				state.customFormatForm = emptyCustomFormatForm();
			}
			state.customFormats = state.customFormats.filter((format) => format.id !== id);
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
