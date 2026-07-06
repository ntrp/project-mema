<script lang="ts">
	import CheckIcon from '@lucide/svelte/icons/check';
	import XIcon from '@lucide/svelte/icons/x';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type {
		MediaComponentCompatibilityReviewState,
		MediaComponentSource
	} from '$lib/settings/types';
	import {
		fileName,
		sourceDisplayName,
		sourceSummary,
		statusTone
	} from './mediaComponentAssemblyView';

	interface Props {
		sources: MediaComponentSource[];
		canManage: boolean;
		reviewingDecisionId?: string;
		onReview: (
			_source: MediaComponentSource,
			_decisionId: string,
			_reviewState: MediaComponentCompatibilityReviewState
		) => void;
	}

	let { sources, canManage, reviewingDecisionId, onReview }: Props = $props();
</script>

<div class="grid gap-2.5">
	{#each sources as source (source.id)}
		<div class="grid gap-3 rounded-md border p-4">
			<div class="flex flex-wrap items-start justify-between gap-3">
				<div class="grid min-w-0 gap-1">
					<strong class="truncate">{sourceDisplayName(source)}</strong>
					<span class="text-sm text-muted-foreground">{fileName(source.retainedPath)}</span>
				</div>
				<div class="flex flex-wrap gap-1.5">
					<Badge>{source.sourceRole}</Badge>
					<Badge variant={statusTone(source.retentionState)}>{source.retentionState}</Badge>
				</div>
			</div>
			<p class="m-0 text-sm text-muted-foreground">{sourceSummary(source)}</p>
			{#if source.artifacts?.length}
				<div class="flex flex-wrap gap-1.5">
					{#each source.artifacts as artifact (artifact.id)}
						<Badge variant={statusTone(artifact.status)}>
							{artifact.streamType}
							{artifact.language ?? artifact.streamId}: {artifact.status}
						</Badge>
					{/each}
				</div>
			{/if}
			{#if source.compatibility?.length}
				<div class="grid gap-2">
					{#each source.compatibility as decision (decision.id)}
						<div
							class="flex flex-wrap items-center justify-between gap-2 rounded border bg-muted/30 px-3 py-2"
						>
							<span class="text-sm">
								{decision.confidenceState} confidence · {decision.automationState} · {decision.reviewState}
							</span>
							{#if canManage && decision.reviewState === 'pending'}
								<div class="flex gap-1">
									<Tooltip.Root>
										<Tooltip.Trigger>
											<Button
												size="icon"
												variant="ghost"
												disabled={reviewingDecisionId === decision.id}
												aria-label="Approve compatibility"
												onclick={() => onReview(source, decision.id, 'approved')}
											>
												<CheckIcon aria-hidden="true" />
											</Button>
										</Tooltip.Trigger>
										<Tooltip.Content>Approve compatibility</Tooltip.Content>
									</Tooltip.Root>
									<Tooltip.Root>
										<Tooltip.Trigger>
											<Button
												size="icon"
												variant="ghost"
												disabled={reviewingDecisionId === decision.id}
												aria-label="Reject compatibility"
												onclick={() => onReview(source, decision.id, 'rejected')}
											>
												<XIcon aria-hidden="true" />
											</Button>
										</Tooltip.Trigger>
										<Tooltip.Content>Reject compatibility</Tooltip.Content>
									</Tooltip.Root>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/each}
</div>
