<script lang="ts">
	import SystemLogsSettings from '$lib/components/settings/SystemLogsSettings.svelte';
	import type { SystemSection } from '$lib/settings/types';

	interface Props {
		activeSection: SystemSection;
	}

	let { activeSection }: Props = $props();
</script>

<section aria-labelledby="system-title">
	{#if activeSection === 'logs'}
		<div class="page-heading">
			<p>System</p>
			<h1 id="system-title">Logs</h1>
		</div>
		<div class="settings-stack">
			<SystemLogsSettings />
		</div>
	{/if}
</section>
