<script lang="ts">
	/* global EventSource, MessageEvent */
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { listSystemEvents } from '$lib/settings/api';
	import { formatCompactDateTime } from '$lib/settings/dateFormat';
	import type { SystemEvent } from '$lib/settings/types';

	type StreamEnvelope<T> = { data: T };

	const maxEvents = 20;

	let open = $state(false);
	let events = $state<SystemEvent[]>([]);
	let errorTotal = $state(0);
	let loaded = $state(false);

	const visibleEvents = $derived(events.slice(0, maxEvents));

	onMount(() => {
		void load();
		const source = new EventSource('/api/events', { withCredentials: true });
		source.addEventListener('system.event.created', (event) => {
			const nextEvent = parseEvent<SystemEvent>(event);
			if (!nextEvent || (nextEvent.severity !== 'warning' && nextEvent.severity !== 'error')) {
				return;
			}
			events = [nextEvent, ...events.filter((item) => item.id !== nextEvent.id)];
			errorTotal = events.filter((item) => item.severity === 'error').length;
		});
		source.addEventListener('system.event.deleted', (event) => {
			const deleted = parseEvent<{ id: string }>(event);
			if (deleted?.id) {
				events = events.filter((item) => item.id !== deleted.id);
				errorTotal = events.filter((item) => item.severity === 'error').length;
			}
		});
		return () => source.close();
	});

	async function load() {
		try {
			events = (await listSystemEvents()).filter(
				(event) => event.severity === 'warning' || event.severity === 'error'
			);
			errorTotal = events.filter((event) => event.severity === 'error').length;
		} catch {
			events = [];
			errorTotal = 0;
		} finally {
			loaded = true;
		}
	}

	function parseEvent<T>(event: Event) {
		try {
			return (JSON.parse((event as MessageEvent<string>).data) as StreamEnvelope<T>).data;
		} catch {
			return undefined;
		}
	}
</script>

<div class="event-notifications">
	<button
		type="button"
		class="icon-button notification-button"
		aria-label="Warning and error events"
		aria-haspopup="menu"
		aria-expanded={open}
		title="Events"
		onclick={() => (open = !open)}
	>
		<span class="app-icon" aria-hidden="true">notifications</span>
		{#if errorTotal > 0}
			<span class="notification-badge">{errorTotal}</span>
		{/if}
	</button>
	{#if open}
		<div class="notification-dropdown" role="menu">
			{#if visibleEvents.length > 0}
				{#each visibleEvents as event (event.id)}
					<a
						class={`notification-item event-${event.severity}`}
						href={resolve('/system/events')}
						role="menuitem"
					>
						<span class="status-pill">{event.severity}</span>
						<strong>{event.message}</strong>
						<small>{event.category} - {formatCompactDateTime(event.createdAt)}</small>
					</a>
				{/each}
			{:else}
				<p>{loaded ? 'No warning or error events.' : 'Loading events'}</p>
			{/if}
		</div>
	{/if}
</div>
