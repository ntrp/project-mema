<script lang="ts">
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import type { MediaFileDetailRow } from '$lib/components/app/media/files/mediaFileDetails';

	interface Props {
		track: MediaFileDetailRow;
		canManage: boolean;
		onRequestDelete: (_event: Event, _track: MediaFileDetailRow) => void;
	}

	let { track, canManage, onRequestDelete }: Props = $props();
</script>

{#if track.deleteRequest}
	<Tooltip.Root>
		<Tooltip.Trigger>
			{#snippet child({ props })}
				<Button
					{...props}
					type="button"
					variant="destructive"
					size="icon-sm"
					aria-label="Delete embedded track"
					disabled={!canManage}
					onclick={(event) => onRequestDelete(event, track)}
					onkeydown={(event) => event.stopPropagation()}
				>
					<TrashIcon aria-hidden="true" />
				</Button>
			{/snippet}
		</Tooltip.Trigger>
		<Tooltip.Content>Delete embedded track</Tooltip.Content>
	</Tooltip.Root>
{/if}
