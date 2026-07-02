<script lang="ts">
	import DownloadIcon from '@lucide/svelte/icons/download';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaItem, ReleaseCandidate, ReleaseOverrideDetails } from '$lib/settings/types';
	import ReleaseOverrideIcon from './ReleaseOverrideIcon.svelte';

	interface Props {
		item: MediaItem;
		release: ReleaseCandidate;
		grabbingKey?: string;
		onGrab: (
			_item: MediaItem,
			_release: ReleaseCandidate,
			_overrideMatch?: boolean,
			_details?: ReleaseOverrideDetails
		) => void;
	}

	let { item, release, grabbingKey, onGrab }: Props = $props();

	const currentReleaseKey = $derived(`${item.id}:${release.id}`);
	const grabbing = $derived(grabbingKey === currentReleaseKey);
	const showGrab = $derived(release.match.severity !== 'error');

	const grabTooltip = $derived.by(() => {
		if (grabbing) {
			return 'Queueing release';
		}
		return 'Grab release';
	});

	const overrideGrabTooltip = $derived(grabbing ? 'Queueing release' : 'Grab with override');
</script>

<div class="flex justify-end gap-1">
	{#if showGrab}
		<Tooltip.Root>
			<Tooltip.Trigger>
				{#snippet child({ props })}
					<Button
						{...props}
						type="button"
						size="icon-sm"
						class="bg-emerald-600 text-white hover:bg-emerald-700"
						aria-label="Grab release"
						disabled={grabbing}
						onclick={() => onGrab(item, release)}
					>
						<DownloadIcon aria-hidden="true" />
					</Button>
				{/snippet}
			</Tooltip.Trigger>
			<Tooltip.Content>
				{grabTooltip}
			</Tooltip.Content>
		</Tooltip.Root>
	{/if}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					size="icon-sm"
					class="bg-amber-400 text-amber-950 hover:bg-amber-500"
					aria-label="Grab with override"
					disabled={grabbing}
					onclick={() => onGrab(item, release, true)}
				>
					<ReleaseOverrideIcon />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>
			{overrideGrabTooltip}
		</Tooltip.Content>
	</Tooltip.Root>
</div>
