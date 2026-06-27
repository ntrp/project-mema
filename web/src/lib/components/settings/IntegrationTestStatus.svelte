<script lang="ts">
	import type { IntegrationTestResponse } from '$lib/settings/types';

	interface Props {
		enabled: boolean;
		result?: IntegrationTestResponse;
		testing?: boolean;
	}

	let { enabled, result, testing = false }: Props = $props();
</script>

<div class="status-stack" aria-live="polite">
	<span class:status-enabled={enabled} class:status-disabled={!enabled}>
		{enabled ? 'Enabled' : 'Disabled'}
	</span>
	{#if testing}
		<span class="test-status pending">Testing</span>
	{:else if result}
		<span class:test-ok={result.success} class:test-failed={!result.success}>
			{result.success ? 'Test OK' : 'Test failed'}
		</span>
		<span class="test-detail">{result.message} - {result.latencyMs} ms</span>
	{/if}
</div>
