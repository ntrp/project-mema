<script lang="ts">
	import CopyIcon from '@lucide/svelte/icons/copy';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { ReleaseCandidate } from '$lib/settings/types';

	interface Props {
		release: ReleaseCandidate;
		copiedReleaseId?: string;
		onCopy: (_release: ReleaseCandidate) => void;
	}

	let { release, copiedReleaseId, onCopy }: Props = $props();
</script>

<div class="group/title relative w-full min-w-0">
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<span {...props} class="block min-w-0 truncate pr-7">{release.title}</span>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content class="max-w-160">{release.title}</Tooltip.Content>
	</Tooltip.Root>
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="ghost"
					size="icon-sm"
					aria-label="Copy release title"
					class="pointer-events-none absolute top-1/2 right-0 -translate-y-1/2 opacity-0 transition-opacity group-hover/title:pointer-events-auto group-hover/title:opacity-100 group-focus-within/title:pointer-events-auto group-focus-within/title:opacity-100"
					onclick={() => onCopy(release)}
				>
					<CopyIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content class="max-w-160">
			{copiedReleaseId === release.id ? 'Copied title' : release.title}
		</Tooltip.Content>
	</Tooltip.Root>
</div>
