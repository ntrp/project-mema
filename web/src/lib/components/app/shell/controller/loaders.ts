import {
	getMediaCollection as getMediaCollectionRequest,
	getMediaMetadataDetails as getMediaMetadataDetailsRequest,
	getPersonDetails as getPersonDetailsRequest,
	loadSettings as loadSettingsRequest
} from '$lib/settings/api';
import type { MediaType, MetadataProviderType } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

export function createLoadActions(state: AppShellState) {
	async function loadSettings() {
		try {
			const settings = await loadSettingsRequest();
			state.downloadClients = settings.downloadClients;
			state.indexers = settings.indexers;
			state.indexerSearch = settings.indexerSearch;
			state.metadataProviders = settings.metadataProviders;
			state.subtitleProviders = settings.subtitleProviders;
			state.metadataCache = settings.metadataCache;
			state.libraryFolders = settings.libraryFolders;
			state.pathMappings = settings.pathMappings;
			state.mediaProfiles = settings.mediaProfiles;
			state.customFormats = settings.customFormats;
			state.users = settings.users;
			state.tags = settings.tags;
			state.languages = settings.languages;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load settings');
		}
	}

	async function loadMetadataDetail() {
		if (
			!state.route.metadataProvider ||
			!state.route.metadataType ||
			!state.route.metadataExternalId
		) {
			return;
		}
		state.loadingMetadataDetail = true;
		try {
			state.metadataDetail = await getMediaMetadataDetailsRequest(
				state.route.metadataProvider as MetadataProviderType,
				state.route.metadataType as MediaType,
				state.route.metadataExternalId
			);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load media details');
		} finally {
			state.loadingMetadataDetail = false;
		}
	}

	async function loadMediaCollection() {
		if (!state.route.collectionProvider || !state.route.collectionId) {
			return;
		}
		state.loadingMediaCollection = true;
		try {
			state.mediaCollection = await getMediaCollectionRequest(
				state.route.collectionProvider as MetadataProviderType,
				state.route.collectionId
			);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load media collection');
		} finally {
			state.loadingMediaCollection = false;
		}
	}

	async function loadPersonDetail() {
		if (!state.route.personProvider || !state.route.personId) {
			return;
		}
		state.loadingPersonDetail = true;
		try {
			state.personDetail = await getPersonDetailsRequest(
				state.route.personProvider as MetadataProviderType,
				state.route.personId
			);
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load person details');
		} finally {
			state.loadingPersonDetail = false;
		}
	}

	return {
		loadSettings,
		loadMetadataDetail,
		loadPersonDetail,
		loadMediaCollection
	};
}
