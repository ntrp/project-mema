<script lang="ts">
	import SystemEventsSettings from '$lib/components/settings/SystemEventsSettings.svelte';
	import SystemLogFilesSettings from '$lib/components/settings/SystemLogFilesSettings.svelte';
	import SystemLogsSettings from '$lib/components/settings/SystemLogsSettings.svelte';
	import SystemStatusSettings from '$lib/components/settings/SystemStatusSettings.svelte';
	import type { SystemSection } from '$lib/settings/types';

	interface Props {
		activeSection: SystemSection;
	}

	let { activeSection }: Props = $props();
</script>

<section aria-labelledby="system-title">
	{#if activeSection === 'status'}
		<div class="page-heading">
			<p>System</p>
			<h1 id="system-title">Status</h1>
		</div>
		<div class="settings-stack">
			<SystemStatusSettings />
		</div>
	{:else if activeSection === 'events'}
		<div class="page-heading">
			<p>System</p>
			<h1 id="system-title">Events</h1>
		</div>
		<div class="settings-stack">
			<SystemEventsSettings />
		</div>
	{:else if activeSection === 'log-files'}
		<div class="page-heading">
			<p>System</p>
			<h1 id="system-title">Log files</h1>
		</div>
		<div class="settings-stack">
			<SystemLogFilesSettings />
		</div>
	{:else}
		<div class="page-heading">
			<p>System</p>
			<h1 id="system-title">Live logs</h1>
		</div>
		<div class="settings-stack">
			<SystemLogsSettings />
		</div>
	{/if}
</section>
