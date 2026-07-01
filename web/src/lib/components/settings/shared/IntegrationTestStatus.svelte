<script lang="ts">
	import StatusPill from '$lib/components/shared/StatusPill.svelte';
	import type { IntegrationTestResponse } from '$lib/settings/types';

	interface Props {
		enabled: boolean;
		result?: IntegrationTestResponse;
		testing?: boolean;
	}

	let { enabled, result, testing = false }: Props = $props();
</script>

<div class="grid min-w-30 gap-1" aria-live="polite">
	<StatusPill tone={enabled ? 'success' : 'muted'}>{enabled ? 'Enabled' : 'Disabled'}</StatusPill>
	{#if testing}
		<StatusPill tone="pending">Testing</StatusPill>
	{:else if result}
		<StatusPill tone={result.success ? 'success' : 'error'}>
			{result.success ? 'Test OK' : 'Test failed'}
		</StatusPill>
		<span class="max-w-55 text-xs text-muted-foreground"
			>{result.message} - {result.latencyMs} ms</span
		>
	{/if}
</div>
