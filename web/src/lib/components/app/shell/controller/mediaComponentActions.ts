import { enqueueMediaComponentAssembly } from '$lib/settings/api/mediaComponentAssemblies';
import { reviewMediaComponentCompatibility } from '$lib/settings/api/mediaComponentSources';
import type {
	MediaComponentCompatibilityReviewState,
	MediaComponentSource,
	MediaItem
} from '$lib/settings/types';
import { errorMessageFrom } from './helpers';
import type { AppShellState } from './state.svelte';

interface MediaComponentDeps {
	clearNotice: () => void;
	loadMediaItems: () => Promise<void>;
}

export function createMediaComponentActions(state: AppShellState, deps: MediaComponentDeps) {
	async function reviewComponentCompatibility(
		item: MediaItem,
		source: MediaComponentSource,
		decisionId: string,
		reviewState: MediaComponentCompatibilityReviewState
	) {
		state.reviewingComponentDecisionId = decisionId;
		deps.clearNotice();

		try {
			await reviewMediaComponentCompatibility(item.id, source.id, decisionId, {
				reviewState,
				reason: 'Reviewed from media detail'
			});
			await deps.loadMediaItems();
			state.message =
				reviewState === 'approved'
					? 'Component compatibility approved'
					: 'Component compatibility rejected';
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not review component compatibility');
		} finally {
			state.reviewingComponentDecisionId = undefined;
		}
	}

	async function assembleMediaComponents(
		item: MediaItem,
		baseSourceId: string,
		artifactIds: string[]
	) {
		state.assemblingMediaItemId = item.id;
		deps.clearNotice();

		try {
			const result = await enqueueMediaComponentAssembly(item.id, { baseSourceId, artifactIds });
			await deps.loadMediaItems();
			state.message = `${result.message} (#${result.jobId})`;
		} catch (error) {
			state.errorMessage = errorMessageFrom(error, 'Could not queue component assembly');
		} finally {
			state.assemblingMediaItemId = undefined;
		}
	}

	return {
		reviewComponentCompatibility,
		assembleMediaComponents
	};
}
