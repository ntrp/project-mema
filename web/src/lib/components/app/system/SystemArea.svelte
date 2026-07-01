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
	let eventsConnected = $state(false);
	let logsConnected = $state(false);
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
			<div class="page-title-row">
				<h1 id="system-title">Events</h1>
				<span
					class:connected={eventsConnected}
					class="event-connection-dot"
					aria-label={eventsConnected ? 'Event stream connected' : 'Event stream reconnecting'}
				></span>
			</div>
		</div>
		<div class="settings-stack">
			<SystemEventsSettings onConnectionChange={(connected) => (eventsConnected = connected)} />
		</div>
	{:else}
		<div class="page-heading">
			<p>System</p>
			<div class="page-title-row">
				<h1 id="system-title">Logs</h1>
				<span
					class:connected={logsConnected}
					class="event-connection-dot"
					aria-label={logsConnected ? 'Log stream connected' : 'Log stream reconnecting'}
				></span>
			</div>
		</div>
		<div class="settings-stack">
			<SystemLogsSettings onConnectionChange={(connected) => (logsConnected = connected)} />
			<SystemLogFilesSettings />
		</div>
	{/if}
</section>
