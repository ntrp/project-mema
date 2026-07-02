import {
	getMediaCollection as getMediaCollectionRequest,
	getMediaMetadataDetails as getMediaMetadataDetailsRequest,
	listDownloadActivity as listDownloadActivityRequest,
	listMediaItems as listMediaItemsRequest,
	listMediaRequests as listMediaRequestsRequest,
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
			state.metadataCache = settings.metadataCache;
			state.libraryFolders = settings.libraryFolders;
			state.pathMappings = settings.pathMappings;
			state.mediaProfiles = settings.mediaProfiles;
			state.customFormats = settings.customFormats;
			state.users = settings.users;
			state.tags = settings.tags;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load settings');
		}
	}

	async function loadLibrary() {
		await Promise.all([loadMediaItems(), loadMediaRequests(), loadDownloadActivity()]);
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

	async function loadMediaItems() {
		state.loadingMediaItems = true;
		try {
			state.mediaItems = await listMediaItemsRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load media items');
		} finally {
			state.loadingMediaItems = false;
		}
	}

	async function loadMediaRequests() {
		try {
			state.mediaRequests = await listMediaRequestsRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load media requests');
		}
	}

	async function loadDownloadActivity() {
		state.loadingActivity = true;
		try {
			state.activities = await listDownloadActivityRequest();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not load download activity');
		} finally {
			state.loadingActivity = false;
		}
	}

	return {
		loadSettings,
		loadLibrary,
		loadMetadataDetail,
		loadMediaCollection,
		loadMediaItems,
		loadMediaRequests,
		loadDownloadActivity
	};
}
