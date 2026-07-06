<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import MediaAssemblyProvenanceTable from './MediaAssemblyProvenanceTable.svelte';
	import type { MediaComponentAssemblyRun } from '$lib/settings/types';
	import { fileName, statusTone } from './mediaComponentAssemblyView';

	interface Props {
		run?: MediaComponentAssemblyRun;
		canAssemble: boolean;
		assembleLabel: string;
		assembling: boolean;
		onAssemble: () => void;
	}

	let { run, canAssemble, assembleLabel, assembling, onAssemble }: Props = $props();
</script>

<div class="grid gap-3 rounded-md border p-4">
	<div class="flex flex-wrap items-start justify-between gap-3">
		<div class="grid gap-1">
			<h3 class="m-0 text-base font-semibold">Assembly</h3>
			<p class="m-0 text-sm text-muted-foreground">
				{#if run}
					{fileName(run.outputPath)}
				{:else}
					No mux run queued yet.
				{/if}
			</p>
		</div>
		<div class="flex items-center gap-2">
			{#if run}
				<Badge variant={statusTone(run.status)}>{run.status}</Badge>
			{/if}
			<Button size="sm" disabled={!canAssemble || assembling} onclick={onAssemble}>
				{assembling ? 'Queueing' : assembleLabel}
			</Button>
		</div>
	</div>
	{#if run?.errorMessage}
		<p class="m-0 text-sm text-destructive">{run.errorMessage}</p>
	{/if}
	{#if run?.status === 'running' || run?.status === 'queued'}
		<p class="m-0 text-sm text-muted-foreground">
			Mux job {run.status}; this panel updates from live events.
		</p>
	{/if}
	{#if run?.status === 'succeeded'}
		<MediaAssemblyProvenanceTable inputs={run.inputs} />
	{:else if run?.toolSummary}
		<p class="m-0 text-xs text-muted-foreground">{run.toolSummary}</p>
	{/if}
</div>
