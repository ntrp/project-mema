<script lang="ts">
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import { formatCompactDateTime } from '$lib/settings/dateFormat';
	import type { Indexer, IntegrationTestResponse } from '$lib/settings/types';

	interface Props {
		indexer: Indexer;
		result?: IntegrationTestResponse;
		checking?: boolean;
	}

	let { indexer, result, checking = false }: Props = $props();

	function healthLabel(indexer: Indexer) {
		if (!indexer.enabled) return 'Disabled';
		if (indexer.healthStatus === 'temporary_disabled') return 'Backing off';
		if (indexer.healthStatus === 'disabled') return 'Disabled until check';
		return 'Healthy';
	}

	function detail(indexer: Indexer) {
		if (checking) return 'Checking now';
		if (result) return `${result.message} - ${result.latencyMs} ms`;
		if (indexer.nextCheckAt) return `Next check ${formatCompactDateTime(indexer.nextCheckAt)}`;
		if (indexer.lastError) return indexer.lastError;
		if (indexer.lastQueryAt) return `Last query ${formatCompactDateTime(indexer.lastQueryAt)}`;
		return 'No query yet';
	}

	function healthTone(indexer: Indexer): 'success' | 'muted' | 'pending' {
		if (!indexer.enabled || indexer.healthStatus === 'disabled') return 'muted';
		if (indexer.healthStatus === 'temporary_disabled') return 'pending';
		return 'success';
	}
</script>

<div class="grid min-w-30 gap-1" aria-live="polite">
	<StatusPill tone={healthTone(indexer)}>{healthLabel(indexer)}</StatusPill>
	<span class="max-w-55 text-xs text-muted-foreground">{detail(indexer)}</span>
	{#if indexer.lastStatusCode}
		<span class="max-w-55 text-xs text-muted-foreground">HTTP {indexer.lastStatusCode}</span>
	{/if}
</div>
