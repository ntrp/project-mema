<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import MediaAssemblyStatusPanel from './MediaAssemblyStatusPanel.svelte';
	import MediaComponentSourceList from './MediaComponentSourceList.svelte';
	import { mediaComponentAssemblyView } from './mediaComponentAssemblyView';
	import type {
		MediaComponentCompatibilityReviewState,
		MediaComponentSource,
		MediaItem
	} from '$lib/settings/types';

	interface Props {
		item: MediaItem;
		canManage: boolean;
		assemblingMediaItemId?: string;
		reviewingComponentDecisionId?: string;
		onAssemble: (_item: MediaItem, _baseSourceId: string, _artifactIds: string[]) => void;
		onReview: (
			_item: MediaItem,
			_source: MediaComponentSource,
			_decisionId: string,
			_reviewState: MediaComponentCompatibilityReviewState
		) => void;
	}

	let {
		item,
		canManage,
		assemblingMediaItemId,
		reviewingComponentDecisionId,
		onAssemble,
		onReview
	}: Props = $props();

	const view = $derived(mediaComponentAssemblyView(item));

	function assemble() {
		if (!view.baseSource || view.allowedArtifacts.length === 0) return;
		onAssemble(
			item,
			view.baseSource.id,
			view.allowedArtifacts.map((artifact) => artifact.id)
		);
	}
</script>

{#if view.retainedSources.length > 0 || view.latestRun}
	<section aria-labelledby="media-component-assembly-title">
		<div class="flex flex-wrap items-center justify-between gap-2">
			<h2 id="media-component-assembly-title" class="m-0 text-3xl font-semibold text-foreground">
				Components
			</h2>
			<div class="flex flex-wrap gap-1.5">
				<Badge variant="secondary">{view.retainedSources.length} retained</Badge>
				<Badge variant={view.blockedCount > 0 ? 'destructive' : 'secondary'}>
					{view.blockedCount} blocked
				</Badge>
			</div>
		</div>
		<div class="grid gap-3">
			<MediaAssemblyStatusPanel
				run={view.latestRun}
				canAssemble={canManage && view.canAssemble}
				assembleLabel={view.assembleLabel}
				assembling={assemblingMediaItemId === item.id}
				onAssemble={assemble}
			/>
			<MediaComponentSourceList
				sources={view.retainedSources}
				{canManage}
				reviewingDecisionId={reviewingComponentDecisionId}
				onReview={(source, decisionId, reviewState) =>
					onReview(item, source, decisionId, reviewState)}
			/>
		</div>
	</section>
{/if}
