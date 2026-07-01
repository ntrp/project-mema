<script lang="ts">
	import SystemEventsSettings from '$lib/components/settings/SystemEventsSettings.svelte';
	import SystemLogFilesSettings from '$lib/components/settings/SystemLogFilesSettings.svelte';
	import SystemLogsSettings from '$lib/components/settings/SystemLogsSettings.svelte';
	import SystemStatusSettings from '$lib/components/settings/SystemStatusSettings.svelte';
	import PageHeading from '$lib/components/shared/PageHeading.svelte';
	import type { SystemSection } from '$lib/settings/types';

	interface Props {
		activeSection: SystemSection;
	}

	let { activeSection }: Props = $props();
	let eventsConnected = $state(false);
	let logsConnected = $state(false);

	function connectionDotClass(connected: boolean) {
		const base = 'ml-2 inline-block size-3 translate-y-[-2px] rounded-full';
		return connected
			? `${base} animate-pulse bg-primary ring-4 ring-primary/20`
			: `${base} bg-muted-foreground ring-4 ring-muted`;
	}
</script>

<section aria-labelledby="system-title">
	{#if activeSection === 'status'}
		<PageHeading eyebrow="System" title="Status" titleId="system-title" />
		<div class="space-y-4">
			<SystemStatusSettings />
		</div>
	{:else if activeSection === 'events'}
		<PageHeading eyebrow="System" title="Events" titleId="system-title">
			<span
				class={connectionDotClass(eventsConnected)}
				aria-label={eventsConnected ? 'Event stream connected' : 'Event stream reconnecting'}
			></span>
		</PageHeading>
		<div class="space-y-4">
			<SystemEventsSettings onConnectionChange={(connected) => (eventsConnected = connected)} />
		</div>
	{:else}
		<PageHeading eyebrow="System" title="Logs" titleId="system-title">
			<span
				class={connectionDotClass(logsConnected)}
				aria-label={logsConnected ? 'Log stream connected' : 'Log stream reconnecting'}
			></span>
		</PageHeading>
		<div class="space-y-4">
			<SystemLogsSettings onConnectionChange={(connected) => (logsConnected = connected)} />
			<SystemLogFilesSettings />
		</div>
	{/if}
</section>
