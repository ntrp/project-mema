<script lang="ts">
	import ConfirmActionButton from '$lib/components/shared/ConfirmActionButton.svelte';
	import InlineSpinner from '$lib/components/shared/InlineSpinner.svelte';
	import { Button } from '$lib/components/ui/button';

	interface Props {
		canImport: boolean;
		loading: boolean;
		importing: boolean;
		duplicateRemovalCount: number;
		onImport: () => void;
	}

	let { canImport, loading, importing, duplicateRemovalCount, onImport }: Props = $props();
	const disabled = $derived(!canImport || loading || importing);
</script>

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
		{#if importing}
			<InlineSpinner label="Importing" />
		{:else}
			<span>Import Selected</span>
		{/if}
	</ConfirmActionButton>
{:else}
	<Button type="button" class="whitespace-nowrap" {disabled} onclick={onImport}>
		{#if importing}
			<InlineSpinner label="Importing" />
		{:else}
			<span>Import Selected</span>
		{/if}
	</Button>
{/if}
