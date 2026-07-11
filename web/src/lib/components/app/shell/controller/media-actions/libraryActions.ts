import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import {
	rescanMediaItemFiles as rescanMediaItemFilesRequest,
	updateMediaItem as updateMediaItemRequest
} from '$lib/features/library/commands';
import {
	deleteMediaItemFile as deleteMediaItemFileRequest,
	deleteMediaItemFileTrack as deleteMediaItemFileTrackRequest
} from '$lib/features/library/filesApi';
import { deleteMediaItem as deleteMediaItemRequest } from '$lib/features/library/commands';
import { enqueueMediaAutomaticSearch as enqueueMediaAutomaticSearchRequest } from '$lib/features/releases/api';
import type {
	MediaFileTrackDeleteRequest,
	MediaItem,
	MediaItemUpdateRequest
} from '$lib/settings/types';
import { errorMessageFrom } from '../helpers';
import { mediaUpdateMessage, optimisticMediaItem } from '../mediaOptimisticUpdate';
import type { AppShellState } from '../state.svelte';
import type { MediaDeps } from './types';

export function createMediaLibraryActions(state: AppShellState, deps: MediaDeps) {
	async function autoSearchMedia(item: MediaItem) {
		state.searchingItemId = item.id;
		deps.clearNotice();
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
		deps.clearNotice();
		try {
			const updated = await rescanMediaItemFilesRequest(item.id);
			deps.upsertMediaItem(updated);
			state.message = `File scan completed: ${updated.filePaths.length} media, ${updated.metadataFilePaths.length} metadata`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not rescan media folder');
		} finally {
			state.scanningMediaItemId = undefined;
		}
	}

	async function deleteMediaFile(item: MediaItem, path: string) {
		deps.clearNotice();
		try {
			deps.upsertMediaItem(await deleteMediaItemFileRequest(item.id, path));
			state.message = 'Media file deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete media file');
		}
	}

	async function deleteMediaFileTrack(item: MediaItem, request: MediaFileTrackDeleteRequest) {
		deps.clearNotice();
		try {
			deps.upsertMediaItem(await deleteMediaItemFileTrackRequest(item.id, request));
			state.message = 'Embedded track deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete embedded track');
		}
	}

	async function saveMediaItemOptions(item: MediaItem, request: MediaItemUpdateRequest) {
		state.savingMediaItemOptionsId = item.id;
		deps.clearNotice();
		const previous = deps.mediaItems().find((entry) => entry.id === item.id) ?? item;
		const optimistic = optimisticMediaItem(previous, request);
		const message = mediaUpdateMessage(previous, optimistic, request);
		deps.mapMediaItems((entry) => (entry.id === optimistic.id ? optimistic : entry));
		state.message = message;
		try {
			const updated = await updateMediaItemRequest(item.id, request);
			deps.mapMediaItems((entry) => (entry.id === updated.id ? updated : entry));
			state.message = message;
		} catch (error) {
			deps.mapMediaItems((entry) => (entry.id === previous.id ? previous : entry));
			state.message = '';
			state.errorMessage = errorMessageFrom(error, 'Could not save media settings');
		} finally {
			state.savingMediaItemOptionsId = undefined;
		}
	}

	function deleteMediaItem(item: MediaItem) {
		deps.clearNotice();
		state.mediaDeleteCandidate = item;
	}

	function closeMediaDelete() {
		if (!state.deletingMediaItemId) state.mediaDeleteCandidate = undefined;
	}

	async function confirmMediaDelete(keepFiles: boolean) {
		if (state.mediaDeleteCandidate) await removeMediaItem(state.mediaDeleteCandidate, keepFiles);
	}

	async function removeMediaItem(item: MediaItem, keepFiles: boolean) {
		state.deletingMediaItemId = item.id;
		deps.clearNotice();
		try {
			await deleteMediaItemRequest(item.id, { keepFiles });
			deps.removeMediaItem(item.id);
			deps.removeReleaseResults(item.id);
			deps.removeActivityForMedia(item.id);
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

	return {
		autoSearchMedia,
		rescanMediaFiles,
		deleteMediaFile,
		deleteMediaFileTrack,
		saveMediaItemOptions,
		deleteMediaItem,
		closeMediaDelete,
		confirmMediaDelete
	};
}
