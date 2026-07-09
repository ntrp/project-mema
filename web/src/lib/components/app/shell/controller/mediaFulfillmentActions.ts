import { enqueueMediaFulfillmentAction } from '$lib/settings/api';
import type { MediaFulfillmentActionRequest, MediaItem } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface MediaFulfillmentDeps {
	clearNotice: () => void;
	loadMediaItems: () => Promise<void>;
}

export function createMediaFulfillmentActions(state: AppShellState, deps: MediaFulfillmentDeps) {
	async function enqueueMediaFulfillment(item: MediaItem, request: MediaFulfillmentActionRequest) {
		deps.clearNotice();
		try {
			const job = await enqueueMediaFulfillmentAction(item.id, request);
			state.message = `${job.message} (#${job.jobId})`;
			await deps.loadMediaItems();
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not enqueue fulfillment job');
		}
	}

	return { enqueueMediaFulfillment };
}
