import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	customFormatFormFromFormat,
	downloadClientFormFromClient,
	indexerFormFromIndexer,
	languageFormFromLanguage,
	mediaProfileFormFromProfile,
	userFormFromUser
} from '$lib/settings/forms';
import type {
	CustomFormat,
	DownloadClient,
	Indexer,
	Language,
	ManagedUser,
	MediaProfile,
	Tag
} from '$lib/settings/types';
import type { AppShellState } from './state.svelte';

export function createSettingsEditActions(state: AppShellState) {
	function editDownloadClient(client: DownloadClient) {
		state.downloadForm = downloadClientFormFromClient(client);
		state.activeSettingsSection = 'download-clients';
		void goto(resolve('/settings/download-clients'));
	}

	function editIndexer(indexer: Indexer) {
		state.indexerForm = indexerFormFromIndexer(indexer);
		state.activeSettingsSection = 'indexers';
		void goto(resolve('/settings/indexers'));
	}

	function editUser(user: ManagedUser) {
		state.userForm = userFormFromUser(user);
		state.activeSettingsSection = 'users';
		void goto(resolve('/settings/users'));
	}

	function editTag(tag: Tag) {
		state.tagForm = { id: tag.id, name: tag.name };
		state.activeSettingsSection = 'tags';
		void goto(resolve('/settings/tags'));
	}

	function editLanguage(language: Language) {
		state.languageForm = languageFormFromLanguage(language);
		state.activeSettingsSection = 'languages';
		void goto(resolve('/settings/languages'));
	}

	function editMediaProfile(profile: MediaProfile) {
		state.mediaProfileForm = mediaProfileFormFromProfile(profile);
		state.activeSettingsSection = 'profiles';
		void goto(resolve('/settings/profiles'));
	}

	function editCustomFormat(format: CustomFormat) {
		state.customFormatForm = customFormatFormFromFormat(format);
		state.activeSettingsSection = 'custom-formats';
		void goto(resolve('/settings/custom-formats'));
	}

	return {
		editDownloadClient,
		editIndexer,
		editUser,
		editTag,
		editLanguage,
		editMediaProfile,
		editCustomFormat
	};
}
