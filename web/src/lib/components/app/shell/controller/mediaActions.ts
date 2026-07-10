import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	approveMediaRequest as approveMediaRequestRequest,
	createMediaItem as createMediaItemRequest,
	createMediaRequest as createMediaRequestRequest,
	deleteMediaItem as deleteMediaItemRequest,
	deleteMediaItemFile as deleteMediaItemFileRequest,
	deleteMediaItemFileTrack as deleteMediaItemFileTrackRequest,
	enqueueMediaAutomaticSearch as enqueueMediaAutomaticSearchRequest,
	rescanMediaItemFiles as rescanMediaItemFilesRequest,
	updateMediaItem as updateMediaItemRequest
} from '$lib/settings/api';
import type { MediaActionSelection } from '$lib/components/app/media/actions/mediaActionTypes';
import type {
	MediaItem,
	MediaFileTrackDeleteRequest,
	MediaItemUpdateRequest,
	MediaRequest,
	MediaRequestApproveRequest,
	MediaSearchResult
} from '$lib/settings/types';
import { candidateKey, errorMessageFrom, omitResult } from './helpers';
import { mediaUpdateMessage, optimisticMediaItem } from './mediaOptimisticUpdate';
import type { AppShellState } from './state.svelte';

interface MediaDeps {
	clearNotice: () => void;
	loadMediaItems: () => Promise<void>;
	loadSettings: () => Promise<void>;
}

export function createMediaActions(state: AppShellState, deps: MediaDeps) {
	const clearNotice = deps.clearNotice;
	const loadMediaItems = deps.loadMediaItems;
	const loadSettings = deps.loadSettings;
	function addMedia(candidate: MediaSearchResult) {
		state.activeMediaCandidate = candidate;
		clearNotice();
	}

	function closeMediaAction() {
		if (!state.savingMediaAction) {
			state.activeMediaCandidate = undefined;
		}
	}

	async function confirmMediaAction(selection: MediaActionSelection) {
		const candidate = state.activeMediaCandidate;
		if (!candidate) {
			return;
		}
		state.addingKey = candidateKey(candidate);
		state.savingMediaAction = true;
		clearNotice();

		try {
			if (state.isAdmin) {
				if (!selection.qualityProfileId || !selection.libraryFolderId) {
					throw new Error('Quality profile and library folder are required');
				}
				await createMediaItemRequest({
					title: candidate.title,
					type: candidate.type,
					year: candidate.year,
					monitored: selection.monitorMode !== 'none',
					monitorMode: selection.monitorMode,
					seriesType: candidate.type === 'serie' ? selection.seriesType : undefined,
					minimumAvailability: selection.minimumAvailability,
					startSearch: selection.startSearch,
					externalProvider: candidate.externalProvider,
					externalId: candidate.externalId,
					overview: candidate.overview,
					posterPath: candidate.posterPath,
					qualityProfileId: selection.qualityProfileId,
					libraryFolderId: selection.libraryFolderId,
					tags: selection.tags
				});
				await loadMediaItems();
				await loadSettings();
				state.message =
					selection.monitorMode === 'none'
						? 'Media item added to library'
						: selection.monitorMode === 'collection'
							? 'Media collection added to monitored'
							: 'Media item added to monitored';
				state.activeHomeSection = candidate.type === 'movie' ? 'movies' : 'series';
				state.activeMediaCandidate = undefined;
				void goto(resolve(candidate.type === 'movie' ? '/movies' : '/series'));
				return;
			}

			const request = await createMediaRequestRequest({
				title: candidate.title,
				type: candidate.type,
				year: candidate.year,
				externalProvider: candidate.externalProvider,
				externalId: candidate.externalId,
				overview: candidate.overview,
				posterPath: candidate.posterPath
			});
			state.mediaRequests = [
				request,
				...state.mediaRequests.filter((item) => item.id !== request.id)
			];
			state.message = 'Media request created';
			state.activeHomeSection = 'requests';
			state.activeMediaCandidate = undefined;
			void goto(resolve('/requests'));
		} catch (error) {
			state.errorMessage = errorMessageFrom(
				error,
				state.isAdmin ? 'Could not add media item' : 'Could not create media request'
			);
		} finally {
			state.addingKey = undefined;
			state.savingMediaAction = false;
		}
	}

	async function autoSearchMedia(item: MediaItem) {
		state.searchingItemId = item.id;
		clearNotice();
		try {
			const job = await enqueueMediaAutomaticSearchRequest(item.id);
			state.message = `${job.message} (#${job.jobId})`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue automatic search');
		} finally {
			state.searchingItemId = undefined;
		}
	}

	async function rescanMediaFiles(item: MediaItem) {
		state.scanningMediaItemId = item.id;
		clearNotice();

		try {
			const updated = await rescanMediaItemFilesRequest(item.id);
			state.mediaItems = [
				updated,
				...state.mediaItems.filter((mediaItem) => mediaItem.id !== updated.id)
			];
			state.message = `File scan completed: ${updated.filePaths.length} media, ${updated.metadataFilePaths.length} metadata`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not rescan media folder');
		} finally {
			state.scanningMediaItemId = undefined;
		}
	}

	async function deleteMediaFile(item: MediaItem, path: string) {
		clearNotice();
		try {
			const updated = await deleteMediaItemFileRequest(item.id, path);
			state.mediaItems = [
				updated,
				...state.mediaItems.filter((mediaItem) => mediaItem.id !== updated.id)
			];
			state.message = 'Media file deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete media file');
		}
	}

	async function deleteMediaFileTrack(item: MediaItem, request: MediaFileTrackDeleteRequest) {
		clearNotice();
		try {
			const updated = await deleteMediaItemFileTrackRequest(item.id, request);
			state.mediaItems = [
				updated,
				...state.mediaItems.filter((mediaItem) => mediaItem.id !== updated.id)
			];
			state.message = 'Embedded track deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete embedded track');
		}
	}

	async function saveMediaItemOptions(item: MediaItem, request: MediaItemUpdateRequest) {
		state.savingMediaItemOptionsId = item.id;
		clearNotice();
		const previous = state.mediaItems.find((mediaItem) => mediaItem.id === item.id) ?? item;
		const optimistic = optimisticMediaItem(previous, request);
		const message = mediaUpdateMessage(previous, optimistic, request);
		state.mediaItems = state.mediaItems.map((mediaItem) =>
			mediaItem.id === optimistic.id ? optimistic : mediaItem
		);
		state.message = message;

		try {
			const updated = await updateMediaItemRequest(item.id, request);
			state.mediaItems = state.mediaItems.map((mediaItem) =>
				mediaItem.id === updated.id ? updated : mediaItem
			);
			state.message = message;
		} catch (error) {
			state.mediaItems = state.mediaItems.map((mediaItem) =>
				mediaItem.id === previous.id ? previous : mediaItem
			);
			state.message = '';
			state.errorMessage = errorMessageFrom(error, 'Could not save media settings');
		} finally {
			state.savingMediaItemOptionsId = undefined;
		}
	}

	function deleteMediaItem(item: MediaItem) {
		clearNotice();
		state.mediaDeleteCandidate = item;
	}

	function closeMediaDelete() {
		if (!state.deletingMediaItemId) {
			state.mediaDeleteCandidate = undefined;
		}
	}

	async function confirmMediaDelete(keepFiles: boolean) {
		const item = state.mediaDeleteCandidate;
		if (!item) {
			return;
		}
		await removeMediaItem(item, keepFiles);
	}

	async function removeMediaItem(item: MediaItem, keepFiles: boolean) {
		state.deletingMediaItemId = item.id;
		clearNotice();

		try {
			await deleteMediaItemRequest(item.id, { keepFiles });
			state.mediaItems = state.mediaItems.filter((mediaItem) => mediaItem.id !== item.id);
			state.releaseResults = omitResult(state.releaseResults, item.id);
			state.activities = state.activities.filter((activity) => activity.mediaItemId !== item.id);
			state.mediaDeleteCandidate = undefined;
			state.message = keepFiles ? 'Media item removed; files kept' : 'Media item and files removed';
			if (state.selectedMediaItemId === item.id) {
				state.selectedMediaItemId = undefined;
				void goto(resolve(item.type === 'movie' ? '/movies' : '/series'));
			}
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not remove media item');
		} finally {
			state.deletingMediaItemId = undefined;
		}
	}

	async function approveMediaRequest(request: MediaRequest, approval: MediaRequestApproveRequest) {
		state.approvingRequestId = request.id;
		clearNotice();

		try {
			const result = await approveMediaRequestRequest(request.id, approval);
			state.mediaRequests = state.mediaRequests.map((item) =>
				item.id === result.request.id ? result.request : item
			);
			state.mediaItems = [
				result.mediaItem,
				...state.mediaItems.filter((item) => item.id !== result.mediaItem.id)
			];
			state.message = 'Media request approved and added to monitored';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not approve media request');
		} finally {
			state.approvingRequestId = undefined;
		}
	}

	return {
		addMedia,
		closeMediaAction,
		confirmMediaAction,
		autoSearchMedia,
		rescanMediaFiles,
		deleteMediaFile,
		deleteMediaFileTrack,
		saveMediaItemOptions,
		deleteMediaItem,
		closeMediaDelete,
		confirmMediaDelete,
		approveMediaRequest
	};
}
