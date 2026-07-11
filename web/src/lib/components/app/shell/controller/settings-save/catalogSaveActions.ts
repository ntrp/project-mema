import {
	saveCustomFormat as saveCustomFormatRequest,
	saveLanguage as saveLanguageRequest,
	saveMediaProfile as saveMediaProfileRequest
} from '$lib/settings/api';
import { saveTag as saveTagFocusedRequest } from '$lib/components/settings/tags/api';
import {
	emptyCustomFormatForm,
	emptyLanguageForm,
	emptyMediaProfileForm
} from '$lib/settings/forms';
import type { CustomFormatForm as CustomFormatFormValue } from '$lib/settings/types';
import { emptyTagForm, errorMessageFrom } from '../helpers';
import type { SettingsSaveContext } from './types';

export function createCatalogSaveActions({
	state,
	clearNotice,
	loadSettings,
	loadMediaItems,
	mediaItems
}: SettingsSaveContext) {
	async function saveTag(event: SubmitEvent) {
		event.preventDefault();
		state.savingTag = true;
		clearNotice();

		try {
			await saveTagFocusedRequest(state.tagForm);
			state.tagForm = emptyTagForm();
			state.message = 'Tag saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save tag');
		} finally {
			state.savingTag = false;
		}
	}

	async function saveLanguage(event: SubmitEvent) {
		event.preventDefault();
		state.savingLanguage = true;
		clearNotice();

		try {
			await saveLanguageRequest(state.languageForm);
			state.languageForm = emptyLanguageForm();
			state.message = 'Language saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save language');
		} finally {
			state.savingLanguage = false;
		}
	}

	async function saveMediaProfile(event: SubmitEvent) {
		event.preventDefault();
		state.savingMediaProfile = true;
		clearNotice();
		const profileId = state.mediaProfileForm.id;
		const refreshAffectedMedia = Boolean(
			profileId && mediaItems?.().some((item) => item.qualityProfileId === profileId)
		);

		try {
			await saveMediaProfileRequest(state.mediaProfileForm);
			state.mediaProfileForm = emptyMediaProfileForm();
			state.message = 'Profile saved';
			await loadSettings();
			if (refreshAffectedMedia && loadMediaItems) {
				await loadMediaItems();
			}
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save profile');
		} finally {
			state.savingMediaProfile = false;
		}
	}

	async function saveCustomFormat(event: SubmitEvent) {
		event.preventDefault();
		state.savingCustomFormat = true;
		clearNotice();

		try {
			await saveCustomFormatRequest(state.customFormatForm);
			state.customFormatForm = emptyCustomFormatForm();
			state.message = 'Custom format saved';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not save custom format');
		} finally {
			state.savingCustomFormat = false;
		}
	}

	async function importCustomFormat(format: CustomFormatFormValue) {
		state.savingCustomFormat = true;
		clearNotice();

		try {
			await saveCustomFormatRequest(format);
			state.message = 'Custom format imported';
			await loadSettings();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not import custom format');
			throw error;
		} finally {
			state.savingCustomFormat = false;
		}
	}

	return {
		saveTag,
		saveLanguage,
		saveMediaProfile,
		saveCustomFormat,
		importCustomFormat
	};
}
