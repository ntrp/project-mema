import {
	deleteMediaItemSubtitle as deleteMediaItemSubtitleRequest,
	enqueueMediaSubtitleSearch as enqueueMediaSubtitleSearchRequest
} from '$lib/settings/api';
import type { MediaItem, SubtitleSearchRequest } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface MediaSubtitleDeps {
	clearNotice: () => void;
}

export function createMediaSubtitleActions(state: AppShellState, deps: MediaSubtitleDeps) {
	const clearNotice = deps.clearNotice;

	async function searchMediaSubtitle(item: MediaItem, request: SubtitleSearchRequest = {}) {
		clearNotice();
		try {
			const job = await enqueueMediaSubtitleSearchRequest(item.id, request);
			state.message = `${job.message} (#${job.jobId})`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue subtitle search');
		}
	}

	async function deleteMediaSubtitle(item: MediaItem, subtitleId: string) {
		clearNotice();
		try {
			const updated = await deleteMediaItemSubtitleRequest(item.id, subtitleId);
			state.mediaItems = [
				updated,
				...state.mediaItems.filter((mediaItem) => mediaItem.id !== updated.id)
			];
			state.message = 'Subtitle deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete subtitle');
		}
	}

	return {
		searchMediaSubtitle,
		deleteMediaSubtitle
	};
}
