<script lang="ts">
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
		if (indexer.nextCheckAt) return `Next check ${formatDate(indexer.nextCheckAt)}`;
		if (indexer.lastError) return indexer.lastError;
		if (indexer.lastQueryAt) return `Last query ${formatDate(indexer.lastQueryAt)}`;
		return 'No query yet';
	}

	function formatDate(value: string) {
		return new Intl.DateTimeFormat(undefined, {
			month: 'short',
			day: '2-digit',
			hour: '2-digit',
			minute: '2-digit'
		}).format(new Date(value));
	}
</script>

<div class="status-stack" aria-live="polite">
	<span
		class:status-enabled={indexer.enabled && indexer.healthStatus === 'healthy'}
		class:pending={indexer.enabled && indexer.healthStatus === 'temporary_disabled'}
		class:status-disabled={!indexer.enabled || indexer.healthStatus === 'disabled'}
	>
		{healthLabel(indexer)}
	</span>
	<span class="test-detail">{detail(indexer)}</span>
	{#if indexer.lastStatusCode}
		<span class="test-detail">HTTP {indexer.lastStatusCode}</span>
	{/if}
</div>
