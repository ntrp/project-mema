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

<div class="flex max-w-96 items-center gap-1.5">
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<span {...props} class="block min-w-0 flex-1 truncate">{release.title}</span>
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
					onclick={() => onCopy(release)}
				>
					<CopyIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>{copiedReleaseId === release.id ? 'Copied title' : 'Copy title'}</Tooltip.Content>
	</Tooltip.Root>
</div>
