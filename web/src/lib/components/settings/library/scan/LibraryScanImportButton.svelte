<script lang="ts">
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Tooltip from '$lib/components/ui/tooltip';

	interface Props {
		canImport: boolean;
		importDisabledReason: string;
		loading: boolean;
		importing: boolean;
		duplicateRemovalCount: number;
		onImport: () => void;
	}

	let {
		canImport,
		importDisabledReason,
		loading,
		importing,
		duplicateRemovalCount,
		onImport
	}: Props = $props();
	let tooltipOpen = $state(false);
	const disabled = $derived(!canImport || loading || importing);
	const tooltip = $derived(
		importing ? 'Import in progress' : loading ? 'Library scan is loading' : importDisabledReason
	);
</script>

<Tooltip.Root bind:open={tooltipOpen}>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<span {...props} class="inline-flex">
				{#if duplicateRemovalCount > 0}
					<ConfirmActionButton
						label="Import selected"
						title="Remove files"
						description={`Import selected rows and remove ${duplicateRemovalCount} file${duplicateRemovalCount === 1 ? '' : 's'}?`}
						confirmLabel="Import selected"
						confirmingLabel="Importing"
						variant="default"
						class="whitespace-nowrap"
						{disabled}
						onConfirm={onImport}
					>
						{#if importing}<InlineSpinner label="Importing" />{:else}Import Selected{/if}
					</ConfirmActionButton>
				{:else}
					<Button type="button" class="whitespace-nowrap" {disabled} onclick={onImport}>
						{#if importing}<InlineSpinner label="Importing" />{:else}Import Selected{/if}
					</Button>
				{/if}
			</span>
		{/snippet}
	</Tooltip.Trigger>
	{#if disabled && tooltipOpen}<Tooltip.Content>{tooltip}</Tooltip.Content>{/if}
</Tooltip.Root>
