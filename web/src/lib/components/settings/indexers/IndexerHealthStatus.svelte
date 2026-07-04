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
		if (indexer.healthStatus === 'temporary_disabled') return 'Temp blocked';
		if (indexer.healthStatus === 'disabled') return 'Permanently blocked';
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

	function healthTone(indexer: Indexer): 'success' | 'muted' | 'pending' | 'error' {
		if (indexer.healthStatus === 'disabled') return 'error';
		if (!indexer.enabled) return 'muted';
		if (indexer.healthStatus === 'temporary_disabled') return 'pending';
		return 'success';
	}

	function healthClass(indexer: Indexer) {
		if (indexer.healthStatus === 'disabled') return 'bg-destructive/10 text-destructive';
		if (!indexer.enabled) return 'bg-muted text-muted-foreground';
		if (indexer.healthStatus === 'temporary_disabled') {
			return 'bg-yellow-500/10 text-yellow-700 dark:text-yellow-300';
		}
		return 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-300';
	}
</script>

<div class="flex w-full min-w-0 items-center gap-2 whitespace-nowrap" aria-live="polite">
	<StatusPill tone={healthTone(indexer)} class={healthClass(indexer)}>
		{healthLabel(indexer)}
	</StatusPill>
	<span class="min-w-0 flex-1 truncate text-xs text-muted-foreground">{detail(indexer)}</span>
	{#if indexer.lastStatusCode}
		<span class="shrink-0 text-xs text-muted-foreground">HTTP {indexer.lastStatusCode}</span>
	{/if}
</div>
