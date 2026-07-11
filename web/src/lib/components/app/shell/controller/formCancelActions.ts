import {
	emptyCustomFormatForm,
	emptyDownloadClientForm,
	emptyIndexerForm,
	emptyLanguageForm,
	emptyMediaProfileForm,
	emptyUserForm
} from '$lib/settings/forms';
import { emptyTagForm } from './helpers';
import type { AppShellState } from './state.svelte';

export function createFormCancelActions(state: AppShellState) {
	return {
		cancelDownloadClient: () => (state.downloadForm = emptyDownloadClientForm()),
		cancelIndexer: () => (state.indexerForm = emptyIndexerForm()),
		cancelMediaProfile: () => (state.mediaProfileForm = emptyMediaProfileForm()),
		cancelCustomFormat: () => (state.customFormatForm = emptyCustomFormatForm()),
		cancelTag: () => (state.tagForm = emptyTagForm()),
		cancelLanguage: () => (state.languageForm = emptyLanguageForm()),
		cancelUser: () => (state.userForm = emptyUserForm())
	};
}
