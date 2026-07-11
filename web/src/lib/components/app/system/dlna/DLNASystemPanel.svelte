<script lang="ts">
	import RefreshCwIcon from '@lucide/svelte/icons/refresh-cw';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import DLNADiagnosticsTables from '$lib/components/settings/dlna/DLNADiagnosticsTables.svelte';
	import { createDLNAStatusQuery } from '$lib/components/settings/dlna/dlnaStatus.svelte';
	import type { DLNAStatus } from '$lib/settings/types';

	const statusQuery = createDLNAStatusQuery();
	const status = $derived<DLNAStatus | undefined>(statusQuery.data?.status);
	const loading = $derived(statusQuery.isPending || statusQuery.isFetching);
	const errorMessage = $derived(queryErrorMessage(statusQuery.error));
	const statusCells = $derived([
		{ label: 'State', value: status?.running ? 'Running' : 'Stopped' },
		{ label: 'SSDP', value: status?.lastSsdpEvent ?? 'None' },
		{ label: 'Last SOAP', value: status?.lastSoapAction ?? 'None' },
		{ label: 'Last error', value: status?.lastError ?? 'None' }
	]);

	function queryErrorMessage(error: unknown) {
		if (!error) return '';
		return error instanceof Error ? error.message : 'Could not load DLNA status';
	}
</script>

<Card.Root aria-label="DLNA runtime diagnostics">
	<Card.Header class="border-b border-border">
		<Card.Title>DLNA</Card.Title>
		<Card.Action>
			<Button
				type="button"
				variant="secondary"
				size="sm"
				disabled={loading}
				onclick={() => void statusQuery.refetch()}
			>
				<RefreshCwIcon class={loading ? 'animate-spin' : ''} />
				Refresh
			</Button>
		</Card.Action>
	</Card.Header>
	<Card.Content class="grid gap-5 pt-5">
		{#if errorMessage}
			<p class="text-sm font-medium text-destructive">{errorMessage}</p>
		{/if}
		{#if loading && !status}
			<p class="text-sm text-muted-foreground">Loading DLNA status...</p>
		{:else if status}
			<div class="grid gap-3 sm:grid-cols-4">
				{#each statusCells as cell (cell.label)}
					<div class="grid gap-1 rounded-md border border-border p-3">
						<span class="text-xs font-medium uppercase text-muted-foreground">{cell.label}</span>
						<span class="break-words text-sm font-medium text-foreground">{cell.value}</span>
					</div>
				{/each}
			</div>
			<DLNADiagnosticsTables {status} />
		{/if}
	</Card.Content>
</Card.Root>
