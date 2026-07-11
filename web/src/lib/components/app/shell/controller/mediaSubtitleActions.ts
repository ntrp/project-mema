import {
	deleteMediaItemSubtitle as deleteMediaItemSubtitleRequest,
	updateMediaItemSubtitle as updateMediaItemSubtitleRequest
} from '$lib/features/library/filesApi';
import {
	grabMediaSubtitle as grabMediaSubtitleRequest,
	enqueueMediaSubtitleSearch as enqueueMediaSubtitleSearchRequest
} from '$lib/features/releases/api';
import type {
	GrabSubtitleRequest,
	MediaItem,
	MediaItemSubtitleSelectionRequest,
	SubtitleSearchRequest
} from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';

interface MediaSubtitleDeps {
	clearNotice: () => void;
	runMutation?: RunCommandMutation;
	upsertMediaItem: (_item: MediaItem) => void;
}

export function createMediaSubtitleActions(state: AppShellState, deps: MediaSubtitleDeps) {
	const clearNotice = deps.clearNotice;
	const runMutation = deps.runMutation ?? ((command) => command());

	async function searchMediaSubtitle(item: MediaItem, request: SubtitleSearchRequest = {}) {
		clearNotice();
		try {
			const job = await runMutation(() => enqueueMediaSubtitleSearchRequest(item.id, request));
			state.message = `${job.message} (#${job.jobId})`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue subtitle search');
		}
	}

	async function deleteMediaSubtitle(item: MediaItem, subtitleId: string) {
		clearNotice();
		try {
			const updated = await runMutation(() => deleteMediaItemSubtitleRequest(item.id, subtitleId));
			deps.upsertMediaItem(updated);
			state.message = 'Subtitle deleted';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not delete subtitle');
		}
	}

	async function updateMediaSubtitle(
		item: MediaItem,
		subtitleId: string,
		request: MediaItemSubtitleSelectionRequest
	) {
		clearNotice();
		try {
			const updated = await runMutation(() =>
				updateMediaItemSubtitleRequest(item.id, subtitleId, request)
			);
			deps.upsertMediaItem(updated);
			state.message = 'Subtitle updated';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not update subtitle');
		}
	}

	async function grabMediaSubtitle(item: MediaItem, request: GrabSubtitleRequest) {
		clearNotice();
		try {
			const updated = await runMutation(() => grabMediaSubtitleRequest(item.id, request));
			deps.upsertMediaItem(updated);
			state.message = 'Subtitle grabbed';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not grab subtitle');
			throw error;
		}
	}

	return {
		searchMediaSubtitle,
		deleteMediaSubtitle,
		updateMediaSubtitle,
		grabMediaSubtitle
	};
}
