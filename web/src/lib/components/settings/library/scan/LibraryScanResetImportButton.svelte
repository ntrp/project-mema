<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import XIcon from '@lucide/svelte/icons/x';
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import type { LibraryScanItem } from '$lib/settings/types';

	interface Props {
		item: LibraryScanItem;
		resetting?: boolean;
		onResetImport: (_item: LibraryScanItem) => void | Promise<void>;
	}

	let { item, resetting = false, onResetImport }: Props = $props();
</script>

<ConfirmActionButton
	label={`Reset import for ${item.fileName}`}
	title="Reset imported file"
	description={`Remove the media and metadata for "${item.fileName}" without touching files on disk?`}
	confirmLabel="Reset import"
	confirmingLabel="Resetting"
	size="icon-sm"
	variant="outline"
	class="mt-0.5 shrink-0"
	disabled={resetting}
	tooltip="Reset import"
	onConfirm={() => onResetImport(item)}
>
	{#if resetting}
		<RefreshCwIcon aria-hidden="true" />
	{:else}
		<XIcon aria-hidden="true" />
	{/if}
</ConfirmActionButton>
