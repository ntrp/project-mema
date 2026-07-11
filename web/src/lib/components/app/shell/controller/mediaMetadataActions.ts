import { refreshMediaItemMetadata as refreshMediaItemMetadataRequest } from '$lib/settings/api';
import type { MediaItem } from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';
import type { RunCommandMutation } from '$lib/app/query/commandMutation.svelte';

interface MediaMetadataDeps {
	clearNotice: () => void;
	runMutation?: RunCommandMutation;
	upsertMediaItem: (_item: MediaItem) => void;
}

export function createMediaMetadataActions(state: AppShellState, deps: MediaMetadataDeps) {
	const clearNotice = deps.clearNotice;
	const runMutation = deps.runMutation ?? ((command) => command());

	async function refreshMediaMetadata(item: MediaItem) {
		state.refreshingMetadataItemId = item.id;
		clearNotice();

		try {
			const updated = await runMutation(() => refreshMediaItemMetadataRequest(item.id));
			deps.upsertMediaItem(updated);
			state.message = 'Media metadata refreshed';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not refresh media metadata');
		} finally {
			state.refreshingMetadataItemId = undefined;
		}
	}

	return { refreshMediaMetadata };
}
